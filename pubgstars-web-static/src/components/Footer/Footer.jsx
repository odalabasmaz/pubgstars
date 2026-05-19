import React from "react";
import { Container, Row } from "react-bootstrap";
import { SocialIcon } from 'react-social-icons';

import "./Footer.css";

class Footer extends React.Component {
  render() {
    return (
      <footer className="Footer">
        <div style={{paddingBottom: '10px'}}>
          <Container>
            <Row className="justify-content-center">
              <ul className="social-media-list">
                <li><SocialIcon url="https://twitter.com/StarsPubg" fgColor="white" target="_blank" style={{height: 40, width: 40}}/></li>
                <li><SocialIcon url="https://instagram.com/pubgstarsofficial" fgColor="white" target="_blank" style={{height: 40, width: 40}}/></li>
                <li><SocialIcon url="https://youtube.com/channel/UCsQB8Xe4DMttuYx_t8OOE-w?view_as=subscriber" fgColor="white" target="_blank" style={{height: 40, width: 40}}/></li>
                <li><SocialIcon url="https://facebook.com/pubgstarrs" fgColor="white" target="_blank" style={{height: 40, width: 40}}/></li>
              </ul>
            </Row>
          </Container>
        </div>
        <div className="footer-copyright text-center">
          <Container>
            <div>&copy; Copyright 2019 PUBGSTARS. Tüm Hakları Saklıdır.</div>
          </Container>
        </div>
      </footer>
    );
  }
}

export default Footer;
