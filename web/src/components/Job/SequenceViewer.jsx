import React from "react";
import "../../styles/components/sequence-viewer.scss";
import PositionMapper from "./PositionMapper";

export default class SequenceViewer extends React.Component {
  constructor(props) {
    super(props);
    this.state = { pos: 0 };

    this.posMap = new PositionMapper(this.props.res);

    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseLeave = this.handleMouseLeave.bind(this);
  }

  handleMouseMove() {
    var posText = this.state.zoomPositionElement.innerText;
    if (posText != this.state.pos) {
      var pos = parseInt(posText.slice(0, posText.length - 1));
      this.props.highlightResidues(this.posMap.unpToPDB(pos));
      this.setState(() => ({ pos: posText }));
    }
  }

  handleMouseLeave() {
    this.props.clearHighlight();
  }

  componentDidMount() {
    let res = this.props.res;

    var FeatureViewer = require("feature-viewer");
    this.fv = new FeatureViewer(res.UniProt.Sequence, "#fv", {
      showAxis: true,
      showSequence: true,
      brushActive: true,
      toolbar: true,
      bubbleHelp: false,
      zoomMax: 20,
    });

    let that = this;
    this.fv.onFeatureSelected(function (d) {
      let chain = d.detail.description.split(" ")[1];
      that.props.select(
        chain,
        d.detail.start + that.posMap.pdbOffsets[chain],
        d.detail.end + that.posMap.pdbOffsets[chain]
      );
    });

    this.setState(() => ({
      zoomPositionElement: document.getElementById("zoomPosition"),
    }));

    this.loadFeatures();
  }

  loadFeatures() {
    let res = this.props.res;

    this.posMap.chains.forEach((chain) => {
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
                x: that.posMap.pdbToUnp(r.Chain, r.Position),
                y: that.posMap.pdbToUnp(r.Chain, r.Position),
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
