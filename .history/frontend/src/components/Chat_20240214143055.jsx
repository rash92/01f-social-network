import React from "react";
import {ListGroup, Image} from "react-bootstrap";
import moment from "moment";
import RoundImage from "./RoundedImg";

export default function Chats({avatar, username, message, timeSent}) {
  const formattedTime = moment(timeSent).format("h:mm A");

  return (
    <ListGroup.Item>
      <RoundImage src={avatar} alt="Avatar" size={30} />
      <span>{username}</span>
      <span>{message}</span>
      <span>{formattedTime}</span>
    </ListGroup.Item>
  );
}