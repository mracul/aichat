// favorites.go
// Type definitions for favorites-related flows (add, remove, list, reorder)

package flows

import "aichat/types"

// AddFavoriteFlow for marking a chat as favorite
// May use a confirmation or notice modal
// Implements the strategy pattern for step logic
type AddFavoriteFlow struct {
	ConfirmModal types.ViewState
	NoticeModal  types.ViewState
	Strategy     FlowStrategy
}

// RemoveFavoriteFlow for unmarking a chat as favorite
type RemoveFavoriteFlow struct {
	ConfirmModal types.ViewState
	NoticeModal  types.ViewState
	Strategy     FlowStrategy
}

// ListFavoritesFlow for listing and selecting favorites
type ListFavoritesFlow struct {
	ListModal types.ViewState
	Strategy  FlowStrategy
}

// ReorderFavoritesFlow for reordering favorites (if supported)
type ReorderFavoritesFlow struct {
	ListModal   types.ViewState
	NoticeModal types.ViewState
	Strategy    FlowStrategy
}
