import { Chip } from "@material-ui/core";
import React from "react";
import { ResultsContext } from "./ResultsContext";

class FeatureChip extends React.Component {
  constructor(props) {
    super(props);
    this.handleMouseEnter = this.handleMouseEnter.bind(this);
    this.handleMouseLeave = this.handleMouseLeave.bind(this);
  }

  handleMouseEnter() {
    this.context.structure.current.highlightResidues(this.props.residues);
  }

  handleMouseLeave() {
    this.context.structure.current.clearHighlight();
  }

  render() {
    return (
      <Chip
        label={this.props.label}
        size="small"
        variant="outlined"
        onMouseEnter={this.handleMouseEnter}
        onMouseLeave={this.handleMouseLeave}
        className="chip"
      />
    );
  }
}

export class ChipRes extends FeatureChip {
  constructor(props) {
    super(props);
  }

  handleMouseEnter() {
    this.context.structure.current.highlightResidues(this.props.residues);
  }
}
ChipRes.contextType = ResultsContext;

export class ChipHet extends FeatureChip {
  constructor(props) {
    super(props);
  }

  handleMouseEnter() {
    this.context.structure.current.highlightHet(this.props.hetID);
  }
}
ChipHet.contextType = ResultsContext;
