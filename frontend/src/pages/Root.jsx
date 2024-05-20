import React, { useContext } from "react";
import {Outlet} from "react-router-dom";
import NavigationBar from "../components/NavBar";
import NotificationsToest from "../components/NotificationsToest";

// import Footer from "../components/Footer"
export default function Root() {
  return (
    <>
      <NavigationBar />
      <NotificationsToest />
      <Outlet />
      {/* <Footer /> */}
    </>
  );
}
