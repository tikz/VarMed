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
import { Variants } from "./Variants";
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
      variants: [],
      error: false,
      errorMsg: "",
      redirect: "",
    };

    this.setUnpData = this.setUnpData.bind(this);
    this.setPDBs = this.setPDBs.bind(this);
    this.setVars = this.setVars.bind(this);
    this.setAnnotated = this.setAnnotated.bind(this);
    this.submit = this.submit.bind(this);
    this.handleErrorClose = this.handleErrorClose.bind(this);
  }

  setUnpData(unpData) {
    this.setState({
      unpData: unpData,
      pdbs: [],
      variants: [],
    });
  }

  setPDBs(pdbs) {
    this.setState({ pdbs: pdbs });
  }

  setVars(vars) {
    vars.sort((a, b) => (a.pos + a.aa > b.pos + b.aa ? 1 : -1));
    this.setState({ variants: vars });
  }

  setAnnotated(toggle) {
    let vars = this.state.variants;
    let keys = vars.map((v) => v.key);
    if (toggle) {
      this.state.unpData.variants.forEach((v) => {
        if (!keys.includes(v.change)) {
          vars.push({
            key: v.change,
            pos: v.position,
            aa: v.toAa,
            label: v.position + " " + v.fromAa + "→" + v.toAa,
          });
        }
      });
    } else {
      let annVars = this.state.unpData.variants.map((v) => v.change);
      vars = vars.filter((v) => !annVars.includes(v.key));
    }
    vars.sort((a, b) => (a.pos + a.aa > b.pos + b.aa ? 1 : -1));
    this.setState({ variants: vars });
  }

  submit(email) {
    let that = this;
    axios
      .post(API_URL + "/api/new-job", {
        name: this.state.unpData.entryName,
        uniprotId: this.state.unpData.id,
        pdbIds: this.state.pdbs,
        email: email,
        variants: this.state.variants.map((v) => v.key),
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
    let dataOk =
      unpOk && this.state.pdbs.length > 0 && this.state.variants.length > 0;
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
                <Grid item xs={12} lg={2}>
                  <UniProtInput setUnpData={this.setUnpData} />
                </Grid>

                <Grid item>
                  <Grid
                    container
                    alignItems="center"
                    justify="flex-start"
                    spacing={2}
                  >
                    <Grid item>{unpOk && <ArrowForwardIosIcon />}</Grid>
                    <Grid item xs>
                      <Grow in={unpOk}>
                        <Box>
                          {unpOk && (
                            <PDBPicker
                              unpID={this.state.unpData.id}
                              pdbs={this.state.unpData.pdbs}
                              setPDBs={this.setPDBs}
                            />
                          )}
                        </Box>
                      </Grow>
                    </Grid>
                  </Grid>
                </Grid>

                <Grid item xs={12} lg={4}>
                  <Grid
                    container
                    alignItems="center"
                    justify="flex-start"
                    spacing={2}
                  >
                    <Grid item>
                      {unpOk && structOk && <ArrowForwardIosIcon />}
                    </Grid>
                    <Grid item xs>
                      <Grow in={unpOk}>
                        <Box>
                          {unpOk && structOk && (
                            <Variants
                              unpID={this.state.unpData.id}
                              sequence={this.state.unpData.sequence}
                              variants={this.state.variants}
                              setVariants={this.setVars}
                              setAnnotated={this.setAnnotated}
                              hasAnnotated={
                                this.state.unpData.variants !== null
                              }
                            />
                          )}
                        </Box>
                      </Grow>
                    </Grid>
                  </Grid>
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
