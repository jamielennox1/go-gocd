package main

import (
	"fmt"
	"github.com/mhanygin/go-gocd/gocd"
)

func main() {
	client := gocd.New("https://go.inn.ru", "sa_dev_go_bot", "7(y3(65#cN*86szT")
	fmt.Print(client.GetPipeline("broforce", 1))
}
