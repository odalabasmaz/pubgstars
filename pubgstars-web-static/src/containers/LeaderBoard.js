import React, {Component} from "react";
import Table from 'react-bootstrap/Table';
import {Auth} from "aws-amplify";
import Moment from 'react-moment';

import API from '../utils/api'

import "./LeaderBoard.css";

export default class LeaderBoard extends Component {
    constructor(props) {
        super(props);

        this.state = {
            leaderBoardList: [],
            gamesHistory: []
        };
    }

    async componentDidMount() {
        Auth.currentAuthenticatedUser().then((user) => {
            const options = {
                headers: {
                    Authorization: user.signInUserSession.idToken.jwtToken
                }
            };

            API.getLeaderboard(options)
                .then(res => {
                    let data = res.data.body;
                    if (data === null) data = [];
                    this.setState({leaderBoardList: data});
                });

            API.getGamesHistory(options)
                .then(res => {
                    let data = res.data.body;
                    if (data === null) data = [];
                    this.setState({gamesHistory: data});
                });
        });
    }

    render() {
        return (
            <div className="LeaderBoard">
                <h4>En çok kazananlar!</h4>
                <Table responsive>
                    <thead>
                        <tr>
                            <th>Kullanıcı</th>
                            <th>Toplam Kazanc</th>
                        </tr>
                    </thead>
                    <tbody>
                        {this.state.leaderBoardList.map((user, i) =>
                            <tr key={i}>
                                <td><i className="fa fa-user"/>&nbsp;{user.username}</td>
                                <td>{user.gain}₺</td>
                            </tr>
                        )}
                    </tbody>
                </Table>

                <h4>Oyunların galipleri!</h4>
                <Table responsive>
                    <thead>
                        <tr>
                            <th>Tarih</th>
                            <th>Harita</th>
                            <th>Platform</th>
                            <th>Mod</th>
                            <th>#1</th>
                            <th><i className="fa fa-trophy"/></th>
                            <th>#2</th>
                            <th><i className="fa fa-trophy"/></th>
                            <th>#3</th>
                            <th><i className="fa fa-trophy"/></th>
                        </tr>
                    </thead>
                    <tbody>
                        {this.state.gamesHistory.map((game, i) =>
                            <tr key={i}>
                                <td>
                                    <i className="fa fa-calendar-alt"/>&nbsp;
                                    <Moment parse="YYYYMMDDHHmm" format="DD/MM/YYYY HH:mm" date={game.gameDate}/>
                                </td>
                                <td>{game.map}</td>
                                <td>{game.platform}</td>
                                <td>{game.type}</td>
                                <td>{game.winner1stName}</td>
                                <td>{game.award1st}₺</td>
                                <td>{game.winner2ndName}</td>
                                <td>{game.award2nd}₺</td>
                                <td>{game.winner3rdName}</td>
                                <td>{game.award3rd}₺</td>
                            </tr>
                        )}
                    </tbody>
                </Table>
            </div>
        );
    }
}
