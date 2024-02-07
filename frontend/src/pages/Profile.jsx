// YourComponent.js
import React, {useState, useRef} from "react";
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
import {useRouteLoaderData, useRouteError} from "react-router";

const options1 = [
  {id: 1, username: "user", name: "User 1"},
  {id: 2, username: "user", name: "User 2"},
  {id: 3, username: "user", name: "User 3"},
  {id: 4, username: "user", name: "User 4"},
  {id: 5, username: "user", name: "User 5"},
  {id: 6, username: "user", name: "User 6"},
  {id: 7, username: "user", name: "User 7"},
  {id: 8, username: "user", name: "User 8"},
  {id: 9, username: "user", name: "User 9"},
  {id: 10, username: "user", name: "User 10"},
];

const options2 = [
  {id: 1, username: "user", name: "User 1"},
  {id: 2, username: "user", name: "User 2"},
  {id: 3, username: "user", name: "User 3"},
  {id: 4, username: "user", name: "User 4"},
  {id: 5, username: "user", name: "User 5"},
  {id: 6, username: "user", name: "User 6"},
  {id: 7, username: "user", name: "User 7"},
  {id: 8, username: "user", name: "User 8"},
  {id: 9, username: "user", name: "User 9"},
  {id: 10, username: "user", name: "User 10"},
];

const Profile = () => {
  const {user} = React.useContext(AuthContext);
  const userData = useRouteLoaderData();
  const routeError = useRouteError();

  console.log(routeError, "useRouteError");
  const dummyUser = {
    id: "4e56t78y98u0ii9i90i",
    relStatus: "fellow",
    isOwner: true,
    Isprivate: true,
    followers: options1,
    following: options2,
    numberOfFollowers: 60,
    numberOfFollowing: 60,
    numberOfPosts: 60,
    posts: dummyPosts,
    firstName: "John",
    lastName: "Doe",
    dateOfBirth: "1990-01-01",
    avatar: user.profileImg,
    nickname: "JD",
    aboutMe:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur et tristique libero.",
  };

  const [Isprivate, setIsPrivate] = useState(true);
  const postRef = useRef(null);
  const [show, setShow] = useState(false);
  const [isActive, setIsActive] = useState("");
  const [data, setData] = useState(dummyUser);
  const [hasMoreFellowers, setHasMoreFellowers] = useState({
    followers: true,
    following: true,
  });
  const [hasMorePosts, setHasMorePosts] = useState(true);
  const toggleProfileVisibility = () => {
    setIsPrivate(!Isprivate);
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
    setIsActive(clickButton.toLowerCase());
    handleShow();
  };

  const fetchMoreFellowers = () => {
    if (data[isActive].length <= 50) {
      setData((prev) => ({
        ...prev,
        [isActive]: [...prev[isActive], ...options1],
      }));
    } else {
      setHasMoreFellowers((pre) => ({...pre, [isActive]: false}));
    }
  };

  const fetchMorePosts = () => {
    if (data.posts.length <= 50) {
      setData((prev) => ({
        ...prev,
        posts: [...prev.posts, ...dummyPosts],
      }));
      console.log(data);
    } else {
      console.log(data);
      setHasMorePosts(false);
    }
  };

  //  const getUserData = ()=>{
  //   try{}
  //  }

  return routeError ? (
    routeError.message
  ) : (
    <Container>
      <div className={classes.profile}>
        {!data.Isprivate &&
        data.id !== user.id &&
        data.relStatus === "fellow" ? (
          <PrivateProfile
            name={"bilal"}
            userName={user.username}
            fellowUserHandler={() => {}}
          />
        ) : (
          <InfiniteScroll
            dataLength={data.posts?.length}
            next={fetchMorePosts}
            hasMore={hasMorePosts}
            loader={<h4>Loading...</h4>}
            endMessage={
              <p style={{textAlign: "center"}}>
                <b>Yay! You have seen it all</b>
              </p>
            }
          >
            <ProfileCard
              user={data}
              toggleAction={toggleAction}
              toggleProfileVisibility={toggleProfileVisibility}
            />
            <ProfileActions
              user={data}
              handleClose={handleClose}
              handleShow={handleShow}
              show={show}
              flag={false}
              isActive={isActive}
              toggleAction={toggleActionModel}
              fetchMoreFellowers={fetchMoreFellowers}
              hasMoreFellowers={hasMoreFellowers}
            />

            <div>
              <section id="posts" ref={postRef}>
                {<Posts posts={data.posts} postref={postRef} />}
              </section>
            </div>
          </InfiniteScroll>
        )}
      </div>
    </Container>
  );
};

export async function profileLoader({request, params}) {
  // console.log(params.id, request);
  return getJson("profile", {
    // signal: request.signal,/
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params.id),
  });
}

export default Profile;
