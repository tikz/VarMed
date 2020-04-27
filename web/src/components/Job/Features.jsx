import { Box, Grid, Switch, FormControlLabel } from "@material-ui/core";
import React from "react";
import FeatureChip from "./FeatureChip";
import "../../styles/components/features.scss";
import { SurfaceSwitch } from "./SurfaceSwitch";

export class Features extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    const props = this.props;
    const res = props.res;
    const chip = function (label, residues) {
      return residues !== null ? (
        <FeatureChip
          label={label}
          key={label}
          highlightResidues={props.highlightResidues}
          clearHighlight={props.clearHighlight}
          residues={residues}
        />
      ) : null;
    };

    const interaction = chip("Interface", res.Interaction.Residues);
    const buried = chip("Buried", res.Exposure.Residues);
    const catalytic = chip("Catalytic", res.Binding.Catalytic.Residues);

    return (
      <Grid container className="features">
        <Grid item xs className="chips">
          {/* {Object.keys(res.PDB.SIFTS.UniProt).map((unpID) => {
            return chip(
              unpID,
              res.PDB.SIFTS.UniProt[unpID].mappings.map((chain) => {
                return {
                  Chain: chain.chain_id,
                  Position: 1,
                  PositionEnd: res.PDB.ChainEndResNumber[chain.chain_id],
                };
              })
            );
          })} */}
          {interaction}
          {buried}
          {catalytic}
          {res.Binding.Pockets !== null &&
            res.Binding.Pockets.map((pocket) => {
              return chip("Pocket", pocket.Residues);
            })}
          {Object.keys(res.Binding.Ligands).map((ligand) => {
            return chip("Near " + ligand, res.Binding.Ligands[ligand]);
          })}
        </Grid>
        <Grid item>
          <SurfaceSwitch showSurface={this.props.showSurface} />
        </Grid>
      </Grid>
    );
  }
}
