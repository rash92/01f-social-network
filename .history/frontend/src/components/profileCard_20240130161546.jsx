// ProfileCard.js
import React from "react";
import {Card, Form, Button} from "react-bootstrap";
import RoundImage from "./RoundedImg";
import classes from "./ProfileCard.module.css";
import Action from "./Action";

const ProfileCard = ({user}) => {
  const toggleAction = () => {};
  return (
    <Card className={classes.card}>
      <div className={classes.profileContainer}>
        {user.isOwner ? (
          <Form>
            <Form.Check
              type="switch"
              id="profileVisibilitySwitch"
              label={`Make Profile ${user.IsPrivate ? "Public" : "Private"}`}
              checked={user.IsPrivate}
              onChange={user.toggleProfileVisibility}
            />
          </Form>
        ) : (
          <div>
            <Button>{user.relStatus}</Button>
          </div>
        )}
        <Card.Body className={classes.body}>
          <ul className={classes["profile-action"]}>
            <Action
              numberAction={user.followers}
              actionName={"Followers"}
              toggleAction={toggleAction}
            />
            <Action
              numberAction={user.fellowing}
              actionName={"Following"}
              toggleAction={toggleAction}
            />
          </ul>
           <div>
          ßå


           </div>
        </Card.Body>
      </div>
    </Card>
  );
};

export default ProfileCard;