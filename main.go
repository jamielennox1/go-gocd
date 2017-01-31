package main

import (
	"fmt"
	"github.com/mhanygin/go-gocd/gocd"
)

func main() {
	client := gocd.New("https://go.inn.ru", "sa_dev_go_bot", "7(y3(65#cN*86szT")
	//fmt.Println(client.GetPipeline("broforce", 1))
	//
	//if pp, err := client.GetPipelineHistory("broforce"); err == nil {
	//	for _, p := range pp.Pipelines {
	//		fmt.Println(p)
	//	}
	//} else {
	//	fmt.Println(err)
	//}
	//
	//if envs, err := client.GetEnvironments(); err == nil {
	//	for _, env := range envs.Embeded.Environments {
	//		fmt.Println(env)
	//	}
	//} else {
	//	fmt.Println(err)
	//}

	//if env, err := client.GetEnvironment("DEV"); err == nil {
	//	fmt.Println(env)
	//	client.SetEnvironment(env)
	//} else {
	//	fmt.Println(err)
	//}

	if p, err := client.GetPipelineConfig("broforce"); err == nil {
		fmt.Println(p)
	} else {
		fmt.Println(err)
	}

}
