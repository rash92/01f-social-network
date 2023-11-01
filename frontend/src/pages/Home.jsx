import Login from "../components/Login";
import Register from "../components/Register";
import { Container } from "react-bootstrap";
import { useState } from "react";

function Home() {
  const [loggedIn, setLogin] = useState(false);
  const [showRegisterForm, setShowRegisterForm] = useState(false);
  return (
    <Container>
      {!loggedIn &&
        (showRegisterForm ? (
          <Register setShowRegisterForm={setShowRegisterForm} />
        ) : (
          <Login setShowRegisterForm={setShowRegisterForm} />
        ))}
    </Container>
  );
}

export default Home;
