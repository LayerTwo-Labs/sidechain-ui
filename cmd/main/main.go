package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ismyhc/sidechain-ui/internal/ui"
	icanvas "github.com/ismyhc/sidechain-ui/internal/ui/canvas"
	itheme "github.com/ismyhc/sidechain-ui/internal/ui/theme"
)

const (
	INIT_WIN_WIDTH  = 800
	INIT_WIN_HEIGHT = 480
)

var (
	headerContainer  *fyne.Container
	mainTabContainer *fyne.Container
	contentContainer *fyne.Container
	footerContainer  *fyne.Container
)

var mainTabItems = []ui.TabItemData{
	{ID: "parent_chain", Name: "Parent Chain", IconName: itheme.ParentIcon},
	{ID: "overview", Name: "Overview", IconName: itheme.HomeIcon},
	{ID: "send", Name: "Send", IconName: itheme.WithdrawIcon},
	{ID: "receive", Name: "Receive", IconName: itheme.DepositIcon},
	{ID: "transactions", Name: "Transactions", IconName: itheme.UpDownIcon},
}

var secondaryTabItems = map[string][]ui.TabItemData{
	"parent_chain": {
		{ID: "transfer", Name: "Transfer", IconName: itheme.UpDownIcon},
		{ID: "withdraw_explorer", Name: "Withdraw Explorer", IconName: itheme.SearchIcon},
		{ID: "bmm", Name: "BMM", IconName: itheme.MineIcon},
	},
}

var (
	a fyne.App
	w fyne.Window
	t *itheme.SidechainTheme
)

var selectedMainTabID string = "parent_chain"

func main() {
	a = app.NewWithID("com.ismyhc.sidechain-ui")
	w = a.NewWindow("Sidechain UI")
	w.SetPadded(false)

	t = itheme.NewSidechainTheme()
	a.Settings().SetTheme(t)

	w.Resize(fyne.NewSize(INIT_WIN_WIDTH, INIT_WIN_HEIGHT))

	mm := fyne.NewMainMenu(&fyne.Menu{
		Label: "File",
		Items: []*fyne.MenuItem{
			{Label: "Quit", Action: func() { a.Quit() }},
		},
	})

	w.SetMainMenu(mm)

	initUI()

	w.ShowAndRun()
}

func setSelectedMainTab(id string) {
	selectedMainTabID = id
	for i, item := range mainTabItems {
		itemID := item.ID
		if itemID == selectedMainTabID {
			mainTabContainer.Objects[i].(*widget.Button).Importance = widget.MediumImportance
		} else {
			mainTabContainer.Objects[i].(*widget.Button).Importance = widget.LowImportance
			// mainTabContainer.Objects[i].(*widget.Button).Disable()
		}
	}
	mainTabContainer.Refresh()
}

func initUI() {
	headerContainer = container.NewVBox()

	// Setup main tab bar
	mainTabContainer = container.NewHBox()
	for _, item := range mainTabItems {
		itemID := item.ID
		b := widget.NewButtonWithIcon(item.Name, t.Icon(item.IconName), func() {
			setSelectedMainTab(itemID)
			setSecondaryContent(itemID)
		})
		if item.ID == "parent_chain" {
			b.Importance = widget.MediumImportance
		} else {
			b.Importance = widget.LowImportance
		}
		mainTabContainer.Add(b)
	}
	headerContainer.Add(mainTabContainer)
	setSelectedMainTab(selectedMainTabID)

	// Content area
	contentContainer = container.NewStack()

	contentBackground := icanvas.NewThemedRectangle(theme.ColorNameBackground)
	contentBackground.CornerRadius = 6
	contentBackground.BorderWidth = 1
	contentBackground.BorderColorName = theme.ColorNameSeparator
	contentBackground.Refresh()

	contentContainer.Add(contentBackground)

	setSecondaryContent(selectedMainTabID)

	// Footer
	version := widget.NewRichTextWithText("Version: 0.1.0")
	version.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		Alignment: fyne.TextAlignLeading,
		SizeName:  theme.SizeNameCaptionText,
		ColorName: theme.ColorNameForeground,
		TextStyle: fyne.TextStyle{Italic: false, Bold: false},
	}

	footerContainer = container.NewHBox(layout.NewSpacer(), version)
	w.SetContent(container.NewStack(icanvas.NewThemedRectangle(theme.ColorNameMenuBackground), container.NewPadded(container.NewBorder(headerContainer, footerContainer, nil, nil, contentContainer))))
}

func setSecondaryContent(id string) {
	// if _, ok := secondaryTabItems[id]; !ok {

	// secondaryTabContainer.Items = nil
	// for _, item := range secondaryTabItems[id] {
	// 	secondaryTabContainer.Append(container.NewTabItemWithIcon(item.Name, t.Icon(item.IconName), widget.NewLabel(item.Name)))
	// }
	contentContainer.Refresh()
}
