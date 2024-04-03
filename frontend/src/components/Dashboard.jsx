import React, {useCallback, useContext, useEffect, useState} from "react";
import {Row, Col, Nav} from "react-bootstrap";
import Posts from "./Posts";
import {dummyPosts, groups} from "../store/dummydata";
import Groups from "./Groups";
import Chats from "./ChatLists";
import {getJson} from "../helpers/helpers";
import AuthContext from "../store/authContext";
function Dashboard() {
  const [activeSection, setActiveSection] = useState("posts");
  const [dashBoardData, setDashBoardData] = useState({});
  const {user} = useContext(AuthContext);
  const fetchDashboard = useCallback(async () => {
    try {
      console.log(user.Id, "dashboard");
      const data = await getJson("dashboard", {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
        },

        body: JSON.stringify(user.Id),
      });

      setDashBoardData(data);
    } catch (err) {
      console.log(err);
    }
  }, [user]);

  useEffect(() => {
    fetchDashboard();
  }, [fetchDashboard]);

  const handleSectionClick = (section, e) => {
    setActiveSection(section);

    // setActiveSection(section);
  };

  return (
    <>
      <Row className="mt-3">
        <Col>
          <Nav
            onSubmit={handleSectionClick.bind(null, "posts")}
            variant="pills"
            defaultActiveKey="posts"
            style={{
              gap: "10rem",
            }}
            className="d-flex justify-content-between "
          >
            <Nav.Item>
              <Nav.Link
                eventKey="posts"
                onClick={handleSectionClick.bind(null, "posts")}
              >
                Posts
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link
                eventKey="chats"
                onClick={handleSectionClick.bind(null, "chats")}
              >
                Chats
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link
                eventKey="groups"
                onClick={handleSectionClick.bind(null, "groups")}
              >
                Groups
              </Nav.Link>
            </Nav.Item>
          </Nav>
        </Col>
      </Row>

      <Row
        className="mt-3 d-flex justify-content-center align-items-center "
        style={{
          flexDirection: "column",
        }}
      >
        <Col>
          {activeSection === "posts" && <Posts posts={dummyPosts} />}
          {activeSection === "chats" && <Chats chats={dashBoardData?.chats} />}
          {activeSection === "groups" && <Groups groups={groups} />}
        </Col>
      </Row>
    </>
  );
}

export default Dashboard;
