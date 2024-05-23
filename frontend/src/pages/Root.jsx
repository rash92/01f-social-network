import React from "react";
import {Outlet} from "react-router-dom";
import NavigationBar from "../components/NavBar";
import NotificationsToest from "../components/NotificationsToest";
import ChatComponent from "../components/Chat";

// import Footer from "../components/Footer"
export default function Root() {
  return (
    <>
      <NavigationBar />
      <NotificationsToest />
      <ChatComponent />
      <Outlet />
      {/* <Footer /> */}
    </>
  );
}
