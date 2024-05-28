import {Form, Button, Image} from "react-bootstrap";
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

const Comment = ({comments, id: postId, addComment}) => {
  const [addCommentErr, setAddCommentErr] = useState({
    message: "",
    isError: true,
  });

  console.log(comments);

  const [avatarValue, setAvatarValue] = useState({value: "", file: null});
  const {user} = useContext(AuthContext);
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
  const allowedImageTypes = {
    "image/png": true,
    "image/jpeg": true,
    "image/giff": true,
  };
  const handleSubmit = async (event) => {
    event.preventDefault();
    if (!formIsValid) return;
    if (
      avatarValue?.file?.type &&
      !allowedImageTypes[avatarValue?.file?.type]
    ) {
      console.log(avatarValue?.file?.type);

      return;
    }

    // Id              string    `json:"Id"`
    // Body            string    `json:"Body"`
    // CreatorId       string    `json:"CreatorId"`
    // PostID          string    `json:"PostId"`
    // CreatedAt       time.Time `json:"CreatedAt"`
    // Likes           int       `json:"Likes"`
    // Dislikes        int       `json:"Dislikes"`
    // NickName        string    `json:"Nickname"`
    // Image           string    `json:"Image"`
    // CreatorNickname string    `json:"CreatorNickname"

    const data = {
      Id: "",
      Body: entereComment,
      CreatorId: user.Id,
      PostId: postId,
      Image: avatarValue.file
        ? await convertImageToBase64(avatarValue.file)
        : "",
      CreatedAt: new Date(),
      Likes: 0,
      Dislikes: 0,
      CreatorNickname: user.Nickname,
    };
    try {
      const res = await getJson("add-comment", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(data),
      });
      if (res.success) {
        addComment(res.comment);
        setAddCommentErr({message: "", isError: false});
        resetCommentInput();
        setAvatarValue({value: "", file: null});
      }
    } catch (error) {
      setAddCommentErr({message: error.message, isError: true});
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
          value={avatarValue.value}
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
          (
            {Id, CreatedAt, CreatorNickname, Body, Likes, Lislikes, Image: img},
            i
          ) => (
            <div className={classes.comment} key={Id}>
              <div className={classes["comment-header"]}>
                <span className="user name">
                  posted by {CreatorNickname} {moment(CreatedAt).fromNow()}
                </span>
              </div>

              <div className={classes["comment-body"]}>{Body}</div>
              <div>
                {img && (
                  <Image
                    src={`http://localhost:8000/images/${img}`}
                    width={"100%"}
                    height={"100%"}
                  />
                )}
              </div>

              {/* <div className={classes["react-status"]}>
                <span>{likes} likes</span>
                <span>{dislikes} dislikes</span>
              </div> */}

              {/* <div className={classes["comment-reaction"]}> */}
              {/* <button
                  onClick={likeHandler.bind(null, {commentId: id, postId})}
                >
                  <AiOutlineLike /> like
                </button> */}
              {/* <button
                  onClick={disLikeHandler.bind(null, {commentId: id, postId})}
                >
                  <AiOutlineDislike /> dislike
                </button> */}
              {/* </div> */}
            </div>
          )
        )}

        {/* {comments?.map((el, i)=> <p key={i}>{el.body}</p>)} */}
      </div>
    </div>
  );
};

export default Comment;
