package main

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const appTitle = "Testchain"

var (
	as  *AppState
	mui *MainUI
)

func main() {
	as = NewAppState("com.ismyhc.sidechain-ui", "")
	mui = NewMainUi(as)
	ConfInit(as)

	// Launch Chain
	LaunchChain(&as.scd, &as.scs)

	// Start rpc loops
	StartSidechainStateUpdate(as, mui)

	// Intercept close so that we can shutdown
	// properly
	mui.as.w.SetCloseIntercept(func() {
		err := StopChain(&as.scd, &as.scs)
		if err != nil {
			println(err.Error())
		}
		mui.as.w.Close()
	})

	mui.as.w.SetTitle(cases.Title(language.English).String(appTitle))
	mui.as.w.ShowAndRun()
}
