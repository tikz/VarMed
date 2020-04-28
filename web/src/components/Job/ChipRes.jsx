import { Chip } from "@material-ui/core";
import React from "react";

export default class ChipRes extends React.Component {
  constructor(props) {
    super(props);
    this.handleMouseEnter = this.handleMouseEnter.bind(this);
    this.handleMouseLeave = this.handleMouseLeave.bind(this);
  }

  handleMouseEnter() {
    this.props.highlightResidues(this.props.residues);
  }

  handleMouseLeave() {
    this.props.clearHighlight();
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
