import { Box, Container, Divider, FormControl, Grid, InputLabel, MenuItem, Paper, Select, Toolbar, Typography } from '@material-ui/core';
import React from 'react';
import NavBar from '../NavBar';
import SequenceViewer from './SequenceViewer';
import StructureViewer from './StructureViewer';

export default class Results extends React.Component {
    constructor(props) {
        super(props);
        this.structureRef = React.createRef();

        this.state = { pdb: "3CON" }

        this.highlightStructure = this.highlightStructure.bind(this);
        this.selectStructure = this.selectStructure.bind(this);
        this.pdbChange = this.pdbChange.bind(this);
    }

    componentDidMount() {
    }

    highlightStructure(start, end) {
        if (start == 0 && end == 0) {
            this.structureRef.current.clearHighlight();
        } else {
            this.structureRef.current.highlight(start + 18, end + 18);
        }
    }

    selectStructure(start, end) {
        this.structureRef.current.focus(start + 18, end + 18);
        this.structureRef.current.highlight(start + 18, end + 18);
        if (start - end == 0) {
            this.structureRef.current.select(start + 18, end + 18);
        }
    }

    pdbChange(e) {
        this.setState({ pdb: e.target.value })
    }

    render() {
        return (
            <Box>
                <NavBar />
                <Toolbar />
                <Container>
                    <Box className="over">
                        <Typography variant="h4" className="title">GTPase NRas</Typography>
                        <Divider />
                        <Grid container spacing={2} alignItems="center">
                            <Grid item>
                                <FormControl variant="outlined" >
                                    <InputLabel>PDB</InputLabel>
                                    <Select
                                        label="PDB"
                                        value={this.state.pdb}
                                        onChange={this.pdbChange}
                                    >
                                        <MenuItem value={"3CON"}>3CON</MenuItem>
                                        <MenuItem value={"123X"}>123X</MenuItem>
                                    </Select>
                                </FormControl>
                            </Grid>
                            <Grid item>
                                <Typography>Crystal structure of the human NRAS GTPase bound with GDP</Typography>
                            </Grid>
                        </Grid>
                    </Box>
                </Container>

                <StructureViewer ref={this.structureRef} />

                <Container>
                    <Box>
                        <Grid container spacing={3}>
                            <Grid item xs={5}>

                            </Grid>
                            <Grid item xs={7}>
                                <Paper>

                                </Paper>
                            </Grid>
                        </Grid>
                    </Box>
                    <Box my={2}>
                        <SequenceViewer highlight={this.highlightStructure} select={this.selectStructure} />
                    </Box>
                </Container>
            </Box>
        )
    }

}