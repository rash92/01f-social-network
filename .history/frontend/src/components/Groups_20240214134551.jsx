import React from "react";
import {ListGroup, Badge} from "react-bootstrap";
import {Link} from 'react-router-dom'

function Groups({groups}) {
  return (
    <div>
      <ListGroup style={{marginTop: "4rem"}}>
        {groups.map((group, index) => (
          <ListGroup.Item
            key={group.id}
            style={{marginTop: index !== 0 ? "2rem" : " 0"}}
          >
            <Link to={`groups/`}>
              <div
                className="d-flex justify-content-between align-items-center"
                style={{
                  gap: "8rem",
                }}
              >
                <div>
                  <h5>{group.title}</h5>
                  <p>{group.description}</p>
                </div>
                <div>
                  <span>{group.status}</span>
                </div>
              </div>
            </Link>
          </ListGroup.Item>
        ))}
      </ListGroup>
    </div>
  );
}

export default Groups;