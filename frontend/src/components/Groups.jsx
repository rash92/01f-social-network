import React from "react";
import {ListGroup, Card, Button} from "react-bootstrap";
import {Link} from "react-router-dom";

function Groups({groups}) {
  return (
    <div className="container mt-5">
      <ListGroup>
        {groups?.map((group, index) => (
          <ListGroup.Item key={group.Id} className="border-0 p-0 mb-4">
            <Card className="shadow-sm">
              <Card.Body>
                <div className="d-flex justify-content-between align-items-center">
                  <div>
                    <Card.Title className="text-primary">
                      {group.Title}
                    </Card.Title>
                    <Card.Text>
                      <strong>Name:</strong> {group.BasicInfo.Title}
                    </Card.Text>
                    <Card.Text>
                      <strong>Status:</strong> {group.Status}
                    </Card.Text>
                    <Card.Text>
                      <strong>Description:</strong>
                      {group.BasicInfo.Description}
                    </Card.Text>
                  </div>
                  <Link to={`groups/${group.BasicInfo.Id}`}>
                    <Button variant="primary">View Details</Button>
                  </Link>
                </div>
              </Card.Body>
            </Card>
          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  );
}

export default Groups;
