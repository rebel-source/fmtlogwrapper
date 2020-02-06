// Utilities: Logging
package fmtlogwrapper

/*
	@author Arjun Dhar
*/

import (
	"fmt"
	"io"
	"os"
	"time"

	rlog "log"

	slog "github.com/sirupsen/logrus"
)

/*
Settings for the Log
*/
type LogSettings struct {
	// Any meta information that has to be present in the record specific log
	// Example: a JobID etc that we want to appear in the log so the caller can ID it
	// with context specific info.
	JobRecMeta string

	// The file path for the log file
	FilePath string

	// An optional Proxy for this logger, additionally where the logs can go like an appender)
	Proxy ProxyLogger
}

// A wrapper structure over core logger and settings, file, writers associated with it
type Logger struct {
	settings LogSettings
	slogger  *slog.Logger
	rlogger  *rlog.Logger
	file     *os.File
	writer   io.Writer

	muted 	 bool
}

var AppLogger *Logger = NewLogger(LogSettings{
	FilePath: ".\\log\\app.log", // Default path
})

// Create New instance of Logger
/*
Return a common Singleton root logger with application specific configs

**Usage**:
```
import ("github.com/sirupsen/logrus")
var slog *logrus.Logger = GetStructuredLogger()

slog.WithFields(logrus.Fields{
				"error": ferr,
			}).Error("getting working directory")
```

```
...
slog.Infof("Can't load config file %s/app.config. Error: %v", resDir, err)
```
*/
func NewLogger(settings LogSettings) *Logger {
	l := &Logger{
		settings: settings,
		slogger:  &slog.Logger{},
		rlogger:  &rlog.Logger{},
	}
	l.initLogger(settings.FilePath)
	return l
}

func getWriter(logPath string) (io.Writer, error, *os.File) {
	f, err := OpenFilePathExists(logPath)
	var w io.Writer = nil
	if err != nil {
		//rlog.Fatalf("[GetRegularLogger] error opening file: %v", err)
		fmt.Printf("\n[Logger][getWriter] error opening file %s: %v. Dumping logs on screen...\n", logPath, err)
		w = os.Stdout //Default to console
	} else {
		w = io.Writer(f)
		fmt.Printf("\n[Logger][getWriter] Writing logs to %s\n", logPath)
	}
	return w, err, f
}

// Close the log file properly
// Contains redundancy to ensure no error in various cases. Loggers can be in inconsistent state
// depending on what the application is doing : JOB, IDLE, using a PROXY so while closing we dont know
// So approach nil conditions with caution (redundancy in code is required)
func (log *Logger) Close() {
	path := ""
	if log != nil {
		path = log.settings.FilePath
	} else {
		fmt.Println("\n[Logger][Close] Not open; nothing to close")
		return
	}
	fmt.Println("\n[Logger][Close]", path)
	if log.rlogger != nil && !log.muted {
		log.rlogger.Println("\n[Logger][Close]", path)
	}
	if log != nil && log.file != nil {
		defer log.file.Close() //defer since we need to close Proxy logger also
	}
	if log.settings.Proxy != nil {
		log.Println("\n[Logger][Close] proxy")
		if e := log.settings.Proxy.Close(); e != nil {
			fmt.Println("\n[Logger][Close] Failed to Close Proxy logger")
			if log.rlogger != nil && !log.muted {
				log.rlogger.Println("\n[Logger][Close] Failed to Close Proxy logger")
			}
		}
		log.settings.Proxy = nil
	}
}

func (log *Logger) GetProxy() ProxyLogger {
	return log.settings.Proxy
}

// Set the Proxy logger and Init() it
func (log *Logger) SetProxy(proxy ProxyLogger) {
	if proxy == nil && log.settings.Proxy == nil {
		return //de-rferenced  anyway, quit
	}
	log.settings.Proxy = proxy
	if proxy != nil {
		if e := proxy.Init(); e != nil {
			log.Println("\n[Logger][SetProxy] Faile to initialize Proxy logger")
		}
	}
}

// Open an existing log that was closed
func (log *Logger) Open() {
	f, err := os.OpenFile(log.settings.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("\n[Logger][Open] error re-opening file %s: %v. Dumping logs on screen...\n", log.settings.FilePath, err)
		log.writer = os.Stdout //Default to console
	} else {
		log.file = f
		log.writer = io.Writer(f)
	}
}

// Can supress logging to File while muted
func (log *Logger) MuteWrite(mute bool) {
	log.muted = mute
}

/*
	slogger.SetFormatter(&slog.JSONFormatter{DisableTimestamp: true})
	slogger.SetFormatter(&slog.JSONFormatter{DisableTimestamp: true})
@return *log.Logger, *io.Writer or os.File
*/
func (log *Logger) initLogger(logPath string) {
	slogger := log.slogger
	slogger.SetFormatter(&slog.JSONFormatter{DisableTimestamp: true})
	slogger.SetReportCaller(false)
	w, _, f := getWriter(logPath)
	slogger.SetOutput(w) //os.Stdout
	log.rlogger.SetOutput(w)
	//....SetFormatter(&slog.TextFormatter{})
	slogger.SetLevel(slog.InfoLevel)
	log.writer = w
	log.file = f
}

// Log the data to the file in the dired format
// Before calling Log, ensure Init logger has been called
func (log *Logger) Log(logRec map[string]interface{}) {
	timestamp := time.Now()
	timeStr := timestamp.Format("02-01-2006 15:04:05" /*time.RFC1123*/) //dd-mm-yyyy hh:MM:ss
	logRec["datetime"] = timeStr
	logRec["epoch"] = timestamp.Unix()
	logRec["meta"] = log.settings.JobRecMeta

	log.slogger.WithFields(logRec).Info("")

	if log.settings.Proxy != nil {
		if err := log.settings.Proxy.Log(logRec); err != nil {
			log.Printf("\n[ERROR][Logger][Log]Could not send this log to server: %v", logRec)
		}
	}
}

func (log *Logger) LogStr(str string) {
	if  !log.muted {
		log.rlogger.Println(str)
	}
}

// Following serves as a convenient replacement for fmt.<...>

// Replacement for fmt.Printf
func (log *Logger) Printf(str string, params ...interface{}) {
	fmt.Printf(str, params...) //Send to std console always
	if  !log.muted {
		log.rlogger.Printf(str, params...)
	}
}

// Replacement for fmt.Println
func (log *Logger) Println(a ...interface{}) {
	fmt.Println(a...) //Send to std console always
	if  !log.muted {
		log.rlogger.Println(a...)
	}
}

func (log *Logger) Errorf(str string, params ...interface{}) {
	fmt.Errorf(str, params...)
	if  !log.muted {
		log.rlogger.Fatalf(str, params...)
	}
}

//A proxy that can log to a network device
type ProxyLogger interface {
	Init() error

	Log(logRec map[string]interface{}) error

	//Any closingoperations on the proxy
	Close() error
}

/////////////////////////////// Multiple Loggers
var loggers map[string]*Logger = make(map[string]*Logger)


/*
 Maintains a reference to the logger within the logger framework
*/
func InitContextLogger(contextId string, settings LogSettings) *Logger {
	logger := NewLogger(settings)
	loggers[contextId] = logger
	return logger
}


/*
 Take al lthe references and do what you want ! Be Happy!
*/
func ContextLoggers() map[string]*Logger {
	return loggers
}