// PrivateRoute.js
import React, {useContext} from "react";
import {Outlet, Navigate} from "react-router-dom";
import AuthContext from "../store/authContext";

const PrivateRoutes = () => {
  const {user} = useContext(AuthContext);
  console.log(user.isLogIn);
  return user.isLogIn ? <Outlet /> : <Navigate to={"/"} />;
};

export default PrivateRoutes;