import React from "react";
import { Outlet } from "react-router-dom";
import NavigationBar from "../components/NavBar";
// import Footer from "../components/Footer"
const RootLayout = () => {
  return (
    <>
      <NavigationBar />
      <Outlet />
      {/* <Footer /> */}
    </>
  );
};
export default RootLayout;
