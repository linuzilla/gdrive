package models

import (
	"database/sql"
	"time"
)

type SyncFileInfo struct {
	Id         string         `gorm:"column:id;primary_key"`
	CheckSum   string         `gorm:"column:checksum;index"`
	FileId     sql.NullString `gorm:"column:file_id;index" sql:"null"`
	RemoteName sql.NullString `gorm:"column:remote_name" sql:"null"`
	FileSize   int64
	Uploaded   bool
	ModTime    time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (SyncFileInfo) TableName() string {
	return "files"
}
