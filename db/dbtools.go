package db

import (
	"fmt"
)

func InitDBase() {
	fmt.Println("Init Data Base...")
	/*
		appPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
		_, err = os.Stat(dbFile)

		var install bool
		if err != nil {
			install = true
		}
		// если install равен true, после открытия БД требуется выполнить
		// sql-запрос с CREATE TABLE и CREATE INDEX
		if install {
			fmt.Print("install true")
		} else {
			fmt.Print("install true")
		}
	*/
}
