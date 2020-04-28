import { Box, Toolbar } from "@material-ui/core";
import axios from "axios";
import React from "react";
import NavBar from "../NavBar";
import Results from "./Results";
import StatusConsole from "./StatusConsole";

export default class Job extends React.Component {
  constructor(props) {
    super(props);

    this.jobID = this.props.match.params.id;
    this.state = { results: {}, error: 0 };
    this.loadResults = this.loadResults.bind(this);
    this.loadResults();
  }

  loadResults() {
    let that = this;
    axios
      .get(API_URL + "/api/job/" + this.jobID)
      .then(function (response) {
        that.setState({ results: response.data, error: 0 });
      })
      .catch(function (error) {
        that.setState({ error: error });
      });
  }

  render() {
    if (this.state.error != 0) {
      return (
        <Box>
          <NavBar />
          <Toolbar />
          <h3>{this.state.error.message}</h3>
        </Box>
      );
    }

    return (
      <Box>
        <NavBar />
        <Toolbar />
        {(this.state.results.Status == 0 || this.state.results.Status == 1) && (
          <StatusConsole jobID={this.jobID} reload={this.loadResults} />
        )}
        {(this.state.results.Status == 2 || this.state.results.Status == 3) && (
          <Results
            results={this.state.results}
            jobID={this.jobID}
            pdbID={this.state.results.Request.pdbs[0]}
          />
        )}
      </Box>
    );
  }
}
