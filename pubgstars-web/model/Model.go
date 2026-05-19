package model

// Models
type Game struct {
	Id                  string  `json:"id"`
	GameDate            string  `json:"gameDate"`
	League              string  `json:"league"`
	Type                string  `json:"type"`
	Map                 string  `json:"map"`
	Price               float64 `json:"price"`
	RoomPassword        string  `json:"roomPassword"`
	Status              string  `json:"status"`
	RegisteredUserCount int32   `json:"registeredUserCount"`
	TotalIncome         float64 `json:"totalIncome"`
	TeamPlayerCount     int32   `json:"teamPlayerCount"`
	Platform            string  `json:"platform"`
	InsertedAt          int64   `json:"insertedAt"`
	InsertedBy          string  `json:"insertedBy"`
	UpdatedAt           int64   `json:"updatedAt"`
	UpdatedBy           string  `json:"updatedBy"`
	Discord             string  `json:"discord"`
	Registered          bool    `json:"registered"`
	ShowPassword        bool    `json:"showPassword"`
	Cancellable         bool    `json:"cancellable"`

	Winner1st string  `json:"winner1st"`
	Winner2nd string  `json:"winner2nd"`
	Winner3rd string  `json:"winner3rd"`
	Award1st  float64 `json:"award1st"`
	Award2nd  float64 `json:"award2nd"`
	Award3rd  float64 `json:"award3rd"`

	Winner1stName string `json:"winner1stName"`
	Winner2ndName string `json:"winner2ndName"`
	Winner3rdName string `json:"winner3rdName"`
}

func (game *Game) Detail() string {
	return game.League + "/" + game.Map + "/" + game.Type + "/" + game.Platform + " @" + game.GameDate
}

type User struct {
	Id             string  `json:"id"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	SecretQuestion string  `json:"secretQuestion"`
	SecretAnswer   string  `json:"secretAnswer"`
	Status         string  `json:"status"`
	Balance        float64 `json:"balance"`
	Bonus          float64 `json:"bonus"`
	Gain           float64 `json:"gain"`
	InsertedAt     int64   `json:"insertedAt"`
	InsertedBy     string  `json:"insertedBy"`
	UpdatedAt      int64   `json:"updatedAt"`
	UpdatedBy      string  `json:"updatedBy"`
}

func (user *User) GetAvailableBalance() float64 {
	return user.Balance + user.Bonus
}

type GameUser struct {
	GameId string   `json:"gameId"`
	Users  []string `json:"users"`
}

type UserGame struct {
	UserId string   `json:"userId"`
	Games  []string `json:"games"`
}

// Transaction logs
type TransactionType int
type SubTransactionType int

type Message struct {
	Id         string `json:"id"`
	DateTime   int64  `json:"dateTime"`
	IsCustomer bool   `json:"isCustomer"`
	Status     string `json:"status"`
	From       string `json:"from"`
	Message    string `json:"message"`
	Comment    string `json:"comment"`
	UpdatedAt  int64  `json:"updatedAt"`
	UpdatedBy  string `json:"updatedBy"`
}

/*
Balance (deposit, withdraw, registerGame, unregisterGame, winGame, bonus)
Game (registerGame, unregisterGame, winGame)
Account (activateAccount, deactivateAccount, changePassword, forgetPassword)
Admin ()
*/
const (
	BALANCE                 TransactionType    = 100
	BALANCE_DEPOSIT         SubTransactionType = 101
	BALANCE_WITHDRAW        SubTransactionType = 102
	BALANCE_REGISTER_GAME   SubTransactionType = 103
	BALANCE_UNREGISTER_GAME SubTransactionType = 104
	BALANCE_LOAD            SubTransactionType = 105

	GAME            TransactionType    = 200
	GAME_REGISTER   SubTransactionType = 201
	GAME_UNREGISTER SubTransactionType = 202
	GAME_WIN        SubTransactionType = 203

	ACCOUNT                 TransactionType    = 300
	ACCOUNT_ACTIVATE        SubTransactionType = 301
	ACCOUNT_DEACTIVATE      SubTransactionType = 302
	ACCOUNT_CHANGE_PASSWORD SubTransactionType = 303
	ACCOUNT_FORGET_PASSWORD SubTransactionType = 304

	ADMIN TransactionType = 1000
)

var TRANSACTION_MAP = map[TransactionType]string{
	BALANCE: "Bakiye",
	GAME:    "Oyun",
	ACCOUNT: "Hesap",
	ADMIN:   "Admin",
}

var SUB_TRANSACTION_MAP = map[SubTransactionType]string{
	BALANCE_DEPOSIT:         "Para yukleme",
	BALANCE_WITHDRAW:        "Para cekme",
	BALANCE_REGISTER_GAME:   "Oyun kayit",
	BALANCE_UNREGISTER_GAME: "Oyun iptal",
	BALANCE_LOAD:            "Bakiye yukleme",
	GAME_REGISTER:           "Oyun kayit",
	GAME_UNREGISTER:         "Oyun iptal",
	GAME_WIN:                "Oyun kazanan",
	ACCOUNT_ACTIVATE:        "Hesap aktiflestirme",
	ACCOUNT_DEACTIVATE:      "Hesap kapama",
	ACCOUNT_CHANGE_PASSWORD: "Parola degisikligi",
	ACCOUNT_FORGET_PASSWORD: "Parolami unuttum",
}

type TransactionLog struct {
	Id                 string             `json:"id"`
	UserId             string             `json:"userId"`
	Operator           string             `json:"operator"`
	InsertedAt         int64              `json:"insertedAt"`
	TransactionType    TransactionType    `json:"transactionType"`
	SubTransactionType SubTransactionType `json:"subTransactionType"`
	Detail             string             `json:"detail"`
}

type TransactionLogResponse struct {
	Id                 string `json:"id"`
	UserName           string `json:"userName"`
	DateTime           string `json:"dateTime"`
	TransactionType    string `json:"transactionType"`
	SubTransactionType string `json:"subTransactionType"`
	Detail             string `json:"detail"`
}
