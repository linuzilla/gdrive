package test_code

import (
	"database/sql"
	"fmt"
	"github.com/linuzilla/gdrive/models"
	"testing"
	"time"
)

func getOne() *models.GoogleDriveConfig {
	return &models.GoogleDriveConfig{
		Id:            "A",
		Password:      sql.NullString{String: "abc", Valid: true},
		FolderId:      "folderId",
		Encoder:       sql.NullString{},
		Decoder:       sql.NullString{},
		TrashFolderId: sql.NullString{},
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}

}

func TestCopy(t *testing.T) {
	str := `2020-01-27T07:01:05.945Z`
	fmt.Println(str)
	if parse, err := time.Parse(time.RFC3339, str); err != nil {
		fmt.Println(err)
	} else {
		in := parse.In(time.Now().Location())
		fmt.Println(in)
	}

	var slice []string

	slice = append(slice, "a")
	slice = append(slice, "b")
	slice = append(slice, "c")
	slice = append(slice, "d")

	fmt.Println(slice)
	for i, item := range slice[1:] {
		fmt.Printf("%d. [%s]\n", i, item)
	}

	//fmt.Println("testing")
	//
	//var stored models.GoogleDriveConfig
	//pointer := getOne()
	//stored = *pointer
	//
	//stored.Password = sql.NullString{String:"abc", Valid:true}
	//stored.FolderId = fmt.Sprintf("%sder%s","fol", "Id")
	//
	//if stored != *pointer {
	//	t.Error("should be equal")
	//}
	//
	//stored.Id = "B"
	//
	//if stored == *pointer {
	//	t.Error("should not be equal")
	//}
	//
	//fmt.Println(utils.Detail(&stored))
	//fmt.Println(utils.Detail(pointer))
}
