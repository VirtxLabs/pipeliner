package pipeliner

type ActionResult[T any] struct {
	StageOfOrigin *Stage[T]
	Success       bool
	Continue      bool
	Message       string
	Data          any
}

func Error[T any](message string, data any) ActionResult[T] {
	return ActionResult[T]{Success: false, Continue: false, Message: message, Data: data}
}

func Success[T any](message string, data any) ActionResult[T] {
	return ActionResult[T]{Success: true, Continue: true, Message: message, Data: data}
}

func Done[T any](message string, data any) ActionResult[T] {
	return ActionResult[T]{Success: true, Continue: false, Message: message, Data: data}
}

type Stage[T any] struct {
	Name           string                  // stores the name of the stage
	PatternMatcher func(T) (bool, error)   // stores the pattern matcher to be used by the stage
	Action         func(T) ActionResult[T] // stores the action to be performed by the stage
}

func NewStage[T any](name string, pattern_matcher func(T) (bool, error), action func(T) ActionResult[T]) Stage[T] {
	return Stage[T]{
		Name:           name,
		PatternMatcher: pattern_matcher,
		Action:         action,
	}
}
