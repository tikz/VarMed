import React from 'react';
import NavBar from './NavBar';
import { Box, Container, Toolbar, Typography, Grid, Paper, createMuiTheme, ThemeProvider, Divider } from '@material-ui/core';
import { FormControl, InputLabel, Select, MenuItem } from '@material-ui/core'
import SequenceViewer from './SequenceViewer';
import StructureViewer from './StructureViewer';

const darkTheme = createMuiTheme({
    palette: {
        type: 'dark',
    },
});

export default class Results extends React.Component {
    constructor(props) {
        super(props);
        this.structureRef = React.createRef();
        this.highlightStructure = this.highlightStructure.bind(this);
        this.selectStructure = this.selectStructure.bind(this);
    }
    componentDidMount() {
        console.log(this.structureRef)
    }
    highlightStructure(start, end) {
        if (start == 0 && end == 0) {
            this.structureRef.current.clearHighlight();
        } else {
            this.structureRef.current.highlight(start + 17, end + 17);
        }
    }
    selectStructure(start, end) {
        this.structureRef.current.select(start + 17, end + 17);
    }
    render() {
        // return (<StructureViewer />)
        return (
            <ThemeProvider theme={darkTheme} >
                <Box>
                    <NavBar />
                    <Toolbar />
                    <Container>
                        <Box className="over">
                            <Typography variant="h4">GTPase NRas</Typography>
                            <Divider />
                            <Grid container spacing={2} alignItems="center">
                                <Grid item>
                                    <FormControl variant="outlined" >
                                        <InputLabel>PDB</InputLabel>
                                        <Select
                                            label="PDB"
                                        >
                                            <MenuItem value={10}>IZ1N</MenuItem>
                                            <MenuItem value={20}>3CON</MenuItem>
                                            <MenuItem value={30}>1A93</MenuItem>
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

            </ThemeProvider>
        )
    }

}
