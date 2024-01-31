// YourComponent.js
import React, {useState, useRef} from "react";
import classes from "./Profile.module.css";
import ProfileCard from "../components/ProfileCard";
import AuthContext from "../store/authContext";
import {Container} from "react-bootstrap";

import Posts from "../components/Posts";
import {dummyPosts} from "../store/dummydata";

const Profile = () => {
  const {user} = React.useContext(AuthContext);
  const [Isprivate, setIsPrivate] = useState(true);
  const postRef = useRef(null);
  const toggleAction = (clickButton, e) => {
    if (clickButton === "Posts") {
      postRef.current.scrollIntoView({behavior: "smooth"});
    }

    
  };
  const toggleProfileVisibility = () => {
    setIsPrivate(!Isprivate);
  };

  const dummyUser = {
    toggleProfileVisibility,
    relStatus: "fellow",
    isOwner: true,
    Isprivate: true,
    followers: 20,
    following: 35,
    posts: 23,
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

        <div>
          {/* <ul className={classes["profile-action"]}>
            <Action
              numberAction={dummyUser.posts}
              actionName={"Posts"}
              toggleAction={toggleAction}
              active={isActive}
            />
            <Action
              numberAction={dummyUser.followers}
              actionName={"Followers"}
              toggleAction={toggleAction}
              active={isActive}
            />
            <Action
              numberAction={dummyUser.following}
              actionName={"Following"}
              toggleAction={toggleAction}
              active={isActive}
            />
          </ul> */}
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
