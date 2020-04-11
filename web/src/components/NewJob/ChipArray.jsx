import { Chip } from "@material-ui/core";
import { makeStyles } from '@material-ui/core/styles';
import React from 'react';



const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
        justifyContent: 'center',
        flexWrap: 'wrap',
        marginTop: theme.spacing(1),
        '& > *': {
            margin: theme.spacing(0.5),
        },
    },
}));

export default function ChipArray(props) {
    const classes = useStyles();

    return (
        <div className={classes.root}>
            {props.variants.map((data) => {
                return (
                    <Chip
                        key={data.key}
                        label={data.label}
                        onDelete={() => props.handleDelete(data)}
                        size="small"
                    />
                );
            })}
        </div>
    );
}