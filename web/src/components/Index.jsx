import { Box, Button, Container, Grid, Link, makeStyles, Typography } from '@material-ui/core';
import React from 'react';
import { Link as LinkRouter } from 'react-router-dom';

const useStyles = makeStyles((theme) => ({
    presentation: {
        width: '100%',
        minHeight: '600px',
        background: 'radial-gradient(circle, #1f2b2f 0%, #1c1e20 25%)',
        // position: 'relative',
    },
    hero: {
        background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
        color: 'white',
    },
    newJob: {
        background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
        borderRadius: 3,
        boxShadow: '0 2px 10px 1px rgba(33, 203, 243, .3)',
        color: 'white',
    },
    logo: {
        userSelect: 'none',
        width: 200,
    },
    footer: {
        color: '#758694',
        paddingTop: 30,
        paddingBottom: 30,
        backgroundColor: "#20232a"
    }
}));

export default function Index() {
    const classes = useStyles();
    return (
        <Box>

            <Grid className={classes.presentation} container direction="column" justify="center" spacing={3}>
                <Grid item>
                    <Grid container spacing={6} direction="row" alignItems="center" justify="center">
                        <Grid item>
                            <img src="/assets/varq.svg" alt="VarQ" className={classes.logo} />
                        </Grid>
                        <Grid item xs={9} sm={4} lg={3}>
                            <Grid container direction="column" alignItems="flex-start" justify="center">
                                <Typography variant="h1" align="left" className={classes.name}>VarQ</Typography>
                                <Typography variant="h5" align="left" className={classes.desc}>
                                    A tool for the structural and functional analysis of protein variants.
                            </Typography>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item>
                    <Typography align="center" className={classes.desc}>
                        Start a <LinkRouter to="/new-job"><Button className={classes.newJob}>New Job</Button></LinkRouter> or view <LinkRouter to="/job/fe2423053f1a75a300e4074b1609ed972e3e2eaeae149f21d9c5fd79b4ef3d5c"><Button variant="outlined">Sample Results</Button></LinkRouter>
                    </Typography>
                </Grid>
                <Grid item>
                    <Typography align="center" className={classes.desc}>
                        If you find our work useful, please cite us: <br /> -
                </Typography>
                </Grid>
            </Grid>
            <Grid container className={classes.footer}>
                <Container>
                    <Grid container direction="row" justify="space-between">
                        <Grid item>
                            <Grid container direction="column">
                                <Grid item>
                                    <Typography variant="caption">Bioinformática Estructural y Biofisicoquímica de Proteínas.</Typography>
                                </Grid>
                                <Grid item>
                                    <Typography variant="caption">IQUIBICEN, Departamento de Química Biológica.</Typography>
                                </Grid>
                                <Grid item>
                                    <Typography variant="caption">Facultad de Ciencias Exactas y Naturales, Universidad de Buenos Aires.</Typography>
                                </Grid>
                            </Grid>
                        </Grid>
                        <Grid item>
                            <Grid container direction="column" justify="space-between" alignItems="flex-end">
                                <Grid item>
                                    <Typography variant="caption">VarQ <Link href="#">source code</Link> is released under the <Link href="#">MIT license</Link>.</Typography>
                                </Grid>
                                <Grid item>
                                    <Typography variant="caption">External tools and libraries may have different licenses.</Typography>
                                </Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </Container>
            </Grid>
        </Box>
    )
}
