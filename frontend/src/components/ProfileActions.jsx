// import classes from "./Profile.module.css";
import Action from "./Action";
import MyModal from "./Modal";
import User from "./User";
const ProfileActions = ({
  user,
  handleClose,
  handleShow,
  show,
  flag,
  toggleAction,
  isActive,
}) => {
  return (
    <MyModal
      handleClose={handleClose}
      handleShow={handleShow}
      show={show}
      flag={flag}
    >
      <ul>
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
        <ul>
          {user[isActive]?.map((user, i) => (
            <li key={i}>
              <User
                userName={user.userName}
                isLoggedIn={true}
                name={user.name}
              />
            </li>
          ))}
        </ul>
      </div>
    </MyModal>
  );
};

export default ProfileActions;
