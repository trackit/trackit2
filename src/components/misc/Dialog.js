import React, {Component} from 'react';
import Dialog from '@material-ui/core/Dialog';
import DialogTitle from '@material-ui/core/DialogTitle';
import DialogContent from '@material-ui/core/DialogContent';
import DialogActions from '@material-ui/core/DialogActions';
import PropTypes from "prop-types";

class DialogComponent extends Component {

  constructor(props) {
    super(props);
    this.state = {
      open: false
    };
    this.openDialog = this.openDialog.bind(this);
    this.closeDialog = this.closeDialog.bind(this);
    this.executeAction = this.executeAction.bind(this);
  }

  openDialog = (e) => {
    e.preventDefault();
    if (this.props.onOpen)
      this.props.onOpen();
    this.setState({open: true});
  };

  closeDialog = (e) => {
    e.preventDefault();
    if (this.props.onClose)
      this.props.onClose();
    this.setState({open: false});
  };

  executeAction = (e) => {
    e.preventDefault();
    this.closeDialog(e);
    this.props.actionFunction();
  };

  render() {
    const title = (this.props.title ? (
      <DialogTitle disableTypography>
        <h1 className="dialog-title">
          {this.props.title}
          <div className="dialog-title-children">
            {this.props.titleChildren}
          </div>
        </h1>
        <hr/>
      </DialogTitle>
    ) : null);

    const content = (this.props.children ? (
      <DialogContent>
        {this.props.children}
      </DialogContent>
    ) : null);

    const mainAction = (this.props.actionName && this.props.actionFunction ? (
      <button className="btn btn-default" onClick={this.executeAction}>
        {this.props.actionName}
      </button>
    ) : null);

    return(
      <div>

        <button className={"btn btn-" + this.props.buttonType} onClick={this.openDialog} disabled={this.props.disabled}>
          {this.props.buttonName}
        </button>

        <Dialog open={this.state.open} fullWidth>

          {title}

          {content}

          <DialogActions>

            <button className="btn btn-default" onClick={this.closeDialog}>
              {this.props.secondActionName}
            </button>

            {mainAction}

          </DialogActions>

        </Dialog>

      </div>
    );
  }

}

DialogComponent.propTypes = {
  buttonType: PropTypes.string,
  buttonName: PropTypes.node.isRequired,
  title: PropTypes.string,
  titleChildren: PropTypes.node,
  actionName: PropTypes.string,
  actionFunction: PropTypes.func,
  secondActionName: PropTypes.string,
  children: PropTypes.node,
  onOpen: PropTypes.func,
  onClose: PropTypes.func,
  disabled: PropTypes.bool
};

DialogComponent.defaultProps = {
  buttonType: "default",
  secondActionName: "Cancel",
  titleChildren: null,
  disabled: false
};

export default DialogComponent;
