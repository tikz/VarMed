import React from "react";
import "../../styles/components/sequence-viewer.scss";
import { ResultsContext } from "./ResultsContext";

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

  componentDidMount() {
    const structure = this.context.structure.current;
    const res = this.context.results;
    const posMap = this.context.posMap;

    const FeatureViewer = require("feature-viewer");
    this.fv = new FeatureViewer(res.UniProt.Sequence, "#fv", {
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

    this.loadFeatures();
  }

  loadFeatures() {
    const posMap = this.context.posMap;
    const res = this.context.results;

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
            .filter((r) => r.Chain == chain.id)
            .map((r) => {
              return {
                x: posMap.pdbToUnp(r.Chain, r.Position),
                y: posMap.pdbToUnp(r.Chain, r.Position),
                description: name,
              };
            }),
          name: chain.id + " - " + title,
          color: "#1aacdb",
          type: "rect",
        });
      };

      if (res.Exposure.Residues !== null) {
        markResidues(this, res.Exposure.Residues, "Buried");
      }
      if (res.Binding.Catalytic.Residues !== null) {
        markResidues(this, res.Binding.Catalytic.Residues, "Catalytic");
      }
      if (Object.keys(res.Binding.Ligands).length != 0) {
        Object.keys(res.Binding.Ligands).forEach((lig) => {
          markResidues(this, res.Binding.Ligands[lig], "Near " + lig);
        });
      }
      if (res.Interaction.Residues !== null) {
        markResidues(this, res.Interaction.Residues, "Interface");
      }
      if (res.Binding.Pockets !== null) {
        res.Binding.Pockets.forEach((p) => {
          markResidues(this, p.Residues, "Pocket");
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
