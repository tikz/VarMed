import { Grid, Typography, Chip, Box, IconButton } from "@material-ui/core";
import ExpandLessIcon from "@material-ui/icons/ExpandLess";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import React from "react";
import EvidenceItem from "./EvidenceItem";

export default class Evidence extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <Box>
        <Grid container direction="column" alignItems="center">
          <Grid item>
            <EvidenceItem />
          </Grid>
          <Grid item>
            <IconButton aria-label="collapse">
              <ExpandMoreIcon />
            </IconButton>
          </Grid>
        </Grid>
      </Box>
    );
  }
}
