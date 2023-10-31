import NavigationBar from "./components/NavBar";
import React from "react";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import RootLayout from "./pages/Root";
import Home from "./pages/Home";

const router = createBrowserRouter([
  {
    path: "/",
    element: <RootLayout />,
    children: [
      {
        path: "/",
        element: <Home />,
      },
      {
        path: "/:id",
        element: <h1>hello world</h1>,
      },
      // {
      //   path: "/login",
      //   element: <LoginForm />,
      // },
      // {
      //   path: "/signup",
      //   element: <SignUp />,
      // },
    ],
  },
]);
function App() {
  return <RouterProvider router={router} />;
}

export default App;
