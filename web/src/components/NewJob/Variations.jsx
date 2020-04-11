import { Box, Checkbox, FormControlLabel, Typography } from '@material-ui/core';
import React from 'react';


export class Variations extends React.Component {
    render() {
        return (
            <Box>
                <Typography variant="h5" gutterBottom>3. Add variations</Typography>
                <FormControlLabel
                    control={<Checkbox />}
                    label="Include ClinVar variants"
                />
            </Box>
        )
    }
}