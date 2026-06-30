// stores/event.ts
import { defineStore } from 'pinia'
import { useAuthStore } from './auth'

export interface Event {
  id: string
  title: string
  description?: string
  category_id?: string
  price?: number
  is_free?: boolean
  venue?: string
  address?: string
  latitude?: number
  longitude?: number
  event_date?: string
  start_time?: string
  end_time?: string
  status?: string
  images?: string[]
  created_at?: string
  updated_at?: string
  category?: {
    id: string
    name: string
  }
}

interface EventState {
  events: Event[]
  currentEvent: Event | null
  loading: boolean
  error: string | null
}

export const useEventStore = defineStore('events', {
  state: (): EventState => ({
    events: [],
    currentEvent: null,
    loading: false,
    error: null
  }),

  getters: {
    getUpcomingEvents: (state) => {
      const now = new Date().toISOString().split('T')[0]
      return state.events.filter(event => event.event_date && event.event_date >= now)
    },
    getPastEvents: (state) => {
      const now = new Date().toISOString().split('T')[0]
      return state.events.filter(event => event.event_date && event.event_date < now)
    },
    getEventsByCategory: (state) => (categoryId: string) => {
      return state.events.filter(event => event.category_id === categoryId)
    },
    getFreeEvents: (state) => {
      return state.events.filter(event => event.is_free === true)
    }
  },

  actions: {
    get authToken(): string {
      const authStore = useAuthStore()
      return authStore.token || ''
    },

    async fetchEvents(filters?: any) {
      this.loading = true
      this.error = null
      try {
        const response: any = await $fetch('http://localhost:8080/v1/graphql', {
          method: 'POST',
          body: {
            query: `
              query GetEvents($where: events_bool_exp) {
                events(where: $where, order_by: [{ event_date: asc }]) {
                  id
                  title
                  description
                  category_id
                  price
                  is_free
                  venue
                  address
                  latitude
                  longitude
                  event_date
                  start_time
                  end_time
                  status
                  created_at
                  updated_at
                  category {
                    id
                    name
                  }
                }
              }
            `,
            variables: {
              where: filters || {}
            }
          },
          headers: {
            Authorization: `Bearer ${this.authToken}`
          }
        })

        this.events = response?.data?.events || []
        return { success: true, data: this.events }
      } catch (error: any) {
        this.error = error.message
        console.error('Fetch events error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async fetchEvent(id: string) {
      this.loading = true
      this.error = null
      try {
        const response: any = await $fetch('http://localhost:8080/v1/graphql', {
          method: 'POST',
          body: {
            query: `
              query GetEventByPk($id: uuid!) {
                events_by_pk(id: $id) {
                  id
                  title
                  description
                  category_id
                  price
                  is_free
                  venue
                  address
                  latitude
                  longitude
                  event_date
                  start_time
                  end_time
                  status
                  created_at
                  updated_at
                  category {
                    id
                    name
                  }
                }
              }
            `,
            variables: { id }
          },
          headers: {
            Authorization: `Bearer ${this.authToken}`
          }
        })

        this.currentEvent = response?.data?.events_by_pk || null
        return { success: true, data: this.currentEvent }
      } catch (error: any) {
        this.error = error.message
        console.error('Fetch event error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async createEvent(eventData: any) {
      this.loading = true
      this.error = null
      try {
        // Use your Go backend for event creation with image upload
        const formData = new FormData()
        Object.keys(eventData).forEach(key => {
          if (key === 'images' && Array.isArray(eventData[key])) {
            eventData[key].forEach((file: File) => {
              formData.append('images', file)
            })
          } else if (eventData[key] !== null && eventData[key] !== undefined) {
            formData.append(key, eventData[key].toString())
          }
        })

        const response: any = await $fetch('http://localhost:4000/api/events/create', {
          method: 'POST',
          body: formData,
          headers: {
            Authorization: `Bearer ${this.authToken}`
          }
        })

        if (response?.success) {
          // Refresh events list
          await this.fetchEvents()
          return { success: true, id: response.id, message: response.message }
        }
        return { success: false, error: response?.message || 'Failed to create event' }
      } catch (error: any) {
        this.error = error.message
        console.error('Create event error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async deleteEvent(id: string) {
      this.loading = true
      this.error = null
      try {
        const response: any = await $fetch('http://localhost:8080/v1/graphql', {
          method: 'POST',
          body: {
            query: `
              mutation DeleteEventByPk($id: uuid!) {
                delete_events_by_pk(id: $id) {
                  id
                }
              }
            `,
            variables: { id }
          },
          headers: {
            Authorization: `Bearer ${this.authToken}`
          }
        })

        if (response?.data?.delete_events_by_pk?.id) {
          this.events = this.events.filter(event => event.id !== id)
          return { success: true, message: 'Event deleted successfully' }
        }
        return { success: false, error: 'Failed to delete event' }
      } catch (error: any) {
        this.error = error.message
        console.error('Delete event error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    }
  }
})