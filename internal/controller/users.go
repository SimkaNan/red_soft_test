package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golibrary/internal/model"
	"net/http"
	"reflect"
	"strconv"
)

type UserBook struct {
	UserId string `form:"userId" json:"userId"`
	BookId string `form:"bookId" json:"bookId"`
}

// @Summary ListUsers
// @Tags users
// @Description listUsers
// @ID list-users
// @Produce json
// @Success 200 {object} []model.User "Users list"
// @Failure 500 {object} Error "Internal error"
// @Router /users/list [get]
func (c *Controller) ListUsers(ctx *gin.Context) {
	users, err := c.service.User.ListUsers(ctx.Request.Context())
	if err != nil {
		c.logger.Error("ListUsers error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"usersList": users})
}

// @Summary GetUserByID
// @Tags users
// @Description getUsersByID
// @ID get-user-Byid
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} model.User "Get user data"
// @Failure 400 {object} Error "Bad request"
// @Failure 500 {object} Error "Internal error"
// @Router /users/getByID [get]
func (c *Controller) GetUserByID(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		c.logger.Error("GetUserByID error", zap.Error(errors.New("UserID parameter is required")))
		NewErrorResponse(ctx, http.StatusBadRequest, "UserID parameter is required")
		return
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		c.logger.Error("GetUserByID error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.service.User.GetUserByID(ctx.Request.Context(), userId)
	if err != nil {
		c.logger.Error("GetUserByID error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": user})
}

// @Summary GetUserBySurname
// @Tags users
// @Description getUsersBySurname
// @ID get-user-by-surname
// @Produce json
// @Param surname query string true "surname"
// @Success 200 {object} model.User "Get user data"
// @Failure 400 {object} Error "Bad request"
// @Failure 500 {object} Error "Internal error"
// @Router /users/getBySurname [get]
func (c *Controller) GetUserBySurname(ctx *gin.Context) {
	surname := ctx.Query("surname")
	if surname == "" {
		c.logger.Error("GetUserBySurname error", zap.Error(errors.New("Surname parameter is required")))
		NewErrorResponse(ctx, http.StatusBadRequest, "Surname parameter is required")
		return
	}

	user, err := c.service.User.GetUserBySurname(ctx.Request.Context(), surname)
	if err != nil {
		c.logger.Error("GetUserBySurname error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": user})
}

// @Summary CreateUser
// @Tags users
// @Description CreateUser
// @ID create-user
// @Produce json
// @Accept json
// @Param input body model.User true "User data"
// @Success 200 {string} string "User successfully created"
// @@Failure 400 {object} Error "Bad request"
// @Failure 500 {object} Error "Internal error"
// @Router /users/create [post]
func (c *Controller) CreateUser(ctx *gin.Context) {
	var user model.User

	err := ctx.BindJSON(&user)
	if err != nil {
		c.logger.Error("CreateUser error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := c.service.User.CreateUser(ctx.Request.Context(), &user)
	if err != nil {
		c.logger.Error("CreateUser error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User successfully created with id: " + strconv.Itoa(id)})
}

// @Summary UpdateUser
// @Tags users
// @Description UpdateUser
// @ID update-user
// @Produce json
// @Accept json
// @Param id query string true "id"
// @Param input body model.UserUpdate true "Data for Update user"
// @Success 200 {string} string "User successfully updated"
// @@Failure 400 {object} Error "Bad request"
// @Failure 500 {object} Error "Internal error"
// @Router /users/update [put]
func (c *Controller) UpdateUser(ctx *gin.Context) {
	var user model.UserUpdate

	id := ctx.Query("id")
	if id == "" {
		c.logger.Error("UpdateUser error", zap.Error(errors.New("UserID parameter is required")))
		NewErrorResponse(ctx, http.StatusBadRequest, "UserID parameter is required")
		return
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		c.logger.Error("UpdateUser error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = ctx.BindJSON(&user)
	if err != nil {
		c.logger.Error("UpdateUser error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if reflect.DeepEqual(user, model.UserUpdate{}) {
		c.logger.Error("UpdateUser error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusBadRequest, "At least one field must be provided")
		return
	}

	err = c.service.User.UpdateUser(ctx.Request.Context(), userId, &user)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User with ID: " + id + "successfully updated."})
}

// @Summary ListUserFriends
// @Tags friendships
// @Description listUserFriends
// @ID list-users-friends
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} []model.User "User friends list"
// @Failure 500 {object} Error "Internal error"
// @Router /friendships/list [get]
func (c *Controller) ListUserFriends(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		c.logger.Error("ListUserFriends error", zap.Error(errors.New("UserID parameter is required")))
		NewErrorResponse(ctx, http.StatusBadRequest, "UserID parameter is required")
		return
	}

	userId, err := strconv.Atoi(id)
	if err != nil {
		c.logger.Error("ListUserFriends error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	friendships, err := c.service.User.ListFriendships(ctx.Request.Context(), userId)
	if err != nil {
		c.logger.Error("ListUserFriends error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"friendshipsList": friendships})
}

type Friendship struct {
	FirstID  int `json:"first_id" db:"user_id1" binding:"required"`
	SecondID int `json:"second_id" db:"user_id2" binding:"required"`
}

// @Summary CreateUsersFriendships
// @Tags friendships
// @Description CreateUsersFriendships
// @ID create-user-friendships
// @Produce json
// @Accept json
// @Param input body Friendship true "Data for create friendships"
// @Success 200 {string} string "User friendships successfully created"
// @@Failure 400 {object} Error "Bad request"
// @Failure 500 {object} Error "Internal error"
// @Router /friendships/create [post]
func (c *Controller) CreateFriendship(ctx *gin.Context) {
	var friendship Friendship

	err := ctx.BindJSON(&friendship)
	if err != nil {
		c.logger.Error("CreateFriendship error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	err = c.service.User.CreateFriendship(ctx.Request.Context(), friendship.FirstID, friendship.SecondID)
	if err != nil {
		c.logger.Error("CreateFriendship error", zap.Error(err))
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Friendship successfully created"})
}
