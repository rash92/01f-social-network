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
  const [showChat, setShowChat] = useState(false);
  const [openChatDetails, setOpenChatDetails] = useState({});
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

  console.log(dashBoardData);
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

  const openChat = (data) => {
    setOpenChatDetails(data);
    setShowChat(true);
  };

  const closeChat = () => {
    setShowChat(false);
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
    console.log();
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
    [user.Id, profileData.data?.Owner?.PrivacySetting, fetchProfileData]
  );

  const currentProfilePrivacyChanged = useCallback(
    (data) => {
      if (
        profileData.isComponentVisible &&
        data.Id === profileData.data.Owner.Id
      ) {
        fetchProfileData(profileData.data.Owner.Id);
      }
    },
    [
      profileData.isComponentVisible,
      profileData.data?.Owner?.Id,
      fetchProfileData,
    ]
  );

  const addPostToFeed = useCallback(
    (data) => {
      const {
        body: Body,
        title: Title,
        createAt: CreatedAt,
        creatorId: CreatorId,
        groupId: GroupId,
        id: Id,
        privacyLevel: PrivacyLevel,
        image: Image,
      } = data.message;
      if (
        profileData.isComponentVisible &&
        CreatorId === profileData.data.Owner.Id
      ) {
        console.log(data.message.CreatorId);
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
                image: Image,
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
              image: Image,
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
          console.log(data, "this the post got back from websocket");
          addPostToFeed(data);
          break;
        case "privacyChanged":
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
  ]);

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
        onlineUsers,
        fetchProfileData,
        profileData,
        resetIsProfileComponentVisible,
        fetchDashboard,
        dashBoardData,
        toggleProfilePrivacy,
        setDashBoardDataOutside,
      }}
    >
      {props.children}
    </AuthContext.Provider>
  );
};
export default AuthContext;
