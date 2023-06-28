package ui

import "fyne.io/fyne/v2"

type TabItemData struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	IconName fyne.ThemeIconName `json:"icon_name"`
}
