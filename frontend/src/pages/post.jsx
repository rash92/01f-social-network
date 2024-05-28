import {useContext, useEffect, useState} from "react";
import Post from "../components/post";
import AuthContext from "../store/authContext";
import {useParams} from "react-router-dom";
import {getJson, reactlikeDislike} from "../helpers/helpers";
import {Container} from "react-bootstrap";

const PostPage = () => {
  const {user} = useContext(AuthContext);
  const [post, setPost] = useState({});
  const likeDislikeHandler = async ({id, query}) => {
    try {
      const res = await reactlikeDislike({id: user.Id, query, postId: id});

      setPost((prev) => ({
        ...prev,
        Likes: res.Likes,
        Dislikes: res.Dislikes,
        UserLikeDislike: res.UserLikeDislike,
      }));
    } catch (err) {
      console.log(err);
    }
  };
  const {postId} = useParams();

  useEffect(() => {
    const fetchPost = async (id) => {
      try {
        const response = await getJson("get-post", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          credentials: "include",
          body: JSON.stringify({
            post_id: id,
            user_id: user.Id,
          }),
        });

        setPost(response);
      } catch (err) {
        console.log("error", err);
      }
    };

    fetchPost(postId);
  }, [postId, user.Id]);
  const addComment = async (data) => {
    setPost((prev) => ({
      ...prev,
      Comments: Array.isArray(prev.Comments)
        ? [data, ...prev.Comments]
        : [data],
    }));
  };

  return (
    <Container style={{maxWidth: "800px"}}>
      <div className="d-flex justify-content-center align-items-center">
        <Post
          type={true}
          title={post.Title}
          body={post.Body}
          CreatedAt={post.CreatedAt}
          comments={post?.Comments}
          likes={post?.Likes}
          dislikes={post?.Dislikes}
          CreatorNickname={post?.CreatorNickname}
          likeDislikeHandler={likeDislikeHandler}
          UserLikeDislike={post.UserLikeDislike}
          addComment={addComment}
          Image={post.Image}
          id={post.Id}
          key={post.Id}
        />
      </div>
    </Container>
  );
};

export default PostPage;
