
# `fmt` `Log` Wrapper

## In one line what is it?
Log that looks like a simple `fmt.Print<xxx>` etc. Complex and extensible logging behind standard `fmt` methods.

## Usage
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

## Why do I need it?
I see developers in `golang` struggling to use and then switch between Standard GoLang [log](https://golang.org/pkg/log/), a framework like [logrus](https://github.com/sirupsen/logrus) and even writing prototype code using `fmt`. 

Would it not be nice to log using the `fmt` package standard functions like `Printf`, `Println` in a log agnostic manner that can be backed by the standard `log` package or a framework like `logrus`?


## Getting Started
See [Usage](#usage)

### Prerequisites
* [GoLang](https://golang.org/) . Tested on version `1.13`
* [You know how to set a `GOPATH`](https://github.com/golang/go/wiki/SettingGOPATH)

## Running stuff
Ususal commands once your `GOPATH` is set

```bash
go get github.com/rebel-source/fmtlogwrapper -v
```

### Run the test
Executing in root folder:
```bash
go test -v fmtlogwrappertest
```

## Deployment
Just use it. Its one file ! :)


## Contributing
`TODO`

## Versioning
`TODO`

## Authors

* [Arjun Dhar](http://arjun-dhar.neurosys.biz/Arjun_Dhar.html) - *Initial work*

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments
`TODO`

## Additional references
* [Logo Credit to creator @ flaticon.com](https://www.flaticon.com/free-icon/fist_128921)
