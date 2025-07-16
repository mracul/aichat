# TUI Chat Application — In-Depth Design & Architecture Guide

---

## Table of Contents
1. [Project Overview](#project-overview)
2. [File & Directory Structure](#file--directory-structure)
3. [Core Architectural Principles](#core-architectural-principles)
4. [Navigation & Controller System](#navigation--controller-system)
5. [Context & Dependency Injection](#context--dependency-injection)
6. [ViewState Protocol & UI Region System](#viewstate-protocol--ui-region-system)
7. [Menu System](#menu-system)
8. [Modal & Flow System](#modal--flow-system)
9. [Composite Chat View Architecture](#composite-chat-view-architecture)
10. [Extensibility & OOP Patterns](#extensibility--oop-patterns)
11. [Design Patterns & Anti-Patterns](#design-patterns--anti-patterns)
12. [Feature Set & User Experience](#feature-set--user-experience)
13. [Review & Forward Planning](#review--forward-planning)
14. [References & Cross-Links](#references--cross-links)

---

## Project Overview
A modern, extensible, and robust terminal-based chat application. The system is built for modularity, testability, and future growth, with a focus on:
- Multi-region, keyboard-driven TUI
- Stack-based navigation and modal management
- Modular flows and advanced error handling
- OOP and dependency injection throughout

---

## File & Directory Structure
**Canonical: [structure.md](documentation/structure.md)**

See [structure.md](documentation/structure.md) for a full, up-to-date file-by-file breakdown, directory tree, and relationships. This file is the canonical reference for the `/src/` directory, with detailed descriptions and cross-references to all major architectural areas.

---

## Core Architectural Principles
**Canonical: [considerations.md](documentation/considerations.md)**

See [considerations.md](documentation/considerations.md) for the rationale behind all major design decisions, including single responsibility, dependency injection, immutability, type safety, and explicit state. This file also documents anti-patterns avoided and the reasoning behind architectural trade-offs.

---

## Navigation & Controller System
**Canonical: [navigation.md](documentation/navigation.md)**

See [navigation.md](documentation/navigation.md) for a detailed explanation of the stack-based navigation/controller system, including stack management, routing logic, and command patterns. All navigation is explicit, testable, and decoupled, with modals managed as a separate overlay stack.

---

## Context & Dependency Injection
**Canonical: [considerations.md](documentation/considerations.md)**

Dependency injection is handled via the AppContext, which provides access to app, GUI, storage, config, and logger. All dependencies are injected, never global, enabling easy mocking and testability. See [considerations.md](documentation/considerations.md#context-object).

---

## ViewState Protocol & UI Region System
**Canonical: [types.md](documentation/types.md), [chatview.md](documentation/chatview.md)**

All UI regions, menus, modals, and flows implement the `ViewState` interface, enabling polymorphism and stack-based navigation. Each region is modular and testable. See [types.md](documentation/types.md) for interface definitions and [chatview.md](documentation/chatview.md) for region orchestration.

---

## Menu System
**Canonical: [menus.md](documentation/menus.md)**

Menus are defined as data, not hardcoded logic, and support dynamic, context-aware logic. All menu and submenu definitions are centralized and type-safe. See [menus.md](documentation/menus.md) for menu entry definitions, dynamic menus, and action patterns.

---

## Modal & Flow System
**Canonical: [modals.md](documentation/modals.md), [flows.md](documentation/flows.md)**

All modals (confirmation, input, notice, error, selection, help, about, editor, custom) implement the Modal interface and ViewState. Multi-step flows are sequences of modals with state, onExit, and onSuccess handlers. See [modals.md](documentation/modals.md) for modal types and management, and [flows.md](documentation/flows.md) for flow definitions and orchestration.

---

## Composite Chat View Architecture
**Canonical: [chatview.md](documentation/chatview.md)**

The chat view is composed of modular regions (sidebar, chat window, input area), each a ViewState. Focus management, event delegation, and responsive layout are handled by the CompositeChatViewState. See [chatview.md](documentation/chatview.md) for region orchestration and extensibility.

---

## Extensibility & OOP Patterns
**Canonical: [considerations.md](documentation/considerations.md), [types.md](documentation/types.md)**

All interfaces are centralized in `types/` to avoid import cycles and maximize modularity. New features are added by defining new structs implementing the relevant interface and registering with the controller/modal manager. See [considerations.md](documentation/considerations.md) and [types.md](documentation/types.md) for extensibility patterns.

---

## Design Patterns & Anti-Patterns
**Canonical: [considerations.md](documentation/considerations.md), [chatview.md](documentation/chatview.md), [modals.md](documentation/modals.md)**

See [considerations.md](documentation/considerations.md) for a summary of all design patterns (stack-based navigation, composite pattern, dependency injection, immutability, command pattern, observer/reactive) and anti-patterns avoided.

---

## Feature Set & User Experience
**Canonical: [todo.md](documentation/todo.md), [structure.md](documentation/structure.md)**

See [todo.md](documentation/todo.md) for the current and planned feature set, technical debt, and migration steps. The feature set includes multi-region chat layout, menu-driven navigation, multi-step flows, modal dialogs, favorites, keyboard shortcuts, responsive layout, error handling, persistent storage, plugin architecture, and backup/restore support.

---

## Review & Forward Planning
**Canonical: [todo.md](documentation/todo.md), [considerations.md](documentation/considerations.md), [chatview.md](documentation/chatview.md)**

See [todo.md](documentation/todo.md) and [considerations.md](documentation/considerations.md) for the latest review, refactor plans, and forward-looking topics. This includes migration of legacy logic, modularization, improved error handling, accessibility, plugin APIs, performance optimization, and testing/CI improvements.

---

## References & Cross-Links
- [structure.md](documentation/structure.md): File-by-file structure and relationships
- [chatview.md](documentation/chatview.md): Composite chat view architecture and region orchestration
- [modals.md](documentation/modals.md): Modal and flow system details
- [flows.md](documentation/flows.md): Multi-step flow definitions
- [menus.md](documentation/menus.md): Menu system design and action patterns
- [types.md](documentation/types.md): Core interfaces, types, and data structures
- [errors.md](documentation/errors.md): Error handling strategies and recovery flows
- [considerations.md](documentation/considerations.md): Architectural rationale and design patterns
- [todo.md](documentation/todo.md): Current/planned features, migration steps, technical debt
- [navigation.md](documentation/navigation.md): Navigation/controller system and stack management

---

*This document is a living reference. For canonical details, always consult the referenced support files. All sections are kept in sync with the latest code and documentation.* 