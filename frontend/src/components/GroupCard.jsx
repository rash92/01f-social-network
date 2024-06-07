// ProfileCard.js

import {Card, Button} from "react-bootstrap";
import classes from "./ProfileCard.module.css";
import Action from "./Action";
import AuthContext from "../store/authContext";
import {useContext} from "react";

const GroupCard = ({group, toggleAction, owner, showChat}) => {
  const {user} = useContext(AuthContext);

  return (
    <Card className={classes.card}>
      <div className={classes.profileContainer}>
        <Card.Body>
          <div className={classes.image}>
            <Card.Title>
              <span>{group?.BasicInfo?.Title}</span>
            </Card.Title>
            <ul className={`${classes["profile-action"]}   `}>
              <Action
                numberAction={group?.Posts?.length || 0}
                actionName={"Posts"}
                toggleAction={toggleAction}
              />

              <Action
                numberAction={group?.Members?.length || 0}
                actionName={"Members"}
                toggleAction={toggleAction}
              />

              <Action
                numberAction={group?.Events?.length || 0}
                actionName={"Events"}
                toggleAction={toggleAction}
              />
            </ul>
          </div>

          <Card.Text>{group?.BasicInfo?.Description}</Card.Text>
          <div className={classes.interact}>
            <Button
              onClick={toggleAction?.bind(
                null,
                `${owner ? "RequestedMembers" : "create event"}`
              )}
            >
              manage
            </Button>

            <Button
              onClick={showChat.bind(null, {
                id: group?.BasicInfo?.Id,
                Nickname: user.Nickname,
                Avatar: user.Avatar,
                type: "groupMessage",
              })}
            >
              chat
            </Button>
          </div>
        </Card.Body>
      </div>
    </Card>
  );
};

export default GroupCard;
