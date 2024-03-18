package pipeliner

const DefaultMaxStages = 200

type Pipeline[T any] struct {
	fail_on_error bool
	Stages        []Stage[T]
	MaxStages     int
}

func NewPipeline[T any]() *Pipeline[T] {
	return newPipeline[T](true)
}

func newPipeline[T any](fail_on_pattern_error bool) *Pipeline[T] {
	return &Pipeline[T]{Stages: []Stage[T]{}, fail_on_error: fail_on_pattern_error, MaxStages: DefaultMaxStages}
}

func (p *Pipeline[T]) canAddStage(stage Stage[T]) error {
	if len(p.Stages) == p.MaxStages {
		return ErrMaxStagesExceeded
	}
	for _, s := range p.Stages {
		if s.Name == stage.Name {
			return ErrStageNameExists
		}
	}
	return nil
}

func (p *Pipeline[T]) RegisterStageLast(stage Stage[T]) error {
	if err := p.canAddStage(stage); err != nil {
		return err
	}

	p.Stages = append(p.Stages, stage)
	return nil
}

func (p *Pipeline[T]) RegisterStageFirst(stage Stage[T]) error {
	if err := p.canAddStage(stage); err != nil {
		return err
	}
	p.Stages = append([]Stage[T]{stage}, p.Stages...)
	return nil
}

func (p *Pipeline[T]) RegisterStageAt(index int, stage Stage[T]) error {
	if err := p.canAddStage(stage); err != nil {
		return err
	}
	p.Stages = append(p.Stages[:index], append([]Stage[T]{stage}, p.Stages[index:]...)...)
	return nil
}

func (p *Pipeline[T]) RemoveStage(name string) error {
	for i, stage := range p.Stages {
		if stage.Name == name {
			p.Stages = append(p.Stages[:i], p.Stages[i+1:]...)
			return nil
		}
	}
	return ErrStageNotFound
}

func (p *Pipeline[T]) execute(data T, keepActions bool) ([]ActionResult[T], error) {
	var results []ActionResult[T]
	for i, stage := range p.Stages {
		matched, err := stage.PatternMatcher(data)
		if err != nil {
			if p.fail_on_error {
				return nil, err
			}
			continue
		}
		if matched {
			result := stage.Action(data)
			if keepActions {
				result.StageOfOrigin = &p.Stages[i]
				results = append(results, result)
			}
			if !result.Continue {
				break
			}
		}
	}
	return results, nil
}

func (p *Pipeline[T]) Execute(data T) error {
	_, err := p.execute(data, false)
	return err
}

func (p *Pipeline[T]) ExecuteWithActions(data T) ([]ActionResult[T], error) {
	return p.execute(data, true)
}
