import React from "react";
import {Container, Card, Button} from "react-bootstrap";

function EventsList({event}) {
  console.log(event.id);
  return (
    <Container>
      <h2>Events</h2>

      <Card className="mb-3">
        <Card.Body>
          <Card.Title>{event.title}</Card.Title>
          <Card.Text>{event.description}</Card.Text>
          <Card.Text>
            Day/Time: {new Date(event.dateTime).toLocaleString()}
          </Card.Text>
          <Card.Text>Options: {event.options.join(", ")}</Card.Text>
          <Button variant="primary">RSVP</Button>
        </Card.Body>
      </Card>
    </Container>
  );
}

export default EventsList;