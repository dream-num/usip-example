// file: datasource/users.go

package datasource

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadDB() (db *gorm.DB, err error) {
	switch viper.GetString("database.driver") {
	case "postgresql":
		db, err = gorm.Open(postgres.Open(viper.GetString("database.dsn")))
	case "mysql":
		db, err = gorm.Open(mysql.Open(viper.GetString("database.dsn")))
	default:
		panic(fmt.Sprintf("Unsupported database driver: %s", viper.GetString("database.driver")))
	}

	return
}
