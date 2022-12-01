package run

import "runtime"

func Callstack(depth int) (cs []runtime.Frame) {
	var skip int
	switch {
	case depth > 0:
		skip = 1
	case depth < 0:
		skip -= depth // Increment
		fallthrough
	default:
		depth = 1
	}
	pc := make([]uintptr, depth)
	pc = pc[:runtime.Callers(skip, pc)]
	if len(pc) == 0 {
		return // Should only occur when skip > stack
	}
	cf := runtime.CallersFrames(pc)
	// A fixed number of PCs can expand to an indefinite number of Frames.
	for {
		f, more := cf.Next()
		cs = append(cs, f)
		if !more || len(cs) == depth {
			return
		}
	}
}
