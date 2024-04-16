import React, {useCallback, useContext, useState} from "react";
import {Link} from "react-router-dom";
import {Nav, Navbar, Container, NavItem} from "react-bootstrap";
import AuthContext from "../store/authContext";
import SearchUser from "./SearchUser"; // Import the Search component
import classes from "./NavBar.module.css";
import {getJson} from "../helpers/helpers";

const NavigationBar = () => {
  const {user, onLogout} = useContext(AuthContext);
  const [seachList, setSeachList] = useState([]);
  const [typingTimeout, setTypingTimeout] = useState(null);

  // Add a state variable to store the search list
  const logoutHandler = () => {
    onLogout();
  };

  const fetchSearch = async (query) => {
    try {
      const data = await getJson("search", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({search: query}),
      });
      setSeachList(data);
    } catch (err) {
      console.log(err);
    }
  };

  const handleSearch = useCallback(
    (query) => {
      if (typingTimeout) {
        clearTimeout(typingTimeout);
      }

      setTypingTimeout(setTimeout(fetchSearch.bind(null, query), 500));
    },
    [typingTimeout]
  );

  return (
    <div className="position-relative">
      {/* Add a position-relative container */}
      <Navbar bg="dark" expand="lg" style={{fontSize: "1em", zIndex: 1000}}>
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
                    <Link
                      className={classes.NavLink}
                      to={`/profile/${user.Id}`}
                    >
                      <div className={`${classes["profile-image"]}`}>
                        {user.Avatar ? (
                          <img
                            src={`http://localhost:8000/images/${user.Avatar}`}
                            alt="Profile"
                            className={`${classes["img-fluid"]} ${classes["rounded-circle"]}`}
                          />
                        ) : (
                          <div className={`${classes["initials-circle"]}`}>
                            {user.Nickname}
                          </div>
                        )}
                      </div>
                    </Link>
                  </NavItem>
                  <NavItem>
                    <button
                      onClick={logoutHandler}
                      className={classes.NavLink}
                      style={{
                        color: "white",
                        background: "none",
                        border: "none",
                      }}
                    >
                      Logout
                    </button>
                  </NavItem>
                  <NavItem className=" mt-sm-3 ">
                    <SearchUser
                      onSearch={handleSearch}
                      searchList={seachList}
                      className={classes.NavLink}
                    />
                  </NavItem>
                </>
              )}
            </Nav>
          </Navbar.Collapse>
        </Container>
      </Navbar>
    </div>
  );
};

export default NavigationBar;
