import { FormControlLabel, Switch } from "@material-ui/core";
import React from "react";
import { ResultsContext } from "./ResultsContext";

export class MutantSwitch extends React.Component {
  constructor(props) {
    super(props);
    this.state = { mutant: false };
    this.handleChange = this.handleChange.bind(this);
  }
  handleChange(e) {
    this.setState({ mutant: e.target.checked });
    this.context.structure.current.showMutated(e.target.checked);
    this.context.structure.current.setState({ mutated: e.target.checked });
  }
  render() {
    return (
      <FormControlLabel
        control={
          <Switch checked={this.state.mutant} onChange={this.handleChange} />
        }
        label="Mutant"
      />
    );
  }
}
MutantSwitch.contextType = ResultsContext;
