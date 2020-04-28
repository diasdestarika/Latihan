package main

import (
	"fmt"

	"go-tutorial-2020/internal/boot"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if err := boot.HTTP(); err != nil {
		fmt.Println("[HTTP] failed to boot http server due to " + err.Error())
	}
}
