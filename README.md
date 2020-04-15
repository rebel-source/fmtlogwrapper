
# `fmt` `Log` Wrapper

<!-- shields.io - https://shields.io/category/analysis -->
![GoSec Verified](https://img.shields.io/badge/gosec-verified-brightgreen)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/rebel-source/fmtlogwrapper)
![GitHub Last Commit](https://img.shields.io/github/last-commit/rebel-source/fmtlogwrapper)
![Maintenance](https://img.shields.io/maintenance/yes/2020)
![golang-security-action](https://github.com/rebel-source/fmtlogwrapper/workflows/golang-security-action/badge.svg?event=push)

## In one line what is it?
Log that looks like a simple `fmt.Print<xxx>` etc. Abstracting complex and extensible logging behind standard `fmt` methods.<br />
Have multiple log destinations in an application, log objects ....

## How stable is it?
* Been used in Real time IoT situations dealing with Hardware event related logs.
* Big Data situations, multi threaded applications of complex workflow audits showing coherent logs for each channel.
* The first version was production used by April 2019, Multi threaded one, end Jan 2020.
* ...On going reviews and testing. Please test and improve ! :)

## Why do you need it?
I see developers in `golang` struggling to use and then switch between Standard GoLang [log](https://golang.org/pkg/log/), a framework like [logrus](https://github.com/sirupsen/logrus) and even writing prototype code using `fmt`. 

Would it not be nice to log using the `fmt` package standard functions like `Printf`, `Println` in a log agnostic manner that can be backed by the standard `log` package or a framework like `logrus`?

In addition it also allows grouping using:
* **Context**: Logs and settings specific to a certain logger instance
* **WriteToBuffer**: Write continuous logs without having your application to maintain context in multi-threaded code
* **Control log to console and context**: Dont have to repeat logs, sometimes its nice to see something on your screen and file; sometimes one or the other. Easy to toggle.


## Next Steps
* File writes can be `async`, though currently, one can use `WriteToBuffer` as a decent substitute also to optimize writes.
* Add `Appender` Interface.
* Move current Default File writes to a `File Appender` and plug that in as a default.
* Build a `MongoDB Appender`.
* Build a `Fanout Appender` to support joining multiple `Appender`s.

## Usage

### Simple
```golang
import(
	log "github.com/rebel-source/fmtlogwrapper"
)

func TestUsage1() {

	// Note One can define multiple logger instances. This is the default a "singleton" instance; 
	// but nothing stopping you from defining many.if you have a specific reason to have another instance.
	// TIP: Re-use logger instances if they are outputting to the same destination (file or DB)

	log.AppLogger = log.NewLogger(log.LogSettings{
		FilePath: ".\\log\\app.log",
	})

	// Simple Printf ...& Println, Errorf  similarly
	log.AppLogger.Printf("\nRunning version %s\n", "xyz") // See how `log.AppLogger` is a simple replacement for `fmt` so Replace-ALL away :)

	//Log an Object - JSON
	logRec := make(map[string]interface{})
	logRec["key-1"] = "nothing"
	keyAsMap := make(map[string]interface{})
	keyAsMap["data"] = "ABC-123"
	keyAsMap["msg"] = "Nice to see you here :)"
	logRec["key-map"] = keyAsMap
	log.AppLogger.Log(logRec)	

	defer log.AppLogger.Close()
}
```

Expected Output Log
```
[Logger][getWriter] Writing logs to .\log\app.log
=== RUN   TestUsage1

[Logger][getWriter] Writing logs to .\log\app.log

Running version xyz

[Logger][Close] .\log\app.log
--- PASS: TestUsage1 (0.00s)
PASS
ok      fmtlogwrapper   0.341s
```

## Multiple Context Loggers
Maintain multiple contexts to different files, debug levels etc. .. use your imagination. A context is an abstract to group your logs.

```golang
func TestMultipleContextLoggers(t *testing.T) {
	info := /*log.*/InitContextLogger("Info1", /*log.*/LogSettings{
		FilePath: ".\\log\\app-info.log",
	})

	debug := /*log.*/InitContextLogger("Debug1", /*log.*/LogSettings{
		FilePath: ".\\log\\app-debug.log",
	})

	info.Println("STEP 11")
	debug.Println("STEP 12")

	info.MuteWrite(true)
	debug.MuteWrite(true)
	info.Println("STEP 21")  // This wont be written only on console
	debug.Println("STEP 22")  // This wont be written only on console

	// Test can get it afresh from the context
	info = /*log.*/ContextLoggers()["Info1"]
	debug = /*log.*/ContextLoggers()["Debug1"]

	info.MuteWrite(false)
	debug.MuteWrite(false)
	info.Println("STEP 31")
	debug.Println("STEP 32")	

	defer /*log.*/info.Close()
	defer /*log.*/debug.Close()
}
```

## Logging in a Multi Threaded environment
```golang
func TestBufferedLogger(t *testing.T) {
	fmt.Println("[TesBufferedLogger]");
	/*
	 We have 2 loggers writing to the same file
	 We want to ensure a Atomic operation in each does not mix with the other
	*/

	log1 := /*log.*/InitContextLogger("L1", /*log.*/LogSettings{
		FilePath: ".\\log\\app-same.log",
	})
	log2 := /*log.*/InitContextLogger("L2", /*log.*/LogSettings{
		FilePath: ".\\log\\app-same.log",
	})

	log1.WriteToBuffer(true)
	log2.WriteToBuffer(true)

	log1.Println("STEP 1-A")
	log2.Println("STEP 2-A")
	log1.Println("STEP 1-B")
	log2.Println("STEP 2-B")
	log1.Println("STEP 1-C")
	log2.Println("STEP 2-C")	

	log1.WriteToBuffer(false)
	log2.WriteToBuffer(false)
	// At this point buffer so far should be comitted but logs for 1 & 2 should be visible in continious lines for each group
	// We can also explicitly call log<x>.CommitBuffer()

	//Now on will appear in sequence of this being written
	log1.Println("STEP 1-D")
	log2.Println("STEP 2-D")
	log1.Println("STEP 1-E")
	log2.Println("STEP 2-E")

	defer /*log.*/log1.Close() //If buffered was true, will also automatically CommitBuffer() any pending stuff. FYI
	defer /*log.*/log2.Close() //Note: Multiple calls to Close() even if they share the same file is ok.
}
```

Expected Output of above:
```
STEP 1-A 
STEP 1-B 
STEP 1-C 

STEP 2-A 
STEP 2-B 
STEP 2-C 

STEP 1-D
STEP 2-D
STEP 1-E
STEP 2-E

[Logger][Close] .\log\app-same.log

[Logger][Close] .\log\app-same.log
```

### Writing context safe logs
Below is an example, where even the same context and buffer for a cotext may not save you; if within the context your threads write to the same buffer.
To get around this, @ your code level you can employ a simple trick as follows

```golang
/*
 @param contextID string file name depends on this

 @param taskId - Allows multiple processes to write to same file;
 however if we plan to use buffered mode, then each can independently maintain its own buffer.
 Use "" for common.
*/
func InitMyLogger(contextID string, taskId string) *log.Logger {
	logger := log.NewLogger(log.LogSettings{
		FilePath: ".\\log\\my-process-" + contextID + ".json",
	})
	log.ContextLoggers()[contextID+"."+taskId] = logger
	return logger
}

/*
 @param contextID string file name depends on this

 @param taskId - Allows multiple processes to write to same file;
 however if we plan to use buffered mode, then each can independently maintain its own buffer.
 Use "" for common.
*/
func MyLogger(contextID string, taskId string) *log.Logger {
	logger := log.ContextLoggers()[contextID+"."+taskId]
	if logger == nil {
		logger = InitMyLogger(contextID, taskId)
	}
	return logger
}
```
The above will ensure a safe virtual environment and buffer for each sub-task, even though they maybe writing to the same file in a process/context.

## Getting Started
See [Usage](#usage)

### Prerequisites
* [GoLang](https://golang.org/) . Tested on version `1.13`
* [You know how to set a `GOPATH`](https://github.com/golang/go/wiki/SettingGOPATH)

## Running stuff
Usual commands once your `GOPATH` is set

```bash
go get github.com/rebel-source/fmtlogwrapper -v
```

### Run the test
Executing in root folder:
```bash
go test -v fmtlogwrapper
```

#### Troubleshoot GOAPTH
Incase having issues finding running the package; can configure exclusively.
Sample GOPATH setting on windows.
```bash
set GOPATH=%GOPATH%;<path till base folder parent of src>
```

## Deployment
Just use it. Its one file ! :)


## Contributing
`TODO`

## Versioning
`TODO`

## Authors

* [Arjun Dhar](http://arjun-dhar.neurosys.biz/Arjun_Dhar.html) - *Initial work*

See also the list of [contributors](https://github.com/rebel-source/fmtlogwrapper/graphs/contributors) who participated in this project.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments
`TODO`

## Additional references
* [Logo Credit to creator @ flaticon.com](https://www.flaticon.com/free-icon/fist_128921)
