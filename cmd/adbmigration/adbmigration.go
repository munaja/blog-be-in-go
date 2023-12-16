package main

import (
	"fmt"
	"log"

	m "github.com/munaja/blog-practice-be-using-go/internal/migration"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DbConf struct {
	Dsn     string
	Dialect string
}

var model = ""

// load from yml files to struct config
func (c *DbConf) loadConfig() (err error) {
	// main viper config
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// read the file
	if err = viper.ReadInConfig(); err != nil {
		err = fmt.Errorf("error reading config file: %s", err)
		return
	}

	// map to app
	if err = viper.Unmarshal(c); err != nil {
		err = fmt.Errorf("unable to decode into struct: %s", err)
		return
	}

	// done
	log.Printf("config loaded successfully")
	return
}

func InitDB(input DbConf) (*gorm.DB, error) {
	var gormD gorm.Dialector
	if input.Dialect == "mysql" {
		gormD = mysql.Open(input.Dsn)
	} else if input.Dialect == "postgres" {
		gormD = postgres.Open(input.Dsn)
	}

	DB, err := gorm.Open(gormD, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})

	if err != nil {
		return nil, err
	}
	return DB, nil
}

func main() {
	var dbConf DbConf

	err := dbConf.loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := InitDB(dbConf)
	if err != nil {
		log.Fatal(err)
	}

	var modelList []interface{}
	modelList = m.GetModelList()
	db.AutoMigrate(modelList...)
	log.Printf("migrate complete")
}
