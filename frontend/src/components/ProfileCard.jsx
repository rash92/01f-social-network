// ProfileCard.js

import {Card, Form, Button} from "react-bootstrap";
import RoundImage from "./RoundedImg";
import classes from "./ProfileCard.module.css";
import Action from "./Action";
import {Link} from "react-router-dom";

const ProfileCard = ({
  user,
  toggleAction,
  toggleProfileVisibility,
  owner,
  isPrivate,
  isFollowed,
}) => {
  // console.log(owner, "owner");
  return (
    <Card className={classes.card}>
      <div className={classes.profileContainer}>
        {owner && (
          <Form>
            <Form.Check
              type="switch"
              id="profileVisibilitySwitch"
              label={`Make Profile ${isPrivate ? "Public" : "Private"}`}
              checked={isPrivate}
              onChange={toggleProfileVisibility}
            />
          </Form>
        )}
        <Card.Body>
          <div className={classes.image}>
            {user.Profile && (
              <RoundImage
                src={user.Profile}
                alt="Profile Picture"
                size={"100px"}
              />
            )}
                     
            <ul className={`${classes["profile-action"]}`}>
              <Action
                numberAction={user.Posts?.length || 0}
                actionName={"Posts"}
                toggleAction={toggleAction}
              />

              <Action
                numberAction={user.Followers?.length || 0}
                actionName={"Followers"}
                toggleAction={toggleAction}
              />

              <Action
                numberAction={user.Following?.length || 0}
                actionName={"Following"}
                toggleAction={toggleAction}
              />
              {owner && (
                <Action
                  numberAction={user?.Requests?.length || 0}
                  actionName={"Requests"}
                  toggleAction={toggleAction}
                />
              )}
            </ul>
          </div>
          <Card.Title>
            <span>{`${user.FirstName} ${user.LastName}`}</span>
          </Card.Title>
          <Card.Subtitle className="mb-2 text-muted">
            Date of Birth: {user.age}
          </Card.Subtitle>
          {user.Nickname && <Card.Text>Nickname: {user.Nickname}</Card.Text>}
          {user.AboutMe && <Card.Text>About Me: {user.AboutMe}</Card.Text>}

          {!owner && (
            <div className={classes.interact}>
              <Button>{isFollowed ? "Ufellow" : "Fellow"}</Button>

              <Link to={`/chats/${user.Id}`}>
                <Button> message </Button>
              </Link>
            </div>
          )}
        </Card.Body>
      </div>
    </Card>
  );
};

export default ProfileCard;
