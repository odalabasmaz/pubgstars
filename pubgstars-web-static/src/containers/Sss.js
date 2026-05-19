import React, { Component } from "react";
import { Card,Accordion,Container } from "react-bootstrap";
import "./Sss.css";

export default class Sss extends Component {
  constructor(props) {
    super(props);

    this.state = {
      faqList: [
        {key:"Maçlara katılmak için bir yaş sınırı  var mı?", value:"Oyuncuların maçlara katılmak için en az 15 yaşını doldurmuş olmaları gerekmektedir."},
        {key:"Maçlar Solo mu Duo mu Squad olarak mı yapılacak?", value:"Maçlar format olarak Duo ve Squad olarak yapılacaktır."},
        {key:"Maçlarda oynamak için bir yerde mi buluşulacak yoksa herkes kendi evinden mi katılacak?", value:"Maçlar tamamen Online ortamda düzenlenen bir organizasyondur."},
        {key:"Maçlara katılım için Discord sunucusuna katılıp şifre almak yeterli mi?", value:"Takım kaptanlarının aynı zamanda maç esnasında Discord kanalında bulunması gerekmektedir."},
        {key:"Maçlara başvuran takımdaki tüm oyuncular bireysel olarak başvuru yapmalı mı?", value:"Maçlara başvuruları sadece takım kaptanları yapmalıdır."},
        {key:"Maçlara katılım ücretleri bireysel olarak mı yapılmaktadır?", value:"Maçlara katılım ücretlerini toplu şekilde takım kaptanlarının bakiyesiden alınmaktadır."},
        {key:"Belirlenen maç günlerinde katılıma engel bir durum olması durumunda ne yapılıyor?", value:" Maç saatleri organizasyon tarafından belirlenir ve değiştirilemez."},
        {key:"Maçlarda yedek oyuncu bulundurma zorunluluğu var mı?", value:"Böyle bir zorunluluk bulunmamaktadır ama yedek oyuncu bulundurmanız takım yararına olmaktadır."},
        {key:"Maçlar TPP mi FPP mi yapılacak?", value:"Maçlar TPP yapılacaktır."},
        {key:"Resmi Discord kanalınızın bilgileri nedir?", value:"Resmi Discord kanalımızın adresi https://discord.gg/J4Usfu"},
        {key:"Hesaba bakiye nasıl yükleriz?", value:"Hesabınızın altında bulunan ödeme yöntemlerini kullanarak dilediğiniz miktarda hesabınıza bakiye yükleyebilirsiniz."},
        {key:"Kazanılan ödüller hesaba ne zaman yüklenir?", value:"Kazanılan ödül başvuru yapan takım kaptanının hesabına en geç 3 saat içerisinde yüklenir."},
        {key:"Hesaptaki bakiyeyi farklı insanlara gönderebilir miyiz?", value:"Hesabınızda bulunan bakiyeyi sadece kayıt olurken vermiş olduğunuz kişisel bilgiler (ad-soyad) ile aynı bilgilere sahip olan başka bir hesaba gönderebilirsiniz."},
        {key:"Bakiyemizdeki para banka hesabımıza ne zaman geçer?", value:"Bakiyenizdeki para en geç 3 iş günü içerisinde hesabınızda olur."},
        {key:"Bakiyemizdeki parayı çekme limiti var mıdır?", value:"Bakiyenizdeki parayı çekme limiti şuan için 100TL kadardır."},
        {key:"Maçlar ne zaman ve ne şekilde açıklanıyor?", value:"Maç saatinden 1 saat önce custom game oluşturulmaktadır ve bilgiler site üzerinden maç DETAY kısmından görüntülenebilmektedir."},
        {key:"Maçtan ne kadar önce takım hazır olmalıdır?", value:"Maç tam saatinde başlatılmaktadır. Bu nedenle takım olarak bunu göz önünde bulundurup öncesinden hazır olmanız gerekmektedir."},
        {key:"Maç saati geldi ve takımda eksik oyuncu var maç başlatılır mı?", value:"Belirlenen maç saati geldiğinde eğer takımınızda eksik oyuncu var ise maç başlatılır. Maçtan önce size tanınan sürede tüm hazırlıklarınızı yapmanız gerekmektedir."},
        {key:"Günlük maçlara katılmak için bir limit var mı?", value:"Hayır."},
        {key:"Maçlarda hile yapıldığı tespit edilirse ne olur?", value:"Hile yapıldığı tespit edilen maçlar Remake ile tekrardan başlatılır."},
        {key:"Maçlarda hile yapan oyuncular için herhangi bir yaptırım var mı?", value:"Hile yaptığı belirlenen oyuncular takımları ile birlikte diskalifiye edilir. Diskalifiye edilen takımın maç için yatırdığı ücret maçın toplam tutarına dâhil edilir ve kazananlar arasında paylaşılır. Diskalifiye edilen takımlar hiçbir koşulda PUBGSTARS organizasyonlarına dâhil olamazlar."},
        {key:"Maçlarda video kaydı alınmakta mıdır?", value:"Maçlarda video kaydı alınmakta ve günlük olarak youtube kanalımızda paylaşılmaktadır."},
        {key:"Maçlara arkadaşlarımız spectator (gözlemci) olarak girebilir mi?", value:" Maçların seyrini etkileyebileceği (ghostlanmak vs gibi) için oyunlara spectator(gözlemci) kabul edilmemektedir."},
        {key:"Maçlara katılım için nasıl başvuru yapılır?", value:"Web sitemiz üzerinden Turnuvalar kısmındaki maçlardan size uygun olan için katıl butonuna basmanız yeterlidir."},
        {key:"Başvuru yapılan maçlar herhangi bir sebepten ötürü iptal edilebilir mi?", value:"Maça katıldığınız gibi iptal de edebilirsiniz. Hesabınızın altında bulunan maçlarım bölümüne girerek ilgili maçı iptal edebilirsiniz. Maç iptal işlemleri sadece maçlardan 1 saat önce yapılabilmektedir."}
      ]
    };
  }

  render() {
    return (
      <div className="Sss">
        <div className="header-container">
          <h3 className="banner-title">Sık Sorulan Sorular</h3>
        </div>
        <Container>
          <Accordion defaultActiveKey="0">
            {
              this.state.faqList.map(
                  (faq, i) =>
                      <Card key={i}>
                        <Accordion.Toggle as={Card.Header} eventKey={i}>{faq.key}</Accordion.Toggle>
                        <Accordion.Collapse eventKey={i}>
                          <Card.Body>{faq.value}</Card.Body>
                        </Accordion.Collapse>
                      </Card>
              )
            }
          </Accordion>
        </Container>
      </div>
    );
  }
}
