package mysql_db

import (
	"context"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uyouii/oauth/common"
	"github.com/uyouii/oauth/dao/db_base"
	"xorm.io/xorm"
)

type OAuthDao struct {
	Config MysqlDbConfig
	Engine *xorm.Engine
}

func GetNewOauthDb(config *MysqlDbConfig) db_base.OauthDbInterface {
	dao := &OAuthDao{}
	err := dao.Init(config)
	if err != nil {
		panic(err)
	}
	return dao
}

func (d *OAuthDao) Init(config *MysqlDbConfig) error {
	infof, errorf := common.GetLogFuns(context.Background())

	if config.Addr == "" && config.Port == 0 {
		config.Addr = "127.0.0.1"
		config.Port = 3306
	}

	infof("init mysql with config: %+v", config)

	datasource := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
		config.User, config.Password, config.Addr, config.Port, config.DatabaseName)
	engine, err := xorm.NewEngine("mysql", datasource)
	if err != nil {
		errorf("init mysql failed, err: %v", err)
		return err
	}

	d.Config = *config
	d.Engine = engine

	infof("init mysql db engine success, source: %v", datasource)

	return nil
}

func (d *OAuthDao) insert(ctx context.Context, data interface{}) error {
	infof, errorf := common.GetLogFuns(ctx)

	session := d.Engine.Context(ctx)

	_, err := session.Insert(data)

	if err != nil {
		errorf("insert failed, err: %v, data: %+v", err, data)
		return err
	}

	infof("insert success, data: %+v", data)

	return nil
}

func (d *OAuthDao) update(ctx context.Context, data interface{}, cond interface{}) error {
	infof, errorf := common.GetLogFuns(ctx)

	session := d.Engine.Context(ctx)

	_, err := session.Update(data, cond)
	if err != nil {
		errorf("update failed, err: %v, data: %+v", err, data)
		return err
	}

	infof("update success, data: %+v", data)

	return nil
}

func (d *OAuthDao) delete(ctx context.Context, data interface{}) error {
	infof, errorf := common.GetLogFuns(ctx)

	session := d.Engine.Context(ctx)

	_, err := session.Delete(data)
	if err != nil {
		errorf("delete failed, err: %v, data: %+v", err, data)
		return err
	}

	infof("delete success, data: %+v", data)

	return nil
}
