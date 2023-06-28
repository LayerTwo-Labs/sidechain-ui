package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	chain "sidechain-ui/chain"
	ui "sidechain-ui/ui"
	ccanvas "sidechain-ui/ui/canvas"
	ctheme "sidechain-ui/ui/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/biter777/processex"
)

//go:embed binaries/linux-testchaind
var linuxBytes []byte

//go:embed chain/default.conf
var defaultConfBytes []byte

const (
	INIT_WIN_WIDTH  = 800
	INIT_WIN_HEIGHT = 480
	DATA_DIR        = ".testchain"
	BIN_NAME        = "testchaind"
	CONF_NAME       = "testchain.conf"
)

var (
	headerContainer         *fyne.Container
	mainNavigationContainer *fyne.Container
	contentContainer        *fyne.Container
	footerContainer         *fyne.Container
)

var mainNavigationItems = []ui.TabItemData{
	{ID: "parent_chain", Name: "Parent Chain", IconName: ctheme.ParentIcon},
	{ID: "overview", Name: "Overview", IconName: ctheme.HomeIcon},
	{ID: "send", Name: "Send", IconName: ctheme.WithdrawIcon},
	{ID: "receive", Name: "Receive", IconName: ctheme.DepositIcon},
	{ID: "transactions", Name: "Transactions", IconName: ctheme.UpDownIcon},
}

var parentChainTabItems = []ui.TabItemData{
	{ID: "transfer", Name: "Transfer", IconName: ctheme.UpDownIcon},
	{ID: "withdraw_explorer", Name: "Withdraw Explorer", IconName: ctheme.SearchIcon},
	{ID: "bmm", Name: "BMM", IconName: ctheme.MineIcon},
}

var (
	a fyne.App
	w fyne.Window
	t *ctheme.SidechainTheme
)

var (
	selectedMainTabID string = "parent_chain"
	chainDir          string
	confDir           string
	chainData         chain.ChainData
)

func main() {
	dirSetup()
	a = app.NewWithID("com.ismyhc.sidechain-ui")
	w = a.NewWindow("Sidechain UI")
	w.SetPadded(false)

	t = ctheme.NewSidechainTheme()
	a.Settings().SetTheme(t)

	w.Resize(fyne.NewSize(INIT_WIN_WIDTH, INIT_WIN_HEIGHT))

	initUI()
	launchChain()

	w.ShowAndRun()
}

func setSelectedMainTab(id string) {
	selectedMainTabID = id
	for i, item := range mainNavigationItems {
		itemID := item.ID
		if itemID == selectedMainTabID {
			mainNavigationContainer.Objects[i].(*widget.Button).Importance = widget.MediumImportance
		} else {
			mainNavigationContainer.Objects[i].(*widget.Button).Importance = widget.LowImportance
			// mainTabContainer.Objects[i].(*widget.Button).Disable()
		}
	}
	mainNavigationContainer.Refresh()
	setContainerContent(id)
}

func initUI() {
	// setup menu items
	menus := fyne.NewMainMenu(&fyne.Menu{
		Label: "File",
		Items: []*fyne.MenuItem{
			{Label: "Open URI", Action: func() {}},
			{Label: "Backup Wallet", Action: func() {}},
			{Label: "Sign Message", Action: func() {}},
			{Label: "Verify Message", Action: func() {}},
			{Label: "", IsSeparator: true, Action: func() {}},
			{Label: "Sending Address", Action: func() {}},
			{Label: "Receiving Address", Action: func() { a.Quit() }},
		},
	},
		&fyne.Menu{
			Label: "Tools",
			Items: []*fyne.MenuItem{
				{Label: "Hash Caclulator", Icon: t.Icon(ctheme.CalculatorIcon), Action: func() { a.Quit() }},
				{Label: "Block Explorer", Icon: t.Icon(ctheme.SearchIcon), Action: func() { a.Quit() }},
			},
		},
	)

	w.SetMainMenu(menus)

	headerContainer = container.NewVBox()

	// Setup main tab bar
	mainNavigationContainer = container.NewHBox()
	for _, item := range mainNavigationItems {
		itemID := item.ID
		b := widget.NewButtonWithIcon(item.Name, t.Icon(item.IconName), func() {
			setSelectedMainTab(itemID)
		})
		if item.ID == "parent_chain" {
			b.Importance = widget.MediumImportance
		} else {
			b.Importance = widget.LowImportance
		}
		mainNavigationContainer.Add(b)
	}
	headerContainer.Add(container.NewPadded(mainNavigationContainer))
	headerContainer.Add(widget.NewSeparator())

	// Content area
	contentContainer = container.NewStack()

	// Footer
	version := widget.NewRichTextWithText("Version: 0.1.0")
	version.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		Alignment: fyne.TextAlignLeading,
		SizeName:  theme.SizeNameCaptionText,
		ColorName: theme.ColorNameForeground,
		TextStyle: fyne.TextStyle{Italic: false, Bold: false},
	}

	footerContainer = container.NewHBox(layout.NewSpacer(), version)
	w.SetContent(container.NewStack(ccanvas.NewThemedRectangle(theme.ColorNameMenuBackground), container.NewPadded(container.NewBorder(headerContainer, footerContainer, nil, nil, container.NewPadded(contentContainer)))))
	setSelectedMainTab(selectedMainTabID)
}

func setContainerContent(id string) {
	switch id {
	case "parent_chain":
		setParentChainContent()
	case "overview":
		setOverviewContent()
	case "send":
		setSendContent()
	case "receive":
		setReceiveContent()
	case "transactions":
		setTransactionsContent()
	default:
		break
	}
	contentContainer.Refresh()
}

func setParentChainContent() {
	contentContainer.RemoveAll()

	contentBackground := ccanvas.NewThemedRectangle(theme.ColorNameBackground)
	contentBackground.CornerRadius = 6
	contentBackground.BorderWidth = 1
	contentBackground.BorderColorName = theme.ColorNameSeparator
	contentBackground.Refresh()
	contentContainer.Add(contentBackground)

	contentBody := container.NewVBox()

	appTabs := container.NewAppTabs()
	for _, item := range parentChainTabItems {
		appTabs.Append(container.NewTabItemWithIcon(item.Name, t.Icon(item.IconName), widget.NewLabel(item.Name)))
	}
	contentBody.Add(container.NewPadded(appTabs))

	contentContainer.Add(contentBody)
}

func setOverviewContent() {
	contentContainer.RemoveAll()
	contentContainer.Add(widget.NewLabel("Overview"))
}

func setSendContent() {
	contentContainer.RemoveAll()
	contentContainer.Add(widget.NewLabel("Send"))
}

func setReceiveContent() {
	contentContainer.RemoveAll()
	contentContainer.Add(widget.NewLabel("Receive"))
}

func setTransactionsContent() {
	contentContainer.RemoveAll()
	contentContainer.Add(widget.NewLabel("Transactions"))
}

func dirSetup() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	chainDir = homeDir + string(os.PathSeparator) + DATA_DIR
	if _, err := os.Stat(chainDir); os.IsNotExist(err) {
		os.MkdirAll(chainDir, 0o755)
	}

	confDir = chainDir + string(os.PathSeparator) + CONF_NAME
	if _, err := os.Stat(confDir); os.IsNotExist(err) {
		err = os.WriteFile(confDir, defaultConfBytes, 0o755)
		if err != nil {
			log.Fatal(err)
		}
	}

	// load
	readFile, err := os.Open(confDir)
	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var fileLines []string

	for fileScanner.Scan() {
		fileLines = append(fileLines, fileScanner.Text())
	}

	readFile.Close()

	confMap := make(map[string]interface{})

	for _, line := range fileLines {
		a := strings.Split(line, "=")
		if len(a) == 2 {
			k := strings.TrimSpace(a[0])
			v := strings.TrimSpace(a[1])
			if k != "" {
				iv, err := (strconv.ParseInt(v, 0, 64))
				if err != nil {
					confMap[k] = v
				} else {
					confMap[k] = int(iv)
				}
			}
		}
	}

	jsonData, _ := json.Marshal(confMap)

	var chainData chain.ChainData
	err = json.Unmarshal(jsonData, &chainData)
	if err != nil {
		log.Fatal(err)
	}

	target := runtime.GOOS
	switch target {
	case "darwin":
		break
	case "linux":
		if _, err := os.Stat(chainDir + BIN_NAME); os.IsNotExist(err) {
			err = os.WriteFile(chainDir+BIN_NAME, linuxBytes, 0o755)
			if err != nil {
				log.Fatal(err)
			}
		}
	case "windows":
		break
	}
}

func launchChain() {
	_, err := getChainProcess()
	if err != nil {
		return
	}
}

func stopChain() {
	_, err := getChainProcess()
	if err != nil {
		return
	}
	// Shutdown chain gracefully via rpc
}

func killChain() {
	process, err := getChainProcess()
	if err != nil {
		return
	}
	process.Kill()
}

func getChainProcess() (*os.Process, error) {
	process, _, err := processex.FindByName(BIN_NAME)
	if err == processex.ErrNotFound {
		fmt.Printf("Process %v not running", BIN_NAME)
		return nil, err
	}
	if err != nil {
		fmt.Printf("Process %v find error: %v", BIN_NAME, err)
		return nil, err
	}
	fmt.Printf("Process %v PID: %v", BIN_NAME, process[0].Pid)
	if len(process) > 0 {
		return process[0], nil
	}
	return nil, fmt.Errorf("something went wrong finding process")
}
