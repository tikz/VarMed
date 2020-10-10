import { Grid, Typography, Chip } from "@material-ui/core";
import React from "react";

export default class EvidenceItem extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <Grid container direction="column" className="evidence">
        <Grid item>
          <Typography>
            {/* {this.props.title} */}
            An atypical variant of Fabry's disease in men with left ventricular
            RT hypertrophy.
          </Typography>
          <Typography className="authors">
            {/* {this.props.authors} */}
            Nakao S., Takenaka T., Maeda M., Kodama C., Tanaka A., Tahara M., RA
            Yoshida A., Kuriyama M., Hayashibe H., Sakuraba H., Tanaka H.
          </Typography>
          <Typography className="journal">
            {/* {this.props.journal} */}
            RL N. Engl. J. Med. 333:288-293(1995)
          </Typography>
        </Grid>
        <Grid item></Grid>
      </Grid>
    );
  }
}
