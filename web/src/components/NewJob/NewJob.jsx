import { Checkbox, Container, Divider, FormControlLabel, Grid, Typography, Fade, Box, Grow } from '@material-ui/core';
import ArrowForwardIosIcon from '@material-ui/icons/ArrowForwardIos';
import React from 'react';
import { EmailInput } from './EmailInput';
import PDBPicker from './PDBPicker';
import { UniProtInput } from './UniProtInput';
import { QueueInfo } from './QueueInfo';
import SendBar from './SendBar';
import { Variations } from './Variations';

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
        let showInputs = Object.keys(this.state.unpData).length > 0;
        let showSend = showInputs && this.state.pdbs.length > 0
            && (this.state.variations.length > 0 || this.state.clinvar);
        return (
            <Container>
                <Typography variant="h2" gutterBottom>New Job</Typography>
                <Grid container spacing={4} direction="column">
                    <Grid item>
                        <Grid container spacing={2} alignItems="center">
                            <Grid item xs={2}>
                                <UniProtInput setUnpData={this.setUnpData} />
                            </Grid>
                            <Grid item>
                                {showInputs && <ArrowForwardIosIcon />}
                            </Grid>
                            <Grid item>
                                <Grow in={showInputs}>
                                    <Box>
                                        {showInputs &&
                                            <PDBPicker
                                                unpID={this.state.unpData.id}
                                                pdbs={this.state.unpData.pdbs}
                                                setPDBs={this.setPDBs} />}
                                    </Box>
                                </Grow>
                            </Grid>
                            <Grid item>
                                {showInputs && <ArrowForwardIosIcon />}
                            </Grid>
                            <Grid item xs={4}>
                                <Grow in={showInputs}>
                                    <Box>
                                        {showInputs &&
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
                    {showSend && <Divider />}
                    <Grid item>
                        <Grow in={showSend}>
                            <Box>
                                <SendBar />
                            </Box>
                        </Grow>
                    </Grid>
                    {showSend && <Divider />}
                </Grid>
            </Container>
        );
    }
}