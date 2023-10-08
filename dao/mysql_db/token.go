package mysql_db

import (
	"context"
	"time"

	"github.com/uyouii/oauth/common"
	"github.com/uyouii/oauth/dao/db_base"
	"github.com/uyouii/oauth/utils"
)

func (d *OAuthDao) convertToken(token *TokenTab) *db_base.OauthToken {
	res := &db_base.OauthToken{
		PartnerKey: token.PartnerKey,
		Token:      token.Token,
		ExpireTime: time.Unix(token.ExpireTimestamp, 0),
		CreateTime: time.Unix(token.CreateTimestamp, 0),
		UpdateTime: time.Unix(token.UpdateTimestamp, 0),
	}
	return res
}

// #TODO: add rate limit
// #TODO: add redis cache for partner info
// #TODO: add local cache for partner info
func (d *OAuthDao) GenOauthToken(ctx context.Context, partnerKey string) (*db_base.OauthToken, error) {
	infof, errorf := common.GetLogFuns(ctx)

	if partnerKey == "" {
		errorf("invalid params, partnerKey: %v", partnerKey)
		return nil, common.GetError(common.INVALID_PARAMS)
	}

	partner, err := d.getPartnerByKey(ctx, partnerKey)
	if err != nil {
		errorf("get partner info failed, err: %v", err)
		return nil, err
	}

	token := &TokenTab{
		PartnerKey:      partnerKey,
		Token:           utils.GetUuid() + utils.GetUuid(),
		ExpireTimestamp: time.Now().Unix() + partner.Expire,
		CreateTimestamp: time.Now().Unix(),
		UpdateTimestamp: time.Now().Unix(),
	}

	err = d.insert(ctx, token)
	if err != nil {
		errorf("insert new oauth failed, err: %v, token: %+v", err, token)
		return nil, err
	}

	infof("gen new oauth success, token: %+v", token)

	return d.convertToken(token), nil
}

func (d *OAuthDao) GetTokenInfo(ctx context.Context, token string) (*db_base.OauthToken, error) {
	infof, errorf := common.GetLogFuns(ctx)

	if token == "" {
		errorf("invalid params, token: %v", token)
		return nil, common.GetError(common.INVALID_PARAMS)
	}

	tokenInfo := &TokenTab{}

	session := d.Engine.Context(ctx).Where("token = ?", token)

	exists, err := session.Get(tokenInfo)
	if err != nil {
		errorf("get token failed, err: %v, token: %v", err, token)
		return nil, err
	}

	if !exists {
		errorf("get token failed, token not existes, token: %v", token)
		return nil, common.GetErrorWithMsg(common.ERROR_EMPTY, "token not existes")
	}

	infof("get token success, token: %+v", tokenInfo)

	return d.convertToken(tokenInfo), nil
}

func (d *OAuthDao) DeleteOauthToken(ctx context.Context, token string) error {
	infof, errorf := common.GetLogFuns(ctx)

	if token == "" {
		errorf("invalid params, token: %v", token)
		return common.GetError(common.INVALID_PARAMS)
	}

	tokenInfo := &TokenTab{
		Token: token,
	}

	err := d.delete(ctx, tokenInfo)
	if err != nil {
		errorf("delete token failed, err: %v, token: %v", err, token)
		return err
	}

	infof("delete token success, token: %v", token)

	return nil
}

func (d *OAuthDao) DeleteExpireTokens(ctx context.Context) error {
	infof, errorf := common.GetLogFuns(ctx)

	tokenInfo := &TokenTab{}

	session := d.Engine.Context(ctx).Where("expire_timestamp < ?", time.Now().Unix())

	effected, err := session.Delete(tokenInfo)
	if err != nil {
		errorf("delete expire token failed, err: %v", err)
		return err
	}

	infof("delete %v expired tokens successfully", effected)

	return nil
}
