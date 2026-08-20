---
name: design-taste-frontend
description: Anti-slop frontend skill. Reads the brief, infers design direction, ships non-templated UI. For ERP-Go React admin use Linear-style B2B density presets.
---

# design-taste-frontend (ERP-Go vendor note)

Vendored because `npx skills add` could not clone GitHub in this environment. Upstream: Leonxlnx/taste-skill (`skills/taste-skill`).

## ERP-Go Design Read (locked)

**Reading this as: B2B ERP admin for operators, with a Linear / Vercel Dashboard language, leaning toward Tailwind + shadcn/Radix + restrained motion.**

### Dials (override Taste landing defaults)

- `DESIGN_VARIANCE: 4` — structured, not artsy
- `MOTION_INTENSITY: 3` — static-leaning; interaction feedback only
- `VISUAL_DENSITY: 8` — cockpit / packed data for tables and forms

## Hard anti-slop (always)

- Do NOT default to: AI-purple gradients, dark mesh heroes, three equal feature cards, glassmorphism everywhere, Inter + slate-900 as the whole identity
- Prefer one accent; borders and typography hierarchy over decoration
- Dashboards/tables: prioritize scanability, `tabular-nums`, clear empty states with one next action
- Respect `baseline-ui` constraints when both skills apply
- For destructive actions use AlertDialog; loading uses skeletons

## Stack alignment for this repo

- React + Vite + Tailwind + shadcn/ui in `@erp/shared-react`
- Do not invent a second component system inside app pages
- Fonts: use project CSS variables (e.g. Geist / IBM Plex), not browser defaults or Inter stacks
