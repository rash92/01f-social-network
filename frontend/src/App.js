import {createBrowserRouter, RouterProvider} from "react-router-dom";
import Root from "./pages/Root";
import Home from "./pages/Home";
import PrivateRoutes from "./components/PrivateRoutes";
import Profile, {profileLoader} from "./pages/Profile";
const FakeComponet = () => {
  return <div>Fake Component</div>;
};
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
        path: "/groups",
        element: <PrivateRoutes component={FakeComponet} />,
      },
      {
        path: "/profile/:id",
        errorElement: <Profile />,
        element: <PrivateRoutes component={Profile} />,
        loader: profileLoader,
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
