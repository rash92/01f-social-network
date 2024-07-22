import React, {useContext, useState, useRef, useEffect} from "react";
import {Form, Button, ListGroup} from "react-bootstrap";
import User from "./User";
import MyModal from "./Modal";
import AuthContext from "../store/authContext";
import moment from "moment";

const ChatComponent = () => {
  const [newMessage, setNewMessage] = useState("");

  const endMessageRef = useRef(null);

  const {
    onlineUsers,
    isWsReady,
    wsMsgToServer,
    OpenChat: handleShow,
    closeChat: handleClose,
    user: loggedUser,
    groupId,
    groupData,
    openChatDetails: {messages, isChatOpen, user},
  } = useContext(AuthContext);

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (newMessage.trim() !== "" && isWsReady) {
      // type GroupMessage struct {
      //   Id        string    ⁠ json:"Id" ⁠
      //   Type      string    ⁠ json:"type" ⁠
      //   SenderId  string    ⁠ json:"SenderId" ⁠
      //   GroupId   string    ⁠ json:"GroupId" ⁠
      //   Message   string    ⁠ json:"Message" ⁠
      //   CreatedAt time.Time ⁠ json:"CreatedAt" ⁠
      // }
      if (user.type === "profileMessage") {
        wsMsgToServer(
          JSON.stringify({
            type: user.type,
            message: {
              SenderId: loggedUser.Id,
              ReceiverId: user.id,
              message: newMessage,
              type: user.type,
              createAt: "",
              Nickname: loggedUser.Nickname,
              Avatar: loggedUser.Avatar,
            },
          })
        );
      } else {
        wsMsgToServer(
          JSON.stringify({
            type: user.type,
            message: {
              Id: "",
              SenderId: loggedUser.Id,
              ReceiverId: user.id,
              Message: newMessage,
              GroupId: groupId || "",
              type: user.type,
              CreateAt: "",
              Nickname: loggedUser.Nickname,
              Avatar: loggedUser.Avatar,
            },
          })
        );
      }

      setNewMessage("");
    }
  };

  const formatCreatedAt = (createdAt) => {
    return moment(createdAt).format("MMM DD, YYYY hh:mm A");
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
        {user.type === "privateMessage" ? (
          <User
            Nickname={user?.Nickname}
            Avatar={user?.Avatar}
            IsOnline={onlineUsers.includes(user.id)}
          />
        ) : (
          <div>
            <h1> {groupData?.data?.BasicInfo?.Title} </h1>
          </div>
        )}
      </div>
      <div>
        <ListGroup
          id="scrollableDiv"
          style={{
            display: "flex",
            overflowY: "auto",
            height: "400px",
            flexDirection: "column-reverse !important",
          }}
        >
          {messages &&
            messages?.map((message) => (
              <React.Fragment key={message.Id}>
                <div className="chats">
                  <span style={{textAlign: "center", color: "#898888"}}>
                    {formatCreatedAt(message.CreateAt)}
                  </span>
                  <ListGroup.Item action as="li">
                    <User
                      Nickname={message?.Nickname}
                      Avatar={message?.Avatar}
                      IsOnline={onlineUsers.includes(message.SenderId)}
                    />

                    <p style={{margin: "1rem 0"}}>{message.Message}</p>
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
