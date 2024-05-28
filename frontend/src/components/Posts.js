import Post from "./post.js";

import {reactlikeDislike} from "../helpers/helpers";
import {useContext} from "react";
import AuthContext from "../store/authContext.js";

const Posts = ({posts}) => {
  const {user, onAddLikeDislikePost} = useContext(AuthContext);
  const likeDislikeHandler = async ({id, query}) => {
    try {
      const res = await reactlikeDislike({id: user.Id, query, postId: id});

      onAddLikeDislikePost(id, res);
    } catch (err) {
      console.log(err);
    }
  };

  return (
    <>
      {posts?.map((el) => (
        <Post
          type={false}
          title={el.Title}
          body={el.Body}
          CreatedAt={el.CreatedAt}
          comments={el?.Comments}
          likes={el?.Likes}
          dislikes={el?.Dislikes}
          CreatorNickname={el?.CreatorNickname}
          likeDislikeHandler={likeDislikeHandler}
          UserLikeDislike={el.UserLikeDislike}
          Image={el.Image}
          id={el.Id}
          key={el.Id}
        />
      ))}
    </>
  );
};
export default Posts;
