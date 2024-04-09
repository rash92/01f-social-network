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
          numberAction={Followers.length}
          actionName={"Followers"}
          toggleAction={toggleAction}
          active={isActive}
        />
        <Action
          numberAction={Following.length}
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
          {/* <InfiniteScroll
            dataLength={
              user[
                `NumberOf${isActive.charAt(0).toUpperCase()}${isActive.slice(
                  1
                )}`
              ] || 0
            }
            next={fetchMoreFellowers}
            hasMore={hasMoreFellowers[isActive]}
            loader={<h4>Loading...</h4>}
            height={400}
            endMessage={
              <p style={{textAlign: "center"}}>
                <b>Yay! You have seen it all</b>
              </p>
            }
            className={classes.followers}
          > */}
          {active?.map((user, i) => (
            <li key={i} className={classes.item}>
              <Link to={`/profile/${user.UserId}`} className={classes.links}>
                <User
                  Avatar={user.Avatar}
                  Nickname={`${user.Nickname}`}
                  isLoggedIn={true}
                />
              </Link>
              {
                isActive === "PendingFollowers" && (
                  <>
                    <Button
                      className={classes.itemButton}
                      onClick={accepOrRejectRequestHandler.bind(null, {
                        id: user.Id,
                        flag: "yes",
                      })}
                    >
                      Accept
                    </Button>

                    <Button
                      onClick={accepOrRejectRequestHandler.bind(null, {
                        id: user.Id,
                        flag: "no",
                      })}
                      className={classes.itemButton}
                    >
                      Reject
                    </Button>
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
