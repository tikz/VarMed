import { Divider, Grid } from "@material-ui/core";
import React from "react";
import "../../styles/components/features.scss";
import ChipRes from "./ChipRes";
import ChipHet from "./ChipHet";
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
        <ChipRes
          label={label}
          key={label}
          highlightResidues={props.highlightResidues}
          clearHighlight={props.clearHighlight}
          residues={residues}
        />
      ) : null;
    };

    const resRange = function (chain, start, end) {
      return {
        Chain: chain,
        Position: start,
        PositionEnd: end,
      };
    };

    const unpChains = [];
    Object.keys(res.PDB.SIFTS.UniProt).forEach((unpID) => {
      const unp = res.PDB.SIFTS.UniProt[unpID];
      const chains = unp.mappings.map((chain) => {
        return resRange(chain.chain_id, 1, chain.end.residue_number);
      });
      unpChains.push(
        [chip(unpID + " " + unp.identifier.replace("_", " "), chains)].concat(
          chains.map((chain) => {
            return chip("Chain " + chain.Chain, [chain]);
          })
        )
      );
    });

    const fams = Object.values(res.PDB.SIFTS.Pfam).map((fam) => {
      return chip(
        fam.Identifier,
        fam.Mappings.map((chain) => {
          return resRange(
            chain.chain_id,
            chain.start.residue_number,
            chain.end.residue_number
          );
        })
      );
    });

    const hets = res.PDB.HetGroups.map((hetID) => {
      if (hetID != "HOH") {
        return (
          <ChipHet
            label={hetID}
            key={hetID}
            highlightHet={props.highlightHet}
            clearHighlight={props.clearHighlight}
            hetID={hetID}
          />
        );
      }
    });

    const interaction = chip("Interface", res.Interaction.Residues);
    const buried = chip("Buried", res.Exposure.Residues);
    const catalytic = chip("Catalytic", res.Binding.Catalytic.Residues);

    return (
      <Grid container className="features">
        <Grid container xs className="chips" xs={11} wrap="wrap">
          {fams}
          <Divider orientation="vertical" flexItem />
          {unpChains.map((unp, index) => {
            return unp.concat(
              <Divider key={index} orientation="vertical" flexItem />
            );
          })}
          <Divider orientation="vertical" flexItem />
          {hets}
          <Divider orientation="vertical" flexItem />
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
        <Grid item xs={1}>
          <SurfaceSwitch showSurface={this.props.showSurface} />
        </Grid>
      </Grid>
    );
  }
}
