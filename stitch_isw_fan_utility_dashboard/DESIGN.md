# Design System Document: Arch Utility Editorial

## 1. Overview & Creative North Star: "The Terminal Architect"
This design system is not a mere utility skin; it is a high-end editorial interpretation of the Arch Linux philosophy. The Creative North Star, **"The Terminal Architect,"** balances the raw, brutalist efficiency of the CLI with the sophisticated spatial awareness of modern luxury interfaces.

Instead of a standard "dashboard" grid, we embrace **Intentional Asymmetry**. We break the "template" look by using wide-margin gutters (Scaling 20-24) against dense data clusters. We treat system telemetry not as "labels and values," but as rhythmic data-poetry. The goal is to make the user feel like they are looking at a high-resolution blueprint of their machine—precise, deep, and authoritative.

---

## 2. Colors: Tonal Depth & The "No-Line" Rule
The palette is a deep-space descent into `#0c0d18`, punctuated by sharp, high-frequency blues.

### The "No-Line" Rule
Standard 1px solid borders are strictly prohibited for structural sectioning. Boundaries must be defined through **Background Color Shifts**. 
*   **Surface-to-Surface Transition:** Place a `surface_container_low` section directly against a `surface` background. The subtle shift in hex value is enough to cue the eye.
*   **Surface Hierarchy:** 
    *   `surface_container_lowest` (#000000) for the deepest background layers.
    *   `surface_container_highest` (#202341) for active, interactive "elevated" panels.
    *   *Nesting Tip:* A `primary_container` card should sit inside a `surface_container` to create a "recessed" look.

### The "Glass & Gradient" Rule
To escape the "flat app" feel, use **Backdrop Blur** on floating overlays.
*   **Glassmorphism:** Use `surface_variant` at 60% opacity with a `20px` backdrop blur for modals and dropdowns. This allows the system's "neon" accents to bleed through the frost.
*   **Signature Textures:** Apply a subtle linear gradient (from `primary` to `primary_container` at 15% opacity) to the background of main system hero-stats (e.g., RAM usage) to provide "visual soul."

---

## 3. Typography: Monospace Authority
We pair the geometric precision of **Space Grotesk** with the utilitarian clarity of **Inter**.

*   **Display (Space Grotesk):** Used for "System Truths"—large CPU percentages or Partition sizes. It feels engineered and modern.
*   **Headline & Title (Space Grotesk):** Commands the user's attention. Use `headline-lg` for primary navigation nodes.
*   **Body (Inter):** Reserved for descriptions and logs. It provides the "human" layer to the machine data.
*   **Label (Space Grotesk Mono-variant/Monospace):** All technical data (Hex codes, IP addresses, PID numbers) must use monospace formatting. This ensures vertical alignment and a "console" heritage.

---

## 4. Elevation & Depth: Tonal Layering
We reject traditional drop shadows in favor of **Ambient Tonal Stacking**.

*   **The Layering Principle:** Depth is achieved by stacking `surface-container` tiers. A `surface_container_high` card on a `surface_dim` background creates a natural lift.
*   **Ambient Shadows:** If an element must "float" (e.g., a critical system alert), use an extra-diffused shadow: `box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4)`. The shadow color should be a tinted version of `surface_container_lowest`.
*   **The "Ghost Border" Fallback:** For high-density data tables where separation is critical, use a **Ghost Border**: `outline_variant` (#424666) at **15% opacity**. This creates a suggestion of a container without breaking the minimalism.

---

## 5. Components: Precision Primitive

### Buttons & Interaction
*   **Primary:** A solid block of `primary` (#8bceff) with `on_primary` text. Use `sm` (0.125rem) roundedness for a sharp, technical feel.
*   **Secondary/Ghost:** No background. Use a `Ghost Border` (outline-variant @ 20%) that transitions to 100% opacity on hover.
*   **States:** On `hover`, buttons should not just lighten; they should "glow" using a `primary_dim` outer glow (4px blur).

### Cards & Lists (The Divider-Free Approach)
*   **Forbid Divider Lines:** Separate list items using the **Spacing Scale**. Use `spacing-4` (0.9rem) between items.
*   **Alternating Tones:** Use a subtle background shift (e.g., `surface_container_low` for even rows, `surface_container` for odd rows) to guide the eye.

### Input Fields
*   **The "Console" Input:** Inputs should resemble a terminal prompt. Use a bottom-only `outline` (#707396) that expands to a full-color `primary` ghost-border when focused.
*   **Monospace Data:** All text within inputs defaults to monospace to match the Arch Linux aesthetic.

### System-Specific Components
*   **Telemetry Micro-Gradients:** Use 2px tall progress bars for CPU/GPU loads using a gradient from `primary` to `primary_fixed_dim`.
*   **Kernel Terminal Chips:** Action chips for package management (`pacman -S`) should use `surface_bright` with `label-sm` monospace text.

---

## 6. Do's and Don'ts

### Do:
*   **Embrace Negative Space:** Use `spacing-16` or `spacing-20` to separate major system modules. Let the "void" of the Arch-inspired dark theme feel intentional.
*   **Align to the Data:** Ensure all numeric values are right-aligned when in columns to maintain terminal-style readability.
*   **Use Subtle Animation:** Layers should slide in with a `200ms` "ease-out-expo" to mimic a high-performance system response.

### Don't:
*   **No Pure White:** Never use `#FFFFFF`. Use `on_surface` (#e3e3ff) to maintain the "Cool Dark" atmosphere and reduce eye strain.
*   **No Heavy Rounding:** Avoid `xl` or `full` rounding for containers. Stick to `sm` (0.125rem) and `DEFAULT` (0.25rem) to preserve the "engineered" look.
*   **No Standard Grids:** Avoid placing four equal-sized cards in a row. Try a `2/3` vs `1/3` layout to create editorial tension.

---

*End of Document. This system is designed to be felt, not just seen—a digital architecture for those who master their own machines.*