import { Box, Button, Container, Grid, Typography } from "@material-ui/core";
import React from "react";
import { Link as LinkRouter } from "react-router-dom";
import "../styles/components/index.scss";
import SplashBackground from "./SplashBackground.jsx";

export default function Index() {
  return (
    <Box>
      <SplashBackground />
      <Box className="index">
        <Container>
          <Grid
            className="presentation"
            container
            direction="column"
            justify="center"
            alignItems="center"
            spacing={3}
          >
            <Grid item xs={3} md={12}>
              <Grid
                container
                spacing={3}
                direction="row"
                alignItems="center"
                justify="center"
              >
                <Grid item>
                  <img src="/assets/varmed.svg" alt="VarMed" className="logo" />
                </Grid>
                <Grid item>
                  <Grid
                    container
                    direction="column"
                    alignItems="center"
                    justify="center"
                  >
                    <Typography variant="h1" align="center" className="name">
                      VarMed
                    </Typography>
                    <Typography
                      variant="h5"
                      align="center"
                      className="desc"
                    ></Typography>
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
            <Grid item>
              <Typography align="center" className="desc">
                Start a{" "}
                <LinkRouter to="/new-job">
                  <Button className="glowButton">New Job</Button>
                </LinkRouter>{" "}
                or view{" "}
                <LinkRouter to="/job/ba2388afd68a4b467fc2c1b6e81a301a8341f98bb5ea564d2d0b483d165f9c4c">
                  <Button variant="outlined">Sample Results</Button>
                </LinkRouter>
              </Typography>
            </Grid>
            <Grid item className="cite-us">
              <Typography align="center" className="desc">
                If you find our work useful, please cite us:
              </Typography>
              <Typography align="center" className="desc paper">
                Pending
              </Typography>
              <Typography align="center" className="desc authors">
              </Typography>
              <Typography align="center" className="desc authors">
              </Typography>
            </Grid>
          </Grid>
        </Container>

        <Box className="footer">
          <Container className="footer">
            <Grid container direction="row" justify="space-between" spacing={5}>
              <Grid item lg={6} xl={12}>
                <Grid container direction="column">
                  <Grid item>
                    <Typography variant="caption">
                      Bioinformática Estructural y Biofisicoquímica de
                      Proteínas.
                    </Typography>
                  </Grid>
                  <Grid item>
                    <Typography variant="caption">
                      IQUIBICEN, Departamento de Química Biológica.
                    </Typography>
                  </Grid>
                  <Grid item>
                    <Typography variant="caption">
                      Facultad de Ciencias Exactas y Naturales, Universidad de
                      Buenos Aires.
                    </Typography>
                  </Grid>
                </Grid>
              </Grid>
              <Grid item lg={6} xl={12}>
                <Grid
                  container
                  direction="column"
                  alignItems="flex-end"
                  spacing={1}
                >
                  {/* <Grid item>
                    <Typography variant="caption">
                      VarMed <Link href="#">source code</Link> is released under
                      the <Link href="#">MIT license</Link>.
                    </Typography>
                  </Grid>
                  <Grid item>
                    <Typography variant="caption">
                      External tools and libraries may have different licenses.
                    </Typography>
                  </Grid> */}
                </Grid>
              </Grid>
            </Grid>
          </Container>
        </Box>
      </Box>
    </Box>
  );
}
