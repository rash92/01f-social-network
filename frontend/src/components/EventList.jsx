import React from "react";
import {Container, Card, Button} from "react-bootstrap";
import {TiTick} from "react-icons/ti";
import {FaTimesCircle} from "react-icons/fa";

function EventsList({event}) {
  const handleRSVP = (status) => {
    // Handle RSVP logic here, you can send the status to a parent component
    console.log(
      `User is ${status === "going" ? "going" : "not going"} to the event`
    );
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

          <Card.Text>Going: 0</Card.Text>
          <Card.Text>Not going: 0</Card.Text>

          <div className="d-flex justify-content-between align-items-center">
            <Button variant="primary" onClick={() => handleRSVP("going")}>
              {event.Going && <TiTick />}
              Going
            </Button>
            <Button variant="secondary" onClick={() => handleRSVP("not going")}>
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
