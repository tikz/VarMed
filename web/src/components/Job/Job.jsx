import { Box, Toolbar } from "@material-ui/core";
import axios from 'axios';
import React from 'react';
import NavBar from '../NavBar';
import Results from './Results';
import StatusConsole from './StatusConsole';

export default class Job extends React.Component {
    constructor(props) {
        super(props);
        this.jobID = this.props.match.params.id;

        this.state = { results: {} }

        this.loadResults = this.loadResults.bind(this)
        this.loadResults()
    }

    loadResults() {
        let that = this;
        axios.get('/api/job/' + this.jobID)
            .then(function (response) {
                that.setState({ results: response.data })
            })
    }

    render() {
        return (
            <Box>
                <NavBar />
                <Toolbar />
                {(this.state.results.Status == 1) &&
                    <StatusConsole jobID={this.jobID} reload={this.loadResults} />}
                {(this.state.results.Status == 2) &&
                    <Results results={this.state.results} jobID={this.jobID} pdbID={this.state.results.Request.pdbs[0]} />}
            </Box>
        )
    }
}