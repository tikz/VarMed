import {
  Grid,
  Typography,
  Chip,
  Box,
  TextField,
  Divider,
  IconButton,
  Tooltip,
  Grow,
} from "@material-ui/core";
import Autocomplete from "@material-ui/lab/Autocomplete";
import GridOn from "@material-ui/icons/GridOn";
import React from "react";

import { ResultsContext } from "./ResultsContext";
import Aminoacid from "./Aminoacid";
import Publications from "./Publications";
import "../../styles/components/variant-viewer.scss";

export default class VariantViewer extends React.Component {
  constructor(props) {
    super(props);

    this.loadVariants(this.props);
    this.state = {
      selected: this.variants[0],
    };
  }

  componentWillReceiveProps(nextProps) {
    this.loadVariants(nextProps);
    this.setState({ selected: this.variants[0] });
  }

  loadVariants(props) {
    this.variants = props.variants
      .map((v) => ({
        variant: v,
        name: v.position + " " + v.fromAa + "⟶" + v.toAa,
      }))
      .sort(function (a, b) {
        return a.variant.position - b.variant.position;
      });
  }

  setVariant(v) {
    this.focusPos(v.variant.position);
    this.setState({ selected: v });

    this.context.structure.current.loadMutated(
      this.props.pdb,
      v.variant.changeDir
    );
  }

  focusPos(pos) {
    const structure = this.context.structure.current;
    const posMap = this.context.posMap;
    const res = posMap.unpToPDB(pos);
    if (res.length > 0) {
      structure.selectFocus(res[0].chain, res[0].position, res[0].position);
    }
  }

  positionChip(pos, tag, label, optTag, optLabel) {
    const included = this.props.posFeatures[pos].includes(tag);
    if (!included) {
      tag = optTag;
      label = optLabel;
    }

    if (included || optTag) {
      return (
        <Grid item>
          <Chip
            variant="outlined"
            size="small"
            label={label}
            className={"propchip " + tag}
          />
        </Grid>
      );
    }
  }

  render() {
    const v = this.state.selected.variant;

    return (
      <Box className="variants">
        <Grid container alignItems="center" spacing={1}>
          <Grid item xs={4}>
            <Autocomplete
              disableClearable
              value={this.state.selected}
              options={this.variants}
              getOptionLabel={(v) => v.name}
              getOptionSelected={(o, v) => o.name == v.name}
              renderInput={(params) => (
                <TextField {...params} label="Variant" variant="outlined" />
              )}
              onChange={(event, newValue) => {
                this.setVariant(newValue);
              }}
            />
          </Grid>
          <Grid item xs={1}>
            <Tooltip title="Download as CSV" arrow>
              <a href={this.props.csvUrl}>
                <IconButton aria-label="collapse">
                  <GridOn />
                </IconButton>
              </a>
            </Tooltip>
          </Grid>
          <Grid item xs container direction="column" alignItems="flex-end">
            <Grid item>
              <Typography variant="overline">Predicted outcome</Typography>
            </Grid>
            <Grid item>
              <Typography
                variant="button"
                className={v.outcome.startsWith("potential") ? "yellow" : "red"}
              >
                {v.outcome}
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
            <Grow in={true} key={v.position}>
              <Grid container direction="column">
                <Grid item>
                  <a
                    onClick={() => {
                      this.focusPos(v.position);
                    }}
                  >
                    <Typography variant="h3">{v.position}</Typography>
                  </a>
                </Grid>
                <Grid item container direction="column">
                  {this.positionChip(
                    v.position,
                    "high-conservation",
                    "Highly conserved"
                  )}
                  {this.positionChip(
                    v.position,
                    "binding-site",
                    "Binding site"
                  )}
                  {this.positionChip(
                    v.position,
                    "buried",
                    "Buried",
                    "exposed",
                    "Exposed"
                  )}
                  {this.positionChip(v.position, "interface", "Interface")}
                  {this.positionChip(
                    v.position,
                    "high-aggregability",
                    "High aggregability"
                  )}
                  {this.positionChip(
                    v.position,
                    "high-switchability",
                    "High switchability"
                  )}
                </Grid>
              </Grid>
            </Grow>
          </Grid>
          <Grid item xs>
            <Aminoacid aa={v.fromAa} />
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
                  <p>
                    ΔΔG ={" "}
                    <b className={v.ddg > 2 ? "red" : ""}>{v.ddg.toFixed(1)}</b>
                  </p>
                  <p className="unit">kcal/mol</p>
                </div>
              </Grid>
            </Grid>
          </Grid>

          <Grid item xs>
            <Aminoacid aa={v.toAa} right />
          </Grid>
        </Grid>

        <Divider />
        {v.dbSNPId && (
          <Box>
            <Grow in={true} key={v.dbSNPId}>
              <Grid container justify="space-between" alignItems="center">
                <Grid item>
                  <Typography variant="h6">ClinVar</Typography>
                </Grid>
                <Grid item>
                  <a
                    href={"https://www.ncbi.nlm.nih.gov/snp/" + v.dbSNPId}
                    target="_blank"
                  >
                    <Chip
                      label={v.dbSNPId}
                      size="small"
                      variant="outlined"
                      className="chip"
                    />
                  </a>
                </Grid>
              </Grid>
            </Grow>

            {v.cvClinSig && (
              <Grow in={true} key={"cs" + v.dbSNPId}>
                <Grid
                  container
                  direction="column"
                  alignItems="center"
                  className="clin-sig"
                >
                  <Grid item>
                    <Typography variant="overline">
                      Clinical significance
                    </Typography>
                  </Grid>
                  <Grid item>
                    <Typography variant="caption">{v.cvClinSig}</Typography>
                  </Grid>
                </Grid>
              </Grow>
            )}

            {v.cvPhenotypes && (
              <Grow in={true} key={"p" + v.dbSNPId}>
                <Grid container direction="column" alignItems="center">
                  <Grid item>
                    <Typography variant="overline">Phenotypes</Typography>
                  </Grid>
                  <Grid item>
                    <Typography variant="caption">
                      {v.cvPhenotypes.split("|").join(", ")}
                    </Typography>
                  </Grid>
                </Grid>
              </Grow>
            )}

            <Divider />
          </Box>
        )}
        {v.pubmedIds && (
          <Box>
            <Grow in={true} key={v.id}>
              <Grid container justify="space-between" alignItems="center">
                <Grid item>
                  <Typography variant="h6">Publications</Typography>
                </Grid>
                <Grid item>
                  {v.id && (
                    <a
                      href={
                        "https://web.expasy.org/variant_pages/" + v.id + ".html"
                      }
                      target="_blank"
                    >
                      <Chip
                        label={v.id}
                        size="small"
                        variant="outlined"
                        className="chip"
                      />
                    </a>
                  )}
                </Grid>
              </Grid>
            </Grow>

            <Publications
              publications={this.props.publications}
              pubmeds={v.pubmedIds}
            />
            <Divider />
          </Box>
        )}
      </Box>
    );
  }
}
VariantViewer.contextType = ResultsContext;
