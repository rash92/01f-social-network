import ChatList from "./ChatL";
import {ListGroup} from "react-bootstrap";
import AuthContext from "../store/authContext";
import {useContext} from "react";
export default function Chats() {
  const {user} = useContext(AuthContext);
  const chatList = [
    {
      avatar: user.profileImg,
      nickname: "John Doe",
      lastMessage: "Hello, how are you?",
      timeSent: `2024-02-14T14:30:36.878396Z`,
    },
    {
      avatar: user.profileImg,
      nickname: "Jane Smith",
      lastMessage: "I'm good, thanks!",
      timeSent: `2024-02-13T14:30:36.878396Z`,
    },
    {
      avatar: user.profileImg,
      nickname: "Jane Smith",
      lastMessage: "I'm good, thanks!",
      timeSent: `2024-02-07T14:30:36.878396Z`,
    },
  ];
  console.log(chatList);
  return (
    <ListGroup
      style={{
        marginTop: "4rem",
        display: "flex",
        flexDirection: "column",
        gap: "2rem",
      }}
    >
      {chatList.map((chat, index) => (
        <Chat
          key={index}
          avatar={chat.avatar}
          username={chat.nickname}
          message={chat.lastMessage}
          timeSent={chat.timeSent}
          isOnline={true}
        />
      ))}
    </ListGroup>
  );
}