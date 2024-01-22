import Login from "../components/Login";
import Register from "../components/Register";
import { Container } from "react-bootstrap";
import { useContext, useState } from "react";
import AuthContext from "../store/authContext";

function Home() {
  const {user} = useContext(AuthContext)

  const [showRegisterForm, setShowRegisterForm] = useState(false);
  return (
    <Container>
      {!user.isLogIn &&
        (showRegisterForm ? (
          <Register setShowRegisterForm={setShowRegisterForm} />
        ) : (
          <Login setShowRegisterForm={setShowRegisterForm} />
        ))}
    </Container>
  );
}

export default Home;
