import React, {Component} from "react";
import {FormGroup, FormControl, FormLabel} from "react-bootstrap";
import {Auth} from "aws-amplify";
import LoaderButton from "../components/LoaderButton";

import "./ChangePassword.css";

export default class ChangePassword extends Component {

    constructor(props) {
        super(props);

        this.state = {
            password: "",
            oldPassword: "",
            isChanging: false,
            confirmPassword: ""
        };
    }

    validateForm() {
        return (
            this.state.oldPassword.length > 0 &&
            this.state.password.length > 0 &&
            this.state.password === this.state.confirmPassword
        );
    }

    handleChange = event => {
        this.setState({
            [event.target.id]: event.target.value
        });
    };

    handleChangeClick = async event => {
        event.preventDefault();

        this.setState({ isChanging: true });

        try {
            const currentUser = await Auth.currentAuthenticatedUser();
            await Auth.changePassword(
                currentUser,
                this.state.oldPassword,
                this.state.password
            );
            this.setState({ isChanging: false, password: "", oldPassword: "", confirmPassword: ""});
            console.log("Şifreniz değiştirilmiştir")
        } catch (e) {
            console.log(e.message)
            this.setState({ isChanging: false });
        }
    };

    render() {
        return (
            <div className="ChangePassword">
                <h5>Şifreniz</h5>
                <form onSubmit={this.handleChangeClick}>
                    <FormGroup controlId="oldPassword">
                        <FormLabel>Eski Şifre</FormLabel>
                        <FormControl
                            type="password"
                            onChange={this.handleChange}
                            value={this.state.oldPassword}
                        />
                    </FormGroup>
                    <hr />
                    <FormGroup controlId="password">
                        <FormLabel>Yeni Şifre</FormLabel>
                        <FormControl
                            type="password"
                            value={this.state.password}
                            onChange={this.handleChange}
                        />
                    </FormGroup>
                    <FormGroup  controlId="confirmPassword">
                        <FormLabel>Yeni Şifre</FormLabel>
                        <FormControl
                            type="password"
                            onChange={this.handleChange}
                            value={this.state.confirmPassword}
                        />
                    </FormGroup>
                    <LoaderButton
                        block
                        disabled={!this.validateForm()}
                        type="submit"
                        isLoading={this.state.isChanging}
                        text="Şifreyi Değiştir"
                        loadingText="Şifre Değiştiriliyor…"
                    />
                </form>
            </div>
        );
    }
}
