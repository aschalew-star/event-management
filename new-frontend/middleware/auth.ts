// middleware/auth.ts
import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()

  // 1. Force state loading immediately if we are executing on the browser
  if (process.client && !authStore.isAuthenticated) {
    authStore.loadAuth()
  }

  // 2. Fallback check: Read token directly from storage to eliminate 
  // any microsecond timing delays during Pinia's initial mounting loop.
  let hasToken = authStore.isAuthenticated
  if (process.client && !hasToken) {
    hasToken = !!localStorage.getItem('auth_token')
  }

  // 3. SSR Safety Gate: If this is running on the server, skip redirect commands.
  // The server lacks local storage access, so it must defer checking to the client.
  if (!process.client) {
    return
  }

  // 4. Handle route access restrictions on the client side
  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !hasToken) {
    return navigateTo('/login')
  }

  // Prevent logged-in users from hitting login or registration gates
  const authPages = ['/login', '/register']
  const targetPath = to.path.replace(/\/$/, '')

  if (authPages.includes(targetPath) && hasToken) {
    return navigateTo('/events')
  }
})