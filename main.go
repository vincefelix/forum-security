package main

import (
	"fmt"

	db "forum/Database"
	hdle "forum/Handlers"
)

func main() {
	tab, err := db.Init_db()
	if err != nil {
		fmt.Println(err)
		return
	}
	hdle.Handlers(tab)
}
