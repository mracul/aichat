package mvc
package mvc

import (
	"aichat/models"
	"aichat/types"
	"aichat/types/render"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

type ChatView struct {
	State       *models.ChatViewState
	MessageList list.Model
	Help        help.Model
	Paginator   paginator.Model
	ShowHelp    bool
	ThemeMap    render.ThemeMap
	Strategies  map[string]render.RenderStrategy
	Spinner     spinner.Model
}

type chatMessageDelegate struct {
	ThemeMap   render.ThemeMap
	Strategies map[string]render.RenderStrategy
}

func (d chatMessageDelegate) Height() int  { return 2 }
func (d chatMessageDelegate) Spacing() int { return 1 }
func (d chatMessageDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}
func (d chatMessageDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	chatItem, ok := item.(models.ChatMessageItem)
	if !ok {
		return
	}
	var content string
	if !chatItem.Msg.IsUser {
		rendered, err := glamour.Render(chatItem.Msg.Content, "dark")
		if err != nil {
			content = chatItem.Msg.Content
		} else {
			content = rendered
		}
	} else {
		content = chatItem.Msg.Content
	}
	prefix := "  "
	if index == m.Index() {
		prefix = "> "
	}
	strategy := d.Strategies["message"]
	theme := d.ThemeMap[strategy.ThemeKey]
	styled := render.ApplyStrategy(prefix+content, strategy, theme)
	w.Write([]byte(styled))
}

// Observer pattern: implement types.Observer
func (v *ChatView) Notify(event interface{}) {
	if ev, ok := event.(types.Event); ok {
		switch ev.Type {
		case "message_added":
			if msg, ok := ev.Payload.(models.ChatMessage); ok {
				v.MessageList.InsertItem(len(v.State.Messages)-1, models.ChatMessageItem{Msg: msg})
			}
		case "chat_title_changed":
			v.MessageList.Title = v.State.ChatTitle
		}
	}
}

func NewChatView(state *models.ChatViewState, themeMap render.ThemeMap, strategies map[string]render.RenderStrategy) *ChatView {
	items := make([]list.Item, len(state.Messages))
	for i, msg := range state.Messages {
		items[i] = models.ChatMessageItem{Msg: msg}
	}
	sp := spinner.New()
	pg := paginator.New()
	pg.PerPage = 10
	hp := help.New()
	delegate := chatMessageDelegate{ThemeMap: themeMap, Strategies: strategies}
	l := list.New(items, delegate, 60, 20)
	l.Title = state.ChatTitle
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)
	l.SetFilteringEnabled(false)
	view := &ChatView{
		State:       state,
		MessageList: l,
		Help:        hp,
		Paginator:   pg,
		ShowHelp:    false,
		ThemeMap:    themeMap,
		Strategies:  strategies,
		Spinner:     sp,
	}
	state.RegisterObserver(view)
	return view
}

func (v *ChatView) View() string {
	view := ""
	if v.State.WaitingForResponse {
		spinnerStr := v.Spinner.View() + " Waiting for response..."
		strategy := v.Strategies["spinner"]
		theme := v.ThemeMap[strategy.ThemeKey]
		view += render.ApplyStrategy(spinnerStr, strategy, theme) + "\n"
	}
	msgStrategy := v.Strategies["messageList"]
	msgTheme := v.ThemeMap[msgStrategy.ThemeKey]
	view += render.ApplyStrategy(v.MessageList.View(), msgStrategy, msgTheme)
	if v.State.Streaming {
		streamingStr := "[Streaming...]"
		strategy := v.Strategies["streaming"]
		theme := v.ThemeMap[strategy.ThemeKey]
		view += "\n" + render.ApplyStrategy(streamingStr, strategy, theme) + "\n"
	}
	if v.ShowHelp {
		helpStr := v.Help.View(nil)
		strategy := v.Strategies["help"]
		theme := v.ThemeMap[strategy.ThemeKey]
		view += "\n" + render.ApplyStrategy(helpStr, strategy, theme) + "\n"
	}
	pagStr := v.Paginator.View()
	pagStrategy := v.Strategies["paginator"]
	pagTheme := v.ThemeMap[pagStrategy.ThemeKey]
	view += "\n" + render.ApplyStrategy(pagStr, pagStrategy, pagTheme) + "\n"
	return view
}

