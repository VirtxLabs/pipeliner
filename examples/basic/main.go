package main

import (
	"fmt"

	"github.com/VirtxLabs/pipeliner"
)

func main() {
	mypipeline := pipeliner.NewPipeline[int]()

	mypipeline.RegisterStageLast(pipeliner.NewStage("Even", func(i int) (bool, error) {
		// Check if the number is even
		return i%2 == 0, nil
	}, func(i int) pipeliner.ActionResult[int] {
		// Print the even number
		fmt.Println("Even:", i)
		// The action was successful
		return pipeliner.Done[int]("Success", i)
	}))

	mypipeline.RegisterStageLast(pipeliner.NewStage("GreaterThan5", func(i int) (bool, error) {
		// Check if the number is greater than 5
		return i > 5, nil
	}, func(i int) pipeliner.ActionResult[int] {
		// Print the number greater than 5
		fmt.Println("Greater than 5:", i)
		// The action was successful
		return pipeliner.Success[int]("Success", i)
	}))

	for i := 0; i < 10; i++ {
		mypipeline.Execute(i)
	}

}
