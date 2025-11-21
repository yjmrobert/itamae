# Documentation Theme

## Overview

The Itamae documentation uses a minimal, clean theme based on the Hugo Relearn theme with simple customizations for readability.

## Configuration

### Hugo Configuration (`hugo.toml`)

The documentation uses:
- `themeVariant = 'relearn-light'` for a clean, light theme
- GitHub-style syntax highlighting for code blocks
- Disabled features: previous/next navigation, shortcuts menu, visited link tracking
- Simple, uncluttered navigation

### Custom Styling (`static/css/custom.css`)

Minimal custom CSS for:
- Clean, readable typography (16px base font)
- Proper spacing for headings
- Basic code block styling
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
- **Readable**: Focus on content with proper typography and spacing
- **Fast**: No custom fonts, minimal CSS, no animations
- **Clean**: Light theme with standard styling from the base Relearn theme
