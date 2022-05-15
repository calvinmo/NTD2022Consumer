package main

import (
	"fmt"
	"gray.net/app-onedeploy/pkg/core/rabbit"
	asynctask "gray.net/lib-go-async-task-manager"
	"os"
)

func main() {
	rabbitURL := getEnv("RABBIT_URL")
	fmt.Println(rabbitURL)
	rabbitGateway := rabbit.New(rabbitURL)
	rabbitGateway.Register(func(data []byte) error {
		fmt.Printf("Message with body %s\n", string(data[:]))
		return rabbitGateway.Publish(data)
	})
	manager := asynctask.Manager{}
	manager.Add(rabbitGateway)
	manager.RunTasks()
}

func getEnv(name string) string {
	if x, ok := os.LookupEnv(name); !ok {
		panic(fmt.Sprintf("%s is required", name))
	} else {
		return x
	}
}
