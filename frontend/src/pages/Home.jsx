import Login from "../components/Login";
import Register from "../components/Register";
import {Container} from "react-bootstrap";
import {useContext, useState} from "react";
import AuthContext from "../store/authContext";
import AddPost from "../components/AddPost";

function Home() {
  const {user} = useContext(AuthContext);
  const [showRegisterForm, setShowRegisterForm] = useState(false);
  return (
    <Container>
      {!user.isLogIn &&
        (showRegisterForm ? (
          <Register setShowRegisterForm={setShowRegisterForm} />
        ) : (
          <Login setShowRegisterForm={setShowRegisterForm} />
        ))}

      {user.isLogIn && <AddPost />}
    </Container>
  );
}

export default Home;
