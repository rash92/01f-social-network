import React, {useState} from "react";
import {InputGroup, FormControl, Button} from "react-bootstrap";
import clsses from "./PostInput.module.css";
import AddPost from "./AddPost";
export default function PostInput({src}) {
  const [show, setShow] = useState(false);
  c

  return (
    <>
      <InputGroup
        className={`mb-3 ${clsses["custom-input-group"]}`}
        style={{marginTop: "2rem", marginBottom: "2rem"}}
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
        />

        <Button variant="primary">Post</Button>
      </InputGroup>
      <AddPost show={show} setShow={setShow} />
    </>
  );
}