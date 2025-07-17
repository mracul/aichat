# Core Types & Interfaces

## Update Instructions
- This file is the canonical reference for all core interfaces, types, and data structures used throughout the codebase.
- When updating, review all type/interface definitions in `/src/types/` and related files.
- For each type, provide a description, fields, and usage examples.
- Cross-reference with `structure.md`, `design.md`, and other documentation.
- Keep this file in sync with new or changed types and interfaces.

---

## Table of Contents
1. [Overview](#overview)
2. [Core Interfaces](#core-interfaces)
3. [Data Structures](#data-structures)
4. [Type Relationships](#type-relationships)
5. [Technical Diagram](#technical-diagram)
6. [Code Examples](#code-examples)
7. [References](#references)

---

## Overview

All interfaces and types are centralized in `/src/types/` to avoid import cycles and maximize modularity. Types are strongly typed and documented for maintainability. The main interfaces are ViewState, MenuAction, and Modal, each with clear responsibilities and usage patterns. See [structure.md](./structure.md#types) and [design.md](../design.md#core-architectural-principles) for context.

---

## Core Interfaces

| Interface   | Purpose                | File(s)              |
|-------------|------------------------|----------------------|
| ViewState   | UI state protocol      | `view_state.go`      |
| MenuAction  | Menu action handler    | `menu.go`            |
| Modal       | Modal protocol         | `modals.go`          |
| ...         | ...                    | ...                  |

- **ViewState**: All UI regions (menus, modals, chats, flows) implement this interface for stack-based navigation and polymorphism.
- **MenuAction**: Encapsulates menu entry actions, always context/controller-driven for testability.
- **Modal**: All modal dialogs implement this for consistent stack management and navigation.

---

## Data Structures

- **MenuEntry, MenuEntrySet**: Define menu options, actions, and navigation (see `menuentryset.go`).
- **ModalOption**: Options for modal dialogs (see `modals.go`).
- **Context, Controller**: Dependency injection and navigation (see `interfaces.go`).
- **Types for chat, flows, prompts, models**: See subfolders for domain-specific types.

---

## Type Relationships

- **ViewState** is implemented by all UI regions, menus, modals, and flows.
- **MenuAction** is used by all menu entries and actions.
- **Modal** is implemented by all modal types and flows.
- **Context/Controller** are injected everywhere for testability and decoupling.

---

## Technical Diagram

### Type Relationships and Flow
```mermaid
graph TD;
  ViewState <|-- MenuViewState
  ViewState <|-- ChatViewState
  ViewState <|-- ModalViewState
  MenuViewState -- uses --> MenuEntrySet
  MenuEntrySet -- contains --> MenuEntry
  ModalViewState -- uses --> ModalOption
  MenuViewState -- uses --> MenuAction
  ModalViewState -- uses --> Modal
  Context -- injects --> ViewState
  Controller -- injects --> ViewState
  ChatViewState --|> Subject
  Observer --|> UIComponent
  Command --|> MenuAction
```

---

## Code Examples

### 1. ViewState Interface
```go
// src/types/view_state.go#L1-10
type ViewState interface {
    Type() ViewType
    UpdateWithContext(msg tea.Msg, ctx Context, nav Controller) (tea.Model, tea.Cmd)
    View() string
    Init() tea.Cmd
}
```

### 2. MenuEntrySet Definition
```go
// src/types/menuentryset.go#L1-20
type MenuEntrySet []MenuEntry

// Example usage:
var MainMenuEntries MenuEntrySet = []MenuEntry{
    {Text: "Chats", Action: ...},
    {Text: "Favorites", Action: ...},
    // ...
}
```

### 3. Modal Interface
```go
// src/types/modals.go#L1-10
type Modal interface {
    ViewState
    ModalType() string
}
```

---

## References
- [structure.md](./structure.md#types)
- [design.md](../design.md#core-architectural-principles)
- [menus.md](./menus.md#menu-entry-definitions)
- [modals.md](./modals.md#modal-types) 

---

# Interface Segregation Structure

## ViewState Interfaces
- **Renderable**: Provides `View() string`
- **Updatable**: Provides `Update(msg tea.Msg) (tea.Model, tea.Cmd)`, `UpdateWithContext(msg tea.Msg, ctx Context, nav Controller) (tea.Model, tea.Cmd)`, `Init() tea.Cmd`
- **Serializable**: Provides `MarshalState() ([]byte, error)`, `UnmarshalState([]byte) error`
- **Navigable**: Provides `IsMainMenu() bool`, `Type() ViewType`, `ViewType() ViewType`
- **ViewState**: Embeds all above

## Modal Interfaces
- **ModalLifecycle**: Provides `OnShow()`, `OnHide()`
- **ModalRenderable**: Provides `View() string`
- **ModalUpdatable**: Provides `Update(msg interface{}, ctx Context, nav Controller) (Modal, interface{})`
- **Closable**: Provides `IsClosable() bool`, `CloseSelf()`
- **Modal**: Embeds all above and adds `Type() ModalType`

## Controller Interfaces
- **NavigationController**: Provides navigation stack methods
- **ModalController**: Provides modal management methods
- **Controller**: Embeds both

## Context Interfaces
- **AppProvider**: Provides `App() interface{}`
- **GUIProvider**: Provides `GUI() interface{}`
- **StorageProvider**: Provides `Storage() interface{}`
- **ConfigProvider**: Provides `Config() interface{}`
- **LoggerProvider**: Provides `Logger() interface{}`
- **Context**: Embeds all above

## Rationale
- Each interface is focused on a single concern, improving modularity and testability.
- Components and functions can depend on the smallest interface required.
- Aligns with Go idioms and best practices.

## Usage Notes
- Implement only the interfaces needed for each component.
- Use the smallest interface possible in function signatures for flexibility and easier testing. 

## Observer, Subject, and Command Patterns

### Observer Interface
```go
// src/types/interfaces.go#L1-10
type Observer interface {
    Notify(event interface{})
}
```

### Subject Interface
```go
// src/types/interfaces.go#L11-20
type Subject interface {
    RegisterObserver(o Observer)
    UnregisterObserver(o Observer)
    NotifyObservers(event interface{})
}
```

### Command Interface
```go
// src/types/interfaces.go#L21-30
type Command interface {
    Execute(ctx Context, nav Controller) error
}
```

### Event Struct
```go
// src/types/types.go#L1-10
type Event struct {
    Type    string
    Payload interface{}
}
```

## Factory Pattern for ViewState
- All ViewState instantiations (menus, modals, chat regions) use factory functions for consistent construction and dependency injection.
- Example:
```go
// src/components/chat/composite.go#L132-226
chatView := chat.NewCompositeChatViewStateFactory(ctx, nav, ...)
```

## Code Example: Observer in Model
```go
// src/models/chat.go#L20-40
type ChatViewState struct {
    observers []types.Observer
}
func (c *ChatViewState) RegisterObserver(o types.Observer) { ... }
func (c *ChatViewState) NotifyObservers(event interface{}) { ... }
```

## Code Example: Command Pattern
```go
// src/models/chat.go#L41-60
func (m *ChatModel) ExecuteCommand(cmd types.Command, ctx types.Context, nav types.Controller) error {
    return cmd.Execute(ctx, nav)
}
```

## Technical Diagram
```mermaid
graph TD;
  ViewState <|-- MenuViewState
  ViewState <|-- ChatViewState
  ViewState <|-- ModalViewState
  MenuViewState -- uses --> MenuEntrySet
  MenuEntrySet -- contains --> MenuEntry
  ModalViewState -- uses --> ModalOption
  MenuViewState -- uses --> MenuAction
  ModalViewState -- uses --> Modal
  Context -- injects --> ViewState
  Controller -- injects --> ViewState
  ChatViewState --|> Subject
  Observer --|> UIComponent
  Command --|> MenuAction
```

## Cross-References
- [structure.md](./structure.md#types)
- [design.md](../design.md#core-architectural-principles)
- [chatview.md](./chatview.md#focus--event-handling) 