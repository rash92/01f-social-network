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
  const dummyUser = {
    // toggleProfileVisibility,
    relStatus: "fellow",
    isOwner: true,
    Isprivate: true,
    followers: options1,
    following: options2,
    numberOfFollowers: 60,
    numberOfFollowing: 60,
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

  const [data, setData] = useState(dummyUser);
  const [hasMore, setHasMore] = useState({followers: true, following: true});

  const fetchMoreData = () => {
    if (data[isActive].length <= 50) {
      setData((prev) => ({
        ...prev,
        [isActive]: [...prev[isActive], ...options1],
      }));
    } else {
      setHasMore((pre) => ({...pre, [isActive]: false}));
    }
  };

  return (
    <Container>
      <div className={classes.profile}>
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
          fetchMoreData={fetchMoreData}
          hasMore={hasMore}
        />

        <div>
          <section id="posts" ref={postRef}>
            <h4 style={{textAlign: "center"}}>Posts</h4>
            {<Posts posts={dummyPosts} postref={postRef} />}
          </section>
        </div>
      </div>
    </Container>
  );
};

export function profileLoader({request, params}) {
  // console.log(request, params);
  return getJson("profile", {
    method: "POST",
    mode: "no-cors",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(params.id),
  });
}

export default Profile;
