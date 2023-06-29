package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type TabItemData struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	IconName fyne.ThemeIconName `json:"icon_name"`
	Disabled bool               `json:"disabled"`
}

type NavigationState struct {
	mainNavigationItems       *[]TabItemData
	parentChainTabItems       *[]TabItemData
	selectedMainNavigationTab string
}

type AppState struct {
	a   fyne.App
	w   fyne.Window
	t   SidechainTheme
	ns  NavigationState
	pcd ChainData
	scd ChainData
	pcs ChainState
	scs ChainState
}

func NewAppState(id string, title string) *AppState {
	a := app.NewWithID(id)
	w := a.NewWindow(title)
	t := NewSidechainTheme()
	a.Settings().SetTheme(t)

	return &AppState{
		a:  a,
		w:  w,
		t:  *t,
		ns: *NewNavigationState(),
	}
}

func NewNavigationState() *NavigationState {
	mainNavigationItems := []TabItemData{
		{ID: "parent_chain", Name: "Parent Chain", IconName: ParentIcon},
		{ID: "overview", Name: "Overview", IconName: HomeIcon, Disabled: true},
		{ID: "send", Name: "Send", IconName: WithdrawIcon, Disabled: true},
		{ID: "receive", Name: "Receive", IconName: DepositIcon, Disabled: true},
		{ID: "transactions", Name: "Transactions", IconName: UpDownIcon, Disabled: true},
	}

	parentChainTabItems := []TabItemData{
		{ID: "transfer", Name: "Transfer", IconName: UpDownIcon, Disabled: false},
		{ID: "withdraw_explorer", Name: "Withdraw Explorer", IconName: SearchIcon, Disabled: false},
		{ID: "bmm", Name: "BMM", IconName: MineIcon, Disabled: false},
	}

	return &NavigationState{
		mainNavigationItems:       &mainNavigationItems,
		parentChainTabItems:       &parentChainTabItems,
		selectedMainNavigationTab: "parent_chain",
	}
}
