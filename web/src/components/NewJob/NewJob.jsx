import { Toolbar, Checkbox, Container, Divider, FormControlLabel, Grid, Typography, Fade, Box, Grow } from '@material-ui/core';
import ArrowForwardIosIcon from '@material-ui/icons/ArrowForwardIos';
import React from 'react';
import { EmailInput } from './EmailInput';
import PDBPicker from './PDBPicker';
import { UniProtInput } from './UniProtInput';
import { QueueInfo } from './QueueInfo';
import SendBar from './SendBar';
import { Variations } from './Variations';
import NavBar from '../NavBar'

export default class NewJob extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            unpData: {},
            pdbs: [],
            variations: [],
            clinvar: false
        }

        this.setUnpData = this.setUnpData.bind(this);
        this.setPDBs = this.setPDBs.bind(this);
        this.setVars = this.setVars.bind(this);
        this.setClinVar = this.setClinVar.bind(this);
    }

    setUnpData(unpData) {
        this.setState({
            unpData: unpData, pdbs: [], variations: [], clinvar: false
        })
    }

    setPDBs(pdbs) {
        this.setState({ pdbs: pdbs })
    }

    setVars(vars) {
        this.setState({ variations: vars })
    }

    setClinVar(flag) {
        this.setState({ clinvar: flag })
    }

    render() {
        let unpOk = Object.keys(this.state.unpData).length > 0;
        let structOk = this.state.unpData.pdbs !== null;
        let dataOk = unpOk && this.state.pdbs.length > 0
            && (this.state.variations.length > 0 || this.state.clinvar);
        return (
            <Box>
                <NavBar />
                <Toolbar />
                <Container>
                    <Typography variant="h2" gutterBottom className="title">New Job</Typography>
                    <Grid container spacing={4} direction="column">
                        <Grid item>
                            <Grid container spacing={2} alignItems="center">
                                <Grid item xs={2}>
                                    <UniProtInput setUnpData={this.setUnpData} />
                                </Grid>
                                <Grid item>
                                    {unpOk && <ArrowForwardIosIcon />}
                                </Grid>
                                <Grid item>
                                    <Grow in={unpOk}>
                                        <Box>
                                            {unpOk &&
                                                <PDBPicker
                                                    unpID={this.state.unpData.id}
                                                    pdbs={this.state.unpData.pdbs}
                                                    setPDBs={this.setPDBs} />}
                                        </Box>
                                    </Grow>
                                </Grid>
                                <Grid item>
                                    {unpOk && structOk && <ArrowForwardIosIcon />}
                                </Grid>
                                <Grid item xs={4}>
                                    <Grow in={unpOk}>
                                        <Box>
                                            {unpOk && structOk &&
                                                <Variations
                                                    unpID={this.state.unpData.id}
                                                    sequence={this.state.unpData.sequence}
                                                    setVariations={this.setVars}
                                                    setClinVar={this.setClinVar} />}
                                        </Box>
                                    </Grow>
                                </Grid>
                            </Grid>
                        </Grid>
                        {dataOk && <Divider />}
                        <Grid item>
                            <Grow in={dataOk}>
                                <Box>
                                    <SendBar />
                                </Box>
                            </Grow>
                        </Grid>
                        {dataOk && <Divider />}
                    </Grid>
                </Container>
            </Box>
        );
    }
}