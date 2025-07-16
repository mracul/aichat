# Enhanced Composite Chat View — Modular Architecture & Implementation Guide

---

**This document is the canonical reference for the composite chat view architecture and region orchestration.**
- For overall architecture and rationale, see [design.md](design.md).
- For file-by-file structure, see [structure.md](structure.md).
- For modal/flow system, see [modals.md](modals.md).
- For design patterns and decisions, see [considerations.md](considerations.md).

---

## Table of Contents
1. [Purpose & Scope](#purpose--scope)
2. [File Implementation Plan](#file-implementation-plan)
3. [Data Flow & Focus Management](#data-flow--focus-management)
4. [Update & View Logic](#update--view-logic)
5. [Region Modules](#region-modules)
6. [Composite Initialization](#composite-initialization)
7. [Design Patterns & Best Practices](#design-patterns--best-practices)
8. [Summary Table](#summary-table)
9. [Cross-References](#cross-references)

---

## Purpose & Scope
This document details the design, implementation, and extensibility of the composite chat view system. It is the canonical source for region orchestration, focus management, and layout logic. For system-wide context, see [design.md](design.md#composite-chat-view-architecture).

---

## File Implementation Plan
| File           | Path               | Purpose                                                           |
| -------------- | ------------------ | ----------------------------------------------------------------- |
| `composite.go` | `components/chat/` | Implements `CompositeChatViewState`, orchestrates the full layout |
| `region.go`    | `components/chat/` | Holds `RegionType`, focus state, layout ratios                    |
| `view.go`      | `components/chat/` | Main view render logic for composite layout                       |
| `update.go`    | `components/chat/` | Event delegation & focus dispatching                              |
| `layout.go`    | `components/chat/` | Responsive layout computation                                     |
| `state.go`     | `components/chat/` | State container for the composite system                          |

---

## Data Flow & Focus Management
- **CompositeChatViewState** orchestrates all four regions and manages focus, layout, and event delegation.
- **Focus cycling** is managed by navigation/controller (Tab/Shift+Tab), with explicit region focus.
- **Unidirectional data flow**: All state updates flow from input → Update → View, never mutate shared state directly.

---

## Update & View Logic
- **UpdateWithContext** delegates input to the focused region, updates state, and handles focus cycling.
- **View** composes the output of all regions using Lipgloss for layout.
- **Init** batches initialization commands for all regions.

---

## Region Modules
Each region is a modular ViewState, responsible for its own state, rendering, and input handling:
- **SidebarTopModal**: Active chat list (see `components/sidebar/chat_list.go`)
- **SidebarBottomModal**: Tabbed chat history (recent/favorites, see `components/chat/tabs.go`)
- **ChatWindowModal**: Chat transcript display (see `components/chatwindow/state.go`)
- **InputAreaModal**: Advanced text editor (see `components/input/editor.go`)

All regions are pluggable and can be extended or replaced independently.

---

## Composite Initialization
- **Init**: Each subcomponent implements `Init()` and is called in `CompositeChatViewState.Init()`.
- **Extensibility**: New panels/widgets can be added as new ViewState regions with minimal changes to orchestration logic.

---

## Design Patterns & Best Practices
- **Composite ViewState Pattern**: CompositeChatViewState orchestrates modular regions, each a ViewState.
- **Unidirectional Data Flow**: All state updates flow from input → Update → View.
- **Dependency Injection**: Context and controller are always passed, never global.
- **Immutability**: Treat models as immutable; return new structs for state changes.
- **Explicit Focus Management**: Track focus using an explicit enum or field.
- **Separation of Concerns**: Each region is responsible for its own logic and rendering.

For more on patterns and anti-patterns, see [design.md](design.md#design-patterns--anti-patterns) and [considerations.md](considerations.md).

---

## Summary Table
| Element                  | Role                                                   |
| ------------------------ | ------------------------------------------------------ |
| `CompositeChatViewState` | Orchestrator                                           |
| `SidebarTopModal`        | Active Chat List                                       |
| `SidebarBottomModal`     | Tabbed Chat History                                    |
| `ChatWindowModal`        | Chat Transcript Display                                |
| `InputAreaModal`         | Advanced Text Editor                                   |
| `Controller`             | Manages focus stack, modal overlays, region navigation |
| `Context`                | Injects storage, app state, config, etc.               |

---

## Cross-References
- [design.md](design.md#composite-chat-view-architecture): System-wide context and rationale
- [structure.md](structure.md): File-by-file structure
- [modals.md](modals.md): Modal and flow system details
- [considerations.md](considerations.md): Design patterns and decisions

---

*Sections of this document are mirrored or summarized in [design.md](design.md). For canonical details on composite chat view architecture, use this file. For system-wide context, always consult design.md and structure.md.* 