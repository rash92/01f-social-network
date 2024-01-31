import {useState} from "react";
import classes from "./Profile.module.css";
import Action from "./Action";
import MyModal from "./Modal";
const ProfileActions = ({user}) => {
  const [isActive, setIsActive] = useState("Posts");
  const toggleAction = (active) => {
    setIsActive(active);
  };

  return (
    <MyModal handleClose={handleClose} handleShow={handleShow} show={show}>
      <ul className={classes["profile-action"]}>
        <Action
          numberAction={user.followers}
          actionName={"Followers"}
          toggleAction={toggleAction}
          active={isActive}
        />
        <Action
          numberAction={user.fellowing}
          actionName={"Following"}
          toggleAction={toggleAction}
          active={isActive}
        />
      </ul>
    </MyModal>
  );
};

export default ProfileActions;
