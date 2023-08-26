package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"path/filepath"
	"time"

	inkassback "github.com/Husenjon/InkassBack"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const maxFileSize = int64(20 * 1024 * 1024) // 10 MB
func (h *Handler) signUp(c *gin.Context) {
	if c.Request.ContentLength > maxFileSize {
		NewErrorResponse(c, http.StatusBadRequest, "File size exceeds limit")
		return
	}
	var user inkassback.User
	form, err := c.MultipartForm()
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	body := form.Value["body"]
	if body == nil {
		NewErrorResponse(c, http.StatusBadRequest, "body empty")
		return
	}
	if err := json.Unmarshal([]byte(body[0]), &user); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		for _, err := range err.(validator.ValidationErrors) {
			NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
	}
	image := form.File["image"]
	if image != nil {
		dir, err := filepath.Abs("./images/user")
		if err != nil {
			fmt.Println(err)
		}
		hashToken := h.services.Authoration.GenerateHash(fmt.Sprint(time.Now().Unix() + int64(rand.Intn(100))))
		user.Image = hashToken
		filename := fmt.Sprintf("%s/%s", dir, hashToken)
		go func(image map[string][]*FileHeader, name string) {
			if err := c.SaveUploadedFile(image, name); err != nil {
				NewErrorResponse(c, http.StatusBadRequest, err.Error())
				return
			}
		}(image, filename)
	} else {
		user.Image = "image"
	}
	_user, err := h.services.CreateUser(user)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":           _user.Id,
		"ism":          _user.Ism,
		"familya":      _user.Familya,
		"otasini_ismi": _user.OtasiniIsmi,
		"phone":        _user.Phone,
		"username":     _user.Username,
		"branch_id":    _user.BranchId,
		"image":        _user.Image,
		"is_active":    _user.IsActive,
		"created_time": _user.CreatedTime,
	})
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.BindJSON(&input); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	user, err := h.services.Authoration.GetUser(input.Username, input.Password)
	if err != nil {
		NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.SetCookie("session", user.Token.String, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id":           user.Id,
		"ism":          user.Ism,
		"familya":      user.Familya,
		"otasini_ismi": user.OtasiniIsmi,
		"phone":        user.Phone,
		"username":     user.Username,
		"branch_id":    user.BranchId,
		"image":        user.Image,
		"is_active":    user.IsActive,
		"created_time": user.CreatedTime,
	})
	/*filePath := "images/user/2.mp4"
	c.Header("Content-Disposition", "attachment; filename=downloaded_file.txt")
	c.File(filePath)*/
}
