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

func TestMultipleContextLoggers(t *testing.T) {
	log1 := /*log.*/InitContextLogger("L1", /*log.*/LogSettings{
		FilePath: ".\\log\\app-1.log",
	})

	log2 := /*log.*/InitContextLogger("L2", /*log.*/LogSettings{
		FilePath: ".\\log\\app-2.log",
	})

	log1.Println("STEP 11")
	log2.Println("STEP 12")

	log1.MuteWrite(true)
	log2.MuteWrite(true)
	log1.Println("STEP 21")  // This wont be written only on console
	log2.Println("STEP 22")  // This wont be written only on console

	// Test can get it afresh from the context
	log1 = /*log.*/ContextLoggers()["L1"]
	log2 = /*log.*/ContextLoggers()["L2"]

	log1.MuteWrite(false)
	log2.MuteWrite(false)
	log1.Println("STEP 31")
	log2.Println("STEP 32")	

	defer /*log.*/log1.Close()
	defer /*log.*/log2.Close()
}