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