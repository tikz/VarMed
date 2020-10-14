import { Grid } from "@material-ui/core";
import React from "react";

export class QueueInfo extends React.Component {
  render() {
    return (
      <Grid container spacing={2}>
        <Grid container item>
          <Grid item xs={8}>
            Jobs in queue:
          </Grid>
          <Grid item>-</Grid>
        </Grid>
        <Grid container item>
          <Grid item xs={8}>
            Estimated time:
          </Grid>
          <Grid item>-</Grid>
        </Grid>
      </Grid>
    );
  }
}
