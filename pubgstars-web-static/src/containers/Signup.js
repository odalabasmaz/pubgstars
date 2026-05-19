import React, {Component} from "react";
import Form from 'react-bootstrap/Form';
import {Link} from "react-router-dom";
import {Auth} from "aws-amplify";
import LoaderButton from "../components/LoaderButton";
import Popup from "../components/Popup/Popup";

import "./Signup.css";

export default class Signup extends Component {
    constructor(props) {
        super(props);

        this.state = {
            isLoading: false,
            email: "",
            password: "",
            username: "",
            confirmPassword: "",
            confirmationCode: "",
            secretQuestion: "Annenizin kızlık soyadı?",
            secretAnswer: ""
        };
    }

    validateForm() {
        return (
            this.state.email.length > 0 &&
            this.state.password.length > 0 &&
            this.state.username.length > 0 &&
            this.state.secretAnswer.length > 0 &&
            this.state.password === this.state.confirmPassword
        );
    }

    validateConfirmationForm() {
        return this.state.confirmationCode.length > 0;
    }

    handleChange = event => {
        this.setState({
            [event.target.id]: event.target.value
        });
    };

    handleSubmit = async event => {
        event.preventDefault();

        this.setState({isLoading: true});

        try {
            const newUser = await Auth.signUp({
                username: this.state.email,
                password: this.state.password,
                attributes: {
                    "name": this.state.username,
                    "custom:secretQuestion": this.state.secretQuestion,
                    "custom:secretAnswer": this.state.secretAnswer
                }
            });
            this.setState({newUser: newUser, isLoading: false});
        } catch (e) {
            if (e.code === "InvalidPasswordException")
                this.setState({error: true, errorText: "Belirlediğiniz şifre en az 1 büyük harf, özel karakter içermelidir!"});
            else if (e.code === "UsernameExistsException")
                this.setState({error: true, errorText: "Girilen E-posta adresi sistemde kayıtlı!"});
            else
                this.setState({error: true, errorText: e.message});
            this.setState({isLoading: false});
        }
    };

    handleConfirmationSubmit = async event => {
        event.preventDefault();

        this.setState({isLoading: true});

        try {
            await Auth.confirmSignUp(this.state.email, this.state.confirmationCode);
            await Auth.signIn(this.state.email, this.state.password);

            this.props.userHasAuthenticated(true);
            this.props.history.push("/");
        } catch (e) {
            this.setState({isLoading: false, errorText: e.message});
        }
    };

  renderCheckboxLabel() {
    return (
        <div>
          <a href={"javascript:void(0)"} className={'thing'}
             onClick={() => this.setState({showTermsOfUseContract: true})}>Üyelik koşullarını</a> ve
          <a href={"javascript:void(0)"} className={'thing'}
             onClick={() => this.setState({showProtectionOfPersonalDataContract: true})}> kişisel verilerimin kullanılmasını </a>kabul ediyorum
        </div>
    );
  }

    renderConfirmationForm() {
        return (
            <div className="SignupForm">
                <Form onSubmit={this.handleConfirmationSubmit}>
                    <Form.Group controlId="confirmationCode">
                        <Form.Label>Şifre Doğrulama</Form.Label>
                        <Form.Control
                            autoFocus
                            type="tel"
                            value={this.state.confirmationCode}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <LoaderButton
                        block
                        disabled={!this.validateConfirmationForm()}
                        type="submit"
                        isLoading={this.state.isLoading}
                        text="Doğrula"
                        loadingText="Doğrulanıyor…"
                    />
                </Form>
            </div>
        );
    }

    renderForm() {
        return (
            <div className="SignupForm">
                <div className="form-title">Üye Ol</div>
                <Form onSubmit={this.handleSubmit}>
                    <Form.Group controlId="email">
                        <Form.Label>E-Posta Adresiniz</Form.Label>
                        <Form.Control
                            autoFocus
                            type="email"
                            value={this.state.email}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Form.Group controlId="username">
                        <Form.Label>Kullanıcı Adınız</Form.Label>
                        <Form.Control
                            type="input"
                            value={this.state.username}
                            onChange={this.handleChange}
                        />
                    </Form.Group>
                    <Form.Group controlId="password">
                        <Form.Label>Şifre</Form.Label>
                        <Form.Control
                            value={this.state.password}
                            onChange={this.handleChange}
                            type="password"
                        />
                    </Form.Group>
                    <Form.Group controlId="confirmPassword">
                        <Form.Label>Şifre Tekrar</Form.Label>
                        <Form.Control
                            value={this.state.confirmPassword}
                            onChange={this.handleChange}
                            type="password"
                        />
                    </Form.Group>
                    <Form.Group controlId="secretQuestion">
                        <Form.Label>Gizli Soru</Form.Label>
                        <Form.Control
                            onChange={this.handleChange}
                            as="select">
                            <option key={1} value="Annenizin kızlık soyadı?">Annenizin kızlık soyadı?</option>
                            <option key={2} value="Doğum yerin?">Doğum yerin?</option>
                            <option key={3} value="En iyi arkadaşın?">En iyi arkadaşın?</option>
                            <option key={4} value="Tuttuğun takım?">Tuttuğun takım?</option>
                        </Form.Control>
                    </Form.Group>
                    <Form.Group controlId="secretAnswer">
                        <Form.Label>Gizli Cevap</Form.Label>
                        <Form.Control
                            value={this.state.secretAnswer}
                            onChange={this.handleChange}
                            type="input"
                        />
                    </Form.Group>
                    {this.state.error &&
                    <div style={{color: 'red', fontSize: '13px', marginBottom: '10px'}}>{this.state.errorText}</div>}
                    <LoaderButton
                        block
                        disabled={!this.validateForm()}
                        type="submit"
                        isLoading={this.state.isLoading}
                        text="Üye Ol"
                        loadingText="Üye olunuyor…"
                    />
                    <Form.Group style={{paddingTop: '10px'}} controlId="formBasicCheckbox">
                        <Form.Check required label={this.renderCheckboxLabel()}/>
                    </Form.Group>
                </Form>
                <div style={{textAlign: 'center', paddingTop: '10px', marginTop: '20px', borderTop: '1px solid #E8E8E8', fontSize:'15px'}}>Hesabın var mı? <Link to='login'>Giriş yap</Link></div>
            </div>
        );
    }

    render() {
        let modalClose = () => this.setState({
            showProtectionOfPersonalDataContract: false,
            showTermsOfUseContract: false
        });
        return (
            <div className="Signup">
                {this.state.newUser ? this.renderConfirmationForm() : this.renderForm()}
                <Popup show={this.state.showProtectionOfPersonalDataContract || this.state.showTermsOfUseContract} onHide={modalClose}>
                    {
                        this.state.showTermsOfUseContract
                            ?
                            <>
                                <div>SİTE KULLANIM ŞARTLARI</div>
                                <p>Pubgstars.com internet sitesini kullanarak, işbu Kullanım Koşulları’nı, Gizlilik Politikası’nı ve Kişisel Verilerin Korunmasına İlişkin Aydınlatma Metni’ni bir bütün olarak kabul etmiş olursunuz</p>
                                <p>Pubgstars.com hizmetlerinden yararlanan kullanıcı;</p>
                                <ol>
                                    <li>Pubgstars.com’u kullanırken yürürlükteki mevzuatı ihlal etmeyeceğini,</li>
                                    <li>Pubgstars.com’da yayınlanan hiçbir içeriği Pubgstars.com’dan izin almadan ticari amaçla olsun veya olmasın kullanmayacağını; editörler tarafından üretilen her neviden içerikleri başka internet sitelerinde, konvansiyonel ortamlarda ve diğer mecralarda sahibinin adı ve bağlantı (link) belirtmeksizin yayınlayamayacağını,</li>
                                    <li>Pubgstars.com markalarının tescilli markalar olduğunu ve izinsiz kullanılamayacağını,</li>
                                    <li>Herhangi bir yazılım, donanım veya iletişim unsuruna zarar vermek, işlevini aksatmak maksadıyla virüs içeren yazılım veya başka bir bilgisayar kodu, dosyası oluşturmayacağını, yetkisi olmayan herhangi bir sisteme ve/veya veriye ulaşmaya çalışmayacağını,</li>
                                    <li>Ayrıca direkt veya dolaylı olarak, verilen hizmetlerdeki algoritmaları ve kodları deşifre edecek, işlevlerini bozacak davranışlarda bulunmayacağını, İçerikleri değiştirme, dönüştürme, çevirme, alıntı göstermeksizin başka sitelerde yayınlama gibi davranışlarda bulunmayacağını,</li>
                                    <li>Pubgstars.com adresinde yayınlanan içeriklerin paylaşılması sebebiyle doğacak olan tüm hukuki sorumluluğun paylaşan kişiye ait olduğunu, Pubgstars.com’un hiçbir sorumluluğunun olmadığını,</li>
                                    <li>Pubgstars.com birçok içeriğinden başka sitelere bağlantı (link) verilebileceği kabulü ile Pubgstars.com tarafından bağlantı (link) verilen, tavsiye edilen diğer sitelerin bilgi kullanımı, gizlilik ilkeleri ve içeriğinden Pubgstars.com’un sorumlu olmadığını,</li>
                                    <li>Pubgstars.com’un kendi ürettiği veya dışardan aldığı bilgi, belge, yazılım, tasarım, grafik vb. eserlerin 5846 Sayılı Fikir ve Sanat Eserleri Kanunu kapsamında korunduğunu ve eser hakkının ihlali halinde bundan dolayı sorumlu olunacağını</li>
                                    <li>Pubgstars.com‘un kullanıcı üyeliği gerektirmeyen hizmetleri zaman içinde üyelik gerektiren bir hale dönüştürebileceğini, ilave hizmetler açabileceğini, bazı hizmetlerini kısmen veya tamamen değiştirebileceği veya ücretli hale dönüştürebileceğini,</li>
                                    <li>Kullanıcının içerik oluşturmasına izin verilen yorumlarda ya da forumlarda ya da diğer platformlarda içerik oluşturan kişilerin bu içerikten dolayı tamamen ve tek başına sorumlu olduğunu,</li>
                                    <li>Pubgstars.com’un içerik oluşturan kullanıcı ile ilgili kısıtlama ve kullanıcı engelleme, silme hakkına sahip olduğunu</li>
                                    <li>Pubgstars.com’daki bilgilerin güncelliği, doğruluğu, şartları, kalitesi, performansı, pazarlanabilirliği, belli bir amaca uygunluğu ve diğer bilgi, hizmet veya ürünlere etkisi ile tamlığı ve/veya kesintisiz devamlılık, güncelleme, işlevsellik, doğruluk, hatasızlık hakkında herhangi bir şekilde, Pubgstars.com tarafından açık ya da zımni olarak garanti verilmediğini ve taahhütte bulunulmadığını, Pubgstars.com’un gerekli gördüğü zamanlarda hizmetleri geçici bir süre askıya alabileceğini veya tamamen durdurabileceğini, hizmetlerin geçici bir süre askıya alınması veya tamamen durdurulmasından dolayı kullanıcılara karşı herhangi bir sorumluluğunun olmadığını,</li>
                                    <li>Pubgstars.com’un hizmetleri, tasarımı ve içeriği her zaman değiştirebilme hakkını saklı tuttuğunu ve sunulan hizmetlerin kullanıcılara kazanılmış hak tahsis etmeyeceğini,</li>
                                    <li>Pubgstars.com’a incelenmek üzere gönderilen bilgi, belge ve veriler arasında Türkiye Cumhuriyeti devletinin gizli tutulması gereken bilgilerinden olan, ticaret şirketleri ve sair tüzel
                                        kişilerin ticari sırrı niteliğinde olan ve/veya üçüncü kişilerin kişilik haklarını zedeleyen bilgi, belge veya verilerin bulunması durumunda Pubgstars.com’un hukuki ve/veya cezai sorumluluğunun olmadığını, bu konudaki sorumluluğun ilgili bilgi, belge veya veriyi gönderen kullanıcıda olduğunu, bu türden bilgi, belge ve/veya verilerin kaydedilip kaydedilmemesi konusunda veyahut da her ne şekilde olursa olsun kullanılıp kullanılmaması hususunda tüm inisiyatifin Pubgstars.com’a ait olduğunu,</li>
                                    <li>Pubgstars.com’un kullanım koşulları ve kuralları her zaman tek taraflı değiştirme hakkının saklı olduğunu</li>
                                </ol>

                                <p>ve bu kapsamda, işbu Kullanım Koşulları’nın, Gizlilik Politikası’nın ve Kişisel Verilerin Korunmasına İlişkin Aydınlatma Metni’nin tamamını okuduğunu ve bir bütün olarak onayladığını, Pubgstars.com’a mümkün olan tüm mecralar vasıtası ile ulaşarak ayrıca onay gerekmeksizin belirtilen tüm metinleri ve de metinlerde belirtilen kural ve koşulları bir bütün olarak kabul, beyan ve taahhüt eder.</p>
                            </>
                            :
                            <>
                                <div>KİŞİSEL VERİ RIZA METNİ</div>
                                <p>Değerli faydalanıcımız,</p>
                                <p>Pubgstars.com platformu olarak kişisel verilerinizin korunmasını ve güvenliğini önemsiyoruz.</p>
                                <p>Kullanıcılarımızın Pubgstars.com sitesi üzerinden sunduğumuz hizmetlere güvenli bir şekilde ulaşması maksadıyla sağladığınız kişisel verilerinizin işlenmesi ve üçüncü kişilere aktarılması gerekmektedir. Bu kapsamda, Pubgstars.com ile paylaştığınız 6698 sayılı Kişisel Verilerin Korunması Kanunu kapsamındaki kişisel veriler, mevzuata uygun şekilde, faaliyet konumuz ve hizmet amaçlarımızla bağlantılı ve de ölçülü olarak işlenebilecek, gerekli olması durumunda üçüncü kişilere aktarılabilecek ve mevzuata uygun süreler boyunca saklanabilecektir.</p>
                                <p>Kişisel verilerin işlenmesi bakımından 6698 sayılı Kişisel Verilerin Korunması Kanunu kapsamında veri sorumlusu, ayrı ayrı olmak üzere Pubgstars.com’dir.</p>
                                <p>Kişisel verileriniz;</p>
                                <ol>
                                    <li>Hukuka ve dürüstlük kuralının öngördüğü biçimde,</li>
                                    <li>Doğru ve güncel olarak,</li>
                                    <li>Belirli, açık ve meşru amaçlar için ve de</li>
                                    <li>İşlenme amaçları ile bağlantılı, sınırlı ve ölçülü olarak</li>
                                </ol>
                                <p>işlenecektir.</p>
                                <p>Kişisel verileriniz, Pubgstars.com’da sunulan hizmetlerden yararlanabilmeniz için sizleri hizmetlerimizden haberdar etmek amacıyla işlenmekte ve muhafaza edilmektedir. İşlenen kişisel verileriniz, Pubgstars.com tarafından sunulan hizmete bağlı olarak değişkenlik gösterebilmekle; internet sitesi ve benzeri vasıtalarla sözlü, yazılı ya da elektronik olarak toplanabilmektedir. Pubgstars.com’un hizmetlerinden yararlandığınız süre boyunca, 6698 sayılı Kişisel Verilerin Korunması Kanunu’nun yürürlük tarihinden önce verdiğiniz rıza veya yürürlük tarihinden sonra verdiğiniz açık rıza veyahut da 6698 sayılı Kişisel Verilerin Korunması Kanunu’nda belirtilen kural ve koşullar çerçevesinde kişisel verileriniz işlenebilecektir.</p>
                                <p>Toplanan kişisel verileriniz, Pubgstars.com tarafından sunulan hizmetlerden kolayca ve güvenilir şekilde faydalanabilmeniz için gerekli çalışmaların yürütülmesi, hizmetlerimizin taleplerinizi, ihtiyaçlarınızı ve kullanım alışkanlıklarınızı karşılayacak şekilde özelleştirilerek sizlere önerilebilmesi, Pubgstars.com’un ticari stratejilerinin belirlenmesi ve değerlendirilmesi, sizlerle Pubgstars.com arasındaki hukuki süreçlerin hızlı ve etkili bir şekilde yürütülebilmesinin temini amaçlarıyla, Kanun’da belirtilen koşullara uygun şekilde işlenecektir.</p>
                                <p>Özel nitelikli kişisel verileriniz sizden talep edilmeyecek ve işlenmeyecektir.</p>
                                <p>Kişisel verileriniz yalnızca Pubgstars.com’un faaliyet konusu kapsamında ve toplanma amaçları doğrultusunda toplanacak, işlenme amaçlarının gerektirdiği süreler boyunca saklanabilecek, 6698 sayılı Kişisel Verilerin Korunması Kanunu’nda belirtilen kuralları aşar şekilde işlenmeyecek ve işlenmesini gerektiren sebeplerin ortadan kalkması hâlinde, yürürlükteki diğer mevzuattan doğan saklama mecburiyeti
                                    bulunan haller saklı olmak üzere, resen veya ilgili kişi olarak talebiniz üzerine silinecek veya yok edilecek veyahut da anonim hale getirilecektir.</p>
                                <p>Kişisel verileriniz, Pubgstars.com’un bu Aydınlatma Metni’nde belirtilen amaçları dahilinde, yürürlükteki mevzuata uygun olarak, sayılanlar ile sınırlı olmamak üzere; yukarıda belirtilen veri sorumluları arasında aktarılacaktır. Bunun yanında, kişisel veri işleme amaçları doğrultusunda gerekli güvenlik önlemlerini de almak sureti ile kişisel veriler üçüncü kişilere, iş ortaklarımıza, ifa yardımcılarımıza, ödeme hizmeti sağlayan kuruluşlara ve sunduğumuz hizmet ve faaliyet amacımız doğrultusunda ya da ilgili mevzuatın öngördüğü durumlarda düzenleyici denetleyici kurumlara ve de resmi mercilere aktarılabilecektir.</p>
                                <p>Kişisel veri sahibi olarak, kişisel verilerinizin işlenip işlenmediğinin sorgulama, işlenmişse bu konuda bilgi talep etme, kişisel verilerinizin işlenme amacını ve amacına uygun kullanılıp kullanılmadığını, yurt içi veya yurtdışına aktarılıp aktarılmadığını sorgulama, eksik veya yanlış işlenmişse düzeltilmesini isteme, kişisel verilerinizin kanuna aykırı işlenmesi nedeniyle zarara uğramanız halinde zararınızı giderilmesini isteme, silinmesini veya yok edilmesini talep etme haklarınız bulunmaktadır. Belirtilen bu haklarınızı kullanmak için, bu yöndeki talebinizi Kanun’un 13. maddesi doğrultusunda ve Veri Sorumlusuna Başvuru Usul ve Esasları Hakkında Tebliğ’de belirtilen kurallara uygun olarak, Pubgstars.com’un aşağıda belirtilen mail adresine, kimliğinizi tespite yarayan gerekli bilgi ve belgeleri de eklemek sureti ile yazılı olarak başvuruda bulunabilirsiniz. Başvurunuz Pubgstars.com tarafından değerlendirilecek ve otuz gün içinde ücretsiz olarak sonuçlandırılacaktır.</p>
                                <p>Mail: info@pubgstars.com</p>
                            </>
                    }
                </Popup>
            </div>
        );
    }
}
