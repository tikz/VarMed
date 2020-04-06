import React from 'react';
import NavBar from './NavBar';
import { Box, Container, Toolbar, Typography, Grid, Paper, createMuiTheme, ThemeProvider, Divider } from '@material-ui/core';
import SequenceViewer from './SequenceViewer';
import StructureViewer from './StructureViewer';

const darkTheme = createMuiTheme({
    palette: {
        type: 'dark',
    },
});

export default class Results extends React.Component {
    componentDidMount() {
    }
    render() {
        // return (<StructureViewer />)
        return (
            <ThemeProvider theme={darkTheme} >
                <Box>
                    <NavBar />
                    <Toolbar />
                    <Container>
                        <Box>
                            <Typography variant="h4">Methylglyoxal synthase</Typography>
                            <Divider />

                            <Grid container spacing={3}>
                                <Grid item xs={5}>
                                    <Paper>xs=12</Paper>
                                </Grid>
                                <Grid item xs={7}>
                                    <Paper>
                                        <StructureViewer />
                                    </Paper>
                                </Grid>
                            </Grid>
                        </Box>
                        <Box my={2}>
                            <SequenceViewer />
                        </Box>

                    </Container>
                </Box>

            </ThemeProvider>
        )
    }

}
