package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type MainUI struct {
	headerContainer         *fyne.Container
	mainNavigationContainer *fyne.Container
	contentContainer        *fyne.Container
	footerContainer         *fyne.Container
	as                      *AppState
}

func NewMainUi(as *AppState) *MainUI {
	mui := &MainUI{
		headerContainer:         container.NewVBox(),
		mainNavigationContainer: container.NewHBox(),
		contentContainer:        container.NewStack(),
		footerContainer:         container.NewHBox(layout.NewSpacer()),
		as:                      as,
	}

	menus := fyne.NewMainMenu(&fyne.Menu{
		Label: "File",
		Items: []*fyne.MenuItem{
			{Label: "Open URI", Action: func() {}},
			{Label: "Backup Wallet", Action: func() {}},
			{Label: "Sign Message", Action: func() {}},
			{Label: "Verify Message", Action: func() {}},
			{Label: "", IsSeparator: true, Action: func() {}},
			{Label: "Sending Address", Action: func() {}},
			{Label: "Receiving Address", Action: func() { mui.as.a.Quit() }},
		},
	},
		&fyne.Menu{
			Label: "Tools",
			Items: []*fyne.MenuItem{
				{Label: "Hash Caclulator", Icon: mui.as.t.Icon(CalculatorIcon), Action: func() { mui.as.a.Quit() }},
				{Label: "Block Explorer", Icon: mui.as.t.Icon(SearchIcon), Action: func() { mui.as.a.Quit() }},
			},
		},
	)

	as.w.SetMainMenu(menus)

	// Setup main tab bar
	for _, item := range *as.ns.mainNavigationItems {
		itemID := item.ID
		b := widget.NewButtonWithIcon(item.Name, mui.as.t.Icon(item.IconName), func() {
			mui.SelectedMainNavigationTab(itemID)
			// setSelectedMainTab(itemID) TODO://
		})
		if item.ID == "parent_chain" {
			b.Importance = widget.MediumImportance
		} else {
			b.Importance = widget.LowImportance
		}
		if item.Disabled {
			b.Disable()
		}
		mui.mainNavigationContainer.Add(b)
	}
	mui.headerContainer.Add(container.NewPadded(mui.mainNavigationContainer))
	mui.headerContainer.Add(widget.NewSeparator())

	// Footer
	version := widget.NewRichTextWithText("Version: 0.1.0")
	version.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
		Alignment: fyne.TextAlignLeading,
		SizeName:  theme.SizeNameCaptionText,
		ColorName: theme.ColorNameForeground,
		TextStyle: fyne.TextStyle{Italic: false, Bold: false},
	}

	mui.footerContainer.Add(version)
	as.w.SetContent(container.NewStack(NewThemedRectangle(theme.ColorNameMenuBackground), container.NewPadded(container.NewBorder(mui.headerContainer, mui.footerContainer, nil, nil, container.NewPadded(mui.contentContainer)))))
	mui.SelectedMainNavigationTab("parent_chain")
	as.w.Resize(fyne.NewSize(800, 600))
	as.w.SetPadded(false)
	return mui
}

func (mui *MainUI) SelectedMainNavigationTab(id string) {
	mui.as.ns.selectedMainNavigationTab = id
	for i, item := range *mui.as.ns.mainNavigationItems {
		itemID := item.ID
		if itemID == mui.as.ns.selectedMainNavigationTab {
			mui.mainNavigationContainer.Objects[i].(*widget.Button).Importance = widget.MediumImportance
		} else {
			mui.mainNavigationContainer.Objects[i].(*widget.Button).Importance = widget.LowImportance
		}
	}
	mui.mainNavigationContainer.Refresh()
	mui.SetContainerContent(id)
}

func (mui *MainUI) SetContainerContent(id string) {
	switch id {
	case "parent_chain":
		mui.SetParentChainContent()
	case "overview":
		mui.SetOverviewContent()
	case "send":
		mui.SetSendContent()
	case "receive":
		mui.SetReceiveContent()
	case "transactions":
		mui.SetTransactionsContent()
	default:
		break
	}
	mui.contentContainer.Refresh()
}

func (mui *MainUI) SetParentChainContent() {
	mui.contentContainer.RemoveAll()

	contentBackground := NewThemedRectangle(theme.ColorNameBackground)
	contentBackground.CornerRadius = 6
	contentBackground.BorderWidth = 1
	contentBackground.BorderColorName = theme.ColorNameSeparator
	contentBackground.Refresh()
	mui.contentContainer.Add(contentBackground)

	contentBody := container.NewVBox()

	appTabs := container.NewAppTabs()
	for _, item := range *mui.as.ns.parentChainTabItems {
		appTabs.Append(container.NewTabItemWithIcon(item.Name, mui.as.t.Icon(item.IconName), widget.NewLabel(item.Name)))
	}
	contentBody.Add(container.NewPadded(appTabs))

	mui.contentContainer.Add(contentBody)
}

func (mui *MainUI) SetOverviewContent() {
	mui.contentContainer.RemoveAll()
	mui.contentContainer.Add(widget.NewLabel("Overview"))
}

func (mui *MainUI) SetSendContent() {
	mui.contentContainer.RemoveAll()
	mui.contentContainer.Add(widget.NewLabel("Send"))
}

func (mui *MainUI) SetReceiveContent() {
	mui.contentContainer.RemoveAll()
	mui.contentContainer.Add(widget.NewLabel("Receive"))
}

func (mui *MainUI) SetTransactionsContent() {
	mui.contentContainer.RemoveAll()
	mui.contentContainer.Add(widget.NewLabel("Transactions"))
}
