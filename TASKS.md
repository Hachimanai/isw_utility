# ISW Utility - Project Roadmap & Tasks

This document tracks the progress of the ISW Utility implementation, following the **Terminal Architect** design system and **Clean Architecture** principles.

## Phase 1: Architecture & Technical Foundation
- [x] **Initial Project Setup**: Configure `go.mod` and install Fyne v2.
- [x] **Domain Definition (`/internal/domain`)**:
    - [x] Define `Telemetry` entity (CPU/GPU Temp/Load).
    - [x] Define `FanStatus` entity (RPM, Idle/Max range).
    - [x] Define `SystemInfo` entity (Kernel, Uptime, Freq).
    - [x] Define Interfaces (Ports): `SensorRepository` and `ControlRepository`.
- [x] **Theme System Implementation**: Create a custom `fyne.Theme` matching the `#0c0d18` palette and Space Grotesk typography.

## Phase 2: Infrastructure & System Access
- [x] **Sensor Implementation (`/internal/repository`)**:
    - [x] Develop parser for `isw` output or `/sys/class/hwmon` data.
    - [x] Implement `SystemInfo` provider (reading `/proc` or using `uname`).
- [x] **Boost Mode Control**: Implement shell command execution (e.g., `isw -b on/off`) with proper privilege handling.
- [x] **Telemetry Service (`/internal/service`)**: Implement an asynchronous polling loop to update application state without blocking the UI.

## Phase 3: UI Development (Fyne)
- [x] **Main Layout Structure**: Implement the "Intentional Asymmetry" layout (Header, 2/3 - 1/3 grid, Analytics section).
- [x] **Custom Widgets**:
    - [x] **Circular Fan Gauges**: Custom Fyne component for CPU/GPU RPM (implemented as linear gauges for now).
    - [x] **Temperature Histograms**: Bar chart component for history.
    - [x] **Boost Mode Switch**: Stylized toggle with "Standby/Active" state.
- [x] **System Micro-Panel**: Display for Kernel, Uptime, and static metrics.

## Phase 4: Data Binding & Polishing
- [x] **Data Binding**: Connect telemetry services to Fyne widgets for smooth real-time updates.
- [x] **Animations & Effects**: Add 200ms transitions (via LERP) and "Glow" effects on interactions.
- [x] **Error Handling UI**: Status bar for reporting errors (isw missing, permissions).

## Phase 5: QA & Validation
- [x] **Unit Testing**: Validated system data parsing and fallback logic.
- [x] **UI Testing**: Verified interface responsiveness, "Boost Mode" activation, and transition effects.
- [x] **Performance Audit**: Optimized animation loop (30ms ticker) with low CPU footprint.

## Final Notes
- **Circular Fan Gauges**: Successfully implemented using `canvas.Arc` (available in Fyne v2.7+).
- **Resilience**: A `CompositeRepository` has been implemented to allow the application to function with or without the `isw` tool, falling back to `/sys/class/hwmon`, `/proc`, and `nvidia-smi`.
- **NVIDIA Support**: Added direct telemetry for NVIDIA GPUs via `nvidia-smi`.
- **EC Communication**: Fixed `isw` command arguments (added `MSI_ADDRESS_DEFAULT`) and increased authentication timeouts for a seamless "Boost Mode" experience.
- **Design System**: Strict adherence to the "Terminal Architect" aesthetic.
