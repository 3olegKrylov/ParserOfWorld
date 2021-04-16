package internal

import (
	"database/sql"
	"fmt"
	"github.com/testSpace/internal/db"
	"log"
	"strings"
)

func SendData(text string, dbConnect *sql.DB) {
	lines := strings.Split(text, "\n\n")
	DataOfUsers := make(map[string]string)

	for num, value := range lines {
		if strings.HasSuffix(value, "Подписчики") {
			DataOfUsers[lines[num-1]] = value
		}
	}
	id := int32(0)

	for name, comment := range DataOfUsers {
		fmt.Println("Name:", name, "\nValue:", comment, "\n")
		db.DBAddUser(id, name, comment, "", dbConnect)
		id++

	}

	log.Print("Количество пользователей: ", len(DataOfUsers))
}
