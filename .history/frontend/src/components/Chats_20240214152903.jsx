import Chat from "./Chat";
import {ListGroup} from "react-bootstrap";
export default function Chats() {
  const chatList = [
    {
      vatar: "user1.jpg",
      nickname: "John Doe",
      lastMessage: "Hello, how are you?",
      timeSent: `2024-02-08T14:30:36.878396Z`,
    },
    {
      avatar: "user2.jpg",
      nickname: "Jane Smith",
      lastMessage: "I'm good, thanks!",
      timeSent: `2024-02-07T14:30:36.878396Z`,
    },
    {
      avatar: "user2.jpg",
      nickname: "Jane Smith",
      lastMessage: "I'm good, thanks!",
      timeSent: `2024-02-07T14:30:36.878396Z`,
    },
  ];
  console.log(chatList);
  return (
    <ListGroup>
      {chatList.map((chat, index) => (
        <Chat
          key={index}
          avatar={chat.avatar}
          username={chat.nickname}
          message={chat.lastMessage}
          timeSent={chat.timeSent}
        />
      ))}
    </ListGroup>
  );
}