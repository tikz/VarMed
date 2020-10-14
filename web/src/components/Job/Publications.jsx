import {
  Grid,
  Typography,
  Chip,
  Box,
  IconButton,
  Grow,
} from "@material-ui/core";
import ExpandLessIcon from "@material-ui/icons/ExpandLess";
import ExpandMoreIcon from "@material-ui/icons/ExpandMore";
import React from "react";
import PublicationItem from "./PublicationItem";

export default class Evidence extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <Box>
        <Grid container direction="column" alignItems="flex-start">
          {this.props.pubmeds.map((id) => {
            const pub = this.props.publications[id];
            return (
              <Grow in={true} key={id}>
                <Grid item>
                  <PublicationItem
                    title={pub.title}
                    authors={pub.authors}
                    journal={pub.journal}
                    doi={pub.doi}
                    pubmed={pub.pubmed}
                  />
                </Grid>
              </Grow>
            );
          })}

          {/* <Grid item>
            <IconButton aria-label="collapse">
              <ExpandMoreIcon />
            </IconButton>
          </Grid> */}
        </Grid>
      </Box>
    );
  }
}
