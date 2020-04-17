import { Dialog, DialogContent, DialogContentText, DialogTitle, Grow, LinearProgress } from "@material-ui/core";
import React from 'react';

export default class StatusConsole extends React.Component {
    constructor(props) {
        super(props);
        this.ws = new WebSocket('ws://ip169.qb.fcen.uba.ar:8888/ws/' + this.props.jobID);

        this.state = { messages: [], connected: false }
    }

    componentDidMount() {
        this.ws.onopen = () => {
            this.setState({
                connected: true
            })
        }
        this.ws.onmessage = evt => {
            if (evt.data == "SUCCESS" || evt.data == "FAILED") { this.ws.close() }
            if (evt.data == "SUCCESS") {
                this.props.reload();
            }
            this.setState({
                messages: this.state.messages.concat(evt.data)
            })
        }
        this.ws.onerror = () => {
            this.props.reload();
        }
    }

    render() {
        if (!this.state.connected) {
            return (
                <Dialog open={true} maxWidth="sm" fullWidth={true}>
                    <DialogTitle>Connecting...</DialogTitle>
                    <LinearProgress variant="query" />
                </Dialog>
            )
        }
        return (
            <Dialog open={true} maxWidth="sm" fullWidth={true}>
                <DialogTitle>Your job is pending</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        This is the real time status:
                        </DialogContentText>
                    <div className="status-console">
                        <div className="contents">
                            {this.state.messages.map((msg, index) => {
                                return (
                                    <Grow in={true} key={index}>
                                        <div>{msg}</div>
                                    </Grow>
                                );
                            })}
                        </div>
                    </div>
                </DialogContent>
            </Dialog>
        )
    }
}