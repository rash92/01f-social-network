import React, {useState, useContext, useEffect} from "react";
import MyModal from "./Modal";
import {Form, Button} from "react-bootstrap";
import useInput from "../hooks/use-input";
import classes from "./AddPost.module.css";
import AuthContext from "../store/authContext";
import {getJson} from "../helpers/helpers";
import User from "./User";
import {Link} from "react-router-dom";

const AddPost = ({show, setShow}) => {
  const [searchTerm, setSearchTerm] = useState("");
  const [suggestions, setSuggestions] = useState([]);
  const [chosenFollowers, setChosenFollowers] = useState([]);
  const [isRighprivacy, setIsRighprivacy] = useState(false);
  const [searchInputTouch, setsearchInputTouch] = useState(false);
  // const [show, setShow] = useState(false);
  // const {OnAddPost} = useContext(AuthContext);
  const [errorAddPost, setErrorAddpost] = useState({
    message: "",
    isError: false,
  });

  const [privacy, setPrivacy] = useState("public");

  const searchBlurHandler = (e) => {
    setsearchInputTouch(true);
  };

  const followerClickHandler = (user, e) => {
    e.preventDefault();
    if (chosenFollowers.some((el) => el.id === user.id)) {
      return;
    }

    setChosenFollowers((prev) => [...prev, user]);
  };

  useEffect(() => {
    // Check if form is valid whenever privacy or suggestions change
    if (privacy === "almost" && suggestions.length < 1) {
      setIsRighprivacy(false);
    } else {
      setIsRighprivacy(true);
    }
  }, [privacy, suggestions]);

  const options = [
    {id: 1, username: "abdi2", name: "abdi"},
    {id: 2, username: "ahmed34", name: "ahmed"},
  ];

  const handleSearch = (event) => {
    const searchTerm = event.target.value;
    setSearchTerm(searchTerm);

    // Filter options based on the search term
    const filteredSuggestions = options.filter((option) =>
      option.username.toLowerCase().includes(searchTerm.toLowerCase())
    );

    setSuggestions(filteredSuggestions);
  };

  const handlePrivacyChange = (event) => {
    setPrivacy(event.target.value);
  };
  // const handleClose = () => setShow(false);
  // const handleShow = () => {
  //   setShow(true);
  // };

  const {
    isValid: titleIsValid,
    value: titleValue,
    hassError: titleInputHassError,
    valueChangeHandler: titleChangeHandler,
    valueInputBlurHandler: titleBlurHandler,
    reset: resetTitleInput,
  } = useInput((value) => value.trim() !== "");

  const {
    isValid: enterePostIsValid,
    value: enterePost,
    hassError: postInputHassError,
    valueChangeHandler: postChangeHandler,
    valueInputBlurHandler: postBlurHandler,
    reset: resetPostInput,
  } = useInput((value) => value.trim() !== "");

  const formIsValid =
    titleIsValid && enterePostIsValid && isRighprivacy ? true : false;

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!formIsValid) return;
    const data = {
      title: titleValue,
      body: enterePost,
    };
    try {
      const res = await getJson("add-post", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          user_token: document.cookie,
        },
        credentials: "include",
        body: JSON.stringify(data),
      });
      if (res.success) {
        resetTitleInput();
        resetPostInput();
        setIsRighprivacy(false);
        setPrivacy("puplic");
        setsearchInputTouch([]);
        chosenFollowers([]);
        // OnAddPost(res.post);
        handleClose();
        setErrorAddpost({message: "", isError: false});
      }
    } catch (error) {
      setErrorAddpost({message: error.message, isError: true});
    }
  };

  const titleInputClasses = titleInputHassError ? `${classes.invalid} ` : "";
  const pastInputClasses = postInputHassError ? `${classes.invalid} ` : "";
  const searchInputClasses =
    chosenFollowers.length < 0 && searchInputTouch ? `${classes.invalid} ` : "";

  const handleClose = () => setShow(false);
  // const handleShow = () => {
  //   setShow(true);
  // };
  return (
    <MyModal handleClose={handleClose} show={show} flag={false}>
      <Form onSubmit={handleSubmit} className="py-3">
        <Form.Group controlId="formBasicTitle" className={titleInputClasses}>
          <Form.Label className="mt-2">title</Form.Label>
          <Form.Control
            type="text"
            placeholder="Enter Email"
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
            placeholder="please add post"
            value={enterePost}
            onChange={postChangeHandler}
            onBlur={postBlurHandler}
          />

          {pastInputClasses && (
            <p className={classes["error-text"]}>
              input text must not be emty.
            </p>
          )}
        </Form.Group>

        <Form.Group controlId="privacy">
          <Form.Label>Choose Privacy:</Form.Label>
          <Form.Select value={privacy} onChange={handlePrivacyChange}>
            <option selected value="public">
              Public
            </option>
            <option value="private">Private</option>
            <option value="almost">Almost private</option>
          </Form.Select>
        </Form.Group>
        {privacy === "almost" && (
          <>
            <Form.Group
              controlId="formBasicText"
              className={searchInputClasses}
            >
              <Form.Label className="mt-2">Search follower</Form.Label>
              <Form.Control
                as="input"
                type="text"
                placeholder="please search follower to add"
                value={searchTerm}
                onChange={handleSearch}
                onBlur={searchBlurHandler}
              />

              {searchInputClasses && (
                <p className={classes["error-text"]}>
                  input text must not be search.
                </p>
              )}
            </Form.Group>

            <ul className={classes.sectionsList}>
              {suggestions.map((suggestion, index) => (
                <li
                  key={index}
                  onClick={followerClickHandler.bind(null, suggestion)}
                >
                  <User name={suggestion.username} isLoggedIn={true} />
                </li>
              ))}
            </ul>

            {
              <div classes={classes.fellowers}>
                <div>Chosen followers</div>
                {chosenFollowers.map((chosen, index) => (
                  <Link
                    to={`/profile/${chosen.id}`}
                    className={classes.chosen}
                    key={index}
                  >
                    <User
                      name={chosen.name}
                      userName={chosen.username}
                      isLoggedIn={true}
                    />
                  </Link>
                ))}
              </div>
            }
          </>
        )}

        <Button variant="primary my-3 " type="submit" disabled={!formIsValid}>
          Submit
        </Button>
      </Form>

      {errorAddPost.isError && (
        <p className={classes["error-text"]}>{errorAddPost.message}.</p>
      )}
    </MyModal>
  );
};

export default AddPost;
