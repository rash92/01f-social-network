import React, {
  useContext,
  useState,
  useRef,
  useEffect,
  useCallback,
} from "react";
import {Form, Button, ListGroup} from "react-bootstrap";
import User from "./User";
import MyModal from "./Modal";
import AuthContext from "../store/authContext";
import moment from "moment";
import {throttle} from "lodash";

const dummyMessages = Array.from({length: 10}, (_, index) => ({
  id: index + 1,
  text: "Hello, World!",
  time: new Date(),
}));

const ChatComponent = () => {
  const [messages, setMessages] = useState(dummyMessages);
  const [newMessage, setNewMessage] = useState("");
  const [hasMore, setHasMore] = useState(true);
  const [isLoading, setIsLoading] = useState(false);

  const prevScrollPosRef = useRef(0);
  const endMessageRef = useRef(null);
  const fistOfTheLastTen = useRef(null);

  const {
    OpenChat: handleShow,
    closeChat: handleClose,
    showChat: show,
    openChatDetails: user,
  } = useContext(AuthContext);

  const handleSendMessage = (e) => {
    e.preventDefault();
    if (newMessage.trim() !== "") {
      setMessages([
        ...messages,
        {id: messages.length + 1, text: newMessage, time: new Date()},
      ]);
      setNewMessage("");
    }
  };

  const helperTimeFormater = (time) => {
    return moment(time).calendar(null, {
      sameDay: "h:mm A",
      lastWeek: "dddd",
      sameElse: "MM/DD/YYYY",
    });
  };

  const fetchMoreMessages = async () => {
    if (messages.length <= 30) {
      setIsLoading(true); // Set isLoading to true while loading

      console.log("fetching more messages");

      setTimeout(() => {
        const newMessages = Array.from({length: 10}, (_, index) => ({
          id: messages.length + index + 1,
          text: "Hello, earth! ",
          time: new Date(),
        }));

        setMessages((prevMessages) => [...prevMessages, ...newMessages]);

        setIsLoading(false); // Set isLoading to false after loading
      }, 500);
    } else if (messages.length >= 30) {
      setHasMore(false);
    }
  };

  const throttleChatScrollHandler = useCallback(
    throttle(async (e) => {
      const el = e.target;
      const containerHeight = el.scrollHeight;
      const scrollTop = el.scrollTop;
      const percentageScrolledUp = (scrollTop / containerHeight) * 100;

      if (
        percentageScrolledUp <= 10 &&
        scrollTop < prevScrollPosRef.current?.prevScrollPos &&
        !isLoading
      ) {
        await fetchMoreMessages();
      }

      prevScrollPosRef.current.prevScrollPos = scrollTop;
    }, 500),
    [isLoading, messages]
  );

  useEffect(() => {
    if (show || newMessage) {
      endMessageRef.current.scrollIntoView({behavior: "smooth"});
    }
  }, [show, newMessage]); // Scroll to the bottom when messages change

  return (
    <MyModal
      handleClose={handleClose}
      handleShow={handleShow}
      show={show}
      flag={false}
    >
      <div style={{margin: "1rem 0 2rem 0"}}>
        {user?.type === "user" ? (
          <User
            isOnline={user?.isOnline}
            userName={user?.username}
            avatar={user?.avatar}
          />
        ) : (
          <div> {user?.userName}</div>
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
            gap: "1rem",
          }}
          ref={prevScrollPosRef}
          onScroll={throttleChatScrollHandler}
        >
          {!hasMore && (
            <p style={{textAlign: "center"}}>
              <b>Yay! You have seen it all</b>
            </p>
          )}
          {isLoading && <h4>Loading...</h4>}
          {messages.map((message, i) => (
            <React.Fragment key={message.id}>
              <div ref={i === 9 ? fistOfTheLastTen : null} className="chats">
                <span style={{textAlign: "center", color: "#898888"}}>
                  {helperTimeFormater(message.time)}
                </span>
                <ListGroup.Item action as="li">
                  <p style={{margin: "0"}}>{message.text + (i + 1)}</p>
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
