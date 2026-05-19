import React, {Component} from "react";
import {Container} from "react-bootstrap";

import "./About.css";

export default class About extends Component {

    render() {
        return (
            <div className="About">
                <div className="header-container">
                    <h3 className="banner-title">Hakkımızda</h3>
                </div>
                <Container>
                    <h4 className="aboutTitle">PUBGSTARS NEDİR?</h4>
                    <p> PUBGSTARS, bir PUBG Corp ürünü olan ve dünya çapında büyük bir üne sahip PLAYERUNKNOWN’S
                        BATTLEGROUNDS oyunu için ödüllü turnuvalar düzenleyen bir organizasyondur.
                    </p>
                    <br/>
                    <h4 className="aboutTitle">NEDEN PUBGSTARS’A Katılmalıyım ?</h4>
                    <ul>
                        <li><i className="fas fa-angle-right"/>Her maçta hayatta kalarak ilk 3 e giren takımlar büyük
                            ödülden pay sahibi olacak.
                        </li>
                        <li><i className="fas fa-angle-right"/>Oyuncular maç başına ödül kazanacak.</li>
                        <li><i className="fas fa-angle-right"/>Sadece son 3 takım ödül ödemelerine hak kazanacak.</li>
                        <li><i className="fas fa-angle-right"/>Kısacası takımlar oyun oynarken para kazanacak.</li>
                    </ul>
                </Container>
            </div>
        );
    }
}
