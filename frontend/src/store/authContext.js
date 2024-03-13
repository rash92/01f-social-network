import React, {useState, useEffect, useCallback, useRef} from "react";
import {getJson} from "../helpers/helpers";
const userObj = {
  Id: "",
  isLogIn: false,
  Nickname: "",
  Email: "",
  FirstName: "",
  LastName: "",
  Profile: "",
  DOB: "",
  Privicy_setting: "",
  CreatedAt: "",
  AboutMe: "",
};

// const logoutHandler = useCallback(() => {
//   setUser(userObj);
//   if (ws.current && ws.current.readyState === WebSocket.OPEN) {
//     ws.current.send(JSON.stringify({ type: "logout", message: "" }));
//     ws.current.close();
//   }
// }, []);
const AuthContext = React.createContext({
  user: userObj,
  OnLogin: () => {},
  onLogout: () => {},
  isWsReady: false,
  wsVal: null,
  wsMsgToServer: (msg) => {},
  showChat: false,
  openChat: () => {},
  closeChat: () => {},
  openChatDetails: {},

  // OnAddCommentToPost: () => {},
  // posts: [],
  // catogaries: [],
  // username: "",
  // selectedPosts: [],
  // setSelectedPosts: () => {},
  // OnAddPost: () => {},
});

export const AuthContextProvider = (props) => {
  const [user, setUser] = useState(userObj);

  const [selectedPosts, setSelectedPosts] = useState([]);
  const logintHandler = (user) => {
    setUser({...user.user, isLogIn: true});
  };
  const [onlineUsers, setOnlineUsers] = useState([]);

  const [showChat, setShowChat] = useState(false);
  const [openChatDetails, setOpenChatDetails] = useState({});
  const openChat = (data) => {
    setOpenChatDetails(data);
    setShowChat(true);
  };

  const closeChat = () => {
    setShowChat(false);
  };

  // web socket states
  const [isWsReady, setIsWsReady] = useState(false);
  const [wsVal, setWsVal] = useState(null);
  const ws = useRef(null);

  const checkSession = useCallback(async () => {
    try {
      const res = await getJson("checksession", {
        method: "GET",
        credentials: "include",
      });

      if (!res.success || res.status === 401) {
        logintHandler(user);
        throw new Error("something went wrong presist login");
      }

      logintHandler(res);
    } catch (err) {
      console.log("error", err);
    }
  }, [user]);

  const logoutHandler = useCallback((flag = true) => {
    setUser(userObj);

    if (flag) {
      ws.current?.send.bind(ws.current)(
        JSON.stringify({type: "logout", message: ""})
      );
    }
    ws.current?.close.bind(ws.current)();
  }, []);
  useEffect(() => {
    // Check session only when the component mounts
    checkSession();
  }, []); // Empty dependency array ensures useEffect runs only once on mount

  useEffect(() => {
    if (user.isLogIn) {
      const socket = new WebSocket("ws://localhost:8000/ws");
      socket.onopen = () => setIsWsReady(true);
      socket.onclose = () => setIsWsReady(false);
      socket.onmessage = (event) => setWsVal(event.data);

      ws.current = socket;

      return () => {
        if (isWsReady) {
          socket.close();
        }
      };
    }

    // Clean up WebSocket connection
  }, [isWsReady, user.isLogIn]);

  useEffect(() => {
    if (isWsReady) {
      const data = JSON.parse(wsVal);
      console.log(data);
      if (data?.type === "logout") {
        console.log("logout");
        logoutHandler(false);
      }

      if (data?.type === "online") {
        setOnlineUsers(data);
      }
    }
  }, [isWsReady, wsVal, logoutHandler]);
  return (
    <AuthContext.Provider
      value={{
        user,
        OnLogin: logintHandler,
        onLogout: logoutHandler,
        isWsReady,
        wsVal,
        wsMsgToServer: ws.current?.send.bind(ws.current),

        // catogaries: catogaries,
        // posts: posts,
        selectedPosts,
        setSelectedPosts,
        // OnAddCommentToPost: OnAddCommentToPost,
        // OnAddPost,
        // onAddLikeDislikePost,
        // onAddLikeDislikeComment,
        // onRemovePost,
        closeChat,
        showChat,
        openChat,
        openChatDetails,
      }}
    >
      {props.children}
    </AuthContext.Provider>
  );
};
export default AuthContext;
