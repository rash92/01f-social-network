import {Button} from "react-bootstrap";
import User from "./User";
const PrivateProfile = ({name, userName, fellowUserHandler}) => {
  return (
    <div
      style={{
        marginTop: "6rem",
        textAlign: "center",
      }}
    >
      <p>{name}'s profile is Private . Please Fellow.</p>
      <div style={{display: "flex", gap: "3rem"}}>
        <User userName={userName} isLoggedIn={true} name={name} />
        <Button onClick={fellowUserHandler}>Follow</Button>
      </div>
    </div>
  );
};

export default PrivateProfile;
