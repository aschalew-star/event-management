<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 py-8 transition-colors duration-300">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header with Stats -->
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">My Bookmarks</h1>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {{ bookmarks.length }} events saved for later
          </p>
        </div>
        <div class="flex items-center gap-3">
          <button 
            @click="manualRefresh" 
            :disabled="loading || manualRefreshing"
            class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <Icon :name="manualRefreshing ? 'lucide:loader-2' : 'lucide:refresh-cw'" class="w-4 h-4" :class="manualRefreshing ? 'animate-spin' : ''" />
            Refresh
          </button>
          <NuxtLink
            to="/events"
            class="px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg transition-colors flex items-center gap-2"
          >
            <Icon name="lucide:search" class="w-4 h-4" />
            Browse Events
          </NuxtLink>
        </div>
      </div>

      <!-- Stats Row -->
      <div v-if="bookmarks.length > 0" class="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-6">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Total Bookmarks</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ bookmarks.length }}</p>
            </div>
            <div class="w-10 h-10 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:bookmark" class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
            </div>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Upcoming Events</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ upcomingCount }}</p>
            </div>
            <div class="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:calendar-clock" class="w-5 h-5 text-green-600 dark:text-green-400" />
            </div>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Free Events</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ freeCount }}</p>
            </div>
            <div class="w-10 h-10 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:gift" class="w-5 h-5 text-amber-600 dark:text-amber-400" />
            </div>
          </div>
        </div>
      </div>
      
      <ClientOnly>
        <!-- Loading State -->
        <div v-if="loading" class="flex justify-center py-16">
          <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
        </div>

        <!-- Not Logged In State -->
        <div v-else-if="!authStore.isAuthenticated" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700">
          <div class="text-5xl mb-4 text-gray-300 dark:text-gray-600">
            <Icon name="lucide:lock" class="w-16 h-16" />
          </div>
          <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-200 mb-2">Please Log In</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Log in to view and manage your bookmarks</p>
          <NuxtLink 
            to="/login"
            class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium rounded-xl hover:shadow-lg transition-all transform hover:scale-[1.02]"
          >
            <Icon name="lucide:log-in" class="w-4 h-4" />
            Log In
          </NuxtLink>
        </div>

        <!-- Logged In but No Bookmarks -->
        <div v-else-if="bookmarks.length === 0 && !loading" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700">
          <div class="text-5xl mb-4 text-gray-300 dark:text-gray-600">
            <Icon name="lucide:bookmark" class="w-16 h-16" />
          </div>
          <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-200 mb-2">No Bookmarks Yet</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Start bookmarking events you're interested in!</p>
          <NuxtLink 
            to="/events"
            class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium rounded-xl hover:shadow-lg transition-all transform hover:scale-[1.02]"
          >
            <Icon name="lucide:search" class="w-4 h-4" />
            Browse Events
          </NuxtLink>
        </div>

        <!-- Display Bookmarks with Delete Button on Page -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div 
            v-for="bookmark in bookmarks" 
            :key="bookmark.id"
            class="relative group"
          >
            <!-- Delete Button - Positioned at top right of the card container -->
            <button 
              @click="removeBookmark(bookmark.event.id)"
              :disabled="deletingBookmark === bookmark.event.id"
              class="absolute -top-2 -right-2 z-10 p-1.5 bg-red-500 hover:bg-red-600 disabled:bg-red-300 text-white rounded-full shadow-lg transition-all hover:scale-110 disabled:cursor-not-allowed"
              title="Remove Bookmark"
              aria-label="Remove bookmark"
            >
              <Icon :name="deletingBookmark === bookmark.event.id ? 'lucide:loader-2' : 'lucide:x'" class="w-4 h-4" :class="deletingBookmark === bookmark.event.id ? 'animate-spin' : ''" />
            </button>

            <!-- Event Card -->
            <EventCard
              :event="enrichEventWithBookmarkData(bookmark)"
              @click="navigateToEvent(bookmark.event.id)"
            />
          </div>
        </div>

        <template #fallback>
          <div class="flex justify-center py-16">
            <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
          </div>
        </template>
      </ClientOnly>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch, onMounted, onActivated, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery, useMutation } from '@vue/apollo-composable'
import { useAuthStore } from '~/stores/auth'
import { useToast } from 'vue-toastification'
import { GET_USER_BOOKMARKS } from '~/graphql/eventQueries'
import { UNBOOKMARK_EVENT } from '~/graphql/eventMutations'
import EventCard from '~/components/events/EventCard.vue'

definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()
const router = useRouter()
const toast = useToast()

const manualRefreshing = ref(false)
const deletingBookmark = ref<string | null>(null)

const userId = computed(() => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    return authStore.user.id
  }
  return null
})

// Fetch User Bookmarks
const { 
  result: bookmarksResult, 
  loading, 
  refetch,
  onError
} = useQuery(
  GET_USER_BOOKMARKS,
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
  console.error('GraphQL Query Error:', error)
  toast.error('Failed to load bookmarks: ' + error.message)
})

const bookmarks = computed(() => {
  return bookmarksResult.value?.bookmarks || []
})

// Computed stats
const upcomingCount = computed(() => {
  const now = new Date()
  return bookmarks.value.filter((b: any) => {
    const eventDate = new Date(b.event.event_date)
    return eventDate >= now && b.event.status !== 'cancelled'
  }).length
})

const freeCount = computed(() => {
  return bookmarks.value.filter((b: any) => b.event.is_free).length
})

// Unbookmark Mutation
const { mutate: unbookmarkMutation, loading: unbookmarkLoading } = useMutation(UNBOOKMARK_EVENT)

const removeBookmark = async (eventId: string) => {
  try {
    if (!eventId) {
      toast.error('Invalid event ID')
      return
    }

    if (unbookmarkLoading.value) {
      return
    }

    // Set loading state for this specific bookmark
    deletingBookmark.value = eventId

    // Get the current user ID
    const currentUserId = authStore.user?.id
    if (!currentUserId) {
      toast.error('User not authenticated')
      return
    }

    console.log('Removing bookmark - Event ID:', eventId, 'User ID:', currentUserId)

    const result = await unbookmarkMutation({
      eventId: eventId,
      userId: currentUserId
    })
    
    console.log('Unbookmark result:', result)
    
    if (result?.data?.delete_bookmarks?.affected_rows > 0) {
      await refetch()
      toast.success('Removed from your bookmarks')
    } else {
      toast.error('Failed to remove bookmark')
    }
  } catch (error: any) {
    console.error('Error removing bookmark:', error)
    toast.error(error.message || 'Failed to remove bookmark')
  } finally {
    deletingBookmark.value = null
  }
}

const navigateToEvent = (eventId: string) => {
  router.push(`/events/${eventId}`)
}

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/images/event-placeholder.jpg'
}

const manualRefresh = async () => {
  if (authStore.isAuthenticated && userId.value) {
    manualRefreshing.value = true
    try {
      await refetch()
      toast.success('Bookmarks refreshed!')
    } catch (error) {
      console.error('Error refreshing:', error)
      toast.error('Failed to refresh bookmarks')
    } finally {
      manualRefreshing.value = false
    }
  }
}

// Enrich event with bookmark data for the EventCard
const enrichEventWithBookmarkData = (bookmark: any) => {
  const event = bookmark.event
  return {
    ...event,
    // Add bookmark count from aggregate
    bookmarks_count: event.bookmarks_aggregate?.aggregate?.count || 0,
    // Ensure images array is properly structured
    images: event.event_images || [],
    // Set featured image from event_images
    featured_image: event.event_images?.find((img: any) => img.is_featured)?.image_url || 
                     event.event_images?.[0]?.image_url || 
                     null
  }
}

// Watch for authentication changes
watch(
  () => authStore.isAuthenticated,
  async (authenticated) => {
    if (!authenticated && process.client) {
      // Don't redirect, just show the login state
    } else if (authenticated && userId.value) {
      await refetch()
    }
  }
)

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

/* Smooth hover transitions */
.relative {
  transition: all 0.2s ease;
}

.absolute {
  opacity: 0.9;
  transition: all 0.2s ease;
}

.absolute:hover {
  opacity: 1;
}
</style>