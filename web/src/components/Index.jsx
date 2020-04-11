import { Container, Box, Typography, makeStyles, Grid, Button } from '@material-ui/core';
import React from 'react';
import { Link } from "react-router-dom";

const useStyles = makeStyles((theme) => ({
    presentation: {
        width: '100%',
        minHeight: '600px',
        position: 'relative',
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
}));

export default function Index() {
    const classes = useStyles();
    return (
        <Grid className={classes.presentation} container direction="column" justify="center" spacing={3}>
            <Grid item>
                <Grid container spacing={6} direction="row" alignItems="center" justify="center">
                    <Grid item>
                        <img src="assets/varq.svg" alt="VarQ" className={classes.logo} />
                    </Grid>
                    <Grid item xs={9} sm={3}>
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
                    Start a <Button className={classes.newJob}>New Job</Button> or view <Button variant="outlined">Sample Results</Button>
                </Typography>
            </Grid>
            <Grid item>
                <Typography align="center" className={classes.desc}>
                    If you find our work useful, please cite:
                </Typography>
            </Grid>
        </Grid>

        // <Container>
        //     <h1>Index</h1>
        //     <p>test views:</p>
        //     <Link to="/results">results</Link>
        //     <br />
        //     <Link to="/new-job">new job</Link>
        // </Container>
    )
}
