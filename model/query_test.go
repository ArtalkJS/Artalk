package model

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/ArtalkJS/ArtalkGo/lib"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mock sqlmock.Sqlmock
var mockDB *gorm.DB

// TestMain 初始化 Mock
func TestMain(m *testing.M) {
	var db *sql.DB
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		panic(err)
	}
	lib.DB, err = gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	mockDB = lib.DB
	defer db.Close()
	m.Run()
}

func TestFindComment(t *testing.T) {
	asserts := assert.New(t)

	// 创建评论
	newComment := Comment{
		Model: gorm.Model{
			ID: 1,
		},
		Content:  "testing content...",
		PageKey:  "/artalk-demo/1.html",
		SiteName: "ArtalkDemo",
		UserID:   1,
		Rid:      0,
	}
	lib.DB.Save(&newComment)

	// 预期数据
	rows := sqlmock.NewRows([]string{"id", "deleted_at", "content", "page_key", "site_name", "user_id", "rid"}).
		AddRow(1, nil, "testing content...", "/artalk-demo/1.html", "ArtalkDemo", 1, 0)
	mock.ExpectQuery("SELECT(.+)").WillReturnRows(rows)

	// 查询评论
	comment := FindComment(1)
	// asserts.NoError(err)

	// 找到时
	asserts.Equal(newComment, comment)
	asserts.Equal(comment.IsEmpty(), false)

	// 未找到时
	mock.ExpectQuery("^SELECT (.+)").WillReturnError(errors.New("not found"))
	comment = FindComment(1)
	// asserts.Error(err)
	asserts.Equal(comment.IsEmpty(), true)
	asserts.Equal(Comment{}, comment)
}
