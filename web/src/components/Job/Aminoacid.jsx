import { Grid, Typography, Chip, Grow } from "@material-ui/core";
import React from "react";
import aminoacids from "./Aminoacids.js";

export default class Aminoacid extends React.Component {
  constructor(props) {
    super(props);
  }

  aaProperties(aa) {
    return aa.properties.map((p) => {
      return (
        <Grid item key={p}>
          <Chip
            variant="outlined"
            size="small"
            label={p}
            className={"propchip " + p}
          />
        </Grid>
      );
    });
  }

  render() {
    const aa = aminoacids[this.props.aa];
    if (this.props.right) {
      return (
        <div className="aa">
          <Grow in={true} key={this.props.aa}>
            <Grid
              container
              direction="row"
              spacing={1}
              alignItems="center"
              justify="flex-end"
            >
              <Grid item xs>
                <Grid container direction="column" alignItems="flex-end">
                  <Grid item>
                    <Typography variant="h4">{aa.name}</Typography>
                  </Grid>
                  <Grid item>
                    <Grid container direction="row" justify="flex-end">
                      {this.aaProperties(aa)}
                    </Grid>
                  </Grid>
                </Grid>
              </Grid>
              <Grid item>
                <img
                  src={"/assets/aa/" + this.props.aa.toLowerCase() + ".svg"}
                  alt=""
                />
              </Grid>
            </Grid>
          </Grow>
        </div>
      );
    }

    return (
      <div className="aa">
        <Grow in={true} key={this.props.aa}>
          <Grid container direction="row" spacing={1} alignItems="center">
            <Grid item>
              <img
                src={"/assets/aa/" + this.props.aa.toLowerCase() + ".svg"}
                alt=""
              />
            </Grid>
            <Grid item xs>
              <Grid container direction="column">
                <Grid item>
                  <Typography variant="h4">{aa.name}</Typography>
                </Grid>
                <Grid item>
                  <Grid container direction="row" justify="flex-start">
                    {this.aaProperties(aa)}
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Grow>
      </div>
    );
  }
}
