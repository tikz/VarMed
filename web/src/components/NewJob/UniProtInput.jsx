import {
  Box,
  CircularProgress,
  InputAdornment,
  Typography,
} from "@material-ui/core";
import TextField from "@material-ui/core/TextField";
import axios from "axios";
import React from "react";

export class UniProtInput extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      unpID: "",
      unpData: {},
      loading: false,
      error: false,
      errorMsg: "",
      entryName: "",
    };

    this.handleChange = this.handleChange.bind(this);
  }

  handleChange(e) {
    this.props.setUnpData({});
    this.setState({
      unpID: e.target.value.toUpperCase(),
      error: false,
      errorMsg: "",
      entryName: "",
    });

    if (e.target.value.length > 5) {
      let that = this;
      that.setState({
        loading: true,
      });
      axios
        .get(API_URL + "/api/uniprot/" + e.target.value)
        .then(function (response) {
          let data = response.data;
          that.setState({
            entryName: [data.gene, data.name, data.organism].join(" - "),
            unpData: data,
          });
          that.props.setUnpData(data);
        })
        .catch(function (error) {
          that.setState({
            error: true,
            errorMsg: "Network error",
          });
        })
        .then(function () {
          that.setState({
            loading: false,
          });
        });
    } else {
      this.setState({
        error: false,
        errorMsg: "",
        entryName: "",
      });
    }
  }

  render() {
    return (
      <Box>
        <Typography variant="h5" gutterBottom>
          1. Enter a protein
        </Typography>
        <div>
          <TextField
            id="filled-basic"
            variant="filled"
            autoFocus
            autoComplete="off"
            fullWidth
            label="UniProt Accession ID"
            onChange={this.handleChange}
            value={this.state.unpID}
            error={this.state.error}
            helperText={this.state.errorMsg}
            InputProps={{
              endAdornment: (
                <InputAdornment position="end">
                  {this.state.loading && <CircularProgress />}
                </InputAdornment>
              ),
            }}
          />
        </div>
        <Typography variant="overline" gutterBottom>
          {this.state.entryName}
        </Typography>
      </Box>
    );
  }
}
