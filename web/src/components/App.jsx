import { createMuiTheme, ThemeProvider } from "@material-ui/core";
import React from "react";
import { Route, Switch } from "react-router-dom";
import Index from "./Index";
import NewJob from "./NewJob/NewJob";
import Job from "./Job/Job";

const theme = createMuiTheme({
  palette: {
    type: "dark",
    primary: {
      main: "#1aacdb",
      contrastText: "#fff",
    },
    secondary: {
      main: "#2196F3",
    },
  },
});

export default class App extends React.Component {
  render() {
    return (
      <ThemeProvider theme={theme}>
        <Switch>
          <Route exact path="/" component={Index} />
          <Route path="/new-job" component={NewJob} />
          <Route path="/job/:id" component={Job} />
          <Route path="*">
            <h1>Not found</h1>
          </Route>
        </Switch>
      </ThemeProvider>
    );
  }
}
