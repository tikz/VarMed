import { Button, Grid, TextField } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';
import { QueueInfo } from './QueueInfo';

const useStyles = makeStyles((theme) => ({
    sendJob: {
        background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
        borderRadius: 3,
        boxShadow: '0 2px 10px 1px rgba(33, 203, 243, .3)',
        color: 'white',
    }
}));

export default function SendBar() {
    const classes = useStyles();
    return (
        <Grid container spacing={2} alignItems="center">
            <Grid item xs>
                <TextField
                    id="name" label="Email address (optional)"
                    margin="dense" type="email" fullWidth
                />
            </Grid>
            <Grid item xs>
                <QueueInfo />
            </Grid>
            <Grid item xs={2}>
                <Button className={classes.sendJob}>Send Job</Button>
            </Grid>
        </Grid>
    )
}
