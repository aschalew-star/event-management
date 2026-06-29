import { defineStore } from 'pinia'

interface Event {
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
}

export const useEventStore = defineStore('events', {
  state: (): EventState => ({
    events: [],
    currentEvent: null,
    loading: false
  }),

  getters: {
    getUpcomingEvents: (state) => {
      const now = new Date().toISOString()
      return state.events.filter(event => event.event_date && event.event_date >= now)
    },
    getPastEvents: (state) => {
      const now = new Date().toISOString()
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
      try {
        const response = await $fetch('/api/graphql', {
          method: 'POST',
          body: {
            query: `
              query GetEvents($where: events_bool_exp, $order: [events_order_by!]) {
                events(where: $where, order_by: $order) {
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
                  images
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
              where: filters || {},
              order: [{ event_date: 'asc' }]
            }
          },
          headers: {
            'Authorization': `Bearer ${this.authToken}`
          }
        })

        this.events = response?.data?.events || []
        return { success: true, data: this.events }
      } catch (error: any) {
        console.error('Fetch events error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async fetchEvent(id: string) {
      this.loading = true
      try {
        const response = await $fetch('/api/graphql', {
          method: 'POST',
          body: {
            query: `
              query GetEvent($id: ID!) {
                event(id: $id) {
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
                  images
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
            'Authorization': `Bearer ${this.authToken}`
          }
        })

        this.currentEvent = response?.data?.event || null
        return { success: true, data: this.currentEvent }
      } catch (error: any) {
        console.error('Fetch event error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async createEvent(eventData: any) {
      this.loading = true
      try {
        const response = await $fetch('/api/graphql', {
          method: 'POST',
          body: {
            query: `
              mutation CreateEvent($input: CreateEventInput!) {
                createEvent(input: $input) {
                  id
                  message
                }
              }
            `,
            variables: {
              input: eventData
            }
          },
          headers: {
            'Authorization': `Bearer ${this.authToken}`
          }
        })

        const data = response?.data?.createEvent
        if (data?.id) {
          return { success: true, id: data.id, message: data.message }
        }
        return { success: false, error: 'Failed to create event' }
      } catch (error: any) {
        console.error('Create event error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async updateEvent(id: string, eventData: any) {
      this.loading = true
      try {
        const response = await $fetch('/api/graphql', {
          method: 'POST',
          body: {
            query: `
              mutation UpdateEvent($id: ID!, $input: UpdateEventInput!) {
                updateEvent(id: $id, input: $input) {
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
                  images
                  updated_at
                  category {
                    id
                    name
                  }
                }
              }
            `,
            variables: {
              id,
              input: eventData
            }
          },
          headers: {
            'Authorization': `Bearer ${this.authToken}`
          }
        })

        const data = response?.data?.updateEvent
        if (data?.id) {
          return { success: true, data }
        }
        return { success: false, error: 'Failed to update event' }
      } catch (error: any) {
        console.error('Update event error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async deleteEvent(id: string) {
      this.loading = true
      try {
        const response = await $fetch('/api/graphql', {
          method: 'POST',
          body: {
            query: `
              mutation DeleteEvent($id: ID!) {
                deleteEvent(id: $id) {
                  id
                  message
                }
              }
            `,
            variables: { id }
          },
          headers: {
            'Authorization': `Bearer ${this.authToken}`
          }
        })

        const data = response?.data?.deleteEvent
        if (data?.id) {
          this.events = this.events.filter(event => event.id !== id)
          return { success: true, message: data.message }
        }
        return { success: false, error: 'Failed to delete event' }
      } catch (error: any) {
        console.error('Delete event error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async uploadEventImages(eventId: string, images: string[]) {
      this.loading = true
      try {
        const response = await $fetch('/api/graphql', {
          method: 'POST',
          body: {
            query: `
              mutation UploadEventImages($id: ID!, $images: [String!]!) {
                uploadEventImages(id: $id, images: $images) {
                  success
                  message
                  urls
                }
              }
            `,
            variables: {
              id: eventId,
              images
            }
          },
          headers: {
            'Authorization': `Bearer ${this.authToken}`
          }
        })

        const data = response?.data?.uploadEventImages
        return { 
          success: data?.success || false, 
          message: data?.message,
          urls: data?.urls || []
        }
      } catch (error: any) {
        console.error('Upload images error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    }
  }
})