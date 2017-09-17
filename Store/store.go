package store

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	models "github.com/tjw0051/log-go/Models"
)

var db *gorm.DB

func Connect(host string, port string, user string, dbName string, password string) error {
	var err error
	dbUrl := "host=" + host + " port=" + port + " user=" + user + " dbname=" + dbName + " sslmode=disable password=" + password
	println(dbUrl)
	db, err = gorm.Open("postgres", dbUrl)
	if err != nil && db == nil {
		return err
	}
	err = db.DB().Ping()
	if err != nil {
		AutoReconnect(host, port, user, dbName, password)
	}

	db.AutoMigrate(&models.MessageModel{}, &models.KeyModel{})
	return nil
}

func AutoReconnect(host string, port string, user string, dbName string, password string) {
	print("Error connecting to DB. attemping reconnect in 10 seconds...")
	timer1 := time.NewTimer(time.Second * 10)

	<-timer1.C
	var err error
	db, err = gorm.Open("postgres", "host="+host+" port="+port+" user="+user+" dbname="+dbName+" sslmode=disable password="+password)
	err = db.DB().Ping()
	if err != nil {
		AutoReconnect(host, port, user, dbName, password)
	} else {
		print("Successfully reconnected.")
	}
}

func CreateLog(log models.MessagesModel) error {
	for i := 0; i < len(log); i++ {
		err := db.Create(&log[i]).Error
		return err
	}
	return nil
}

func QueryLog(queryMap models.MessageModel, count int, page int) (models.QueryModel, error) {
	var messages models.MessagesModel

	query := db

	if queryMap.Platform != "" {
		println("platform: " + queryMap.Platform)
		//searchMap["platform"] = queryMap.Platform
		query = db.Where("platform = ?", queryMap.Platform)
	}
	if queryMap.Project != "" {
		println("project: " + queryMap.Project)
		//searchMap["project"] = queryMap.Project
		query = query.Where("project = ?", queryMap.Project)
	}
	if queryMap.Component != "" {
		println("component: " + queryMap.Component)
		//searchMap["component"] = queryMap.Component
		query = query.Where("component = ?", queryMap.Component)
	}
	if queryMap.LogType != "" {
		println("logtype: " + queryMap.LogType)
		//searchMap["logtype"] = queryMap.LogType
		query = query.Where("log_type = ?", queryMap.LogType)
	}
	if queryMap.Severity != "" {
		println("severity: " + queryMap.Severity)
		//searchMap["severity"] = queryMap.Severity
		query = query.Where("severity = ?", queryMap.Severity)
	}
	if queryMap.UserID != "" {
		println("userid: " + queryMap.UserID)
		//searchMap["userid"] = queryMap.UserID
		query = query.Where("user_id = ?", queryMap.UserID)
	}
	if queryMap.Message != "" {
		println("message: " + queryMap.Message)
		query = query.Where("message LIKE ?", "%"+queryMap.Message+"%")
	}

	// Set count
	validCount := count
	if count == 0 {
		validCount = 100
	}
	query = query.Limit(validCount)

	println("ValidCount: ", validCount, " Offset: ", page*validCount)
	query = query.Offset(page * validCount)

	// Find
	err := query.Find(&messages).Error

	queryResults := models.QueryModel{
		Messages: messages,
		Page:     page,
	}

	return queryResults, err

	/*
		queryMessage := queryMap.Message
		println("message: " + queryMessage)

		var err error
		query := db.Where(searchMap)

		if queryMessage != "" {
			err = query.Where("message LIKE ?", queryMessage).Find(&messages).Error
		} else {
			err = query.Find(&messages).Error
		}
		return messages, err
	*/
}
