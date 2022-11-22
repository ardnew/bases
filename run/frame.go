package run

import "runtime"

func Callstack(depth int) (cs []runtime.Frame) {
	pc := make([]uintptr, depth)
	cf := runtime.CallersFrames(pc[:runtime.Callers(1, pc)])
	// A fixed number of PCs can expand to an indefinite number of Frames.
	for {
		f, more := cf.Next()
		cs = append(cs, f)
		if !more || len(cs) == depth {
			return
		}
	}
}
