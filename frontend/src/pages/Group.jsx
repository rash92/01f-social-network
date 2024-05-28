// YourComponent.js
import React, {useState, useRef, useEffect} from "react";
import classes from "./Profile.module.css";
import GroupCard from "../components/GroupCard";
import AuthContext from "../store/authContext";
import {Container} from "react-bootstrap";
import GroupActions from "../components/GroupActions";
import Posts from "../components/Posts";
import {useRouteError, useParams} from "react-router";

const Group = () => {
  const {user, fetchGroupData, groupData} = React.useContext(AuthContext);
  const routeError = useRouteError();
  const postRef = useRef(null);
  const [show, setShow] = useState(false);
  const [isActive, setIsActive] = useState("");

  const status = isActive === "Members" || isActive === "Events" ? true : false;
  console.log(groupData);
  const toggleActionModel = (active) => {
    setIsActive(active);
  };
  const handleClose = () => {
    setShow(false);
    setIsActive("");
  };
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
  console.log(isActive);
  const owner = user.Id === groupData?.data?.BasicInfo?.CreatorId;
  const {id} = useParams();
  useEffect(() => {
    fetchGroupData(id);
  }, [id, fetchGroupData]);

  return routeError ? (
    routeError.message
  ) : (
    <Container>
      <div className={classes.profile}>
        <GroupCard
          group={groupData.data}
          toggleAction={toggleAction}
          owner={owner}
        />
        <GroupActions
          group={groupData.data}
          handleClose={handleClose}
          handleShow={handleShow}
          show={show}
          flag={false}
          status={status}
          isActive={isActive}
          toggleAction={toggleActionModel}
          owner={owner}
        />

        <div>
          <section id="posts" ref={postRef}>
            {<Posts posts={groupData.data.Posts} postref={postRef} />}
          </section>
        </div>
      </div>
    </Container>
  );
};

export default Group;
