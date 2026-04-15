# Design System Specification: The Sovereign Ledger

## 1. Overview & Creative North Star
The objective of this design system is to transform a high-volume wholesale environment into a high-performance editorial experience. We are moving away from the "generic SaaS dashboard" and toward a philosophy we call **"Architectural Precision."**

This system treats data not as a list to be managed, but as an asset to be displayed with authority. By utilizing intentional asymmetry, expansive white space, and a rejection of traditional containment lines, we create a sense of limitless scale and institutional trust. The interface should feel like a premium financial ledger—light, airy, yet structurally immovable.

## 2. Color & Surface Philosophy
The palette is rooted in a "Tech Blue" foundation, but its application is nuanced to avoid the visual fatigue common in B2B platforms.

### Surface Hierarchy & Nesting
To achieve a premium feel, we adhere to the **"No-Line" Rule**: 1px solid borders are prohibited for sectioning. Layouts must be defined through tonal shifts in background tokens.
*   **Base Layer:** Use `surface` (#f8f9fa) for the primary application background.
*   **The Nesting Principle:** Deeply nested content (like pricing tiers within a product card) should move up the hierarchy: `surface_container_low` → `surface_container` → `surface_container_highest`. 
*   **Glassmorphism:** For floating overlays or navigation sidebars, use `surface_container_lowest` at 80% opacity with a `20px` backdrop blur. This creates a "frosted glass" effect that keeps the user grounded in the overall context.

### Signature Textures
*   **The Momentum Gradient:** For primary Action Buttons and Hero highlights, use a linear gradient from `primary` (#005bbf) to `primary_container` (#1a73e8) at a 135-degree angle. This adds a "soul" and depth that flat hex codes cannot replicate.

## 3. Typography: The Editorial Edge
This system uses a dual-typeface strategy to balance technical utility with brand authority.

*   **Display & Headlines (Manrope):** We use Manrope for all `display` and `headline` tokens. Its geometric yet open character provides an "Institutional Modern" feel. High-contrast sizing (e.g., `display-lg` at 3.5rem) should be used for account balances and total stock counts to demand attention.
*   **Data & UI (Inter):** All `title`, `body`, and `label` roles utilize Inter. It is chosen for its superior legibility in dense data tables and pricing grids.
*   **Price Hierarchy:** Wholesale prices should always use `title-lg` with a medium weight, while "Stock Remaining" counts use `label-md` in `on_surface_variant` to provide a secondary but clear data point.

## 4. Elevation & Depth
In this system, depth is a function of light and layering, not artificial outlines.

*   **Tonal Layering:** Instead of a shadow, place a `surface_container_lowest` card on a `surface_container_low` background. The subtle 2% shift in brightness provides a more sophisticated "lift" than a drop shadow.
*   **Ambient Shadows:** Where physical elevation is required (e.g., a hovering virtual card), use an extra-diffused shadow: `offset: 0, 12px; blur: 32px; color: rgba(25, 28, 29, 0.06)`. This mimics natural, ambient light.
*   **The Ghost Border:** If accessibility requires a boundary, use the `outline_variant` token at **15% opacity**. It should be felt, not seen.

## 5. Components

### High-Performance Cards
*   **Style:** No borders. Use `surface_container_lowest` for the card body. 
*   **Radius:** `xl` (0.75rem) for external corners; `md` (0.375rem) for internal nested elements.
*   **Layout:** Use vertical whitespace (32px+) instead of dividers to separate product names from pricing tables.

### The Borderless Data Table
*   **Structure:** Forbid horizontal and vertical divider lines.
*   **Separation:** Use alternating row fills with `surface_container_low` and `surface` to guide the eye.
*   **Headers:** Use `label-sm` in all-caps with `0.05rem` letter spacing for a professional, "ticker-tape" aesthetic.

### Precision Inputs
*   **State:** Default state uses `surface_container_highest` fill with no border. 
*   **Focus:** Transition to a 2px `primary` bottom-border only, creating a clean, "architectural" focus state that doesn't box the user in.

### Buttons
*   **Primary:** Momentum Gradient (Primary to Primary-Container), white text, `md` (0.375rem) corner radius.
*   **Secondary:** `surface_container_high` fill with `on_surface` text. No border.
*   **Tertiary:** Transparent background with `on_secondary_container` text. Use for low-priority bulk actions.

## 6. Do’s and Don’ts

### Do
*   **Do** use asymmetrical layouts. For example, left-align headlines while right-aligning action groups with significant "breathing room" between them.
*   **Do** use `tertiary` (#006d2c) for "In Stock" indicators to provide a calm, trustworthy signal of availability.
*   **Do** prioritize "Tonal Nesting" over shadows whenever possible to keep the interface looking "flat-premium."

### Don't
*   **Don't** use 1px solid borders to separate sections. This is the quickest way to make the platform look like a legacy template.
*   **Don't** use pure black (#000000) for text. Always use `on_surface` (#191c1d) to maintain a soft, professional contrast.
*   **Don't** overcrowd cards. If a product card feels full, increase the spacing scale rather than shrinking the typography.