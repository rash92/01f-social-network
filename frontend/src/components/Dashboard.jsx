import React, {useCallback, useContext, useEffect, useState} from "react";
import {Row, Col, Nav} from "react-bootstrap";
import Posts from "./Posts";
import {dummyPosts, groups} from "../store/dummydata";
import Groups from "./Groups";
import Chats from "./ChatLists";
import {getJson} from "../helpers/helpers";
import AuthContext from "../store/authContext";
import NotificationContainer from "./Notifications";

const notifications = [
  {message: "Notification 1", variant: "success"},
  {message: "Notification 2", variant: "info"},
  {message: "Notification 3", variant: "warning"},
  {message: "Notification 4", variant: "danger"},
];
function Dashboard() {
  const [activeSection, setActiveSection] = useState("posts");

  const {fetchDashboard, dashBoardData} = useContext(AuthContext);

  useEffect(() => {
    fetchDashboard();
  }, []);

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

            <Nav.Item>
              <Nav.Link
                eventKey="Notifications"
                onClick={handleSectionClick.bind(null, "notifications")}
              >
                Notifications
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
          {activeSection === "posts" && <Posts posts={dashBoardData.Posts} />}
          {activeSection === "chats" && <Chats chats={dashBoardData?.chats} />}
          {activeSection === "groups" && (
            <Groups groups={dashBoardData?.groups} />
          )}
          {activeSection === "notifications" && (
            <NotificationContainer
              notifications={dashBoardData?.notifications}
            />
          )}
        </Col>
      </Row>
    </>
  );
}

export default Dashboard;
