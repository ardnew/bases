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

// FmtSpec represents a printf-style format specifier ('%'-escaped rune).
//
// No padding or width arguments are supported. The specifiers are simply
// text replacement tokens performed in a single scan of the format string.
type FmtSpec rune

const (
	Message  FmtSpec = 's' // The actual log message
	Date     FmtSpec = 'd' // The date in the local time zone: 2009/01/23
	Time     FmtSpec = 't' // The time in the local time zone: 01:23:23
	Micros   FmtSpec = 'u' // Same as Time (to microseconds): 01:23:23.123123
	FilePath FmtSpec = 'F' // Full file path: /path/to/file/name.go
	FileBase FmtSpec = 'f' // Base file name: name.go
	Line     FmtSpec = 'n' // Line number: 23
)

var FmtDefault = "%d %t ┆ %f:%n ┆ %s"

// A Log represents an active logging object that generates lines of output to
// an io.Writer. Each logging operation makes a single call to the Writer's
// Write method. A Logger can be used simultaneously from multiple goroutines;
// it guarantees to serialize access to the Writer.
type Log struct {
	mu  sync.Mutex // Ensures atomic writes; protects the following fields
	fmt []rune     // Output format
	out io.Writer  // Destination for output
	buf []byte     // For accumulating text to write
	nul int32      // Atomic boolean: whether out == io.Discard
}

// New creates a new Log. The out variable sets the destination to which log
// data will be written. The format of the log message is defined via fmt which
// contains printf-style specifiers ('%'-escaped runes).
func New(out io.Writer, fmt string) *Log {
	l := &Log{out: out, fmt: []rune(fmt)}
	if out == io.Discard {
		l.nul = 1
	}
	return l
}

// Output writes the output for a logging event.
// A newline is appended if the last character of s is not already a newline.
func (l *Log) Output(calldepth int, s string) error {
	l.format(calldepth, s)
	_, err := l.out.Write(l.buf)
	return err
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Log) Printf(format string, v ...any) {
	if atomic.LoadInt32(&l.nul) != 0 {
		return
	}
	l.Output(2, fmt.Sprintf(format, v...))
}

// Print calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Print.
func (l *Log) Print(v ...any) {
	if atomic.LoadInt32(&l.nul) != 0 {
		return
	}
	l.Output(2, fmt.Sprint(v...))
}

// Println calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Println.
func (l *Log) Println(v ...any) {
	if atomic.LoadInt32(&l.nul) != 0 {
		return
	}
	l.Output(2, fmt.Sprintln(v...))
}

// Fatal is equivalent to l.Print() followed by a call to os.Exit(1).
func (l *Log) Fatal(v ...any) {
	l.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Log) Fatalf(format string, v ...any) {
	l.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Fatalln is equivalent to l.Println() followed by a call to os.Exit(1).
func (l *Log) Fatalln(v ...any) {
	l.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Log) Panic(v ...any) {
	s := fmt.Sprint(v...)
	l.Output(2, s)
	panic(s)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Log) Panicf(format string, v ...any) {
	s := fmt.Sprintf(format, v...)
	l.Output(2, s)
	panic(s)
}

// Panicln is equivalent to l.Println() followed by a call to panic().
func (l *Log) Panicln(v ...any) {
	s := fmt.Sprintln(v...)
	l.Output(2, s)
	panic(s)
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid
// zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

var spec = map[FmtSpec]func(b *[]byte, s string, t time.Time, f string, n int){
	Message: func(b *[]byte, s string, t time.Time, f string, n int) {
		*b = append(*b, s...)
	},
	Date: func(b *[]byte, s string, t time.Time, f string, n int) {
		year, month, day := t.Date()
		itoa(b, year, 4)
		*b = append(*b, '/')
		itoa(b, int(month), 2)
		*b = append(*b, '/')
		itoa(b, day, 2)
	},
	Time: func(b *[]byte, s string, t time.Time, f string, n int) {
		hour, min, sec := t.Clock()
		itoa(b, hour, 2)
		*b = append(*b, ':')
		itoa(b, min, 2)
		*b = append(*b, ':')
		itoa(b, sec, 2)
	},
	Micros: func(b *[]byte, s string, t time.Time, f string, n int) {
		// spec[Time](b, t, file, line)
		hour, min, sec := t.Clock()
		itoa(b, hour, 2)
		*b = append(*b, ':')
		itoa(b, min, 2)
		*b = append(*b, ':')
		itoa(b, sec, 2)
		*b = append(*b, '.')
		itoa(b, t.Nanosecond()/1e3, 6)
	},
	FilePath: func(b *[]byte, s string, t time.Time, f string, n int) {
		*b = append(*b, f...)
	},
	FileBase: func(b *[]byte, s string, t time.Time, f string, n int) {
		short := f
		for i := len(f) - 1; i > 0; i-- {
			if f[i] == '/' {
				short = f[i+1:]
				break
			}
		}
		// spec[FilePath](b, t, short, line)
		*b = append(*b, short...)
	},
	Line: func(b *[]byte, s string, t time.Time, f string, n int) {
		itoa(b, n, -1)
	},
}

// format replaces printf-style format specifiers found in fmt with their
// corresponding values.
func (l *Log) format(calldepth int, message string) {
	now, line := time.Now(), -1
	l.mu.Lock()
	defer l.mu.Unlock()
	var file string
	l.buf = l.buf[:0]
	for i := 0; i < len(l.fmt)-1; i++ {
		if l.fmt[i] == '%' {
			i++
			if f, ok := spec[FmtSpec(l.fmt[i])]; ok {
				if line < 0 { // runtime callstack has never been retrieved
					if FmtSpec(l.fmt[i]) == FileBase ||
						FmtSpec(l.fmt[i]) == FilePath {
						// release lock while getting caller info - it's expensive.
						l.mu.Unlock()
						var ok bool
						if _, file, line, ok = runtime.Caller(calldepth); !ok {
							file = "???"
							line = 0
						}
						l.mu.Lock()
					}
				}
				f(&l.buf, message, now, file, line)
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
