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
      const chain = d.detail.description.split(" ")[1];
      structure.selectFocus(
        chain,
        d.detail.start + posMap.pdbOffsets[chain],
        d.detail.end + posMap.pdbOffsets[chain]
      );
    });

    this.setState(() => ({
      zoomPositionElement: document.getElementById("zoomPosition"),
    }));

    posMap.chains.forEach((chain) => {
      let name = "Chain " + chain.id;
      this.fv.addFeature({
        data: [
          {
            x: chain.start,
            y: chain.end,
            description: name,
          },
        ],
        name: name,
        color: "#1aacdb",
        type: "rect",
      });

      let markResidues = function (that, residues, title) {
        that.fv.addFeature({
          data: residues
            .filter((r) => r.chain == chain.id)
            .map((r) => {
              return {
                x: posMap.pdbToUnp(r.chain, r.position),
                y: posMap.pdbToUnp(r.chain, r.position),
                description: name,
              };
            }),
          name: chain.id + " - " + title,
          color: "#1aacdb",
          type: "rect",
        });
      };

      if (res.exposure.residues !== null) {
        markResidues(this, res.exposure.residues, "Buried");
      }
      if (res.binding.catalytic.residues !== null) {
        markResidues(this, res.binding.catalytic.residues, "Catalytic");
      }
      if (Object.keys(res.binding.ligands).length != 0) {
        Object.keys(res.binding.ligands).forEach((lig) => {
          markResidues(this, res.binding.ligands[lig], "Near " + lig);
        });
      }
      if (res.interaction.residues !== null) {
        markResidues(this, res.interaction.residues, "Interface");
      }
      if (res.binding.pockets !== null) {
        res.binding.pockets.forEach((p) => {
          markResidues(this, p.residues, "Pocket");
        });
      }
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
