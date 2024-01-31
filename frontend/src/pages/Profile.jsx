// YourComponent.js
import React, {useState, useRef} from "react";
import classes from "./Profile.module.css";
import ProfileCard from "../components/ProfileCard";
import AuthContext from "../store/authContext";
import {Container} from "react-bootstrap";
import ProfileActions from "../components/ProfileActions";
import Posts from "../components/Posts";
import {dummyPosts} from "../store/dummydata";

const options1 = [
  {id: 1, username: "abdi2", name: "abdi"},
  {id: 2, username: "ahmed34", name: "ahmed"},
  {id: 3, username: "user3", name: "User 3"},
  {id: 4, username: "user4", name: "User 4"},
  {id: 5, username: "user5", name: "User 5"},
];

const options2 = [
  {id: 6, username: "john_doe", name: "John"},
  {id: 7, username: "jane_smith", name: "Jane"},
  {id: 8, username: "user8", name: "User 8"},
  {id: 9, username: "user9", name: "User 9"},
  {id: 10, username: "user10", name: "User 10"},
];

const Profile = () => {
  const {user} = React.useContext(AuthContext);
  const [Isprivate, setIsPrivate] = useState(true);
  const postRef = useRef(null);
  const [show, setShow] = useState(false);
  const [isActive, setIsActive] = useState("");
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
    setIsActive("clickButton");
    
  };

  const toggleProfileVisibility = () => {
    setIsPrivate(!Isprivate);
  };

  const dummyUser = {
    toggleProfileVisibility,
    relStatus: "fellow",
    isOwner: true,
    Isprivate: true,
    followers: options1,
    following: options2,
    posts: dummyPosts,
    firstName: "John",
    lastName: "Doe",
    dateOfBirth: "1990-01-01",
    avatar: user.profileImg,
    nickname: "JD",
    aboutMe:
      "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur et tristique libero.",
  };

  return (
    <Container>
      <div className={classes.profile}>
        <ProfileCard user={dummyUser} toggleAction={toggleAction} />

        <ProfileActions
          user={dummyUser}
          handleClose={handleClose}
          handleShow={handleShow}
          show={show}
          flag={false}
          isActive={isActive}
          toggleAction={toggleActionModel}
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

export default Profile;
