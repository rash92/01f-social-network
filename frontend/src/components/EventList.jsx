import React, {useContext} from "react";
import {Container, Card, Button} from "react-bootstrap";
import {TiTick} from "react-icons/ti";
import {FaTimesCircle} from "react-icons/fa";
import AuthContext from "../store/authContext";

function EventsList({event}) {
  const {user, isWsReady, wsMsgToServer} = useContext(AuthContext);
  const handleRSVP = (status) => {
    console.log(isWsReady);
    if (isWsReady) {
      console.log(
        {
          type: "toggleAttendEvent",
          message: {
            SenderId: user.Id,
            EventId: event?.event.Id,
            GroupId: event.event.GroupId,
          },
        },
        Date.now(),
        "sending"
      );
      wsMsgToServer(
        JSON.stringify({
          type: "toggleAttendEvent",
          message: {
            SenderId: user.Id,
            EventId: event?.event.Id,
            GroupId: event.event.GroupId,
          },
        })
      );
    }
  };

  return (
    <Container>
      <Card className="mb-3">
        <Card.Body>
          <Card.Title>{event?.event?.Title}</Card.Title>
          <Card.Text>{event?.event?.Description}</Card.Text>
          <Card.Text>
            Day/Time: {new Date(event.event.Time).toLocaleString()}
          </Card.Text>

          <Card.Text>Going: {event?.event?.Going}</Card.Text>
          <Card.Text>Not going: {event?.event?.NotGoing}</Card.Text>

          <div className="d-flex justify-content-between align-items-center">
            <Button
              disabled={event.Going}
              variant="primary"
              onClick={() => handleRSVP("going")}
            >
              {event.Going && <TiTick />}
              Going
            </Button>
            <Button
              disabled={!event.Going}
              variant="secondary"
              onClick={() => handleRSVP("not going")}
            >
              {!event?.Going && <FaTimesCircle size={20} />}
              <span
                style={{
                  marginLeft: "0.5rem",
                }}
              >
                Not Going
              </span>
            </Button>
          </div>
        </Card.Body>
      </Card>
    </Container>
  );
}

export default EventsList;
