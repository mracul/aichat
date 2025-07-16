// region.go
// Defines RegionType and region order for CompositeChatViewState

package chat

// RegionType identifies each UI region in the composite chat view
// Used for focus management and region mapping

type RegionType int

const (
	SidebarTop RegionType = iota
	SidebarBottom
	ChatWindow
	InputArea
)

// RegionOrder defines the order for focus cycling
var RegionOrder = []RegionType{
	SidebarTop,
	SidebarBottom,
	ChatWindow,
	InputArea,
}

func (r RegionType) String() string {
	switch r {
	case SidebarTop:
		return "SidebarTop"
	case SidebarBottom:
		return "SidebarBottom"
	case ChatWindow:
		return "ChatWindow"
	case InputArea:
		return "InputArea"
	default:
		return "UnknownRegion"
	}
}
