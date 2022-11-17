package sym

import pkg "github.com/ardnew/bases/log"

const logID = "SYM"

var log = pkg.LookupNew(logID)

func logf(format string, v ...any) {
	log.Call(1, func() { log.Printf(format, v...) })
}
