import React from "react";
import { Link as LinkRouter } from "react-router-dom";
import {
  Grid,
  Typography,
  Box,
  Badge,
  IconButton,
  Popover,
  LinearProgress,
  Divider,
} from "@material-ui/core";
import DynamicFeed from "@material-ui/icons/DynamicFeed";
import { withStyles } from "@material-ui/core/styles";
import "../styles/components/queue-status.scss";
import { MoreVert } from "@material-ui/icons";

const StyledBadge = withStyles((theme) => ({
  badge: {
    right: -3,
    top: 13,
    border: `2px solid ${theme.palette.background.paper}`,
    padding: "0 4px",
  },
}))(Badge);

export default class QueueIcon extends React.Component {
  constructor(props) {
    super(props);

    this.state = { status: {}, anchorEl: null };

    this.handleClick = this.handleClick.bind(this);
    this.handleClose = this.handleClose.bind(this);
  }

  connectWS() {
    var url;
    if (API_URL == "") {
      let protocol = window.location.protocol === "https:" ? "wss://" : "ws://";
      url = protocol + window.location.host + "/ws/queue";
    } else {
      url = API_URL.replace("http", "ws") + "/ws/queue";
    }

    this.ws = new WebSocket(url);

    this.ws.onmessage = (evt) => {
      const status = JSON.parse(evt.data);
      this.setState({ status: status });
    };

    this.ws.onclose = () => {
      setTimeout(function () {
        connectWS();
      }, 1000);
    };

    this.ws.onerror = function (err) {
      ws.close();
    };
  }

  componentDidMount() {
    this.connectWS();
  }

  handleClick(e) {
    this.setState({ anchorEl: e.currentTarget });
  }

  handleClose() {
    this.setState({ anchorEl: null });
  }

  render() {
    const status = this.state.status;

    const isProcessing = function (pos) {
      return status.myJobs
        ? status.jobs.filter((j) => j.position == pos).length != 0
        : [];
    };

    const isMyJob = function (shortId) {
      return status.myJobs
        ? status.myJobs.filter((j) => j.shortId == shortId).length != 0
        : [];
    };

    const myJobId = function (shortId) {
      var id;
      if (status.myJobs) {
        status.myJobs.forEach((j) => {
          if (j.shortId == shortId) {
            id = j.id;
          }
        });
      }
      return id;
    };

    return (
      <Box>
        {status.myJobs && <LinearProgress className="queue-icon" />}

        <IconButton
          aria-label="queue"
          onClick={this.handleClick}
          className="queue-icon"
        >
          <StyledBadge badgeContent={status.totalJobs} color="primary">
            <DynamicFeed />
          </StyledBadge>
        </IconButton>
        <Popover
          open={Boolean(this.state.anchorEl)}
          anchorEl={this.state.anchorEl}
          onClose={this.handleClose}
          anchorOrigin={{
            vertical: "bottom",
            horizontal: "center",
          }}
          transformOrigin={{
            vertical: "top",
            horizontal: "center",
          }}
        >
          <Box className="queue-status">
            <Typography variant="h5">Queue</Typography>
            <Divider />

            <Grid
              container
              direction="column"
              spacing={2}
              className="queue-job"
            >
              {status.jobs &&
                status.jobs.map((j) => {
                  return (
                    <Grid
                      item
                      key={j.shortId}
                      container
                      alignItems="center"
                      justify="center"
                      spacing={3}
                    >
                      <Grid item>
                        <Typography variant="h6">#{j.position}</Typography>
                      </Grid>
                      <Grid item xs={4}>
                        {isMyJob(j.shortId) && (
                          <LinkRouter to={"job/" + myJobId(j.shortId)}>
                            <Typography variant="overline">
                              {j.shortId} (my job)
                            </Typography>
                          </LinkRouter>
                        )}

                        {!isMyJob(j.shortId) && (
                          <Typography variant="overline">
                            {j.shortId}
                          </Typography>
                        )}

                        <br />
                        <Typography variant="caption">
                          {j.pdbs} PDBs, {j.variants} variants
                        </Typography>
                      </Grid>
                      <Grid item xs>
                        <Grid container direction="column" justify="center">
                          <Grid item>
                            <LinearProgress
                              variant="buffer"
                              value={j.progress * 100}
                              valueBuffer={j.progressPdb * 100}
                              className="progress"
                            />
                          </Grid>
                          <Grid item>
                            <center>
                              <Typography variant="caption">
                                {j.elapsed} elapsed
                              </Typography>
                            </center>
                          </Grid>
                        </Grid>
                      </Grid>
                    </Grid>
                  );
                })}

              {status.myJobs && status.totalJobs > status.myJobs.length && (
                <Grid item>
                  <center>
                    <MoreVert />
                  </center>
                </Grid>
              )}

              {status.myJobs &&
                status.myJobs.map((j) => {
                  if (isProcessing(j.position)) {
                    return;
                  }
                  return (
                    <Grid
                      item
                      key={j.shortId}
                      container
                      alignItems="center"
                      justify="center"
                      spacing={3}
                    >
                      <Grid item>
                        <Typography variant="h6">#{j.position}</Typography>
                      </Grid>
                      <Grid item xs={4}>
                        <LinkRouter to={"job/" + myJobId(j.id)}>
                          <Typography variant="overline">
                            {j.shortId} (my job)
                          </Typography>
                        </LinkRouter>
                        <br />
                        <Typography variant="caption">
                          {j.pdbs} PDBs, {j.variants} variants
                        </Typography>
                      </Grid>
                      <Grid item xs>
                        <Grid container direction="column" justify="center">
                          <Grid item>
                            <center>
                              <Typography variant="caption">
                                Waiting in queue.
                              </Typography>
                            </center>
                          </Grid>
                        </Grid>
                      </Grid>
                    </Grid>
                  );
                })}

              {status.myJobs === null &&
                status.jobs &&
                status.totalJobs > status.jobs.length && (
                  <Grid item>
                    <center>
                      <MoreVert />
                    </center>
                  </Grid>
                )}

              {status.myJobs === null &&
                status.jobs &&
                status.totalJobs > status.jobs.length && (
                  <Grid item>
                    <center>
                      <Typography variant="caption">
                        +{status.totalJobs - status.jobs.length} jobs waiting
                      </Typography>
                    </center>
                  </Grid>
                )}

              {status.totalJobs == 0 && (
                <Grid item>
                  <center>
                    <Typography variant="h6">No jobs in queue.</Typography>
                  </center>
                </Grid>
              )}
            </Grid>
          </Box>
        </Popover>
      </Box>
    );
  }
}
