package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/artalkjs/artalk/v2/internal/core"
	"github.com/artalkjs/artalk/v2/internal/entity"
	"github.com/artalkjs/artalk/v2/internal/log"
	"github.com/artalkjs/artalk/v2/internal/sync"
	"github.com/artalkjs/artalk/v2/server/common"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type RequestAuthDataMergeApply struct {
	UserName string `json:"user_name" validate:"required"`
}

type ResponseAuthDataMergeApply struct {
	UpdatedComment int64  `json:"update_comments_count"`
	UpdatedNotify  int64  `json:"update_notifies_count"`
	UpdatedVote    int64  `json:"update_votes_count"`
	DeletedUser    int64  `json:"deleted_user_count"`
	UserToken      string `json:"user_token"` // Empty if login user is target user no need to re-login
}

// @Id           ApplyDataMerge
// @Summary      Apply data merge
// @Description  This function is to solve the problem of multiple users with the same email address, should be called after user login and then check, and perform data merge.
// @Tags         Auth
// @Security     ApiKeyAuth
// @Param        data  body  RequestAuthDataMergeApply  true  "The data"
// @Success      200  {object}  ResponseAuthDataMergeApply
// @Failure      400  {object}  Map{msg=string}
// @Failure      500  {object}  Map{msg=string}
// @Accept       json
// @Produce      json
// @Router       /auth/merge  [post]
func AuthMergeApply(app *core.App, router fiber.Router) {
	mutexMap := sync.NewKeyMutex[uint]()

	router.Post("/auth/merge", common.LoginGuard(app, func(c *fiber.Ctx, user entity.User) error {
		// Mutex for each user to avoid concurrent merge operation
		mutex := mutexMap.GetLock(user.ID)
		mutex.Lock()
		defer mutex.Unlock()

		if user.Email == "" {
			return common.RespError(c, 500, "User email is empty")
		}

		var p RequestAuthDataMergeApply
		if isOK, resp := common.ParamsDecode(c, &p); !isOK {
			return resp
		}

		// Get all users with same email
		sameEmailUsers := app.Dao().FindUsersByEmail(user.Email)
		if len(sameEmailUsers) == 0 {
			return common.RespError(c, 500, "No user with same email")
		}

		targetUser, err := app.Dao().FindCreateUser(p.UserName, user.Email, user.Link)
		if err != nil {
			return common.RespError(c, 500, "Failed to create user")
		}

		// Check target if admin and recover
		isAdmin := false
		for _, u := range sameEmailUsers {
			if u.IsAdmin {
				isAdmin = true
				break
			}
		}
		if targetUser.IsAdmin != isAdmin {
			targetUser.IsAdmin = isAdmin
			app.Dao().UpdateUser(&targetUser)
		}

		resp := ResponseAuthDataMergeApply{}
		otherUsers := lo.Filter(sameEmailUsers, func(u entity.User, _ int) bool {
			return u.ID != targetUser.ID
		})

		// Functions for log
		getMergeLogSummary := func() string {
			getUserInfo := func(u entity.User) string {
				return fmt.Sprintf("[%d, %s, %s]", u.ID, strconv.Quote(u.Name), strconv.Quote(u.Email))
			}
			getUsersInfo := func(otherUsers []entity.User) string {
				return strings.Join(lo.Map(otherUsers, func(u entity.User, _ int) string { return getUserInfo(u) }), ", ")
			}
			return " | " + getUsersInfo(otherUsers) + " -> " + getUserInfo(targetUser)
		}

		// Begin a transaction to Merge all user data to target user
		if err := app.Dao().DB().Transaction(func(tx *gorm.DB) error {
			// Merge all user data to target user
			// @note Search code files under `./internal/entity` keyword 'UserID` to find all related tables
			for _, u := range otherUsers {
				// comments
				if tx := app.Dao().DB().Model(&entity.Comment{}).
					Where("user_id = ?", u.ID).Update("user_id", targetUser.ID); tx.Error != nil {
					return tx.Error // if error the whole transaction will be rollback
				} else {
					resp.UpdatedComment += tx.RowsAffected
				}

				// notifies
				if tx := app.Dao().DB().Model(&entity.Notify{}).
					Where("user_id = ?", u.ID).Update("user_id", targetUser.ID); tx.Error != nil {
					return tx.Error
				} else {
					resp.UpdatedNotify += tx.RowsAffected
				}

				// votes
				if tx := app.Dao().DB().Model(&entity.Vote{}).
					Where("user_id = ?", u.ID).Update("user_id", targetUser.ID); tx.Error != nil {
					return tx.Error
				} else {
					resp.UpdatedVote += tx.RowsAffected
				}

				// auth_identities
				if tx := app.Dao().DB().Model(&entity.AuthIdentity{}).
					Where("user_id = ?", u.ID).Update("user_id", targetUser.ID); tx.Error != nil {
					return tx.Error
				}
			}

			return nil
		}); err != nil {
			log.Error("Failed to merge user data: ", err.Error(), getMergeLogSummary())
			return common.RespError(c, 500, "Failed to merge data")
		}

		// Delete other users except target user
		for _, u := range otherUsers {
			if err := app.Dao().DelUser(&u); err != nil {
				log.Error("Failed to delete other user [id=", u.ID, "]: ", err.Error(), getMergeLogSummary())
			} else {
				resp.DeletedUser++
			}
		}

		// Re-login
		jwtToken, err := common.LoginGetUserToken(targetUser, app.Conf().AppKey, app.Conf().LoginTimeout)
		if err != nil {
			return common.RespError(c, 500, "Failed to re-login")
		}
		resp.UserToken = jwtToken

		// Log
		log.Info("User data merged successfully", getMergeLogSummary(), " | ",
			"Updated Comments: ", resp.UpdatedComment, " | ",
			"Updated Notifies: ", resp.UpdatedNotify, " | ",
			"Updated Votes: ", resp.UpdatedVote, " | ",
			"Deleted Users: ", resp.DeletedUser)

		return common.RespData(c, resp)
	}))
}
