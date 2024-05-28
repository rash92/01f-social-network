import React, {useState, useCallback, useContext} from "react";
import MyModal from "./Modal";
import {Form, Button} from "react-bootstrap";
import useInput from "../hooks/use-input";
import classes from "./AddPost.module.css";
import {getJson} from "../helpers/helpers";
import FormGroup from "./FormGroup";
import SearchUser from "./SearchUser";
import AuthContext from "../store/authContext";
import {Link} from "react-router-dom";
import User from "./User";
import {convertImageToBase64} from "../helpers/helpers";

const AddPost = ({show, setShow, type = "profile", groupId = ""}) => {
  const {
    user,
    isWsReady,

    wsMsgToServer,
  } = useContext(AuthContext);

  const [errorAddPost, setErrorAddpost] = useState({
    message: "",
    isError: false,
  });
  const [postImgValue, setPostImgValue] = useState({value: "", file: null});
  const [privacy, setPrivacy] = useState("public");
  const handlePrivacyChange = (e) => {
    setPrivacy(e.target.value);
  };
  const [seachList, setSeachList] = useState([]);
  const [typingTimeout, setTypingTimeout] = useState(null);
  const [chosenFollowers, setChosenFollowers] = useState([]);
  const [searchInputTouch, setSearchInputTouch] = useState(false);
  const fetchSearch = async ({query, id}) => {
    try {
      const data = await getJson("search-Follower", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({search: query, id}),
      });
      setSeachList(data);
    } catch (err) {
      console.log(err);
    }
  };
  const searchInputTouchHandler = (input) => {
    setSearchInputTouch(input);
  };

  const handleSearch = useCallback(
    (query) => {
      if (typingTimeout) {
        clearTimeout(typingTimeout);
      }

      setTypingTimeout(
        setTimeout(fetchSearch.bind(null, {query, id: user.Id}), 500)
      );
    },
    [typingTimeout, user.Id]
  );

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

  let formIsValid = titleIsValid && enterePostIsValid ? true : false;

  formIsValid =
    privacy === "superprivate" && chosenFollowers.length === 0
      ? false
      : formIsValid;

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!formIsValid) return;
    const allowedImageTypes = {
      "image/png": true,
      "image/jpeg": true,
      "image/giff": true,
    };

    if (
      postImgValue?.file?.type &&
      !allowedImageTypes[postImgValue?.file?.type]
    ) {
      console.log(postImgValue?.file?.type);
      console.log("this here");
      setErrorAddpost({
        message: "image type  allowed png , jpeg giff ",
        isError: false,
      });
      return;
    }

    const data = {
      type: "post",
      message: {
        title: titleValue,
        body: enterePost,
        PostChosenFollowers: chosenFollowers.map((el) => el.Id),
        privacyLevel: privacy,
        groupId: "",
        creatorId: user.Id,
        createdAt: new Date(),
        image: postImgValue.file
          ? await convertImageToBase64(postImgValue.file)
          : "",
        isWholeForum: type === "profile" ? true : false,
        id: "",
      },
    };

    if (isWsReady) {
      console.log(data, " post that we sending");
      wsMsgToServer(JSON.stringify(data));
    }

    resetTitleInput();
    resetPostInput();
    setChosenFollowers([]);
    setPrivacy("puplic");
    handleClose();
    setErrorAddpost({message: "", isError: false});
  };

  const titleInputClasses = titleInputHassError ? `${classes.invalid} ` : "";
  const pastInputClasses = postInputHassError ? `${classes.invalid} ` : "";
  const searchInputClasses =
    !chosenFollowers.length && searchInputTouch ? `${classes.invalid} ` : "";

  const handleClose = () => setShow(false);
  const addChosen = (id) => {
    if (chosenFollowers.some((el) => el.Id === id)) return;
    setChosenFollowers([
      ...chosenFollowers,
      seachList.find((el) => el.Id === id),
    ]);
  };
  return (
    <MyModal handleClose={handleClose} show={show} flag={false}>
      <Form onSubmit={handleSubmit} className="py-3">
        <Form.Group controlId="formBasicTitle" className={titleInputClasses}>
          <Form.Label className="mt-2">title</Form.Label>
          <Form.Control
            type="text"
            placeholder="Enter Title"
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
            placeholder="please add Text"
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

        {type === "profile" && (
          <Form.Group controlId="privacy">
            <Form.Label>Choose Privacy:</Form.Label>
            <Form.Select value={privacy} onChange={handlePrivacyChange}>
              <option selected value="public">
                Public
              </option>
              <option value="private">Private</option>
              <option value="superprivate">Almost private</option>
            </Form.Select>
          </Form.Group>
        )}

        {privacy === "superprivate" && type !== "group" && (
          <>
            {searchInputClasses && (
              <p className={classes["error-text"]}>
                please choose a Followers.
              </p>
            )}
            <SearchUser
              addChosen={addChosen}
              onSearch={handleSearch}
              searchList={seachList}
              className={classes.NavLink}
              Blur={searchInputTouchHandler}
              type={"follower"}
            />

            <div
              style={{
                margin: "4rem 0",
              }}
            >
              <div>Chosen followers</div>
              {chosenFollowers?.map((chosen, index) => (
                <Link
                  to={`/profile/${chosen.Id}`}
                  className={classes.chosen}
                  key={index}
                >
                  <User
                    Avatar={chosen.Avatar}
                    Nickname={chosen.Nickname}
                    isLoggedIn={true}
                  />
                </Link>
              ))}
            </div>
          </>
        )}

        <div style={{margin: "1rem 0rem"}}>
          <FormGroup
            value={postImgValue}
            setValue={setPostImgValue}
            type="file"
            accept="image/*"
            Label="add photo"
            Text="JPG, PNG, GIF."
          />
        </div>

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
