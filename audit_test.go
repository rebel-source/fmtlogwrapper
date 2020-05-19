package fmtlogwrapper

// Also run this through : https://goreportcard.com/report/github.com/rebel-source/fmtlogwrapper

import (
	"testing"

	"fmt"
	"strconv"
	"sync"
	"time"
)

func TestAudit(t *testing.T) {

	threadId := /*fmtlogwrapper.*/ GenSHA("SomethingThatIdenifiesThisThread") //Context id
	audit := /*fmtlogwrapper.*/ AuditLogger("AuditFunction", "audit-log", threadId, "")

	defer func() {
		/*fmtlogwrapper.*/ Persist(audit)
	}()

	auditData := make(map[string]interface{})
	auditData["start"] = /*fmtlogwrapper.*/ Now()
	//...
	auditData["end"] = /*fmtlogwrapper.*/ Now()
	auditData["result"] = "some result"

	audit.Println(Audit(auditData))
	audit.Println("[end]")

	return
}

func TestOpenCloseAudit(t *testing.T) {
	fmt.Println("\n[TestOpenCloseAudit]")

	var wg sync.WaitGroup

	for x := 1; x < 160; x++ {
		xStr := strconv.Itoa(x)
		wg.Add(1)
		go func() {
			fmt.Println("\n[TestOpenCloseAudit] Starting Thread " + xStr)
			log := /*fmtlogwrapper.*/ AuditLogger("AuditFunction", "audit-log", xStr, "")
			//log.WriteToBuffer(true)
			log.Println("[A" + xStr + "] A")
			time.Sleep(5 * time.Millisecond)
			log.Println("[A" + xStr + "] B")
			time.Sleep(5 * time.Millisecond)

			log = /*log.*/ ContextLoggers()[getContextId("AuditFunction", "audit-log", xStr)]
			log = InitAuditLogger("AuditFunction", "audit-log", xStr, "")

			log.Println("[A" + xStr + "] C")
			time.Sleep(5 * time.Millisecond)
			log.Println("[A" + xStr + "] D")
			time.Sleep(5 * time.Millisecond)
			log.Println("[A" + xStr + "] E")
			time.Sleep(5 * time.Millisecond)
			log.Println("[A" + xStr + "] F")
			time.Sleep(5 * time.Millisecond)

			Persist(log)

			wg.Done()
		}()
	} //end-for

	wg.Wait()
}
