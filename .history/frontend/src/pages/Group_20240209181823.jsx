// YourComponent.js
import React, {useState, useRef} from "react";
import classes from "./Profile.module.css";
import GroupCard from "../components/GroupCard";
import AuthContext from "../store/authContext";
import {Container} from "react-bootstrap";
import ProfileActions from "../components/ProfileActions";
import Posts from "../components/Posts";
import {dummyPosts, group} from "../store/dummydata";
import {getJson} from "../helpers/helpers";
import InfiniteScroll from "react-infinite-scroll-component";
import {useLoaderData, useRouteError} from "react-router";

const Group = () => {
  const {user} = React.useContext(AuthContext);
  const routeError = useRouteError();
  const postRef = useRef(null);
  const [show, setShow] = useState(false);
  const [isActive, setIsActive] = useState("");
  const [data, setData] = useState(group);
  const [hasMoreActions, setHasMoreAction] = useState({
    members: true,
    events: true,
  });

  // const isprivate = data.Owner.privacySetting === "private" ? true : false;

  const [hasMorePosts, setHasMorePosts] = useState(true);
  const toggleProfileVisibility = () => {};

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
    if (
      data[`NumberOf${isActive.charAt(0).toUpperCase()}${isActive.slice(1)}`]
    ) {
      handleShow();
    }
  };

  const fetchMoreFellowers = () => {
    if (
      data[isActive].length <
      data[`NumberOf${isActive.charAt(0).toUpperCase()}${isActive.slice(1)}`]
    ) {
      setData((prev) => ({
        ...prev,
        [isActive]: [...prev[isActive], ...[]],
      }));
    } else {
      setHasMoreAction((pre) => ({...pre, [isActive]: false}));
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

  return routeError ? (
    routeError.message
  ) : (
    <Container>
      <div className={classes.profile}>
        <InfiniteScroll
          dataLength={data.posts?.length || 0}
          next={fetchMorePosts}
          hasMore={hasMorePosts}
          loader={<h4>Loading...</h4>}
          endMessage={
            <p style={{textAlign: "center"}}>
              <b>Yay! You have seen it all</b>
            </p>
          }
        >
          <GroupCard
            user={data}
            toggleAction={toggleAction}
            toggleProfileVisibility={toggleProfileVisibility}
            owner={user.id === data.id}
            numberOfFollowers={data.NumberOfFollowers}
            numberOfFollowing={data.NumberOfFollowing}
            numberOfPosts={data.NumberOfPosts}
          />
          {/* <ProfileActions
            user={data}
            handleClose={handleClose}
            handleShow={handleShow}
            show={show}
            flag={false}
            isActive={isActive}
            toggleAction={toggleActionModel}
            fetchMoreFellowers={fetchMoreFellowers}
            hasMoreFellowers={hasMoreActions}
          /> */}

          {/* <div>
              <section id="posts" ref={postRef}>
                {<Posts posts={data.Posts} postref={postRef} />}
              </section>
            </div> */}
        </InfiniteScroll>
      </div>
    </Container>
  );
};

export async function GroupLoader({request, params}) {
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

export default Group;