package handlers

import (
	"github.com/gin-gonic/gin"
	models "github.com/tjw0051/log-go/Models"
	store "github.com/tjw0051/log-go/Store"
)

func CreateKeys(c *gin.Context) {
	var json models.KeysModel
	err := c.BindJSON(&json)
	if err != nil {
		returnBadJson(c, err)
		return
	}

	err = store.CreateKeys(json)
	if err != nil {
		ReturnError(c, 500, false, "Error creating key", err.Error())
		return
	}
	c.Status(201)
}

func GetKeys(c *gin.Context) {
	keys, err := store.GetKeys()
	if err != nil {
		ReturnError(c, 500, false, "Error getting keys", err.Error())
		return
	}
	c.JSON(200, keys)
}

func DeleteKeys(c *gin.Context) {
	var json models.KeysModel
	err := c.BindJSON(&json)
	if err != nil {
		returnBadJson(c, err)
		return
	}

	err = store.DeleteKeys(json)
	if err != nil {
		ReturnError(c, 500, false, "Error deleting keys", err.Error())
		return
	}
	c.Status(200)
}
