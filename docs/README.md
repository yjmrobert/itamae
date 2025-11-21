# Itamae Documentation

This directory contains the Hugo-based documentation for Itamae using the [Relearn theme](https://mcshelby.github.io/hugo-theme-relearn/).

## Local Development

### Prerequisites

Hugo Extended version 0.141.0 or later is required. Check your version:

```bash
hugo version
```

### Installing Hugo

#### Ubuntu/Debian

Download and install Hugo Extended:

```bash
# Download Hugo Extended 0.141.0
wget https://github.com/gohugoio/hugo/releases/download/v0.141.0/hugo_extended_0.141.0_linux-amd64.deb

# Install
sudo dpkg -i hugo_extended_0.141.0_linux-amd64.deb

# Verify
hugo version
```

#### Using Snap

```bash
sudo snap install hugo
```

#### Using Go

```bash
go install github.com/gohugoio/hugo@latest
```

### Build and Serve

#### Development Server

Run the Hugo development server with live reload:

```bash
cd docs
hugo server
```

Then open http://localhost:1313/ in your browser.

#### Production Build

Build the static site:

```bash
cd docs
hugo
```

Output will be in `docs/public/`.

## Theme

This documentation uses the [Hugo Relearn theme](https://github.com/McShelby/hugo-theme-relearn) with a custom **Tokyo Night Cyber Retro** color scheme. The theme is included as a git submodule.

### Custom Styling

The Tokyo Night cyber retro theme includes:
- **Color Palette**: Deep blue-purple backgrounds with neon cyan, purple, and green accents inspired by Tokyo Night
- **Typography**: Orbitron and Rajdhani fonts for a retro-futuristic feel, Space Mono for code
- **Visual Effects**: Subtle neon glows, scanline overlays, animated gradients, and cyber-style borders
- **Custom CSS**: Located in `static/css/theme-tokyo-night.css` and `static/css/custom.css`
- **Custom Layout**: `layouts/partials/custom-header.html` loads the fonts and theme CSS

### Update Theme

To update the base theme to the latest version:

```bash
cd docs/themes/hugo-theme-relearn
git pull origin main
```

## Documentation Structure

```
docs/
├── content/
│   ├── _index.md              # Home page
│   └── docs/
│       ├── installation.md     # Installation guide
│       ├── usage.md           # Usage guide
│       ├── troubleshooting.md # Troubleshooting
│       ├── contributing.md    # Contributing guide
│       └── developers/        # Developer documentation
│           ├── adding-plugins.md
│           ├── testing.md
│           ├── tui.md
│           └── releases.md
├── hugo.toml                  # Hugo configuration
└── themes/
    └── hugo-theme-relearn/    # Theme (submodule)
```

## Writing Documentation

### Frontmatter

Each page should have frontmatter:

```markdown
---
title: Page Title
weight: 10
---
```

### Shortcodes

The Relearn theme provides several useful shortcodes:

#### Notice Boxes

```markdown
{{% notice style="info" title="Title" %}}
Content here
{{% /notice %}}
```

Styles: `info`, `warning`, `note`, `tip`

#### Expand/Collapse

```markdown
{{% expand title="Click to expand" %}}
Hidden content
{{% /expand %}}
```

#### Buttons

```markdown
{{% button href="https://example.com" %}}Click Me{{% /button %}}
```

### Navigation

Pages are automatically organized by:
- Directory structure
- `weight` in frontmatter (lower numbers appear first)

## Deployment

Documentation is automatically deployed to GitHub Pages when changes are pushed to the `master` branch.

The deployment workflow is defined in `.github/workflows/deploy-docs.yml`.

## Troubleshooting

### Hugo Version Mismatch

If you see errors about Hugo version:

```
WARN Module "hugo-theme-relearn" is not compatible with this Hugo version
```

Upgrade Hugo to version 0.141.0 or later.

### Theme Not Found

If the theme directory is empty:

```bash
git submodule update --init --recursive
```

### Build Errors

Clear Hugo cache and rebuild:

```bash
hugo mod clean
hugo
```
