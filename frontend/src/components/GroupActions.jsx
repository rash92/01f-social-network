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
          {isActive === "create event" && <AddEvent />}

          {isActive === "RequestedMembers" &&
            group[isActive]?.map((user, i) => (
              <li key={i} className={classes.item}>
                <Link to={`/profile/${i}`} className={classes.links}>
                  <User
                    userName={`${user.nickname} ${i + 1}`}
                    isLoggedIn={true}
                    name={user.name}
                  />
                </Link>
                <Button>approve</Button>
              </li>
            ))}

          {isActive === "Members" &&
            group[isActive]?.map((user, i) => (
              <li key={i} className={classes.item}>
                <Link to={`/profile/${i}`} className={classes.links}>
                  <User
                    userName={`${user.nickname} ${i + 1}`}
                    isLoggedIn={true}
                    name={user.name}
                    avatar={user.avatar}
                    isOnline={user.isOnline}
                  />
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
