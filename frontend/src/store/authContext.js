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
});

export const AuthContextProvider = (props) => {
  const [user, setUser] = useState(userObj);
  const [groupId, setGroupId] = useState("");

  const updateGroupId = (id) => {
    setGroupId(id);
  };

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
    user: {},
    isChatOpen: false,
  });

  const [profileData, setProfileData] = useState({
    isComponentVisible: false,
    data: {},
    error: {type: "", message: ""},
  });

  const [groupData, setGroupData] = useState({
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

  const fetchGroupData = useCallback(async (id) => {
    try {
      const res = await getJson("group", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(id),
      });

      setGroupData((prev) => ({...prev, data: res}));
    } catch (err) {
      console.log(err);
    }
  }, []);

  const toggleProfilePrivacy = async () => {
    let s =
      profileData?.data?.Owner?.PrivacySetting === "private"
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
          id: profileData?.data?.Owner?.Id,
          setting: s,
        }),
      });
      if (res.message) {
        setProfileData((prev) => ({
          ...prev,
          data: {
            ...prev.data,
            Owner: {...prev?.data?.Owner, PrivacySetting: s},
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
  const resetIsGroupComponentVisible = useCallback((value) => {
    setGroupData((prev) => ({...prev, isComponentVisible: value}));
  }, []);

  const [selectedPosts, setSelectedPosts] = useState([]);
  const logintHandler = (user) => {
    setUser({...user.user, isLogIn: true});
  };

  const openChat = async (data) => {
    try {
      const res = await getJson("get-messages", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },

        credentials: "include",
        body: JSON.stringify({
          currUser: user.Id,
          otherUser: data.id,
          type: data.type,
        }),
      });

      setOpenChatDetails({
        messages: res.messages,
        isChatOpen: true,
        user: data,
      });
    } catch (er) {
      console.log();
    }
  };

  const closeChat = () => {
    setOpenChatDetails({messages: [], user: {}, isChatOpen: false});
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

  const httpLogout = async () => {
    try {
      const res = await getJson("logout", {
        method: "Get",
        credentials: "include",
      });

      if (res.success) {
        console.log("we logout ");
      }
    } catch (err) {
      console.log("error", err);
    }
  };

  const logoutHandler = useCallback(async (flag = true) => {
    setUser(userObj);
    if (flag) {
      ws.current?.send.bind(ws.current)(
        JSON.stringify({type: "logout", message: ""})
      );
    }

    await httpLogout();
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
      socket.onmessage = (event) => {
        setWsVal(event?.data);
      };

      ws.current = socket;

      return () => {
        socket.close();
      };
    }

    // Clean up WebSocket connection
  }, [user.isLogIn]);

  // notifications

  const handleWebsocketNotification = useCallback(() => {
    if (
      dashBoardData?.notifications?.length === 0 ||
      dashBoardData?.notifications === null
    )
      return;
    switch (dashBoardData?.notifications[0]?.type) {
      case "notification requestToFollow":
        if (profileData?.isComponentVisible) {
          if (profileData?.data?.Owner?.PrivacySetting === "public") {
            setProfileData((prev) => ({
              ...prev,
              data: {
                ...prev?.data,

                Followers: Array.isArray(prev?.data?.Followers)
                  ? [
                      dashBoardData?.notifications[0]?.payload.Data,
                      ...prev?.data?.Followers,
                    ]
                  : [dashBoardData?.notifications[0]?.Body?.Data],
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
                  ...prev?.data,
                  PendingFollowers: [
                    dashBoardData?.notifications[0]?.payload?.Data,
                  ],
                },
              };
            }
            return {
              ...prev,
              data: {
                ...prev.data,
                PendingFollowers: Array.isArray(prev?.data?.PendingFollowers)
                  ? [
                      dashBoardData?.notifications[0]?.payload.Data,
                      ...prev?.data?.PendingFollowers,
                    ]
                  : [dashBoardData?.notifications[0]?.payload.Data],
              },
            };
          });
        }
        break;

      case "notification answerRequestToFollow":
        if (profileData?.isComponentVisible) {
          if (dashBoardData?.notifications[0]?.payload.Data?.Reply === "no") {
            setProfileData((prev) => ({
              ...prev,
              data: {...prev?.data, IsPending: false},
            }));
            return;
          }
          fetchProfileData(
            dashBoardData?.notifications[0]?.payload?.Data?.SenderId
          );
        }
        break;

      case "notification unfollow":
        setProfileData((prev) => ({
          ...prev,
          data: {
            ...prev.data,
            Followers: Array.isArray(prev?.data?.Followers)
              ? prev.data.Followers.filter(
                  (item) =>
                    item.Id !==
                    dashBoardData?.notifications[0]?.payload?.Data?.Id
                )
              : [],
          },
        }));

        break;
      case "notification requestToJoinGroup":
        if (dashBoardData?.notifications[0]?.payload?.groupId !== groupId)
          return;
        setGroupData((prev) => ({
          ...prev,
          data: {
            ...prev.data,
            RequestedMembers: Array.isArray(prev?.data?.RequestedMembers)
              ? [
                  ...prev.data?.RequestedMembers,
                  dashBoardData?.notifications[0].payload,
                ]
              : [dashBoardData?.notifications[0]?.payload],
          },
        }));
        break;
      case "notification answerRequestToJoinGroup":
        if (dashBoardData?.notifications[0]?.payload.groupId !== groupId)
          return;

        if (dashBoardData?.notifications[0]?.payload.type) {
          setGroupData((prev) => ({
            ...prev,
            data: dashBoardData.notifications[0].payload.group,
          }));
        } else {
          setGroupData((prev) => ({
            ...prev,
            data: {...prev.data, Status: "none"},
          }));
        }

        break;

      case "notification inviteToJoinGroup":
        if (groupId !== dashBoardData?.notifications[0]?.payload.groupId)
          return;
        setGroupData((prev) => ({
          ...prev,
          data: {...prev.data, Status: "invited"},
        }));
        break;

      case "notification createEvent":
        if (
          dashBoardData?.notifications[0]?.payload?.EventCard.event?.GroupId !==
          groupId
        )
          return;
        setGroupData((prev) => {
          return {
            ...prev,
            data: {
              ...prev?.data,
              Events: Array.isArray(prev?.data.Events)
                ? [
                    ...prev?.data.Events,
                    dashBoardData?.notifications[0]?.payload?.EventCard,
                  ]
                : [dashBoardData?.notifications[0]?.payload?.EventCard],
            },
          };
        });

        break;

      default:
        console.log(
          "unknow nofication type",
          dashBoardData?.notifications,
          "hehhere"
        );
        break;
    }
  }, [
    dashBoardData?.notifications,
    profileData?.isComponentVisible,
    fetchProfileData,
    profileData?.data?.Owner?.PrivacySetting,
    groupId,
  ]);

  useEffect(() => {
    handleWebsocketNotification();
  }, [dashBoardData?.notifications, handleWebsocketNotification]);

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
          if (profileData?.data?.Owner?.PrivacySetting === "public") {
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
        case "requestToJoinGroup":
          setGroupData((prev) => ({
            ...prev,
            data: {
              ...prev.data,
              Status: "requested",
            },
          }));

          break;
        case "answerRequestToJoinGroup":
          if (data.whatever.accept) {
            setGroupData((prev) => {
              const user = prev?.data?.RequestedMembers?.find(
                (el) => el.Id === data.whatever.applicantId
              );

              return {
                ...prev,
                data: {
                  ...prev.data,
                  Members: Array.isArray(prev?.data?.Members)
                    ? [user, ...prev?.data?.Members]
                    : [user],

                  RequestedMembers: prev?.data?.RequestedMembers.filter(
                    (el) => el.Id !== data.whatever.applicantId
                  ),
                },
              };
            });
          } else {
            setGroupData((prev) => ({
              ...prev,
              data: {
                ...prev.data,
                RequestedMembers: prev?.data?.RequestedMembers.filter(
                  (el) => el.Id !== data.whatever.applicantId
                ),
              },
            }));
          }

          break;
        case "inviteToJoinGroup":
          if (groupId !== data.whatever.groupId) return;
          setGroupData((prev) => ({
            ...prev,
            data: {
              ...prev.data,
              Invite: prev.data.Invite.map((el) => {
                if (el.Id === data.whatever.receiverId) {
                  return {...el, isInvited: true};
                }
                return el;
              }),
            },
          }));

          break;
        case "answerInvitationToJoinGroup":
          if (groupId !== data.whatever) return;
          fetchGroupData(groupId);
          break;

        case "createEvent":
          if (data.whatever.payload.event.GroupId !== groupId) return;
          setGroupData((prev) => {
            return {
              ...prev,
              data: {
                ...prev?.data,
                Events: Array.isArray(prev?.data.Events)
                  ? [...prev?.data.Events, data.whatever.payload]
                  : [data.whatever.payload],
              },
            };
          });
          break;

        default:
          console.log("unknown websocket message");
          break;
      }
    },
    [
      user.Id,
      profileData?.data?.Owner?.PrivacySetting,
      fetchProfileData,
      profileData.data?.Owner?.Id,
      user.Avatar,
      user.FirstName,
      user.LastName,
      user.Nickname,
      user.Privicy_setting,
      groupId,
      fetchGroupData,
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
      profileData?.isComponentVisible,
      fetchProfileData,
      profileData?.data?.Owner?.Id,
      user?.Id,
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
        profileData?.isComponentVisible &&
        CreatorId === profileData?.data?.Owner?.Id
      ) {
        setProfileData((prev) => ({
          ...prev,
          data: {
            ...prev.data,
            Posts: Array.isArray(prev?.data?.Posts)
              ? [
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
                  ...prev?.data?.Posts,
                ]
              : [
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
                ],
          },
        }));
      } else if (groupId === GroupId) {
        setGroupData((prev) => ({
          ...prev,
          data: {
            ...prev?.data,
            Posts: Array.isArray(prev?.data.Posts)
              ? [
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
                  ...prev?.data.Posts,
                ]
              : [
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
      profileData?.isComponentVisible,
      profileData.data?.Owner?.Id,
      setProfileData,
      setDashBoardData,
      groupId,
    ]
  );

  useEffect(() => {
    if (isWsReady) {
      const data = JSON.parse(wsVal);
      console.log(data, "data from websocket", Date.now());
      if (data?.type === "error") {
        handleWebsocketErrors(data);
        return;
      }

      if (data?.type === "success") {
        handleWebsocketSucess(data);
        return;
      }

      if (
        data?.type?.includes("notification")
        // dashBoardData.notifications.some((el) => el.Id === data.Id) === false
      ) {
        setDashBoardData((prev) => ({
          ...prev,
          notifications: Array.isArray(prev?.notifications)
            ? [data, ...prev?.notifications]
            : [data],
        }));
        return;
      }

      switch (data?.type) {
        case "logout":
          httpLogout();
          setUser(userObj);
          break;
        case "groupMessage":
          if (groupId === data.GroupId) {
            setOpenChatDetails((prev) => ({
              ...prev,
              messages: Array.isArray(prev.messages)
                ? [...prev.messages, data]
                : [data],
            }));
          }
          break;

        case "privateMessage":
          setOpenChatDetails((prev) => ({
            ...prev,
            messages: Array.isArray(prev.messages)
              ? [...prev.messages, data]
              : [data],
          }));
          break;
        case "online-user":
          setOnlineUsers(data.data);

          break;
        case "post":
          addPostToFeed(data);
          break;
        case "togglePrivacySetting":
          currentProfilePrivacyChanged(data);
          break;
        case "createGroup":
          const obj = {
            BasicInfo: data.payload.Data,
            Status:
              data.payload.Data.CreatorId === user.Id ? "accepted" : "none",
          };
          setDashBoardData((prev) => {
            return {
              ...prev,
              groups: Array.isArray(prev.groups)
                ? [...prev.groups, obj]
                : [obj],
            };
          });

          break;
        case "answerInvitationToJoinGroup":
          if (data.groupId !== groupId) return;
          console.log("answerInvitationToJoinGroup");
          setGroupData((prev) => {
            return {
              ...prev,
              data: {
                ...prev?.data,
                Members: Array.isArray(prev?.data.Members)
                  ? [...prev?.data.Members, data?.newMember]
                  : [data?.newMember],
              },
            };
          });

          break;
        case "toggleAttendEvent":
          if (groupId !== data.GroupId) return;
          // console.log(data, "we are making here");
          setGroupData((prev) => {
            return {
              ...prev,
              data: {
                ...prev.data,
                Events: prev.data.Events.map((el) => {
                  if (data.EventId === el.event.Id) {
                    return {
                      event: {
                        ...el.event,
                        Going: data.Going,
                        NotGoing: data.NotGoing,
                      },

                      Going: data.hasOwnProperty("IsAttending")
                        ? data.IsAttending
                        : el.Going,
                    };
                  }
                  return el;
                }),
              },
            };
          });

          break;

        default:
          // console.log("no action", data);
          break;
      }
    }
  }, [
    isWsReady,
    wsVal,
    // logoutHandler,
    // addPostToFeed,
    // currentProfilePrivacyChanged,
    // handleWebsocketErrors,
    // handleWebsocketSucess,

    user.Id,
    groupId,
  ]);

  const onAddLikeDislikePost = (id, data) => {
    if (profileData?.isComponentVisible) {
      setProfileData((prev) => ({
        ...prev,
        data: {
          ...prev.data,
          Posts: prev.data.Posts.map((el) =>
            el.Id === id
              ? {
                  ...el,
                  Likes: data.Likes,
                  Dislikes: data.Dislikes,
                  UserLikeDislike: data.UserLikeDislike,
                }
              : el
          ),
        },
      }));
    } else {
      setDashBoardData((prev) => ({
        ...prev,
        Posts: prev.Posts.map((el) =>
          el.Id === id
            ? {
                ...el,
                Likes: data.Likes,
                Dislikes: data.Dislikes,
                UserLikeDislike: data.UserLikeDislike,
              }
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
        fetchGroupData,
        groupData,
        resetIsGroupComponentVisible,
        updateGroupId,
        groupId,
      }}
    >
      {props.children}
    </AuthContext.Provider>
  );
};
export default AuthContext;
