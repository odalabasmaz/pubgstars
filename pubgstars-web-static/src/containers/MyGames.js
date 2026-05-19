import React, {Component} from "react";

import {Auth} from "aws-amplify";
import API from "../utils/api";
import {Button, Card, Modal} from "react-bootstrap";
import GameInfo from "../components/GameInfo/GameInfo";

import "./MyGames.css";

export default class MyGames extends Component {

    constructor(props) {
        super(props);

        this.handleShow = this.handleShow.bind(this);
        this.handleClose = this.handleClose.bind(this);
        this.handleClose2 = this.handleClose2.bind(this);

        this.state = {
            isLoading: true,
            games: [],
            gameUsers: []
        };
    }

    async componentDidMount() {
        this.getUserGameList();
        this.setState({isLoading: false});
    }

    getUserGameList() {
        Auth.currentAuthenticatedUser().then((user) => {
            const options = {
                headers: {
                    Authorization: user.signInUserSession.idToken.jwtToken
                }
            };
            API.getUserGames(options)
                .then(res => {
                    const games = res.data.body.filter(function (game) {
                        return game.registered;
                    });

                    if (games)
                        this.setState({games});
                });
        });
    }

    showGamePassword(id, e) {
        const game = {
            id: id
        };
        Auth.currentAuthenticatedUser().then((user) => {
            const options = {
                headers: {
                    Authorization: user.signInUserSession.idToken.jwtToken
                }
            };
            API.getGamePassword(game, options)
                .then(res => {
                    this.setState({showPassword: true, cancellable: false, passwordValue: res.data.body.roomPassword, discordValue: res.data.body.discord});
                });
            API.getGameUsers(game, options)
                .then(res => {
                    this.setState({gameUsers: res.data.body});
                });
        });
    }

    unregisterFromGame(id, e) {
        const game = {
            id: id
        };
        Auth.currentAuthenticatedUser().then((user) => {
            const options = {
                headers: {
                    Authorization: user.signInUserSession.idToken.jwtToken
                }
            };
            API.unregisterGame(game, options)
                .then(res => {
                    if (res.data.statusCode !== 200) {
                        alert(res.data.errorMessage);
                    } else {
                        this.getUserGameList();
                    }
                });
        });
        this.setState({show: false});
    }

    handleClose() {
        this.setState({show: false});
    }

    handleClose2() {
        this.setState({showPassword: false, cancellable: true});
    }

    handleShow(selGame, e) {
        this.setState({show: true});
        this.setState({selectedGame: selGame});

        this.setState({selectedId: selGame.id});
        this.setState({selectedRegistered: selGame.registered});
    }

    renderUserGamesList(games) {
        if (games.length === 0) {
            return (
                <h4>Kayıtlı olduğunuz oyun bulunmuyor.</h4>
            );
        }
        return games.map(
            (game, i) =>
                <Card key={i} style={{marginBottom: "20px", fontSize: '0.9rem', width: '250px'}}>
                    <Card.Img variant="top" src={require('../images/gameImages/map_' + game.map.toLowerCase() + '1.jpg')}/>
                    <Card.Body>
                        <Card.Title>{game.map}</Card.Title>
                        <GameInfo game={game}/>
                        {
                            game.showPassword
                                ?
                                <Button style={{borderRadius: "0", width: "100%", backgroundColor: "orange"}}
                                        variant="primary"
                                        onClick={this.showGamePassword.bind(this, game.id)}>Detay</Button>
                                :
                                game.cancellable
                                    ?
                                    <Button style={{borderRadius: "0", width: "100%", backgroundColor: "red"}}
                                            variant="primary" onClick={this.handleShow.bind(this, game)}>Oyundan Çık</Button>
                                    :
                                    <Button style={{borderRadius: "0", width: "100%", backgroundColor: "gray"}}
                                            variant="primary" disabled>Geçmis oyun</Button>
                        }
                    </Card.Body>
                </Card>
        );
    }

    render() {
        return (
            <div className="MyGames">
                <div>
                    {!this.state.isLoading && this.renderUserGamesList(this.state.games)}
                    <Modal show={this.state.show} onHide={this.handleClose}>
                        <Modal.Header closeButton/>
                        <Modal.Body>
                            <p>Oyundan çıkmak istediğinize emin misiniz?</p>
                        </Modal.Body>
                        <Modal.Footer>
                            <Button variant="secondary" onClick={this.handleClose}>Kapat</Button>
                            <Button variant="primary" onClick={this.unregisterFromGame.bind(this, this.state.selectedId)}>Oyundan Çık</Button>
                        </Modal.Footer>
                    </Modal>
                    <Modal show={this.state.showPassword} onHide={this.handleClose2}>
                        <Modal.Header closeButton/>
                        <Modal.Body>
                            <div>Oda adı ve şifresi: {this.state.passwordValue}</div>
                            <div>Discord: {this.state.discordValue}</div>
                            <div style={{padding: '10px 0'}}>Kayıtlı Oyuncular</div>
                            <ul>
                                {this.state.gameUsers.map((value, index) => {
                                    return <li key={index}>{value.username}</li>
                                })}
                            </ul>
                        </Modal.Body>
                        <Modal.Footer>
                            <Button variant="secondary" onClick={this.handleClose2}>Kapat</Button>
                        </Modal.Footer>
                    </Modal>
                </div>
            </div>
        );
    }
}
