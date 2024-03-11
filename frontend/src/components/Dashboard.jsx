import React, {useState} from "react";
import {Row, Col, Nav} from "react-bootstrap";
import Posts from "./Posts";
import {dummyPosts, groups} from "../store/dummydata";
import Groups from "./Groups";
import Chats from "./ChatLists";

function Dashboard() {
  const [activeSection, setActiveSection] = useState("posts");
  const handleSectionClick = (section) => {
    setActiveSection(section);
  };

  return (
    <>
      <Row className="mt-3">
        <Col>
          <Nav
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
                onClick={() => handleSectionClick("posts")}
              >
                Posts
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link
                eventKey="chats"
                onClick={() => handleSectionClick("chats")}
              >
                Chats
              </Nav.Link>
            </Nav.Item>
            <Nav.Item>
              <Nav.Link
                eventKey="groups"
                onClick={() => handleSectionClick("groups")}
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
          {activeSection === "chats" && <Chats />}
          {activeSection === "groups" && <Groups groups={groups} />}
        </Col>
      </Row>
    </>
  );
}

export default Dashboard;
