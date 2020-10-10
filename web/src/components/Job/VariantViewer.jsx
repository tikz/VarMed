import {
  Grid,
  Typography,
  Chip,
  Box,
  TextField,
  Divider,
} from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import React from "react";

import Aminoacid from "./Aminoacid";
import Evidence from "./Evidence";
import "../../styles/components/variant.scss";
import { ResultsContext } from "./ResultsContext";

export default class VariantViewer extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <Box>
        <Grid container>
          <Grid item xs={4}>
            <Autocomplete
              options={["89 M⟶K", "129 F⟶W"]}
              renderInput={(params) => (
                <TextField {...params} label="Variant" variant="outlined" />
              )}
            />
          </Grid>
          <Grid item xs container direction="column" alignItems="flex-end">
            <Grid item>
              <Typography variant="overline">Predicted outcome</Typography>
            </Grid>
            <Grid item>
              <Typography variant="button" className="orange">
                Potentially disrupts protein function
              </Typography>
            </Grid>
          </Grid>
        </Grid>

        <Divider />
        <Grid
          container
          className="substitution"
          alignItems="center"
          spacing={2}
        >
          <Grid item>
            <Grid container direction="column">
              <Grid item>
                <Typography variant="h3">666</Typography>
              </Grid>
              <Grid item container direction="column">
                <Grid item>
                  <Chip
                    variant="outlined"
                    size="small"
                    label="High conservation"
                    className="propchip high-conservation"
                  />
                </Grid>
                <Grid item>
                  <Chip
                    variant="outlined"
                    size="small"
                    label="Buried"
                    className="propchip buried"
                  />
                </Grid>

                <Grid item>
                  <Chip
                    variant="outlined"
                    size="small"
                    label="Interface"
                    className="propchip interface"
                  />
                </Grid>
                <Grid item>
                  <Chip
                    variant="outlined"
                    size="small"
                    label="Disulfide"
                    className="propchip disulfide"
                  />
                </Grid>
                <Grid item>
                  <Chip
                    variant="outlined"
                    size="small"
                    label="High switchability"
                    className="propchip high-switchability"
                  />
                </Grid>
              </Grid>
            </Grid>
          </Grid>
          <Grid item xs>
            <Aminoacid aa="H" />
          </Grid>

          <Grid item>
            <Grid container direction="column" alignItems="center">
              <Grid item>
                <Typography variant="h4" className="arrow">
                  ⟶
                </Typography>
              </Grid>
              <Grid item>
                <div className="ddg">
                  <p>ΔΔG = 5.32</p>
                  <p className="unit">kcal/mol</p>
                </div>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs>
            <Aminoacid aa="E" right />
          </Grid>
        </Grid>
        <Divider />
        <Typography variant="h6">Evidence</Typography>
        <Evidence />

        <Divider />
      </Box>
    );
  }
}
VariantViewer.contextType = ResultsContext;
