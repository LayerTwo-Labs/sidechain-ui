package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type TabItemData struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	IconName    fyne.ThemeIconName `json:"icon_name"`
	Disabled    bool               `json:"disabled"`
	ContentBody ContentBodyUI      `json:"content_body"`
}

type NavigationState struct {
	mainNavigationItems  *[]TabItemData
	selectedTabItemIndex int
}

type AppState struct {
	a       fyne.App
	w       fyne.Window
	t       SidechainTheme
	ns      NavigationState
	pcd     ChainData
	scd     ChainData
	pcs     ChainState
	scs     ChainState
	scbmmtd []BMMTableItem
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
		{ID: "parent_chain", Name: "Parent Chain", IconName: ParentIcon, Disabled: false, ContentBody: NewParentChainContentBodyUI()},
		{ID: "overview", Name: "Overview", IconName: HomeIcon, Disabled: true},
		{ID: "send", Name: "Send", IconName: WithdrawIcon, Disabled: true},
		{ID: "receive", Name: "Receive", IconName: DepositIcon, Disabled: true},
		{ID: "transactions", Name: "Transactions", IconName: UpDownIcon, Disabled: true},
	}

	return &NavigationState{
		mainNavigationItems:  &mainNavigationItems,
		selectedTabItemIndex: -1, // No tab selected
	}
}
