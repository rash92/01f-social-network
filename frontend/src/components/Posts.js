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
      {posts?.map((el) => (
        <Post
          title={el.Title}
          body={el.Body}
          CreatedAt={el.CreatedAt}
          comments={el?.Comments}
          likes={el?.Likes}
          dislikes={el?.Dislikes}
          CreatorNickname={el?.CreatorNickname}
          likeHandler={likeHandler}
          disLikeHandler={disLikeHandler}
          id={el.Id}
          key={el.Id}
        />
      ))}
    </>
  );
};
export default Posts;
