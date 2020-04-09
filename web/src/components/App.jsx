import { Box, createMuiTheme, ThemeProvider, Toolbar } from '@material-ui/core';
import React from 'react';
import { Route, Switch } from "react-router-dom";
import Index from './Index';
import NavBar from './NavBar';
import Results from './Results';

const darkTheme = createMuiTheme({
    palette: {
        type: 'dark',
    },
});

export default class App extends React.Component {
    render() {
        return (
            <ThemeProvider theme={darkTheme} >
                <NavBar />
                <Toolbar />
                <Box>
                    <Switch>
                        <Route exact path="/" component={Index} />
                        <Route path="/results" component={Results} />
                        <Route path="*">
                            <h1>Not found</h1>
                        </Route>
                    </Switch>
                </Box>
            </ThemeProvider>
        )
    }
}