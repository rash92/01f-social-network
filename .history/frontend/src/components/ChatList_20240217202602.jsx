import React from "react";
import {Button, ListGroup} from "react-bootstrap";
import moment from "moment";
import User from "./User";

export default function Chats({avatar, username, message, timeSent, isOnline}) {
  const formattedTime = moment(timeSent).calendar(null, {
    sameDay: "h:mm A",
    lastWeek: "dddd",
    sameElse: "MM/DD/YYYY",
  });
  open()

  return (
    <ListGroup.Item>
      <Button onClick={}>
        <div style={{padding: "1rem"}}>
          <User avatar={avatar} userName={username} isOnline={isOnline} />
        </div>

        <div className="d-flex justify-content-between gap-4 ">
          <span style={{color: "#6e6e6e"}}>{message}</span>
          <span style={{color: "#6e6e6e", fontSize: "0.8rem"}}>
            {formattedTime}
          </span>
        </div>
      </Button>
    </ListGroup.Item>
  );
}