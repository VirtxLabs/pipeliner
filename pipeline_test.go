package pipeliner

import (
	"slices"
	"testing"
)

func TestPipelineMiddleware(t *testing.T) {
	p := NewPipeline[int]()
	if p == nil {
		t.Error("Pipeline not created")
	}
	results := []int{}

	var action = func(i int) ActionResult[int] {
		results = append(results, i)
		return Success[int]("Success", i)
	}

	p.RegisterStageLast(NewStage("Even", func(i int) (bool, error) {
		return i%2 == 0, nil
	}, action))
	p.RegisterStageLast(NewStage("GreaterThanFive", func(i int) (bool, error) {
		return i > 5, nil
	}, action))

	for i := 0; i < 10; i++ {
		if err := p.Execute(i); err != nil {
			t.Error(err)
		}
	}
	var expected = []int{0, 2, 4, 6, 6, 7, 8, 8, 9}

	if !slices.Equal(results, expected) {
		t.Errorf("Expected %v but got %v", expected, results)
	}
}

func TestPipelineEarlyFinish(t *testing.T) {
	p := NewPipeline[int]()
	if p == nil {
		t.Error("Pipeline not created")
	}
	results := []int{}

	p.RegisterStageLast(NewStage("Even", func(i int) (bool, error) {
		return i%2 == 0, nil
	}, func(i int) ActionResult[int] {
		results = append(results, i)
		return Done[int]("Success", i)
	}))

	p.RegisterStageLast(NewStage("GreaterThanFive", func(i int) (bool, error) {
		return i > 5, nil
	}, func(i int) ActionResult[int] {
		results = append(results, i)
		return Done[int]("Success", i)
	}))

	for i := 0; i < 10; i++ {
		if err := p.Execute(i); err != nil {
			t.Error(err)
		}
	}
	var expected = []int{0, 2, 4, 6, 7, 8, 9}

	if !slices.Equal(results, expected) {
		t.Errorf("Expected %v but got %v", expected, results)
	}
}

func TestPipelineMiddlewareActionTracking(t *testing.T) {
	p := NewPipeline[int]()
	if p == nil {
		t.Error("Pipeline not created")
	}
	results := []int{}
	var evenMatcher = func(i int) (bool, error) {
		return i%2 == 0, nil
	}
	var greaterThanFiveMatcher = func(i int) (bool, error) {
		return i > 5, nil
	}

	var action = func(i int) ActionResult[int] {
		results = append(results, i)
		return Success[int]("Success", i)
	}

	evenStage := NewStage("Even", evenMatcher, action)
	greaterThanFiveStage := NewStage("GreaterThanFive", greaterThanFiveMatcher, action)
	p.RegisterStageLast(evenStage)
	p.RegisterStageLast(greaterThanFiveStage)

	var actionTracker = []ActionResult[int]{}
	for i := 0; i < 10; i++ {
		a, err := p.ExecuteWithActions(i)
		if err != nil {
			t.Error(err)
		}
		actionTracker = append(actionTracker, a...)
	}
	expected := []ActionResult[int]{}

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			expected = append(expected, ActionResult[int]{
				StageOfOrigin: &evenStage,
				Success:       true,
				Continue:      true,
				Message:       "Success",
				Data:          i,
			})
		}
		if i > 5 {
			expected = append(expected, ActionResult[int]{
				StageOfOrigin: &greaterThanFiveStage,
				Success:       true,
				Continue:      true,
				Message:       "Success",
				Data:          i,
			})
		}
	}
	t.Logf("%p:%p", &greaterThanFiveStage, &evenStage)

	for _, a := range actionTracker {
		t.Logf("Stage: %s, Success: %v, Continue: %v, Message: %v, Data: %v", a.StageOfOrigin.Name, a.Success, a.Continue, a.Message, a.Data)
	}
	if len(actionTracker) != len(expected) {
		t.Errorf("Expected %v but got %v", expected, actionTracker)
	}

	for i, a := range actionTracker {
		if a.StageOfOrigin.Name != expected[i].StageOfOrigin.Name {
			t.Errorf("Expected %v but got %v", expected[i].StageOfOrigin, a.StageOfOrigin)
		}
		if a.Success != expected[i].Success {
			t.Errorf("Expected %v but got %v", expected[i].Success, a.Success)
		}
		if a.Continue != expected[i].Continue {
			t.Errorf("Expected %v but got %v", expected[i].Continue, a.Continue)
		}
		if a.Message != expected[i].Message {
			t.Errorf("Expected %v but got %v", expected[i].Message, a.Message)
		}
		if a.Data != expected[i].Data {
			t.Errorf("Expected %v but got %v", expected[i].Data, a.Data)
		}
	}
}
