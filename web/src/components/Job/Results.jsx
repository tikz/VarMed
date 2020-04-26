import {
  Box,
  Container,
  Divider,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Typography,
} from "@material-ui/core";
import axios from "axios";
import React from "react";
import SequenceViewer from "./SequenceViewer";
import StructureViewer from "./StructureViewer";

export default class Results extends React.Component {
  constructor(props) {
    super(props);
    this.structureRef = React.createRef();

    this.state = { pdbID: this.props.pdbID, jobID: this.props.jobID, res: {} };

    this.highlightResidues = this.highlightResidues.bind(this);
    this.select = this.select.bind(this);
    this.clearHighlight = this.clearHighlight.bind(this);
    this.pdbChange = this.pdbChange.bind(this);

    this.pdbLoad(this.props.pdbID);
  }

  highlightResidues(residues) {
    if (residues.length == 0) {
      this.structureRef.current.clearHighlight();
    } else {
      this.structureRef.current.highlightResidues(residues);
    }
  }

  select(chain, start, end) {
    this.structureRef.current.focus(chain, start, end);
    this.structureRef.current.highlight(chain, start, end);
    if (start - end == 0) {
      this.structureRef.current.select(chain, start, end);
    }
  }

  clearHighlight() {
    this.structureRef.current.clearHighlight();
  }

  pdbChange(e) {
    let id = e.target.value;
    this.pdbLoad(id);
  }

  pdbLoad(id) {
    let that = this;
    axios
      .get("http://localhost:8888/api/job/" + this.state.jobID + "/" + id)
      .then(function (response) {
        that.setState({ res: response.data, pdb: id });
        that.structureRef.current.load(response.data);
      });
  }

  render() {
    if (this.state.res.PDB === undefined) {
      return <Box />;
    }
    return (
      <Box>
        <Container>
          <Box className="over">
            <Typography variant="h4" className="title">
              {this.state.res.UniProt.Name}
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
                <Typography>{this.state.res.PDB.Title}</Typography>
              </Grid>
            </Grid>
          </Box>
        </Container>

        <StructureViewer
          ref={this.structureRef}
          pdbID={this.state.pdbID}
          res={this.state.res}
        />

        <Container>
          <Box>
            <Grid container spacing={3}>
              <Grid item xs={5}></Grid>
              <Grid item xs={7}>
                <Paper></Paper>
              </Grid>
            </Grid>
          </Box>
          <Box my={2}>
            <SequenceViewer
              highlightResidues={this.highlightResidues}
              select={this.select}
              clearHighlight={this.clearHighlight}
              res={this.state.res}
              key={this.state.res.PDB.ID}
            />
          </Box>
        </Container>
      </Box>
    );
  }
}
