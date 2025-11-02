package db

import (
	"log"
	"os"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Bucket struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type Object struct {
	ID         uint   `gorm:"primaryKey"`
	BucketID   uint   `gorm:"index"`
	Key        string
	Size       int64
	ETag       string
	VersionID  string `gorm:"default:(substr(md5(randomblob(8)),1,12))"`
	Metadata   string `gorm:"type:json"`
	IsLatest   bool   `gorm:"default:true"`
	CreatedAt  int64  `gorm:"autoCreateTime"`
}

type Policy struct {
	ID       uint   `gorm:"primaryKey"`
	Subject  string
	Resource string
	Action   string
	Effect   string `gorm:"check:effect IN ('allow','deny')"`
}

func InitDB(file string) {
	os.MkdirAll("./data", 0755)
	var err error
	DB, err = gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil { log.Fatal(err) }
	DB.AutoMigrate(&Bucket{}, &Object{}, &Policy{})
	DB.FirstOrCreate(&Policy{Subject: "admin", Resource: ".*", Action: ".*", Effect: "allow"})
}
