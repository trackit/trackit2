package tagginglambda

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/trackit/jsonlog"

	"github.com/trackit/trackit/tagging/utils"
)

type function struct {
	Name   string              `json:"name"`
	Region string              `json:"region"`
	Tags   []utils.TagDocument `json:"tags"`
}

type source struct {
	Function function `json:"function"`
}

const sourceIndexName = "lambda-reports"
const urlFormat = "https://%s.console.aws.amazon.com/lambda/home?region=%s#/functions/%s"

// Process generates tagging reports from Lambda reports
func Process(ctx context.Context, account int, awsAccount string) ([]utils.TaggingReportDocument, error) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	logger.Info("Processing Lambda reports.", nil)

	hits, err := fetchReports(ctx, account)
	if err != nil {
		return nil, err
	}

	var documents []utils.TaggingReportDocument
	for _, hit := range hits {
		document, success := processHit(ctx, hit, awsAccount)
		if success {
			documents = append(documents, document)
		}
	}

	logger.Info(fmt.Sprintf("%d Lambda reports processed.", len(documents)), nil)
	return documents, nil
}

// processHit converts an elasticSearch hit into a TaggingReportDocument
// Second argument is true if operation is a success
func processHit(ctx context.Context, hit *elastic.SearchHit, awsAccount string) (utils.TaggingReportDocument, bool) {
	logger := jsonlog.LoggerFromContextOrDefault(ctx)
	var source source
	err := json.Unmarshal(*hit.Source, &source)
	if err != nil {
		logger.Error("Could not process a Lambda report.", nil)
		return utils.TaggingReportDocument{}, false
	}

	regionForURL := utils.GetRegionForURL(source.Function.Region)
	fmt.Printf("https://%s.console.aws.amazon.com/lambda/home?region=%s#/functions/%s\n", regionForURL, regionForURL, source.Function.Name)

	document := utils.TaggingReportDocument{
		Account:      awsAccount,
		ResourceID:   source.Function.Name,
		ResourceType: "lambda",
		Region:       source.Function.Region,
		URL:          fmt.Sprintf(urlFormat, regionForURL, regionForURL, source.Function.Name),
		Tags:         source.Function.Tags,
	}
	return document, true
}
