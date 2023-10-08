package oauth

import (
	"context"
	"sync"
	"time"

	"github.com/uyouii/oauth/common"
	"github.com/uyouii/oauth/dao/db_base"
	"github.com/uyouii/oauth/dao/mysql_db"
)

type OauthManger struct {
	db                db_base.OauthDbInterface
	deleteProcessOnce sync.Once
}

func GetOauthManager() *OauthManger {
	return &OauthManger{}
}

func (g *OauthManger) WithMysql(config *mysql_db.MysqlDbConfig) *OauthManger {
	g.db = mysql_db.GetNewOauthDb(config)
	return g
}

// type OauthDbInterface interface {
// 	GenPartner(ctx context.Context, PartnerName string, expire int64) (*Partner, error)
// 	GetPartnerByKey(ctx context.Context, partnerKey string) (*Partner, error)
// 	DeletePartner(ctx context.Context, partnerKey string) error
// 	GenOauthToken(ctx context.Context, partnerKey string) (*OauthToken, error)
// 	GetTokenInfo(ctx context.Context, token string) (*OauthToken, error)
// 	DeleteOauthToken(ctx context.Context, token string) error
// 	DeleteExpireTokens(ctx context.Context) error
// }

func (g *OauthManger) GenPartner(ctx context.Context, PartnerName string, expire int64) (*db_base.Partner, error) {
	if g.db == nil {
		return nil, common.GetError(common.ERROR_SYSTEM)
	}
	return g.db.GenPartner(ctx, PartnerName, expire)
}

func (g *OauthManger) GetPartnerByKey(ctx context.Context, partnerKey string) (*db_base.Partner, error) {
	if g.db == nil {
		return nil, common.GetError(common.ERROR_SYSTEM)
	}
	return g.db.GetPartnerByKey(ctx, partnerKey)
}

func (g *OauthManger) DeletePartner(ctx context.Context, partnerKey string) error {
	if g.db == nil {
		return common.GetError(common.ERROR_SYSTEM)
	}
	return g.db.DeletePartner(ctx, partnerKey)
}

func (g *OauthManger) GenOauthToken(ctx context.Context, partnerKey string, paratnerSecret string) (*db_base.OauthToken, error) {
	_, errorf := common.GetLogFuns(ctx)
	if g.db == nil {
		return nil, common.GetError(common.ERROR_SYSTEM)
	}

	partner, err := g.db.GetPartnerByKey(ctx, partnerKey)
	if err != nil {
		errorf("get partner info failed, err: %v", err)
		return nil, err
	}

	if partner.PartnerSecret != paratnerSecret {
		errorf("invalid secret, req secret: %v, partner secret: %v", paratnerSecret, partner.PartnerSecret)
		return nil, common.GetError(common.INVALID_SECRET)
	}

	return g.db.GenOauthToken(ctx, partnerKey)
}

func (g *OauthManger) CheckOauthToken(ctx context.Context, partnerKey string, token string) error {
	_, errorf := common.GetLogFuns(ctx)

	if g.db == nil {
		return common.GetError(common.ERROR_SYSTEM)
	}

	if partnerKey == "" || token == "" {
		errorf("invalid params, partkerkey: %v, token: %v", partnerKey, token)
		return common.GetError(common.INVALID_PARAMS)
	}

	tokenInfo, err := g.db.GetTokenInfo(ctx, token)
	if err != nil {
		errorf("get token info failed, err: %v, partkerkey: %v, token: %v", err, partnerKey, token)
		return err
	}

	if tokenInfo.PartnerKey != partnerKey {
		errorf("partner key not match, req partkerkey: %v, token partnerkey: %v, token: %v",
			err, partnerKey, tokenInfo.PartnerKey, token)
		return common.GetError(common.INVALID_TOKEN)
	}

	if tokenInfo.ExpireTime.Before(time.Now()) {
		errorf("token already expired, token: %v, expiredTime: %+v", token, tokenInfo.ExpireTime)
		return common.GetError(common.TOKEN_EXPIRED)
	}

	return nil
}

func (g *OauthManger) StartDeleteExpireTokenProcess() error {
	if g.db == nil {
		return common.GetError(common.ERROR_SYSTEM)
	}

	deleteProcess := func() {
		ctx := context.Background()
		infof, errorf := common.GetLogFuns(ctx)
		for {
			time.Sleep(time.Second * 3600 * 12)
			err := g.db.DeleteExpireTokens(ctx)
			if err != nil {
				errorf("delete expire token failed, err: %v", err)
			} else {
				infof("delete expired token success")
			}
		}
	}

	g.deleteProcessOnce.Do(func() {
		go deleteProcess()
	})
	return nil
}
