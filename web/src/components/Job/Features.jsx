import { Divider, Grid } from "@material-ui/core";
import React from "react";
import "../../styles/components/features.scss";
import { ChipRes, ChipHet } from "./FeatureChips";
import { SurfaceSwitch } from "./SurfaceSwitch";
import { ResultsContext } from "./ResultsContext";

export class Features extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    const res = this.context.results;
    const chip = function (label, residues) {
      return residues !== null ? (
        <ChipRes label={label} key={label} residues={residues} />
      ) : null;
    };

    const resRange = function (chain, start, end) {
      return {
        chain: chain,
        position: start,
        positionEnd: end,
      };
    };

    const unpChains = [];
    Object.keys(res.pdb.SIFTS.UniProt).forEach((unpID) => {
      const unp = res.pdb.SIFTS.UniProt[unpID];
      const chains = unp.mappings.map((chain) => {
        return resRange(chain.chain_id, 1, chain.end.residue_number);
      });
      unpChains.push(
        [chip(unpID + " " + unp.identifier.replace("_", " "), chains)].concat(
          chains.map((chain) => {
            return chip("Chain " + chain.chain, [chain]);
          })
        )
      );
    });

    const fams = Object.values(res.pdb.SIFTS.Pfam).map((fam) => {
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

    const hets = res.pdb.hetGroups;

    const interaction = chip("Interface", res.interaction.residues);
    const buried = chip("Buried", res.exposure.residues);
    const binding = chip("Binding", res.binding.residues);

    return (
      <Grid container className="features">
        <Grid item xs>
          <Grid container className="chips" wrap="wrap">
            {fams}
            <Divider orientation="vertical" flexItem />
            {unpChains.map((unp, index) => {
              return unp.concat(
                <Divider key={index} orientation="vertical" flexItem />
              );
            })}
            <Divider orientation="vertical" flexItem />
            {hets &&
              hets.map((hetId) => {
                if (hetId != "HOH") {
                  return <ChipHet label={hetId} key={hetId} hetID={hetId} />;
                }
              })}
            <Divider orientation="vertical" flexItem />
            {interaction}
            {buried}
            {binding}
            {/* {res.binding.pockets !== null &&
              res.binding.pockets.map((pocket, index) => {
                return chip("Pocket " + index, pocket.residues);
              })} */}
            {Object.keys(res.binding.ligands).map((ligand) => {
              return chip("Near " + ligand, res.binding.ligands[ligand]);
            })}
          </Grid>
        </Grid>
        <Grid item xs={1}>
          <SurfaceSwitch />
        </Grid>
      </Grid>
    );
  }
}
Features.contextType = ResultsContext;
