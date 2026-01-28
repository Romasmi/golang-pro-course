package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, s := range stages {
		out = runStage(out, done, s)
	}
	return out
}

func runStage(in, done In, stage Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)
		stageOut := stage(in)
		if stageOut == nil {
			return
		}
		for {
			select {
			case <-done:
				go drainStageOut(stageOut)
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				select {
				case <-done:
					go drainStageOut(stageOut)
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func drainStageOut(out Out) {
	if out == nil {
		return
	}
	for v := range out {
		_ = v
	}
}
