import { Button } from '@material-ui/core';
import { makeStyles } from '@material-ui/core/styles';
import React, { PropTypes } from 'react';

const useStyles = makeStyles((theme) => ({
    sendJob: {
        background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
        borderRadius: 3,
        boxShadow: '0 2px 10px 1px rgba(33, 203, 243, .3)',
        color: 'white',
    }
}));

export default function SendButton(props) {
    const classes = useStyles();
    return (
        <Button className={classes.sendJob} onClick={props.submit}>Send Job</Button>
    )
}
