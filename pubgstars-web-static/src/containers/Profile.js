import React, {Component} from "react";
import ChangePassword from "./ChangePassword";

import "./Profile.css";
import {Auth} from "aws-amplify";
import API from "../utils/api";

export default class Profile extends Component {
    constructor(props) {
        super(props);

        this.state = {};
    }

    async componentDidMount() {
        let user = await Auth.currentAuthenticatedUser();
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

    render() {
        return (
            <div className="Profile">
                <h5>E-Posta adresi</h5>
                <div style={{paddingBottom:'20px'}}>{this.state.userInfo && this.state.userInfo.email}</div>
                <ChangePassword/>
            </div>
        );
    }
}