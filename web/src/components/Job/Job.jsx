import { Box, Snackbar, Toolbar } from "@material-ui/core";
import MuiAlert from "@material-ui/lab/Alert";
import axios from "axios";
import React from "react";
import Results from "./Results";
import StatusConsole from "./StatusConsole";
import NavBar from "../NavBar";

function Alert(props) {
  return <MuiAlert elevation={11} variant="filled" {...props} />;
}

export default class Job extends React.Component {
  constructor(props) {
    super(props);

    this.state = { results: {}, error: 0 };
    this.loadResults = this.loadResults.bind(this);
    this.handleAddedClose = this.handleAddedClose.bind(this);
  }

  componentDidMount() {
    this.loadResults(this.props);
  }

  componentWillReceiveProps(nextProps) {
    this.setState({ results: {} });
    this.loadResults(nextProps);
  }

  loadResults(props) {
    props = props || this.props;
    const that = this;
    axios
      .get(API_URL + "/api/job/" + props.match.params.id)
      .then(function (response) {
        that.setState({ results: response.data, error: 0 }, () => {
          that.addToMyJobs();
        });
      })
      .catch(function (error) {
        that.setState({ error: error });
      });
  }

  addToMyJobs() {
    const jobID = this.props.match.params.id;
    let jobs = JSON.parse(window.localStorage.getItem("jobs"));
    if (jobs === null) {
      jobs = [];
    }
    if (
      jobs.filter((j) => {
        return j.id == jobID;
      }).length > 0
    ) {
      return;
    }

    const req = this.state.results.Request;
    jobs.unshift({
      id: jobID,
      name: req.name,
      pdbs: req.pdbs,
      date: Date.now(),
    });

    this.setState({ added: true });
    window.localStorage.removeItem("jobs");
    window.localStorage.setItem("jobs", JSON.stringify(jobs));
  }

  handleAddedClose() {
    this.setState({ added: false });
  }

  render() {
    const jobID = this.props.match.params.id;

    if (this.state.error != 0) {
      return <h3>{this.state.error.message}</h3>;
    }

    return (
      <Box>
        <NavBar />
        <Toolbar />
        {(this.state.results.Status == 0 || this.state.results.Status == 1) && (
          <StatusConsole jobID={jobID} reload={this.loadResults} />
        )}
        {(this.state.results.Status == 2 || this.state.results.Status == 3) && (
          <Results jobID={jobID} jobResults={this.state.results} />
        )}
        <Snackbar
          open={this.state.added}
          autoHideDuration={3000}
          onClose={this.handleAddedClose}
          anchorOrigin={{ vertical: "top", horizontal: "center" }}
        >
          <Alert onClose={this.handleAddedClose} severity="success">
            Added to My Jobs
          </Alert>
        </Snackbar>
      </Box>
    );
  }
}
