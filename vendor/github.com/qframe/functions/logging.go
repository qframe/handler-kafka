package qfunctions

import (
	"github.com/qframe/types/helper"
	"fmt"
	"log"
	"strings"
	"github.com/qframe/types/interfaces"
)

func Log(p qtypes_interfaces.QPlugin, logLevel, msg string) {
	_,pkg,name := p.GetInfo()

	if len(p.GetLogOnlyPlugs()) != 0 && ! qtypes_helper.IsItem(p.GetLogOnlyPlugs(), name) {
		return
	}
	// TODO: Setup in each Log() invocation seems rude
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	dL := p.CfgStringOr("log.level", "info")
	dI := qtypes_helper.LogStrToInt(dL)
	lI := qtypes_helper.LogStrToInt(logLevel)
	lMsg := fmt.Sprintf("[%+6s] %15s Name:%-10s >> %s", strings.ToUpper(logLevel), pkg, name, msg)
	if lI == 0 {
		log.Panic(lMsg)
	} else if dI >= lI {
		log.Println(lMsg)
	}
}
