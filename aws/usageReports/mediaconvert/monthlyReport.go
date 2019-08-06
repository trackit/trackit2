//   Copyright 2019 MSolution.IO
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package mediaconvert

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediaconvert"
	"github.com/trackit/jsonlog"
	"gopkg.in/olivere/elastic.v5"

	taws "github.com/trackit/trackit/aws"
	"github.com/trackit/trackit/aws/usageReports"
	"github.com/trackit/trackit/config"
	"github.com/trackit/trackit/errors"
)

// getElasticSearchMediaConvertJob prepares and run the request to retrieve the a report of an instance
// It will return the data and an error.
func getElasticSearchMediaConvertJob(ctx context.Context, account, instance string, client *elastic.Client, index string) (*elastic.SearchResult, error) {
	l := jsonlog.LoggerFromContextOrDefault(ctx)
	query := elastic.NewBoolQuery()
	query = query.Filter(elastic.NewTermQuery("account", account))
	query = query.Filter(elastic.NewTermQuery("instance.id", instance))
	search := client.Search().Index(index).Size(1).Query(query)
	res, err := search.Do(ctx)
	if err != nil {
		if elastic.IsNotFound(err) {
			l.Warning("Query execution failed, ES index does not exists", map[string]interface{}{
				"index": index,
				"error": err.Error(),
			})
			return nil, errors.GetErrorMessage(ctx, err)
		} else if cast, ok := err.(*elastic.Error); ok && cast.Details.Type == "search_phase_execution_exception" {
			l.Error("Error while getting data from ES", map[string]interface{}{
				"type":  fmt.Sprintf("%T", err),
				"error": err,
			})
		} else {
			l.Error("Query execution failed", map[string]interface{}{"error": err.Error()})
		}
		return nil, errors.GetErrorMessage(ctx, err)
	}
	return res, nil
}
/*
// getJobInfoFromEs gets information about an instance from previous report to put it in the new report
func getJobInfoFromES(ctx context.Context, instance utils.CostPerResource, account string, userId int) Job {
	var docType JobReport
	var inst = Job{
		JobBase: JobBase{
			Id:         "N/A",
			Region:     "N/A",
			State:      "N/A",
			Purchasing: "N/A",
			KeyPair:    "",
			Type:       "N/A",
			Platform:   "Linux/UNIX",
		},
		Tags:  make([]utils.Tag, 0),
		Costs: make(map[string]float64, 0),
		Stats: Stats{
			Cpu: Cpu{
				Average: -1,
				Peak:    -1,
			},
			Network: Network{
				In:  -1,
				Out: -1,
			},
			Volumes: make([]Volume, 0),
		},
	}
	inst.Costs["instance"] = instance.Cost
	res, err := getElasticSearchMediaConvertJob(ctx, account, instance.Resource,
		es.Client, es.IndexNameForUserId(userId, IndexPrefixMediaConvertReport))
	if err == nil && res.Hits.TotalHits > 0 && len(res.Hits.Hits) > 0 {
		err = json.Unmarshal(*res.Hits.Hits[0].Source, &docType)
		if err == nil {
			inst.Region = docType.Job.Region
			inst.Purchasing = docType.Job.Purchasing
			inst.KeyPair = docType.Job.KeyPair
			inst.Type = docType.Job.Type
			inst.Platform = docType.Job.Platform
			inst.Tags = docType.Job.Tags
		}
	}
	return inst
}*/

// fetchMonthlyJobsList sends in instanceInfoChan the instances fetched from DescribeJobs
// and filled by DescribeJobs and getJobStats.
func fetchMonthlyJobsList(ctx context.Context, creds *credentials.Credentials,
	account, region string, instanceChan chan Job, startDate, endDate time.Time, userId int) error {
	defer close(instanceChan)
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	}))
	svc := mediaconvert.New(sess)
	listJobs, err := svc.ListJobs(&mediaconvert.ListJobsInput{})
	if err != nil {
		//instanceChan <- getJobInfoFromES(ctx, inst, account, userId)
		return err
	}
	for _, job := range listJobs.Jobs {
		instanceChan <- Job{
			ReportBase: ReportBase{
				Id: aws.StringValue(job.Id),
				Arn: aws.StringValue(job.Arn),
				BillingTagsSource: aws.StringValue(job.BillingTagsSource),
				CreatedAt: aws.TimeValue(job.CreatedAt),
				CurrentPhase: aws.StringValue(job.CurrentPhase),
			},
			Tags:    nil,
			Costs:   nil,
			Stats:   Stats{},
		}
		log.Printf("in CHANNEL =========", map[string]interface{}{
			"chan": instanceChan,
			"ID": job.Id,
		})
	}
	return nil
}

// getMediaConvertMetrics gets credentials, accounts and region to fetch MediaConvert instances stats
func fetchMonthlyJobsStats(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) ([]JobReport, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	creds, err := taws.GetTemporaryCredentials(aa, MonitorJobStsSessionName)
	if err != nil {
		logger.Error("Error when getting temporary credentials", err.Error())
		return nil, err
	}
	defaultSession := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(config.AwsRegion),
	}))
	account, err := utils.GetAccountId(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when getting account id", err.Error())
		return nil, err
	}
	regions, err := utils.FetchRegionsList(ctx, defaultSession)
	if err != nil {
		logger.Error("Error when fetching regions list", err.Error())
		return nil, err
	}
	jobChans := make([]<-chan Job, 0, len(regions))
		for _, region := range regions {
				jobChan := make(chan Job)
				go fetchMonthlyJobsList(ctx, creds, account, region, jobChan, startDate, endDate, aa.UserId)
				jobChans = append(jobChans, jobChan)
		}
	instancesList := make([]JobReport, 0)
	for instance := range merge(jobChans...) {
		instancesList = append(instancesList, JobReport{
			ReportBase: utils.ReportBase{
				Account:    account,
				ReportDate: startDate,
				ReportType: "monthly",
			},
			Job: instance,
		})
	}
	return instancesList, nil
}

// PutMediaConvertMonthlyReport puts a monthly report of MediaConvert instance in ES
func PutMediaConvertMonthlyReport(ctx context.Context, aa taws.AwsAccount, startDate, endDate time.Time) (bool, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Starting MediaConvert monthly report", map[string]interface{}{
		"awsAccountId": aa.Id,
		"startDate":    startDate.Format("2006-01-02T15:04:05Z"),
		"endDate":      endDate.Format("2006-01-02T15:04:05Z"),
	})
	already, err := utils.CheckMonthlyReportExists(ctx, startDate, aa, IndexPrefixMediaConvertReport)
	if err != nil {
		return false, err
	} else if already {
		logger.Info("There is already an MediaConvert monthly report", nil)
		return false, nil
	}
	instances, err := fetchMonthlyJobsStats(ctx, aa, startDate, endDate)
	if err != nil {
		return false, err
	}
	err = importJobsToEs(ctx, aa, instances)
	if err != nil {
		return false, err
	}
	return true, nil
}
