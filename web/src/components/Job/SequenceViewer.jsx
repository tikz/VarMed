import React from "react";

export default class SequenceViewer extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      pos: 0,
    };

    this.handleMouseMove = this.handleMouseMove.bind(this);
    this.handleMouseLeave = this.handleMouseLeave.bind(this);
  }

  handleMouseMove() {
    var posText = this.state.zoomPositionElement.innerText;
    if (posText != this.state.pos) {
      var pos = parseInt(posText.slice(0, posText.length - 1));
      this.props.highlightResidues(this.unpToPDB(pos));
      this.setState((state) => ({
        pos: posText,
      }));
    }
  }

  handleMouseLeave() {
    this.props.highlight(0, 0);
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

    let selectFunc = this.props.select;
    let that = this;
    this.fv.onFeatureSelected(function (d) {
      let chain = d.detail.description.split(" ")[1];
      console.log(d.detail.start + "->");
      selectFunc(
        chain,
        d.detail.start + that.pdbOffsets[chain],
        d.detail.end + that.pdbOffsets[chain]
      );
    });

    this.setState((state) => ({
      zoomPositionElement: document.getElementById("zoomPosition"),
    }));

    this.loadOffsets(res);
    this.loadFeatures(res);
  }

  loadOffsets(res) {
    let unpID = res.UniProt.ID;

    // Calculate position offsets
    this.seqResOffsets = res.PDB.SeqResOffsets;
    this.chainStartResN = res.PDB.ChainStartResNumber;
    this.unpOffsets = {};
    res.PDB.SIFTS.UniProt[unpID].mappings.forEach((chain) => {
      this.unpOffsets[chain.chain_id] =
        chain.unp_start +
        this.seqResOffsets[chain.chain_id] -
        chain.start.residue_number -
        this.chainStartResN[chain.chain_id] +
        1;
    });

    this.pdbOffsets = {};
    res.PDB.SIFTS.UniProt[unpID].mappings.forEach((chain) => {
      this.pdbOffsets[chain.chain_id] =
        -chain.unp_start +
        this.seqResOffsets[chain.chain_id] -
        this.chainStartResN[chain.chain_id] +
        2;
    });
  }

  pdbToUnp(chain, pos) {
    return pos + this.unpOffsets[chain];
  }

  unpToPDB(pos) {
    let unpID = this.props.res.UniProt.ID;
    let residues = [];
    this.props.res.PDB.SIFTS.UniProt[unpID].mappings.forEach((chain) => {
      residues.push({
        chain: chain.chain_id,
        pos: pos + this.pdbOffsets[chain.chain_id],
      });
    });
    return residues;
  }

  loadFeatures(res) {
    let unpID = res.UniProt.ID;
    let unpChains = []; // chains in unpID

    // Chains
    res.PDB.SIFTS.UniProt[unpID].mappings.forEach((chain) => {
      let name = "Chain " + chain.chain_id;
      unpChains.push(chain.chain_id);
      this.fv.addFeature({
        data: [
          {
            x:
              chain.unp_start +
              this.seqResOffsets[chain.chain_id] -
              chain.start.residue_number +
              1,
            y:
              chain.unp_end +
              this.seqResOffsets[chain.chain_id] -
              chain.start.residue_number +
              1,
            description: name,
          },
        ],
        name: name,
        className: "test1",
        color: "#2196F3",
        type: "rect",
        filter: "type1",
      });

      let markResidues = function (that, residues, title) {
        that.fv.addFeature({
          data: residues
            .filter((r) => r.Chain == chain.chain_id)
            .map((r) => {
              return {
                x: that.pdbToUnp(r.Chain, r.Position),
                y: that.pdbToUnp(r.Chain, r.Position),
                description: name,
              };
            }),
          name: chain.chain_id + " - " + title,
          className: "test1",
          color: "#2196F3",
          type: "rect",
          filter: "type1",
        });
      };
      if (res.Exposure.Residues != null) {
        markResidues(this, res.Exposure.Residues, "Buried");
      }
      if (res.Binding.Catalytic != null) {
        markResidues(this, res.Binding.Catalytic.Residues, "Catalytic");
      }
      if (Object.keys(res.Binding.Ligands).length !== 0) {
        Object.keys(res.Binding.Ligands).forEach((lig) => {
          markResidues(this, res.Binding.Ligands[lig], "Near " + lig);
        });
      }
      if (res.Interaction.Residues != null) {
        markResidues(this, res.Interaction.Residues, "Interaction");
      }
    });

    // // Pfam
    // for (const [id, fam] of Object.entries(res.PDB.SIFTS.Pfam)) {
    //     fam.Mappings.forEach(m => {
    //         let famsData = []
    //         if (unpChains.includes(m.chain_id)) {
    //             let off = this.unpOffsets[m.chain_id]
    //             let desc = id + " - " + fam.Description
    //             famsData.push({ x: m.start.residue_number - off, y: m.end.residue_number - off, description: desc });

    //             this.fv.addFeature({
    //                 data: famsData,
    //                 name: "Pfam",
    //                 className: "test1",
    //                 color: "#2196F3",
    //                 type: "rect",
    //                 filter: "type1"
    //             });
    //         }
    //     })
    // }

    // TODO: repeated code, refactor
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
