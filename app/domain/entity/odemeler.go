package entity

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	tr_translations "gopkg.in/go-playground/validator.v9/translations/tr"
)

const (
	OdemelerGeciciDurum               = 1 ///TODO: onemli diger durumlar silindi gerekirse durum eklemek icin, burasi rezerve
	OdemelerDurumIkiGrupYerDegistirdi = 3 //  iki grup arasi degisim yapildi --- // TODO : iptal burada olmamali belki kurbanda  olabilir, belki burada durmasinin zarari olmaz

	OdemelerBorcDurumIlkEklenenFiyat        = 1 // ilk eklenen fiyat değeri
	OdemelerBorcDurumTaksitOdemesi          = 2 // taksit eklemiş
	OdemelerBorcDurumKasaBorcluDurumda      = 3 //  kasa borçlu kalmışsa---TODO: burada event surekli degisir bu mantikli mi kontrol edilmeli
	OdemelerBorcDurumKaporaOdemesiHayvanBos = 4 //  kapora odendi ama hayvan atanmamışdır
	OdemelerBorcDurumKaporaOdemesiHayvanVar = 5 //  kapora odendi ve hayvan atanmışdır
	OdemelerBorcDurumIkiGrupYerDegistirdi   = 5 //  iki grup arasi degisim yapildi --- iptal burada olmamali belki durumda olabilir
	OdemelerBorcDurumHesapKapandi           = 6 //  tüm borcunu ödemiş
	OdemelerBorcDurumHicOdemeYapmadi        = 7 //  hic odeme yapmamis hizli ekleme gibi durumlar
	//OdemelerBorcDurumHayvanKesildi                      = 7  //  hayvan kesildi
	//OdemelerBorcDurumHayvanKesildiTeslimEdildi          = 8  //  hayvan kesildi teslim edildi
	OdemelerBorcDurumFiyatManuelDegistirildiEsitFiyat   = 9  //  kurbanın düzenleme alanında fiyatı değiştirilmiş
	OdemelerBorcDurumFiyatManuelDegistirildiYuksekFiyat = 10 //  kurbanın düzenleme alanında fiyatı değiştirilmiş
	OdemelerBorcDurumFiyatManuelDegistirildiDusukFiyat  = 11 //  kurbanın düzenleme alanında fiyatı değiştirilmiş
	//OdemelerBorcDurumHayvanEklendi                      = 12 //  hayvanEklemesi yapılmış
	OdemelerBorcDurumTaksitSilindi         = 13 // taksit silindi
	OdemelerBorcDurumTaksitSilindiBildirim = 14 // taksit silindi rapor alanı içindir
)

// Odemeler strcut
type SacrificeOdemeler struct {
	ID             uint64     `gorm:"primary_key;auto_increment" json:"id"`
	KurbanID       uint64     `gorm:"not null;" json:"kurban_id"`
	UserID         uint64     `gorm:"not null;" json:"user_id"`
	Aciklama       string     `gorm:"type:text ;" json:"aciklama"`
	Makbuz         string     `gorm:"type:varchar(255) ; null;"  json:"makbuz" `
	Durum          int        `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" `
	BorcDurum      int        `gorm:"type:smallint ;NOT NULL;DEFAULT:'1'" validate:"required"`
	KurbanFiyati   float64    `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';"  validate:"required,numeric" `
	VerilenUcret   float64    `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';"  validate:"required,numeric" `
	KasaBorcu      float64    `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';" validate:"numeric" `
	Alacak         float64    `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';" validate:"numeric" `
	Bakiye         float64    `gorm:"type:decimal(10,2); NOT NULL; DEFAULT:'0';" validate:"numeric" `
	VerildigiTarih time.Time  `gorm:"type:date ; NULL;"  json:"VerildigiTarih"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

// OdemelerTableName table name
var OdemelerTableName string = "sacrifice_odemeler"

// BeforeSave init
func (gk *SacrificeOdemeler) BeforeSave() {
	gk.Aciklama = html.EscapeString(strings.TrimSpace(gk.Aciklama))
}

// Prepare init
func (gk *SacrificeOdemeler) Prepare() {
	gk.Aciklama = html.EscapeString(strings.TrimSpace(gk.Aciklama))
	gk.CreatedAt = time.Now()
	gk.UpdatedAt = time.Now()
}

// TableName override
func (gk *SacrificeOdemeler) TableName() string {
	return OdemelerTableName
}

// Validate fluent validation
func (gk *SacrificeOdemeler) Validate() map[string]string {
	var (
		validate *validator.Validate
		uni      *ut.UniversalTranslator
	)
	tr := en.New()
	uni = ut.New(tr, tr)
	trans, _ := uni.GetTranslator("tr")
	validate = validator.New()
	tr_translations.RegisterDefaultTranslations(validate, trans)

	errorLog := make(map[string]string)

	err := validate.Struct(gk)
	fmt.Println(err)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		fmt.Println(errs)
		for _, e := range errs {
			// can translate each error one at a time.
			lng := strings.Replace(e.Translate(trans), e.Field(), "Burası", 1)
			errorLog[e.Field()+"_error"] = e.Translate(trans)
			// errorLog[e.Field()] = e.Translate(trans)
			errorLog[e.Field()] = lng
			errorLog[e.Field()+"_valid"] = "is-invalid"
		}
	}
	return errorLog
}
