import React from "react";
import {Image} from "react-bootstrap";
import classes from "./User.module.css";
function User({name, isLoggedIn, userName}) {
  return (
    <div style={{display: "flex", alignItems: "center"}}>
      <div style={{position: "relative", marginRight: "10px"}}>
        <Image
          src="https://via.placeholder.com/50"
          roundedCircle
          style={{marginRight: "5px"}}
        />
        {isLoggedIn && (
          <div
            style={{
              position: "absolute",
              width: "15px",
              height: "15px",
              backgroundColor: "green",
              borderRadius: "50%",
              bottom: "6px",
              right: "6px",
              border: "2px solid white",
            }}
          />
        )}
      </div>
      <div className={classes.userinfo}>
        <span className={classes.name}>{name}</span>
        <span className={classes.username}>{userName} </span>
      </div>
    </div>
  );
}

export default User;
