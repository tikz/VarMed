import { Container } from '@material-ui/core';
import React from 'react';
import { Link } from "react-router-dom";
import NewJob from './NewJob';

export default class Index extends React.Component {
    render() {
        return (
            <Container>
                <h1>Index</h1>
                <p>views:</p>
                <Link to="/results">results</Link>
                <br />
                <Link to="/new-job">new job</Link>
            </Container>
        )
    }
}