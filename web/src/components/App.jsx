import { Box, createMuiTheme, ThemeProvider, Toolbar } from '@material-ui/core';
import React from 'react';
import { Route, Switch } from "react-router-dom";
import Index from './Index';
import NavBar from './NavBar';
import NewJob from './NewJob/NewJob';
import Results from './Results';

const theme = createMuiTheme({
    palette: {
        type: 'dark',
        primary: {
            main: '#2196F3',
            contrastText: '#fff',
        },
        secondary: {
            main: '#21CBF3',
        },
    },
});

export default class App extends React.Component {
    render() {
        return (
            <ThemeProvider theme={theme} >
                <Switch>
                    <Route exact path="/" component={Index} />
                    <Route path="/new-job" component={NewJob} />
                    <Route path="/results" component={Results} />
                    <Route path="*">
                        <h1>Not found</h1>
                    </Route>
                </Switch>
                {/* <NavBar />
                <Toolbar />
                <Box>
                    
                </Box> */}
            </ThemeProvider>
        )
    }
}