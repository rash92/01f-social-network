import React from "react";
import {ListGroup, Image} from "react-bootstrap";

const NotificationList = ({notifications}) => {
  return (
    <div
      style={{
        marginTop: "3rem",
      }}
    >
      <ListGroup>
        {notifications.map((notification) => (
          <ListGroup.Item
            key={notification.Id}
            style={{
              display: "flex",
              alignItems: "center",
              marginBottom: "40px",
            }}
          >
            <Image
              src={`http://localhost:8000/images/${notification.SenderAvatar}`}
              roundedCircle
              style={{marginRight: "10px", width: "50px", height: "50px"}}
            />
            <span
              style={{
                fontSize: "0.9rem",
              }}
            >
              {notification.Body}
            </span>
          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  );
};

export default NotificationList;
