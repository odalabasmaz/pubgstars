import React, { Component } from "react";
import { Link } from "react-router-dom";
import { Auth } from "aws-amplify";
import { Form, FormGroup, FormControl, FormLabel } from "react-bootstrap";
import LoaderButton from "../components/LoaderButton";
import "./Login.css";

export default class Login extends Component {
  constructor(props) {
    super(props);

    this.state = {
      isLoading: false,
      email: "",
      password: ""
    };
  }

  validateForm() {
    return this.state.email.length > 0 && this.state.password.length > 0;
  }

  handleChange = event => {
    this.setState({
      [event.target.id]: event.target.value
    });
  };

  handleSubmit = async (event) => {
    event.preventDefault();
    this.setState({isLoading: true});
    try {
      await Auth.signIn(this.state.email, this.state.password);
      this.props.userHasAuthenticated(true)
    } catch (e) {
      this.setState({error: true});
      this.setState({isLoading: false})
    }
  };

  render() {
    return (
      <div className="Login">
        <div className="LoginForm">
          <div className="form-title">Üye Girişi</div>
          <Form onSubmit={this.handleSubmit}>
            <FormGroup controlId="email">
              <FormLabel>E-posta Adresi</FormLabel>
              <FormControl
                  autoFocus
                  type="email"
                  value={this.state.email}
                  onChange={this.handleChange}
              />
            </FormGroup>
            <FormGroup controlId="password">
              <FormLabel>Şifre</FormLabel>
              <FormControl
                value={this.state.password}
                onChange={this.handleChange}
                type="password"
              />
            </FormGroup>
            {this.state.error && <div style={{color:'red', fontSize: '13px', marginBottom:'10px'}}>E-posta adresiniz veya şifreniz hatalı!</div>}
            <LoaderButton
              block
              disabled={!this.validateForm()}
              type="submit"
              isLoading={this.state.isLoading}
              text="Giriş Yap"
              loadingText="Giriş Yapılıyor…"
            />
          </Form>
          <Link to="/login/reset">Şifremi Unuttum</Link>
          <div style={{textAlign:'center', paddingTop:'10px', marginTop:'20px', borderTop: '1px solid #E8E8E8', fontSize:'15px'}}>
            Hesabın yok mu?<Link to='signup'> Üye Ol</Link>
          </div>
        </div>
      </div>
    );
  }
}
