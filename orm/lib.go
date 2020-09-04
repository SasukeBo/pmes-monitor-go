package orm

import (
	"fmt"
	"github.com/SasukeBo/configer"
	"github.com/SasukeBo/log"
	"github.com/jinzhu/gorm"
	"time"
	// set db driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func init() {
	var err error
	var uri = createUriWithDBName("mysql")
	var dbname = configer.GetString("db_name")

	reconnectLimit := 5
	for {
		conn, err := gorm.Open("mysql", uri)
		if err != nil && reconnectLimit > 0 {
			log.Errorln(err)
			reconnectLimit--
			time.Sleep(time.Duration(5-reconnectLimit) * 2 * time.Second)
			log.Info("open connection with %s failed, try again ...\n", uri)
			continue
		}
		conn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))
		conn.Close()
		break
	}

	DB, err = gorm.Open("mysql", createUriWithDBName(dbname))
	if err != nil {
		panic(err)
	}
	DB.LogMode(false)
	env := configer.GetString("env")

	err = DB.AutoMigrate(
		&Dashboard{},
		&DeviceErrorCode{},
		&Device{},
		&Attachment{},
		&DeviceErrorLog{},
		&ProduceLog{},
		&DeviceType{},
	).Error
	if err != nil {
		panic(fmt.Errorf("migrate to db error: \n%v", err.Error()))
	}

	if env != "test" && env != "TEST" {
		tableNames := []string{"device_errors", "devices", "attachments", "device_error_logs", "produce_logs"}
		setupUTF8GeneralCI(tableNames)
	}
	setupIndex()
	DB.LogMode(true)
}

func setupIndex() {
	return
}

func createUriWithDBName(name string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configer.GetString("db_user"),
		configer.GetString("db_pass"),
		configer.GetString("db_host"),
		configer.GetString("db_port"),
		name,
	)
}

func setupUTF8GeneralCI(tableNames []string) {
	Exec("SET collation_connection = 'utf8_general_ci'")
	Exec(fmt.Sprintf("ALTER DATABASE %s CHARACTER SET utf8 COLLATE utf8_general_ci", configer.GetString("db_name")))
	for _, name := range tableNames {
		alterTableUtf8(name)
	}
}

func alterTableUtf8(tbName string) {
	Exec(fmt.Sprintf("ALTER TABLE %s CONVERT TO CHARACTER SET utf8 COLLATE utf8_general_ci", tbName))
}
