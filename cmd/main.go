package main

import (
	"bookService/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)
}
