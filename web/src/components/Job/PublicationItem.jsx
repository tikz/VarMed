import { Grid, Typography, Chip } from "@material-ui/core";
import React from "react";

export default class PublicationItem extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    let title = "";
    if (this.props.pubmed) {
      title = (
        <a
          href={"https://pubmed.ncbi.nlm.nih.gov/" + this.props.pubmed}
          target="_blank"
        >
          <Typography>{this.props.title}</Typography>
        </a>
      );
    } else if (this.props.doi) {
      title = (
        <a href={"https://doi.org/" + this.props.doi} target="_blank">
          <Typography>{this.props.title}</Typography>
        </a>
      );
    } else {
      title = <Typography>{this.props.title}</Typography>;
    }
    return (
      <Grid container direction="column" className="evidence">
        <Grid item>
          {title}
          <Typography className="authors">{this.props.authors}</Typography>
          <Typography className="journal">{this.props.journal}</Typography>
        </Grid>
        <Grid item></Grid>
      </Grid>
    );
  }
}
