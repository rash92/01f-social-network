import Login from "../components/Login";
import Register from "../components/Register";
import {Container} from "react-bootstrap";
import {useContext, useState} from "react";
import AuthContext from "../store/authContext";
import PostInput from "../components/PostInput";
import Dashboard from "../components/Dashboard";
import Chat from "../components/Chat";
import {useEffect} from "react";

function Home() {
  const {user, isWsReady, wsVal, onLogout, toggleDashboard} =
    useContext(AuthContext);
  const [showRegisterForm, setShowRegisterForm] = useState(false);

  useEffect(() => {
    toggleDashboard(true);
    return () => toggleDashboard(false);
  }, [toggleDashboard]);

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
              <PostInput src={user.Avatar} id={user.Id} />
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
