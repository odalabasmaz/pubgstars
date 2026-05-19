import React from "react";
import {Auth} from "aws-amplify";
import {NavLink, withRouter} from "react-router-dom";
import {Navbar, Container, Dropdown} from "react-bootstrap";
import {FiLogOut, FiUser} from "react-icons/fi";
import UserToggle from "../../components/UserToggle/UserToggle";
import Routes from "../../Routes";
import Footer from "../../components/Footer/Footer";
import logo from "../../images/pubg.png";

import Nav from 'react-bootstrap/Nav'

import "./SigninLayout.css";
import API from "../../utils/api";

class SigninLayout extends React.Component {

    constructor(props, context) {
        super(props, context);
        this.state = {}
    }

    async componentDidMount() {
        let user = await Auth.currentAuthenticatedUser();
        this.setState({user: user});
        const options = {
            headers: {
                Authorization: user.signInUserSession.idToken.jwtToken
            }
        };

        API.getUser(options)
            .then(res => {
                let userInfo = res.data.body;
                this.setState({userInfo});
            });
    }

    handleLogout = async event => {
        await Auth.signOut();
        this.props.childProps.userHasAuthenticated(false);
    };

    handleProfilePage = event => {
        this.props.history.push("/profile")
    };

    render() {
        return (
            <div className="App">
                <header>
                    <Container>
                        <Navbar expand="lg">
                            {/*<HamburgerButton/>*/}
                            <Navbar.Brand>
                                <a href='/'>
                                    <img src={logo} className="navbar-brand-img" alt="PubgStars"/>
                                </a>
                            </Navbar.Brand>
                            <div className="ml-auto">
                                <div style={{display: 'flex', flexFlow: 'row'}}>
                                    <div style={{display: 'flex', justifyContent: 'center', flexFlow: 'column', fontSize: '13px', marginRight: '5px'}}>
                                        <div>Hoş geldin <b>{this.state.user && this.state.user.attributes.name}!</b>
                                        </div>
                                        <div>Bakiye: <b>{this.state.userInfo && (this.state.userInfo.balance + this.state.userInfo.bonus)}₺</b>
                                        </div>
                                    </div>
                                    <Dropdown alignRight>
                                        <Dropdown.Toggle as={UserToggle} id="dropdown-custom-components"/>
                                        <Dropdown.Menu>
                                            <Dropdown.Header>
                                                <div style={{display: 'flex', alignItems: 'center'}}>
                                                    <i className="fa fa-user-circle fa-2x" aria-hidden="true"/>
                                                    <div style={{marginLeft: '10px'}}>
                                                        <p style={{fontSize: "12px", margin: "0"}}>
                                                            <b>{this.state.user && this.state.user.attributes.name}</b>
                                                        </p>
                                                        <p style={{fontSize: "12px", margin: "0"}}>
                                                            {this.state.user && this.state.user.attributes.email}
                                                        </p>
                                                    </div>
                                                </div>
                                            </Dropdown.Header>
                                            <Dropdown.Divider/>
                                            <Dropdown.Item eventKey="1" onClick={this.handleProfilePage}>
                                                <FiUser/> Profilim
                                            </Dropdown.Item>
                                            <Dropdown.Item eventKey="2" onClick={this.handleLogout}>
                                                <FiLogOut/> Çıkış
                                            </Dropdown.Item>
                                        </Dropdown.Menu>
                                    </Dropdown>
                                </div>
                            </div>
                        </Navbar>
                    </Container>
                </header>
                <section className="header-bottom" style={{backgroundColor: '#185aa3', color:'white'}}>
                    <Nav className="justify-content-center" activeKey="/">
                        <Nav.Link as={NavLink} to="/" eventKey="1" exact>Anasayfa</Nav.Link>
                        <Nav.Link as={NavLink} to="/oyunlarim" eventKey="2">Oyunlarım</Nav.Link>
                        <Nav.Link as={NavLink} to="/leaderboard" eventKey="3">Kazananlar</Nav.Link>
                        <Nav.Link as={NavLink} to="/balance" eventKey="4">Bakiye İşlemleri</Nav.Link>
                        <Nav.Link as={NavLink} to="/transactionLog" eventKey="5">İşlem Geçmişi</Nav.Link>
                        <Nav.Link as={NavLink} to="/app/contact" eventKey="6">İletişim</Nav.Link>
                        <Nav.Link as={NavLink} to="/app/rules" eventKey="7">Kurallar</Nav.Link>
                    </Nav>
                </section>
                <Container>
                    <div className='mainSection'>
                        <Routes childProps={this.props.childProps}/>
                    </div>
                </Container>
                <Footer/>
            </div>
        );
    }
}

export default withRouter(SigninLayout);
