import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core'
import { provideApolloClient } from '@vue/apollo-composable'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  
  const httpLink = createHttpLink({
    uri: config.public.hasuraEndpoint,
    headers: {
      'x-hasura-admin-secret': 'myadminsecretkey'
    }
  })

  const apolloClient = new ApolloClient({
    link: httpLink,
    cache: new InMemoryCache()
  })

  provideApolloClient(apolloClient)
})