package main

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface
}

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n========================\n", sql)
}

var db *gorm.DB

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/db_gorm?parseTime=true"
	dial := mysql.Open(dsn)

	var err error
	db, err = gorm.Open(dial, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false, //FUngsinya hanya debug bukan create asli
	})
	if err != nil {
		panic(err)
	}

	// db.Migrator().CreateTable(Test{})
	// db.AutoMigrate(Gender{}, Test{})
	// CreateGender("Monster")
	// GetGenders()
	// GetGender(1)
	GetGenderByName("bencong")
}

func GetGenders() {
	genders := []Gender{}
	tx := db.Order("id").Find(&genders)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}
func GetGender(id uint) {
	genders := []Gender{}
	tx := db.First(&genders, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}
func GetGenderByName(name string) {
	genders := []Gender{}
	tx := db.First(&genders, "name=?", name)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(genders)
}
func CreateGender(name string) {
	gender := Gender{Name: name}
	tx := db.Create(&gender)

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

type Gender struct {
	ID   uint
	Name string `gorm:"unique;size(10)"`
}

type Test struct {
	gorm.Model
	Code string `gorm:"primaryKey;comment:This is Code"`
	Name string `gorm:"column:myname;size:20;unique:default:Hello;not null"`
}
