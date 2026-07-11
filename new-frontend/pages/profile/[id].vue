<!-- pages/profile/[id].vue -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center min-h-screen">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-16 w-16 border-4 border-indigo-500 border-t-transparent"></div>
        <p class="mt-4 text-gray-600 dark:text-gray-400 text-lg">Loading profile...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex items-center justify-center min-h-screen px-4">
      <div class="text-center max-w-md">
        <div class="text-6xl mb-4">😅</div>
        <h3 class="text-2xl font-bold text-gray-800 dark:text-gray-200 mb-2">User not found</h3>
        <p class="text-gray-600 dark:text-gray-400 mb-6">The user you're looking for doesn't exist.</p>
        <NuxtLink
          to="/events"
          class="inline-flex items-center gap-2 px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
        >
          <Icon name="lucide:arrow-left" class="w-4 h-4" />
          Back to Events
        </NuxtLink>
      </div>
    </div>

    <!-- Profile Content -->
    <div v-else-if="user" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Back Button -->
      <button
        @click="router.back()"
        class="inline-flex items-center gap-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 mb-6 transition-colors group"
      >
        <Icon name="lucide:arrow-left" class="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
        Back
      </button>

      <!-- Profile Header -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl overflow-hidden mb-8">
        <!-- Cover Image -->
        <div class="relative h-48 md:h-64 bg-gradient-to-r from-indigo-500 via-purple-500 to-pink-500">
          <div class="absolute inset-0 bg-[url('/images/pattern.svg')] opacity-10"></div>
          <div class="absolute -bottom-16 left-8 flex items-end gap-6">
            <!-- Avatar -->
            <div class="relative">
              <img
                :src="user.avatar_url || '/images/default-avatar.jpg'"
                :alt="user.name"
                class="w-32 h-32 rounded-full object-cover border-4 border-white dark:border-gray-800 shadow-lg"
                @error="handleAvatarError"
              />
              <div v-if="isCurrentUser" class="absolute -bottom-1 -right-1 bg-green-500 text-white text-xs font-bold px-2 py-0.5 rounded-full border-2 border-white dark:border-gray-800">
                You
              </div>
            </div>
            <div class="pb-4">
              <h1 class="text-2xl md:text-3xl font-bold text-white">{{ user.name }}</h1>
              <p class="text-white/80 text-sm flex items-center gap-2">
                <Icon name="lucide:calendar" class="w-3 h-3" />
                Joined {{ formatDate(user.created_at) }}
              </p>
            </div>
          </div>
        </div>
        <div class="h-16"></div>
      </div>

      <!-- Stats Row -->
      <div class="grid grid-cols-1 sm:grid-cols-4 gap-4 mb-8">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Events</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ user.events?.length || 0 }}</p>
            </div>
            <div class="w-10 h-10 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:calendar" class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
            </div>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Followers</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ followersCount }}</p>
            </div>
            <div class="w-10 h-10 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:users" class="w-5 h-5 text-purple-600 dark:text-purple-400" />
            </div>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Following</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ followingCount }}</p>
            </div>
            <div class="w-10 h-10 bg-pink-100 dark:bg-pink-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:heart" class="w-5 h-5 text-pink-600 dark:text-pink-400" />
            </div>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Total Views</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ totalViews }}</p>
            </div>
            <div class="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:eye" class="w-5 h-5 text-green-600 dark:text-green-400" />
            </div>
          </div>
        </div>
      </div>

      <!-- Bio -->
      <div v-if="user.bio" class="bg-white dark:bg-gray-800 rounded-xl shadow p-6 mb-8">
        <h3 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider mb-2">About</h3>
        <p class="text-gray-700 dark:text-gray-300 whitespace-pre-wrap">{{ user.bio }}</p>
      </div>

      <!-- Follow Button -->
      <div v-if="!isCurrentUser && authStore.isAuthenticated" class="mb-8">
        <button
          @click="toggleFollow"
          :disabled="followLoading"
          class="px-6 py-3 rounded-xl font-semibold transition-all flex items-center gap-2"
          :class="isFollowing 
            ? 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600' 
            : 'bg-gradient-to-r from-indigo-600 to-purple-600 text-white hover:shadow-lg hover:-translate-y-0.5'"
        >
          <Icon :name="isFollowing ? 'lucide:check' : 'lucide:user-plus'" class="w-5 h-5" />
          {{ isFollowing ? 'Following' : 'Follow' }}
        </button>
      </div>

      <!-- Events Section -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h2 class="text-xl font-bold text-gray-900 dark:text-white">Events by {{ user.name }}</h2>
            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
              {{ user.events?.length || 0 }} events created
            </p>
          </div>
          <div class="flex gap-2">
            <button
              @click="viewMode = 'grid'"
              class="p-2 rounded-lg transition-colors"
              :class="viewMode === 'grid' ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'"
            >
              <Icon name="lucide:layout-grid" class="w-5 h-5" />
            </button>
            <button
              @click="viewMode = 'list'"
              class="p-2 rounded-lg transition-colors"
              :class="viewMode === 'list' ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' : 'text-gray-400 hover:text-gray-600 dark:hover:text-gray-300'"
            >
              <Icon name="lucide:list" class="w-5 h-5" />
            </button>
          </div>
        </div>

        <!-- Loading Events -->
        <div v-if="loadingEvents" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
        </div>

        <!-- Events Grid -->
        <div v-else-if="enrichedEvents.length > 0" :class="viewMode === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6' : 'space-y-4'">
          <EventCard
            v-for="event in enrichedEvents"
            :key="event.id"
            :event="event"
            @click="navigateToEvent(event.id)"
          />
        </div>

        <!-- Empty State -->
        <div v-else class="text-center py-12">
          <Icon name="lucide:calendar" class="w-16 h-16 text-gray-300 dark:text-gray-600 mx-auto mb-4" />
          <p class="text-gray-500 dark:text-gray-400">{{ user.name }} hasn't created any events yet</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useQuery, useMutation } from '@vue/apollo-composable';
import { useAuthStore } from '~/stores/auth';
import { useToast } from 'vue-toastification';
import { GET_USER_PROFILE, GET_USER_EVENTS } from '~/graphql/eventQueries';
import { FOLLOW_USER, UNFOLLOW_USER } from '~/graphql/eventMutations';
import EventCard from '~/components/events/EventCard.vue';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

const viewMode = ref<'grid' | 'list'>('grid');
const followLoading = ref(false);
const isFollowing = ref(false);

const userId = computed(() => {
  const id = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id;
  return id || null;
});

const isCurrentUser = computed(() => {
  return authStore.isAuthenticated && authStore.user?.id === userId.value;
});

// Fetch User Profile
const {
  result: userResult,
  loading,
  error,
  refetch
} = useQuery(
  GET_USER_PROFILE,
  () => ({ id: userId.value }),
  {
    fetchPolicy: 'network-only',
    skip: () => !userId.value,
  }
);

const user = computed(() => userResult.value?.users_by_pk);

// Fetch User's Events
const {
  result: eventsResult,
  loading: loadingEvents,
  refetch: refetchEvents
} = useQuery(
  GET_USER_EVENTS,
  () => ({ userId: userId.value }),
  {
    fetchPolicy: 'network-only',
    skip: () => !userId.value,
  }
);

// Enrich events with images
const enrichedEvents = computed(() => {
  const events = eventsResult.value?.events || [];
  return events.map((event: any) => {
    const images = event.event_images || [];
    const featuredImage = images.find((img: any) => img.is_featured)?.image_url || 
                          images[0]?.image_url || 
                          null;
    return {
      ...event,
      featured_image: featuredImage,
      images: images,
      bookmarks_count: event.bookmarks_aggregate?.aggregate?.count || 0,
      tickets_count: event.tickets_aggregate?.aggregate?.count || 0
    };
  });
});

// Stats
const followersCount = computed(() => {
  return user.value?.followers_aggregate?.aggregate?.count || 0;
});

const followingCount = computed(() => {
  return user.value?.following_aggregate?.aggregate?.count || 0;
});

const totalViews = computed(() => {
  const events = eventsResult.value?.events || [];
  return events.reduce((sum: number, e: any) => sum + (e.view_count || 0), 0);
});

// Check if current user follows this user
const { refetch: refetchFollowStatus } = useQuery(
  GET_FOLLOW_STATUS,
  () => ({
    followerId: authStore.user?.id || '',
    followedUserId: userId.value || ''
  }),
  {
    skip: () => !authStore.isAuthenticated || !userId.value || isCurrentUser.value,
    fetchPolicy: 'network-only',
  }
);

// Mutations
const { mutate: followMutation } = useMutation(FOLLOW_USER);
const { mutate: unfollowMutation } = useMutation(UNFOLLOW_USER);

const toggleFollow = async () => {
  if (!authStore.isAuthenticated) {
    toast.warning('Please sign in to follow users');
    router.push('/login');
    return;
  }

  const currentUserId = authStore.user?.id;
  if (!currentUserId) return;

  followLoading.value = true;
  try {
    if (isFollowing.value) {
      const result = await unfollowMutation({
        followerId: currentUserId,
        followedUserId: userId.value
      });
      if (result?.data?.delete_follows?.affected_rows > 0) {
        isFollowing.value = false;
        toast.success(`Unfollowed ${user.value?.name}`);
        await refetch();
        await refetchFollowStatus();
      }
    } else {
      const result = await followMutation({
        followerId: currentUserId,
        followedUserId: userId.value
      });
      if (result?.data?.insert_follows_one) {
        isFollowing.value = true;
        toast.success(`Now following ${user.value?.name}!`);
        await refetch();
        await refetchFollowStatus();
      }
    }
  } catch (error: any) {
    console.error('Follow error:', error);
    toast.error(error.message || 'Failed to update follow status');
  } finally {
    followLoading.value = false;
  }
};

const navigateToEvent = (eventId: string) => {
  router.push(`/events/${eventId}`);
};

const handleAvatarError = (event: Event) => {
  const img = event.target as HTMLImageElement;
  img.src = '/images/default-avatar.jpg';
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return 'N/A';
  try {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'long',
      year: 'numeric'
    });
  } catch {
    return 'N/A';
  }
};

// Watch for userId changes
watch(userId, async (newId) => {
  if (newId) {
    await refetch();
    await refetchEvents();
    if (authStore.isAuthenticated) {
      await refetchFollowStatus();
    }
  }
}, { immediate: true });

onMounted(async () => {
  if (userId.value) {
    await refetch();
    await refetchEvents();
  }
});
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