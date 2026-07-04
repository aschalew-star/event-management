<!-- pages/my-following.vue -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-50 via-white to-pink-50">
    <!-- Header -->
    <header class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Following</h1>
            <p class="text-sm text-gray-500">Events you're following</p>
          </div>
          <NuxtLink
            to="/my-followers"
            class="flex items-center gap-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
          >
            <Icon name="lucide:users" class="w-4 h-4" />
            My Followers
          </NuxtLink>
        </div>
      </div>
    </header>

    <!-- Stats -->
    <div class="max-w-7xl mx-auto px-4 pt-6">
      <div class="bg-white rounded-xl shadow-sm p-4 mb-6">
        <div class="flex items-center justify-between">
          <div>
            <span class="text-sm text-gray-500">Following</span>
            <p class="text-2xl font-bold text-gray-900">{{ followingCount }}</p>
          </div>
          <div>
            <span class="text-sm text-gray-500">Events</span>
            <p class="text-2xl font-bold text-gray-900">{{ followings.length }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 pb-8">
      <!-- Loading -->
      <div v-if="loading" class="flex justify-center py-16">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
      </div>

      <!-- Following Grid -->
      <div v-else-if="followings.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="follow in followings"
          :key="follow.id"
          class="group bg-white rounded-2xl shadow-xl overflow-hidden hover:shadow-2xl transition-all duration-300 hover:-translate-y-1"
        >
          <NuxtLink :to="`/events/${follow.event.id}`" class="block">
            <!-- Image -->
            <div class="relative h-48 overflow-hidden">
              <img
                :src="follow.event.featured_image || '/images/placeholder-event.jpg'"
                :alt="follow.event.title"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                @error="handleImageError"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
              
              <div class="absolute top-3 left-3 flex gap-2">
                <span
                  v-if="follow.event.is_free"
                  class="px-2.5 py-1 bg-green-500/90 backdrop-blur-sm text-white rounded-full text-xs font-medium"
                >
                  Free
                </span>
                <span
                  v-else
                  class="px-2.5 py-1 bg-purple-500/90 backdrop-blur-sm text-white rounded-full text-xs font-medium"
                >
                  ${{ Number(follow.event.price || 0).toFixed(2) }}
                </span>
              </div>

              <!-- Unfollow button -->
              <button
                @click.prevent="unfollow(follow.event.id)"
                class="absolute top-3 right-3 p-2 bg-white/90 backdrop-blur-sm rounded-full hover:bg-red-50 transition-colors shadow-lg"
                title="Unfollow"
              >
                <Icon name="lucide:heart-off" class="w-4 h-4 text-red-500" />
              </button>
            </div>

            <!-- Content -->
            <div class="p-5">
              <div class="flex items-start justify-between">
                <h3 class="font-semibold text-gray-900 line-clamp-1 flex-1">
                  {{ follow.event.title }}
                </h3>
              </div>

              <div class="flex items-center gap-1 mt-1">
                <span class="text-xs text-gray-500">by</span>
                <span class="text-xs font-medium text-purple-600">
                  {{ follow.event.user?.name || 'Unknown' }}
                </span>
              </div>

              <p class="text-sm text-gray-600 line-clamp-2 mt-2">
                {{ follow.event.description }}
              </p>

              <div class="space-y-1.5 mt-3 text-sm text-gray-500">
                <div class="flex items-center gap-2">
                  <Icon name="lucide:calendar" class="w-4 h-4 text-purple-600 flex-shrink-0" />
                  <span>{{ formatDate(follow.event.event_date) }}</span>
                </div>
                <div class="flex items-center gap-2">
                  <Icon name="lucide:map-pin" class="w-4 h-4 text-purple-600 flex-shrink-0" />
                  <span class="line-clamp-1">{{ follow.event.venue || 'Venue TBD' }}</span>
                </div>
              </div>

              <div class="mt-3 pt-3 border-t border-gray-100 flex items-center justify-between">
                <span class="text-xs text-gray-400">
                  Following since {{ formatRelativeTime(follow.created_at) }}
                </span>
                <NuxtLink
                  :to="`/follower/${follow.event.user_id}/events`"
                  class="text-xs text-purple-600 hover:text-purple-700 font-medium"
                >
                  View Creator →
                </NuxtLink>
              </div>
            </div>
          </NuxtLink>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-16 bg-white rounded-2xl shadow-xl">
        <Icon name="lucide:heart" class="w-20 h-20 text-gray-300 mx-auto mb-4" />
        <h3 class="text-2xl font-semibold text-gray-700 mb-2">Not Following Any Events</h3>
        <p class="text-gray-500 mb-6 max-w-md mx-auto">
          Start following events to get updates and never miss out!
        </p>
        <NuxtLink
          to="/events"
          class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-xl hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <Icon name="lucide:search" class="w-5 h-5" />
          Discover Events
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useQuery, useMutation } from '@vue/apollo-composable'
import { useToast } from 'vue-toastification'
import { GET_USER_FOLLOWS } from '~/graphql/eventQueries'
import { UNFOLLOW_EVENT } from '~/graphql/eventMutations'

definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()
const toast = useToast()

const userId = computed(() => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    return authStore.user.id
  }
  return null
})

// Fetch user follows
const { 
  result: followsResult, 
  loading, 
  refetch,
  onError 
} = useQuery(
  GET_USER_FOLLOWS,
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
  console.error('Error fetching follows:', error)
  toast.error('Failed to load following list')
})

const followings = computed(() => {
  return followsResult.value?.follows || []
})

const followingCount = computed(() => followings.value.length)

// Unfollow mutation
const { mutate: unfollowMutation, loading: unfollowLoading } = useMutation(UNFOLLOW_EVENT)

const unfollow = async (eventId: string) => {
  if (!eventId) {
    toast.error('Invalid event ID')
    return
  }

  if (unfollowLoading.value) {
    return
  }

  try {
    const result = await unfollowMutation({
      eventId: eventId
    })
    
    console.log('Unfollow result:', result)
    
    // Refetch the list
    await refetch()
    
    toast.success('Unfollowed event successfully')
  } catch (error: any) {
    console.error('Error unfollowing:', error)
    toast.error(error.message || 'Failed to unfollow event')
  }
}

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/images/placeholder-event.jpg'
}

const formatDate = (date: string) => {
  if (!date) return ''
  try {
    return new Date(date).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    })
  } catch {
    return ''
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

// Watch for authentication changes
watch(
  () => authStore.isAuthenticated,
  async (authenticated) => {
    if (!authenticated && process.client) {
      // Optionally redirect or handle logout
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
</script>