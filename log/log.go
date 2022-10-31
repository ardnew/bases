package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// FormatSpec represents a printf-style format specifier ('%'-escaped rune).
//
// No padding or width arguments are supported. The specifiers are simply
// text replacement tokens performed in a single scan of the format string.
type FormatSpec rune

const (
	Message  FormatSpec = 's' // The actual log message
	Date     FormatSpec = 'd' // The date in the local time zone: 2009/01/23
	Time     FormatSpec = 't' // The time in the local time zone: 01:23:23
	Micros   FormatSpec = 'u' // Same as Time (to microseconds): 01:23:23.123123
	FilePath FormatSpec = 'F' // Full file path: /path/to/file/name.go
	FileBase FormatSpec = 'f' // Base file name: name.go
	Line     FormatSpec = 'n' // Line number: 23
)

var (
	DefaultWriter = os.Stderr
	DefaultFormat = "%d %t ┆ %f:%n ┆ %s"
)

// A Log represents an active logging object that generates lines of output to
// an [io.Writer]. Each logging operation makes a single call to the method
// [io.Writer.Write]. A Logger can be used simultaneously from multiple
// goroutines; it guarantees to serialize access to the Writer.
type Log struct {
	mut sync.Mutex // Ensures atomic writes; protects the following fields
	fmt []rune     // Output format
	out io.Writer  // Destination for output
	buf []byte     // For accumulating text to write
	nul int32      // Atomic boolean: whether out == io.Discard
}

// New creates a new Log. The out variable sets the destination to which log
// data will be written. The format of the log message is defined via fmt which
// contains printf-style specifiers ('%'-escaped runes).
func New(output io.Writer, format string) *Log {
	l := &Log{out: output, fmt: []rune(format)}
	if output == io.Discard {
		l.nul = 1
	}
	return l
}

// Writer returns the output destination for the Log.
func (l *Log) Writer() io.Writer {
	l.mut.Lock()
	defer l.mut.Unlock()
	return l.out
}

// SetWriter sets the output destination for the Log.
func (l *Log) SetWriter(w io.Writer) {
	l.mut.Lock()
	defer l.mut.Unlock()
	l.out = w
	var nul int32
	if w == io.Discard {
		nul = 1
	}
	atomic.StoreInt32(&l.nul, nul)
}

// Format returns the printf-style format string that describes the structure
// and content of log messages.
//
// See type [FormatSpec] for available specifiers and see const [DefaultFormat]
// for an example.
func (l *Log) Format() string {
	l.mut.Lock()
	defer l.mut.Unlock()
	return string(l.fmt)
}

// SetFormat sets the printf-style format string that describes the structure
// and content of log messages.
//
// See type [FormatSpec] for available specifiers and see const [DefaultFormat]
// for an example.
func (l *Log) SetFormat(format string) {
	l.mut.Lock()
	defer l.mut.Unlock()
	l.fmt = []rune(format)
}

// Output writes the output for a logging event.
// A newline is appended if the last character of s is not already a newline.
func (l *Log) Output(calldepth int, s string) error {
	l.format(calldepth, s)
	_, err := l.out.Write(l.buf)
	return err
}

// Printf calls method Output for writing to the Log.
// Arguments are handled in the manner of [fmt.Printf].
func (l *Log) Printf(format string, v ...any) {
	if atomic.LoadInt32(&l.nul) != 0 {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

// Print calls method Output for writing to the Log.
// Arguments are handled in the manner of [fmt.Print].
func (l *Log) Print(v ...any) {
	if atomic.LoadInt32(&l.nul) != 0 {
		return
	}
	l.Output(2, fmt.Sprint(v...))
}

// Println calls method Output for writing to the Log.
// Arguments are handled in the manner of [fmt.Println].
func (l *Log) Println(v ...any) {
	if atomic.LoadInt32(&l.nul) != 0 {
		return
	}
	l.Output(2, fmt.Sprintln(v...))
}

// Fatal is equivalent to Print followed by a call to [os.Exit] (code=1).
func (l *Log) Fatal(v ...any) {
	l.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to Printf followed by a call to [os.Exit] (code=1).
func (l *Log) Fatalf(format string, v ...any) {
	l.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to Println followed by a call to [os.Exit] (code=1).
func (l *Log) Fatalln(v ...any) {
	l.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to Print followed by a call to [panic].
func (l *Log) Panic(v ...any) {
	s := fmt.Sprint(v...)
	l.Output(2, s)
	panic(s)
}

// Panicf is equivalent to Printf followed by a call to [panic].
func (l *Log) Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, s)
	panic(s)
}

// Panicln is equivalent to Println followed by a call to [panic].
func (l *Log) Panicln(v ...any) {
	s := fmt.Sprintln(v...)
	l.Output(2, s)
	panic(s)
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid
// zero-padding.
func itoa(buf *[]byte, i int, width int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	offset := len(b) - 1
	for i >= 10 || width > 1 {
		width--
		q := i / 10
		b[offset] = byte('0' + i - q*10)
		offset--
		i = q
	}
	// i < 10
	b[offset] = byte('0' + i)
	*buf = append(*buf, b[offset:]...)
}

// specArgs encapsulates the set of all potentially-used arguments by any given
// format specifier handler function.
//
// This encapsulation is required to trick the static analyzer into not
// identifying potentially unused arguments in the map functions defined in spec
// below.
type specArgs struct {
	mesg string
	time time.Time
	file string
	line int
}

// specFunc defines a format specifier handler function that formats data
// derived from specArgs and writes the formatted data to the output byte slice
// referenced in out.
//
// The output byte slice referenced in out is not a field of specArgs because
// a.) it is a required (non-nil) argument for all formatting functions, and
// b.) it is an output parameter — not an input argument controlling format.
type specFunc func(out *[]byte, a specArgs)

var spec = map[FormatSpec]func(out *[]byte, a specArgs){
	Message: func(out *[]byte, a specArgs) {
		*out = append(*out, a.mesg...)
	},
	Date: func(out *[]byte, a specArgs) {
		year, month, day := a.time.Date()
		itoa(out, year, 4)
		*out = append(*out, '/')
		itoa(out, int(month), 2)
		*out = append(*out, '/')
		itoa(out, day, 2)
	},
	Time: func(out *[]byte, a specArgs) {
		hour, min, sec := a.time.Clock()
		itoa(out, hour, 2)
		*out = append(*out, ':')
		itoa(out, min, 2)
		*out = append(*out, ':')
		itoa(out, sec, 2)
	},
	Micros: func(out *[]byte, a specArgs) {
		// spec[Time](b, t, file, line)
		hour, min, sec := a.time.Clock()
		itoa(out, hour, 2)
		*out = append(*out, ':')
		itoa(out, min, 2)
		*out = append(*out, ':')
		itoa(out, sec, 2)
		*out = append(*out, '.')
		itoa(out, a.time.Nanosecond()/1e3, 6)
	},
	FilePath: func(out *[]byte, a specArgs) {
		*out = append(*out, a.file...)
	},
	FileBase: func(out *[]byte, a specArgs) {
		short := a.file
		for i := len(a.file) - 1; i > 0; i-- {
			if a.file[i] == '/' {
				short = a.file[i+1:]
				break
			}
		}
		// spec[FilePath](b, t, short, line)
		*out = append(*out, short...)
	},
	Line: func(out *[]byte, a specArgs) {
		itoa(out, a.line, -1)
	},
}

// format replaces printf-style format specifiers found in fmt with their
// corresponding values.
func (l *Log) format(calldepth int, message string) {
	now, line := time.Now(), -1
	l.mut.Lock()
	defer l.mut.Unlock()
	var file string
	l.buf = l.buf[:0]
	for i := 0; i < len(l.fmt)-1; i++ {
		if l.fmt[i] == '%' {
			i++
			if f, ok := spec[FormatSpec(l.fmt[i])]; ok {
				if line < 0 { // runtime callstack has never been retrieved
					if FormatSpec(l.fmt[i]) == FileBase ||
						FormatSpec(l.fmt[i]) == FilePath {
						// release lock while getting caller info - it's expensive.
						l.mut.Unlock()
						var ok bool
						if _, file, line, ok = runtime.Caller(calldepth); !ok {
							file = "???"
							line = 0
						}
						l.mut.Lock()
					}
				}
				f(&l.buf, specArgs{mesg: message, time: now, file: file, line: line})
			} else {
				// format specifier unrecognized, write it out literally
				l.buf = append(l.buf, '%')
				l.buf = append(l.buf, []byte(string(l.fmt[i]))...)
			}
		} else {
			// write any extraneous runes in the format string literally to output
			l.buf = append(l.buf, []byte(string(l.fmt[i]))...)
		}
	}
	if n := len(l.buf); n == 0 || l.buf[n-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
}
