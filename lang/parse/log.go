package parse

import pkg "github.com/ardnew/bases/log"

const logID = "PARSE"

var log = pkg.LookupNew(logID)

func init() { log.SetCallerOffset(2) }

func logf(format string, v ...any) { log.Printf(format, v...) }
