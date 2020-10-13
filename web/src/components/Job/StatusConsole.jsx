import {
  Dialog,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Grow,
  LinearProgress,
} from "@material-ui/core";
import React from "react";
import "../../styles/components/status-console.scss";

export default class StatusConsole extends React.Component {
  constructor(props) {
    super(props);
    var url;
    if (API_URL == "") {
      let protocol = window.location.protocol === "https:" ? "wss://" : "ws://";
      url = protocol + window.location.host + "/ws/job/" + this.props.jobId;
    } else {
      url = API_URL.replace("http", "ws") + "/ws/job/" + this.props.jobId;
    }
    this.ws = new WebSocket(url);

    this.state = { messages: [], connected: false, error: false };
  }

  componentDidMount() {
    this.ws.onopen = () => {
      this.setState({
        connected: true,
        messages: this.state.messages.concat("Connected."),
      });
    };
    this.ws.onmessage = (evt) => {
      if (evt.data.startsWith("ERROR")) {
        this.setState({ error: true });
      }
      this.setState({
        messages: this.state.messages.concat(evt.data),
      });
    };
    this.ws.onclose = () => {
      this.props.reload();
    };
    this.ws.onerror = () => {
      this.props.reload();
    };
  }

  render() {
    if (!this.state.connected) {
      return (
        <Dialog open={true} maxWidth="sm" fullWidth={true}>
          <DialogTitle>Connecting...</DialogTitle>
          <LinearProgress variant="query" />
        </Dialog>
      );
    }
    return (
      <Dialog open={true} maxWidth="sm" fullWidth={true}>
        <DialogTitle>Your job is pending</DialogTitle>
        <DialogContent>
          <DialogContentText>This is the real time status:</DialogContentText>
          <div className="status-console">
            <div className="contents">
              {this.state.messages.map((msg, index) => {
                return (
                  <Grow in={true} key={index}>
                    <div>{msg}</div>
                  </Grow>
                );
              })}
            </div>
          </div>
        </DialogContent>
        {!this.state.error && <LinearProgress variant="query" />}
      </Dialog>
    );
  }
}
