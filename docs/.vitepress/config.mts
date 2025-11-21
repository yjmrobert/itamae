import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Itamae Documentation",
  description: "A command-line tool for setting up Linux development workstations",
  base: '/',
  
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    logo: '/logo.svg',
    
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Documentation', link: '/installation' },
      { text: 'GitHub', link: 'https://github.com/yjmrobert/itamae' }
    ],

    sidebar: [
      {
        text: 'Getting Started',
        items: [
          { text: 'Installation', link: '/installation' },
          { text: 'Usage', link: '/usage' },
          { text: 'Troubleshooting', link: '/troubleshooting' },
          { text: 'Contributing', link: '/contributing' }
        ]
      },
      {
        text: 'Developer Guide',
        items: [
          { text: 'Overview', link: '/developers/' },
          { text: 'Adding Plugins', link: '/developers/adding-plugins' },
          { text: 'Testing', link: '/developers/testing' },
          { text: 'TUI Development', link: '/developers/tui' },
          { text: 'Releases', link: '/developers/releases' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/yjmrobert/itamae' }
    ],

    search: {
      provider: 'local'
    },

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024-present yjmrobert'
    }
  },

  markdown: {
    theme: 'monokai'
  }
})
