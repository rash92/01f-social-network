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
      <div className="d-flex w-100 gap-4 mb-3">
        <RoundImage src={avatar} alt="Avatar" size={40} />
        <span style={{fontSize: "1rem"}} className="fw-bold ">
          {username}
        </span>
      </div>


     <div>
      
     </div>
      <span>{message}</span>
      <span>{formattedTime}</span>
    </ListGroup.Item>
  );
}