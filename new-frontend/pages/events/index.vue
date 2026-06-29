<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold dark:text-white">Events</h1>
      <NuxtLink 
        to="/events/create" 
        class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700"
      >
        Create Event
      </NuxtLink>
    </div>

    <!-- Filters -->
    <div class="mb-6 flex flex-wrap gap-4">
      <input
        type="text"
        v-model="filters.title"
        placeholder="Search events..."
        class="px-3 py-2 border rounded-md dark:bg-gray-800 dark:border-gray-600 dark:text-white"
        @input="applyFilters"
      />
      
      <select
        v-model="filters.category_id"
        class="px-3 py-2 border rounded-md dark:bg-gray-800 dark:border-gray-600 dark:text-white"
        @change="applyFilters"
      >
        <option value="">All Categories</option>
        <option v-for="category in categories" :key="category.id" :value="category.id">
          {{ category.name }}
        </option>
      </select>
      
      <select
        v-model="filters.status"
        class="px-3 py-2 border rounded-md dark:bg-gray-800 dark:border-gray-600 dark:text-white"
        @change="applyFilters"
      >
        <option value="">All Status</option>
        <option value="upcoming">Upcoming</option>
        <option value="ongoing">Ongoing</option>
        <option value="completed">Completed</option>
      </select>
      
      <button
        @click="resetFilters"
        class="px-4 py-2 bg-gray-500 text-white rounded-md hover:bg-gray-600"
      >
        Reset
      </button>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
    </div>

    <!-- Events Grid -->
    <div v-else-if="events.length" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div 
        v-for="event in events" 
        :key="event.id"
        class="bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden hover:shadow-lg transition"
      >
        <div class="p-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
            {{ event.title }}
          </h3>
          <p class="text-gray-600 dark:text-gray-400 text-sm mb-2">
            {{ event.description?.slice(0, 100) }}...
          </p>
          <div class="flex justify-between items-center text-sm text-gray-500 dark:text-gray-400">
            <span>{{ formatDate(event.event_date) }}</span>
            <span>{{ event.is_free ? 'Free' : `$${event.price}` }}</span>
          </div>
          <div class="mt-3 flex space-x-2">
            <NuxtLink 
              :to="`/events/${event.id}`" 
              class="px-3 py-1 text-sm bg-primary-600 text-white rounded hover:bg-primary-700"
            >
              View
            </NuxtLink>
            <NuxtLink 
              :to="`/events/edit/${event.id}`" 
              class="px-3 py-1 text-sm bg-blue-500 text-white rounded hover:bg-blue-600"
            >
              Edit
            </NuxtLink>
            <button 
              @click="deleteEvent(event.id)" 
              class="px-3 py-1 text-sm bg-red-500 text-white rounded hover:bg-red-600"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- No Events -->
    <div v-else class="text-center py-12">
      <p class="text-gray-600 dark:text-gray-400">No events found</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useEventStore } from '~/stores/events'

// Middleware
definePageMeta({
  requiresAuth: true
})

const eventStore = useEventStore()
const toast = useToast()
const router = useRouter()

const events = computed(() => eventStore.events)
const loading = ref(false)

// Filters
const filters = reactive({
  title: '',
  category_id: '',
  status: ''
})

// Categories (you'd fetch these from Hasura)
const categories = ref([])

// Load events on mount
onMounted(async () => {
  await loadEvents()
  await loadCategories()
})

// Load events with filters
const loadEvents = async () => {
  loading.value = true
  try {
    const where: any = {}
    
    if (filters.title) {
      where.title = { _ilike: `%${filters.title}%` }
    }
    
    if (filters.category_id) {
      where.category_id = { _eq: filters.category_id }
    }
    
    if (filters.status) {
      where.status = { _eq: filters.status }
    }
    
    await eventStore.fetchEvents(where)
  } catch (error) {
    toast.error('Failed to load events')
  } finally {
    loading.value = false
  }
}

// Load categories
const loadCategories = async () => {
  try {
    const response = await $fetch('/api/graphql', {
      method: 'POST',
      body: {
        query: `
          query GetCategories {
            categories(order_by: {name: asc}) {
              id
              name
            }
          }
        `
      }
    })
    categories.value = response?.data?.categories || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

// Apply filters
const applyFilters = () => {
  loadEvents()
}

// Reset filters
const resetFilters = () => {
  filters.title = ''
  filters.category_id = ''
  filters.status = ''
  loadEvents()
}

// Delete event
const deleteEvent = async (id: string) => {
  if (!confirm('Are you sure you want to delete this event?')) return
  
  const result = await eventStore.deleteEvent(id)
  if (result.success) {
    toast.success('Event deleted successfully')
    await loadEvents()
  } else {
    toast.error(result.error || 'Failed to delete event')
  }
}

// Format date helper
const formatDate = (dateString: string) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}
</script>