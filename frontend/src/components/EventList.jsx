import React from "react";
import {Container, Card, Button} from "react-bootstrap";

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
          <Card.Title>{event.title}</Card.Title>
          <Card.Text>{event.description}</Card.Text>
          <Card.Text>
            Day/Time: {new Date(event.dateTime).toLocaleString()}
          </Card.Text>
          <Card.Text>Options: {event.options.join(", ")}</Card.Text>
          <div className="d-flex justify-content-between align-items-center">
            <Button variant="primary" onClick={() => handleRSVP("going")}>
              RSVP Going
            </Button>
            <Button variant="secondary" onClick={() => handleRSVP("not going")}>
              RSVP Not Going
            </Button>
          </div>
        </Card.Body>
      </Card>
    </Container>
  );
}

export default EventsList;
