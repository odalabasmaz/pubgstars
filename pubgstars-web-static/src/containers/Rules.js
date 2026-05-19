import React, { Component } from "react";
import { Card,Accordion,Container } from "react-bootstrap";
import "./Rules.css";

export default class Rules extends Component {

  render() {
    return (
      <div className="Rules">
        <div className="header-container">
          <h3 className="banner-title">Site kuralları</h3>
        </div>
        <Container>
          <Accordion defaultActiveKey="0">
            <Card>
              <Accordion.Toggle as={Card.Header} eventKey="0">
                DAVRANIŞ KURALLARI
              </Accordion.Toggle>
              <Accordion.Collapse eventKey="0">
                <Card.Body>
                  Oyuncuların birbirleri ile olan iletişimlerinde saygı kuralları gözetilmektedir. Oyuncu ve yönetici arasında da aynı kurallar geçerlidir. Saygı kuralları çerçevesinde olmayan her türlü iletişim (küfür, argo, kişilik haklarını küçük düşürücü davranış, siyasi ve dini görüşlere saygısızlık, politik söylemler vs) yöneticiler tarafından izlenmekte olup gerekli yaptırımlar uygulanacaktır.
                  Yukarda belirtilen kurallar oyun içinde, web sitesi üzerinde ve Discord kanalımızda da geçerlidir.</Card.Body>
              </Accordion.Collapse>
            </Card>
            <Card>
              <Accordion.Toggle as={Card.Header} eventKey="1">
                OYUNCU VE TAKIM İSMİ KURALLARI
              </Accordion.Toggle>
              <Accordion.Collapse eventKey="1">
                <Card.Body>
                  •	Oyuncu ve takım isimleri onay alınmamış klüp ve marka isimlerinden oluşamaz.<br />
                  •	Oyuncu ve takım isimleri herhangi bir oluşuma hakaret niteliğinde olamaz.<br />
                  •	Oyuncu ve takım isimleri siyasi, ırkçı, cinsiyetçi, dinsel ve aşırı milliyetçi olamaz.</Card.Body>
              </Accordion.Collapse>
            </Card>
            <Card>
              <Accordion.Toggle as={Card.Header} eventKey="2">
                YAYIN KURALLARI
              </Accordion.Toggle>
              <Accordion.Collapse eventKey="2">
                <Card.Body>
                  Oyuncular turnuvada canlı yayın yapabilir. Ancak yayın başlığında “PUBGSTARS Turnuvası” ibaresini eklemek ve 10dk gecikme koymak zorundadır.
                </Card.Body>
              </Accordion.Collapse>
            </Card>
            <Card>
              <Accordion.Toggle as={Card.Header} eventKey="3">
                GENEL ZAMANLAMA KURALLARI
              </Accordion.Toggle>
              <Accordion.Collapse eventKey="3">
                <Card.Body>
                  •	Oyuncular turnuva saatinden 30 dakika önce PUBGSTARS Discord kanalında bulunmak zorundadır.<br />
                  •	“Custom Game” oda adı ve şifresi maçtan en geç 30dk önce katıldığınız oyunun DETAY kısmında duyurulacaktır.<br />
                  •	Discorddaki tüm duyuruları takip etmek takım kaptanının sorumluluğundadır.
                </Card.Body>
              </Accordion.Collapse>
            </Card>
            <Card>
              <Accordion.Toggle as={Card.Header} eventKey="4">
                TEKNİK PROBLEMLER
              </Accordion.Toggle>
              <Accordion.Collapse eventKey="4">
                <Card.Body>
                  •	Maça katılım sağlayan oyuncuların bağlantıda yaşadıkları her türlü problem kendi sorumluluklarındadır.<br />
                  •	Bağlantı problemleri yaşanması durumunda maça devam etmek ya da yeniden başlatmak PUBGSTARS yöneticisinin insiyatifindedir.
                </Card.Body>
              </Accordion.Collapse>
            </Card>
            <Card>
              <Accordion.Toggle as={Card.Header} eventKey="5">
                TAKIM KURALLARI
              </Accordion.Toggle>
              <Accordion.Collapse eventKey="5">
                <Card.Body>
                  •	Takım kaptanının maçtan önce ve maç esnasında Discord kanalında bulunması zorunludur.<br />
                  •	Tüm oyuncuların Discord kanalında perm alması  zorunludur. Aksi takdirde maça katılım gösteremezler.
                </Card.Body>
              </Accordion.Collapse>
            </Card>
            <Card>
              <Accordion.Toggle as={Card.Header} eventKey="6">
                OYUN KURALLARI
              </Accordion.Toggle>
              <Accordion.Collapse eventKey="6">
                <Card.Body>
                  •	Lobide bulunan oyuncuların %10undan fazlasının bağlantı problemi yaşaması durumunda,<br />
                  •	Olası herhangi bir hata durumunda, <br />
                  Lobi yeniden başlatılır. Bunlar dışında yaşanabilecek problemlerin çözümü için PUBGSTARS yöneticileri aksiyon alacaktır.
                </Card.Body>
              </Accordion.Collapse>
            </Card>
            <Card>
              <Accordion.Toggle as={Card.Header} eventKey="7">
                HİLE VE BUG
              </Accordion.Toggle>
              <Accordion.Collapse eventKey="7">
                <Card.Body>
                  •	Maçlara katılan oyuncuların 3. Parti yazılımlar kullanması kesinlikle yasaktır. Bunlardan bazıları;<br />
                  ESP, Radar hilesi, Wallhack, Speedhack, Aimhack, Hitbox hilesi, Teleport hilesi<br />
                  •	Maçlara katılan oyuncuların yardımcı yazılımlar kullanması kesinlikle yasaktır. Bunlardan bazıları;<br />
                  ReShade, SweetFX, VibranceGUI<br />
                  •	Maça katılan oyuncuların son 2 yıl içerisinde PUBGSTARS bünyesinde düzenlenen turnuvalardan herhangi birinde diskalifiye edilmemiş olması gerekmektedir. Aksi durum tespit edilirse ilgili oyuncu yeni maça katılamaz.<br />
                  •	Maç esnasında hile yapıldığı tespit edilir ise lobi yeniden başlatılır. Maçtan sonra fark edilmesi durumunda maç yeniden başlatılmaz ve sıralama yeniden oluşturulur.
                </Card.Body>
              </Accordion.Collapse>
            </Card>
            <Card>
              <Accordion.Toggle as={Card.Header} eventKey="8">
                CHECK-IN İLETİŞİM
              </Accordion.Toggle>
              <Accordion.Collapse eventKey="8">
                <Card.Body>
                  •	Duyuruları takip edebilmek için takım kaptanı perm alması  gerekmektedir.<br />
                  •	Duyuruları takip etmek takım kaptanlarının sorumluluğundadır.<br />
                  •	Duyuruların PUBGSTARS Discord kanalından takip edilmesi gerekmektedir.
                </Card.Body>
              </Accordion.Collapse>
            </Card>
          </Accordion>

          <br />
          <br />
          <p>
            Maçlara katılan tüm takımlar yayımlanan bu kuralları kabul etmiş sayılır.
            PUBGSTARS ödül ve kuralları değiştirme hakkını saklı tutar.
          </p>
        </Container>
      </div>
    );
  }
}
