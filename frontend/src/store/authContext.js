import React, {useState, useEffect, useCallback, useRef} from "react";
import {getJson} from "../helpers/helpers";
import User from "../components/User";

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

  const [dashBoardData, setDashBoardData] = useState({
    notifications: [],
    groups: [],
    chat: [],
    Posts: [],
  });

  const [onlineUsers, setOnlineUsers] = useState([]);
  const [isWsReady, setIsWsReady] = useState(false);
  const [wsVal, setWsVal] = useState(null);
  const ws = useRef(null);

  const [openChatDetails, setOpenChatDetails] = useState({
    messages: [],
    isChatOpen: false,
    openChatId: "",
    type: "",
  });
  const [profileData, setProfileData] = useState({
    isComponentVisible: false,
    data: {},
    error: {type: "", message: ""},
  });

  const setDashBoardDataOutside = () => {
    setDashBoardData({
      notifications: [],
      groups: [],
      chat: [],
      Posts: [],
    });
  };

  // dash broad stuff
  const fetchDashboard = useCallback(async () => {
    try {
      const data = await getJson("dashboard", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },

        body: JSON.stringify(user.Id),
      });

      setDashBoardData(data);
    } catch (err) {
      console.log(err);
    }
  }, [user.Id]);

  // profile

  const fetchProfileData = useCallback(async (id) => {
    try {
      const res = await getJson("profile", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(id),
      });

      setProfileData((prev) => ({...prev, data: res}));
    } catch (err) {
      console.log(err);
    }
  }, []);

  const toggleProfilePrivacy = async () => {
    let s =
      profileData.data?.Owner?.PrivacySetting === "private"
        ? "public"
        : "private";

    try {
      const res = await getJson("toggle-privacy", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          id: profileData.data.Owner.Id,
          setting: s,
        }),
      });
      if (res.message) {
        setProfileData((prev) => ({
          ...prev,
          data: {
            ...prev.data,
            Owner: {...prev.data.Owner, PrivacySetting: s},
          },
        }));
      }
    } catch (err) {
      console.log(err);
    }
  };

  const resetIsProfileComponentVisible = useCallback((value) => {
    setProfileData((prev) => ({...prev, isComponentVisible: value}));
  }, []);

  const [selectedPosts, setSelectedPosts] = useState([]);
  const logintHandler = (user) => {
    setUser({...user.user, isLogIn: true});
  };

  const openChat = async (id) => {
    try {
      const res = await getJson("get-messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          currUser: User.Id
           

        }),
      });

      setOpenChatDetails({
        messages: res,
        isChatOpen: true,
        openChatId: id,
      });
    } catch (er) {
      console.log();
    }
  };

  const closeChat = () => {
    setOpenChatDetails({messages: [], isChatOpen: false, openChatId: ""});
  };

  // web socket states

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
  }, []);

  // websocket setup;
  useEffect(() => {
    if (user.isLogIn) {
      const socket = new WebSocket("ws://localhost:8000/ws");
      socket.onopen = () => setIsWsReady(true);
      socket.onclose = () => setIsWsReady(false);
      socket.onmessage = (event) => setWsVal(event.data);

      ws.current = socket;

      return () => {
        socket.close();
      };
    }

    // Clean up WebSocket connection
  }, [user.isLogIn]);

  // notifications

  const handleWebsocketNotification = useCallback(() => {
    if (dashBoardData?.notifications.length === 0) return;
    switch (dashBoardData.notifications[0].type) {
      case "notification requestToFollow":
        if (profileData.isComponentVisible) {
          if (profileData?.data?.Owner?.PrivacySetting === "public") {
            setProfileData((prev) => ({
              ...prev,
              data: {
                ...prev.data,

                Followers: Array.isArray(prev.data.Followers)
                  ? [
                      dashBoardData.notifications[0].payload.Data,
                      ...prev.data.Followers,
                    ]
                  : [dashBoardData.notifications[0].Body.Data],
                IsFollowed: true,
              },
            }));

            return;
          }

          setProfileData((prev) => {
            if (prev.data?.PendingFollowers?.length === 0) {
              return {
                ...prev,
                data: {
                  ...prev.data,
                  PendingFollowers: [
                    dashBoardData.notifications[0].payload.Data,
                  ],
                },
              };
            }
            return {
              ...prev,
              data: {
                ...prev.data,
                PendingFollowers: Array.isArray(prev.data.PendingFollowers)
                  ? [
                      dashBoardData.notifications[0].payload.Data,
                      ...prev.data.PendingFollowers,
                    ]
                  : [dashBoardData.notifications[0].payload.Data],
              },
            };
          });
        }
        break;

      case "notification answerRequestToFollow":
        if (profileData.isComponentVisible) {
          if (dashBoardData.notifications[0].payload.Data.Reply === "no") {
            setProfileData((prev) => ({
              ...prev,
              data: {...prev.data, IsPending: false},
            }));
            return;
          }
          fetchProfileData(
            dashBoardData.notifications[0].payload.Data.SenderId
          );
        }
        break;

      case "notification unfollow":
        setProfileData((prev) => ({
          ...prev,
          data: {
            ...prev.data,
            Followers: Array.isArray(prev.data.Followers)
              ? prev.data.Followers.filter(
                  (item) =>
                    item.Id !== dashBoardData.notifications[0].payload.Data.Id
                )
              : [],
          },
        }));

        break;
      default:
        console.log(
          "unknow nofication type",
          dashBoardData.notifications,
          "hehhere"
        );
        break;
    }
    // setNotification({
    //   type: "requestToFollow",
    //   message: dashBoardData.notifications[0].Body.message,
    // });
  }, [
    dashBoardData.notifications,
    profileData.isComponentVisible,
    fetchProfileData,
    profileData?.data?.Owner?.PrivacySetting,
  ]);

  useEffect(() => {
    handleWebsocketNotification();
  }, [dashBoardData.notifications, handleWebsocketNotification]);

  //websocket actions errors

  const handleWebsocketErrors = useCallback((data) => {
    switch (data.message) {
      case "requestToFollow":
        setProfileData((prev) => ({
          ...prev,
          error: {type: "requestToFollow", message: data.message},
        }));
        break;
      case "":
        break;

      default:
        break;
    }
  }, []);

  const handleWebsocketSucess = useCallback(
    (data) => {
      const followers = {
        Id: user.Id,
        Nickname: user.Nickname,
        Avatar: user.Avatar,
        FirstName: user.FirstName,
        LastName: user.LastName,
        PrivacySetting: user.Privicy_setting,
      };
      switch (data.message) {
        case "requestToFollow":
          //  we need to broadcast the when ever privacy changes any profile we viewing then that profilr needs
          // to be updated and refetched.
          // fetchProfileData(profileData?.data?.Owner?.Id);
          if (profileData.data?.Owner?.PrivacySetting === "public") {
            setProfileData((prev) => ({
              ...prev,
              data: {
                ...prev.data,
                Followers: Array.isArray(prev.data.Followers)
                  ? [followers, ...prev.data.Followers]
                  : [followers],
                IsFollowed: true,
              },
            }));
            return;
          }

          setProfileData((prev) => ({
            ...prev,
            data: {...prev.data, IsPending: true},
          }));
          break;
        case "answerRequestToFollow no":
          //  remove the user from the pending list

          setProfileData((prev) => ({
            ...prev,
            data: {
              ...prev.data,
              PendingFollowers: prev.data.PendingFollowers.filter(
                (item) => item.Id !== data.whatever
              ),
            },
          }));

          break;

        case "answerRequestToFollow yes":
          //  add the user to the followers list and remove from the pending list

          fetchProfileData(profileData?.data?.Owner?.Id);

          break;
        case "unfollow":
          // fetchProfileData(profileData?.data?.Owner?.Id);
          if (profileData.data?.Owner?.PrivacySetting === "public") {
            setProfileData((prev) => ({
              ...prev,
              data: {
                ...prev.data,
                Followers: data.Followers?.filter((el) => el.Id !== user.Id),
                IsFollowed: false,
              },
            }));
            return;
          }

          setProfileData((prev) => ({
            ...prev,
            data: {
              ...prev.data,
              IsFollowed: false,

              Followers: Array.isArray(prev.data.Followers)
                ? prev.data.Followers.filter(
                    (item) => item.Id !== data.whatever
                  )
                : [],
            },
          }));

          break;
        default:
          console.log("unknown websocket message");
          break;
      }
    },
    [
      user.Id,
      profileData.data?.Owner?.PrivacySetting,
      fetchProfileData,
      profileData.data?.Owner?.Id,
    ]
  );

  const currentProfilePrivacyChanged = useCallback(
    (data) => {
      if (!profileData.isComponentVisible) {
        return;
      }

      if (data.body.senderId === user.Id) {
        setProfileData((prev) => ({
          ...prev,
          data: {
            ...prev.data,
            Owner: {
              ...prev.data.Owner,
              PrivacySetting: data.body.privacySetting,
            },
          },
        }));
      } else {
        fetchProfileData(profileData?.data?.Owner?.Id);
      }
    },
    [
      profileData.isComponentVisible,
      fetchProfileData,
      profileData?.data?.Owner?.Id,
      user.Id,
    ]
  );

  const addPostToFeed = useCallback(
    (data) => {
      const {
        Body,
        Title,
        CreatedAt,
        CreatorId,
        GroupId,
        Id,
        PrivacyLevel,
        Image,
      } = data.message;

      if (
        profileData.isComponentVisible &&
        CreatorId === profileData.data.Owner.Id
      ) {
        setProfileData((prev) => ({
          ...prev,
          data: {
            ...prev.data,
            Posts: [
              {
                Body,
                Title,
                CreatedAt,
                CreatorId,
                GroupId,
                Id,
                privacyLevel: PrivacyLevel,
                Image: Image,
              },
              ...prev.data.Posts,
            ],
          },
        }));
      } else {
        setDashBoardData((prev) => ({
          ...prev,
          Posts: [
            {
              Body,
              Title,
              CreatedAt,
              CreatorId,
              GroupId,
              Id,
              privacyLevel: PrivacyLevel,
              Image: Image,
            },
            ...prev?.Posts,
          ],
        }));
      }
    },
    [
      profileData.isComponentVisible,
      profileData.data?.Owner?.Id,
      setProfileData,
      setDashBoardData,
    ]
  );

  useEffect(() => {
    if (isWsReady) {
      const data = JSON.parse(wsVal);
      console.log(data, "data from websocket");

      if (data?.type === "error") {
        handleWebsocketErrors(data);
        return;
      }

      if (data?.type === "success") {
        handleWebsocketSucess(data);
        return;
      }

      if (
        data?.type?.includes("notification") &&
        dashBoardData.notifications.some((el) => el.Id === data.Id) === false
      ) {
        setDashBoardData((prev) => ({
          ...prev,
          notifications: [data, ...prev.notifications],
        }));
        return;
      }

      switch (data?.type) {
        case "logout":
          logoutHandler(false);
          break;
        case "online":
          setOnlineUsers(data);
          break;
        case "post":
          console.log(data.message);
          addPostToFeed(data);
          break;
        case "togglePrivacySetting":
          currentProfilePrivacyChanged(data);
          break;

        default:
          // console.log("no action", data);
          break;
      }
    }
  }, [
    isWsReady,
    wsVal,
    logoutHandler,
    addPostToFeed,
    currentProfilePrivacyChanged,
    handleWebsocketErrors,
    handleWebsocketSucess,
    dashBoardData.notifications,
  ]);

  const onAddLikeDislikePost = (id, data) => {
    if (profileData.isComponentVisible) {
      console.log(data, id, "if profile PostlikeDislike");
      setProfileData((prev) => ({
        ...prev,
        data: {
          ...prev.data,
          Posts: prev.data.Posts.map((el) =>
            el.Id === id
              ? {...el, Likes: data.Likes, Dislikes: data.Dislikes}
              : el
          ),
        },
      }));
    } else {
      setDashBoardData((prev) => ({
        ...prev,
        Posts: prev.Posts.map((el) =>
          el.Id === id
            ? {...el, Likes: data.Likes, Dislikes: data.Dislikes}
            : el
        ),
      }));
    }
  };

  //   if(quary === "like"){
  //   const { likes, dislikes , userlikes,  ...rest} = posts.find(el => el.id ===id)
  //      SetPosts( posts.map(el=> el.id === id ? {likes: data.likes, dislikes: data.dislikes, userlikes: data.userlikes, ...rest}   : el  ))
  //   }else{
  //     const { likes, dislikes,  userlikes,  ...rest} = posts.find(el => el.id ===id)
  //     SetPosts( posts.map(el=> el.id === id ? {likes: data.likes, dislikes: data.dislikes,  userlikes: data.userlikes, ...rest}   : el  ))
  //   }
  // }

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
        openChat,
        openChatDetails,
        onlineUsers,
        fetchProfileData,
        profileData,
        resetIsProfileComponentVisible,
        fetchDashboard,
        dashBoardData,
        toggleProfilePrivacy,
        setDashBoardDataOutside,
        onAddLikeDislikePost,
      }}
    >
      {props.children}
    </AuthContext.Provider>
  );
};
export default AuthContext;
