package store

import (
	"go.uber.org/zap"
	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

var db *gorm.DB

func Db() *gorm.DB {
	var err error
	if db != nil {
		return db
	}
	db, err = gorm.Open(sqlite.Open("raspan.db"), &gorm.Config{})
	if err != nil {
		zap.S().Fatal(err)
	}
	zap.S().Info("successfully connected to db")
	models := []interface{}{
		&Lang{},
		&Category{},
		&Answer{},
		&Question{},
	}
	if err := db.AutoMigrate(models...); err != nil {
		zap.S().Errorf("db migration error: %v", err)
	} else {
		zap.S().Info("db migration successfully completed")
	}
	insertDefaultLangs()
	return db
}

func insertDefaultLangs() {

	if res := Db().FirstOrCreate(&Lang{}, &Lang{LangCode: "he"}); res.Error != nil {
		zap.S().Error(res.Error)
	}
	if res := Db().FirstOrCreate(&Lang{}, &Lang{LangCode: "ru"}); res.Error != nil {
		zap.S().Error(res.Error)
	}

}

type Answer struct {
	gorm.Model
	Answer      string
	Right       bool
	AnswerIndex string
	LangID      uint
	Lang        Lang
	QuestionID  uint
}

type Lang struct {
	gorm.Model
	LangCode string `gorm:"unique"`
}

type Question struct {
	gorm.Model
	Question   string
	QID        string
	LangID     uint
	Lang       Lang
	CategoryID uint
	Category   Category
	Answer     []Answer
}

type Category struct {
	gorm.Model
	Category    string `gorm:"unique"`
	Description string
	LangID      uint
	Lang        Lang
}
