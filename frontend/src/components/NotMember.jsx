import {Button} from "react-bootstrap";
import User from "./User";
const NoMember = ({Title, GroupAction, status}) => {
  return (
    <div
      style={{
        marginTop: "6rem",
        textAlign: "center",
      }}
    >
      {status === "requested" && (
        <p>
          {" "}
          hold your horses <br />
          <strong>
            {Title}: <br />
          </strong>{" "}
          has not accepted your request yet.
        </p>
      )}

      {status === "none" && (
        <p>
          <strong>
            {Title}: <br />
          </strong>{" "}
          you not a member of this group yet.
        </p>
      )}

      {status === "invited" && (
        <p>
          <strong>
            {" "}
            {Title}: <br />{" "}
          </strong>{" "}
          You have been invited to join this group!
        </p>
      )}

      <div style={{display: "flex", gap: "3rem"}}>
        {/* <User Nickname={Nickname} isLoggedIn={IsOnline} Avatar={Avatar} /> */}

        {status !== "requested" && (
          <Button onClick={GroupAction}>
            {status === "invited" ? "Accept" : "Join"}?
          </Button>
        )}
      </div>
    </div>
  );
};

export default NoMember;
