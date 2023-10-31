// import { Link } from "react-router-dom";
import { Link } from "react-router-dom";
// import AuthContext from "../store/authContext";
// import  React, { useContext } from "react";
import { Nav, Navbar, NavItem, Container } from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";
// import { getJson, Url } from "../helpers/helpers";

const NavigationBar = () => {
  // const {isLoggedIn, username, onLogout} = useContext(AuthContext)
  // const logoutHandler =async ()=>{
  //   try{
  //     const headers = new Headers();
  //     headers.append("user_token", document.cookie);
  // const res =  await getJson(`logout`, {
  //   method: 'POST',
  //   headers: headers,
  //   credentials: "include",
  // })
  // if(res.success){
  //   onLogout()
  // }
  //   }catch(err){
  //     console.log(err)
  //   }
  // }
  return (
    <Navbar bg="dark" expand="lg" style={{ fontSize: "1em" }}>
      <Container>
        <Navbar.Brand style={{ color: "white" }}>
          <Link className="nav-link" style={{ color: "white" }} to="/">
            Forum
          </Link>
        </Navbar.Brand>
        <Navbar.Toggle
          aria-controls="basic-navbar-nav"
          style={{ backgroundColor: "white", borderColor: "black" }}
        />
        <Navbar.Collapse
          id="basic-navbar-nav"
          className="justify-content-end md"
        >
          <Nav>
            <>
              <NavItem>
                <Link
                  className="nav-link"
                  style={{ color: "white" }}
                  to="/login"
                >
                  Login
                </Link>
              </NavItem>
              <NavItem>
                <Link
                  className="nav-link"
                  style={{ color: "white" }}
                  to="/signup"
                >
                  Sign Up
                </Link>
              </NavItem>
            </>
            <NavItem>
              <span className="nav-link" style={{ color: "white" }}>
                John Doe
              </span>
            </NavItem>
            <NavItem>
              <button
                // onClick={logoutHandler}
                className="nav-link"
                style={{ color: "white", background: "none", border: "none" }}
              >
                logout
              </button>
            </NavItem>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};
export default NavigationBar;
