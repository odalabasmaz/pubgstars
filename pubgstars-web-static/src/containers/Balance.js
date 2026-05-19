import React, {Component} from "react";
import {Card, Button, Form} from "react-bootstrap";
import Popup from "../components/Popup/Popup";
import MaskedFormControl from 'react-bootstrap-maskedinput'

import {Auth} from "aws-amplify";
import API from '../utils/api'

import "./Balance.css";

export default class Balance extends Component {

    constructor(props) {
        super(props);

        this.state = {
            modalShow: false,
            validated: false,
            secretQuestion: "Annenizin kızlık soyadı?",
            secretAnswer: ""
        }
    }

    async componentDidMount() {
        Auth.currentAuthenticatedUser().then((user) => {
            const options = {
                headers: {
                    Authorization: user.signInUserSession.idToken.jwtToken
                }
            };

            API.getUser(options)
                .then(res => {
                    this.setState({userInfo: res.data.body});
                });
        });
    }

    handleChange = event => {
        this.setState({
            [event.target.id]: event.target.value
        });
    };

    handleWithdraw = async event => {
        event.preventDefault();
        if (this.state.iban.replace(/_/g, "").replace(/ /g, "").length !== 26) {
            this.setState({validated: false});
            return;
        }

        if (this.state.amountOut <= 0) {
            this.setState({validated: false});
            return;
        }

        const parameters = {
            iban: this.state.iban.replace(/ /g, ""),
            amount: this.state.amountOut,
            nameSurname: this.state.nameSurname,
            secretQuestion: this.state.secretQuestion,
            secretAnswer: this.state.secretAnswer
        };

        try {
            Auth.currentAuthenticatedUser().then((user) => {
                const options = {
                    headers: {
                        Authorization: user.signInUserSession.idToken.jwtToken
                    }
                };
                API.getWithdraw(parameters, options)
                    .then((res) => {
                        if (res.data.statusCode !== 200) {
                            this.setState({popupText: res.data.errorMessage});
                        } else {
                            this.setState({popupText: parameters.amount + "₺ para çekme talebiniz alınmıştır!"});
                        }
                        this.setState({modalShow: true, iban: '', amountOut: '', nameSurname: '', secretAnswer: ''});
                    });
            });
        } catch (e) {
            console.log(e);
        }
    };

    handleDeposit = async event => {
        event.preventDefault();

        if (this.state.amountIn <= 0) {
            this.setState({validated: false});
            return;
        }

        const parameters = {
            amount: this.state.amountIn,
            description: this.state.userInfo.username
        };

        try {
            Auth.currentAuthenticatedUser().then((user) => {
                const options = {
                    headers: {
                        Authorization: user.signInUserSession.idToken.jwtToken
                    }
                };
                API.getDeposit(parameters, options)
                    .then((res) => {
                        if (res.data.statusCode !== 200) {
                            this.setState({popupText: res.data.errorMessage});
                        } else {
                            this.setState({popupText: "Para yatırma talebiniz alınmıştır, kontrol edildikten sonra hesabınıza yansıtılacaktır!"});
                        }
                        this.setState({modalShow: true, amountIn: ''});
                    });
            });
        } catch (e) {
            console.log(e);
        }
    };

    render() {
        let modalClose = () => this.setState({modalShow: false});
        return (
            <div className="Balance">
                <Card style={{marginBottom: '15px'}}>
                    <Card.Header>
                        <p><span>Oynanabilir tutar: {this.state.userInfo && (this.state.userInfo.balance + this.state.userInfo.bonus)}₺</span></p>
                        <p><span>Bakiye: {this.state.userInfo && this.state.userInfo.balance}₺ + Bonus: {this.state.userInfo && this.state.userInfo.bonus}₺</span></p>
                    </Card.Header>
                    <Card.Body>
                        <h4>Hesabıma Para Yükle</h4>
                        <div style={{display: 'flex'}}>
                            <i className="fas fa-arrow-down fa-4x" style={{color: '#32cd32'}}/>
                            <div style={{marginLeft: 25}}>
                                <p>
                                    Bakiye yüklemesi yapmak için aşağıda belirtilen bilgilerle EFT/Havale yapabilirsiniz.<br/>
                                    Ödediğiniz tutar en geç <b>2 iş günü</b> içinde sisteme yansıtılacaktır.<br/>
                                </p>
                                <p>
                                    Banka: <b>Türkiye İş Bankası</b><br/>
                                    IBAN: <b>TR280006400000143780023495</b><br/>
                                    Ad Soyad: <b>Ümit Utku Adak</b><br/>
                                    Açıklama: <b>{this.state.userInfo && this.state.userInfo.username}</b><br/>
                                </p>
                            </div>
                        </div>
                        <div style={{paddingBottom: '30px'}}>
                            <Form id="asc" onSubmit={this.handleDeposit}>
                                <Form.Label>EFT/Havale Miktarı:</Form.Label>
                                <Form.Group>
                                    <Form.Control id="amountIn"
                                                  placeholder="Yükleme Miktarı"
                                                  onChange={this.handleChange} type="number" step="1" min="1"
                                                  pattern="[0-9]*" required/>
                                </Form.Group>
                                <Button variant="primary" type="submit" block>EFT/Havale Onay Talebi Oluştur</Button>
                            </Form>
                        </div>

                        <h4>Bakiyemi Banka Hesabima Aktar</h4>
                        <div style={{display: 'flex'}}>
                            <i className="fas fa-arrow-up fa-4x" style={{color: '#dc143c'}}/>
                            <div style={{marginLeft: 25}}>
                                <p>Kazancınızi çekebilmeniz için minimum 100₺ bakiyeniz olmasi gerekmektedir</p>
                                <p>Sadece bakiyenizdeki miktari çekebilirsiniz!</p>
                            </div>
                        </div>
                        <div style={{paddingBottom: '30px'}}>
                            <Form id="asd" onSubmit={this.handleWithdraw}>
                                <Form.Group>
                                    <Form.Label>Ad Soyad:</Form.Label>
                                    <Form.Control id="nameSurname"
                                                  type="text"
                                                  maxLength="60"
                                                  value={this.state.nameSurname}
                                                  onChange={this.handleChange}
                                                  required/>
                                </Form.Group>
                                <Form.Group>
                                    <Form.Label>Iban:</Form.Label>
                                    <MaskedFormControl id="iban" type='text'
                                                       mask='TR11 1111 1111 1111 1111 1111 11'
                                                       placeholderChar="_"
                                                       value={this.state.iban}
                                                       onChange={this.handleChange}
                                                       required/>
                                </Form.Group>
                                <Form.Group controlId="amountOut">
                                    <Form.Label>Miktar:</Form.Label>
                                    <Form.Control
                                                  placeholder="Aktarmak istediğiniz miktarı giriniz"
                                                  value={this.state.amountOut}
                                                  onChange={this.handleChange} type="number" step="1" min="100"
                                                  pattern="[0-9]*"
                                                  required/>
                                </Form.Group>
                                <Form.Group controlId="secretQuestion">
                                    <Form.Label>Gizli Soru:</Form.Label>
                                    <Form.Control
                                                  onChange={this.handleChange}
                                                  reqired
                                                  as="select">
                                        <option key={1} value="Annenizin kızlık soyadı?">Annenizin kızlık soyadı?</option>
                                        <option key={2} value="Doğum yerin?">Doğum yerin?</option>
                                        <option key={3} value="En iyi arkadaşın?">En iyi arkadaşın?</option>
                                        <option key={4} value="Tuttuğun takım?">Tuttuğun takım?</option>
                                    </Form.Control>
                                </Form.Group>
                                <Form.Group controlId="secretAnswer">
                                    <Form.Label>Gizli Cevap:</Form.Label>
                                    <Form.Control
                                                  placeholder="Gizli cevabınızı giriniz"
                                                  onChange={this.handleChange} type="text"
                                                  required/>
                                </Form.Group>
                                <Button variant="primary" type="submit" block>Banka Hesabıma Aktar</Button>
                            </Form>
                        </div>
                    </Card.Body>
                </Card>
                <Popup show={this.state.modalShow} onHide={modalClose}>{this.state.popupText}</Popup>
            </div>
        );
    }
}
