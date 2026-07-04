<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 py-8 transition-colors duration-300">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">My Bookmarks</h1>
        <button 
          @click="manualRefresh" 
          :disabled="loading"
          class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          <i :class="loading ? 'fas fa-spinner fa-spin' : 'fas fa-sync-alt'"></i>
          Refresh
        </button>
      </div>
      
      <ClientOnly>
        <!-- Loading State -->
        <div v-if="loading" class="flex justify-center py-16">
          <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
        </div>

        <!-- Not Logged In State -->
        <div v-else-if="!authStore.isAuthenticated" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700">
          <div class="text-5xl mb-4 text-gray-300 dark:text-gray-600">
            <i class="fas fa-lock"></i>
          </div>
          <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-200 mb-2">Please Log In</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Log in to view and manage your bookmarks</p>
          <NuxtLink 
            to="/login"
            class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium rounded-xl hover:shadow-lg transition-all transform hover:scale-[1.02]"
          >
            <i class="fas fa-sign-in-alt mr-1"></i>
            Log In
          </NuxtLink>
        </div>

        <!-- Logged In but No Bookmarks -->
        <div v-else-if="bookmarks.length === 0 && !loading" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700">
          <div class="text-5xl mb-4 text-gray-300 dark:text-gray-600">
            <i class="far fa-bookmark"></i>
          </div>
          <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-200 mb-2">No Bookmarks Yet</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Start bookmarking events you're interested in!</p>
          <NuxtLink 
            to="/events"
            class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium rounded-xl hover:shadow-lg transition-all transform hover:scale-[1.02]"
          >
            <i class="fas fa-search mr-1"></i>
            Browse Events
          </NuxtLink>
        </div>

        <!-- Display Bookmarks -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div 
            v-for="bookmark in bookmarks" 
            :key="bookmark.id"
            class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg overflow-hidden border border-gray-100 dark:border-gray-700 hover:shadow-2xl transition-all hover:-translate-y-1 group cursor-pointer"
            @click="navigateToEvent(bookmark.event.id)"
          >
            <div class="relative h-48 bg-gray-200 dark:bg-gray-700">
              <img 
                :src="bookmark.event.featured_image || '/images/event-placeholder.jpg'" 
                :alt="bookmark.event.title"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                @error="handleImageError"
              />
              <div class="absolute top-3 right-3">
                <span 
                  class="px-2.5 py-1 rounded-full text-xs font-semibold backdrop-blur-md shadow-sm"
                  :class="bookmark.event.is_free ? 'bg-green-500/90 text-white' : 'bg-indigo-600/90 text-white'"
                >
                  {{ bookmark.event.is_free ? 'Free' : `$${Number(bookmark.event.price || 0).toFixed(2)}` }}
                </span>
              </div>
            </div>
            
            <div class="p-5">
              <h3 class="font-bold text-gray-900 dark:text-white text-lg mb-2 line-clamp-1 group-hover:text-indigo-500 dark:group-hover:text-indigo-400 transition-colors">
                {{ bookmark.event.title }}
              </h3>
              
              <div class="space-y-2 text-sm text-gray-600 dark:text-gray-300 mb-4">
                <div class="flex items-center gap-2.5">
                  <i class="fas fa-calendar-day text-indigo-500 w-4 text-center"></i>
                  <span>{{ formatDate(bookmark.event.event_date) }}</span>
                </div>
                <div class="flex items-center gap-2.5">
                  <i class="fas fa-map-pin text-indigo-500 w-4 text-center"></i>
                  <span class="line-clamp-1">{{ bookmark.event.venue || 'Venue TBD' }}</span>
                </div>
              </div>
              
              <div class="pt-3 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between">
                <span class="text-xs text-gray-400 dark:text-gray-500">
                  <i class="far fa-clock mr-1"></i>
                  Saved {{ formatRelativeTime(bookmark.created_at) }}
                </span>
                <button 
                  @click.stop="removeBookmark(bookmark.event.id)"
                  class="w-8 h-8 rounded-full bg-red-50 dark:bg-red-950/30 text-red-500 hover:text-white hover:bg-red-500 dark:hover:bg-red-600 flex items-center justify-center transition-all shadow-sm"
                  title="Remove Bookmark"
                  aria-label="Remove bookmark"
                >
                  <i class="fas fa-trash-alt text-xs"></i>
                </button>
              </div>
            </div>
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

definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()
const router = useRouter()
const toast = useToast()

const manualRefreshing = ref(false)

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

    await unbookmarkMutation({ eventId })
    await refetch()
    toast.success('Removed from your bookmarks')
  } catch (error: any) {
    console.error('Error removing bookmark:', error)
    toast.error(error.message || 'Failed to remove bookmark')
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

// Watch for authentication changes
watch(
  () => authStore.isAuthenticated,
  async (authenticated) => {
    if (!authenticated && process.client) {
      router.push('/login')
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

// Refetch when component is activated (coming back from event details)
onActivated(async () => {
  if (authStore.isAuthenticated && userId.value) {
    await refetch()
  }
})

const formatDate = (dateStr: string) => {
  if (!dateStr) return 'TBD'
  try {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    })
  } catch {
    return 'Invalid date'
  }
}

const formatRelativeTime = (dateStr: string) => {
  if (!dateStr) return ''
  
  try {
    const now = new Date()
    const past = new Date(dateStr)
    const diff = Math.floor((now.getTime() - past.getTime()) / (1000 * 60 * 60 * 24))
    
    if (diff <= 0) return 'today'
    if (diff === 1) return 'yesterday'
    if (diff < 7) return `${diff} days ago`
    if (diff < 30) return `${Math.floor(diff / 7)} weeks ago`
    
    return formatDate(dateStr)
  } catch {
    return ''
  }
}
</script>