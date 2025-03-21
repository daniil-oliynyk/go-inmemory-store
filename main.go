package main

import (
	"fmt"
	"os"
)

func main() {

	knockoff_redis_server := Tcp_Server{port: ":8080"}
	err := knockoff_redis_server.Run()
	if err != nil {
		fmt.Println("Error in Server ", err)
		os.Exit(1)
	}
}
