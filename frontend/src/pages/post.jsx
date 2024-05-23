import {useContext, useEffect, useState} from "react";
import Posts from "../components/Posts";
import AuthContext from "../store/authContext";
// import {useHistory} from "react-router-dom";
import {useLocation} from "react-router-dom";

const Post = () => {
  // const {} = useContext(AuthContext)

  const location = useLocation();
  console.log(location, "loacation");
  return <div>hello world</div>;
};

export default Post;
