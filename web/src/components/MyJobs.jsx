import {
  Avatar,
  Button,
  Dialog,
  DialogContent,
  DialogTitle,
  Grid,
  IconButton,
  List,
  ListItem,
  ListItemAvatar,
  ListItemIcon,
  ListItemSecondaryAction,
  Typography,
} from "@material-ui/core";
import DeleteIcon from "@material-ui/icons/Delete";
import NotesIcon from "@material-ui/icons/Notes";
import moment from "moment";
import React from "react";
import { withRouter } from "react-router-dom";

console.log(JOBS_KEY);
function MyJobsDialog(props) {
  const getJobs = () => {
    return JSON.parse(window.localStorage.getItem(JOBS_KEY));
  };

  const { onClose, open } = props;
  const [render, setRender] = React.useState(0);

  const handleClose = () => {
    onClose();
  };

  const handleListItemClick = (jobId) => {
    onClose("/job/" + jobId);
  };

  const deleteJob = (jobId) => {
    let jobsStorage = getJobs();
    window.localStorage.removeItem(JOBS_KEY);
    window.localStorage.setItem(
      JOBS_KEY,
      JSON.stringify(
        jobsStorage.filter((j) => {
          return j.id != jobId;
        })
      )
    );
    setRender(render + 1);
  };

  const jobs = getJobs();
  let jobsList;
  if (jobs === null || jobs.length == 0) {
    jobsList = (
      <DialogContent>
        <Typography variant="button" display="block" gutterBottom>
          No jobs found.
        </Typography>
        <Typography variant="caption" display="block" gutterBottom>
          Start a new job or open a results page, it will show up here.
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
          <ListItemIcon>
            <Grid container direction="column">
              <Grid item>
                <ListItemAvatar>
                  <Avatar>
                    <NotesIcon />
                  </Avatar>
                </ListItemAvatar>
              </Grid>
              <Grid item>
                <Typography variant="caption">{job.id.slice(0, 6)}</Typography>
              </Grid>
            </Grid>
          </ListItemIcon>
          <Grid container>
            <Grid item></Grid>
            <Grid item>
              <Grid container direction="column">
                <Grid item>
                  <Typography variant="button">{job.name}</Typography>
                </Grid>
                <Grid item>
                  <Typography variant="overline">
                    {job.pdbIds.join(", ")}
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

          <ListItemSecondaryAction>
            <IconButton aria-label="delete" onClick={() => deleteJob(job.id)}>
              <DeleteIcon />
            </IconButton>
          </ListItemSecondaryAction>
        </ListItem>
      );
    });
  }

  return (
    <Dialog onClose={handleClose} open={open} maxWidth="xs" fullWidth={true}>
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
