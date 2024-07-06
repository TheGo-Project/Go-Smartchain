import React from "react";
import { Nav } from "react-bootstrap";
import { useAppSelector } from "./store/hooks";
import { Link } from "react-router-dom";

const NavMenu: React.FC = () => {
  const active = useAppSelector((state) => state.nav.active);

  return (
    <Nav activeKey={active}>
      <Nav.Item>
        <Nav.Link as={Link} eventKey="home" to="/">
          Home
        </Nav.Link>
      </Nav.Item>
      <Nav.Item>
        <Nav.Link as={Link} eventKey="explorer" to="/explorer">
          Explorer
        </Nav.Link>
      </Nav.Item>
    </Nav>
  );
};

export default NavMenu;
