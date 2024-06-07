import {createBrowserRouter, RouterProvider} from "react-router-dom";
import Root from "./pages/Root";
import Home from "./pages/Home";
import PrivateRoutes from "./components/PrivateRoutes";
import Profile from "./pages/Profile";
import Group from "./pages/Group";
import Post from "./pages/post";

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    children: [
      {
        path: "/",
        element: <Home />,
      },
      {
        path: "groups/:id",
        element: <PrivateRoutes component={Group} />,

        children: [
          {
            path: "post/:postId",
            element: <PrivateRoutes component={Post} />,
          },
        ],
      },
      {
        path: "profile/:id",
        // element: <ProfileLayout />,
        children: [
          {
            path: "",
            element: <PrivateRoutes component={Profile} />,
          },
          {
            path: "post/:postId",
            element: <PrivateRoutes component={Post} />,
          },
        ],
      },
      {
        path: "/post/:postId",
        element: <PrivateRoutes component={Post} />,
      },
    ],
  },
  {
    path: "*",
    element: <h1>404 Page Not Found</h1>,
  },
]);

function App() {
  return <RouterProvider router={router} />;
}

export default App;
