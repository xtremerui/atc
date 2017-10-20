package exec

import (
	"os"

	"github.com/concourse/atc/worker"
	"github.com/hashicorp/go-multierror"
)

// OnAbortStep will run one step, and then a second step if the first step
// aborts (but not errors).
type OnAbortStep struct {
	stepFactory  StepFactory
	abortFactory StepFactory

	repo *worker.ArtifactRepository

	step  Step
	abort Step
}

// OnAbort constructs an OnAbortStep factory.
func OnAbort(firstStep StepFactory, secondStep StepFactory) OnAbortStep {
	return OnAbortStep{
		stepFactory:  firstStep,
		abortFactory: secondStep,
	}
}

// Using constructs an *OnAbortStep.
func (o OnAbortStep) Using(repo *worker.ArtifactRepository) Step {
	o.repo = repo

	o.step = o.stepFactory.Using(o.repo)
	return &o
}

// Run will call Run on the first step and wait for it to complete. If the
// first step errors, Run returns the error. OnAbortStep is ready as soon as
// the first step is ready.
//
// If the first step aborts (that is, it gets interrupted), the second
// step is executed. If the second step errors, its error is returned.
func (o *OnAbortStep) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	var errors error

	stepRunErr := o.step.Run(signals, ready)
	if stepRunErr != nil {
		errors = multierror.Append(errors, stepRunErr)
		if stepRunErr == ErrInterrupted {
			o.abort = o.abortFactory.Using(o.repo)
			hookRunErr := o.abort.Run(signals, make(chan struct{}))
			errors = multierror.Append(errors, hookRunErr)
		}
		return errors
	}

	return nil
}

// Succeeded is true if the first step doesn't exist, or if it
// completed successfully.
func (o *OnAbortStep) Succeeded() bool {
	return o.step.Succeeded()
}
