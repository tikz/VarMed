import { Grid, TextField, Button } from "@material-ui/core";
import React from "react";
import { QueueInfo } from "./QueueInfo";

export default class SendBar extends React.Component {
  constructor(props) {
    super(props);
    this.state = { email: "" };
    this.handleChange = this.handleChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }
  handleChange(e) {
    this.setState({ email: e.target.value });
  }
  handleSubmit() {
    this.props.submit(this.state.email);
  }
  render() {
    return (
      <Grid container spacing={2} alignItems="center">
        <Grid item xs>
          <TextField
            id="name"
            label="Email address (optional)"
            onChange={this.handleChange}
            margin="dense"
            type="email"
            value={this.state.email}
            fullWidth
          />
        </Grid>
        <Grid item xs>
          <QueueInfo />
        </Grid>
        <Grid item xs={2}>
          <Button className="glowButton" onClick={this.handleSubmit}>
            Send Job
          </Button>
        </Grid>
      </Grid>
    );
  }
}
