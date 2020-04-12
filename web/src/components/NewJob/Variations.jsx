import { Box, Checkbox, FormControlLabel, Typography, Link } from '@material-ui/core';
import React from 'react';
import VariantInput from './VariantInput';

export class Variations extends React.Component {
    constructor(props) {
        super(props);
        this.handleChange = this.handleChange.bind(this);
    }
    handleChange(e) {
        this.props.setClinVar(e.target.checked)
    }
    render() {
        let unpSeqURL = "https://www.uniprot.org/uniprot/" + this.props.unpID + ".fasta";
        return (
            <Box>
                <Typography variant="h5" gutterBottom>3. Add variations</Typography>
                <Typography variant="overline" gutterBottom><Link href={unpSeqURL} target="_blank" rel="noreferrer">sequence</Link> length: {this.props.sequence.length}
                </Typography>
                <Box>
                    <FormControlLabel
                        control={<Checkbox onChange={this.handleChange} />}
                        label="Include ClinVar variants"
                    />
                </Box>

                <VariantInput sequence={this.props.sequence}
                    setVariations={this.props.setVariations} />
            </Box>
        )
    }
}