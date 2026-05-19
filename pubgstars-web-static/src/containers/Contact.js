import React, {Component} from "react";
import {Container, Form, Button, FormLabel, FormControl, FormGroup} from "react-bootstrap";
import API from '../utils/api';

import "./Contact.css";
import Popup from "../components/Popup/Popup";

export default class Contact extends Component {

    constructor(props) {
        super(props);
        this.state = {
            from: "",
            message: "",
            modalShow: false
        };
    }

    handleChange = event => {
        this.setState({
            [event.target.id]: event.target.value
        });
    };

    handleSubmit = async event => {
        event.preventDefault();

        API.sendMessage({
            from: this.state.from,
            message: this.state.message
        }).then(res => {
            if (res.data.statusCode !== 200) {
                this.setState({popupText: "Mesajınız iletilmemiştir. Lütfen daha sonra tekrar deneyiniz"});
            } else {
                this.setState({popupText: "Mesajınız iletilmiştir"});
            }
            this.setState({from: "", message: ""});
            this.setState({modalShow: true});
        });
    };

    validateForm() {
        return this.state.from.length > 0 && this.state.message.length > 0;
    }

    render() {
        let modalClose = () => this.setState({ modalShow: false });
        return (
            <div className="Contact">
                <div className="header-container">
                    <h3 className="banner-title">Bize Ulaşın</h3>
                </div>
                <Container>
                    <p>Görüş ve önerileriniz bizim için değerlidir. Aşağıdaki form aracılığıyla bize ulaşıp,
                        düşüncelerinizi paylaşabilirsiniz.</p>
                    <Form onSubmit={this.handleSubmit}>
                        <FormGroup controlId="from">
                            <FormLabel>E-Posta Adresiniz</FormLabel>
                            <FormControl
                                autoFocus
                                type="email"
                                value={this.state.from}
                                onChange={this.handleChange}
                            />
                        </FormGroup>
                        <Form.Group controlId="message">
                            <Form.Label>Mesajınızı Giriniz</Form.Label>
                            <Form.Control as="textarea" rows="3" value={this.state.message}
                                          onChange={this.handleChange}/>
                        </Form.Group>
                        <Button block type="submit" disabled={!this.validateForm()}>Gönder</Button>
                    </Form>
                </Container>
                <Popup show={this.state.modalShow} onHide={modalClose}>{this.state.popupText}</Popup>
            </div>
        );
    }
}
