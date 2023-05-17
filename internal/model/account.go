package model

type Account struct {
	AccountId   string `json:"accountId" gorm:"column:uuid"`
	AccountName string `json:"accountName" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	CountryCode string `json:"countryCode" gorm:"column:country_code"`
	Status      int32  `json:"status" gorm:"column:status"`
}

func (*Account) TableName() string {
	return "accounts"
}
