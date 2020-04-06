import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import AppBar from '@material-ui/core/AppBar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import Toolbar from '@material-ui/core/Toolbar';

const useStyles = makeStyles((theme) => ({
    root: {
        flexGrow: 1,
    },
    bar: {
        backgroundColor: "#20232a",
    },
    logo: {
        marginRight: theme.spacing(1.5),
        width: 50,
        height: 50,
    },
    myJobs: {
        marginRight: theme.spacing(1.5),
    },
    newJob: {
        background: 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
        borderRadius: 3,
        boxShadow: '0 2px 10px 1px rgba(33, 203, 243, .3)',
        color: 'white',
    },
    title: {
        flexGrow: 1,
    },
}));

export default function NavBar() {
    const classes = useStyles();

    return (
        <div className={classes.root}>
            <AppBar className={classes.bar}>
                <Toolbar>
                    <img className={classes.logo} src="assets/varq.svg" alt="" />
                    <Typography variant="h6" className={classes.title}>
                        VarQ
                    </Typography>
                    <Button className={classes.myJobs} variant="outlined" color="inherit">My Jobs</Button>
                    <Button className={classes.newJob}>New Job</Button>
                </Toolbar>
            </AppBar>
        </div >
    )
}