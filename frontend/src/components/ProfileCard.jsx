// ProfileCard.js

import {Card, Form, Button} from "react-bootstrap";
import RoundImage from "./RoundedImg";
import classes from "./ProfileCard.module.css";
import Action from "./Action";
import {Link} from "react-router-dom";

const ProfileCard = ({user, toggleAction, toggleProfileVisibility}) => {
  return (
    <Card className={classes.card}>
      <div className={classes.profileContainer}>
        {user.isOwner && (
          <Form>
            <Form.Check
              type="switch"
              id="profileVisibilitySwitch"
              label={`Make Profile ${user.IsPrivate ? "Public" : "Private"}`}
              checked={user.IsPrivate}
              onChange={toggleProfileVisibility}
            />
          </Form>
        )}
        <Card.Body>
          <div className={classes.image}>
            {user.avatar && (
              <RoundImage
                src={user.avatar}
                alt="Profile Picture"
                size={"100px"}
              />
            )}
            <ul className={`${classes["profile-action"]}   `}>
              <Action
                numberAction={user.numberOfPosts}
                actionName={"Posts"}
                toggleAction={toggleAction}
              />
              <Action
                numberAction={user.numberOfFollowers}
                actionName={"Followers"}
                toggleAction={toggleAction}
              />
              <Action
                numberAction={user.numberOfFollowing}
                actionName={"Following"}
                toggleAction={toggleAction}
              />
            </ul>
          </div>
          <Card.Title>
            <span>{`${user.firstName} ${user.lastName}`}</span>
          </Card.Title>
          <Card.Subtitle className="mb-2 text-muted">
            Date of Birth: {user.dateOfBirth}
          </Card.Subtitle>
          {user.nickname && <Card.Text>Nickname: {user.nickname}</Card.Text>}
          {user.aboutMe && <Card.Text>About Me: {user.aboutMe}</Card.Text>}
          <div className={classes.interact}>
            {/* {!user.isOwner && ( */}
            <>
              <Button>{user.relStatus}</Button>

              <Link to={`/chats/${user.id}`}>
                <Button> message </Button>
              </Link>
            </>
            {/* )} */}
          </div>
        </Card.Body>
      </div>
    </Card>
  );
};

export default ProfileCard;
