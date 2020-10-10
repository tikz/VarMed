import {
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
import "../../styles/components/results.scss";
import VariantViewer from "./VariantViewer";

export default class Results extends React.Component {
  constructor(props) {
    super(props);
    this.structureRef = React.createRef();
    this.sequenceRef = React.createRef();

    this.state = { results: {} };

    this.handlePDBChange = this.handlePDBChange.bind(this);
  }

  componentDidMount() {
    this.pdbLoadFirst();
  }

  componentDidUpdate(prevProps) {
    if (this.props.jobId != prevProps.jobId) {
      this.pdbLoadFirst();
    }
  }

  pdbLoadFirst() {
    this.pdbLoad(this.props.jobResults.request.pdbIds[0]);
  }

  handlePDBChange(e) {
    this.pdbLoad(e.target.value);
  }

  pdbLoad(pdbId) {
    const jobId = this.props.jobId;
    let that = this;
    axios
      .get(API_URL + "/api/job/" + jobId + "/" + pdbId)
      .then(function (response) {
        that.setState({ results: response.data, pdb: pdbId });
        that.structureRef.current.load(response.data);
        that.sequenceRef.current.load();
      });
  }

  render() {
    if (this.state.results.pdb === undefined) {
      return <div />;
    }

    const ctx = {
      structure: this.structureRef,
      results: this.state.results,
      posMap: new PositionMapper(this.state.results),
    };

    return (
      <ResultsContext.Provider value={ctx}>
        <FPSStats left="auto" top="auto" right="0" bottom="0" />
        <div className="left split">
          <Container>
            <div className="over-title">
              <Typography variant="h4" className="title">
                {this.state.results.uniprot.name}
              </Typography>
              <Divider />
              <Grid container spacing={2} alignItems="center">
                <Grid item>
                  <FormControl variant="outlined">
                    <InputLabel>PDB</InputLabel>
                    <Select
                      label="PDB"
                      value={this.state.pdb}
                      onChange={this.handlePDBChange}
                    >
                      {this.props.jobResults.request.pdbIds.map(
                        (pdbId, index) => {
                          return (
                            <MenuItem key={index} value={pdbId}>
                              {pdbId}
                            </MenuItem>
                          );
                        }
                      )}
                    </Select>
                  </FormControl>
                </Grid>
                <Grid item>
                  <Typography>{this.state.results.pdb.title}</Typography>
                </Grid>
              </Grid>
            </div>
            <div className="over-features">
              <Features />
            </div>
          </Container>

          <StructureViewer ref={this.structureRef} />
        </div>

        <div className="right split">
          <Container>
            <VariantViewer />
            <SequenceViewer ref={this.sequenceRef} />
          </Container>
        </div>
      </ResultsContext.Provider>
    );
  }
}
