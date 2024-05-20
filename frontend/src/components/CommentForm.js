import {Form, Button} from "react-bootstrap";
import useInput from "../hooks/use-input";
import FormGroup from "./FormGroup";
import {getJson} from "../helpers/helpers";
import {useState, useContext} from "react";
import AuthContext from "../store/authContext";
import moment from "moment";
import {AiOutlineLike, AiOutlineDislike} from "react-icons/ai";
import classes from "./CommentForm.module.css";
import {useNavigate} from "react-router-dom";
import {convertImageToBase64} from "../helpers/helpers";
const reactCommentLikeDislike = async ({commentId, postId}, quary) => {
  try {
    return await getJson("react-omment-like-dislike", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        user_token: document.cookie,
      },
      credentials: "include",
      body: JSON.stringify({
        commentId,
        postId,
        quary,
      }),
    });
  } catch (err) {
    throw err;
  }
};

const Comment = ({comments, id: postId}) => {
  const Navigate = useNavigate();
  const [addCommentErr, setAddCommentErr] = useState({
    message: "",
    isError: true,
  });

  const [avatarValue, setAvatarValue] = useState({value: "", file: null});
  const {
    user: {isLogIn: isLoggedIn},
  } = useContext(AuthContext);
  // const timeago = moment(created_at).fromNow();
  const {
    isValid: enteredCommentIsValid,
    value: entereComment,
    hassError: commnetInputHassError,
    valueChangeHandler: commentChangeHandler,
    valueInputBlurHandler: commentBlurHandler,
    reset: resetCommentInput,
  } = useInput((value) => value.trim() !== "");

  const formIsValid = enteredCommentIsValid ? true : false;
  const commnetInputClasses = commnetInputHassError
    ? `${classes.invalid} `
    : "";

  const handleSubmit = async (event) => {
    event.preventDefault();

    if (!formIsValid) return;

    const data = {
      body: entereComment,
      postId: postId,
      image: avatarValue.file
        ? await convertImageToBase64(avatarValue.file)
        : null,
    };
    try {
      const res = await getJson("add-Comment", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          user_token: document.cookie,
        },
        credentials: "include",
        body: JSON.stringify(data),
      });
      if (res.success) {
        // OnAddCommentToPost(res.id, res.comments);
        setAddCommentErr({message: "", isError: false});
        resetCommentInput();
      }
    } catch (error) {
      setAddCommentErr({message: error.message, isError: true});
    }
  };

  const likeHandler = async (option, e) => {
    try {
      if (!isLoggedIn) {
        Navigate("/login");
      } else {
        const res = await reactCommentLikeDislike(option, "like");
        // onAddLikeDislikeComment(option, res, "like");
      }
    } catch (err) {
      console.log(err);
    }
  };

  const disLikeHandler = async (option, e) => {
    try {
      const res = await reactCommentLikeDislike(option, "dislike");
      // onAddLikeDislikeComment(option, res, "dislike");
    } catch (err) {
      console.log(err);
    }
  };

  return (
    <div className="comment mt-3">
      <Form onSubmit={handleSubmit} className="py-3">
        <Form.Group controlId="formBasicText" className={commnetInputClasses}>
          <Form.Control
            as="textarea"
            type="text"
            placeholder="add Comment"
            value={entereComment}
            onChange={commentChangeHandler}
            onBlur={commentBlurHandler}
          />

          {addCommentErr.isError && (
            <p className={classes["error-text"]}>{addCommentErr.message}</p>
          )}
        </Form.Group>

        <FormGroup
          value={avatarValue}
          setValue={setAvatarValue}
          type="file"
          accept="image/*"
          Label="Add image"
          Text="JPG, PNG, GIF."
        />

        <Button variant="primary my-3 " type="submit" disabled={!formIsValid}>
          Submit
        </Button>
      </Form>

      <div className={classes["comment-container"]}>
        {comments?.map(
          ({id, created_at, username, body, likes, dislikes}, i) => (
            <div className={classes.comment} key={id}>
              <div className={classes["comment-header"]}>
                <span className="user name">
                  posted by {username} {moment(created_at).fromNow()}
                </span>
              </div>

              <div className={classes["comment-body"]}>{body}</div>

              <div className={classes["react-status"]}>
                <span>{likes} likes</span>
                <span>{dislikes} dislikes</span>
              </div>

              <div className={classes["comment-reaction"]}>
                <button
                  onClick={likeHandler.bind(null, {commentId: id, postId})}
                >
                  <AiOutlineLike /> like
                </button>
                <button
                  onClick={disLikeHandler.bind(null, {commentId: id, postId})}
                >
                  <AiOutlineDislike /> dislike
                </button>
              </div>
            </div>
          )
        )}

        {/* {comments?.map((el, i)=> <p key={i}>{el.body}</p>)} */}
      </div>
    </div>
  );
};

export default Comment;
