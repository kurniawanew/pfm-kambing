package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kurniawanew/pfm-kambing/api/auth"
	"github.com/kurniawanew/pfm-kambing/api/models"
	"github.com/kurniawanew/pfm-kambing/api/utils/formaterror"
)

func (server *Server) CreateUser(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	userCreated, err := user.SaveUser(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		c.JSON(http.StatusInternalServerError, formattedError)
		return
	}
	c.Header("Location", fmt.Sprintf("%s%s/%d", c.Request.Host, c.Request.RequestURI, userCreated.ID))
	c.JSON(http.StatusCreated, userCreated)
}

func (server *Server) GetUsers(c *gin.Context) {

	user := models.User{}

	users, err := user.FindAllUsers(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)

		return
	}
	c.JSON(http.StatusOK, users)
}

func (server *Server) GetUser(c *gin.Context) {

	vars := c.Param("id")
	uid, err := strconv.ParseUint(vars, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	user := models.User{}
	userGotten, err := user.FindUserByID(server.DB, uint32(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, userGotten)
}

func (server *Server) UpdateUser(c *gin.Context) {

	vars := c.Param("id")
	uid, err := strconv.ParseUint(vars, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != uint32(uid) {
		c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	user.Prepare()
	err = user.Validate("update")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	updatedUser, err := user.UpdateAUser(server.DB, uint32(uid))
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, formattedError)
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (server *Server) DeleteUser(c *gin.Context) {

	vars := c.Param("id")

	user := models.User{}

	uid, err := strconv.ParseUint(vars, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if tokenID != 0 && tokenID != uint32(uid) {
		c.JSON(http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	_, err = user.DeleteAUser(server.DB, uint32(uid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.Header("Entity", fmt.Sprintf("%d", uid))
	c.JSON(http.StatusNoContent, "")
}
