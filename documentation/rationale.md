# Implementation Rationale & Lessons Learned

## Update Instructions
- This file is the canonical reference for the rationale behind architectural and implementation choices that have already been made in the project.
- When updating:
  - For each major decision, document the context, options considered, rationale, trade-offs, and lessons learned.
  - Reference relevant code, issues, and documentation.
  - Keep this file concise, actionable, and cross-referenced with `design.md`, `todo.md`, and other documentation.
  - Use this file for things already implemented. For forward-looking considerations and planned work, see [todo.md](./todo.md).
  - Ensure all cross-references are current and relative.
  - Do not exceed 1500 words.

---

## Table of Contents
1. [Overview](#overview)
2. [Navigation & Controller System](#navigation--controller-system)
3. [Context & Dependency Injection](#context--dependency-injection)
4. [Menu & Modal System](#menu--modal-system)
5. [Composite Chat View Architecture](#composite-chat-view-architecture)
6. [Patterns, Trade-offs, and Lessons Learned](#patterns-trade-offs-and-lessons-learned)
7. [Technical Diagram](#technical-diagram)
8. [Code Examples](#code-examples)
9. [Cross-References](#cross-references)

---

## Overview
This file documents the reasoning and lessons learned from key architectural and implementation choices in the TUI Chat Application. It is the canonical source for why the system is structured as it is. For things yet to be implemented or decided, see [todo.md](./todo.md). For system-wide context, see [design.md](../design.md) and [structure.md](./structure.md).

---

## Navigation & Controller System
- **Context & Motivation:** Needed robust, testable, and extensible navigation for stack-based TUI flows and modal overlays.
- **Options Considered:** Direct stack manipulation, global navigation state, centralized controller.
- **Rationale:** Chose a centralized, thread-safe controller (always injected, never global) for consistency, extensibility, and testability.
- **Trade-offs:** Slightly more boilerplate, but much greater maintainability and testability.
- **References:** [design.md](../design.md#navigation--controller-system), [structure.md](./structure.md#navigation), navigation code in `/src/navigation/`.

---

## Context & Dependency Injection
- **Context & Motivation:** Avoid tight coupling and global state; enable easy testing and future refactors.
- **Options Considered:** Global variables, direct field access, explicit context injection.
- **Rationale:** Chose explicit context injection everywhere for decoupling, testability, and extensibility.
- **Trade-offs:** Slightly more verbose, but enables easy mocking and future-proofing.
- **References:** [design.md](../design.md#context--dependency-injection), [structure.md](./structure.md#navigation), context code in `/src/navigation/`.

---

## Menu & Modal System
- **Context & Motivation:** Needed extensible, type-safe, and testable menu and modal logic.
- **Options Considered:** Hardcoded menus, string-based actions, data-driven menus with typed actions.
- **Rationale:** Chose data-driven, type-safe menu and modal definitions with centralized action logic for extensibility and testability.
- **Trade-offs:** More up-front design, but much easier to extend and maintain.
- **References:** [design.md](../design.md#menu-system), [menus.md](./menus.md#menu-system-design), [modals.md](./modals.md#modal-types), menu/modal code in `/src/types/` and `/src/components/`.

---

## Composite Chat View Architecture
- **Context & Motivation:** Needed modular, testable, and extensible chat UI with multiple regions and focus management.
- **Options Considered:** Monolithic chat view, region-based composite pattern, global state.
- **Rationale:** Chose composite pattern with modular ViewState regions for strict separation of concerns, testability, and future extensibility.
- **Trade-offs:** More interface definitions and orchestration logic, but much greater modularity and maintainability.
- **References:** [design.md](../design.md#composite-chat-view-architecture), [chatview.md](./chatview.md#composite-architecture), chat view code in `/src/components/chat/`.

---

## Patterns, Trade-offs, and Lessons Learned
- **Patterns Used:**
  - Stack-based navigation
  - Composite ViewState pattern
  - Data-driven menus and modals
  - Dependency injection
  - Immutability
  - Command pattern
  - Observer/reactive UI
- **Trade-offs:**
  - Explicit state and context everywhere for testability and extensibility
  - No global state or direct stack manipulation
  - Modular UI regions require more interface definitions
- **Lessons Learned:**
  - Early investment in modularity and testability pays off in maintainability
  - Centralized, data-driven patterns make future refactors and feature additions much easier
  - Clear separation of concerns is critical for large, extensible TUIs

---

## Technical Diagram

### Architectural Pattern Relationships
```mermaid
graph TD;
  Controller --> Stack
  Controller --> ModalManager
  ModalManager --> Modal
  Stack --> ViewState
  ViewState <|-- MenuViewState
  ViewState <|-- ChatViewState
  ViewState <|-- ModalViewState
  MenuViewState -- uses --> MenuEntrySet
  ModalViewState -- uses --> ModalOption
```

---

## Code Examples

### 1. Controller Push/Pop
```go
// src/navigation/controller.go#L1-50
func (c *Controller) Push(state ViewState) { ... }
func (c *Controller) Pop() ViewState { ... }
```

### 2. CompositeChatViewState
```go
// src/components/chat/composite.go#L132-226
func NewCompositeChatViewState(ctx navigation.Context, nav navigation.Controller) *CompositeChatViewState { ... }
```

### 3. MenuEntrySet
```go
// src/types/menuentryset.go#L1-20
type MenuEntrySet []MenuEntry
```

---

## Cross-References
- [todo.md](./todo.md): Forward-looking considerations and planned work
- [design.md](../design.md): System-wide context and rationale
- [structure.md](./structure.md): File-by-file structure
- [chatview.md](./chatview.md): Composite chat view architecture
- [modals.md](./modals.md): Modal and flow system details 