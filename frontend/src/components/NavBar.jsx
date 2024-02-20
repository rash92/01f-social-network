// import { Link } from "react-router-dom";
import {Link} from "react-router-dom";
// import AuthContext from "../store/authContext";
// import  React, { useContext } from "react";
import {Nav, Navbar, NavItem, Container} from "react-bootstrap";
import "bootstrap/dist/css/bootstrap.min.css";
// import { getJson, Url } from "../helpers/helpers";
import classes from "./NavBar.module.css";
import AuthContext from "../store/authContext";
import {useContext} from "react";

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

  const {user} = useContext(AuthContext);

  return (
    <Navbar bg="dark" expand="lg" style={{fontSize: "1em"}}>
      <Container>
        <Navbar.Brand style={{color: "white"}}>
          <Link className="nav-link" style={{color: "white"}} to="/">
            Forum
          </Link>
        </Navbar.Brand>
        <Navbar.Toggle
          aria-controls="basic-navbar-nav"
          style={{backgroundColor: "white", borderColor: "black"}}
        />
        <Navbar.Collapse
          id="basic-navbar-nav"
          className="justify-content-end md"
        >
          <Nav className={classes.NavBarFlex}>
            {user.isLogIn && (
              <>
                <NavItem>
                  <Link className={classes.NavLink} to={`/profile/${user.id}`}>
                    <div className={`${classes["profile-image"]}`}>
                      {user.profileImg ? (
                        <img
                          src={`http://localhost:8000/images/${user.profileImg}`}
                          alt="Profile"
                          className={`${classes["img-fluid"]} ${classes["rounded-circle"]}`}
                        />
                      ) : (
                        <div className={`${classes["initials-circle"]}`}>
                          {user.username[0]}
                        </div>
                      )}
                    </div>
                  </Link>
                </NavItem>
                <NavItem>
                  <button
                    // onClick={logoutHandler}
                    className={classes.NavLink}
                    style={{color: "white", background: "none", border: "none"}}
                  >
                    Logout
                  </button>
                </NavItem>
              </>
            )}
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
};
export default NavigationBar;
