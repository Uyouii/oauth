package mysql_db

type MysqlDbConfig struct {
	User         string
	Password     string
	Addr         string
	Port         int64
	DatabaseName string
}
