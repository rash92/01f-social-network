import React from "react";
import {Image} from "react-bootstrap";

function User({name, isLoggedIn}) {
  return (
    <div style={{display: "flex", alignItems: "center" }}>
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
      <div>
        <div>{name}</div>
      </div>
    </div>
  );
}

export default User;