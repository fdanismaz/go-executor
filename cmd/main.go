package main

import (
	"fmt"
	"github.com/fdanismaz/go-executor/executor"
)

type sampleTask struct {
}

func (sampleTask) Run() error {
	fmt.Println("a sample task is running")
	return nil
}

func (sampleTask) GetName() string {
	return "sample task"
}

func main() {
	exec := executor.NewParallel(3, true)
	exec.Submit(sampleTask{})
	err := exec.Execute()
	if err != nil {
		fmt.Printf("executor failed, reason: %v", err)
	}
}
