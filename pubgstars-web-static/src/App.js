import React, {Component} from "react";
import {Auth} from "aws-amplify";

import PublicLayout from "./components/PublicLayout/PublicLayout";
import SigninLayout from "./components/PublicLayout/SigninLayout";

import "./App.css";

class App extends Component {
  constructor(props) {
    super(props);

    this.state = {
      isAuthenticated: false,
      isAuthenticating: true
    };
  }

  async componentDidMount() {
    try {
      //The Amplify client will refresh the tokens calling Auth.currentSession if they are no longer valid.
      await Auth.currentSession();
      this.userHasAuthenticated(true);
    } catch (e) {
      if (e !== 'No current user') {
      }
    }
    this.setState({isAuthenticating: false});
  }

  userHasAuthenticated = authenticated => {
    this.setState({isAuthenticated: authenticated});
  };

  render() {
    const childProps = {
      isAuthenticated: this.state.isAuthenticated,
      userHasAuthenticated: this.userHasAuthenticated
    };

    return (
      !this.state.isAuthenticating &&
      (this.state.isAuthenticated ? <SigninLayout childProps={childProps}/> : <PublicLayout childProps={childProps}/>)
    );
  }
}

export default App;
