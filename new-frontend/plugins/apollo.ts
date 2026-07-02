import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core'
import { setContext } from '@apollo/client/link/context'
import { provideApolloClient } from '@vue/apollo-composable'
import { useAuthStore } from '~/stores/auth'

export default defineNuxtPlugin(() => {
  const isServer = import.meta.env.SSR
  
  const uri = isServer 
    ? (process.env.HASURA_GRAPHQL_ENDPOINT || 'http://localhost:8080/v1/graphql')
    : '/api/graphql'

  const httpLink = createHttpLink({ uri })

  const authLink = setContext((_, { headers }) => {
    if (isServer) {
      return {
        headers: {
          ...headers,
          'x-hasura-admin-secret': "myadminsecretkey"
        }
      }
    }

    const authStore = useAuthStore()
    
    // 1. Force state hydration if it hasn't happened yet
    if (!authStore.token && window.localStorage.getItem('auth_token')) {
      authStore.loadAuth()
    }

    // 2. Safely grab the clean token string
    const token = authStore.token

    // 3. Strict structural checks to guarantee we don't pass an invalid string
    const hasValidToken = token && 
                          typeof token === 'string' && 
                          token.trim() !== '' && 
                          token !== 'null' && 
                          token !== 'undefined'

    return {
      headers: {
        ...headers,
        // Only attach Authorization if it is a structurally valid JWT
        ...(hasValidToken 
          ? { Authorization: `Bearer ${token.trim()}` } 
          : { 'x-hasura-role': 'anonymous' }
        )
      }
    }
  })

  const apolloClient = new ApolloClient({
    link: authLink.concat(httpLink),
    cache: new InMemoryCache(),
    ssrMode: isServer
  })

  provideApolloClient(apolloClient)
})