<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50">
    <!-- Header -->
    <header class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">My Events</h1>
            <p class="text-sm text-gray-500">Manage your created events</p>
          </div>
          <div class="flex items-center gap-4">
            <button
              @click="refetch"
              :disabled="loading"
              class="flex items-center gap-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors disabled:opacity-50"
            >
              <Icon :name="loading ? 'lucide:loader-2' : 'lucide:refresh-cw'" class="w-4 h-4" :class="loading ? 'animate-spin' : ''" />
              Refresh
            </button>
            <NuxtLink
              to="/events/create"
              class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg hover:shadow-lg transition-all hover:-translate-y-0.5"
            >
              <Icon name="lucide:plus" class="w-4 h-4" />
              Create Event
            </NuxtLink>
          </div>
        </div>
      </div>
    </header>

    <!-- Stats -->
    <div class="max-w-7xl mx-auto px-4 pt-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm text-gray-500">Total Events</span>
              <p class="text-2xl font-bold text-gray-900">{{ eventsCount }}</p>
            </div>
            <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
              <Icon name="lucide:calendar" class="w-5 h-5 text-blue-600" />
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm text-gray-500">Published</span>
              <p class="text-2xl font-bold text-green-600">{{ publishedCount }}</p>
            </div>
            <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
              <Icon name="lucide:check-circle" class="w-5 h-5 text-green-600" />
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm text-gray-500">Drafts</span>
              <p class="text-2xl font-bold text-yellow-600">{{ draftCount }}</p>
            </div>
            <div class="w-10 h-10 bg-yellow-100 rounded-lg flex items-center justify-center">
              <Icon name="lucide:file-text" class="w-5 h-5 text-yellow-600" />
            </div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm text-gray-500">Upcoming</span>
              <p class="text-2xl font-bold text-purple-600">{{ upcomingCount }}</p>
            </div>
            <div class="w-10 h-10 bg-purple-100 rounded-lg flex items-center justify-center">
              <Icon name="lucide:calendar-clock" class="w-5 h-5 text-purple-600" />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="max-w-7xl mx-auto px-4">
      <div class="flex gap-2 border-b border-gray-200 dark:border-gray-700 mb-6 overflow-x-auto">
        <button
          v-for="tab in tabs"
          :key="tab.value"
          @click="activeTab = tab.value"
          class="px-4 py-2 text-sm font-medium transition-colors whitespace-nowrap relative"
          :class="activeTab === tab.value 
            ? 'text-blue-600 border-b-2 border-blue-600' 
            : 'text-gray-500 hover:text-gray-700'"
        >
          {{ tab.label }}
          <span class="ml-2 px-2 py-0.5 text-xs rounded-full" 
            :class="activeTab === tab.value 
              ? 'bg-blue-100 text-blue-600' 
              : 'bg-gray-100 text-gray-600'"
          >
            {{ getTabCount(tab.value) }}
          </span>
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 pb-8">
      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-16">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>

      <!-- Error -->
      <div v-else-if="error" class="text-center py-16 bg-white rounded-2xl shadow-xl">
        <Icon name="lucide:alert-circle" class="w-20 h-20 text-red-400 mx-auto mb-4" />
        <h3 class="text-2xl font-semibold text-gray-700 mb-2">Failed to Load Events</h3>
        <p class="text-gray-500 mb-6">{{ error.message || 'Something went wrong' }}</p>
        <button
          @click="refetch"
          class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition-colors"
        >
          <Icon name="lucide:refresh-cw" class="w-5 h-5" />
          Try Again
        </button>
      </div>

      <!-- Events Grid -->
      <div v-else-if="filteredEvents.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="event in filteredEvents"
          :key="event.id"
          class="group bg-white rounded-2xl shadow-xl overflow-hidden hover:shadow-2xl transition-all duration-300 hover:-translate-y-1"
        >
          <!-- Image Section -->
          <div class="relative h-48 overflow-hidden">
            <!-- Get featured image from event_images table -->
            <img
              :src="getFeaturedImage(event) || '/images/placeholder-event.jpg'"
              :alt="event.title"
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
              @error="handleImageError"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
            
            <!-- Badges -->
            <div class="absolute top-3 left-3 flex gap-2 flex-wrap">
              <span
                v-if="event.is_free"
                class="px-2.5 py-1 bg-green-500/90 backdrop-blur-sm text-white rounded-full text-xs font-medium"
              >
                <Icon name="lucide:gift" class="w-3 h-3 inline mr-1" />
                Free
              </span>
              <span
                v-else
                class="px-2.5 py-1 bg-blue-500/90 backdrop-blur-sm text-white rounded-full text-xs font-medium"
              >
                <Icon name="lucide:dollar-sign" class="w-3 h-3 inline mr-1" />
                ${{ Number(event.price || 0).toFixed(2) }}
              </span>
              <span
                v-if="event.status"
                class="px-2.5 py-1 backdrop-blur-sm text-white rounded-full text-xs font-medium"
                :class="statusColors[event.status] || 'bg-gray-500/90'"
              >
                <Icon :name="event.status === 'published' ? 'lucide:check-circle' : 'lucide:file-text'" class="w-3 h-3 inline mr-1" />
                {{ formatStatus(event.status) }}
              </span>
            </div>

            <!-- Date and Time Badge -->
            <div class="absolute bottom-3 left-3 flex flex-col gap-1">
              <span class="text-white text-xs font-medium bg-black/50 backdrop-blur-sm px-2 py-1 rounded-full">
                <Icon name="lucide:calendar" class="w-3 h-3 inline mr-1" />
                {{ formatDate(event.event_date) }}
              </span>
              <span v-if="event.start_time || event.end_time" class="text-white text-xs font-medium bg-black/50 backdrop-blur-sm px-2 py-1 rounded-full">
                <Icon name="lucide:clock" class="w-3 h-3 inline mr-1" />
                {{ event.start_time || 'TBD' }} {{ event.end_time ? `- ${event.end_time}` : '' }}
              </span>
            </div>

            <!-- Action Buttons Overlay -->
            <div class="absolute top-3 right-3 flex gap-2">
              <NuxtLink
                :to="`/events/${event.id}`"
                class="p-2 bg-white/90 backdrop-blur-sm rounded-full hover:bg-blue-50 transition-colors shadow-lg"
                title="View Event"
              >
                <Icon name="lucide:eye" class="w-4 h-4 text-blue-600" />
              </NuxtLink>
              
              <button
                @click.prevent="editEvent(event.id)"
                class="p-2 bg-white/90 backdrop-blur-sm rounded-full hover:bg-purple-50 transition-colors shadow-lg"
                title="Edit Event"
              >
                <Icon name="lucide:edit-2" class="w-4 h-4 text-purple-600" />
              </button>
              
              <button
                @click.prevent="confirmDelete(event)"
                class="p-2 bg-white/90 backdrop-blur-sm rounded-full hover:bg-red-50 transition-colors shadow-lg"
                title="Delete Event"
              >
                <Icon name="lucide:trash-2" class="w-4 h-4 text-red-500" />
              </button>
            </div>
          </div>

          <!-- Content Section -->
          <div class="p-5">
            <!-- Category Badge -->
            <div v-if="event.category" class="mb-2">
              <span class="inline-flex items-center gap-1 px-2 py-0.5 bg-blue-50 text-blue-700 rounded-full text-xs">
                <Icon v-if="event.category.icon" :name="event.category.icon" class="w-3 h-3" />
                {{ event.category.name }}
              </span>
            </div>

            <h3 class="font-semibold text-gray-900 line-clamp-1 group-hover:text-blue-600 transition-colors">
              {{ event.title }}
            </h3>

            <div class="flex items-center gap-1 mt-1">
              <span class="text-xs text-gray-500">by</span>
              <span class="text-xs font-medium text-blue-600">
                {{ event.user?.name || 'You' }}
              </span>
            </div>

            <p class="text-sm text-gray-600 line-clamp-2 mt-2">
              {{ event.description }}
            </p>

            <!-- Event Details -->
            <div class="space-y-1.5 mt-3 text-sm text-gray-500">
              <div v-if="event.venue || event.address" class="flex items-center gap-2">
                <Icon name="lucide:map-pin" class="w-4 h-4 text-blue-600 flex-shrink-0" />
                <span class="line-clamp-1">{{ event.venue || event.address || 'Venue TBD' }}</span>
              </div>
              <div class="flex items-center gap-2">
                <Icon name="lucide:bookmark" class="w-4 h-4 text-blue-600 flex-shrink-0" />
                <span>{{ getBookmarkCount(event) }} bookmarks</span>
              </div>
              <div class="flex items-center gap-2">
                <Icon name="lucide:eye" class="w-4 h-4 text-blue-600 flex-shrink-0" />
                <span>{{ event.view_count || 0 }} views</span>
              </div>
            </div>

            <!-- Footer -->
            <div class="mt-3 pt-3 border-t border-gray-100 flex items-center justify-between">
              <span class="text-xs text-gray-400">
                <Icon name="lucide:clock" class="w-3 h-3 inline mr-1" />
                Created {{ formatRelativeTime(event.created_at) }}
              </span>
              <div class="flex items-center gap-3">
                <!-- Status Toggle Button -->
                <button
                  @click.prevent="toggleEventStatus(event)"
                  class="text-xs font-medium flex items-center gap-1 px-2 py-1 rounded-lg transition-colors"
                  :class="event.status === 'published' 
                    ? 'text-green-600 bg-green-50 hover:bg-green-100' 
                    : 'text-yellow-600 bg-yellow-50 hover:bg-yellow-100'"
                >
                  <Icon :name="event.status === 'published' ? 'lucide:check-circle' : 'lucide:file-text'" class="w-3 h-3" />
                  {{ formatStatus(event.status) }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-16 bg-white rounded-2xl shadow-xl">
        <Icon :name="activeTab === 'all' ? 'lucide:calendar-x' : 
                     activeTab === 'published' ? 'lucide:check-circle' : 
                     'lucide:file-text'" 
              class="w-20 h-20 text-gray-300 mx-auto mb-4" />
        <h3 class="text-2xl font-semibold text-gray-700 mb-2">No Events Found</h3>
        <p class="text-gray-500 mb-6 max-w-md mx-auto">
          {{ getEmptyStateMessage() }}
        </p>
        <NuxtLink
          to="/events/create"
          class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-xl hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <Icon name="lucide:plus" class="w-5 h-5" />
          Create Your First Event
        </NuxtLink>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center p-4" @click.self="closeDeleteModal">
        <div class="absolute inset-0 bg-black/50 backdrop-blur-sm"></div>
        
        <div class="relative bg-white rounded-2xl shadow-2xl max-w-md w-full p-6 animate-in fade-in zoom-in duration-200">
          <div class="text-center">
            <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <Icon name="lucide:alert-triangle" class="w-8 h-8 text-red-600" />
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-2">Delete Event</h3>
            <p class="text-gray-500 mb-2">
              Are you sure you want to delete this event?
            </p>
            <p class="text-gray-700 font-semibold mb-6">
              "{{ eventToDelete?.title || 'Untitled Event' }}"
            </p>
            <p class="text-sm text-red-500 mb-6">
              ⚠️ This will also delete all associated images, bookmarks, and tickets.
            </p>
            <div class="flex gap-3">
              <button
                @click="closeDeleteModal"
                class="flex-1 px-4 py-2.5 bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-xl transition-colors font-medium"
              >
                Cancel
              </button>
              <button
                @click="deleteEvent"
                :disabled="deleting"
                class="flex-1 px-4 py-2.5 bg-red-600 hover:bg-red-700 text-white rounded-xl transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
              >
                <Icon v-if="deleting" name="lucide:loader-2" class="w-4 h-4 animate-spin" />
                <Icon v-else name="lucide:trash-2" class="w-4 h-4" />
                {{ deleting ? 'Deleting...' : 'Delete Event' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onActivated } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useMutation } from '@vue/apollo-composable'
import { useAuthStore } from '~/stores/auth'
import { useToast } from 'vue-toastification'
import { GET_USER_EVENTS } from '~/graphql/eventQueries'
import { DELETE_EVENT, UPDATE_EVENT_STATUS } from '~/graphql/eventMutations'

definePageMeta({
  middleware: 'auth'
})

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const activeTab = ref('all')
const showDeleteModal = ref(false)
const eventToDelete = ref<any>(null)
const deleting = ref(false)

const tabs = [
  { label: 'All Events', value: 'all' },
  { label: 'Published', value: 'published' },
  { label: 'Drafts', value: 'draft' },
  { label: 'Upcoming', value: 'upcoming' },
]

const userId = computed(() => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    return authStore.user.id
  }
  return null
})

// Fetch user's events
const { 
  result: eventsResult, 
  loading, 
  error,
  refetch,
  onError 
} = useQuery(
  GET_USER_EVENTS,
  () => ({ 
    userId: userId.value || '00000000-0000-0000-0000-000000000000' 
  }),
  {
    fetchPolicy: 'network-only',
    skip: () => {
      return !process.client || !authStore.isAuthenticated || !userId.value
    },
    notifyOnNetworkStatusChange: true,
  }
)

onError((error) => {
  console.error('Error fetching events:', error)
  toast.error('Failed to load your events')
})

const events = computed(() => {
  return eventsResult.value?.events || []
})

// Stats
const eventsCount = computed(() => events.value.length)
const publishedCount = computed(() => events.value.filter(e => e.status === 'published').length)
const draftCount = computed(() => events.value.filter(e => e.status === 'draft' || !e.status).length)
const upcomingCount = computed(() => {
  const now = new Date()
  return events.value.filter(e => {
    const eventDate = new Date(e.event_date)
    return eventDate >= now && e.status !== 'cancelled'
  }).length
})

const getTabCount = (tab: string) => {
  switch (tab) {
    case 'all': return eventsCount.value
    case 'published': return publishedCount.value
    case 'draft': return draftCount.value
    case 'upcoming': return upcomingCount.value
    default: return 0
  }
}

// Filter events based on active tab
const filteredEvents = computed(() => {
  switch (activeTab.value) {
    case 'published':
      return events.value.filter(e => e.status === 'published')
    case 'draft':
      return events.value.filter(e => e.status === 'draft' || !e.status)
    case 'upcoming': {
      const now = new Date()
      return events.value.filter(e => {
        const eventDate = new Date(e.event_date)
        return eventDate >= now && e.status !== 'cancelled'
      })
    }
    default:
      return events.value
  }
})

// Helper functions
const getFeaturedImage = (event: any) => {
  if (event.event_images && event.event_images.length > 0) {
    const featured = event.event_images.find((img: any) => img.is_featured)
    return featured?.image_url || event.event_images[0]?.image_url
  }
  return null
}

const getBookmarkCount = (event: any) => {
  return event.bookmarks_aggregate?.aggregate?.count || 0
}

const formatStatus = (status: string) => {
  if (!status) return 'Draft'
  return status.charAt(0).toUpperCase() + status.slice(1)
}

const getEmptyStateMessage = () => {
  switch (activeTab.value) {
    case 'all': return 'You haven\'t created any events yet.'
    case 'published': return 'You don\'t have any published events.'
    case 'upcoming': return 'You don\'t have any upcoming events.'
    case 'draft': return 'You don\'t have any draft events.'
    default: return 'No events found.'
  }
}

// Mutations
const { mutate: deleteEventMutation } = useMutation(DELETE_EVENT)
const { mutate: updateEventStatusMutation } = useMutation(UPDATE_EVENT_STATUS)

const editEvent = (eventId: string) => {
  router.push(`/events/edit/${eventId}`)
}

const confirmDelete = (event: any) => {
  eventToDelete.value = event
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  eventToDelete.value = null
  deleting.value = false
}

const deleteEvent = async () => {
  if (!eventToDelete.value) {
    toast.error('No event selected for deletion')
    return
  }

  deleting.value = true
  
  try {
    console.log('Deleting event:', eventToDelete.value.id)
    
    const result = await deleteEventMutation({
      eventId: eventToDelete.value.id
    })
    
    console.log('Delete result:', result)
    
    if (result?.data?.delete_events_by_pk?.id) {
      await refetch()
      toast.success(`Event "${eventToDelete.value.title}" deleted successfully`)
      closeDeleteModal()
    } else {
      toast.error('Failed to delete event. Please try again.')
    }
  } catch (error: any) {
    console.error('Error deleting event:', error)
    toast.error(error.message || 'Failed to delete event')
  } finally {
    deleting.value = false
  }
}

const toggleEventStatus = async (event: any) => {
  const newStatus = event.status === 'published' ? 'draft' : 'published'
  
  try {
    const result = await updateEventStatusMutation({
      eventId: event.id,
      status: newStatus
    })
    
    if (result?.data?.update_events?.affected_rows > 0) {
      await refetch()
      toast.success(`Event ${newStatus === 'published' ? 'published' : 'moved to draft'}`)
    } else {
      toast.error('Failed to update event status')
    }
  } catch (error: any) {
    console.error('Error updating event status:', error)
    toast.error(error.message || 'Failed to update event status')
  }
}

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/images/placeholder-event.jpg'
}

const statusColors: Record<string, string> = {
  published: 'bg-green-500/90',
  draft: 'bg-yellow-500/90',
  cancelled: 'bg-red-500/90',
}

const formatDate = (date: string) => {
  if (!date) return 'TBD'
  try {
    return new Date(date).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    })
  } catch {
    return 'Invalid date'
  }
}

const formatRelativeTime = (date: string) => {
  if (!date) return ''
  try {
    const now = new Date()
    const past = new Date(date)
    const diffMinutes = Math.floor((now.getTime() - past.getTime()) / 1000 / 60)
    const diffHours = Math.floor(diffMinutes / 60)
    const diffDays = Math.floor(diffHours / 24)
    
    if (diffMinutes < 1) return 'Just now'
    if (diffMinutes < 60) return `${diffMinutes}m ago`
    if (diffHours < 24) return `${diffHours}h ago`
    if (diffDays < 7) return `${diffDays}d ago`
    return formatDate(date)
  } catch {
    return ''
  }
}

// Refetch on mount
onMounted(async () => {
  if (authStore.isAuthenticated && userId.value) {
    await refetch()
  }
})

// Refetch when component is activated
onActivated(async () => {
  if (authStore.isAuthenticated && userId.value) {
    await refetch()
  }
})

// Watch for authentication changes
watch(
  () => authStore.isAuthenticated,
  async (authenticated) => {
    if (authenticated && userId.value) {
      await refetch()
    }
  }
)
</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes zoomIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-in {
  animation: fadeIn 0.2s ease-out;
}

.fade-in {
  animation: fadeIn 0.2s ease-out;
}

.zoom-in {
  animation: zoomIn 0.2s ease-out;
}
</style>