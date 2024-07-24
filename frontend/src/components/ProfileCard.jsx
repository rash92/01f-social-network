// ProfileCard.js

import {Card, Form, Button} from "react-bootstrap";
import RoundImage from "./RoundedImg";
import classes from "./ProfileCard.module.css";
import Action from "./Action";

const ProfileCard = ({
  data,
  toggleAction,
  toggleProfileVisibility,
  owner,
  isPrivate,
  requestFollow,
  showChat,
}) => {
  const {
    Owner: user,
    PendingFollowers,
    Followers,
    Following,
    Posts,
    IsFollowed,
  } = data;

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
            {user?.Avatar && (
              <RoundImage
                src={user?.Avatar}
                alt="Profile Picture"
                size={"100px"}
              />
            )}

            <ul className={`${classes["profile-action"]}`}>
              <Action
                numberAction={Posts?.length || 0}
                actionName={"Posts"}
                toggleAction={toggleAction}
              />

              <Action
                numberAction={Followers?.length || 0}
                actionName={"Followers"}
                toggleAction={toggleAction}
              />

              <Action
                numberAction={Following?.length || 0}
                actionName={"Following"}
                toggleAction={toggleAction}
              />
              {owner && (
                <Action
                  numberAction={PendingFollowers?.length || 0}
                  actionName={"PendingFollowers"}
                  toggleAction={toggleAction}
                />
              )}
            </ul>
          </div>
          <Card.Title>
            <span>{`${user.FirstName} ${user.LastName}`}</span>
          </Card.Title>
          <Card.Subtitle className="mb-2 text-muted">
            Date of Birth: {user.DOB}
          </Card.Subtitle>
          {user.Nickname && <Card.Text>Nickname: {user.Nickname}</Card.Text>}
          {user.AboutMe && <Card.Text>About Me: {user.AboutMe}</Card.Text>}

          {!owner && (
            <div className={classes.interact}>
              <Button
                variant={`${!IsFollowed ? "primary" : "secondary"}`}
                onClick={requestFollow}
              >
                {IsFollowed ? "Unfollow" : "Follow"}
              </Button>

              <Button
                onClick={showChat.bind(null, {
                  id: user.Id,
                  Nickname: user.Nickname,
                  Avatar: user.Avatar,
                  type: "privateMessage",
                })}
              >
                chat
              </Button>
            </div>
          )}
        </Card.Body>
      </div>
    </Card>
  );
};

export default ProfileCard;
