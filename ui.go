package main

import (
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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

type ContentBodyUI interface {
	Set(mui *MainUI, c *fyne.Container)     // Set the content body by replacing the content container
	Refresh(mui *MainUI, c *fyne.Container) // Refresh the content body updating only changed values
}

func NewMainUi(as *AppState) *MainUI {
	mui := &MainUI{
		headerContainer:         container.NewVBox(),
		mainNavigationContainer: container.NewHBox(),
		contentContainer:        container.NewStack(),
		footerContainer:         container.NewHBox(),
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
	for i, item := range *as.ns.mainNavigationItems {
		index := i
		b := widget.NewButtonWithIcon(item.Name, mui.as.t.Icon(item.IconName), func() {
			mui.SelectedMainNavigationTab(index)
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
	mui.SetFooter()
	as.w.SetContent(container.NewStack(NewThemedRectangle(theme.ColorNameMenuBackground), container.NewPadded(container.NewBorder(mui.headerContainer, mui.footerContainer, nil, nil, container.NewPadded(mui.contentContainer)))))
	mui.SelectedMainNavigationTab(0)
	as.w.Resize(fyne.NewSize(800, 600))
	as.w.SetPadded(false)
	return mui
}

func (mui *MainUI) Refresh() {
	for i, item := range *mui.as.ns.mainNavigationItems {
		if i == mui.as.ns.selectedTabItemIndex {
			item.ContentBody.Refresh(mui, mui.contentContainer)
			println("Refreshed: ", item.Name)
			break
		}
	}
	mui.SetFooter()
}

func (mui *MainUI) SetFooter() {
	if mui.footerContainer.Objects == nil {
		hbox := container.NewHBox()

		version := widget.NewRichTextWithText("Version: 0.1.0")
		version.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
			Alignment: fyne.TextAlignLeading,
			SizeName:  theme.SizeNameCaptionText,
			ColorName: theme.ColorNameForeground,
			TextStyle: fyne.TextStyle{Italic: false, Bold: false},
		}

		hbox.Add(version)
		hbox.Add(widget.NewSeparator())

		blockHeight := widget.NewRichTextWithText(fmt.Sprintf("Blocks: %v", mui.as.scs.Height))
		blockHeight.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
			Alignment: fyne.TextAlignLeading,
			SizeName:  theme.SizeNameCaptionText,
			ColorName: theme.ColorNameForeground,
			TextStyle: fyne.TextStyle{Italic: false, Bold: false},
		}

		hbox.Add(blockHeight)

		mui.footerContainer.RemoveAll()
		mui.footerContainer.Add(hbox)
	} else {
		blockHeight := widget.NewRichTextWithText(fmt.Sprintf("Blocks: %v", mui.as.scs.Height))
		blockHeight.Segments[0].(*widget.TextSegment).Style = widget.RichTextStyle{
			Alignment: fyne.TextAlignLeading,
			SizeName:  theme.SizeNameCaptionText,
			ColorName: theme.ColorNameForeground,
			TextStyle: fyne.TextStyle{Italic: false, Bold: false},
		}
		// TODO: Keep reference to the blockHeight widget
		mui.footerContainer.Objects[0].(*fyne.Container).Objects[2].(*widget.RichText).Segments[0].(*widget.TextSegment).Text = fmt.Sprintf("Blocks: %v", mui.as.scs.Height)
		mui.footerContainer.Refresh()
	}
}

func (mui *MainUI) SelectedMainNavigationTab(index int) {
	if index == mui.as.ns.selectedTabItemIndex {
		return
	}
	mui.as.ns.selectedTabItemIndex = index
	for i, item := range *mui.as.ns.mainNavigationItems {
		if i == mui.as.ns.selectedTabItemIndex {
			mui.mainNavigationContainer.Objects[i].(*widget.Button).Importance = widget.MediumImportance
			item.ContentBody.Set(mui, mui.contentContainer)
		} else {
			mui.mainNavigationContainer.Objects[i].(*widget.Button).Importance = widget.LowImportance
		}
	}
	mui.mainNavigationContainer.Refresh()
}

// Build Inner UI Here
type ParentChainContentBodyUI struct {
	TabItems         *[]TabItemData
	SelectedTabIndex int
}

func NewParentChainContentBodyUI() *ParentChainContentBodyUI {
	tabItems := []TabItemData{
		{ID: "transfer", Name: "Transfer", IconName: UpDownIcon, Disabled: false, ContentBody: NewParentChainTransfersContentUI()},
		{ID: "withdraw_explorer", Name: "Withdraw Explorer", IconName: SearchIcon, Disabled: false},
		{ID: "bmm", Name: "BMM", IconName: MineIcon, Disabled: false},
	}
	return &ParentChainContentBodyUI{
		TabItems:         &tabItems,
		SelectedTabIndex: 0,
	}
}

func (pc *ParentChainContentBodyUI) Set(mui *MainUI, c *fyne.Container) {
	// Remove all content
	// Change to parent chain tab
	c.RemoveAll()

	contentBackground := NewThemedRectangle(theme.ColorNameBackground)
	contentBackground.CornerRadius = 6
	contentBackground.BorderWidth = 1
	contentBackground.BorderColorName = theme.ColorNameSeparator
	contentBackground.Refresh()
	c.Add(contentBackground)

	contentBody := container.NewVBox()

	// Set first tab in constructor. Tab bug in fyne.
	pct := *pc.TabItems
	firstTab := pct[0]
	appTabs := container.NewAppTabs(container.NewTabItemWithIcon(firstTab.Name, mui.as.t.Icon(firstTab.IconName), container.NewVBox()))
	currentTabIndex := 0
	for i, item := range pct {
		if i > 0 {
			appTabs.Append(container.NewTabItemWithIcon(item.Name, mui.as.t.Icon(item.IconName), container.NewVBox()))
		}
		if i == pc.SelectedTabIndex {
			currentTabIndex = i
		}
	}

	appTabs.OnSelected = func(item *container.TabItem) {
	}

	appTabs.SelectIndex(currentTabIndex)
	appTabs.Selected().Content = container.NewVBox()
	for i, item := range pct {
		if i == currentTabIndex {
			item.ContentBody.Set(mui, appTabs.Selected().Content.(*fyne.Container))
		}
	}

	contentBody.Add(container.NewPadded(appTabs))
	c.Add(contentBody)
}

func (pc *ParentChainContentBodyUI) Refresh(mui *MainUI, c *fyne.Container) {
	// Refresh the content body updating only changed values
	for i, item := range *pc.TabItems {
		if i == pc.SelectedTabIndex {
			item.ContentBody.Refresh(mui, c)
		}
	}
}

type ParentChainTransfersContentUI struct {
	TabItems         *[]TabItemData
	SelectedTabIndex int
	AvailableBalance *widget.Label
}

func NewParentChainTransfersContentUI() *ParentChainTransfersContentUI {
	tabItems := []TabItemData{
		{ID: "withdraw_from_sidechain", Name: "Withdraw from Sidechain", IconName: WithdrawIcon, Disabled: false, ContentBody: NewParentChainTransfersWithdrawContentUI()},
		{ID: "deposit_to_sidechain", Name: "Depsoit to Sidechain", IconName: DepositIcon, Disabled: false},
	}
	return &ParentChainTransfersContentUI{
		TabItems:         &tabItems,
		SelectedTabIndex: 0,
	}
}

func (pct *ParentChainTransfersContentUI) Set(mui *MainUI, c *fyne.Container) {
	c.RemoveAll()

	contentBody := container.NewVBox()

	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})

	pct.AvailableBalance = widget.NewLabel(fmt.Sprintf("Your sidechain balance: %s", as.scs.FormatedAvailableBalance(true)))
	// pct.AvailableBalance.TextStyle.Bold = true
	contentBody.Add(container.NewPadded(pct.AvailableBalance))

	contentBody.Add(widget.NewSeparator())

	// Set first tab in constructor. Tab bug in fyne.
	ti := *pct.TabItems
	firstTab := ti[0]
	appTabs := container.NewAppTabs(container.NewTabItemWithIcon(firstTab.Name, as.t.Icon(firstTab.IconName), container.NewVBox()))
	currentTabIndex := 0
	for i, item := range ti {
		if i > 0 {
			appTabs.Append(container.NewTabItemWithIcon(item.Name, as.t.Icon(item.IconName), container.NewVBox()))
		}
		if i == pct.SelectedTabIndex {
			currentTabIndex = i
		}
	}

	appTabs.OnSelected = func(item *container.TabItem) {
	}

	appTabs.SelectIndex(currentTabIndex)
	appTabs.Selected().Content = container.NewVBox()
	for i, item := range ti {
		if i == currentTabIndex {
			item.ContentBody.Set(mui, appTabs.Selected().Content.(*fyne.Container))
		}
	}

	contentBody.Add(container.NewPadded(appTabs))

	contentBackground := NewThemedRectangle(theme.ColorNameBackground)
	contentBackground.CornerRadius = 6
	contentBackground.BorderWidth = 1
	contentBackground.BorderColorName = theme.ColorNameSeparator
	contentBackground.Refresh()

	contentConainer := container.NewStack(contentBackground)
	contentConainer.Add(contentBody)

	c.Add(container.NewPadded(container.NewPadded(contentConainer)))
	c.Refresh()
}

func (pct *ParentChainTransfersContentUI) Refresh(mui *MainUI, c *fyne.Container) {
	pct.AvailableBalance.Text = fmt.Sprintf("Your sidechain balance: %s", as.scs.FormatedAvailableBalance(true))
	pct.AvailableBalance.Refresh()
	for i, item := range *pct.TabItems {
		if i == pct.SelectedTabIndex {
			item.ContentBody.Refresh(mui, c)
		}
	}
}

type ParentChainTransfersWithdrawContentUI struct {
	WithdrawForm *widget.Form
	Address      binding.String
	Amount       binding.String
	MainchainFee binding.String
	SidechainFee binding.String
}

func NewParentChainTransfersWithdrawContentUI() *ParentChainTransfersWithdrawContentUI {
	return &ParentChainTransfersWithdrawContentUI{
		Address:      binding.NewString(),
		Amount:       binding.NewString(),
		MainchainFee: binding.NewString(),
		SidechainFee: binding.NewString(),
	}
}

func (pctw *ParentChainTransfersWithdrawContentUI) Set(mui *MainUI, c *fyne.Container) {
	c.RemoveAll()

	contentBody := container.NewVBox()
	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})
	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})

	pctw.WithdrawForm = widget.NewForm()
	pctw.WithdrawForm.SubmitText = "Withdraw"
	pctw.WithdrawForm.OnSubmit = func() {
		sca, err := GetSidechainDepositAddress(&mui.as.scd)
		if err != nil {
			println(err.Error())
			return
		}
		adr, err := pctw.Address.Get()
		if err != nil {
			println(err.Error())
			return
		}
		err = WithdrawFromSidechain(&mui.as.scd, &mui.as.pcd, adr, *sca, 1.0, mui.as.scd.MinimumFee, mui.as.pcd.MinimumFee)
		if err != nil {
			println(err.Error())
		}
	}

	address := widget.NewEntryWithData(pctw.Address)
	address.SetPlaceHolder("Mainchain bitcoin address")
	address.Validator = func(s string) error {
		// TODO: validate value
		if len(s) == 0 {
			return errors.New("address is required")
		}
		return nil
	}

	getAddrBtn := widget.NewButton("Get Address", func() {
		a, err := GetDrivechainDepositAddress(&mui.as.pcd)
		if err != nil {
			println(err.Error())
		} else if a != nil {
			pctw.Address.Set(*a)
			address.Refresh()
		}
	})
	addri := widget.NewFormItem("Destination", container.NewBorder(nil, nil, nil, getAddrBtn, address))
	pctw.WithdrawForm.AppendItem(addri)

	amount := widget.NewEntryWithData(pctw.Amount)
	amount.Validator = func(s string) error {
		// TODO: validate value
		if len(s) == 0 {
			return errors.New("amount is required")
		}
		return nil
	}
	amount.SetPlaceHolder("0.00000000")
	pctw.WithdrawForm.Append("Amount", amount)

	mcf := widget.NewEntryWithData(pctw.MainchainFee)
	mcf.SetPlaceHolder("0.00000000")
	pctw.WithdrawForm.Append("Mainchain Fee", mcf)

	scf := widget.NewEntryWithData(pctw.SidechainFee)
	scf.SetPlaceHolder("0.00000000")
	pctw.WithdrawForm.Append("Sidechain Fee", scf)

	contentBody.Add(pctw.WithdrawForm)

	c.Add(container.NewPadded(container.NewPadded(contentBody)))
	c.Refresh()
}

func (pctw *ParentChainTransfersWithdrawContentUI) Refresh(mui *MainUI, c *fyne.Container) {
	// TODO: Autoset the drivechain deposit address?
	pctw.MainchainFee.Set(fmt.Sprintf("%.8f", mui.as.pcd.MinimumFee))
	pctw.SidechainFee.Set(fmt.Sprintf("%.8f", mui.as.scd.MinimumFee))

	pctw.WithdrawForm.Items[1].HintText = fmt.Sprintf("Max amount: %s", mui.as.scs.FormatedAvailableBalance(false))
	pctw.WithdrawForm.Items[2].HintText = fmt.Sprintf("Minimum fee: %.8f", mui.as.pcd.MinimumFee)
	pctw.WithdrawForm.Items[3].HintText = fmt.Sprintf("Minimum fee: %.8f", mui.as.scd.MinimumFee)

	pctw.WithdrawForm.Refresh()
}
