import TextField from '@material-ui/core/TextField';
import React from 'react';

export class EmailInput extends React.Component {
    render() {
        return (
            <TextField
                margin="dense"
                id="name"
                label="Email address (optional)"
                type="email"
                fullWidth
            />
        )
    }
}