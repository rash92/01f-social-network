import React from "react";
import {Button, ListGroup} from "react-bootstrap";
import moment from "moment";
import User from "./User";

export default function Chats({
  avatar,
  nickname,
  timeSent,
  isOnline,
  opeChat,
  id,
}) {
  return (
    <ListGroup.Item>
      <Button
        onClick={opeChat.bind(null, {
          id,
          Nickname: nickname,
          Avatar: avatar,
          type: "privateMessage",
        })}
        style={{background: "none", border: "none"}}
      >
        <div style={{padding: "1rem"}}>
          <User Avatar={avatar} Nickname={nickname} IsOnline={isOnline} />
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
