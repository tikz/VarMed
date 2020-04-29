import { AppBar, Button, Toolbar, Typography } from "@material-ui/core";
import React from "react";
import { Link } from "react-router-dom";
import "../styles/components/nav-bar.scss";
import MyJobs from "./MyJobs";

export default function NavBar() {
  return (
    <div className="nav">
      <AppBar className="bar">
        <Toolbar className="bar">
          <Link to="/" className="link">
            <img className="nav-logo" src="/assets/varq.svg" alt="" />
          </Link>
          <Typography variant="h6" className="nav-title">
            VarQ
          </Typography>
          <MyJobs />
          <Link to="/new-job">
            <Button className="glowButton">New Job</Button>
          </Link>
        </Toolbar>
      </AppBar>
    </div>
  );
}
