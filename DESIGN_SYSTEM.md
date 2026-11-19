# Krewire UI Design System Specification

## Overview

The Krewire UI Design System (`kui`) is an accessible, high-contrast, zero-dependency design system designed specifically for the Krewire Go ecosystem. It blends modern neo-brutalist ("pop-brutalist") aesthetics — distinct borders, crisp offset drop shadows, vibrant accents, and monospace cues — with rock-solid semantic HTML and WCAG AA/AAA accessibility.

---

## 1. Design Philosophy

1. **Own Your Components (shadcn Model)**:
   Never bury your UI in an opaque node_modules tree. `.kiw` component files live in your project's `components/` directory where you have full ownership over markup, layout, and behaviour.

2. **Zero Runtime JavaScript Bloat**:
   Interactive components prioritize native browser primitives (`<details>`, `<summary>`, `<dialog>`, semantic tags). Where JavaScript is needed (e.g. click-outside listeners), it is lightweight, framework-free vanilla JS.

3. **High-Contrast Neo-Brutalist Visuals**:
   - Signature pop-borders (`2px` to `2.5px` solid `#0B1F3B` in light, `#FAF8F4` in dark).
   - Hard offset box-shadows (`2px 2px`, `4px 4px`, `6px 6px`) without fuzzy blur.
   - Distinct surface colors that separate interactive elements from background content.

4. **Deterministic Token Inheritance**:
   All component styles inherit from single-sourced CSS variables (`--primary`, `--bg`, `--fg`, `--border`, `--radius`, `--pop-shadow`). Overriding a token at `:root` or per-container instantly cascades through all components.

---

## 2. Color System & Surface Hierarchy

### Primary Brand
- Primary (`--primary`): `#39D353` — Crisp bright emerald green.
- Primary Content (`--primary-content`): `#0B1F3B` — High-contrast midnight navy.
- Primary Hover (`--primary-hover`): `#2ecc50` (light) / `#4ae664` (dark).
- Primary Subtle (`--primary-subtle`): `rgba(57, 211, 83, 0.15)`.

### Secondary & Accents
- Secondary (`--secondary`): `#00D1C1` (vibrant cyan in light) / `#112646` (deep navy in dark).
- Accent (`--accent`): `#FF3B2E` — Crimson red alert.
- Warning (`--warning`): `#ffb703` — Amber gold.
- Error (`--error`): `#FF3B2E` — Vivid alert red.
- Success (`--success`): `#39D353` — Matching primary green.
- Info (`--info`): `#00D1C1` — Cyan information token.

### Surface Tiers
- Base 1 (`--bg`): `#FAF8F4` (Warm alabaster cream) in light mode; `#0B1F3B` (Midnight navy) in dark mode.
- Base 2 (`--bg-2`): `#f2ede4` in light mode; `#112646` in dark mode.
- Base 3 (`--bg-3`): `#e8e2d5` in light mode; `#1a3865` in dark mode.
- Foregrounds (`--fg`): `#0B1F3B` in light mode; `#FAF8F4` in dark mode.
- Muted Text (`--fg-muted`): `#57677d` in light mode; `#9ab0ce` in dark mode.

---

## 3. Typography Scale

- **Sans Stack**: `Inter, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`
- **Mono Stack**: `ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace`
- **Display Stack**: `Outfit, Inter, system-ui, sans-serif`

| Token | Size (rem / px) | Line Height | Tracking | Recommended Use |
|---|---|---|---|---|
| `--ui-text-2xs` | 0.625rem (10px) | 1.25 | 0.1em | Badges, tags, micro-labels |
| `--ui-text-xs` | 0.6875rem (11px) | 1.4 | 0.05em | Code inline, metadata, tooltips |
| `--ui-text-sm` | 0.8125rem (13px) | 1.6 | 0 | Body small, table cells, secondary text |
| `--ui-text-base` | 0.90625rem (14.5px) | 1.625 | 0 | Standard body paragraphs, inputs |
| `--ui-text-lg` | 1.0625rem (17px) | 1.5 | -0.01em | Subheadings, card titles, leads |
| `--ui-text-xl` | 1.25rem (20px) | 1.4 | -0.025em | Section subheaders, dialog titles |
| `--ui-text-2xl` | 1.5rem (24px) | 1.3 | -0.025em | Section titles, feature headers |
| `--ui-text-3xl` | 1.875rem (30px) | 1.25 | -0.03em | Major headings, hero leads |
| `--ui-text-4xl` | 2.25rem (36px) | 1.2 | -0.035em | Display headers, landing headlines |
| `--ui-text-5xl` | 3rem (48px) | 1.1 | -0.04em | Impact display headlines |

---

## 4. Spacing & Geometry Scale

Spacing adheres to an unambiguous 4px / 8px grid:

- `--ui-space-1`: `4px`
- `--ui-space-2`: `8px`
- `--ui-space-3`: `12px`
- `--ui-space-4`: `16px`
- `--ui-space-5`: `20px`
- `--ui-space-6`: `24px`
- `--ui-space-8`: `32px`
- `--ui-space-10`: `40px`
- `--ui-space-12`: `48px`
- `--ui-space-16`: `64px`

### Corner Radii
- `--ui-radius-xs`: `4px` (subtle badges, kbd keys)
- `--ui-radius-sm`: `6px` (tags, buttons small, dropdown items)
- `--ui-radius`: `10px` (standard buttons, inputs, alerts)
- `--ui-radius-md`: `12px` (icon boxes, tabs)
- `--ui-radius-lg`: `14px` (cards, dialogs, popovers)
- `--ui-radius-round`: `9999px` (pills, avatars, switch toggles)

---

## 5. Neo-Brutalist Pop Tokens

The signature Krewire visual fingerprint:

```css
--pop-border: 2.5px solid var(--border);
--pop-border-thin: 2px solid var(--border);
--pop-shadow: 4px 4px 0px var(--border);
--pop-shadow-sm: 2px 2px 0px var(--border);
--pop-shadow-lg: 6px 6px 0px var(--border);
```

### Hover & Active State Physics
- Standard clickable surfaces translate `-1.5px, -1.5px` or `-2px, -2px` on `:hover` and grow the drop shadow from `4px 4px` to `6px 6px`.
- On `:active`, surfaces translate back to `0px, 0px` (or `+1px, +1px`) with a compressed `1px 1px` shadow to deliver tactile tactile button feedback.

---

## 6. Accessibility & Keyboard Discipline

- **Visible Focus Rings**: All interactive controls carry `:focus-visible` with `outline: 2px solid var(--primary)` and `outline-offset: 2px`.
- **Keyboard Navigation**:
  - `Dropdown` & `Popover`: Toggle with `Space` or `Enter`; dismiss with `Escape`.
  - `SearchInput`: Trigger focus globally with `/` or `⌘K`.
  - `Modal`: Closes on `Escape` key press.
  - `Tabs`: Accessible tab navigation with `role="tab"` and `role="tabpanel"`.
- **Reduced Motion Support**:
  Media query `@media (prefers-reduced-motion: reduce)` disables animations and smooth scrolling for sensitive users.
