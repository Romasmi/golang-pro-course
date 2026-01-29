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
		out = WithDone(s, done)(out)
	}
	return out
}

func WithDone(stage Stage, done In) Stage {
	return func(in In) Out {
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
					go drainStageOut(stageOut, done)
					return

				case v, ok := <-stageOut:
					if !ok {
						return
					}

					select {
					case out <- v:
					case <-done:
						go drainStageOut(stageOut, done)
						return
					}
				}
			}
		}()

		return out
	}
}

func drainStageOut(out Out, done In) {
	if out == nil {
		return
	}
	for {
		select {
		case _, ok := <-out:
			if !ok {
				return
			}
		case <-done:
			return
		}
	}
}
