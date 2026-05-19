import React from "react";
import { Nav, Navbar} from "react-bootstrap";
import { Container, Row, Col } from "react-bootstrap";

import "./Header.css";

const AVATAR = 'https://www.gravatar.com/avatar/429e504af19fc3e1cfa5c4326ef3394c?s=240&d=mm&r=pg';

const Header = () => (
    <header>
        <Navbar>
            <Container>
                <Row className>
                    <Col>
                        <img src={AVATAR} alt="avatar" className="img-fluid rounded-circle" style={{ width: 36 }} />
                    </Col>
                    <Col>
                        <img src={AVATAR} alt="avatar" className="img-fluid rounded-circle" style={{ width: 36 }} />
                    </Col>
                    <Col>
                        <img src={AVATAR} alt="avatar" className="img-fluid rounded-circle" style={{ width: 36 }} />
                    </Col>
                </Row>
            </Container>
        </Navbar>
    </header>
);

export default Header;
