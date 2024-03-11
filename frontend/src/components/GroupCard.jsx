// ProfileCard.js

import {Card, Button} from "react-bootstrap";
import classes from "./ProfileCard.module.css";
import Action from "./Action";

const GroupCard = ({group, toggleAction, owner = true}) => {
  return (
    <Card className={classes.card}>
      <div className={classes.profileContainer}>
        <Card.Body>
          <div className={classes.image}>
            <Card.Title>
              <span>{group.title}</span>
            </Card.Title>
            <ul className={`${classes["profile-action"]}   `}>
              <Action
                numberAction={group.NumberOfPosts}
                actionName={"Posts"}
                toggleAction={toggleAction}
              />

              <Action
                numberAction={group.NumberOfMembers}
                actionName={"Members"}
                toggleAction={toggleAction}
              />

              <Action
                numberAction={group.NumberOfEvents}
                actionName={"Events"}
                toggleAction={toggleAction}
              />
            </ul>
          </div>

          <Card.Text>{group.description}</Card.Text>
          <div className={classes.interact}>
            <Button onClick={toggleAction.bind(null, "ToBeInvites")}>
              Invite
            </Button>
            {owner && (
              <Button onClick={toggleAction.bind(null, "Requests")}>
                manage
              </Button>
            )}
          </div>
        </Card.Body>
      </div>
    </Card>
  );
};

export default GroupCard;
