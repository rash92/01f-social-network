import Post from "./post.js";
import AuthContext from "../store/authContext";
import {useContext} from "react";
import {getJson} from "../helpers/helpers";
import {useNavigate} from "react-router-dom";
const reactlikeDislike = async (postId, quary) => {
  try {
    return await getJson("react-Post-like-dislike", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        user_token: document.cookie,
      },
      credentials: "include",
      body: JSON.stringify({
        postId,
        quary,
      }),
    });
  } catch (err) {
    throw err;
  }
};
const Posts = ({posts,  onAddLikeDislikePost}) => {
  const Navigate = useNavigate();
  const {selectedPosts, onAddLikeDislikePost, isLoggedIn} =
    useContext(AuthContext);
  const likeHandler = async (id, e) => {
    try {
      const res = await reactlikeDislike(id, "like");
      onAddLikeDislikePost(id, res, "like");
    } catch (err) {
      console.log(err);
    }
  };
  const disLikeHandler = async (id, e) => {
    try {
      if (!isLoggedIn) {
        Navigate("/login");
      } else {
        const res = await reactlikeDislike(id, "dislike");
        onAddLikeDislikePost(id, res, "dislike");
      }
    } catch (err) {
      console.log(err);
    }
  };
  return (
    <>
      {posts?.map((el) => (
        <Post
          id={el.id}
          title={el.title}
          body={el.body}
          categories={el.categories}
          created_at={el.created_at}
          comments={el.comments}
          likes={el.likes}
          dislikes={el.dislikes}
          username={el.username}
          likeHandler={likeHandler}
          disLikeHandler={disLikeHandler}
          key={el.id}
        />
      ))}
    </>
  );
};
export default Posts;