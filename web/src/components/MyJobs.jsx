import { DialogContent, Grid } from "@material-ui/core";
import Avatar from "@material-ui/core/Avatar";
import Button from "@material-ui/core/Button";
import Dialog from "@material-ui/core/Dialog";
import DialogTitle from "@material-ui/core/DialogTitle";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemAvatar from "@material-ui/core/ListItemAvatar";
import Typography from "@material-ui/core/Typography";
import NotesIcon from "@material-ui/icons/Notes";
import moment from "moment";
import React from "react";
import { withRouter } from "react-router-dom";

function MyJobsDialog(props) {
  const { onClose, open } = props;

  const handleClose = () => {
    onClose();
  };

  const handleListItemClick = (jobID) => {
    onClose("/job/" + jobID);
  };

  let jobs = JSON.parse(window.localStorage.getItem("jobs"));
  var jobsList;
  if (jobs === null) {
    jobsList = (
      <DialogContent>
        <Typography variant="button" display="block" gutterBottom>
          No jobs found.
        </Typography>
        <Typography variant="caption" display="block" gutterBottom>
          Start a new job or open a results page and it will show up here.
        </Typography>
      </DialogContent>
    );
  } else {
    jobsList = jobs.map((job) => {
      return (
        <ListItem
          button
          onClick={() => handleListItemClick(job.id)}
          key={job.id}
        >
          <Grid container>
            <Grid item>
              <Grid container direction="column">
                <Grid item>
                  <ListItemAvatar>
                    <Avatar>
                      <NotesIcon />
                    </Avatar>
                  </ListItemAvatar>
                </Grid>
                <Grid item>
                  <Typography variant="caption">
                    {job.id.slice(0, 6)}
                  </Typography>
                </Grid>
              </Grid>
            </Grid>
            <Grid item>
              <Grid container direction="column">
                <Grid item>
                  <Typography variant="button">{job.name}</Typography>
                </Grid>
                <Grid item>
                  <Typography variant="overline">
                    {job.pdbs.join(", ")}
                  </Typography>
                </Grid>
                <Grid item>
                  <Typography variant="caption">
                    {moment(job.date).fromNow()}
                  </Typography>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </ListItem>
      );
    });
  }

  return (
    <Dialog onClose={handleClose} open={open}>
      <DialogTitle>My jobs</DialogTitle>
      <List>{jobsList}</List>
    </Dialog>
  );
}

function MyJobs(props) {
  const [open, setOpen] = React.useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = (redirectURL) => {
    if (redirectURL) {
      props.history.push(redirectURL);
    }
    setOpen(false);
  };

  return (
    <div>
      <Button
        className="myJobs"
        variant="outlined"
        color="inherit"
        onClick={handleClickOpen}
      >
        My Jobs
      </Button>
      <MyJobsDialog open={open} onClose={handleClose} />
    </div>
  );
}

export default withRouter(MyJobs);
