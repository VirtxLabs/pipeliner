# Pipeliner

Pipeliner is a Go library that allows you to create data processing pipelines. Each pipeline is composed of multiple stages, where each stage has an action to execute and a pattern match function to check if the action should be performed on the data.

## Installation

To install Pipeliner, you can use `go get`:

```sh
go get github.com/VirtxLabs/pipeliner
```

## Usage
First, import the `pipeliner` package in your Go file:

```go
import "github.com/VirtxLabs/pipeliner"
```

Then, you can create a new pipeline with your desired data type (e.g. `int`):

```go
mypipeline := pipeliner.NewPipeline[int]()
```

You can then register stages to the pipeline. Each stage has a name, a pattern match function, and an action function. The pattern match function checks if the action should be performed on the data, and the action function performs the action on the data and returns an `ActionResult`.

The first stage checks if the number is even and prints it if it is:
```go
mypipeline.RegisterStageLast(pipeliner.NewStage("Even", func(i int) (bool, error) {
    // Check if the number is even
    return i%2 == 0, nil
}, func(i int) pipeliner.ActionResult[int] {
    // Print the even number
    fmt.Println("Even:", i)
    // The action was successful
    return pipeliner.Success[int]("Success", i)
}))
```

The second stage checks if the number is greater than 5 and prints it if it is:
```go

mypipeline.RegisterStageLast(pipeliner.NewStage("GreaterThan5", func(i int) (bool, error) {
    // Check if the number is greater than 5
    return i > 5, nil
}, func(i int) pipeliner.ActionResult[int] {
    // Print the number greater than 5
    fmt.Println("Greater than 5:", i)
    // The action was successful
    return pipeliner.Success[int]("Success", i)
}))
```

Finally, you can run the pipeline with your data:

```go
for i := 0; i < 10; i++ {
    mypipeline.Execute(i)
}
```

This will print the following output:

```
Even: 0
Even: 2
Even: 4
Even: 6
Greater than 5: 6
Greater than 5: 7
Even: 8
Greater than 5: 8
Greater than 5: 9
```

If you want the `Even` stage to end the pipeline, you can replace the `pipeliner.Success[int]("Success", i)` with `pipeliner.Done[int]("Success", i)` in the return statement. This will result in the following output:

```
Even: 0
Even: 2
Even: 4
Even: 6
Greater than 5: 7
Even: 8
Greater than 5: 9
```


## Testing



Tests are located in the `pipeline_test.go` file. You can run them using `go test`:

```sh
go test
```

## Contributing

If you find any issues or have any suggestions, feel free to open an issue or a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
