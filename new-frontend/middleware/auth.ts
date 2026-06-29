export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore()
  
  // Load auth state from localStorage if not loaded
  if (!authStore.isAuthenticated) {
    authStore.loadAuth()
  }
  
  // Check if route requires authentication
  const requiresAuth = to.meta.requiresAuth !== false
  
  if (requiresAuth && !authStore.isAuthenticated) {
    // Redirect to login page
    return navigateTo('/login')
  }
  
  // Redirect to dashboard if already authenticated and trying to access auth pages
  const authPages = ['/login', '/register']
  if (authPages.includes(to.path) && authStore.isAuthenticated) {
    return navigateTo('/dashboard')
  }
})