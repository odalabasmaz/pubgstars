import React from "react";

import First from "../../images/first.jpg";
import Second from "../../images/second.jpg";
import Third from "../../images/third.jpg";
import Moment from "react-moment";
import {ProgressBar} from "react-bootstrap";

import "./GameInfo.css";

const GameInfo = (props) => (
    <div style={{paddingBottom:'15px'}}>
        <div>Tarih: <b><Moment parse="YYYYMMDDHHmm" format="DD/MM/YYYY HH:mm" date={props.game.gameDate}/></b></div>
        <div>Platform: <b>{props.game.platform}</b></div>
        <div>Oyun Modu: <b>{props.game.type}</b></div>
        <div>Takım Oyuncu Sayısı: <b>{props.game.teamPlayerCount}</b></div>
        <div>Takım Katılım Ücreti: <b>{props.game.price} ₺</b></div>
        <div>Kayıtlı Oyuncu Sayısı: <b>{props.game.registeredUserCount}/100</b></div>
        <ProgressBar now={props.game.registeredUserCount}/>
        <div style={{display: 'flex', flexDirection: 'row', justifyContent: 'space-around', paddingTop: '10px'}}>
            <div>
                <img src={Second} alt="second" style={{width: 36}}/>
                <div style={{display: 'flex', justifyContent: 'center'}}><b>{props.game.award2nd} ₺</b></div>
            </div>
            <div>
                <img src={First} alt="first" style={{width: 36}}/>
                <div style={{display: 'flex', justifyContent: 'center'}}><b>{props.game.award1st} ₺</b></div>
            </div>
            <div>
                <img src={Third} alt="third" style={{width: 36}}/>
                <div style={{display: 'flex', justifyContent: 'center'}}><b>{props.game.award3rd} ₺</b></div>
            </div>
        </div>
    </div>
);

export default GameInfo;
