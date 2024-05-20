// import classes from "./Profile.module.css";

import Action from "./Action";
import MyModal from "./Modal";
import User from "./User";
import classes from "./ProfileActions.module.css";
import {Link} from "react-router-dom";
import {Button} from "react-bootstrap";
// import InfiniteScroll from "react-infinite-scroll-component";

const ProfileActions = ({
  data: {Owner: user, PendingFollowers, Followers, Following},
  active,
  handleClose,
  handleShow,
  show,
  flag,
  toggleAction,
  isActive,
  owner,
  accepOrRejectRequestHandler,
}) => {
  return (
    <MyModal
      handleClose={handleClose}
      handleShow={handleShow}
      show={show}
      flag={flag}
    >
      <ul className={classes.actions}>
        <Action
          numberAction={Followers?.length}
          actionName={"Followers"}
          toggleAction={toggleAction}
          active={isActive}
        />
        <Action
          numberAction={Following?.length}
          actionName={"Following"}
          toggleAction={toggleAction}
          active={isActive}
        />
        {owner && (
          <Action
            numberAction={PendingFollowers?.length || 0} //  user.numberOfRequests
            actionName={"PendingFollowers"}
            toggleAction={toggleAction}
            active={isActive}
          />
        )}
      </ul>
      <div>
        <ul>
          {active?.map((user, i) => (
            <li key={i} className={classes.item}>
              <Link to={`/profile/${user.Id}`} className={classes.links}>
                <User
                  Avatar={user.Avatar}
                  Nickname={`${user.Nickname}`}
                  isLoggedIn={true}
                />
              </Link>
              {
                isActive === "PendingFollowers" && (
                  <>
                    <div>
                      <Button
                        className={classes.itemButton}
                        onClick={accepOrRejectRequestHandler.bind(null, {
                          id: user.Id,
                          flag: "yes",
                        })}
                      >
                        Accept
                      </Button>
                    </div>

                    <div>
                      <Button
                        onClick={accepOrRejectRequestHandler.bind(null, {
                          id: user.Id,
                          flag: "no",
                        })}
                        className={classes.itemButton}
                      >
                        Reject
                      </Button>
                    </div>
                  </>
                )

                //   <Button className={classes.itemButton}>
                //        {isActive === "followers" ? "Remove" : "Following"}
                // </Button>
              }
            </li>
          ))}
          {/* </InfiniteScroll> */}
        </ul>
      </div>
    </MyModal>
  );
};

export default ProfileActions;
