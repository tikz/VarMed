import { FormControlLabel, Switch } from "@material-ui/core";
import React from "react";

export class SurfaceSwitch extends React.Component {
  constructor(props) {
    super(props);
    this.state = { surface: true };
    this.handleChange = this.handleChange.bind(this);
  }
  handleChange(e) {
    this.setState({ surface: e.target.checked });
    this.props.showSurface(e.target.checked);
  }
  render() {
    return (
      <FormControlLabel
        control={
          <Switch checked={this.state.surface} onChange={this.handleChange} />
        }
        label="Surface"
      />
    );
  }
}
