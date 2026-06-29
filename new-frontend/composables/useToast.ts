// composables/useToast.ts
import { useNuxtApp } from 'nuxt/app'

export const useToast = () => {
  const { $toast } = useNuxtApp()
  const isServer = typeof window === 'undefined'
  
  // Return a mock/no-op function during SSR
  if (isServer) {
    return {
      success: (msg: string) => console.log('[Toast SSR]', msg),
      error: (msg: string) => console.log('[Toast SSR]', msg),
      info: (msg: string) => console.log('[Toast SSR]', msg),
      warning: (msg: string) => console.log('[Toast SSR]', msg),
      clear: () => {},
      // Add other methods as no-ops
    } as any
  }
  
  if (!$toast) {
    throw new Error('Toast plugin not initialized. Make sure the plugin is properly installed.')
  }
  
  return $toast
}