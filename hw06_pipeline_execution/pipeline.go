package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	select {
	case <-done:
		out := make(Bi)
		close(out)
		return out
	default:
	}

	out := in
	for _, s := range stages {
		out = WithDone(s, done)(out)
	}
	return out
}

func WithDone(stage Stage, done In) Stage {
	return func(in In) Out {
		out := make(Bi)
		stageOut := stage(in)

		go func() {
			defer close(out)

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
					case out <- v:
					case <-done:
						go drainStageOut(stageOut)
						return
					}
				}
			}
		}()

		return out
	}
}

func drainStageOut(out Out) {
	if out == nil {
		return
	}
	for v := range out {
		_ = v
	}
}
