package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/skip2/go-qrcode"
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
			// {Label: "Open URI", Action: func() {}},
			// {Label: "Backup Wallet", Action: func() {}},
			// {Label: "Sign Message", Action: func() {}},
			// {Label: "Verify Message", Action: func() {}},
			{Label: "", IsSeparator: true, Action: func() {}},
			// {Label: "Sending Address", Action: func() {}},
			// {Label: "Receiving Address", Action: func() { mui.as.a.Quit() }},
		},
	},
	// &fyne.Menu{
	// 	Label: "Tools",
	// 	Items: []*fyne.MenuItem{
	// 		{Label: "Hash Caclulator", Icon: mui.as.t.Icon(CalculatorIcon), Action: func() { mui.as.a.Quit() }},
	// 		{Label: "Block Explorer", Icon: mui.as.t.Icon(SearchIcon), Action: func() { mui.as.a.Quit() }},
	// 	},
	// },
	)

	as.w.SetMainMenu(menus)

	// Setup main tab bar
	for i, item := range *as.ns.mainNavigationItems {
		index := i
		b := widget.NewButtonWithIcon(item.Name, mui.as.t.Icon(item.IconName), func() {
			mui.SelectedMainNavigationTab(index, false)
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
	mui.SelectedMainNavigationTab(0, false)
	as.w.Resize(fyne.NewSize(800, 600))
	as.w.SetPadded(false)
	return mui
}

func (mui *MainUI) Refresh() {
	for i, item := range *mui.as.ns.mainNavigationItems {
		if i == mui.as.ns.selectedTabItemIndex {
			item.ContentBody.Refresh(mui, mui.contentContainer)
			println("Refreshed UI TAB: ", item.Name)
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

func (mui *MainUI) SelectedMainNavigationTab(index int, forceSet bool) {
	if index == mui.as.ns.selectedTabItemIndex && !forceSet {
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
		{ID: "withdraw_explorer", Name: "Withdraw Explorer", IconName: SearchIcon, Disabled: false, ContentBody: NewParentChainWithdrawExplorerContentUI()},
		{ID: "bmm", Name: "BMM", IconName: MineIcon, Disabled: false, ContentBody: NewParentChainBMMContentUI()},
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
		si := appTabs.SelectedIndex()
		pc.SelectedTabIndex = si
		for i, it := range pct {
			if i == si {
				it.ContentBody.Set(mui, appTabs.Selected().Content.(*fyne.Container))
				it.ContentBody.Refresh(mui, appTabs.Selected().Content.(*fyne.Container))
				break
			}
		}
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
		{ID: "deposit_to_sidechain", Name: "Depsoit to Sidechain", IconName: DepositIcon, Disabled: false, ContentBody: NewParentChainTransfersDepositContentUI()},
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
		si := appTabs.SelectedIndex()
		pct.SelectedTabIndex = si
		for i, it := range ti {
			if i == si {
				it.ContentBody.Set(mui, appTabs.Selected().Content.(*fyne.Container))
				it.ContentBody.Refresh(mui, appTabs.Selected().Content.(*fyne.Container))
				break
			}
		}
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
		sca, err := GetSidechainDepositAddress(&mui.as.scd, false)
		if err != nil {
			println(err.Error())
			return
		}
		adr, err := pctw.Address.Get()
		if err != nil {
			println(err.Error())
			return
		}
		af := 0.0
		if a, err := pctw.Amount.Get(); err == nil {
			af, err = strconv.ParseFloat(a, 64)
			if err != nil {
				println(err.Error())
				return
			}
		} else {
			println(err.Error())
			return
		}
		mcf := mui.as.pcd.MinimumFee
		if umcf, err := pctw.MainchainFee.Get(); err == nil {
			if umcff, err := strconv.ParseFloat(umcf, 64); err == nil {
				mcf = umcff
			}
		}
		scf := mui.as.scd.MinimumFee
		if uscf, err := pctw.SidechainFee.Get(); err == nil {
			if uscff, err := strconv.ParseFloat(uscf, 64); err == nil {
				scf = uscff
			}
		}
		err = WithdrawFromSidechain(&mui.as.scd, &mui.as.pcd, adr, *sca, af, scf, mcf)
		if err != nil {
			println(err.Error())
		} else {
			d := widget.NewModalPopUp(widget.NewLabel("Withdrawal Sent."), mui.as.w.Canvas())
			d.Show()
			time.AfterFunc(time.Duration(2)*time.Second, func() {
				d.Hide()
				pctw.Address = binding.NewString()
				pctw.Amount = binding.NewString()
				mui.SelectedMainNavigationTab(0, true)
				mui.Refresh()
			})
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

	getAddrBtn := widget.NewButton("Get New Address", func() {
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
	amount.Validator = floatValidator
	amount.SetPlaceHolder("0.00000000")
	pctw.WithdrawForm.Append("Amount", amount)

	mcf := widget.NewEntryWithData(pctw.MainchainFee)
	mcf.Validator = floatValidator
	mcf.SetPlaceHolder("0.00000000")
	pctw.WithdrawForm.Append("Mainchain Fee", mcf)

	scf := widget.NewEntryWithData(pctw.SidechainFee)
	scf.Validator = floatValidator
	scf.SetPlaceHolder("0.00000000")
	pctw.WithdrawForm.Append("Sidechain Fee", scf)

	contentBody.Add(pctw.WithdrawForm)

	c.Add(container.NewPadded(container.NewPadded(contentBody)))
	c.Refresh()
}

func floatValidator(s string) error {
	if len(s) == 0 {
		return errors.New("amount is required")
	}
	_, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return errors.New("amount must be a number")
	}
	return nil
}

func (pctw *ParentChainTransfersWithdrawContentUI) Refresh(mui *MainUI, c *fyne.Container) {
	if a, err := pctw.Address.Get(); err == nil && a == "" {
		// Auto set the address
		if a, err := GetDrivechainDepositAddress(&mui.as.pcd); err == nil {
			pctw.Address.Set(*a)
			pctw.WithdrawForm.Refresh()
		} else {
			println(err.Error())
		}
	}

	pctw.MainchainFee.Set(fmt.Sprintf("%.8f", mui.as.pcd.MinimumFee))
	pctw.SidechainFee.Set(fmt.Sprintf("%.8f", mui.as.scd.MinimumFee))

	pctw.WithdrawForm.Items[1].HintText = fmt.Sprintf("Max amount: %s", mui.as.scs.FormatedAvailableBalance(false))
	pctw.WithdrawForm.Items[2].HintText = fmt.Sprintf("Minimum fee: %.8f", mui.as.pcd.MinimumFee)
	pctw.WithdrawForm.Items[3].HintText = fmt.Sprintf("Minimum fee: %.8f", mui.as.scd.MinimumFee)

	pctw.WithdrawForm.Refresh()
}

type ParentChainTransfersDepositContentUI struct {
	DepositEntry *widget.Entry
	Address      binding.String
	QrCode       *canvas.Image
}

func NewParentChainTransfersDepositContentUI() *ParentChainTransfersDepositContentUI {
	return &ParentChainTransfersDepositContentUI{
		Address: binding.NewString(),
	}
}

func (pctd *ParentChainTransfersDepositContentUI) Set(mui *MainUI, c *fyne.Container) {
	c.RemoveAll()
	contentBody := container.NewVBox()
	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})
	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})

	pctd.DepositEntry = widget.NewEntryWithData(pctd.Address)

	getAddrBtn := widget.NewButton("Get New Address", func() {
		a, err := GetSidechainDepositAddress(&mui.as.scd, true)
		if err != nil {
			println(err.Error())
		} else if a != nil {
			qr, err := qrcode.New(*a, qrcode.Medium)
			if err == nil {
				pctd.QrCode.Image = qr.Image(256)
				pctd.QrCode.SetMinSize(fyne.NewSize(192, 192))
				pctd.QrCode.Refresh()
				println("Set QR Code")
				contentBody.Refresh()
				c.Refresh()
			}
			pctd.Address.Set(*a)
			pctd.DepositEntry.Refresh()
		}
	})
	copyBtn := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		if a, err := pctd.Address.Get(); err == nil {
			mui.as.w.Clipboard().SetContent(a)
		}
	})

	// Auto set the address
	if a, err := GetSidechainDepositAddress(&mui.as.scd, true); err == nil {
		pctd.Address.Set(*a)
		if qr, err := qrcode.New(*a, qrcode.Medium); err == nil {
			pctd.QrCode = canvas.NewImageFromImage(qr.Image(192))
			pctd.QrCode.SetMinSize(fyne.NewSize(192, 192))
			pctd.QrCode.Refresh()
		}
		pctd.DepositEntry.Refresh()
	}

	contentBody.Add(container.NewBorder(nil, nil, nil, container.NewHBox(getAddrBtn, copyBtn), pctd.DepositEntry))
	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})
	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})
	contentBody.Add(container.NewCenter(pctd.QrCode))

	c.Add(container.NewPadded(container.NewPadded(contentBody)))
	c.Refresh()
}

func (pctd *ParentChainTransfersDepositContentUI) Refresh(mui *MainUI, c *fyne.Container) {
	// No-op
}

type ParentChainWithdrawExplorerContentUI struct{}

func NewParentChainWithdrawExplorerContentUI() *ParentChainWithdrawExplorerContentUI {
	return &ParentChainWithdrawExplorerContentUI{}
}

func (pcwe *ParentChainWithdrawExplorerContentUI) Set(mui *MainUI, c *fyne.Container) {
	c.RemoveAll()
	contentBody := container.NewVBox()
	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})
	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})
	contentBody.Add(widget.NewLabel("Withdraw Explorer"))
	c.Add(container.NewPadded(container.NewPadded(contentBody)))
	c.Refresh()
}

func (pcwe *ParentChainWithdrawExplorerContentUI) Refresh(mui *MainUI, c *fyne.Container) {
}

type ParentChainBMMContentUI struct {
	RefreshBMM     binding.Bool
	BidAmount      binding.String
	BidAmountEntry *TextRestrictedEntry
	StartBtn       *widget.Button
	StopBtn        *widget.Button
}

func NewParentChainBMMContentUI() *ParentChainBMMContentUI {
	return &ParentChainBMMContentUI{
		RefreshBMM: binding.NewBool(),
		BidAmount:  binding.NewString(),
	}
}

func (pcbmm *ParentChainBMMContentUI) Set(mui *MainUI, c *fyne.Container) {
	c.RemoveAll()

	contentBody := container.NewVBox()

	pcbmm.StartBtn = widget.NewButton(" Start ", func() {
		mui.as.scd.RefreshBMM = true
		mui.Refresh()
	})
	pcbmm.StartBtn.Importance = widget.HighImportance
	if mui.as.scd.RefreshBMM {
		pcbmm.StartBtn.Disable()
		pcbmm.StartBtn.Refresh()
	}

	pcbmm.StopBtn = widget.NewButton(" Stop ", func() {
		mui.as.scd.RefreshBMM = false
		mui.Refresh()
	})
	if !mui.as.scd.RefreshBMM {
		pcbmm.StopBtn.Disable()
		pcbmm.StartBtn.Refresh()
	}

	bal := widget.NewLabel("Bid Amount")

	entryValidator := func(text, selText string, r rune) bool {
		return (unicode.IsDigit(r) || r == '.') && len(text)-len(selText) < 10
	}

	pcbmm.BidAmount.Set("0.001")
	pcbmm.BidAmountEntry = NewTextRestrictedEntry(entryValidator)
	pcbmm.BidAmountEntry.Bind(pcbmm.BidAmount)
	pcbmm.BidAmountEntry.SetMinCharWidth(8)
	pcbmm.BidAmountEntry.SetPlaceHolder("0.00000000")
	pcbmm.BidAmountEntry.OnChanged = func(s string) {
		af, err := strconv.ParseFloat(s, 64)
		if err != nil {
			println(err.Error())
			return
		}
		mui.as.scd.BMMFee = af
	}

	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})
	contentBody.Add(container.NewPadded(container.NewPadded(container.NewHBox(pcbmm.StartBtn, pcbmm.StopBtn, widget.NewSeparator(), bal, pcbmm.BidAmountEntry))))
	contentBody.Add(widget.NewSeparator())

	contentBody.Add(&layout.Spacer{FixHorizontal: false, FixVertical: true})

	contentBody.Add(widget.NewLabel("Your attempts:"))

	// Table of bmm attempts

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

func (pcbmm *ParentChainBMMContentUI) Refresh(mui *MainUI, c *fyne.Container) {
	if mui.as.scd.RefreshBMM {
		pcbmm.StartBtn.Disable()
		pcbmm.StopBtn.Enable()
	} else {
		pcbmm.StartBtn.Enable()
		pcbmm.StopBtn.Disable()
	}
	pcbmm.StartBtn.Refresh()
	pcbmm.StopBtn.Refresh()
}
