package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	models "github.com/tjw0051/log-go/Models"
	store "github.com/tjw0051/log-go/Store"
)

type ErrorResponse struct {
	Success  bool   `json:"success" binding:"required"`
	ErrorMsg string `json:"errorMsg" binding:"required"`
	Debug    string `json:"debug" binding:"required"`
}

/*	Convenience Functions 	*/
func ReturnError(c *gin.Context, statusCode int, success bool, message string, errorMsg string) {
	c.JSON(statusCode, ErrorResponse{Success: success, ErrorMsg: message, Debug: errorMsg})
}

func returnBadJson(c *gin.Context, err error) {
	ReturnError(c, http.StatusBadRequest, false, "Malformed json.", err.Error())
}

func CreateLog(c *gin.Context) {
	var json models.MessagesModel
	err := c.BindJSON(&json)
	if err != nil {
		returnBadJson(c, err)
		return
	}

	err = store.CreateLog(json)
	if err != nil {
		ReturnError(c, 500, false, "Error creating message", err.Error())
	}
	c.Status(201)
}

func Query(c *gin.Context) {
	queryMap := models.MessageModel{
		Platform:  c.Query("platform"),
		Project:   c.Query("project"),
		Component: c.Query("component"),
		LogType:   c.Query("logtype"),
		Severity:  c.Query("severity"),
		UserID:    c.Query("userid"),
		Message:   c.Query("message"),
	}
	var err error
	count, _ := strconv.Atoi(c.Query("count"))
	page, _ := strconv.Atoi(c.Query("page"))

	messages, err := store.QueryLog(queryMap, count, page)
	if err != nil {
		ReturnError(c, 500, false, "Error deleting keys", err.Error())
		return
	}
	c.JSON(200, messages)
}

func QueryPlatforms(c *gin.Context) {

}
