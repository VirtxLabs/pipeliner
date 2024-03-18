package pipeliner

import "sync"

const DefaultMaxStages = 200

type Pipeline[T any] struct {
	fail_on_error bool
	lock          sync.RWMutex
	stages        []Stage[T]
	max_stages    int
}

func NewPipeline[T any]() *Pipeline[T] {
	return newPipeline[T](true)
}

func newPipeline[T any](fail_on_pattern_error bool) *Pipeline[T] {
	return &Pipeline[T]{stages: []Stage[T]{}, fail_on_error: fail_on_pattern_error, max_stages: DefaultMaxStages, lock: sync.RWMutex{}}
}

func (p *Pipeline[T]) SetMaxStages(max int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.max_stages = max
}

func (p *Pipeline[T]) canAddStage(stage Stage[T]) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if len(p.stages) == p.max_stages {
		return ErrMaxStagesExceeded
	}
	for _, s := range p.stages {
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

	p.stages = append(p.stages, stage)
	return nil
}

func (p *Pipeline[T]) RegisterStageFirst(stage Stage[T]) error {
	if err := p.canAddStage(stage); err != nil {
		return err
	}
	p.stages = append([]Stage[T]{stage}, p.stages...)
	return nil
}

func (p *Pipeline[T]) RegisterStageAt(index int, stage Stage[T]) error {
	if err := p.canAddStage(stage); err != nil {
		return err
	}
	p.stages = append(p.stages[:index], append([]Stage[T]{stage}, p.stages[index:]...)...)
	return nil
}

func (p *Pipeline[T]) RemoveStage(name string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	for i, stage := range p.stages {
		if stage.Name == name {
			p.stages = append(p.stages[:i], p.stages[i+1:]...)
			return nil
		}
	}
	return ErrStageNotFound
}

func (p *Pipeline[T]) execute(data T, keepActions bool) ([]ActionResult[T], error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	var results []ActionResult[T]
	for i, stage := range p.stages {
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
				result.StageOfOrigin = &p.stages[i]
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
