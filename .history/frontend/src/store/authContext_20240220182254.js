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

  const logoutHandler =  useCallback () =>  {
    setUser({...user, isLogIn: false});
    ws.current?.send.bind(ws.current)({type: "logout", message: ""});
    ws.current?.close();
  };
  useEffect(() => {
    // Check session only when the component mounts
    checkSession();

    // Clean up WebSocket connection
  }, []); // Empty dependency array ensures useEffect runs only once on mount

  const logoutHere = useCallback(() => {
     logoutHandler()
  },[logoutHandler])
  useEffect(() => {
    if (user.isLogIn) {
      const socket = new WebSocket("ws://localhost:8000/ws");
      socket.onopen = () => setIsWsReady(true);
      socket.onclose = () => setIsWsReady(false);
      socket.onmessage = (event) => setWsVal(event.data);

      ws.current = socket;

      // Clean up WebSocket connection
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
        
      }
    }
  }, [isWsReady, wsVal]);
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