import {
  AiOutlineLike,
  AiOutlineDislike,
  AiOutlineComment,
  AiFillDelete,
} from "react-icons/ai";
import moment from "moment";
import classes from "./Post.module.css";
import {Col} from "react-bootstrap";
import CommentForm from "./CommentForm";
import {useState, useContext} from "react";
import {getJson} from "../helpers/helpers";
import AuthContext from "../store/authContext";

const removePostHandler = async (postId) => {
  try {
    return await getJson("removePost", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        user_token: document.cookie,
      },
      credentials: "include",
      body: JSON.stringify({
        postId,
      }),
    });
  } catch (err) {
    throw err;
  }
};

const Post = ({
  id,
  title,
  body,
  created_at,
  comments,
  likes,
  dislikes,
  username,
  likeHandler,
  disLikeHandler,
}) => {
  const [ShowComment, setShowComment] = useState(false);
  const {
    isLoggedIn,
    username: LoggedInUser,
    onRemovePost,
  } = useContext(AuthContext);
  const deleteHandler = async (e) => {
    try {
      const res = await removePostHandler(id);
      onRemovePost(res.postId);
    } catch (err) {
      console.log(err);
    }
  };

  const timeago = moment(created_at).fromNow();
  return (
    <Col xs={12} key={id}>
      <div className={`${classes.post}`}>
        <div className={classes["post-head"]}>
          <span className="user name">
            posted by {username} {timeago}
          </span>
          {isLoggedIn && LoggedInUser === username && (
            <button onClick={deleteHandler}>
              <AiFillDelete />
            </button>
          )}
        </div>

        <div className={classes["post-title"]}>
          <h4>
            {" "}
            {title} post {id}
          </h4>
        </div>

        <div className={classes["post-body"]}>
          <p> {body}</p>
        </div>

        <div
          className={`${classes["post-reaction-info"]} d-flex justify-content-between`}
        >
          {likes > 0 && <span>{likes} likes</span>}
          {dislikes > 0 && <span> {dislikes} dislikes</span>}
          {comments?.length > 0 && (
            <span onClick={() => setShowComment(true)}>
              {" "}
              {comments?.length} comments
            </span>
          )}
        </div>

        <div
          className={`${classes["post-reaction"]} d-flex justify-content-between`}
        >
          <button onClick={likeHandler.bind(null, id)}>
            <div className={classes["post-reaction-info-reaction"]}>
              <AiOutlineLike size={24} />
              <span>like</span>
            </div>
          </button>
          <button onClick={disLikeHandler.bind(null, id)}>
            <div className={classes["post-reaction-info-reaction"]}>
              <AiOutlineDislike size={24} />
              <span>dislike</span>
            </div>
          </button>

          <button onClick={() => setShowComment(true)}>
            <div className={classes["post-reaction-info-reaction"]}>
              <AiOutlineComment size={24} />
              <span>comment</span>
            </div>
          </button>
        </div>

        {ShowComment && <CommentForm id={id} comments={comments} />}
      </div>
    </Col>
  );
};

export default Post;
