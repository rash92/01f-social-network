import {createBrowserRouter, RouterProvider} from "react-router-dom";
import Root from "./pages/Root";
import Home from "./pages/Home";
import PrivateRoutes from "./components/PrivateRoutes";
import Profile, {profileLoader} from "./pages/Profile";
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
        path: "/:id",
        element: <h1>404 Page Not Found</h1>,
      },

      {
        path: "/groups/:id",
        element: <PrivateRoutes component={Group} />,
      },
      {
        path: "/profile/:id",
        errorElement: <Profile />,
        element: <PrivateRoutes component={Profile} />,
        // loader: profileLoader,
      },

      {
        path: "/post/:id",
        errorElement: <Post />,
        element: <PrivateRoutes component={Post} />,
        // loader: profileLoader,
      },

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
