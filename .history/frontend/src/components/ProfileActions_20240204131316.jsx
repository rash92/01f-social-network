// import classes from "./Profile.module.css";

import {useState} from "react";
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
}) => {
  const [activeFellowers, setActiveFellowers] = useState(user[isActive]);
  const [hasMore, setHasMore] = useState(true);

  const fetchMoreData = () => {
    if (activeFellowers.length <= 50) {
      setActiveFellowers((pre) => [
        ...pre,
        ...user[isActive].slice(pre.length),
      ]);
    } else {
      setHasMore(false);
    }
  };
  console.log(activeFellowers, isActive);
  return (
    <MyModal
      handleClose={handleClose}
      handleShow={handleShow}
      show={show}
      flag={flag}
    >
      <ul className={classes.actions}>
        <Action
          numberAction={user.followers.length}
          actionName={"followers"}
          toggleAction={toggleAction}
          active={isActive}
        />
        <Action
          numberAction={user.following.length}
          actionName={"following"}
          toggleAction={toggleAction}
          active={isActive}
        />
      </ul>
      <div>
        <InfiniteScroll
          dataLength={activeFellowers?.length }
          next={fetchMoreData}
          hasMore={hasMore}
          loader={<h4>Loading...</h4>}
          endMessage={
            <p style={{textAlign: "center"}}>
              <b>Yay! You have seen it all</b>
            </p>
          }
        >
          <ul className={classes.followers}>
            {activeFellowers?.map((user, i) => (
              <li key={i} className={classes.item}>
                <Link to={`/profile/${i}`} className={classes.links}>
                  <User
                    userName={user.userName}
                    isLoggedIn={true}
                    name={user.name}
                  />
                </Link>

                <Button className={classes.itemButton}>
                  {isActive === "followers" ? "Remove" : "Fellowing"}
                </Button>
              </li>
            ))}
          </ul>
        </InfiniteScroll>
      </div>
    </MyModal>
  );
};

export default ProfileActions;