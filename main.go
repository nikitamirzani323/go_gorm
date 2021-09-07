package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	// GetGenderByName("bencong")
	// UpdateGender2(4, "Dracula")
	// DeleteGender(4)
	// CreateTest(0, "Test1")
	// CreateTest(0, "Test2")
	// CreateTest(0, "Test3")
	// GetTests()
	// DeleteTest(2)
	// db.Migrator().CreateTable(Customer{})
	// CreateCustomer("Note", 1)
	// GetCustomers()
	UpdateGender2(1, "asdasdasd")

}

type Customer struct {
	ID       uint
	Name     string
	Gender   Gender
	GenderID uint
}

func GetCustomers() {
	customers := []Customer{}
	tx := db.Preload(clause.Associations).Find(&customers)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	for _, customer := range customers {
		fmt.Printf("%v|%v|%v\n", customer.ID, customer.Name, customer.Gender.Name)
	}
}

func CreateCustomer(name string, genderID uint) {
	customer := Customer{Name: name, GenderID: genderID}
	tx := db.Create(&customer)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customer)
}
func CreateTest(code uint, name string) {
	test := Test{Code: code, Name: name}
	db.Create(&test)
}
func GetTests() {
	tests := []Test{}
	db.Find(&tests)
	for _, t := range tests {
		fmt.Printf("%v|%v\n", t.ID, t.Name)
	}
}
func DeleteTest(id uint) {
	db.Unscoped().Delete(&Test{}, id)
}
func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("Deleted")
	GetGender(id)
}
func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	gender.Name = name
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(id)
}
func UpdateGender2(id uint, name string) {
	gender := Gender{Name: name}
	tx := db.Model(&Gender{}).Where("id=@myid", sql.Named("myid", id)).Updates(gender)

	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(id)
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
	// tx := db.First(&genders, "name=?", name)
	tx := db.Where("name=?", name).First(&genders)
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

func (g Gender) BeforeUpdate(*gorm.DB) error {
	fmt.Printf("Before Update Gender: %v => %v\n", g.ID, g.Name)
	return nil
}

func (g Gender) AfterUpdate(*gorm.DB) error {
	fmt.Printf("After Update Gender: %v => %v\n", g.ID, g.Name)
	return nil
}

type Test struct {
	gorm.Model
	Code uint   `gorm:"comment:This is Code"`
	Name string `gorm:"column:myname;size:20;unique:default:Hello;not null"`
}
