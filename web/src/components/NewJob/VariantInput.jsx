import { Box, Grid, IconButton, TextField } from "@material-ui/core";
import { Add } from "@material-ui/icons";
import React from "react";

export default class VariantInput extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      aaValue: "",
      posError: false,
      posErrorMsg: "",
      aaError: false,
      aaErrorMsg: "",
      pos: "",
      aa: "",
    };

    this.handleAaChange = this.handleAaChange.bind(this);
    this.handlePosChange = this.handlePosChange.bind(this);
    this.handleAdd = this.handleAdd.bind(this);
    this.handleKey = this.handleKey.bind(this);
  }

  handleAdd() {
    if (
      this.state.aa != "" &&
      this.state.pos > 0 &&
      !this.state.posError &&
      !this.state.aaError
    ) {
      let variants = this.props.variants;
      variants.push({
        key:
          this.props.sequence[this.state.pos - 1] +
          this.state.pos +
          this.state.aa,
        pos: this.state.pos,
        aa: this.state.aa,
        label:
          this.state.pos +
          " " +
          this.props.sequence[this.state.pos - 1] +
          "→" +
          this.state.aa,
      });
      this.setState({ pos: "", aa: "" });
      this.props.setVariants(variants);
    }
  }

  handlePosChange(e) {
    let pos = parseInt(e.target.value);
    if (isNaN(pos)) {
      pos = "";
    }
    this.setState({ pos: pos }, () => {
      this.checkFields();
    });
  }

  handleAaChange(e) {
    this.setState({ aa: e.target.value.toUpperCase() }, () => {
      this.checkFields();
    });
  }

  handleKey(e) {
    if (e.key === "Enter") {
      this.handleAdd();
    }
  }

  checkFields() {
    const aminoacids = "ARNDCEQGHILKMFPSTWYV";
    let pos = this.state.pos;
    let aa = this.state.aa;
    let fromAa = this.props.sequence[pos - 1];
    pos > this.props.sequence.length || pos < 1
      ? this.setState({
          posError: true,
          posErrorMsg: "Outside sequence length",
        })
      : this.setState({ posError: false, posErrorMsg: "" });
    aa.length > 1 || !aminoacids.includes(aa)
      ? this.setState({ aaError: true, aaErrorMsg: "Unknown" })
      : this.setState({ aaError: false, aaErrorMsg: "" });
    if (this.props.sequence[pos - 1] == aa) {
      this.setState({
        posError: true,
        posErrorMsg: "No change",
        aaError: true,
        aaErrorMsg: "No change",
      });
    }

    let change = fromAa + pos + aa;
    if (this.props.variants.filter((v) => v.key == change).length > 0) {
      this.setState({ posError: true, posErrorMsg: "Variant already added" });
    }
  }

  render() {
    return (
      <Box>
        <Grid container spacing={1} alignItems="flex-start">
          <Grid item xs>
            <TextField
              label="Position"
              onChange={this.handlePosChange}
              type="number"
              error={this.state.posError}
              helperText={this.state.posErrorMsg}
              value={this.state.pos}
              onKeyPress={this.handleKey}
            />
          </Grid>
          <Grid item>
            <Grid container alignItems="flex-start">
              <Grid item>
                <TextField
                  onChange={this.handleAaChange}
                  label="1-letter Aa"
                  error={this.state.aaError}
                  helperText={this.state.aaErrorMsg}
                  value={this.state.aa}
                  onKeyPress={this.handleKey}
                />
              </Grid>
              <Grid item>
                <IconButton
                  color="primary"
                  aria-label="add variant"
                  onClick={this.handleAdd}
                  style={{ marginTop: 20 }}
                >
                  <Add />
                </IconButton>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Box>
    );
  }
}
