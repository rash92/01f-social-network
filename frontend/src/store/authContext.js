import React, {useState, useEffect, useCallback} from "react";
import {getJson} from "../helpers/helpers";

const userObj = {
  id: "",
  isLogIn: false,
  username: "",
  profileImg: "",
};
const AuthContext = React.createContext({
  user: userObj,
  OnLogin: () => {},
  onLogout: () => {},
  OnAddCommentToPost: () => {},
  posts: [],
  catogaries: [],
  username: "",
  selectedPosts: [],
  setSelectedPosts: () => {},
  OnAddPost: () => {},
});

export const AuthContextProvider = (props) => {
  const [user, setUser] = useState(userObj);

  const [catogaries, SetCatogaries] = useState([]);
  const [posts, SetPosts] = useState([]);
  const [selectedPosts, setSelectedPosts] = useState([]);
  const logintHandler = (user) => {
    setUser({...user, isLogIn: true});
  };
  const logoutHandler = () => {
    setUser(userObj);
  };
  const checkSession = useCallback(async () => {
    try {
      const res = await getJson("checksession", {
        method: "GET",
        credentials: "include",
      });

      if (!res.success || res.status === 401) {
        logintHandler(user);
        throw new Error("something went wrong presist login");
      }

      logintHandler(res);
    } catch (err) {
      console.log("error", err);
    }
  }, []);
  const getCatogries = useCallback(async () => {
    try {
      const res = await getJson("get-catogries");
      SetCatogaries(res);
    } catch (error) {
      console.log(error.message);
    }
  }, []);
  const getPosts = useCallback(async () => {
    try {
      const res = await getJson("get-posts");
      SetPosts(res);
      setSelectedPosts(res);
    } catch (error) {
      console.log(error.message);
    }
  }, []);
  const OnAddCommentToPost = (postId, Comment) => {
    const {comments, ...rest} = posts.find((el) => el.id === postId);
    if (comments === null) {
      SetPosts(
        posts.map((el) =>
          el.id === postId
            ? {
                comments: [Comment],
                ...rest,
              }
            : el
        )
      );
    } else {
      SetPosts(
        posts.map((el) =>
          el.id === postId
            ? {
                comments: [Comment, ...comments],
                ...rest,
              }
            : el
        )
      );
    }
  };
  const OnAddPost = (post) => {
    if (post.length === 0) {
      SetPosts([post]);
    } else {
      SetPosts([post, ...posts]);
    }
  };
  const onAddLikeDislikePost = (id, data, quary) => {
    if (quary === "like") {
      const {likes, dislikes, userlikes, ...rest} = posts.find(
        (el) => el.id === id
      );
      SetPosts(
        posts.map((el) =>
          el.id === id
            ? {
                likes: data.likes,
                dislikes: data.dislikes,
                userlikes: data.userlikes,
                ...rest,
              }
            : el
        )
      );
    } else {
      const {likes, dislikes, userlikes, ...rest} = posts.find(
        (el) => el.id === id
      );
      SetPosts(
        posts.map((el) =>
          el.id === id
            ? {
                likes: data.likes,
                dislikes: data.dislikes,
                userlikes: data.userlikes,
                ...rest,
              }
            : el
        )
      );
    }
  };
  const onAddLikeDislikeComment = ({commentId, postId}, data, quary) => {
    // console.log(data, 'data likes')
    if (quary === "like" || "dislike") {
      const {comments, ...rest} = posts.find((el) => el.id === postId);
      const upDatedComments = comments.map((el) => {
        const {likes, dislikes, ...all} = el;
        return el.id === commentId
          ? {likes: data.likes, dislikes: data.dislikes, ...all}
          : el;
      });
      SetPosts(
        posts.map((el) =>
          el.id === postId ? {comments: upDatedComments, ...rest} : el
        )
      );
    }
  };
  const onRemovePost = (id) => {
    SetPosts(posts.filter((el) => el.id !== id));
  };
  useEffect(() => {
    checkSession();
    // getCatogries();
    // getPosts();
    checkSession();
  }, [checkSession]);
  // console.log(selectedPosts)
  return (
    <AuthContext.Provider
      value={{
        user,
        OnLogin: logintHandler,
        onLogout: logoutHandler,
        catogaries: catogaries,
        posts: posts,
        selectedPosts,
        setSelectedPosts,
        OnAddCommentToPost: OnAddCommentToPost,
        OnAddPost,
        onAddLikeDislikePost,
        onAddLikeDislikeComment,
        onRemovePost,
      }}
    >
      {props.children}
    </AuthContext.Provider>
  );
};
export default AuthContext;
