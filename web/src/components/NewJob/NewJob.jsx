import {
  Box,
  Container,
  Divider,
  Grid,
  Grow,
  Snackbar,
  Typography,
  Toolbar,
} from "@material-ui/core";
import ArrowForwardIosIcon from "@material-ui/icons/ArrowForwardIos";
import MuiAlert from "@material-ui/lab/Alert";
import axios from "axios";
import React from "react";
import { Redirect } from "react-router-dom";
import PDBPicker from "./PDBPicker";
import SendBar from "./SendBar";
import { UniProtInput } from "./UniProtInput";
import { Variations } from "./Variations";
import NavBar from "../NavBar";

function Alert(props) {
  return <MuiAlert elevation={6} variant="filled" {...props} />;
}

export default class NewJob extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      unpData: {},
      pdbs: [],
      variations: [],
      clinvar: false,
      error: false,
      errorMsg: "",
      redirect: "",
    };

    this.setUnpData = this.setUnpData.bind(this);
    this.setPDBs = this.setPDBs.bind(this);
    this.setVars = this.setVars.bind(this);
    this.setClinVar = this.setClinVar.bind(this);
    this.submit = this.submit.bind(this);
    this.handleErrorClose = this.handleErrorClose.bind(this);
  }

  setUnpData(unpData) {
    this.setState({
      unpData: unpData,
      pdbs: [],
      variations: [],
      clinvar: false,
    });
  }

  setPDBs(pdbs) {
    this.setState({ pdbs: pdbs });
  }
  setVars(vars) {
    this.setState({ variations: vars });
  }
  setClinVar(flag) {
    this.setState({ clinvar: flag });
  }

  submit(email) {
    let that = this;
    axios
      .post(API_URL + "/api/new-job", {
        name: this.state.unpData.entryName,
        uniprotId: this.state.unpData.id,
        pdbIds: this.state.pdbs,
        email: email,
        sas: ["M1K"],
        // variationsPos: this.state.variations.map((x) => x.pos),
        // variationsAa: this.state.variations.map((x) => x.aa),
        clinvar: this.state.clinvar,
      })
      .then(function (response) {
        if (response.data.error != "") {
          that.setState({ errorMsg: response.data.error }, () => {
            that.setState({ error: true });
          });
        } else {
          that.setState({ redirect: "/job/" + response.data.id });
        }
      });
  }

  handleErrorClose() {
    this.setState({ error: false });
  }

  render() {
    if (this.state.redirect != "") {
      return <Redirect to={this.state.redirect} />;
    }

    let unpOk = Object.keys(this.state.unpData).length > 0;
    let structOk = this.state.unpData.pdbs !== null;
    let dataOk = unpOk && this.state.pdbs.length > 0; // && (this.state.variations.length > 0 || this.state.clinvar);
    return (
      <Box>
        <NavBar />
        <Toolbar />
        <Container>
          <Typography variant="h2" gutterBottom className="title">
            New Job
          </Typography>
          <Grid container spacing={4} direction="column">
            <Grid item>
              <Grid container spacing={2} alignItems="center">
                <Grid item xs={2}>
                  <UniProtInput setUnpData={this.setUnpData} />
                </Grid>
                <Grid item>{unpOk && <ArrowForwardIosIcon />}</Grid>
                <Grid item>
                  <Grow in={unpOk}>
                    <Box>
                      {unpOk && (
                        <PDBPicker
                          unpID={this.state.unpData.id}
                          pdbs={this.state.unpData.pdbIds}
                          setPDBs={this.setPDBs}
                        />
                      )}
                    </Box>
                  </Grow>
                </Grid>
                <Grid item>
                  {false && unpOk && structOk && <ArrowForwardIosIcon />}
                </Grid>
                <Grid item xs={4}>
                  <Grow in={unpOk}>
                    <Box>
                      {false && unpOk && structOk && (
                        <Variations
                          unpID={this.state.unpData.id}
                          sequence={this.state.unpData.sequence}
                          setVariations={this.setVars}
                          setClinVar={this.setClinVar}
                        />
                      )}
                    </Box>
                  </Grow>
                </Grid>
              </Grid>
            </Grid>
            {dataOk && <Divider />}
            <Grid item>
              <Grow in={dataOk}>
                <Box>
                  <SendBar submit={this.submit} />
                </Box>
              </Grow>
            </Grid>
            {dataOk && <Divider />}
          </Grid>
        </Container>
        <Snackbar
          open={this.state.error}
          autoHideDuration={6000}
          onClose={this.handleErrorClose}
        >
          <Alert onClose={this.handleErrorClose} severity="error">
            {this.state.errorMsg}
          </Alert>
        </Snackbar>
      </Box>
    );
  }
}
