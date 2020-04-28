import {
  Box,
  Container,
  Divider,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Select,
  Typography,
} from "@material-ui/core";
import axios from "axios";
import React from "react";
import FPSStats from "react-fps-stats";
import { Features } from "./Features";
import PositionMapper from "./PositionMapper";
import { ResultsContext } from "./ResultsContext";
import SequenceViewer from "./SequenceViewer";
import StructureViewer from "./StructureViewer";

export default class Results extends React.Component {
  constructor(props) {
    super(props);
    this.structureRef = React.createRef();

    this.state = {
      pdbID: this.props.pdbID,
      jobID: this.props.jobID,
      results: {},
    };

    this.pdbChange = this.pdbChange.bind(this);

    this.pdbLoad(this.props.pdbID);
  }

  pdbChange(e) {
    let id = e.target.value;
    this.pdbLoad(id);
  }

  pdbLoad(id) {
    let that = this;
    axios
      .get(API_URL + "/api/job/" + this.state.jobID + "/" + id)
      .then(function (response) {
        that.setState({ results: response.data, pdb: id });
        that.structureRef.current.load(response.data);
      });
  }

  render() {
    if (this.state.results.PDB === undefined) {
      return <Box />;
    }

    const ctx = {
      structure: this.structureRef,
      results: this.state.results,
      posMap: new PositionMapper(this.state.results),
    };

    return (
      <ResultsContext.Provider value={ctx}>
        <FPSStats left="auto" top="auto" right="0" bottom="0" />
        <Container>
          <Box className="over">
            <Typography variant="h4" className="title">
              {this.state.results.UniProt.Name}
            </Typography>
            <Divider />
            <Grid container spacing={2} alignItems="center">
              <Grid item>
                <FormControl variant="outlined">
                  <InputLabel>PDB</InputLabel>
                  <Select
                    label="PDB"
                    value={this.state.pdb}
                    onChange={this.pdbChange}
                  >
                    {this.props.results.Request.pdbs.map((pdbID, index) => {
                      return (
                        <MenuItem key={index} value={pdbID}>
                          {pdbID}
                        </MenuItem>
                      );
                    })}
                  </Select>
                </FormControl>
              </Grid>
              <Grid item>
                <Typography>{this.state.results.PDB.Title}</Typography>
              </Grid>
            </Grid>
          </Box>
        </Container>

        <StructureViewer
          ref={this.structureRef}
          pdbID={this.state.pdbID}
          res={this.state.results}
        />

        <Container>
          <Box>
            <Features />
          </Box>
          <Box>
            <SequenceViewer />
          </Box>
        </Container>
      </ResultsContext.Provider>
    );
  }
}
