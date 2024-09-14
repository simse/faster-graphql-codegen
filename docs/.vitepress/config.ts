import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "faster-graphql-codegen",
  description: "",
  base: "/",
  srcDir: "./src",
  cleanUrls: true,
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Reference', link: '/reference' }
    ],

    sidebar: [
      {
        text: 'Get Started',
        items: [
          { text: 'Quick Start', link: '/quick-start' },
          { text: 'Installation', link: '/install' },
        ]
      },
      {
        text: 'Configuration',
        items: [
          { text: 'Config File', link: '/config' },
        ]
      },
      {
        text: 'Plugins',
        items: [
          { text: 'Overview', link: '/plugins' },
        ]
      },
      {
        text: 'Reference',
        items: [
          { text: 'codegen.ts', link: '/reference/config' },
        ],
        collapsed: true,
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/simse/faster-graphql-codegen' }
    ],

    search: {
      provider: 'local'
    }
  }
})
