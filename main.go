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


	mui.as.w.SetTitle(cases.Title(language.English).String(appTitle))
	mui.as.w.ShowAndRun()
}
