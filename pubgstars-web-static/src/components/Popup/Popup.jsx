import React from "react";
import {Modal, Button} from "react-bootstrap";

import "./Popup.css";

export default class Popup extends React.Component {

    onClose = e => {
        this.props.onHide && this.props.onHide(e);
    };

    render() {
        return (
            <Modal
                {...this.props}
                size="md"
                aria-labelledby="popup"
                centered
            >
                <Modal.Header closeButton>
                    <Modal.Title id="popup">
                    </Modal.Title>
                </Modal.Header>
                <Modal.Body style={{'maxHeight': 'calc(100vh - 410px)', 'overflowY': 'auto'}}>
                    <p>
                        {this.props.children}
                    </p>
                </Modal.Body>
                <Modal.Footer>
                    <Button onClick={this.onClose}>Kapat</Button>
                </Modal.Footer>
            </Modal>
        );
    }
}
