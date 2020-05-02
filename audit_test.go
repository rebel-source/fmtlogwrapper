package fmtlogwrapper

// Also run this through : https://goreportcard.com/report/github.com/rebel-source/fmtlogwrapper

import (
	"testing"
)

func TestAudit(t *testing.T) {

	threadId := /*fmtlogwrapper.*/ GenSHA("SomethingThatIdenifiesThisThread") //Context id
	audit := /*fmtlogwrapper.*/ AuditLogger("AuditFunction", "audit-log", threadId, "")

	defer func() {
		/*fmtlogwrapper.*/ AuditToDiskAndClose(audit)
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
