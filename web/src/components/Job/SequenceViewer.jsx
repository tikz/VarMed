import React from "react";
import "../../styles/components/sequence-viewer.scss";
import { ResultsContext } from "./ResultsContext";

const FeatureViewer = require("feature-viewer");

export default class SequenceViewer extends React.Component {
  constructor(props) {
    super(props);
    this.state = { pos: 0 };

    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseLeave = this.handleMouseLeave.bind(this);
  }

  handleMouseMove() {
    var posText = this.state.zoomPositionElement.innerText;
    if (posText != this.state.pos) {
      var pos = parseInt(posText.slice(0, posText.length - 1));
      this.context.structure.current.highlightResidues(
        this.context.posMap.unpToPDB(pos)
      );
      this.setState(() => ({ pos: posText }));
    }
  }

  handleMouseLeave() {
    this.context.structure.current.clearHighlight();
  }

  load() {
    let div = document.getElementById("fv");
    if (div !== null) {
      div.innerHTML = "";
    }

    const structure = this.context.structure.current;
    const res = this.context.results;
    const posMap = this.context.posMap;

    this.fv = new FeatureViewer(res.uniprot.sequence, "#fv", {
      showAxis: true,
      showSequence: true,
      brushActive: true,
      toolbar: true,
      bubbleHelp: false,
      zoomMax: 20,
    });

    this.fv.onFeatureSelected(function (d) {
      if (d.detail.description && d.detail.description.includes("Chain")) {
        const chain = d.detail.description.split(" ")[1];
        structure.selectFocus(
          chain,
          posMap.unpToPDB(d.detail.start).filter((c) => c.chain == chain)[0]
            .position,
          posMap.unpToPDB(d.detail.end).filter((c) => c.chain == chain)[0]
            .position
        );
      } else {
        const chain = posMap.chains[0].id;
        structure.selectFocus(
          chain,
          posMap.unpToPDB(d.detail.start)[0].position, // d.detail.start + posMap.pdbOffsets[chain],
          posMap.unpToPDB(d.detail.end)[0].position // d.detail.end + posMap.pdbOffsets[chain]
        );
      }
    });

    this.setState(() => ({
      zoomPositionElement: document.getElementById("zoomPosition"),
    }));

    if (res.variants) {
      const vars = [];
      res.variants.forEach((v) => {
        vars.push({
          x: v.position,
          y: v.position,
          description:
            v.fromAa +
            " → " +
            v.toAa +
            "<br>" +
            "ΔΔG: " +
            v.ddg.toFixed(2) +
            " kcal/mol<br>" +
            v.cvClinSig,
        });
      });
      this.fv.addFeature({
        data: vars,
        name: "Variant",
        color: "#00ffa6",
        type: "rect",
        filter: "type2",
        className: "var",
      });
    }

    var glycos = res.uniprot.ptms.glycosilationSites;
    if (glycos) {
      this.fv.addFeature({
        data: glycos.map((g) => {
          return {
            x: g.position,
            y: g.position,
          };
        }),
        name: "Glycosilation",
        color: "#d1973f",
        type: "rect",
        filter: "type2",
        className: "glyco",
      });
    }

    var disulfides = res.uniprot.ptms.disulfideBonds;
    if (disulfides) {
      this.fv.addFeature({
        data: disulfides.map((d) => {
          return {
            x: d.positions[0],
            y: d.positions[1],
          };
        }),
        name: "Disulfide",
        color: "#B4AF91",
        type: "path",
        className: "disulf",
      });
    }

    res.conservation.families.forEach((fam) => {
      this.fv.addFeature({
        data: [
          {
            x: fam.start,
            y: fam.end,
            description: fam.id + " " + fam.desc,
          },
        ],
        name: "Domain",
        color: "#1aacdb",
        type: "rect",
        className: "fam" + fam.id,
      });

      var consData = [];
      fam.positions.forEach((p) => {
        consData.push({
          x: p.position,
          y: p.bitscore,
        });
      });

      this.fv.addFeature({
        data: consData,
        name: "Conservation",
        color: "#008B8D",
        type: "line",
        filter: "type2",
        height: "5",
        className: "cons",
      });
    });

    posMap.chains.forEach((chain) => {
      let name = "Chain " + chain.id;
      this.fv.addFeature({
        data: posMap.chainRanges[chain.id].map(function (r) {
          return { x: r.start, y: r.end, description: name };
        }),
        name: name,
        color: "#1aacdb",
        type: "rect",
        className: "chain" + chain.id,
      });

      let markResidues = function (that, residues, title) {
        residues = residues
          .filter((r) => r.chain == chain.id)
          .map((r) => {
            let pm = posMap.pdbToUnp(r.chain, r.position);
            return { x: pm, y: pm, description: name };
          })
          .filter((r) => r.x);
        if (residues.length > 0) {
          that.fv.addFeature({
            data: residues,
            name: chain.id + " - " + title,
            color: "#1aacdb",
            type: "rect",
            className: chain.id + title.split(" ").slice(-1)[0],
          });
        }
      };

      if (res.interaction.residues) {
        markResidues(
          this,
          res.interaction.residues.map((r) => r.residue),
          "Interface"
        );
      }

      // if (res.exposure.residues !== null) {
      //   markResidues(
      //     this,
      //     res.exposure.residues.map((r) => r.residue),
      //     "Buried"
      //   );
      // }

      // if (res.binding.residues !== null) {
      //   markResidues(this, res.binding.residues, "Sites");
      // }

      // if (Object.keys(res.binding.ligands).length != 0) {
      //   Object.keys(res.binding.ligands).forEach((lig) => {
      //     markResidues(this, res.binding.ligands[lig], "Near " + lig);
      //   });
      // }

      // if (res.fpocket.pockets) {
      //   res.fpocket.pockets.forEach((p) => {
      //     markResidues(
      //       this,
      //       p.residues.map((r) => r.residue),
      //       "Pocket"
      //     );
      //   });
      // }
    });
  }

  render() {
    return (
      <div
        id="fv"
        onMouseMove={this.handleMouseMove}
        onMouseLeave={this.handleMouseLeave}
      />
    );
  }
}
SequenceViewer.contextType = ResultsContext;
