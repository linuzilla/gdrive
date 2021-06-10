package models

import (
	"database/sql"
	"time"
)

type GoogleDriveConfig struct {
	Id            string         `gorm:"column:id;primary_key"`
	Name          sql.NullString `gorm:"column:name" sql:"null"`
	Password      sql.NullString `gorm:"column:password" sql:"null"`
	FolderId      string         `gorm:"column:folder_id;index" sql:"null"`
	Encoder       sql.NullString `gorm:"column:encoder" sql:"null"`
	Decoder       sql.NullString `gorm:"column:decoder" sql:"null"`
	TrashFolderId sql.NullString `gorm:"column:trash_id" sql:"null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (GoogleDriveConfig) TableName() string {
	return "config"
}
