import { Grid, Box, Chip, TextField, IconButton, Select, MenuItem } from "@material-ui/core";
import { Add, FilterTiltShiftSharp } from '@material-ui/icons';
import React from 'react';
import ChipArray from "./ChipArray";

export default class VariantInput extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            variants: [],
            aaValue: '', posError: false, posErrorMsg: '', aaError: false, aaErrorMsg: '',
            pos: '', aa: ''
        }

        this.handleAaChange = this.handleAaChange.bind(this)
        this.handlePosChange = this.handlePosChange.bind(this);
        this.handleAdd = this.handleAdd.bind(this);
        this.handleKey = this.handleKey.bind(this);
        this.handleDelete = this.handleDelete.bind(this);
    }

    handleDelete(chip) {
        this.setState({
            variants: this.state.variants.filter((c) => c.key !== chip.key),
        }, () => {
            this.props.setVariations(this.state.variants);
        });
    }

    handleAdd() {
        if (this.state.aa != '' && this.state.pos > 0
            && !this.state.posError && !this.state.aaError) {
            let variants = this.state.variants;
            variants.push({
                key: this.state.variants.length,
                pos: this.state.pos,
                aa: this.state.aa,
                label: this.state.pos + ' ' + this.props.sequence[this.state.pos - 1] + '→' + this.state.aa
            })
            this.setState({ variants: variants, pos: '', aa: '' });
            this.props.setVariations(variants);
        }
    }

    handlePosChange(e) {
        let pos = parseInt(e.target.value);
        if (isNaN(pos)) { pos = '' }
        this.setState({ pos: pos }, () => {
            this.checkFields()
        });
    }

    handleAaChange(e) {
        this.setState({ aa: e.target.value.toUpperCase() }, () => {
            this.checkFields()
        });
    }

    handleKey(e) {
        if (e.key === 'Enter') { this.handleAdd() }
    }

    checkFields() {
        const aminoacids = 'ARNDBCEQZGHILKMFPSTWYV';
        let aa = this.state.aa;
        let pos = this.state.pos;
        (pos > this.props.sequence.length || pos < 1) ?
            this.setState({ posError: true, posErrorMsg: 'Outside sequence length' }) :
            this.setState({ posError: false, posErrorMsg: '' });
        (aa.length > 1 || !aminoacids.includes(aa)) ?
            this.setState({ aaError: true, aaErrorMsg: 'Unknown' }) :
            this.setState({ aaError: false, aaErrorMsg: '' });
        if (this.props.sequence[pos - 1] == aa) {
            this.setState({
                posError: true, posErrorMsg: 'Silent mutation',
                aaError: true, aaErrorMsg: 'Silent mutation'
            })
        }
        if (this.state.variants.filter(v => v.pos == pos).length > 0) {
            this.setState({ posError: true, posErrorMsg: 'Position already added' })
        }
    }

    render() {
        return (
            <Box>
                <Grid container spacing={1} alignItems="center">
                    <Grid item xs={7}>
                        <TextField label="Position" onChange={this.handlePosChange} type="number"
                            error={this.state.posError} helperText={this.state.posErrorMsg}
                            value={this.state.pos}
                            onKeyPress={this.handleKey} />
                    </Grid>
                    <Grid item xs={3}>
                        <TextField onChange={this.handleAaChange} label="1-letter Aa"
                            error={this.state.aaError} helperText={this.state.aaErrorMsg}
                            value={this.state.aa}
                            onKeyPress={this.handleKey} />
                    </Grid>
                    <Grid item xs={2}>
                        <IconButton color="primary" aria-label="add variant" onClick={this.handleAdd} style={{ marginTop: 20 }}>
                            <Add />
                        </IconButton>
                    </Grid>
                </Grid>
                <ChipArray variants={this.state.variants} handleDelete={this.handleDelete} />
            </Box>
        );
    }
}