import React, {useState, useEffect, useCallback, useRef} from "react";
import {getJson} from "../helpers/helpers";
const userObj = {
  id: "",
  isLogIn: false,
  username: "",
  profileImg: "",
};
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
    setUser({...user, isLogIn: true});
  };
  const logoutHandler = () => {
    setUser(userObj);
  };
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

  useEffect(() => {
    // Check session only when the component mounts
    checkSession();

    const socket = new WebSocket("wss://echo.websocket.events/");
    socket.onopen = () => setIsWsReady(true);
    socket.onclose = () => setIsWsReady(false);
    socket.onmessage = (event) => setWsVal(event.data);

    ws.current = socket;

    // Clean up WebSocket connection
    return () => {
      socket.close();
    };
  }, []); // Empty dependency array ensures useEffect runs only once on mount

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
