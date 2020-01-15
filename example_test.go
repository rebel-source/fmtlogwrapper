package fmtlogwrapper

import(
	"testing"
	//log "fmtlogwrapper"
)

func TestUsage1(t *testing.T) {

	// Note One can define multiple logger instances. This is the default a "singleton" instance; 
	// but nothing stopping you from defining many.if you have a specific reason to have another instance.
	// TIP: Re-use logger instances if they are outputting to the same destination (file or DB)

	/*log.*/AppLogger = /*log.*/NewLogger(/*log.*/LogSettings{
		FilePath: ".\\log\\app.log",
	})

	// Simple Printf ...& Println, Errorf  similarly
	/*log.*/AppLogger.Printf("\nRunning version %s\n", "xyz") // See how `/*log.*/AppLogger` is a simple replacement for `fmt` so Replace-ALL away :)

	//Log an Object - JSON
	logRec := make(map[string]interface{})
	logRec["key-1"] = "nothing"
	keyAsMap := make(map[string]interface{})
	keyAsMap["data"] = "ABC-123"
	keyAsMap["msg"] = "Nice to see you here :)"
	logRec["key-map"] = keyAsMap
	/*log.*/AppLogger.Log(logRec)	

	defer /*log.*/AppLogger.Close()
}