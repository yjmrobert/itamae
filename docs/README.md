# Itamae Documentation

This directory contains the VitePress-based documentation for Itamae.

## Local Development

### Prerequisites

- Node.js 20 or higher
- npm (comes with Node.js)

### Quick Start

```bash
# Install dependencies
npm install

# Start dev server with hot reload
npm run docs:dev

# Visit http://localhost:5173
```

### Build for Production

```bash
# Build static site
npm run docs:build

# Preview production build
npm run docs:preview
```

## Structure

```
docs/
├── .vitepress/
│   ├── config.mts          # VitePress configuration
│   └── theme/
│       ├── index.ts        # Theme entry
│       └── style.css       # Tokyo Night custom theme
├── index.md                # Home page (hero layout)
├── installation.md         # Installation guide
├── usage.md                # Usage guide
├── troubleshooting.md      # Troubleshooting
├── contributing.md         # Contributing guide
└── developers/             # Developer documentation
    ├── index.md            # Overview
    ├── adding-plugins.md   # Plugin development
    ├── testing.md          # Testing guide
    ├── tui.md              # TUI development
    └── releases.md         # Release process
```

## Theme

The documentation uses a custom **Tokyo Night** dark theme built on VitePress's default theme. The color palette is defined in `.vitepress/theme/style.css`.

## Writing Documentation

### Containers (Callouts)

Use VitePress custom containers for notices:

```markdown
:::info
Information message
:::

:::tip
Helpful tip
:::

:::warning
Warning message
:::

:::danger
Critical warning
:::
```

### Details/Collapsible Sections

Use HTML `<details>` tags:

```markdown
<details>
<summary>Click to expand</summary>

Content here...

</details>
```

### Links

Use relative markdown links:

```markdown
[Link to page](/path/to/page)
[Link to section](#section-anchor)
```

## Deployment

Documentation is automatically deployed to GitHub Pages when changes are pushed to the `master` branch.

The deployment workflow is defined in `.github/workflows/deploy-vitepress.yml`.

## Migration Notes

This documentation was migrated from Hugo (hugo-theme-relearn) to VitePress on November 21, 2024.

### Why VitePress?

- **Simpler installation**: npm-based, no Hugo Extended required
- **Better DX**: Fast hot reload, Vue-based
- **Portable syntax**: Standard Markdown, minimal framework-specific syntax
- **No submodules**: Theme is npm package, not git submodule
- **Modern**: Active development, growing ecosystem

## Troubleshooting

### Module not found errors

```bash
rm -rf node_modules package-lock.json
npm install
```

### Port already in use

VitePress runs on port 5173 by default. Change it:

```bash
npm run docs:dev -- --port 3000
```

### Build fails

Check Node.js version:

```bash
node --version  # Should be 20+
```

## Resources

- [VitePress Documentation](https://vitepress.dev/)
- [VitePress GitHub](https://github.com/vuejs/vitepress)
- [Markdown Extensions](https://vitepress.dev/guide/markdown)
- [Default Theme Config](https://vitepress.dev/reference/default-theme-config)
