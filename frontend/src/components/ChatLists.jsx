import ChatList from "./ChatList";
import {ListGroup} from "react-bootstrap";
import AuthContext from "../store/authContext";
import {useContext} from "react";
export default function Chats({chats}) {
  const {onlineUsers, openChat} = useContext(AuthContext);

  return (
    <ListGroup
      style={{
        marginTop: "4rem",
        display: "flex",
        flexDirection: "column",
        gap: "2rem",
      }}
    >
      {chats?.map((chat, index) => (
        <ChatList
          key={index}
          avatar={chat.Avatar}
          nickname={chat.Nickname}
          isOnline={onlineUsers.includes(chat.Id)}
          opeChat={openChat}
          id={chat.Id}
        />
      ))}
    </ListGroup>
  );
}
