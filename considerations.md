# Navigation & Menu System: Design Considerations

---

**This document is the canonical reference for architectural rationale, design patterns, and anti-patterns.**
- For overall architecture and rationale, see [design.md](design.md).
- For file-by-file structure, see [structure.md](structure.md).
- For composite chat view, see [chatview.md](chatview.md).
- For modal/flow system, see [modals.md](modals.md).

---

## Table of Contents
1. [Core Principles](#core-principles)
2. [Navigation Controller](#navigation-controller)
3. [Context Object](#context-object)
4. [Menu System](#menu-system)
5. [View State Protocol](#view-state-protocol)
6. [TUI/GUI Integration](#tuigui-integration)
7. [General Efficiency & Maintainability](#general-efficiency--maintainability)
8. [New Chat Flow Implementation](#new-chat-flow-implementation)
9. [Cross-References](#cross-references)

---

## Core Principles
- **Single Responsibility:** Each component (navigation, context, view state, menu entry) has a clear, focused purpose.
- **Dependency Injection:** Context and controller are passed explicitly, decoupling logic from global state and enabling testability.
- **Immutability:** View states are immutable; navigation creates new states, supporting undo/redo and reliable stack management.
- **Type Safety:** Strong typing for all interfaces and data structures prevents runtime errors and improves maintainability.
- **Explicit State:** No hidden globals; all state is explicit and testable.

For system-wide context, see [design.md](design.md#core-architectural-principles).

---

## Navigation Controller
- **Centralized, thread-safe navigation controller (`ControllerImpl`)** manages a stack of `ViewState` objects and modal stack.
- All navigation actions (push, pop, replace, show/hide modal) are routed through this controller.
- **Command Pattern:** All navigation actions are commands routed through the controller.
- **Service/Singleton:** Controller acts as a service for navigation, but is always injected for testability.

For more, see [design.md](design.md#navigation--controller-system).

---

## Context Object
- **AppContext struct** implements a `Context` interface, providing access to app, GUI, storage, config, and logger.
- Passed to all menu actions and view states for dependency injection.
- **Interface Segregation:** Context interface can be extended/refined as needed.
- **No global variables:** All dependencies are injected, enabling easy mocking and testability.

For more, see [design.md](design.md#context--dependency-injection).

---

## Menu System
- **MenuEntry and MenuEntrySet** use a `MenuAction` function signature: `func(ctx navigation.Context, nav navigation.Controller) error`.
- All menu and submenu definitions are centralized and type-safe.
- **Data-Driven Design:** Menus are defined as data, not hardcoded logic.
- **Command Pattern:** Menu actions are commands with injected dependencies.

For more, see [design.md](design.md#menu-system).

---

## View State Protocol
- All view states implement a common `ViewState` interface, with `Update(msg, ctx, nav)` and `View()` methods.
- View states are immutable and context/controller-driven.
- **Polymorphism:** All view states are interchangeable in the navigation stack.
- **Immutability:** State transitions always create new instances.

For more, see [design.md](design.md#viewstate-protocol--ui-region-system).

---

## TUI/GUI Integration
- TUI/GUI main loop always renders the current top view state from the navigation controller.
- All input is routed to the current view state’s `Update` method, with context and controller.
- **Observer/Reactive:** UI reacts to changes in navigation state.
- **MVC/MVVM:** Clear separation of model, view, and controller logic.

For more, see [design.md](design.md#design-patterns--anti-patterns).

---

## General Efficiency & Maintainability
- All state is passed explicitly, reducing bugs and making the codebase easier to reason about.
- New features (menus, modals, view types) can be added with minimal changes to existing code.
- The system is highly testable, extensible, and ready for future growth.

---

## New Chat Flow Implementation
- The New Chat Flow is implemented using the Flow Modal System and is accessible from the Chats menu.
- Follows all core principles: single responsibility, dependency injection, immutability, and type safety.
- Navigation and state transitions are routed through the controller, ensuring extensibility and testability.
- The onSuccess handler currently displays a modal with the collected info; next step is to implement full chat creation and view integration.

For more, see [modals.md](modals.md#flow-system-multi-step-flows) and [chatview.md](chatview.md#purpose--scope).

---

## Cross-References
- [design.md](design.md): System-wide context and rationale
- [structure.md](structure.md): File-by-file structure
- [chatview.md](chatview.md): Composite chat view architecture
- [modals.md](modals.md): Modal and flow system details

---

*Sections of this document are mirrored or summarized in [design.md](design.md). For canonical details on architectural rationale and patterns, use this file. For system-wide context, always consult design.md and structure.md.* 