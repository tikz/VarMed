import {
  Box,
  Button,
  Container,
  Grid,
  Link,
  Typography,
} from "@material-ui/core";
import React from "react";
import { Link as LinkRouter } from "react-router-dom";
import "../styles/components/index.scss";

export default function Index() {
  return (
    <Box>
      <Grid
        className="presentation"
        container
        direction="column"
        justify="center"
        spacing={3}
      >
        <Grid item>
          <Grid
            container
            spacing={6}
            direction="row"
            alignItems="center"
            justify="center"
          >
            <Grid item>
              <img src="/assets/varq.svg" alt="VarQ" className="logo" />
            </Grid>
            <Grid item xs={9} sm={4} lg={3}>
              <Grid
                container
                direction="column"
                alignItems="flex-start"
                justify="center"
              >
                <Typography variant="h1" align="left" className="name">
                  VarQ
                </Typography>
                <Typography variant="h5" align="left" className="desc">
                  A tool for the structural and functional analysis of protein
                  variants.
                </Typography>
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
            <LinkRouter to="/job/aa2725a483568c283274c6e551b83ac1c34548736c0dbb2581ba770bb0de21eb">
              <Button variant="outlined">Sample Results</Button>
            </LinkRouter>
          </Typography>
        </Grid>
        <Grid item>
          <Typography align="center" className="desc">
            If you find our work useful, please cite us: <br /> -
          </Typography>
        </Grid>
      </Grid>
      <Grid container className="footer">
        <Container>
          <Grid container direction="row" justify="space-between">
            <Grid item>
              <Grid container direction="column">
                <Grid item>
                  <Typography variant="caption">
                    Bioinformática Estructural y Biofisicoquímica de Proteínas.
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
            <Grid item>
              <Grid
                container
                direction="column"
                justify="space-between"
                alignItems="flex-end"
              >
                <Grid item>
                  <Typography variant="caption">
                    VarQ <Link href="#">source code</Link> is released under the{" "}
                    <Link href="#">MIT license</Link>.
                  </Typography>
                </Grid>
                <Grid item>
                  <Typography variant="caption">
                    External tools and libraries may have different licenses.
                  </Typography>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </Container>
      </Grid>
    </Box>
  );
}
