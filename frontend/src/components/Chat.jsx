import React, {useContext, useState, useRef, useEffect} from "react";
import {Form, Button, ListGroup} from "react-bootstrap";
import User from "./User";
import MyModal from "./Modal";
import AuthContext from "../store/authContext";
import moment from "moment";

const ChatComponent = () => {
  const [newMessage, setNewMessage] = useState("");
  const prevScrollPosRef = useRef(0);
  const endMessageRef = useRef(null);

  const {
    OpenChat: handleShow,
    closeChat: handleClose,
    openChatDetails: {messages, isChatOpen, type, user, group},
  } = useContext(AuthContext);

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (newMessage.trim() !== "") {
    }
  };

  const helperTimeFormater = (time) => {
    return moment(time).calendar(null, {
      sameDay: "h:mm A",
      lastWeek: "dddd",
      sameElse: "MM/DD/YYYY",
    });
  };

  useEffect(() => {
    if (isChatOpen || messages) {
      endMessageRef?.current?.scrollIntoView({behavior: "smooth"});
    }
  }, [isChatOpen, messages]); // Scroll to the bottom when messages change

  return (
    <MyModal
      handleClose={handleClose}
      handleShow={handleShow}
      show={isChatOpen}
      flag={false}
    >
      <div style={{margin: "1rem 0 2rem 0"}}>
        {type === "user" ? (
          <User
            isOnline={user?.isOnline}
            Nickname={user?.username}
            avatar={user?.avatar}
          />
        ) : (
          <div>
            <h1> {group?.title} </h1>
          </div>
        )}
      </div>
      <div>
        <ListGroup
          id="scrollableDiv"
          style={{
            display: "flex",
            flexDirection: "column-reverse  !important",
            overflowY: "auto",
            height: "400px",
          }}
          ref={prevScrollPosRef}
        >
          {messages.map((message) => (
            <React.Fragment key={message.id}>
              <div className="chats">
                <span style={{textAlign: "center", color: "#898888"}}>
                  {helperTimeFormater(message.time)}
                </span>
                <ListGroup.Item action as="li">
                  <p style={{margin: "0"}}>{message.text}</p>
                </ListGroup.Item>
              </div>
            </React.Fragment>
          ))}

          <div ref={endMessageRef} />
        </ListGroup>
      </div>
      <Form onSubmit={handleSendMessage}>
        <Form.Group>
          <Form.Control
            type="text"
            placeholder="Type your message"
            value={newMessage}
            style={{margin: "2rem 0"}}
            onChange={(e) => setNewMessage(e.target.value)}
          />
        </Form.Group>

        <Button variant="primary" type="submit">
          Send
        </Button>
      </Form>
    </MyModal>
  );
};

export default ChatComponent;
