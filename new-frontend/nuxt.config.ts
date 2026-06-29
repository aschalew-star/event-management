export default defineNuxtConfig({
  runtimeConfig: {
    public: {
      hasuraEndpoint: process.env.NUXT_PUBLIC_HASURA_ENDPOINT || '/api/graphql',
      hasuraAdminSecret: process.env.NUXT_PUBLIC_HASURA_ADMIN_SECRET || '',
      apiBaseUrl: process.env.NUXT_PUBLIC_API_BASE_URL || 'http://localhost:8080',
      golangApiUrl: process.env.GOLANG_API_URL || 'http://golang-backend:8081'
    }
  },
  
  modules: [
    '@nuxtjs/tailwindcss',
    '@pinia/nuxt',
  ],
  
  css: [
    '~~/assets/css/main.css'  // Use @ instead of ~
  ],
   nitro: {
    routeRules: {
      '/api/graphql': {
        proxy: process.env.HASURA_GRAPHQL_ENDPOINT || 'http://localhost:8080/v1/graphql'
      }
    }
  },
  
  tailwindcss: {
    cssPath: '~~/assets/css/main.css',  // Use @ instead of ~
    configPath: 'tailwind.config.js',
    exposeConfig: false,
    viewer: true,
  },
  
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },
  
  compatibilityDate: '2025-02-26',
  pages: true,

})
