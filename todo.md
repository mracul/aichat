# Project TODOs, Migration, and Planning

---

**This document is the canonical reference for current/planned features, migration steps, and technical debt.**
- For overall architecture and rationale, see [design.md](design.md).
- For file-by-file structure, see [structure.md](structure.md).
- For composite chat view, see [chatview.md](chatview.md).
- For modal/flow system, see [modals.md](modals.md).
- For design patterns and decisions, see [considerations.md](considerations.md).

---

## Table of Contents
1. [Progress & Steps](#progress--steps)
2. [Legacy Migration Plan](#legacy-migration-plan)
3. [Featureset to Implement](#featureset-to-implement)
4. [Legacy Menu Structure](#legacy-menu-structure)
5. [Modal/Prompt Features](#modalprompt-features)
6. [Navigation & Behavior](#navigation--behavior)
7. [Refactor Plan](#refactor-plan)
8. [Cross-References](#cross-references)

---

## Progress & Steps
This section tracks the major steps in refactoring, migration, and feature implementation. For architectural context, see [design.md](design.md#review--forward-planning).

---

## Legacy Migration Plan
- Extract legacy menu logic and document all menu options, actions, and transitions.
- Design new menu components as modular, reusable modals.
- Integrate main menu modal at app startup.
- Route all navigation between menus/submenus through the modal manager and update app state accordingly.
- Encapsulate all menu logic for reusability and extensibility.

For more, see [design.md](design.md#review--forward-planning) and [structure.md](structure.md).

---

## Featureset to Implement
- Multi-region chat layout (sidebar, chat window, input area)
- Menu-driven navigation (Chats, Favorites, Prompts, Models, API Keys, Help, Exit)
- Multi-step flows for chat creation, API key management, etc.
- Modal dialogs: confirmation, input, error, help, selection, etc.
- Favorites and recent chats with tabbed navigation
- Keyboard shortcuts for all major actions
- Responsive layout and accessibility features
- Error handling and recovery via modals
- Persistent storage for chats, models, prompts, keys
- Extensible plugin architecture for AI providers
- Caching and backup/restore support

For more, see [design.md](design.md#feature-set--user-experience).

---

## Legacy Menu Structure
See [design.md](design.md#menu-system) and [structure.md](structure.md#components) for canonical menu and submenu structure.

---

## Modal/Prompt Features
- Confirmation modals for actions like exit, delete, etc.
- Information modals for help/about, error messages, etc.
- Text input modals for entering chat names, prompts, model names, etc.
- Selection modals for choosing from lists (chats, models, prompts, API keys).

For more, see [modals.md](modals.md#modal-types).

---

## Navigation & Behavior
- Keyboard navigation: Up/Down arrows (or j/k), Enter to select, Esc/q/Ctrl+Q to go back or quit.
- Menus and prompts are always centered.
- ESC/back always returns to the previous menu, not quitting the app unless at the main menu.
- All submenus and modals are managed by a modal manager for stack-based navigation.

For more, see [chatview.md](chatview.md#data-flow--focus-management) and [design.md](design.md#navigation--controller-system).

---

## Refactor Plan
- Finalize migration of all legacy menu/submenu logic to new menu system.
- Modularize and document all modal types and flows.
- Clean up interface definitions and remove legacy/unused types.
- Refactor CompositeChatViewState for even stricter separation of concerns.
- Consolidate error handling and modal management logic.
- Address any technical debt flagged in this file and in [chatview.md](chatview.md).

For more, see [design.md](design.md#review--forward-planning).

---

## Cross-References
- [design.md](design.md): System-wide context and rationale
- [structure.md](structure.md): File-by-file structure
- [chatview.md](chatview.md): Composite chat view architecture
- [modals.md](modals.md): Modal and flow system details
- [considerations.md](considerations.md): Design patterns and decisions

---

*Sections of this document are mirrored or summarized in [design.md](design.md). For canonical details on project TODOs, migration, and planning, use this file. For system-wide context, always consult design.md and structure.md.* 