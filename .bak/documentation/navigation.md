# Navigation & Controller System

## Update Instructions
- This file is the canonical reference for the navigation and controller system, including stack management and routing logic.
- When updating, review all navigation-related code in `/src/navigation/` and related interfaces in `/src/types/`.
- Summarize the architecture, patterns, and how navigation is handled throughout the app.
- Cross-reference with `design.md`, `considerations.md`, and other documentation.
- Use diagrams or tables for clarity and keep stack/command pattern documentation current.

---

## Table of Contents
1. [Overview](#overview)
2. [Controller System](#controller-system)
3. [Stack Management](#stack-management)
4. [Routing Logic](#routing-logic)
5. [Design Patterns](#design-patterns)
6. [Technical Diagram](#technical-diagram)
7. [Code Examples](#code-examples)
8. [References](#references)

---

## Overview

The navigation system is built around a stack-based controller, ensuring all UI state transitions are explicit, testable, and decoupled. The controller manages both the main view stack and a separate modal stack for overlays. See [design.md](../design.md#navigation--controller-system) and [structure.md](./structure.md#navigation) for context.

---

## Controller System

- **Location**: `navigation/controller.go`
- **Responsibilities**: Manages stack of `ViewState` objects (menus, modals, flows), routes all navigation actions (push, pop, replace, show/hide modal).
- **Thread-safety**: Always injected, never global.
- **API**: Exposes methods for stack manipulation and modal management.

**Example:**
```go
// src/navigation/controller.go#L1-50
func (c *Controller) Push(state ViewState) { ... }
func (c *Controller) Pop() ViewState { ... }
func (c *Controller) ShowModal(modalType string, message string) { ... }
```

---

## Stack Management

- **Stack-based UI**: UI states are managed as a stack for undo/redo and reliable navigation.
- **Modal overlays**: Modals are managed as a separate overlay stack using the ModalManager (`components/modals/manager.go`).
- **Implementation**: Stack is implemented in `navigation/stack.go`.
- **Persistence**: Navigation state is serializable for persistence and restoration.

**Example:**
```go
// src/navigation/stack.go#L14-81
func NewNavigationStack(main types.ViewState) *NavigationStack { ... }
func (ns *NavigationStack) Push(v types.ViewState) { ... }
func (ns *NavigationStack) Pop() types.ViewState { ... }
```

---

## Routing Logic

- **Command-based**: All navigation is command-based, using the controller as a service.
- **Context injection**: Context (`navigation/context.go`) provides dependencies to all view states.
- **Event dispatcher**: (`navigation/dispatcher.go`) routes navigation events.
- **Type safety**: Navigation actions are strongly typed and decoupled from UI logic.

---

## Design Patterns

- **Stack-based navigation**: All transitions managed via stack for overlays and undo/redo.
- **Command pattern**: Navigation actions as commands.
- **Dependency injection**: For testability and decoupling.
- **Observer/reactive UI**: UI updates via event-driven patterns.
- **Separation of concerns**: Navigation, view state, and modal management are distinct.

---

## Command Pattern & Observer Integration
- All navigation actions are implemented as commands for modular, testable transitions.
- UI updates are observer-driven, with ViewStates reacting to navigation events.
- All ViewStates are instantiated via factories for consistency.

## Code Example: Command Pattern in Navigation
```go
// src/navigation/controller.go#L1-50
type NavigateCommand struct { ... }
func (c *NavigateCommand) Execute(ctx Context, nav Controller) error { ... }
```

## Code Example: Observer Update
```go
// src/models/chat.go#L61-80
func (c *ChatViewState) Notify(event interface{}) {
    // React to navigation or model events
}
```

## Technical Diagram
```mermaid
graph TD;
  UserInput["User Input"] --> Command["Command"]
  Command --> Controller["Navigation Controller"]
  Controller --> Stack["ViewState Stack"]
  Controller --> ModalManager["ModalManager (Modal Stack)"]
  Controller --> Dispatcher["Event Dispatcher"]
  Stack --> ViewState["Current ViewState"]
  ModalManager --> Modal["Current Modal"]
  ViewState --|> Observer
```

---

## Code Examples

### 1. Controller Push/Pop
```go
// src/navigation/controller.go#L1-50
func (c *Controller) Push(state ViewState) { ... }
func (c *Controller) Pop() ViewState { ... }
```

### 2. Stack Implementation
```go
// src/navigation/stack.go#L14-81
func NewNavigationStack(main types.ViewState) *NavigationStack { ... }
func (ns *NavigationStack) Push(v types.ViewState) { ... }
func (ns *NavigationStack) Pop() types.ViewState { ... }
```

### 3. ModalManager Usage
```go
// src/components/modals/manager.go#L1-50
func (m *ModalManager) Push(modal ViewState) { ... }
func (m *ModalManager) Pop() ViewState { ... }
```

---

## References
- [design.md](../design.md#navigation--controller-system)
- [structure.md](./structure.md#navigation)
- [modals.md](./modals.md#modal-management)
- [types.md](./types.md#core-interfaces) 