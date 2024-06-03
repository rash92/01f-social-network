import {
  AiOutlineLike,
  AiOutlineDislike,
  AiOutlineComment,
  AiFillDelete,
} from "react-icons/ai";
import moment from "moment";
import classes from "./Post.module.css";
import {Col, Image, Button} from "react-bootstrap";

// import CommentForm from "./CommentForm";
import {useState, useContext} from "react";
import {getJson} from "../helpers/helpers";
import AuthContext from "../store/authContext";
import {Link} from "react-router-dom";
import CommentForm from "./CommentForm";
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
  CreatedAt,
  comments,
  likes,
  dislikes,
  likeDislikeHandler,
  CreatorNickname,
  UserLikeDislike,
  Image: image,
  type,
  addComment,
}) => {
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
  const [showComment, setShowComent] = useState();

  const htmlComment = !type ? (
    <Link to={`post/${id}`} style={{textDecoration: "none"}}>
      <span>{comments?.length} comments</span>
    </Link>
  ) : (
    <button>
      <span>
        <AiOutlineComment />
      </span>
      <span>{comments?.length} comments</span>
    </button>
  );

  const htmlComment1 = !type ? (
    <Link to={`post/${id}`} style={{textDecoration: "none"}}>
      <span>
        <AiOutlineComment />
      </span>
      <span>Comment</span>
    </Link>
  ) : (
    <button>
      <span>
        <AiOutlineComment />
      </span>
      <span>Comment</span>
    </button>
  );

  const timeago = moment(CreatedAt).fromNow();
  return (
    <Col xs={12} key={id}>
      <div className={`${classes.post}`}>
        <div className={classes["post-head"]}>
          <span className="user name">
            posted by {CreatorNickname} {timeago}
          </span>
          {isLoggedIn && LoggedInUser === CreatorNickname && (
            <button onClick={deleteHandler}>
              <AiFillDelete />
            </button>
          )}
        </div>

        <div className={classes["post-title"]}>
          <h4>{title}</h4>
        </div>

        <div className={classes["post-body"]}>
          <p> {body}</p>

          {image && (
            <div style={{width: "100%", overflow: "hidden"}}>
              <Image
                src={`http://localhost:8000/images/${image}`}
                width={500}
                height={500}
                style={{width: "100%", objectFit: "cover"}}
              />
            </div>
          )}
        </div>

        <div
          className={`${classes["post-reaction-info"]} d-flex justify-content-between`}
        >
          {likes > 0 && <span>{likes} likes</span>}
          {dislikes > 0 && <span> {dislikes} dislikes</span>}
          {comments?.length > 0 && htmlComment}
        </div>

        <div
          className={`${classes["post-reaction"]} d-flex justify-content-between`}
        >
          <button onClick={likeDislikeHandler.bind(null, {id, query: "like"})}>
            <div className={classes["post-reaction-info-reaction"]}>
              <AiOutlineLike
                color={UserLikeDislike === 1 ? "green" : "black"}
                size={24}
              />
              <span>like</span>
            </div>
          </button>
          <button
            onClick={likeDislikeHandler.bind(null, {id, query: "dislike"})}
          >
            <div className={classes["post-reaction-info-reaction"]}>
              <AiOutlineDislike
                color={UserLikeDislike === -1 ? "red" : "black"}
                size={24}
              />
              <span>dislike</span>
            </div>
          </button>

          {htmlComment1}
        </div>

        {type && (
          <CommentForm id={id} comments={comments} addComment={addComment} />
        )}
      </div>
    </Col>
  );
};

export default Post;
