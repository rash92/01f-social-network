import React, {useState} from "react";
import {InputGroup, FormControl, Button} from "react-bootstrap";
import clsses from "./PostInput.module.css";
import AddPost from "./AddPost";
export default function PostInput({src}) {
  const [show, setShow] = useState(false);

  return (
    <>
      <InputGroup
        className={`mb-3 `}
        style={{marginTop: "2rem", marginBottom: "2rem", borderColor: "none"}}
      >
        <InputGroup.Text>
          <img
            src={`http://localhost:8000/images/${src}`}
            alt="Profile"
            style={{width: "30px", height: "30px", borderRadius: "50%"}}
          />
        </InputGroup.Text>

        <FormControl
          placeholder="What's on your mind?"
          aria-label="Post input"
          style={{
            outline: "none",
            boxShadow: "none",
            borderC: "none",
          }}
          // className={`${clsses["custom-input-group"]}`}
          readOnly
          // onClick={() => setShow(true)}
        />

        <Button variant="primary">Post</Button>
      </InputGroup>
      <AddPost show={show} setShow={setShow} />
    </>
  );
}