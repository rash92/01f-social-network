// YourComponent.js
import React, {useState, useRef, useEffect} from "react";
import classes from "./Profile.module.css";
import ProfileCard from "../components/ProfileCard";
import AuthContext from "../store/authContext";
import {Container} from "react-bootstrap";
import ProfileActions from "../components/ProfileActions";
import Posts from "../components/Posts";
import {dummyPosts} from "../store/dummydata";
import {getJson} from "../helpers/helpers";
import InfiniteScroll from "react-infinite-scroll-component";
import PrivateProfile from "../components/PrivateProfile";
import {useLoaderData, useRouteError} from "react-router";
import PostInput from "../components/PostInput";


const Profile = () => {
  const {user, isWsReady, wsMsgToServer} = React.useContext(AuthContext);
  const userData = useLoaderData();
  const routeError = useRouteError();
  const postRef = useRef(null);
  const [show, setShow] = useState(false);
  const [isActive, setIsActive] = useState("Posts");
  const [data, setData] = useState(userData);
  // const [hasMoreFellowers, setHasMoreFellowers] = useState({
  //   followers: true,
  //   following: true,
  // });

  // const dummyUser = {
  //   id: "4e56t78y98u0ii9i90i",
  //   relStatus: "fellow",
  //   isOwner: true,
  //   Isprivate: true,
  //   followers: options1,
  //   following: options2,
  //   numberOfFollowers: 60,
  //   numberOfFollowing: 60,
  //   numberOfPosts: 60,
  //   posts: dummyPosts,
  //   firstName: "John",
  //   lastName: "Doe",
  //   dateOfBirth: "1990-01-01",
  //   avatar: user.profileImg,
  //   nickname: "JD",
  //   aboutMe:
  //     "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur et tristique libero.",
  // };

  const isPrivate =
    data?.Owner.PrivacySetting === "private" ||
    data?.Owner.PrivacySetting === ""
      ? true
      : false;

  // const [hasMorePosts, setHasMorePosts] = useState(true);
  const toggleProfileVisibility = () => {};

  const toggleActionModel = (active) => {
    setIsActive(active);
  };
  const handleClose = () => setShow(false);
  const handleShow = () => {
    setShow(true);
  };

  const toggleAction = (clickButton, e) => {
    console.log(clickButton);
    if (clickButton === "Posts") {
      postRef.current.scrollIntoView({behavior: "smooth"});
      return;
    }
    setIsActive(clickButton.toLowerCase());
  };

  useEffect(() => {
    if (!(isActive === "Posts")) {
      handleShow();
    }

    return () => {
      setIsActive("Posts");
    };
  }, [isActive]);

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

  const requestFollowHandler = async () => {
    try {
      if (isWsReady) {
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
    } catch (err) {
      console.log(err);
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
          // <InfiniteScroll
          //   dataLength={data.posts?.length || 0}
          //   next={fetchMorePosts}
          //   hasMore={hasMorePosts}
          //   loader={<h4>Loading...</h4>}
          //   endMessage={
          //     <p style={{textAlign: "center"}}>
          //       <b>Yay! You have seen it all</b>
          //     </p>
          //   }
          // >
          <>
            <ProfileCard
              user={data.Owner}
              toggleAction={toggleAction}
              toggleProfileVisibility={toggleProfileVisibility}
              isPrivate={isPrivate}
              owner={user.Id === data.Owner.Id}
              numberOfFollowers={data.NumberOfFollowers}
              numberOfFollowing={data.NumberOfFollowing}
              numberOfPosts={data.NumberOfPosts}
              numberOfRequests={1}
            />
            <ProfileActions
              user={data}
              handleClose={handleClose}
              handleShow={handleShow}
              show={show}
              flag={false}
              isActive={isActive}
              owner={user.Id === data.Owner.Id}
              toggleAction={toggleActionModel}
              // fetchMoreFellowers={fetchMoreFellowers}
              // hasMoreFellowers={hasMoreFellowers}
            />

            <div>
              <div style={{width: "50vw", margin: "2rem 0 3rem 0"}}>
                <PostInput src={user.Profile} id={user.Id} />
              </div>
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

export async function profileLoader({request, params}) {
  return getJson("profile", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(params.id),
  });
}

export default Profile;
