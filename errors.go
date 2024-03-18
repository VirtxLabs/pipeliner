package pipeliner

type PipelineError string

const (
	ErrMaxStagesExceeded PipelineError = "Max stages exceeded"
	ErrStageNotFound     PipelineError = "Stage not found"
	ErrStageNameExists   PipelineError = "Stage name already exists"
)

func (e PipelineError) Error() string {
	return string(e)
}
