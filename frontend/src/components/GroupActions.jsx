import Action from "./Action";
import MyModal from "./Modal";
import User from "./User";
import classes from "./ProfileActions.module.css";
import {Link} from "react-router-dom";
import EventList from "./EventList";
import {Button} from "react-bootstrap";
import AddEvent from ".././components/AddEvent";

const GroupActions = ({
  group,
  handleClose,
  handleShow,
  show,
  flag,
  toggleAction,
  isActive,
  owner,
  status,
  AcceptRefuseGroupRequestHandler,
  inviteHandler,
}) => {
  return (
    <MyModal
      handleClose={handleClose}
      handleShow={handleShow}
      show={show}
      flag={flag}
    >
      <ul className={classes.actions}>
        {status ? (
          <>
            <Action
              numberAction={group.NumberOfMembers}
              actionName={"Members"}
              toggleAction={toggleAction}
              active={isActive}
            />
            <Action
              numberAction={group.NumberOfEvents}
              actionName={"Events"}
              toggleAction={toggleAction}
              active={isActive}
            />
          </>
        ) : (
          <>
            <Action
              numberAction={""}
              actionName={"create event"}
              toggleAction={toggleAction}
              active={isActive}
            />
            <Action
              numberAction={""}
              actionName={"Invite"}
              toggleAction={toggleAction}
              active={isActive}
            />
            {owner && (
              <Action
                numberAction={group.NumberOfRequests}
                actionName={"RequestedMembers"}
                toggleAction={toggleAction}
                active={isActive}
              />
            )}
          </>
        )}
      </ul>
      <div>
        <ul>
          {isActive === "create event" && (
            <AddEvent groupId={group?.BasicInfo?.Id} />
          )}

          {isActive === "RequestedMembers" &&
            group[isActive]?.map((user, i) => (
              <li key={user?.Id} className={classes.item}>
                <Link to={`/profile/${user?.Id}`} className={classes.links}>
                  <User Nickname={`${user?.Nickname}`} Avatar={user.Avatar} />
                </Link>
                <Button
                  onClick={AcceptRefuseGroupRequestHandler.bind(null, {
                    id: user.Id,
                    type: true,
                  })}
                >
                  Accept
                </Button>
                <Button
                  onClick={AcceptRefuseGroupRequestHandler.bind(null, {
                    id: user.Id,
                    type: false,
                  })}
                  variant="secondary"
                >
                  Reject
                </Button>
              </li>
            ))}

          {isActive === "Invite" &&
            group[isActive]?.map((user) => (
              <li key={user.Id} className={classes.item}>
                <Link to={`/profile/${user.Id}`} className={classes.links}>
                  <User Nickname={`${user.Nickname}`} Avatar={user.Avatar} />
                </Link>
                <Button
                  variant={!user.isInvited ? "primary" : "secondary"}
                  disabled={!user?.isInvited ? false : true}
                  onClick={inviteHandler?.bind(null, user.Id)}
                >
                  {user?.isInvited ? "Invited" : "Invite"}
                </Button>
              </li>
            ))}

          {isActive === "Members" &&
            group[isActive]?.map((user) => (
              <li key={user?.Id} className={classes.item}>
                <Link to={`/profile/${User?.Id}`} className={classes.links}>
                  <User Nickname={`${user?.Nickname}`} Avatar={user?.Avatar} />
                </Link>
              </li>
            ))}

          {isActive === "Events" &&
            group[isActive]?.map((event, i) => (
              <EventList key={event.id} event={event} />
            ))}
        </ul>
      </div>
    </MyModal>
  );
};

export default GroupActions;
