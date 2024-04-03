import {Button} from "react-bootstrap";
import User from "./User";
const PrivateProfile = ({
  Nickname,
  fellowUserHandler,
  IsOnline,
  Avatar,
  IsPending,
}) => {
  return (
    <div
      style={{
        marginTop: "6rem",
        textAlign: "center",
      }}
    >
      {IsPending ? (
        <p>{Nickname} has not accepted your request yet. Please wait.</p>
      ) : (
        <p>{Nickname}'s profile is Private. Please Follow.</p>
      )}
      <div style={{display: "flex", gap: "3rem"}}>
        <User Nickname={Nickname} isLoggedIn={IsOnline} Avatar={Avatar} />

        {!IsPending && <Button onClick={fellowUserHandler}>Follow</Button>}
      </div>
    </div>
  );
};

export default PrivateProfile;
