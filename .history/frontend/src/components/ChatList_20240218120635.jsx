import React, {useContext} from "react";
import {Button, ListGroup} from "react-bootstrap";
import moment from "moment";
import User from "./User";
import AuthContext from "../store/authContext";

export default function Chats({avatar, username, message, timeSent, isOnline,openChatHandler }) {
  const formattedTime = moment(timeSent).calendar(null, {
    sameDay: "h:mm A",
    lastWeek: "dddd",
    sameElse: "MM/DD/YYYY",
  });
  
 

  return (
    <ListGroup.Item>
      <Button
        onClick={openChatHandler.bind(null, {avatar, username, isOnline})}
      >
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