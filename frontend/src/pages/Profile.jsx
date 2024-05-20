// YourComponent.js
import React, {useState, useRef, useEffect, useCallback} from "react";
import classes from "./Profile.module.css";
import ProfileCard from "../components/ProfileCard";
import AuthContext from "../store/authContext";
import {Container} from "react-bootstrap";
import ProfileActions from "../components/ProfileActions";
import Posts from "../components/Posts";

import PrivateProfile from "../components/PrivateProfile";
import {useRouteError, useParams} from "react-router";
import PostInput from "../components/PostInput";

const Profile = () => {
  const {
    user,
    isWsReady,

    wsMsgToServer,
    toggleProfilePrivacy,
    profileData: {data, error},
    resetIsProfileComponentVisible,
    fetchProfileData,
  } = React.useContext(AuthContext);

  const {id} = useParams();

  console.log(data.Posts, "posts");
  useEffect(() => {
    fetchProfileData(id);
  }, [id, fetchProfileData]);

  const routeError = useRouteError();
  const postRef = useRef(null);
  const [show, setShow] = useState(false);
  const [isActive, setIsActive] = useState("Posts");

  const isPrivate =
    data?.Owner?.PrivacySetting === "private" ||
    data?.Owner?.PrivacySetting === ""
      ? true
      : false;

  // const [hasMorePosts, setHasMorePosts] = useState(true);
  const toggleProfileVisibility = () => {
    // this need to be romeved when when handle this in the websocket
    // toggleProfilePrivacy();
    // this the code when add code the backend websocket
    if (isWsReady) {
      wsMsgToServer(
        JSON.stringify({
          Type: "togglePrivacy",
          message: {
            SenderId: user.Id,
            PrivacySetting:
              data?.Owner?.PrivacySetting === "private" ? "public" : "private",
          },
        })
      );

      return;
    }
  };

  const toggleActionModel = (active) => {
    setIsActive(active);
  };
  const handleClose = () => setShow(false);
  const handleShow = () => {
    setShow(true);
  };

  const toggleAction = (clickButton, e) => {
    if (clickButton === "Posts") {
      postRef.current.scrollIntoView({behavior: "smooth"});
      return;
    }

    handleShow();
    setIsActive(clickButton);
  };

  useEffect(() => {
    resetIsProfileComponentVisible(true);
    return () => resetIsProfileComponentVisible(false);
  }, [resetIsProfileComponentVisible]);

  // useEffect(() => {
  //   if (!isComponentVisible || !isWsReady) return;
  //   const data = JSON.parse(wsVal);
  //   console.log(data);

  //   switch (data?.message) {
  //     case "requestToFollow":
  //       setData((prw) => ({...prw, IsPending: true}));
  //       break;
  //     default: {
  //       break;
  //     }
  //   }

  //   //  nofications actions
  //   // console.log(data.Body.Data);
  //   if (data?.Type === "notification requestToFollow") {
  //     setData((prev) => {
  //       if (Object.keys(prev).length === 0) {
  //         console.log(" not prev", prev);
  //         return {};
  //       }

  //       let arr = prev?.PendingFollowers;

  //       if (!arr) {
  //         console.log(" not arr", prev);
  //         return {};
  //       }
  //       console.log("prev", prev);
  //       return {
  //         ...prev,
  //         PendingFollowers: [...arr, data?.Body?.Data],
  //       };
  //     });
  //   }
  // }, [wsVal, isWsReady, isComponentVisible]);

  // useEffect(() => {
  //   console.log(isActive);
  //   if (!(isActive === "Posts")) {
  //     handleShow();
  //   }

  //   // return () => {
  //   //   setIsActive("Posts");
  //   // };
  // }, [isActive]);

  // const fetchMoreFellowers = () => {
  //   if (
  //     data[isActive]?.length <
  //     data[`NumberOf${isActive.charAt(0).toUpperCase()}${isActive.slice(1)}`]
  //   ) {
  //     setData((prev) => ({
  //       ...prev,
  //       [isActive]: [...prev[isActive], ...options1],
  //     }));
  //   } else {
  //     setHasMoreFellowers((pre) => ({...pre, [isActive]: false}));
  //   }
  // };

  // const fetchMorePosts = () => {
  //   if (data.Posts.length <= 50) {
  //     setData((prev) => ({
  //       ...prev,
  //       Posts: [...prev.Posts, ...dummyPosts],
  //     }));
  //     console.log(data);
  //   } else {
  //     console.log(data);
  //     setHasMorePosts(false);
  //   }
  // };

  const requestFollowHandler = () => {
    console.log("unfollowing");
    if (isWsReady) {
      if (data.IsFollowed) {
        wsMsgToServer(
          JSON.stringify({
            Type: "unfollow",
            message: {
              FollowerId: user.Id,
              FollowingId: data.Owner.Id,
            },
          })
        );

        return;
      }

      wsMsgToServer(
        JSON.stringify({
          type: "requestToFollow",
          message: {
            FollowerId: user.Id,
            FollowingId: data.Owner.Id,
            Status: "pending",
          },
        })
      );
    }
  };

  const accepOrRejectRequestHandler = ({id, flag}, e) => {
    if (isWsReady) {
      wsMsgToServer(
        JSON.stringify({
          Type: "answerRequestToFollow",
          message: {
            SenderId: user.Id,
            ReceiverId: id,
            Reply: flag,
          },
        })
      );
    }
  };

  return routeError ? (
    routeError.message
  ) : (
    <Container>
      <div className={classes.profile}>
        {isPrivate && data.Owner.Id !== user.Id && !data.IsFollowed ? (
          <PrivateProfile
            IsPending={data.IsPending}
            Avatar={data?.Owner?.avatar}
            IsOnline={true}
            name={"bilal"}
            Nickname={data.Owner.Nickname}
            fellowUserHandler={requestFollowHandler}
          />
        ) : (
          <>
            {data.Owner && (
              <ProfileCard
                data={data}
                toggleAction={toggleAction}
                toggleProfileVisibility={toggleProfileVisibility}
                isPrivate={isPrivate}
                owner={user?.Id === data?.Owner?.Id}
                requestFollow={requestFollowHandler}
              />
            )}

            {data.Owner && (
              <ProfileActions
                data={data}
                active={data[isActive]}
                handleClose={handleClose}
                handleShow={handleShow}
                show={show}
                flag={false}
                isActive={isActive}
                owner={user?.Id === data?.Owner?.Id}
                toggleAction={toggleActionModel}
                accepOrRejectRequestHandler={accepOrRejectRequestHandler}

                // fetchMoreFellowers={fetchMoreFellowers}
                // hasMoreFellowers={hasMoreFellowers}
              />
            )}

            <div>
              {user?.Id === data?.Owner?.Id && (
                <div style={{width: "50vw", margin: "2rem 0 3rem 0"}}>
                  <PostInput src={user.Profile} id={user.Id} />
                </div>
              )}

              <section id="posts" ref={postRef}>
                {<Posts posts={data.Posts} postref={postRef} />}
              </section>
            </div>
          </>
          // </InfiniteScroll>
        )}
      </div>
    </Container>
  );
};

// export async function profileLoader({request, params}) {
//   return getJson("profile", {
//     method: "POST",
//     headers: {
//       "Content-Type": "application/json",
//     },
//     credentials: "include",
//     body: JSON.stringify(params.id),
//   });
// }

export default Profile;
