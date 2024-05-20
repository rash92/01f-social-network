import React, {useState} from "react";
import {ListGroup} from "react-bootstrap";
import {Link} from "react-router-dom";
import AddGroup from "./AddGroup";
import {Button} from "react-bootstrap";

function Groups({groups}) {
  const [show, setShow] = useState(false);

  return (
    <div>
      <Button
        onClick={() => {
          setShow(true);
        }}
      >
        add Group
      </Button>
      <AddGroup show={show} setShow={setShow} />
      <ListGroup style={{marginTop: "4rem"}}>
        {groups?.map((group, index) => (
          <ListGroup.Item
            key={group.Id}
            style={{marginTop: index !== 0 ? "2rem" : " 0"}}
          >
            <Link
              to={`groups/${group.Id}`}
              style={{textDecoration: "none", color: "black"}}
            >
              <div
                className="d-flex justify-content-between align-items-center"
                style={{
                  gap: "8rem",
                }}
              >
                <div>
                  <h5>{group.Title}</h5>
                  <p>{group.Description}</p>
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
