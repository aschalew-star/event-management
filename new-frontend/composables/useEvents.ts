// composables/useEvents.ts
import { ref } from 'vue'
import { useEventStore } from '~/stores/events'

export const useEvents = () => {
  const store = useEventStore()
  const loading = ref(false)
  const error = ref<string | null>(null)

  const createEvent = async (formData: FormData) => {
    loading.value = true
    error.value = null
    
    try {
      // Get auth token
      const authStore = useAuthStore()
      const token = authStore.token

      if (!token) {
        throw new Error('Authentication required')
      }

      const response: any = await $fetch('http://localhost:4000/api/events/create', {
        method: 'POST',
        body: formData,
        headers: {
          Authorization: `Bearer ${token}`
        }
      })

      if (!response?.success) {
        throw new Error(response?.message || 'Failed to create event')
      }

      // Refresh the store
      await store.fetchEvents()
      
      return response
    } catch (err: any) {
      error.value = err.message || 'An unexpected error occurred'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    error,
    createEvent,
    // Pass through store methods
    fetchEvents: store.fetchEvents,
    fetchEvent: store.fetchEvent,
    deleteEvent: store.deleteEvent,
    events: store.events,
    currentEvent: store.currentEvent
  }
}