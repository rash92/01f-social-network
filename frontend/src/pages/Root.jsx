import React from "react";
import { Outlet } from "react-router-dom";
import NavigationBar from "../components/NavBar";
// import Footer from "../components/Footer"
export default function Root() {
  return (
    <>
      <NavigationBar />
      <Outlet />
      {/* <Footer /> */}
    </>
  );
}
