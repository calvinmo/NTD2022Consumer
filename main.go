package main

import (
	"fmt"
	"gray.net/app-onedeploy/pkg/core/rabbit"
	asynctask "gray.net/lib-go-async-task-manager"
	"os"
)

func main() {
	rabbitURL := getEnv("RABBIT_URL")
	fmt.Printf(rabbitURL)
	rabbitGateway := rabbit.New(rabbitURL)
	rabbitGateway.Register("Incoming", func(data []byte) error {
		message := string(data[:])
		fmt.Printf("Message with body %s", message)
		return rabbitGateway.Publish("Outgoing", message)
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
