package sym

import pkg "github.com/ardnew/bases/log"

const logID = "SYM"

var log = pkg.LookupNew(logID)

func init() { log.SetCallerOffset(2) }

func logf(format string, v ...any) { log.Printf(format, v...) }