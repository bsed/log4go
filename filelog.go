// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package log4go

import (
	"fmt"
	"github.com/bsed/log4go/support"
	"os"
	"time"
	"sync"
)

// This log writer sends output to a file
type FileLogWriter struct {
	rec chan *LogRecord
	rot chan bool

	// The opened file
	fileprefix string
	filename string
	file     *os.File

	// The logging format
	format string

	// File header/trailer
	header, trailer string

	// Rotate at linecount
	maxlines          int
	maxlines_curlines int

	// Rotate at size
	maxsize         int64
	maxsize_cursize int64

	// Rotate daily
	daily bool
	// daily_opendate int
	daily_opendaystr string

	// Keep old logfiles (.001, .002, etc)
	rotate    bool
	maxbackup int
}

// This is the FileLogWriter's output method
func (w *FileLogWriter) LogWrite(rec *LogRecord) {
	w.rec <- rec
}


var lock = new(sync.Mutex)
var cond = sync.NewCond(lock)

func (w *FileLogWriter) Close() {
	lock.Lock()
	close(w.rec)
	w.file.Sync()
	cond.Wait()
	lock.Unlock()
}

// NewFileLogWriter creates a new LogWriter which writes to the given file and
// has rotation enabled if rotate is true.
//
// If rotate is true, any time a new log file is opened, the old one is renamed
// with a .### extension to preserve it.  The various Set* methods can be used
// to configure log rotation based on lines, size, and daily.
//
// The standard log-line format is:
//   [%D %T] [%L] (%S) %M
func NewFileLogWriter(fname string, rotate, daily bool) *FileLogWriter {
	w := &FileLogWriter{
		rec:       make(chan *LogRecord, LogBufferLength),
		rot:       make(chan bool),
		filename:  fname,
		format:    "[%D %T] [%L] (%S) %M",
		rotate:    rotate,
		daily:     daily,
		maxbackup: 999,
	}

	w.filename = w.genFileName()

	if _, err := os.Lstat(w.filename); err == nil {
		_, ctime, _, err := support.GetStatTime(w.filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
			return nil
		}
		w.daily_opendaystr = ctime.Format("2006-01-02")
		w.maxlines_curlines = support.GetLines(w.filename)
		w.maxsize_cursize = support.GetSize(w.filename)
	}

	// open the file for the first time
	if err := w.intRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
		return nil
	}

	go func() {
		defer func() {
			if w.file != nil {
				fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
				w.file.Close()
			}
		}()

		for {
			select {
			case <-w.rot:
				if err := w.intRotate(); err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
					return
				}
			case rec, ok := <-w.rec:
				if !ok {
					// signal notification wakeup
					cond.Signal()

					return
				}
				now := time.Now()
				if (w.maxlines > 0 && w.maxlines_curlines > w.maxlines) ||
					(w.maxsize > 0 && w.maxsize_cursize > w.maxsize) ||
					(w.daily && now.Format("2006-01-02") != w.daily_opendaystr) {
					if err := w.intRotate(); err != nil {
						fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
						return
					}
				}

				// Perform the write
				n, err := fmt.Fprint(w.file, FormatLogRecord(w.format, rec))
				if err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
					return
				}

				// Update the counts
				w.maxlines_curlines++
				w.maxsize_cursize += int64(n)
			}
		}
	}()

	return w
}

// Request that the logs rotate
func (w *FileLogWriter) Rotate() {
	w.rot <- true
}

// If this is called in a threaded context, it MUST be synchronized
func (w *FileLogWriter) intRotate() error {
	// Close any log file that may be open
	if w.file != nil {
		fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
		w.file.Close()
	}

	now := time.Now()

	// If we are keeping log files, move it to the next available number
	if w.rotate {
		_, err := os.Lstat(w.filename)
		if err == nil { // file exists
			num := 1
			fname := ""
			todayDate := time.Now().Format("2006-01-02")
			if w.daily && todayDate != w.daily_opendaystr {
				// another day, rename all old log file
				for ; err == nil && num <= 999; num++ {
					fname = w.filename + fmt.Sprintf(".%03d", num)
					nfname := w.filename + fmt.Sprintf(".%s.%03d", w.daily_opendaystr, num)
					_, err = os.Lstat(fname)
					if err == nil {
						os.Rename(fname, nfname)
					}
				}
				// return error if the last file checked still existed
				if err == nil {
					return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", w.filename)
				} else {
					fname = w.filename + fmt.Sprintf(".%s", w.daily_opendaystr)
				}
			} else if (w.maxlines > 0 && w.maxlines_curlines > w.maxlines) ||
				(w.maxsize > 0 && w.maxsize_cursize > w.maxsize) {
				// maxlines or maxsize reached, create new log and rename the old
				num = w.maxbackup - 1
				for ; num >= 1; num-- {
					fname = w.filename + fmt.Sprintf(".%03d", num)
					nfname := w.filename + fmt.Sprintf(".%03d", num+1)
					_, err = os.Lstat(fname)
					if err == nil {
						os.Rename(fname, nfname)
					}
				}
			} else {
				// first time init logger, reuse old log file if exist, here we do nothing
			}

			if w.file != nil { w.file.Close() }

			// Rename the file to its newfound home
			if fname != "" {
				err = os.Rename(w.filename, fname)
				if err != nil {
					return fmt.Errorf("Rotate: %s\n", err)
				}
			}
		}
	}else if (w.daily) {
		// for daily log output
		w.filename = w.genFileName()
	}

	// Open the log file
	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	w.file = fd

	fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: now}))

	// Set the daily open date to the current date
	//	w.daily_opendate = now.Day()
	w.daily_opendaystr = now.Format("2006-01-02")

	// initialize rotation values
	w.maxlines_curlines = 0
	w.maxsize_cursize = 0

	return nil
}

// Set the logging format (chainable).  Must be called before the first log
// message is written.
func (w *FileLogWriter) SetFormat(format string) *FileLogWriter {
	w.format = format
	return w
}

// Set the logfile header and footer (chainable).  Must be called before the first log
// message is written.  These are formatted similar to the FormatLogRecord (e.g.
// you can use %D and %T in your header/footer for date and time).
func (w *FileLogWriter) SetHeadFoot(head, foot string) *FileLogWriter {
	w.header, w.trailer = head, foot
	if w.maxlines_curlines == 0 {
		fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: time.Now()}))
	}
	return w
}

// Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *FileLogWriter) SetRotateLines(maxlines int) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateLines: %v\n", maxlines)
	w.maxlines = maxlines
	return w
}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *FileLogWriter) SetRotateSize(maxsize int64) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateSize: %v\n", maxsize)
	w.maxsize = maxsize
	return w
}

// Set rotate daily (chainable). Must be called before the first log message is
// written.
func (w *FileLogWriter) SetRotateDaily(daily bool) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateDaily: %v\n", daily)
	w.daily = daily
	return w
}

// Set max backup files. Must be called before the first log message
// is written.
func (w *FileLogWriter) SetRotateMaxBackup(maxbackup int) *FileLogWriter {
	w.maxbackup = maxbackup
	return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *FileLogWriter) SetRotate(rotate bool) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotate: %v\n", rotate)
	w.rotate = rotate
	return w
}


func (w *FileLogWriter) SetFilePrefix(prefix string) *FileLogWriter {
	w.fileprefix =prefix
	return w
}

func (w *FileLogWriter) genFileName() string {
	now := time.Now()
	return fmt.Sprintf("%s%d%02d%02d.log", w.fileprefix, now.Year(), now.Month(), now.Day())
}

// NewXMLLogWriter is a utility method for creating a FileLogWriter set up to
// output XML record log messages instead of line-based ones.
func NewXMLLogWriter(fname string, rotate, daily bool) *FileLogWriter {
	return NewFileLogWriter(fname, rotate, daily).SetFormat(
		`	<record level="%L">
		<timestamp>%D %T</timestamp>
		<source>%S</source>
		<message>%M</message>
	</record>`).SetHeadFoot("<log created=\"%D %T\">", "</log>")
}
