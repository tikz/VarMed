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
              <img src="/assets/varmed.svg" alt="VarMed" className="logo" />
            </Grid>
            <Grid item xs={9} sm={4} lg={3}>
              <Grid
                container
                direction="column"
                alignItems="flex-start"
                justify="center"
              >
                <Typography variant="h1" align="left" className="name">
                  VarMed
                </Typography>
                <Typography
                  variant="h5"
                  align="left"
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
            <LinkRouter to="/job/15e20e5f18326d264b60eeaa07c9af8d04b0a6c70f037b7f69b6d40d22fb590b">
              <Button variant="outlined">Sample Results</Button>
            </LinkRouter>
          </Typography>
        </Grid>
        <Grid item>
          <Typography align="center" className="desc">
            If you find our work useful, please cite us:
          </Typography>
          <Typography align="center" className="desc paper">
            pending
          </Typography>
          <Typography align="center" className="desc authors">
            Mauro Song, Florencia Niesi, Demian Avendaño
          </Typography>
          <Typography align="center" className="desc authors">
            Marcelo Martí, Pietro Roversi, Carlos Modenutti
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
                    VarMed <Link href="#">source code</Link> is released under
                    the <Link href="#">MIT license</Link>.
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
