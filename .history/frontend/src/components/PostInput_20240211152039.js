import React from "react";
import {Container, InputGroup, FormControl, Button} from "react-bootstrap";
import clsses from "./post.module.css"
export default function PostInput({src}) {
  return (
    <Container>
      <InputGroup
        className="mb-3"
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
    </Container>
  );
}