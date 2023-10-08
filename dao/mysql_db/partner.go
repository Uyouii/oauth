package mysql_db

import (
	"context"
	"time"

	"uyouii.cool/oauth/common"
	"uyouii.cool/oauth/dao/db_base"
	"uyouii.cool/oauth/utils"
)

func (d *OAuthDao) convertPartner(partner *PartnerTab) *db_base.Partner {
	res := &db_base.Partner{
		PartnerName:   partner.PartnerName,
		PartnerKey:    partner.PartnerKey,
		PartnerSecret: partner.PartnerSecret,
		Expire:        partner.Expire,
		CreateTime:    time.Unix(partner.CreateTimestamp, 0),
		UpdateTime:    time.Unix(partner.UpdateTimestamp, 0),
	}
	return res
}

func (d *OAuthDao) GenPartner(ctx context.Context, partnerName string, expire int64) (*db_base.Partner, error) {
	infof, errorf := common.GetLogFuns(ctx)

	if partnerName == "" || expire == 0 {
		errorf("invalid params, partnerName: %v, expire: %v", partnerName, expire)
		return nil, common.GetError(common.INVALID_PARAMS)
	}

	partner := &PartnerTab{
		PartnerName:     partnerName,
		PartnerKey:      utils.GetUuid(),
		PartnerSecret:   utils.GetRandomSecret(),
		Expire:          expire,
		CreateTimestamp: time.Now().Unix(),
		UpdateTimestamp: time.Now().Unix(),
	}

	err := d.insert(ctx, partner)
	if err != nil {
		errorf("insert partner failed, err: %v, partner: %+v", err, partner)
		return nil, err
	}

	infof("gen new partner success, partner: %+v", partner)

	return d.convertPartner(partner), nil
}

func (d *OAuthDao) getPartnerByKey(ctx context.Context, partnerKey string) (*PartnerTab, error) {
	infof, errorf := common.GetLogFuns(ctx)

	if partnerKey == "" {
		errorf("invalid params, partnerKey: %v", partnerKey)
		return nil, common.GetError(common.INVALID_PARAMS)
	}

	partner := &PartnerTab{}

	session := d.Engine.Context(ctx).Where("partner_key = ?", partnerKey)

	exists, err := session.Get(partner)
	if err != nil {
		errorf("get partner failed, err: %v, partnerKey: %v", err, partnerKey)
		return nil, err
	}

	if !exists {
		errorf("get partner failed, partner key not existes, partner key: %v", partnerKey)
		return nil, common.GetErrorWithMsg(common.ERROR_EMPTY, "partner key not existes")
	}

	infof("get partner success, partenr: %+v", partner)

	return partner, nil
}

func (d *OAuthDao) GetPartnerByKey(ctx context.Context, partnerKey string) (*db_base.Partner, error) {
	partner, err := d.getPartnerByKey(ctx, partnerKey)
	if err != nil {
		return nil, err
	}
	return d.convertPartner(partner), nil
}

func (d *OAuthDao) DeletePartner(ctx context.Context, partnerKey string) error {
	infof, errorf := common.GetLogFuns(ctx)

	if partnerKey == "" {
		errorf("invalid params, partnerKey: %v", partnerKey)
		return common.GetError(common.INVALID_PARAMS)
	}

	partner := &PartnerTab{
		PartnerKey: partnerKey,
	}

	err := d.delete(ctx, partner)
	if err != nil {
		errorf("delete partner failed, err: %v, partnerkey: %v", err, partnerKey)
		return err
	}

	infof("delete partner key success, partnerkey: %v", partnerKey)

	return nil
}
