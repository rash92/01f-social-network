import React from "react";
import {Image} from "react-bootstrap";
import classes from "./User.module.css";
function User({isOnline, userName, avatar}) {
  return (
    <div style={{display: "flex", alignItems: "center"}}>
      <div style={{position: "relative", marginRight: "10px"}}>
        <Image
          src={`http://localhost:8000/images/${avatar}`}
          roundedCircle
          style={{marginRight: "5px"}}
          width={50}
          height={50}
        />
        {isOnline && (
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
        <span className={classes.username}>{userName} </span>
      </div>
    </div>
  );
}

export default User;
