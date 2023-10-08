package mysql_db

type TokenTab struct {
	Id              int64  `xorm:"pk autoincr BIGINT(20)"`
	PartnerKey      string `xorm:"index(pratner_token_index) VARCHAR(64)"`
	Token           string `xorm:"index(pratner_token_index) unique VARCHAR(128)"`
	ExpireTimestamp int64  `xorm:"BIGINT(20)"`
	CreateTimestamp int64  `xorm:"BIGINT(20)"`
	UpdateTimestamp int64  `xorm:"BIGINT(20)"`
}

type PartnerTab struct {
	Id              int64  `xorm:"pk autoincr BIGINT(20)"`
	PartnerName     string `xorm:"VARCHAR(64)"`
	PartnerKey      string `xorm:"unique VARCHAR(64)"`
	PartnerSecret   string `xorm:"VARCHAR(64)"`
	Expire          int64  `xorm:"default 3600 BIGINT(20)"`
	CreateTimestamp int64  `xorm:"BIGINT(20)"`
	UpdateTimestamp int64  `xorm:"BIGINT(20)"`
}
