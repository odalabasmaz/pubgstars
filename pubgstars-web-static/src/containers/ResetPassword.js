import React, {Component} from "react";
import {Auth} from "aws-amplify";
import {Link} from "react-router-dom";
import {FormText, FormGroup, FormControl, FormLabel} from "react-bootstrap";
import LoaderButton from "../components/LoaderButton";
import "./ResetPassword.css";

export default class ResetPassword extends Component {
    constructor(props) {
        super(props);

        this.state = {
            code: "",
            email: "",
            password: "",
            codeSent: false,
            confirmed: false,
            confirmPassword: "",
            isConfirming: false,
            isSendingCode: false
        };
    }

    validateCodeForm() {
        return this.state.email.length > 0;
    }

    validateResetForm() {
        return (
            this.state.code.length > 0 &&
            this.state.password.length > 0 &&
            this.state.password === this.state.confirmPassword
        );
    }

    handleChange = event => {
        this.setState({
            [event.target.id]: event.target.value
        });
    };

    handleSendCodeClick = async event => {
        event.preventDefault();

        this.setState({isSendingCode: true});

        try {
            await Auth.forgotPassword(this.state.email);
            this.setState({codeSent: true});
        } catch (e) {
            this.setState({error: true, errorText: "Girmiş olduğunuz e-posta sistemimizde kayıtlı değildir. Lütfen e-posta adresinizi kontrol ediniz!"});
            this.setState({isSendingCode: false});
        }
    };

    handleConfirmClick = async event => {
        event.preventDefault();

        this.setState({isConfirming: true});

        try {
            await Auth.forgotPasswordSubmit(
                this.state.email,
                this.state.code,
                this.state.password
            );
            this.setState({confirmed: true});
        } catch (e) {
            if (e.code === "InvalidPasswordException")
                this.setState({error: true, errorText: "Belirlediğiniz şifre en az 1 büyük harf, özel karakter içermelidir!"});
            else if (e.code === "UsernameExistsException")
                this.setState({error: true, errorText: "Girilen E-posta adresi sistemde kayıtlı!"});
            else
                this.setState({error: true, errorText: e.message});
            this.setState({isConfirming: false});
        }
    };

    renderRequestCodeForm() {
        return (
            <>
                <form onSubmit={this.handleSendCodeClick}>
                    <FormText>Şifrenizi yeniden oluşturmak için <b>E-Posta adresinizi</b> girin.</FormText>
                    <FormGroup bsSize="large" controlId="email">
                        <FormLabel>E-Posta</FormLabel>
                        <FormControl
                            autoFocus
                            type="email"
                            value={this.state.email}
                            onChange={this.handleChange}
                        />
                    </FormGroup>
                    <LoaderButton
                        block
                        type="submit"
                        bsSize="large"
                        loadingText="Gönderiliyor…"
                        text="Gönder"
                        isLoading={this.state.isSendingCode}
                        disabled={!this.validateCodeForm()}
                    />
                </form>
                {this.state.error && <div style={{color:'red', fontSize: '13px', marginBottom:'10px'}}>{this.state.errorText}</div>}
            </>
        );
    }

    renderConfirmationForm() {
        return (
            <>
                <form onSubmit={this.handleConfirmClick}>
                    <FormGroup controlId="code">
                        <FormLabel>Doğrulama Kodu</FormLabel>
                        <FormControl
                            autoFocus
                            type="tel"
                            value={this.state.code}
                            onChange={this.handleChange}
                        />
                        <FormText>
                            Doğrulama kodu mail adresinize ({this.state.email}) gönderilmiştir.
                        </FormText>
                    </FormGroup>
                    <hr/>
                    <FormGroup controlId="password">
                        <FormLabel>Yeni Şifre</FormLabel>
                        <FormControl
                            type="password"
                            value={this.state.password}
                            onChange={this.handleChange}
                        />
                    </FormGroup>
                    <FormGroup controlId="confirmPassword">
                        <FormLabel>Yeni Şifre Tekrar</FormLabel>
                        <FormControl
                            type="password"
                            onChange={this.handleChange}
                            value={this.state.confirmPassword}
                        />
                    </FormGroup>
                    <LoaderButton
                        block
                        type="submit"
                        text="Doğrula"
                        loadingText="Doğrulanıyor…"
                        isLoading={this.state.isConfirming}
                        disabled={!this.validateResetForm()}
                    />
                </form>
                {this.state.error && <div style={{color:'red', fontSize: '13px', marginBottom:'10px'}}>{this.state.errorText}</div>}
            </>
        );
    }

    renderSuccessMessage() {
        return (
            <div className="success">
                <p>Şifreniz sıfırlanmıştır.</p>
                <p>
                    <Link to="/login">
                        Yeni şifrenizle giriş yapmak için buraya tıkayınız.
                    </Link>
                </p>
            </div>
        );
    }

    render() {
        return (
            <div className="ResetPassword">
                {!this.state.codeSent
                    ? this.renderRequestCodeForm()
                    : !this.state.confirmed
                        ? this.renderConfirmationForm()
                        : this.renderSuccessMessage()}
            </div>
        );
    }
}