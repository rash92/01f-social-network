import React from "react";
import {ListGroup} from "react-bootstrap";
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
import React from "react";
import { ListGroup } from "react-bootstrap";
import moment from "moment";
import RoundImage from "./RoundedImg";

export default function Chats({ avatar, username, message, timeSent }) {
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

export function ChatList({ chats }) {
  return (
    <ListGroup>
      {chats.map((chat) => (
        <Chats
          key={chat.id}
          avatar={chat.avatar}
          username={chat.username}
          message={chat.message}
          timeSent={chat.timeSent}
        />
      ))}
    </ListGroup>
  );
}