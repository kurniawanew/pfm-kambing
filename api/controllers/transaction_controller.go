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

func (server *Server) CreateTransaction(c *gin.Context) {

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
	}
	transaction := models.Transaction{}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	tokenID, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	transaction.Prepare(tokenID)
	err = transaction.Validate()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}
	transactionCreated, err := transaction.SaveTransaction(server.DB)

	if err != nil {

		formattedError := formaterror.FormatError(err.Error())

		c.JSON(http.StatusInternalServerError, formattedError)
		return
	}
	c.Header("Location", fmt.Sprintf("%s%s/%d", c.Request.Host, c.Request.RequestURI, transactionCreated.ID))
	c.JSON(http.StatusCreated, transactionCreated)
}

func (server *Server) GetTransactions(c *gin.Context) {

	transaction := models.Transaction{}

	transactions, err := transaction.FindAllTransactions(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)

		return
	}
	c.JSON(http.StatusOK, transactions)
}

func (server *Server) GetTransaction(c *gin.Context) {

	vars := c.Param("id")
	uid, err := strconv.ParseUint(vars, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	transaction := models.Transaction{}
	transactionGotten, err := transaction.FindTransactionByID(server.DB, uint32(uid))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, transactionGotten)
}

func (server *Server) UpdateTransaction(c *gin.Context) {

	vars := c.Param("id")

	// Check if the post id is valid
	tid, err := strconv.ParseUint(vars, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	transaction := models.Transaction{}
	err = server.DB.Debug().Model(models.Transaction{}).Where("id = ?", tid).Take(&transaction).Error
	if err != nil {
		c.JSON(http.StatusNotFound, errors.New("Transaction not found"))
		return
	}

	// If a user attempt to update a transaction not belonging to him
	if uid != transaction.TransactionUserID {
		c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	transactionUpdate := models.Transaction{}
	err = json.Unmarshal(body, &transactionUpdate)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	//Also check if the request user id is equal to the one gotten from token
	if uid != transactionUpdate.TransactionUserID {
		c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	transactionUpdate.Prepare(uid)
	err = transactionUpdate.Validate()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err)
		return
	}

	transactionUpdate.ID = transaction.ID

	transactionUpdated, err := transactionUpdate.UpdateATransaction(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, formattedError)
		return
	}
	c.JSON(http.StatusOK, transactionUpdated)
}

func (server *Server) DeleteTransaction(c *gin.Context) {

	vars := c.Param("id")

	// Is a valid post id given to us?
	id, err := strconv.ParseUint(vars, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the post exist
	transaction := models.Transaction{}
	err = server.DB.Debug().Model(models.Transaction{}).Where("id = ?", id).Take(&transaction).Error
	if err != nil {
		c.JSON(http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this post?
	if uid != transaction.TransactionUserID {
		c.JSON(http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = transaction.DeleteATransaction(server.DB, id, uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.Header("Entity", fmt.Sprintf("%d", id))
	c.JSON(http.StatusNoContent, "")
}
