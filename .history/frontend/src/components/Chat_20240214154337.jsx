import React from "react";
import {ListGroup} from "react-bootstrap";
import moment from "moment";
import RoundImage from "./RoundedImg";

export default function Chats({avatar, username, message, timeSent}) {
  const formattedTime = moment(timeSent).calendar(null, {
    sameDay: "h:mm A",
    lastWeek: "dddd",
    sameElse: "MM/DD/YYYY",
  });

  return (
    <ListGroup.Item>
      <div d-flex justify-content-between align-items-center >
        <RoundImage src={avatar} alt="Avatar" size={40} />
        <span>{username}</span>
      </div>

      <span>{message}</span>
      <span>{formattedTime}</span>
    </ListGroup.Item>
  );
}