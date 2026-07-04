<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50">
    <!-- Header -->
    <header class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <NuxtLink
              to="/my-following"
              class="flex items-center gap-2 text-gray-600 hover:text-gray-900 transition-colors"
            >
              <Icon name="lucide:arrow-left" class="w-5 h-5" />
              Back to Following
            </NuxtLink>
            <div class="h-6 w-px bg-gray-300"></div>
            <div>
              <h1 class="text-2xl font-bold text-gray-900">{{ creatorName }}'s Events</h1>
              <p class="text-sm text-gray-500">{{ eventsCount }} events created</p>
            </div>
          </div>
          <NuxtLink
            to="/events"
            class="flex items-center gap-2 px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg hover:shadow-lg transition-all hover:-translate-y-0.5"
          >
            <Icon name="lucide:search" class="w-4 h-4" />
            Browse All Events
          </NuxtLink>
        </div>
      </div>
    </header>

    <!-- Stats -->
    <div class="max-w-7xl mx-auto px-4 pt-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="bg-white rounded-xl shadow-sm p-4">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm text-gray-500">Total Events</span>
              <p class="text-2xl font-bold text-gray-900">{{ eventsCount }}</p>
            </div>
            <Icon name="lucide:calendar" class="w-8 h-8 text-blue-500 opacity-50" />
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm text-gray-500">Upcoming</span>
              <p class="text-2xl font-bold text-green-600">{{ upcomingEventsCount }}</p>
            </div>
            <Icon name="lucide:calendar-clock" class="w-8 h-8 text-green-500 opacity-50" />
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm text-gray-500">Completed</span>
              <p class="text-2xl font-bold text-gray-600">{{ completedEventsCount }}</p>
            </div>
            <Icon name="lucide:check-circle" class="w-8 h-8 text-gray-500 opacity-50" />
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4">
          <div class="flex items-center justify-between">
            <div>
              <span class="text-sm text-gray-500">Free Events</span>
              <p class="text-2xl font-bold text-green-600">{{ freeEventsCount }}</p>
            </div>
            <Icon name="lucide:gift" class="w-8 h-8 text-green-500 opacity-50" />
          </div>
        </div>
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
      <div v-else-if="events.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="event in events"
          :key="event.id"
          class="group bg-white rounded-2xl shadow-xl overflow-hidden hover:shadow-2xl transition-all duration-300 hover:-translate-y-1"
        >
          <NuxtLink :to="`/events/${event.id}`" class="block">
            <!-- Image -->
            <div class="relative h-48 overflow-hidden">
              <img
                :src="event.featured_image || '/images/placeholder-event.jpg'"
                :alt="event.title"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                @error="handleImageError"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
              
              <div class="absolute top-3 left-3 flex gap-2">
                <span
                  v-if="event.is_free"
                  class="px-2.5 py-1 bg-green-500/90 backdrop-blur-sm text-white rounded-full text-xs font-medium"
                >
                  Free
                </span>
                <span
                  v-else
                  class="px-2.5 py-1 bg-blue-500/90 backdrop-blur-sm text-white rounded-full text-xs font-medium"
                >
                  ${{ Number(event.price || 0).toFixed(2) }}
                </span>
                <span
                  v-if="event.status"
                  class="px-2.5 py-1 backdrop-blur-sm text-white rounded-full text-xs font-medium"
                  :class="statusColors[event.status] || 'bg-gray-500/90'"
                >
                  {{ event.status.charAt(0).toUpperCase() + event.status.slice(1) }}
                </span>
              </div>

              <div class="absolute bottom-3 left-3">
                <div class="flex items-center gap-2">
                  <span class="text-white text-xs font-medium bg-black/50 backdrop-blur-sm px-2 py-1 rounded-full">
                    <Icon name="lucide:calendar" class="w-3 h-3 inline mr-1" />
                    {{ formatDate(event.event_date) }}
                  </span>
                </div>
              </div>
            </div>

            <!-- Content -->
            <div class="p-5">
              <h3 class="font-semibold text-gray-900 line-clamp-1 group-hover:text-blue-600 transition-colors">
                {{ event.title }}
              </h3>

              <div class="flex items-center gap-1 mt-1">
                <span class="text-xs text-gray-500">by</span>
                <span class="text-xs font-medium text-blue-600">
                  {{ event.user?.name || 'Unknown' }}
                </span>
              </div>

              <p class="text-sm text-gray-600 line-clamp-2 mt-2">
                {{ event.description }}
              </p>

              <div class="space-y-1.5 mt-3 text-sm text-gray-500">
                <div class="flex items-center gap-2">
                  <Icon name="lucide:map-pin" class="w-4 h-4 text-blue-600 flex-shrink-0" />
                  <span class="line-clamp-1">{{ event.venue || 'Venue TBD' }}</span>
                </div>
              </div>

              <div class="mt-3 pt-3 border-t border-gray-100 flex items-center justify-between">
                <span class="text-xs text-gray-400">
                  Created {{ formatRelativeTime(event.created_at) }}
                </span>
                <div class="flex items-center gap-2">
                  <button
                    @click.prevent="toggleFollow(event.id)"
                    class="text-xs text-blue-600 hover:text-blue-700 font-medium flex items-center gap-1"
                  >
                    <Icon :name="isFollowing(event.id) ? 'lucide:heart' : 'lucide:heart-off'" class="w-3 h-3" />
                    {{ isFollowing(event.id) ? 'Following' : 'Follow' }}
                  </button>
                </div>
              </div>
            </div>
          </NuxtLink>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-16 bg-white rounded-2xl shadow-xl">
        <Icon name="lucide:calendar-x" class="w-20 h-20 text-gray-300 mx-auto mb-4" />
        <h3 class="text-2xl font-semibold text-gray-700 mb-2">No Events Found</h3>
        <p class="text-gray-500 mb-6 max-w-md mx-auto">
          This creator hasn't published any events yet.
        </p>
        <NuxtLink
          to="/events"
          class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-xl hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <Icon name="lucide:search" class="w-5 h-5" />
          Browse All Events
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onActivated } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery, useMutation } from '@vue/apollo-composable'
import { useAuthStore } from '~/stores/auth'
import { useToast } from 'vue-toastification'
import { GET_USER_EVENTS } from '~/graphql/eventQueries'
import { FOLLOW_EVENT, UNFOLLOW_EVENT } from '~/graphql/eventMutations'

definePageMeta({
  middleware: 'auth'
})

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

// Get creator ID from route params
const creatorId = computed(() => {
  const id = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id
  return id || null
})

const currentUserId = computed(() => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    return authStore.user.id
  }
  return null
})

// State for following statuses
const followingEvents = ref<Set<string>>(new Set())

// Fetch creator's events
const { 
  result: eventsResult, 
  loading, 
  error,
  refetch,
  onError 
} = useQuery(
  GET_USER_EVENTS,
  () => ({ 
    userId: creatorId.value || '00000000-0000-0000-0000-000000000000' 
  }),
  {
    fetchPolicy: 'network-only',
    skip: () => {
      return !process.client || !creatorId.value
    },
    notifyOnNetworkStatusChange: true,
  }
)

onError((error) => {
  console.error('Error fetching creator events:', error)
  toast.error('Failed to load creator events')
})

const events = computed(() => {
  return eventsResult.value?.events || []
})

const creatorName = computed(() => {
  if (events.value.length > 0 && events.value[0]?.user?.name) {
    return events.value[0].user.name
  }
  return 'Creator'
})

const eventsCount = computed(() => events.value.length)
const upcomingEventsCount = computed(() => {
  const now = new Date()
  return events.value.filter(e => {
    const eventDate = new Date(e.event_date)
    return eventDate >= now && e.status !== 'cancelled'
  }).length
})

const completedEventsCount = computed(() => {
  const now = new Date()
  return events.value.filter(e => {
    const eventDate = new Date(e.event_date)
    return eventDate < now || e.status === 'completed'
  }).length
})

const freeEventsCount = computed(() => {
  return events.value.filter(e => e.is_free).length
})

// Follow/Unfollow mutations
const { mutate: followMutation } = useMutation(FOLLOW_EVENT)
const { mutate: unfollowMutation } = useMutation(UNFOLLOW_EVENT)

const toggleFollow = async (eventId: string) => {
  if (!authStore.isAuthenticated) {
    toast.warning('Please sign in to follow events')
    router.push('/login')
    return
  }

  const userId = currentUserId.value
  if (!userId) {
    toast.error('User ID not found')
    return
  }

  try {
    if (followingEvents.value.has(eventId)) {
      // Unfollow
      await unfollowMutation({
        eventId: eventId
      })
      followingEvents.value.delete(eventId)
      toast.success('Unfollowed event')
    } else {
      // Follow
      await followMutation({
        eventId: eventId,
        userId: userId
      })
      followingEvents.value.add(eventId)
      toast.success('Following event!')
    }
  } catch (error: any) {
    console.error('Error toggling follow:', error)
    toast.error(error.message || 'Failed to update follow status')
  }
}

const isFollowing = (eventId: string) => {
  return followingEvents.value.has(eventId)
}

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/images/placeholder-event.jpg'
}

const statusColors: Record<string, string> = {
  published: 'bg-green-500/90',
  upcoming: 'bg-blue-500/90',
  completed: 'bg-gray-500/90',
  cancelled: 'bg-red-500/90',
  draft: 'bg-yellow-500/90',
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
  if (creatorId.value) {
    await refetch()
  }
})

// Refetch when component is activated
onActivated(async () => {
  if (creatorId.value) {
    await refetch()
  }
})

// Watch for creator ID changes
watch(
  creatorId,
  async (newId) => {
    if (newId) {
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
</style>