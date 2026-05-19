import React from "react";
import {NavLink} from "react-router-dom";
import {Nav, Navbar, Button, Container} from "react-bootstrap";
import Routes from "../../Routes";
import Footer from "../../components/Footer/Footer";
import logo from "../../images/pubg.png";

import "./PublicLayout.css";

export default class PublicLayout extends React.Component {

  render() {
    return (
      <div className="App">
        <header style={{borderBottom: '2px solid #474d66'}}>
          <Container>
            <Navbar collapseOnSelect expand="lg">
              <Navbar.Brand>
                <a href='/'>
                  <img src={logo} className="navbar-brand-img" alt="PubgStars"/>
                </a>
              </Navbar.Brand>
              <Navbar.Toggle aria-controls="responsive-navbar-nav"/>
              <Navbar.Collapse id="responsive-navbar-nav">
                <>
                  <Nav className="mr-auto">
                    <Nav.Link eventKey="0" as={NavLink} to='/' exact>Anasayfa</Nav.Link>
                    <Nav.Link eventKey="1" as={NavLink} to='/about' exact>Hakkında</Nav.Link>
                    <Nav.Link eventKey="2" as={NavLink} to='/rules' exact>Kurallar</Nav.Link>
                    <Nav.Link eventKey="3" as={NavLink} to='/sss' exact>SSS</Nav.Link>
                    <Nav.Link eventKey="4" as={NavLink} to='/contact' exact>İletişim</Nav.Link>
                  </Nav>
                  <Nav>
                    <>
                      <Nav.Link eventKey="5" as={NavLink} to="/signup">
                        <Button variant="outline-success">Üye Ol</Button>
                      </Nav.Link>
                      <Nav.Link eventKey="6" as={NavLink} to="/login">
                        <Button variant="outline-success">Giriş Yap</Button>
                      </Nav.Link>
                    </>
                  </Nav>
                </>
              </Navbar.Collapse>
            </Navbar>
          </Container>
        </header>
        <Routes childProps={this.props.childProps}/>
        <Footer/>
      </div>
    );
  }

}
