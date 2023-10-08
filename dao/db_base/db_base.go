package db_base

import (
	"context"
	"time"
)

type Partner struct {
	PartnerName   string
	PartnerKey    string
	PartnerSecret string
	Expire        int64
	CreateTime    time.Time
	UpdateTime    time.Time
}

type OauthToken struct {
	PartnerKey string
	Token      string
	ExpireTime time.Time
	CreateTime time.Time
	UpdateTime time.Time
}

type OauthDbInterface interface {
	GenPartner(ctx context.Context, PartnerName string, expire int64) (*Partner, error)
	GetPartnerByKey(ctx context.Context, partnerKey string) (*Partner, error)
	DeletePartner(ctx context.Context, partnerKey string) error
	GenOauthToken(ctx context.Context, partnerKey string) (*OauthToken, error)
	GetTokenInfo(ctx context.Context, token string) (*OauthToken, error)
	DeleteOauthToken(ctx context.Context, token string) error
	DeleteExpireTokens(ctx context.Context) error
}
