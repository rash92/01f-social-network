import React from "react";
import MyModal from "./Modal";
import {Form, Button} from "react-bootstrap";
import useInput from "../hooks/use-input";
import classes from "./AddPost.module.css";
import AuthContext from "../store/authContext";

const AddGroup = ({show, setShow}) => {
  const {user, wsMsgToServer} = React.useContext(AuthContext);
  const {
    isValid: titleIsValid,
    value: titleValue,
    hassError: titleInputHassError,
    valueChangeHandler: titleChangeHandler,
    valueInputBlurHandler: titleBlurHandler,
    reset: resetTitleInput,
  } = useInput((value) => value.trim() !== "");

  const {
    isValid: enteredDescriptionIsValid,
    value: enteredDescription,
    hassError: descriptionInputHassError,
    valueChangeHandler: descriptionChangeHandler,
    valueInputBlurHandler: descriptionBlurHandler,
    reset: resetPostInput,
  } = useInput((value) => value.trim() !== "");

  let formIsValid = titleIsValid && enteredDescriptionIsValid ? true : false;

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!formIsValid) return;

    wsMsgToServer(
      JSON.stringify({
        type: "createGroup",
        message: {
          Name: titleValue,
          Description: enteredDescription,
          CreatorId: user.Id,
          CreatedAt: new Date(),
          Id: "",
        },
      })
    );

    resetTitleInput("");
    resetPostInput("");
    handleClose();
  };

  const titleInputClasses = titleInputHassError ? `${classes.invalid} ` : "";
  const pastInputClasses = descriptionInputHassError
    ? `${classes.invalid} `
    : "";

  const handleClose = () => setShow(false);

  return (
    <MyModal handleClose={handleClose} show={show} flag={false}>
      <Form onSubmit={handleSubmit} className="py-3">
        <Form.Group controlId="formBasicTitle" className={titleInputClasses}>
          <Form.Label className="mt-2">title</Form.Label>
          <Form.Control
            type="text"
            placeholder="Enter title"
            value={titleValue}
            onBlur={titleBlurHandler}
            onChange={titleChangeHandler}
          />

          {titleInputClasses && (
            <p className={classes["error-text"]}>
              input title must not be emty.
            </p>
          )}
        </Form.Group>

        <Form.Group controlId="formBasicText" className={pastInputClasses}>
          <Form.Label className="mt-2">Text</Form.Label>
          <Form.Control
            as="textarea"
            type="text"
            placeholder="Enter Discription"
            value={enteredDescription}
            onChange={descriptionChangeHandler}
            onBlur={descriptionBlurHandler}
          />

          {pastInputClasses && (
            <p className={classes["error-text"]}>
              input text must not be emty.
            </p>
          )}
        </Form.Group>

        <Button variant="primary my-3 " type="submit" disabled={!formIsValid}>
          Submit
        </Button>
      </Form>
    </MyModal>
  );
};

export default AddGroup;
