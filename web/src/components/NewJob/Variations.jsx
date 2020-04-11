import { Box, Checkbox, FormControlLabel, Typography } from '@material-ui/core';
import React from 'react';
import VariantInput from './VariantInput';


export class Variations extends React.Component {
    render() {
        return (
            <Box>
                <Typography variant="h5" gutterBottom>3. Add variations</Typography>
                <Typography variant="overline" gutterBottom>
                    UniProt sequence length: {this.props.sequence.length}
                </Typography>
                <FormControlLabel
                    control={<Checkbox />}
                    label="Include ClinVar variants"
                />
                <VariantInput sequence={this.props.sequence} />
            </Box>
        )
    }
}