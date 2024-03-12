import ChatList from "./ChatList";
import {ListGroup} from "react-bootstrap";
import AuthContext from "../store/authContext";
import {useContext} from "react";
export default function Chats({chat}) {
  const {user, openChat} = useContext(AuthContext);

  return (
    <ListGroup
      style={{
        marginTop: "4rem",
        display: "flex",
        flexDirection: "column",
        gap: "2rem",
      }}
    >
      {chat.map((chat, index) => (
        <ChatList
          key={index}
          avatar={chat.avatar}
          username={chat.nickname}
          isOnline={true}
          opeChat={openChat}
        />
      ))}
    </ListGroup>
  );
}
