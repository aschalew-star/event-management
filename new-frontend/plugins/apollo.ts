import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core'
import { provideApolloClient } from '@vue/apollo-composable'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  
  // 1. Determine the correct URL based on environment
  // On Node server: Hit Hasura directly using your Nitro proxy target
  // In Browser: Use the relative proxy path
  const isServer = import.meta.env.SSR
  
  const uri = isServer 
    ? (process.env.HASURA_GRAPHQL_ENDPOINT || 'http://localhost:8080/v1/graphql')
    : '/api/graphql'

  const httpLink = createHttpLink({
    uri: uri,
    headers: isServer ? {
      // 2. Only attach the admin secret on the server-side to keep it secure!
      'x-hasura-admin-secret': "myadminsecretkey"
    } : {}
  })

  const apolloClient = new ApolloClient({
    link: httpLink,
    cache: new InMemoryCache(),
    ssrMode: isServer // Tells Apollo to optimize for SSR when on the server
  })

  provideApolloClient(apolloClient)
})