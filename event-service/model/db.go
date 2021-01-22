package model

import (
	"scala-disaster-adviser/event-service/database"
)

var DB = database.Database{}

func DropAll() {
	DB.Instance.Query("DELETE FROM users")
	DB.Instance.Query("DELETE FROM events")
}
