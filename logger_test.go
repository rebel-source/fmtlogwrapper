package fmtlogwrapper

// Also run this through : https://goreportcard.com/report/github.com/rebel-source/fmtlogwrapper

import (
	"testing"
	//log "fmtlogwrapper"
	"fmt"
)

func TestUsage1(t *testing.T) {

	// Note One can define multiple logger instances. This is the default a "singleton" instance;
	// but nothing stopping you from defining many.if you have a specific reason to have another instance.
	// TIP: Re-use logger instances if they are outputting to the same destination (file or DB)

	/*log.*/
	AppLogger = /*log.*/ NewLogger( /*log.*/ LogSettings{
		FilePath: ".\\log\\app.log",
	})

	// Simple Printf ...& Println, Errorf  similarly
	/*log.*/
	AppLogger.Printf("\nRunning version %s\n", "xyz") // See how `/*log.*/AppLogger` is a simple replacement for `fmt` so Replace-ALL away :)

	//Log an Object - JSON
	logRec := make(map[string]interface{})
	logRec["key-1"] = "nothing"
	keyAsMap := make(map[string]interface{})
	keyAsMap["data"] = "ABC-123"
	keyAsMap["msg"] = "Nice to see you here :)"
	logRec["key-map"] = keyAsMap
	/*log.*/ AppLogger.Log(logRec)

	defer /*log.*/ AppLogger.Close()
}

func TestMultipleContextLoggers(t *testing.T) {
	log1 := /*log.*/ InitContextLogger("L1" /*log.*/, LogSettings{
		FilePath: ".\\log\\app-1.log",
	})

	log2 := /*log.*/ InitContextLogger("L2" /*log.*/, LogSettings{
		FilePath: ".\\log\\app-2.log",
	})

	log1.Println("STEP 1-A")
	log2.Println("STEP 2-A")

	log1.MuteWrite(true)
	log2.MuteWrite(true)
	log1.Println("STEP 1-B") // This wont be written only on console
	log2.Println("STEP 2-B") // This wont be written only on console

	// Test can get it afresh from the context
	log1 = /*log.*/ ContextLoggers()["L1"]
	log2 = /*log.*/ ContextLoggers()["L2"]

	log1.MuteWrite(false)
	log2.MuteWrite(false)
	log1.Println("STEP 1-C")
	log2.Println("STEP 2-C")

	defer /*log.*/ log1.Close()
	defer /*log.*/ log2.Close()
}

func TestBufferedLogger(t *testing.T) {
	fmt.Println("[TesBufferedLogger]")
	/*
	 We have 2 loggers writing to the same file
	 We want to ensure a Atomic operation in each does not mix with the other
	*/

	log1 := /*log.*/ InitContextLogger("L1" /*log.*/, LogSettings{
		FilePath: ".\\log\\app-same.log",
	})
	log2 := /*log.*/ InitContextLogger("L2" /*log.*/, LogSettings{
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
	// At this point buffer so far should be committed but logs for 1 & 2 should be visible in continuous lines for each group
	// We can also explicitly call log<x>.CommitBuffer()

	//Now on will appear in sequence of this being written
	log1.Println("STEP 1-D")
	log2.Println("STEP 2-D")
	log1.Println("STEP 1-E")
	log2.Println("STEP 2-E")

	defer /*log.*/ log1.Close() //If buffered was true, will also automatically CommitBuffer() any pending stuff. FYI
	defer /*log.*/ log2.Close() //Note: Multiple calls to Close() even if they share the same file is ok.
}

func TestChainLogger(t *testing.T) {
	fmt.Println("[TestChainLogger]")

	log := /*log.*/ InitContextLogger("L1" /*log.*/, LogSettings{
		FilePath: ".\\log\\app-same.log",
	})

	s := struct {
		x1 int
		x2 int
	}{
		x1: 1,
		x2: 2,
	}
	s = log.LogChain(s, nil /*formatter*/).(struct {
		x1 int
		x2 int
	})

	m := make(map[string]interface{})
	m["m1"] = 1
	m["m2"] = 2
	m = log.LogChain(m, nil /*formatter*/).(map[string]interface{})
}

// TODO: Add test for When switching from Buffered mode to non-buffered
