import { Chip } from "@material-ui/core";
import React from "react";

export default class ChipHet extends React.Component {
  constructor(props) {
    super(props);
    this.handleMouseEnter = this.handleMouseEnter.bind(this);
    this.handleMouseLeave = this.handleMouseLeave.bind(this);
  }

  handleMouseEnter() {
    this.props.highlightHet(this.props.hetID);
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
