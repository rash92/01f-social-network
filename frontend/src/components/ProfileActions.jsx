// import classes from "./Profile.module.css";

import Action from "./Action";
import MyModal from "./Modal";
import User from "./User";
import classes from "./ProfileActions.module.css";
import {Link} from "react-router-dom";
import {Button} from "react-bootstrap";
import InfiniteScroll from "react-infinite-scroll-component";

const ProfileActions = ({
  user,
  handleClose,
  handleShow,
  show,
  flag,
  toggleAction,
  isActive,
  fetchMoreData,
  hasMore,
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
          numberAction={user.numberOfFollowers}
          actionName={"followers"}
          toggleAction={toggleAction}
          active={isActive}
        />
        <Action
          numberAction={user.numberOfFollowing}
          actionName={"following"}
          toggleAction={toggleAction}
          active={isActive}
        />
      </ul>
      <div>
        <ul>
          <InfiniteScroll
            dataLength={user[isActive]?.length}
            next={fetchMoreData}
            hasMore={hasMore[isActive]}
            loader={<h4>Loading...</h4>}
            height={400}
            endMessage={
              <p style={{textAlign: "center"}}>
                <b>Yay! You have seen it all</b>
              </p>
            }
            className={classes.followers}
          >
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
                  {isActive === "followers" ? "Remove" : "Fellowing"}
                </Button>
              </li>
            ))}
          </InfiniteScroll>
        </ul>
      </div>
    </MyModal>
  );
};

export default ProfileActions;
