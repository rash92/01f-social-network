import Login from "../components/Login";
import Register from "../components/Register";
import {Container} from "react-bootstrap";
import {useCallback, useContext, useEffect, useState} from "react";
import AuthContext from "../store/authContext";
import PostInput from "../components/PostInput";
import Dashboard from "../components/Dashboard";
import Chat from "../components/Chat";
import {getJson} from "../helpers/helpers";
import {Link} from "react-router-dom";
function Home() {
  const {user, isWsReady, wsVal, onLogout} = useContext(AuthContext);
  const [showRegisterForm, setShowRegisterForm] = useState(false);

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
              <PostInput src={user.Profile} id={user.Id} />
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
