import React from "react";
import {ListGroup, Image} from "react-bootstrap";
import moment from "moment";

export default function Chats({username, message, timeSent}) {
  const formattedTime = moment(timeSent).format("h:mm A");

  return (
    <ListGroup>
      <ListGroup.Item>
        <Image src="profile-pic.jpg" roundedCircle />
        <span>{username}</span>
        <span>{message}</span>
        <span>{formattedTime}</span>
      </ListGroup.Item>
    </ListGroup>
  );
}
// export default function Chats() {
//   const chatList = [
//     {
//       profilePic: "user1.jpg",
//       username: "John Doe",
//       lastMessage: "Hello, how are you?",
//       timeSent: "10:30 AM"
//     },
//     {
//       profilePic: "user2.jpg",
//       username: "Jane Smith",
//       lastMessage: "I'm good, thanks!",
//       timeSent: "11:45 AM"
//     },
//     // Add more chat objects as needed
//   ];

//   return (
//     <div>
//       {chatList.map((chat, index) => (
//         <div key={index}>
//           <img src={chat.profilePic} alt="Profile Pic" />
//           <div>
//             <h3>{chat.username}</h3>
//             <p>{chat.lastMessage}</p>
//             <p>{chat.timeSent}</p>
//           </div>
//         </div>
//       ))}
//     </div>
//   );
// }