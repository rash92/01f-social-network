import React from "react";
import {Image} from "react-bootstrap";
import classes from "./User.module.css";
function User({Nickname, IsOnline, Avatar}) {
  return (
    <div style={{display: "flex", alignItems: "center"}}>
      <div style={{position: "relative", marginRight: "10px"}}>
        <Image
          src={`http://localhost:8000/images/${Avatar}`}
          roundedCircle
          style={{marginRight: "5px"}}
          width={50}
          height={50}
        />
        {IsOnline && (
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
        <span className={classes.username}>{Nickname} </span>
      </div>
    </div>
  );
}

export default User;
