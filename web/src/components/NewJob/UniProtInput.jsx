import { Box, CircularProgress, InputAdornment, Typography } from "@material-ui/core";
import TextField from '@material-ui/core/TextField';
import axios from 'axios';
import React from 'react';


export class UniProtInput extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            unpID: "",
            unpData: {},
            loading: false,
            error: false,
            errorMsg: "",
            entryName: "",
        }

        this.handleChange = this.handleChange.bind(this);
    };

    handleChange(e) {
        this.props.setUnpData({})
        this.setState({
            unpID: e.target.value.toUpperCase(),
            error: false,
            errorMsg: "",
            entryName: "",
        })

        if (e.target.value.length > 5) {
            let p = this
            p.setState({
                loading: true,
            })
            axios.get('/api/uniprot/' + e.target.value)
                .then(function (response) {
                    let data = response.data
                    p.setState({
                        entryName: [data.gene, data.name, data.organism].join(" - "),
                        unpData: data
                    })
                    p.props.setUnpData(data)
                }).catch(function (error) {
                    p.setState({
                        error: true,
                        errorMsg: "Network error"
                    })
                }).then(function () {
                    p.setState({
                        loading: false
                    })
                });
        } else {
            this.setState({
                error: false,
                errorMsg: "",
                entryName: ""
            })
        }
    };

    render() {
        return (
            <Box>
                <Typography variant="h5" gutterBottom>1. Enter a protein</Typography>
                <div>
                    <TextField id="filled-basic"
                        label="UniProt Accession ID"
                        variant="filled"
                        autoFocus
                        autoComplete="off"
                        fullWidth
                        onChange={this.handleChange}
                        value={this.state.unpID}
                        error={this.state.error}
                        helperText={this.state.errorMsg}
                        InputProps={{
                            endAdornment: <InputAdornment position="end">
                                {this.state.loading && <CircularProgress />}
                            </InputAdornment>,
                        }} />
                </div>
                <Typography variant="overline" gutterBottom>{this.state.entryName}</Typography>
            </Box>
        )
    }
}