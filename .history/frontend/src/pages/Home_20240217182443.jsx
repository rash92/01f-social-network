import Login from "../components/Login";
import Register from "../components/Register";
import {Container} from "react-bootstrap";
import {useContext, useState} from "react";
import AuthContext from "../store/authContext";
import PostInput from "../components/PostInput";
import Dashboard from "../components/Dashboard";
import {Link} from "react-router-dom";
import Chat from "../components/Chat";

function Home() {
  const {user} = useContext(AuthContext);
  const [showRegisterForm, setShowRegisterForm] = useState(false);
  const [show, setShow] = useState(false);
  const handleClose = () => {
    setShow(false);
  };
  const handleShow = () => {
    setShow(true);
  };

 
  return (
    <Container>
      <div
        style={{
          display: "flex",
          flexDirection: "column",
          justifyContent: "center",
          alignItems: "center",
          backgroundColor: "#fcfcfc",
        }}
      >
        {!user.isLogIn &&
          (showRegisterForm ? (
            <Register setShowRegisterForm={setShowRegisterForm} />
          ) : (
            <Login setShowRegisterForm={setShowRegisterForm} />
          ))}

        {user.isLogIn && (
          <>
            <div style={{width: "50vw", margin: "2rem 0 3rem 0"}}>
              <PostInput src={user.profileImg} id={user.id} />
            </div>

            <Chat />

            <Dashboard />
          </>
        )}
      </div>
    </Container>
  );
}

export default Home;