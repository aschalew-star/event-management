<!-- pages/my-followers.vue -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-purple-50 via-white to-pink-50">
    <!-- Header -->
    <header class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-40">
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">My Followers</h1>
            <p class="text-sm text-gray-500">People who follow your events</p>
          </div>
          <div class="flex items-center gap-4">
            <NuxtLink
              to="/my-following"
              class="flex items-center gap-2 px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
            >
              <Icon name="lucide:heart" class="w-4 h-4" />
              Following
            </NuxtLink>
          </div>
        </div>
      </div>
    </header>

    <!-- Stats -->
    <div class="max-w-7xl mx-auto px-4 pt-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">Total Followers</span>
            <Icon name="lucide:users" class="w-5 h-5 text-purple-600" />
          </div>
          <p class="text-2xl font-bold text-gray-900 mt-2">{{ totalFollowers }}</p>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">Events Followed</span>
            <Icon name="lucide:calendar" class="w-5 h-5 text-blue-600" />
          </div>
          <p class="text-2xl font-bold text-gray-900 mt-2">{{ totalEventsFollowed }}</p>
        </div>
        <div class="bg-white rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">Active Followers</span>
            <Icon name="lucide:activity" class="w-5 h-5 text-green-600" />
          </div>
          <p class="text-2xl font-bold text-gray-900 mt-2">{{ activeFollowers }}</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 pb-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-16">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-16 bg-white rounded-2xl shadow-xl">
        <Icon name="lucide:alert-circle" class="w-20 h-20 text-red-400 mx-auto mb-4" />
        <h3 class="text-2xl font-semibold text-gray-700 mb-2">Failed to Load Followers</h3>
        <p class="text-gray-500 mb-6">{{ error.message || 'Something went wrong' }}</p>
        <button
          @click="refetch"
          class="inline-flex items-center gap-2 px-6 py-3 bg-purple-600 text-white rounded-xl hover:bg-purple-700 transition-colors"
        >
          <Icon name="lucide:refresh-cw" class="w-5 h-5" />
          Try Again
        </button>
      </div>

      <!-- Not Logged In State -->
      <div v-else-if="!authStore.isAuthenticated" class="text-center py-16 bg-white rounded-2xl shadow-xl">
        <Icon name="lucide:lock" class="w-20 h-20 text-gray-300 mx-auto mb-4" />
        <h3 class="text-2xl font-semibold text-gray-700 mb-2">Please Log In</h3>
        <p class="text-gray-500 mb-6 max-w-md mx-auto">
          Log in to see who's following your events.
        </p>
        <NuxtLink
          to="/login"
          class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-xl hover:shadow-lg transition-all"
        >
          <Icon name="lucide:log-in" class="w-5 h-5" />
          Log In
        </NuxtLink>
      </div>

      <!-- Followers Grid -->
      <div v-else-if="followers.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="follower in followers"
          :key="follower.id"
          class="group bg-white rounded-2xl shadow-xl overflow-hidden hover:shadow-2xl transition-all duration-300 hover:-translate-y-1"
        >
          <div class="p-6">
            <div class="flex items-start gap-4">
              <!-- Avatar -->
              <div class="relative">
                <img
                  :src="follower.user?.avatar_url || '/images/default-avatar.jpg'"
                  :alt="follower.user?.name"
                  class="w-16 h-16 rounded-full object-cover border-2 border-purple-100 group-hover:border-purple-400 transition-colors"
                  @error="handleAvatarError"
                />
                <div class="absolute -bottom-1 -right-1 w-4 h-4 bg-green-500 rounded-full border-2 border-white"></div>
              </div>

              <!-- User Info -->
              <div class="flex-1 min-w-0">
                <h3 class="font-semibold text-gray-900 truncate">
                  {{ follower.user?.name || 'Anonymous' }}
                </h3>
                <p class="text-sm text-gray-500 truncate">{{ follower.user?.email }}</p>
                <p class="text-xs text-gray-400 mt-1">
                  Following since {{ formatDate(follower.created_at) }}
                </p>
              </div>
            </div>

            <!-- Bio -->
            <p v-if="follower.user?.bio" class="text-sm text-gray-600 mt-3 line-clamp-2">
              {{ follower.user.bio }}
            </p>

            <!-- Event Info -->
            <div v-if="follower.event" class="mt-3 pt-3 border-t border-gray-100">
              <div class="flex items-center gap-2 text-sm">
                <Icon name="lucide:calendar" class="w-4 h-4 text-purple-600 flex-shrink-0" />
                <span class="text-gray-600">Following:</span>
                <NuxtLink
                  :to="`/events/${follower.event.id}`"
                  class="text-purple-600 hover:text-purple-700 font-medium truncate"
                >
                  {{ follower.event.title }}
                </NuxtLink>
              </div>
            </div>

            <!-- Actions -->
            <div class="mt-4 flex items-center justify-between">
              <div class="flex items-center gap-2">
                <NuxtLink
                  :to="`/follower/${follower.user_id}/events`"
                  class="px-3 py-1.5 text-sm bg-purple-100 text-purple-700 rounded-lg hover:bg-purple-200 transition-colors"
                >
                  View Events
                </NuxtLink>
              </div>
              <button
                @click="viewFollowerDetails(follower)"
                class="text-xs text-gray-400 hover:text-gray-600 transition-colors"
              >
                <Icon name="lucide:chevron-right" class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-16 bg-white rounded-2xl shadow-xl">
        <Icon name="lucide:users" class="w-20 h-20 text-gray-300 mx-auto mb-4" />
        <h3 class="text-2xl font-semibold text-gray-700 mb-2">No Followers Yet</h3>
        <p class="text-gray-500 mb-6 max-w-md mx-auto">
          You don't have any followers yet. Create amazing events to attract followers!
        </p>
        <NuxtLink
          to="/my-events"
          class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-xl hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <Icon name="lucide:plus" class="w-5 h-5" />
          Create Events
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onActivated } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@vue/apollo-composable'
import { useAuthStore } from '~/stores/auth'
import { useToast } from 'vue-toastification'
import { GET_USER_FOLLOWERS } from '~/graphql/eventQueries'

definePageMeta({
  middleware: 'auth'
})

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const userId = computed(() => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    return authStore.user.id
  }
  return null
})

// Fetch followers for the current user
const { 
  result: followersResult, 
  loading, 
  error,
  refetch,
  onError 
} = useQuery(
  GET_USER_FOLLOWERS,
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
  console.error('Error fetching followers:', error)
  toast.error('Failed to load followers')
})

const followers = computed(() => {
  return followersResult.value?.follows || []
})

// Stats
const totalFollowers = computed(() => followers.value.length)

const totalEventsFollowed = computed(() => {
  // Count unique events followed by all followers
  const eventIds = new Set()
  followers.value.forEach(f => {
    if (f.event_id) {
      eventIds.add(f.event_id)
    }
  })
  return eventIds.size
})

const activeFollowers = computed(() => {
  return followers.value.filter(f => f.user).length
})

const handleAvatarError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/images/default-avatar.jpg'
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

const viewFollowerDetails = (follower: any) => {
  router.push(`/follower/${follower.user_id}/events`)
}

// Watch for authentication changes
watch(
  () => authStore.isAuthenticated,
  async (authenticated) => {
    if (authenticated && userId.value) {
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

// Refetch when component is activated (coming back from navigation)
onActivated(async () => {
  if (authStore.isAuthenticated && userId.value) {
    await refetch()
  }
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>