import Post from "./post.js";

import {getJson} from "../helpers/helpers";

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
const Posts = ({posts, onAddLikeDislikePost}) => {
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
      const res = await reactlikeDislike(id, "dislike");
      onAddLikeDislikePost(id, res, "dislike");
    } catch (err) {
      console.log(err);
    }
  };
  return (
    <>
      {posts?.map((el, i) => (
        <Post
          id={i}
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
          key={i}
        />
      ))}
    </>
  );
};
export default Posts;
