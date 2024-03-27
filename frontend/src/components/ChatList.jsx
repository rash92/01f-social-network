import React from "react";
import {Button, ListGroup} from "react-bootstrap";
import moment from "moment";
import User from "./User";

export default function Chats({
  avatar,
  nickname,
  message,
  timeSent,
  isOnline,
  opeChat,
  type,
}) {
  const formattedTime = moment(timeSent).calendar(null, {
    sameDay: "h:mm A",
    lastWeek: "dddd",
    sameElse: "MM/DD/YYYY",
  });
  const openChatHandler = (e) => {
    e.preventDefault();
    opeChat({type, avatar, nickname, isOnline});
  };

  return (
    <ListGroup.Item>
      <Button
        onClick={openChatHandler}
        style={{background: "none", border: "none"}}
      >
        <div style={{padding: "1rem"}}>
          <User avatar={avatar} username={nickname} isOnline={isOnline} />
        </div>

        {/* <div className="d-flex justify-content-between gap-4 ">
          <span style={{color: "#6e6e6e"}}>{message}</span>
          <span style={{color: "#6e6e6e", fontSize: "0.8rem"}}>
            {formattedTime}
          </span>
        </div> */}
      </Button>
    </ListGroup.Item>
  );
}
