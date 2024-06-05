import React from "react";
import MyModal from "./Modal";
import {Form, Button} from "react-bootstrap";
import useInput from "../hooks/use-input";
import classes from "./AddPost.module.css";
import AuthContext from "../store/authContext";

const AddEvent = ({groupId}) => {
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
    isValid: statusIsValid,
    value: statusValue,
    hassError: statusInputHassError,
    valueChangeHandler: statusChangeHandler,
    valueInputBlurHandler: stausBlurHandler,
    reset: statusTitleInput,
  } = useInput((value) => value.trim() !== "");

  const {
    isValid: enteredDateIsValid,
    value: enteredDate,
    hassError: enteredDateInputHassError,
    valueChangeHandler: enteredDateChangeHandler,
    valueInputBlurHandler: enteredDateBlurHandler,
    reset: resetEnteredDateInput,
  } = useInput((value) => value.trim() !== "");

  const {
    isValid: enteredDescriptionIsValid,
    value: enteredDescription,
    hassError: descriptionInputHassError,
    valueChangeHandler: descriptionChangeHandler,
    valueInputBlurHandler: descriptionBlurHandler,
    reset: resetPostInput,
  } = useInput((value) => value.trim() !== "");

  let formIsValid =
    titleIsValid &&
    enteredDescriptionIsValid &&
    enteredDateIsValid &&
    statusIsValid
      ? true
      : false;

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!formIsValid) return;
    wsMsgToServer(
      JSON.stringify({
        type: "createEvent",
        message: {
          Title: titleValue,
          Description: enteredDescription,
          CreatorId: user.Id,
          Time: new Date(enteredDate),
          GroupId: groupId,
          Going: 0,
          NotGoing: 0,
          Id: statusValue,
        },
      })
    );

    resetTitleInput("");
    resetPostInput("");
    statusTitleInput("");
    resetEnteredDateInput("");
  };

  const titleInputClasses = titleInputHassError ? `${classes.invalid} ` : "";
  const pastInputClasses = descriptionInputHassError
    ? `${classes.invalid} `
    : "";
  const dateInputClasses = enteredDateInputHassError
    ? `${classes.invalid} `
    : "";
  const statusInputClasses = statusInputHassError ? `${classes.invalid} ` : "";

  return (
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
          <p className={classes["error-text"]}>input title must not be emty.</p>
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
          <p className={classes["error-text"]}>input text must not be emty.</p>
        )}
      </Form.Group>

      <Form.Group controlId="formBasicText" className={statusInputClasses}>
        <Form.Label className="mt-2">Are You Goiing ?</Form.Label>
        <Form.Control
          as="textarea"
          type="Are You Goiing ?"
          placeholder="Enter going"
          value={statusValue}
          onChange={statusChangeHandler}
          onBlur={stausBlurHandler}
        />

        {statusInputClasses && (
          <p className={classes["error-text"]}>input text must not be emty.</p>
        )}
      </Form.Group>

      <Form.Control
        type="date"
        placeholder="Enter Date"
        value={enteredDate}
        onChange={enteredDateChangeHandler}
        onBlur={enteredDateBlurHandler}
      />

      {dateInputClasses && (
        <p className={classes["error-text"]}>input text must not be emty.</p>
      )}
      <Button variant="primary my-3 " type="submit" disabled={!formIsValid}>
        Submit
      </Button>
    </Form>
  );
};

export default AddEvent;
