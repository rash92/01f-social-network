// import classes from "./Profile.module.css";

import Action from "./Action";
import MyModal from "./Modal";
import User from "./User";
import classes from "./ProfileActions.module.css";
import {Link} from "react-router-dom";
import {Button} from "react-bootstrap";
// import InfiniteScroll from "react-infinite-scroll-component";

const ProfileActions = ({
  user,
  handleClose,
  handleShow,
  show,
  flag,
  toggleAction,
  isActive,
  owner,
  // fetchMoreFellowers,
  // hasMoreFellowers,
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
          numberAction={user.Followers.length}
          actionName={"followers"}
          toggleAction={toggleAction}
          active={isActive}
        />
        <Action
          numberAction={user.Following.length}
          actionName={"following"}
          toggleAction={toggleAction}
          active={isActive}
        />
        {owner && (
          <Action
            numberAction={user?.Requests?.length || 0} //  user.numberOfRequests
            actionName={"Requests"}
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
          {user[isActive]?.map((user, i) => (
            <li key={i} className={classes.item}>
              <Link to={`/profile/${i}`} className={classes.links}>
                <User
                  userName={`${user.username} ${i + 1}`}
                  isLoggedIn={true}
                  name={user.name}
                />
              </Link>

              <Button className={classes.itemButton}>
                {isActive === "followers" ? "Remove" : "Following"}
              </Button>
            </li>
          ))}
          {/* </InfiniteScroll> */}
        </ul>
      </div>
    </MyModal>
  );
};

export default ProfileActions;
