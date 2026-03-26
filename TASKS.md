# ISW Utility - Project Roadmap & Tasks

This document tracks the progress of the ISW Utility implementation, following the **Terminal Architect** design system and **Clean Architecture** principles.

## Phase 1: Architecture & Technical Foundation
- [x] **Initial Project Setup**: Configure `go.mod` and install Fyne v2.
- [ ] **Domain Definition (`/internal/domain`)**:
    - [x] Define `Telemetry` entity (CPU/GPU Temp/Load).
    - [x] Define `FanStatus` entity (RPM, Idle/Max range).
    - [x] Define `SystemInfo` entity (Kernel, Uptime, Freq).
    - [x] Define Interfaces (Ports): `SensorRepository` and `ControlRepository`.
- [ ] **Theme System Implementation**: Create a custom `fyne.Theme` matching the `#0c0d18` palette and Space Grotesk typography.

## Phase 2: Infrastructure & System Access
- [ ] **Sensor Implementation (`/internal/repository`)**:
    - [ ] Develop parser for `isw` output or `/sys/class/hwmon` data.
    - [ ] Implement `SystemInfo` provider (reading `/proc` or using `uname`).
- [ ] **Boost Mode Control**: Implement shell command execution (e.g., `isw -b on/off`) with proper privilege handling.
- [ ] **Telemetry Service (`/internal/service`)**: Implement an asynchronous polling loop to update application state without blocking the UI.

## Phase 3: UI Development (Fyne)
- [ ] **Main Layout Structure**: Implement the "Intentional Asymmetry" layout (Header, 2/3 - 1/3 grid, Analytics section).
- [ ] **Custom Widgets**:
    - [ ] **Circular Fan Gauges**: Custom Fyne component for CPU/GPU RPM.
    - [ ] **Temperature Histograms**: Bar chart component for 15-minute history.
    - [ ] **Boost Mode Switch**: Stylized toggle with "Standby/Active" state.
- [ ] **System Micro-Panel**: Display for Kernel, Uptime, and static metrics.

## Phase 4: Data Binding & Polishing
- [ ] **Data Binding**: Connect telemetry services to Fyne widgets for smooth real-time updates.
- [ ] **Animations & Effects**: Add 200ms transitions and "Glow" effects on interactions.
- [ ] **Error Handling UI**: Graceful degradation if `isw` tools are missing or permissions are denied.

## Phase 5: QA & Validation
- [ ] **Unit Testing**: Validate system data parsing and business logic.
- [ ] **UI Testing**: Verify interface responsiveness and "Boost Mode" activation.
- [ ] **Performance Audit**: Ensure low CPU/Memory footprint for the monitoring utility.
