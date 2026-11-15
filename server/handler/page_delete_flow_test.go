package handler_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/artalkjs/artalk/v2/server/handler"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestPageDeletionPreservesOtherSites(t *testing.T) {
	app, router := NewApiTestApp()
	t.Cleanup(func() { require.NoError(t, app.Cleanup()) })
	handler.PageDelete(app.App, router)
	d := app.Dao()
	page := d.NewPage("/delete-page", "Page", "Site A")
	other := d.NewPage(page.Key, "Other", "Site B")
	comment := entity.Comment{SiteName: page.SiteName, PageKey: page.Key, UserID: 1001, Content: "Delete"}
	require.NoError(t, d.CreateComment(&comment))
	preserved := entity.Comment{SiteName: other.SiteName, PageKey: other.Key, UserID: 1001, Content: "Keep"}
	require.NoError(t, d.CreateComment(&preserved))
	notification := d.NewNotify(1002, comment.ID)
	vote, err := d.NewVote(page.ID, entity.VoteTypePageUp, 1001, "", "")
	require.NoError(t, err)
	admin, err := common.LoginGetUserToken(d.FindUserByID(1000), app.Conf().AppKey, 60)
	require.NoError(t, err)
	path := fmt.Sprintf("/pages/%d", page.ID)
	flowRequest(t, router, "DELETE", path, nil, "", 403)
	flowRequest(t, router, "DELETE", "/pages/999999", nil, admin, 404)
	require.NoError(t, d.DB().Callback().Delete().Before("gorm:delete").Register("test:page_delete", func(tx *gorm.DB) { tx.AddError(errors.New("read-only database")) }))
	flowRequest(t, router, "DELETE", path, nil, admin, 500)
	require.NoError(t, d.DB().Callback().Delete().Remove("test:page_delete"))
	require.False(t, d.FindPageByID(page.ID).IsEmpty())
	flowRequest(t, router, "DELETE", path, nil, admin, 200)
	require.True(t, d.FindPageByID(page.ID).IsEmpty())
	require.True(t, d.FindComment(comment.ID).IsEmpty())
	require.False(t, d.FindPageByID(other.ID).IsEmpty())
	require.Equal(t, preserved.Content, d.FindComment(preserved.ID).Content)
	require.ErrorIs(t, d.DB().First(&notification).Error, gorm.ErrRecordNotFound)
	require.ErrorIs(t, d.DB().First(&vote).Error, gorm.ErrRecordNotFound)
}
