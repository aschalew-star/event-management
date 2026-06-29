import { defineNuxtPlugin } from '#app'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  
  // Set up GraphQL client configuration
  const graphqlEndpoint = config.public.hasuraEndpoint || '/api/graphql'
  
  // Add global headers for GraphQL requests
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    'x-hasura-admin-secret': config.public.hasuraAdminSecret || ''
  }
  
  // Helper function for GraphQL queries
  const graphqlQuery = async (query: string, variables?: any, authToken?: string) => {
    const requestHeaders = { ...headers }
    
    if (authToken) {
      requestHeaders['Authorization'] = `Bearer ${authToken}`
    }
    
    try {
      const response = await $fetch(graphqlEndpoint, {
        method: 'POST',
        headers: requestHeaders,
        body: {
          query,
          variables
        }
      })
      
      if (response?.errors) {
        throw new Error(response.errors[0]?.message || 'GraphQL error')
      }
      
      return response
    } catch (error) {
      console.error('GraphQL request error:', error)
      throw error
    }
  }
  
  return {
    provide: {
      graphql: graphqlQuery
    }
  }
})