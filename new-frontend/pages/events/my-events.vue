<!-- pages/my-events.vue -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-50 via-white to-pink-50">
    <!-- Header -->
    <header class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">My Events</h1>
            <p class="text-sm text-gray-500">Manage your events</p>
          </div>
          <div class="flex items-center gap-4">
            <button
              @click="showCreateModal = true"
              class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg hover:shadow-lg transition-all hover:-translate-y-0.5"
            >
              <Icon name="lucide:plus" class="w-5 h-5" />
              Create Event
            </button>
          </div>
        </div>
      </div>
    </header>

    <!-- Stats -->
    <div class="max-w-7xl mx-auto px-4 pt-6">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">Total Events</span>
            <Icon name="lucide:calendar" class="w-5 h-5 text-purple-600" />
          </div>
          <p class="text-2xl font-bold text-gray-900 mt-2">{{ myEvents.length }}</p>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">Published</span>
            <Icon name="lucide:check-circle" class="w-5 h-5 text-green-600" />
          </div>
          <p class="text-2xl font-bold text-gray-900 mt-2">
            {{ myEvents.filter(e => e.status === 'published').length }}
          </p>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">Drafts</span>
            <Icon name="lucide:file-text" class="w-5 h-5 text-yellow-600" />
          </div>
          <p class="text-2xl font-bold text-gray-900 mt-2">
            {{ myEvents.filter(e => e.status === 'draft').length }}
          </p>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">Total Views</span>
            <Icon name="lucide:eye" class="w-5 h-5 text-blue-600" />
          </div>
          <p class="text-2xl font-bold text-gray-900 mt-2">
            {{ myEvents.reduce((sum, e) => sum + (e.view_count || 0), 0) }}
          </p>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="max-w-7xl mx-auto px-4">
      <div class="bg-white rounded-xl shadow-sm p-4 mb-6">
        <div class="flex flex-wrap gap-3 items-center">
          <!-- Search -->
          <div class="flex-1 min-w-[200px]">
            <div class="relative">
              <Icon name="lucide:search" class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search your events..."
                class="w-full pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-all"
                @input="debouncedSearch"
              />
            </div>
          </div>

          <!-- Status Filter -->
          <select
            v-model="filters.status"
            class="px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 transition-all"
            @change="applyFilters"
          >
            <option value="">All Status</option>
            <option value="published">Published</option>
            <option value="draft">Draft</option>
            <option value="cancelled">Cancelled</option>
          </select>

          <!-- Type Filter -->
          <select
            v-model="filters.type"
            class="px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 transition-all"
            @change="applyFilters"
          >
            <option value="">All Types</option>
            <option value="free">Free</option>
            <option value="paid">Paid</option>
          </select>

          <!-- Sort -->
          <select
            v-model="filters.sortBy"
            class="px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 transition-all"
            @change="applyFilters"
          >
            <option value="event_date">Date</option>
            <option value="title">Title</option>
            <option value="view_count">Views</option>
            <option value="created_at">Created</option>
          </select>

          <button
            @click="resetFilters"
            class="px-3 py-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
            title="Reset filters"
          >
            <Icon name="lucide:rotate-ccw" class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 pb-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-16">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
      </div>

      <!-- Events Grid -->
      <div v-else-if="filteredEvents.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="event in paginatedEvents"
          :key="event.id"
          class="group bg-white rounded-2xl shadow-xl overflow-hidden hover:shadow-2xl transition-all duration-300 hover:-translate-y-1"
        >
          <!-- Image -->
          <div class="relative h-52 overflow-hidden">
            <img
              :src="event.featured_image || '/images/placeholder-event.jpg'"
              :alt="event.title"
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
            
            <!-- Badges -->
            <div class="absolute top-3 left-3 flex gap-2">
              <span
                :class="[
                  'px-2.5 py-1 rounded-full text-xs font-medium backdrop-blur-sm',
                  getStatusColor(event.status)
                ]"
              >
                {{ event.status.charAt(0).toUpperCase() + event.status.slice(1) }}
              </span>
              <span
                v-if="event.is_free"
                class="px-2.5 py-1 bg-green-500/90 backdrop-blur-sm text-white rounded-full text-xs font-medium"
              >
                Free
              </span>
            </div>

            <!-- Action Buttons -->
            <div class="absolute top-3 right-3 flex gap-2">
              <button
                @click="openEditModal(event)"
                class="p-2 bg-white/90 backdrop-blur-sm rounded-lg hover:bg-white shadow-lg transition-all hover:scale-110"
                title="Edit Event"
              >
                <Icon name="lucide:edit-2" class="w-4 h-4 text-gray-700" />
              </button>
              <button
                @click="confirmDelete(event)"
                class="p-2 bg-white/90 backdrop-blur-sm rounded-lg hover:bg-red-50 shadow-lg transition-all hover:scale-110"
                title="Delete Event"
              >
                <Icon name="lucide:trash-2" class="w-4 h-4 text-red-600" />
              </button>
            </div>

            <!-- Price -->
            <div class="absolute bottom-3 right-3">
              <span class="px-3 py-1.5 bg-white/90 backdrop-blur-sm rounded-lg text-sm font-bold text-purple-600 shadow-lg">
                {{ event.is_free ? 'Free' : `$${event.price}` }}
              </span>
            </div>
          </div>

          <!-- Content -->
          <div class="p-5">
            <div class="flex items-start justify-between mb-2">
              <h3 class="font-semibold text-gray-900 line-clamp-1 flex-1 text-lg">
                {{ event.title }}
              </h3>
            </div>

            <p class="text-sm text-gray-600 line-clamp-2 mb-3">
              {{ event.description }}
            </p>

            <div class="space-y-1.5 text-sm text-gray-500">
              <div class="flex items-center gap-2">
                <Icon name="lucide:calendar" class="w-4 h-4 text-purple-600 flex-shrink-0" />
                <span>{{ formatDate(event.event_date) }}</span>
                <span v-if="event.start_time" class="text-xs text-gray-400">
                  {{ event.start_time }} - {{ event.end_time }}
                </span>
              </div>
              <div class="flex items-center gap-2">
                <Icon name="lucide:map-pin" class="w-4 h-4 text-purple-600 flex-shrink-0" />
                <span class="line-clamp-1">{{ event.venue }}</span>
              </div>
              <div class="flex items-center gap-2">
                <Icon name="lucide:eye" class="w-4 h-4 text-purple-600 flex-shrink-0" />
                <span>{{ event.view_count || 0 }} views</span>
                <span class="mx-1">•</span>
                <Icon name="lucide:calendar-clock" class="w-4 h-4 text-purple-600 flex-shrink-0" />
                <span>{{ formatRelativeTime(event.created_at) }}</span>
              </div>
            </div>

            <!-- Tags -->
            <div v-if="event.event_tags?.length" class="flex flex-wrap gap-1 mt-3">
              <span
                v-for="{ tag } in event.event_tags.slice(0, 3)"
                :key="tag.id"
                class="px-2 py-0.5 bg-purple-100 text-purple-700 rounded-full text-xs"
              >
                #{{ tag.name }}
              </span>
              <span
                v-if="event.event_tags.length > 3"
                class="px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full text-xs"
              >
                +{{ event.event_tags.length - 3 }}
              </span>
            </div>

            <!-- Footer Actions -->
            <div class="flex items-center justify-between mt-4 pt-3 border-t border-gray-100">
              <div class="flex items-center gap-2">
                <select
                  :value="event.status"
                  @change="handleStatusChange(event, $event)"
                  class="text-xs px-2 py-1 border border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 bg-white"
                >
                  <option value="draft">Draft</option>
                  <option value="published">Published</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </div>
              <div class="flex items-center gap-3">
                <NuxtLink
                  :to="`/event/${event.id}`"
                  class="text-xs text-purple-600 hover:text-purple-700 font-medium flex items-center gap-1"
                >
                  View Details
                  <Icon name="lucide:chevron-right" class="w-3 h-3" />
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-16 bg-white rounded-2xl shadow-xl">
        <Icon name="lucide:calendar-plus" class="w-20 h-20 text-gray-300 mx-auto mb-4" />
        <h3 class="text-2xl font-semibold text-gray-700 mb-2">No Events Found</h3>
        <p class="text-gray-500 mb-6 max-w-md mx-auto">
          You haven't created any events yet. Start by creating your first event!
        </p>
        <button
          @click="showCreateModal = true"
          class="px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-xl hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <Icon name="lucide:plus" class="w-5 h-5 inline mr-2" />
          Create Your First Event
        </button>
      </div>

      <!-- Pagination -->
      <div v-if="filteredEvents.length > itemsPerPage" class="flex justify-center mt-8">
        <div class="flex gap-2 bg-white rounded-lg shadow-sm p-1">
          <button
            v-for="page in totalPages"
            :key="page"
            @click="goToPage(page)"
            :class="[
              'px-4 py-2 rounded-lg transition-all text-sm font-medium',
              page === currentPage
                ? 'bg-purple-600 text-white shadow-md'
                : 'text-gray-700 hover:bg-gray-100'
            ]"
          >
            {{ page }}
          </button>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <EventFormModal
      v-if="showCreateModal || showEditModal"
      :event="selectedEvent"
      :is-edit="!!selectedEvent"
      @close="closeModal"
      @success="handleFormSuccess"
    />

    <!-- Delete Confirmation Modal -->
    <DeleteConfirmationModal
      v-if="showDeleteModal"
      :event="eventToDelete"
      @confirm="handleDeleteConfirm"
      @cancel="closeDeleteModal"
    />

    <!-- Quick Stats Modal -->
    <EventStatsModal
      v-if="showStatsModal"
      :event="selectedEvent"
      @close="showStatsModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useEventStore } from '~/stores/event'
import { useCategoryStore } from '~/stores/category'

// Components
import EventFormModal from '~/components/events/EventFormModal.vue'
import DeleteConfirmationModal from '~/components/events/DeleteConfirmationModal.vue'
import EventStatsModal from '~/components/events/EventStatsModal.vue'

definePageMeta({
  middleware: 'auth'
})

// Stores
const authStore = useAuthStore()
const eventStore = useEventStore()
const categoryStore = useCategoryStore()

// UI State
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showDeleteModal = ref(false)
const showStatsModal = ref(false)
const selectedEvent = ref<any>(null)
const eventToDelete = ref<any>(null)
const searchQuery = ref('')
const loading = ref(false)
const currentPage = ref(1)
const itemsPerPage = 9

let searchTimeout: NodeJS.Timeout

// Filters
const filters = ref({
  status: '',
  type: '',
  sortBy: 'event_date'
})

// Computed
const userId = computed(() => authStore.user?.id)

const myEvents = computed(() => {
  return eventStore.events.filter(event => event.user_id === userId.value)
})

const filteredEvents = computed(() => {
  let events = [...myEvents.value]

  // Search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    events = events.filter(event =>
      event.title.toLowerCase().includes(query) ||
      event.description.toLowerCase().includes(query) ||
      event.venue.toLowerCase().includes(query)
    )
  }

  // Status filter
  if (filters.value.status) {
    events = events.filter(event => event.status === filters.value.status)
  }

  // Type filter
  if (filters.value.type) {
    events = events.filter(event => 
      filters.value.type === 'free' ? event.is_free : !event.is_free
    )
  }

  // Sort
  const sortField = filters.value.sortBy
  events.sort((a, b) => {
    if (sortField === 'event_date') {
      return new Date(a.event_date).getTime() - new Date(b.event_date).getTime()
    }
    if (sortField === 'view_count') {
      return (b.view_count || 0) - (a.view_count || 0)
    }
    if (sortField === 'title') {
      return a.title.localeCompare(b.title)
    }
    return new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
  })

  return events
})

const paginatedEvents = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return filteredEvents.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(filteredEvents.value.length / itemsPerPage)
})

// Methods
const fetchMyEvents = async () => {
  loading.value = true
  try {
    await eventStore.fetchEvents()
  } catch (error) {
    console.error('Error fetching events:', error)
  } finally {
    loading.value = false
  }
}

const applyFilters = () => {
  currentPage.value = 1
}

const debouncedSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    currentPage.value = 1
  }, 500)
}

const resetFilters = () => {
  filters.value = {
    status: '',
    type: '',
    sortBy: 'event_date'
  }
  searchQuery.value = ''
  currentPage.value = 1
}

const goToPage = (page: number) => {
  currentPage.value = page
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

const openEditModal = (event: any) => {
  selectedEvent.value = event
  showEditModal.value = true
}

const closeModal = () => {
  showCreateModal.value = false
  showEditModal.value = false
  selectedEvent.value = null
}

const handleFormSuccess = () => {
  closeModal()
  fetchMyEvents()
}

const confirmDelete = (event: any) => {
  eventToDelete.value = event
  showDeleteModal.value = true
}

const closeDeleteModal = () => {
  showDeleteModal.value = false
  eventToDelete.value = null
}

const handleDeleteConfirm = async () => {
  if (eventToDelete.value) {
    const result = await eventStore.deleteEvent(eventToDelete.value.id)
    if (result.success) {
      closeDeleteModal()
      await fetchMyEvents()
    }
  }
}

const handleStatusChange = async (event: any, $event: any) => {
  const status = $event.target.value
  const result = await eventStore.updateEventStatus(event.id, status)
  if (result.success) {
    await fetchMyEvents()
  }
}

const formatDate = (date: string) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
}

const formatRelativeTime = (date: string) => {
  if (!date) return ''
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
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    published: 'bg-green-500/90 text-white',
    draft: 'bg-yellow-500/90 text-white',
    cancelled: 'bg-red-500/90 text-white'
  }
  return colors[status] || 'bg-gray-500/90 text-white'
}

// Lifecycle
onMounted(async () => {
  await fetchMyEvents()
  await categoryStore.fetchCategories()
})

// Watch for changes
watch(() => authStore.user?.id, () => {
  fetchMyEvents()
})
</script>

<style scoped>
/* Custom scrollbar */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

::-webkit-scrollbar-thumb {
  background: #c4b5d4;
  border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
  background: #a78bfa;
}

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
</style>