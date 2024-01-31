// PrivateRoute.js
import React, {useContext} from "react";
import {Navigate} from "react-router-dom";
import AuthContext from "../store/authContext";

const PrivateRoutes = ({component: Component, authenticated, ...rest}) => {
  // Add your authentication logic here
  const {user} = useContext(AuthContext);
  return user.isLogIn ? <Component {...rest} /> : <Navigate to="/" replace />;
};
export default PrivateRoutes;
