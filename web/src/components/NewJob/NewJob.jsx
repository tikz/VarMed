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
            unpData: {}
        }

        this.setUnpData = this.setUnpData.bind(this);
    }

    setUnpData(unpData) {
        this.setState({
            unpData: unpData
        })
    }

    render() {
        let showInputs = Object.keys(this.state.unpData).length > 0;
        return (
            <Container>
                <Typography variant="h2" gutterBottom>New Job</Typography>
                <Grid container spacing={4} direction="column">
                    <Grid item>
                        <Grid container spacing={2} alignItems="center">
                            <Grid item>
                                <UniProtInput setUnpData={this.setUnpData} />
                            </Grid>
                            <Grid item>
                                {showInputs && <ArrowForwardIosIcon />}
                            </Grid>
                            <Grid item>
                                <Grow in={showInputs}>
                                    <Box>
                                        {showInputs &&
                                            <PDBPicker pdbs={this.state.unpData.pdbs} />}
                                    </Box>
                                </Grow>
                            </Grid>
                            <Grid item>
                                {showInputs && <ArrowForwardIosIcon />}
                            </Grid>
                            <Grid item>
                                <Grow in={showInputs}>
                                    <Box>
                                        {showInputs &&
                                            <Variations />}
                                    </Box>
                                </Grow>
                            </Grid>
                        </Grid>
                    </Grid>


                    {showInputs && <Divider />}

                    <Grid item>
                        <Grow in={showInputs}>
                            <Box>
                                {showInputs &&
                                    <SendBar />}
                            </Box>
                        </Grow>
                    </Grid>

                    {showInputs && <Divider />}

                </Grid>
            </Container>
        );
    }
}