import { Button, Checkbox, CircularProgress, Container, Divider, FormControlLabel, Grid, InputAdornment, Typography } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import TextField from '@material-ui/core/TextField';
import ArrowForwardIosIcon from '@material-ui/icons/ArrowForwardIos';
import axios from 'axios';
import React from 'react';
import PDBPicker from './PDBPicker';

const useStyles = makeStyles((theme) => ({
    root: {

    },
    sendJob: {
        background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
        borderRadius: 3,
        boxShadow: '0 2px 10px 1px rgba(33, 203, 243, .3)',
        color: 'white',
    },
    buttons: {
        display: 'flex',
        justifyContent: 'flex-end',
    }
}));


export default function NewJob(props) {
    const classes = useStyles();
    const handleChange = (e) => {
        if (e.target.value.length > 5) {
            axios.get('http://127.0.0.1:3000/api/uniprot/' + e.target.value)
                .then(function (response) {
                    console.log(response);
                }).catch(function (error) {
                    console.log(error);
                });
            props.sarasa = "ok";
        }
    };
    return (
        <Container className={classes.root}>



            <Container>
                <Typography variant="h2" gutterBottom>New Job</Typography>

                <Grid container spacing={4} direction="column">
                    <Grid item>
                        <Grid container spacing={2} alignItems="center">
                            <Grid item>
                                <Typography variant="h5" gutterBottom>1. Enter a protein</Typography>
                                <div>
                                    <TextField id="filled-basic" label="UniProt Accession ID" variant="filled" autoFocus
                                        onChange={handleChange}
                                        InputProps={{
                                            endAdornment: <InputAdornment position="end"><CircularProgress /></InputAdornment>,
                                        }} />
                                </div>
                                <Typography variant="overline" gutterBottom>NRAS - GTPase NRas - Homo sapiens (Human)</Typography>
                            </Grid>
                            <Grid item>
                                <ArrowForwardIosIcon />
                            </Grid>
                            <Grid item>
                                <Typography variant="h5" gutterBottom>2. Choose structures</Typography>
                                <PDBPicker />
                            </Grid>
                            <Grid item>
                                <ArrowForwardIosIcon />
                            </Grid>
                            <Grid item>
                                <Typography variant="h5" gutterBottom>3. Add variations</Typography>
                                <FormControlLabel
                                    control={<Checkbox />}
                                    label="Include ClinVar variants"
                                />
                            </Grid>
                        </Grid>
                    </Grid>

                    <Divider />

                    <Grid item>
                        <Grid container spacing={2} alignItems="center">
                            <Grid item xs>
                                <TextField
                                    margin="dense"
                                    id="name"
                                    label="Email address (optional)"
                                    type="email"
                                    value={props.sarasa}
                                    fullWidth
                                />
                            </Grid>
                            <Grid item xs>
                                <Grid container spacing={2}>
                                    <Grid container item>
                                        <Grid item xs={8}>Jobs in queue:</Grid>
                                        <Grid item>0</Grid>
                                    </Grid>
                                    <Grid container item>
                                        <Grid item xs={8}>Estimated time:</Grid>
                                        <Grid item>{"<"}1 minute</Grid>
                                    </Grid>
                                </Grid>
                            </Grid>

                            <Grid item xs={2} className={classes.buttons}>
                                <Button className={classes.sendJob}>Send Job</Button>
                            </Grid>
                        </Grid>
                    </Grid>

                    <Divider />
                </Grid>



            </Container>








        </Container>
    );
}