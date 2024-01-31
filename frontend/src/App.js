import {createBrowserRouter, RouterProvider} from "react-router-dom";
import Root from "./pages/Root";
import Home from "./pages/Home";
import PrivateRoutes from "./components/PrivateRoutes";
import Profile from "./pages/Profile";

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

      // {
      //   path: "/groups",
      //   element: <PrivateRoutes component={f} />,
      // },
      {
        path: "/profile/:id",
        element: <PrivateRoutes component={Profile} />,
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
