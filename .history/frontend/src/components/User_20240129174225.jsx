import React from "react";
import {Image} from "react-bootstrap";
import classes from "./User.module.css";
function User({name, isLoggedIn, d{}}) {
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
              width: "20px",
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
        <span className="name">{name}</span>
        <span className="userName">{userName} </span>
      </div>
    </div>
  );
}

export default User;