import React, {Component} from "react";
import {Card, Table} from "react-bootstrap";

import {Auth} from "aws-amplify";
import API from '../utils/api'

import "./TransactionLog.css";

export default class TransactionLog extends Component {
    constructor(props) {
        super(props);
        this.state = {
            logs: []
        };
    }

    async componentDidMount() {
        Auth.currentAuthenticatedUser().then((user) => {
            const options = {
                headers: {
                    Authorization: user.signInUserSession.idToken.jwtToken
                }
            };
            API.getTransactionLog(options)
                .then(res => {
                    let logs = res.data.body;
                    if (logs === null)
                        logs = [];
                    this.setState({logs});
                });
        });
    }

    render() {
        return (
            <div className="TransactionLog">
                <h4>İşlem Geçmişi!</h4>
                <Table responsive>
                    <thead>
                        <tr>
                            <th>Tarih</th>
                            <th>Tür</th>
                            <th>Alt Tür</th>
                            <th>Detay</th>
                        </tr>
                    </thead>
                    <tbody>
                        {
                        this.state.logs.map((log, i) =>
                            <tr key={i}>
                                <td>{log.dateTime}</td>
                                <td>{log.transactionType}</td>
                                <td>{log.subTransactionType}</td>
                                <td>{log.detail}</td>
                            </tr>
                        )}
                    </tbody>
                </Table>
            </div>
        );
    }
}
