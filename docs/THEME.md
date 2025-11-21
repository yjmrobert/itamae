# Documentation Theme

## Overview

The Itamae documentation uses a subtle Tokyo Night dark theme based on the Hugo Relearn theme with minimal, clean customizations.

## Configuration

### Hugo Configuration (`hugo.toml`)

The documentation uses:
- `themeVariant = 'relearn-dark'` as the base dark theme
- Monokai syntax highlighting for code blocks (complements Tokyo Night)
- Disabled features: previous/next navigation, shortcuts menu, visited link tracking
- Simple, uncluttered navigation

### Tokyo Night Theme (`static/css/theme-tokyo-night.css`)

Subtle Tokyo Night color palette:
- **Background**: Deep blue-black (`#1a1b26`)
- **Text**: Soft blue-white (`#c0caf5`)
- **Primary**: Soft blue (`#7aa2f7`)
- **Accent**: Purple (`#bb9af7`)
- **Secondary**: Teal (`#73daca`)
- Minimal borders and subtle hover effects

### Custom Styling (`static/css/custom.css`)

Minimal enhancements:
- Clean, readable typography (16px base font)
- Tokyo Night colors for headings (H1: purple, H2: blue, H3: teal)
- Proper spacing and padding
- Subtle active menu indicator
- Responsive font sizing for mobile devices

## Building

To build the documentation:

```bash
cd docs
hugo --gc --minify
```

To preview locally:

```bash
cd docs
hugo server --port 1313
```

Then visit: http://localhost:1313/

## Design Principles

- **Minimal**: Stripped down to essential elements only
- **Readable**: Dark theme optimized for reduced eye strain with Tokyo Night palette
- **Fast**: No custom fonts, minimal CSS, no heavy animations
- **Clean**: Subtle color accents that don't distract from content
- **Professional**: Balanced dark theme suitable for technical documentation
