// YourComponent.js
import React, {useState, useRef, useEffect} from "react";
import classes from "./Profile.module.css";
import GroupCard from "../components/GroupCard";
import AuthContext from "../store/authContext";
import {Container} from "react-bootstrap";
import GroupActions from "../components/GroupActions";
import Posts from "../components/Posts";
import NoMember from "../components/NotMember";
import {useRouteError, useParams} from "react-router";
import PostInput from "../components/PostInput";

const Group = () => {
  const {
    user,
    fetchGroupData,
    groupData,
    wsMsgToServer,
    isWsReady,
    openChat,
    resetIsGroupComponentVisible,
    updateGroupId,
  } = React.useContext(AuthContext);
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

  const groupActionHandler = () => {
    if (isWsReady) {
      if (groupData?.data.Status === "none") {
        console.log(groupData?.data.Status);
        wsMsgToServer(
          JSON.stringify({
            Type: "requestToJoinGroup",
            message: {
              userId: user.Id,
              GroupId: groupData?.data?.BasicInfo.Id,
            },
          })
        );
        return;
      }

      wsMsgToServer(
        JSON.stringify({
          type: "answerInvitationToJoinGroup",
          message: {
            UserId: user.Id,
            GroupId: groupData?.data?.BasicInfo.Id,
          },
        })
      );
    }
  };

  useEffect(() => {
    resetIsGroupComponentVisible(true);
    return () => resetIsGroupComponentVisible(false);
  }, [resetIsGroupComponentVisible]);
  useEffect(() => {
    updateGroupId(id);
    return () => updateGroupId("");
  }, [updateGroupId, id]);

  const AcceptRefuseGroupRequestHandler = ({id, type}) => {
    if (isWsReady) {
      console.log(type, "type");
      wsMsgToServer(
        JSON.stringify({
          Type: "answerRequestToJoinGroup",
          message: {
            SenderId: user.Id,
            GroupId: groupData?.data?.BasicInfo.Id,
            Accept: type,
            ReceiverId: id,
          },
        })
      );
    }
  };

  const inviteHandler = (id) => {
    console.log(id, "this where sending the invite ", isWsReady);
    if (isWsReady) {
      wsMsgToServer(
        JSON.stringify({
          Type: "inviteToJoinGroup",
          message: {
            SenderId: user.Id,
            GroupId: groupData?.data?.BasicInfo.Id,
            ReceiverId: id,
          },
        })
      );
    }
  };

  console.log("this message");
  return routeError ? (
    routeError.message
  ) : (
    <Container>
      <div className={classes.profile}>
        {groupData?.data?.Status !== "accepted" ? (
          <NoMember
            Title={groupData?.data?.BasicInfo?.Title}
            status={groupData?.data?.Status}
            GroupAction={groupActionHandler}
          />
        ) : (
          <>
            <GroupCard
              group={groupData.data}
              toggleAction={toggleAction}
              owner={owner}
              showChat={openChat}
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
              inviteHandler={inviteHandler}
              AcceptRefuseGroupRequestHandler={AcceptRefuseGroupRequestHandler}
            />

            <div>
              <div style={{width: "50vw", margin: "2rem 0 3rem 0"}}>
                <PostInput src={user.Avatar} id={user.Id} />
              </div>
              <section id="posts" ref={postRef}>
                {<Posts posts={groupData.data.Posts} postref={postRef} />}
              </section>
            </div>
          </>
        )}
      </div>
    </Container>
  );
};

export default Group;
