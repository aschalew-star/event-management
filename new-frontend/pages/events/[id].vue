<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center min-h-screen">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-16 w-16 border-4 border-indigo-500 border-t-transparent"></div>
        <p class="mt-4 text-gray-600 dark:text-gray-400 text-lg">Loading event details...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex items-center justify-center min-h-screen px-4">
      <div class="text-center max-w-md">
        <div class="text-6xl mb-4">😅</div>
        <h3 class="text-2xl font-bold text-gray-800 dark:text-gray-200 mb-2">Event not found</h3>
        <p class="text-gray-600 dark:text-gray-400 mb-6">{{ error.message || 'The event you\'re looking for doesn\'t exist or has been removed.' }}</p>
        <NuxtLink
          to="/events"
          class="inline-flex items-center gap-2 px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
        >
          <i class="fas fa-arrow-left"></i>
          Back to Events
        </NuxtLink>
      </div>
    </div>

    <!-- Event Details -->
    <div v-else-if="event" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Back Button -->
      <button
        @click="router.back()"
        class="inline-flex items-center gap-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 mb-6 transition-colors group"
      >
        <i class="fas fa-arrow-left group-hover:-translate-x-1 transition-transform"></i>
        Back to Events
      </button>

      <!-- Main Content -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column: Event Details -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Hero Image -->
          <div class="relative rounded-2xl overflow-hidden bg-gray-200 dark:bg-gray-700 shadow-xl">
            <img
              :src="event.featured_image || '/images/event-placeholder.jpg'"
              :alt="event.title"
              class="w-full h-64 sm:h-96 object-cover"
            />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>
            
            <!-- Badges -->
            <div class="absolute top-4 left-4 flex flex-wrap gap-2">
              <span
                v-if="event.status"
                class="px-3 py-1.5 text-xs font-medium rounded-full backdrop-blur-sm"
                :class="statusColors[event.status] || 'bg-gray-500/80 text-white'"
              >
                <i class="fas fa-circle text-[6px] mr-1.5 align-middle"></i>
                {{ event.status.charAt(0).toUpperCase() + event.status.slice(1) }}
              </span>
              <span
                v-if="event.is_free"
                class="px-3 py-1.5 text-xs font-medium bg-green-500/80 backdrop-blur-sm text-white rounded-full"
              >
                <i class="fas fa-gift mr-1.5"></i>
                Free Entry
              </span>
              <span
                v-if="!event.is_free && event.price"
                class="px-3 py-1.5 text-xs font-medium bg-amber-500/80 backdrop-blur-sm text-white rounded-full"
              >
                <i class="fas fa-tag mr-1.5"></i>
                ${{ Number(event.price).toFixed(2) }}
              </span>
            </div>

            <!-- Category Tag -->
            <div class="absolute bottom-4 left-4">
              <span
                class="px-3 py-1.5 bg-black/50 backdrop-blur-sm rounded-full text-white text-sm font-medium flex items-center gap-2"
              >
                <i :class="[event.category?.icon || 'fas fa-tag']"></i>
                {{ event.category?.name || 'Uncategorized' }}
              </span>
            </div>

            <!-- Action Buttons Overlay -->
            <div class="absolute top-4 right-4 flex gap-2">
              <button
                @click="toggleFollow"
                :disabled="followLoading"
                class="p-2.5 rounded-full backdrop-blur-sm transition-all hover:scale-110"
                :class="isFollowing ? 'bg-indigo-500/90 text-white' : 'bg-white/90 text-gray-700 hover:bg-indigo-100'"
                :title="isFollowing ? 'Unfollow Event' : 'Follow Event'"
              >
                <i :class="[isFollowing ? 'fas fa-bell': 'far fa-bell']"></i>
              </button>
              <button
                @click="toggleBookmark"
                :disabled="bookmarkLoading"
                class="p-2.5 rounded-full backdrop-blur-sm transition-all hover:scale-110"
                :class="isBookmarked ? 'bg-amber-500/90 text-white' : 'bg-white/90 text-gray-700 hover:bg-amber-100'"
                :title="isBookmarked ? 'Remove Bookmark' : 'Add Bookmark'"
              >
                <i :class="[isBookmarked ? 'fas fa-bookmark' : 'far fa-bookmark']"></i>
              </button>
            </div>
          </div>

          <!-- Title & Description -->
          <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
            <div class="flex items-start justify-between gap-4">
              <h1 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white flex-1">
                {{ event.title }}
              </h1>
              <div class="flex gap-2 flex-shrink-0">
                <span class="px-3 py-1 bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 rounded-full text-xs font-medium flex items-center gap-1">
                  <i class="fas fa-eye"></i>
                  {{ event.view_count || 0 }}
                </span>
              </div>
            </div>
            
            <!-- Organizer Info -->
            <div v-if="event.user" class="flex items-center gap-3 mb-4 pb-4 border-b border-gray-200 dark:border-gray-700">
              <div class="w-10 h-10 rounded-full bg-gradient-to-r from-indigo-500 to-purple-500 flex items-center justify-center text-white font-bold">
                {{ event.user.name?.charAt(0) || 'U' }}
              </div>
              <div>
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  Organized by {{ event.user.name }}
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  {{ formatDate(event.created_at) }}
                </p>
              </div>
            </div>

            <div class="prose prose-indigo dark:prose-invert max-w-none">
              <p class="text-gray-700 dark:text-gray-300 leading-relaxed whitespace-pre-wrap">
                {{ event.description }}
              </p>
            </div>
          </div>

          <!-- Event Details Grid -->
          <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              <i class="fas fa-info-circle text-indigo-500 mr-2"></i>
              Event Details
            </h3>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-xl">
                <i class="fas fa-calendar-day text-indigo-500 text-lg mt-0.5"></i>
                <div>
                  <p class="text-xs text-gray-500 dark:text-gray-400">Date</p>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ formatDate(event.event_date, true) }}
                  </p>
                </div>
              </div>
              <div v-if="event.start_time" class="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-xl">
                <i class="fas fa-clock text-indigo-500 text-lg mt-0.5"></i>
                <div>
                  <p class="text-xs text-gray-500 dark:text-gray-400">Time</p>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ event.start_time }}
                    <span v-if="event.end_time" class="text-gray-500 dark:text-gray-400">
                      - {{ event.end_time }}
                    </span>
                  </p>
                </div>
              </div>
              <div class="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-xl">
                <i class="fas fa-map-pin text-indigo-500 text-lg mt-0.5"></i>
                <div>
                  <p class="text-xs text-gray-500 dark:text-gray-400">Venue</p>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ event.venue || 'TBD' }}
                  </p>
                </div>
              </div>
              <div class="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-xl">
                <i class="fas fa-location-dot text-indigo-500 text-lg mt-0.5"></i>
                <div>
                  <p class="text-xs text-gray-500 dark:text-gray-400">Address</p>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">
                    {{ event.address || 'TBD' }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Location Map -->
          <div v-if="event.latitude && event.longitude" class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
              <i class="fas fa-map text-indigo-500 mr-2"></i>
              Location
            </h3>
            <div class="rounded-xl overflow-hidden">
              <EventMapSingle
                :latitude="event.latitude"
                :longitude="event.longitude"
                height="300px"
                mode="display"
              />
            </div>
            <div class="mt-3 text-sm text-gray-600 dark:text-gray-400">
              <i class="fas fa-map-marker-alt text-indigo-500 mr-2"></i>
              {{ event.address || `${event.latitude}, ${event.longitude}` }}
            </div>
          </div>
        </div>

        <!-- Right Column: Sidebar -->
        <div class="space-y-6">
          <!-- Action Card -->
          <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6 sticky top-24">
            <div class="text-center mb-4">
              <div class="text-3xl font-bold text-gray-900 dark:text-white">
                {{ event.is_free ? 'FREE' : `$${Number(event.price).toFixed(2)}` }}
              </div>
              <p class="text-sm text-gray-500 dark:text-gray-400">per person</p>
            </div>

            <!-- Stats Row -->
            <div class="grid grid-cols-3 gap-2 mb-4 p-3 bg-gray-50 dark:bg-gray-700/30 rounded-xl">
              <div class="text-center">
                <p class="text-lg font-bold text-gray-900 dark:text-white">{{ attendees.length }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400">Going</p>
              </div>
              <div class="text-center border-x border-gray-200 dark:border-gray-600">
                <p class="text-lg font-bold text-gray-900 dark:text-white">{{ followersCount }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400">Followers</p>
              </div>
              <div class="text-center">
                <p class="text-lg font-bold text-gray-900 dark:text-white">{{ bookmarksCount }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400">Bookmarks</p>
              </div>
            </div>

            <!-- Ticket Button -->
            <button
              @click="handleTicketAction"
              :disabled="ticketLoading || hasTicket"
              class="w-full py-3.5 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 disabled:from-gray-400 disabled:to-gray-400 text-white font-semibold rounded-lg shadow-lg transition-all transform hover:scale-[1.02] disabled:scale-100 disabled:cursor-not-allowed text-base"
            >
              <i v-if="ticketLoading" class="fas fa-spinner fa-spin mr-2"></i>
              <i v-else-if="hasTicket" class="fas fa-check-circle mr-2"></i>
              <template v-else-if="event.is_free">
                <i class="fas fa-ticket-alt mr-2"></i>
                Get Free Ticket
              </template>
              <template v-else>
                <i class="fas fa-credit-card mr-2"></i>
                Buy Ticket - ${{ Number(event.price).toFixed(2) }}
              </template>
              <span v-if="hasTicket">Ticket Confirmed ✓</span>
            </button>

            <!-- Quick Actions -->
            <div class="mt-4 grid grid-cols-2 gap-2">
              <button
                @click="toggleFollow"
                :disabled="followLoading"
                class="py-2 px-3 rounded-lg text-sm font-medium transition-all flex items-center justify-center gap-2"
                :class="isFollowing 
                  ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-200' 
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200'"
              >
                <i :class="[isFollowing ? 'fas fa-bell' : 'far fa-bell']"></i>
                {{ isFollowing ? 'Following' : 'Follow' }}
              </button>
              <button
                @click="toggleBookmark"
                :disabled="bookmarkLoading"
                class="py-2 px-3 rounded-lg text-sm font-medium transition-all flex items-center justify-center gap-2"
                :class="isBookmarked 
                  ? 'bg-amber-100 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400 hover:bg-amber-200' 
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200'"
              >
                <i :class="[isBookmarked ? 'fas fa-bookmark' : 'far fa-bookmark']"></i>
                {{ isBookmarked ? 'Saved' : 'Save' }}
              </button>
            </div>

            <div class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
              <div class="space-y-2 text-sm">
                <div class="flex items-center justify-between text-gray-600 dark:text-gray-400">
                  <span><i class="fas fa-share-alt mr-2"></i>Share</span>
                  <div class="flex gap-2">
                    <button @click="shareEvent('facebook')" class="text-blue-600 hover:text-blue-800 dark:text-blue-400 dark:hover:text-blue-300">
                      <i class="fab fa-facebook"></i>
                    </button>
                    <button @click="shareEvent('twitter')" class="text-sky-500 hover:text-sky-700 dark:text-sky-400 dark:hover:text-sky-300">
                      <i class="fab fa-twitter"></i>
                    </button>
                    <button @click="shareEvent('linkedin')" class="text-blue-700 hover:text-blue-900 dark:text-blue-500 dark:hover:text-blue-400">
                      <i class="fab fa-linkedin"></i>
                    </button>
                    <button @click="shareEvent('copy')" class="text-gray-600 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-300">
                      <i class="fas fa-link"></i>
                    </button>
                  </div>
                </div>
                <div class="flex items-center justify-between text-gray-600 dark:text-gray-400">
                  <span><i class="fas fa-calendar-plus mr-2"></i>Add to Calendar</span>
                  <button @click="addToCalendar" class="text-indigo-600 hover:text-indigo-800 dark:text-indigo-400 dark:hover:text-indigo-300">
                    <i class="fas fa-plus-circle"></i>
                  </button>
                </div>
              </div>
            </div>

            <!-- Attendees Avatars -->
            <div v-if="attendees.length > 0" class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
              <div class="flex items-center gap-2">
                <div class="flex -space-x-2">
                  <div
                    v-for="attendee in attendees.slice(0, 6)"
                    :key="attendee.id"
                    class="w-8 h-8 rounded-full bg-gradient-to-r from-indigo-500 to-purple-500 flex items-center justify-center text-white text-xs font-bold ring-2 ring-white dark:ring-gray-800"
                    :title="attendee.user?.name"
                  >
                    {{ attendee.user?.name?.charAt(0) || 'U' }}
                  </div>
                  <div
                    v-if="attendees.length > 6"
                    class="w-8 h-8 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center text-xs font-medium text-gray-600 dark:text-gray-400 ring-2 ring-white dark:ring-gray-800"
                  >
                    +{{ attendees.length - 6 }}
                  </div>
                </div>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  {{ attendees.length }} people going
                </span>
              </div>
            </div>
          </div>

          <!-- Related Events -->
          <div v-if="relatedEvents.length > 0" class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-6">
            <h4 class="text-sm font-semibold text-gray-900 dark:text-white uppercase tracking-wider mb-4">
              <i class="fas fa-arrow-right text-indigo-500 mr-2"></i>
              Similar Events
            </h4>
            <div class="space-y-3">
              <div
                v-for="related in relatedEvents.slice(0, 3)"
                :key="related.id"
                class="flex gap-3 p-2 hover:bg-gray-50 dark:hover:bg-gray-700 rounded-lg cursor-pointer transition-colors"
                @click="router.push(`/events/${related.id}`)"
              >
                <img
                  :src="related.featured_image || '/images/event-placeholder.jpg'"
                  :alt="related.title"
                  class="w-16 h-16 rounded-lg object-cover flex-shrink-0"
                />
                <div class="flex-1 min-w-0">
                  <p class="text-sm font-medium text-gray-900 dark:text-white truncate">
                    {{ related.title }}
                  </p>
                  <p class="text-xs text-gray-500 dark:text-gray-400 truncate">
                    {{ related.venue || 'TBD' }}
                  </p>
                  <p class="text-xs text-gray-500 dark:text-gray-400">
                    {{ formatDate(related.event_date) }}
                  </p>
                </div>
              </div>
            </div>
          </div>
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
import EventMapSingle from '~/components/EventMapSingle.vue';
import {
  GET_EVENT_BY_ID,
  GET_RELATED_EVENTS,
  GET_EVENT_ATTENDEES,
  CHECK_EVENT_FOLLOW,
  CHECK_EVENT_BOOKMARK,
  GET_EVENT_STATS,
  GET_USER_TICKET,
} from '~/graphql/eventQueries';
import {
  FOLLOW_EVENT,
  UNFOLLOW_EVENT,
  BOOKMARK_EVENT,
  UNBOOKMARK_EVENT,
  CREATE_TICKET,
} from '~/graphql/eventMutations';

definePageMeta({
  middleware: 'auth'
});

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

// Safely extract eventId from dynamic routes
const eventId = computed(() => {
  const id = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id;
  if (!id || id.length < 10) return null;
  return id;
});

const isValidEventId = computed(() => {
  return eventId.value !== null && eventId.value.length > 10;
});

// UI State loaders
const followLoading = ref(false);
const bookmarkLoading = ref(false);
const ticketLoading = ref(false);
const isFollowing = ref(false);
const isBookmarked = ref(false);
const hasTicket = ref(false);

// Helper function to get clean event ID
const getCleanEventId = (): string => {
  const rawId = Array.isArray(route.params.id) ? route.params.id[0] : route.params.id;
  return String(rawId || '').trim();
};

// Get current user ID helper
const getCurrentUserId = (): string | null => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    return authStore.user.id;
  }
  return null;
};

// 1. Fetch Event Details
const {
  result: eventResult,
  loading: eventLoading,
  error: eventError,
  refetch: refetchEvent,
} = useQuery(
  GET_EVENT_BY_ID,
  () => ({ id: eventId.value }),
  {
    fetchPolicy: 'cache-and-network',
    skip: () => !isValidEventId.value,
  }
);

const event = computed(() => eventResult.value?.events_by_pk);
const hasCategory = computed(() => !!event.value?.category_id);

// 2. Fetch Related Events
const { result: relatedResult } = useQuery(
  GET_RELATED_EVENTS,
  () => ({
    categoryId: event.value?.category_id || '00000000-0000-0000-0000-000000000000',
    eventId: eventId.value,
    limit: 4,
  }),
  {
    fetchPolicy: 'cache-and-network',
    prefetch: false,
    skip: () => !hasCategory.value || !isValidEventId.value,
  }
);

const relatedEvents = computed(() => relatedResult.value?.events || []);

// 3. Fetch Attendees
const { result: attendeesResult } = useQuery(
  GET_EVENT_ATTENDEES,
  () => ({ eventId: eventId.value }),
  {
    fetchPolicy: 'cache-and-network',
    skip: () => !isValidEventId.value,
  }
);

const attendees = computed(() => attendeesResult.value?.bookmarks || []);

// 4. Check Follow status
const { result: followResult, refetch: refetchFollow } = useQuery(
  CHECK_EVENT_FOLLOW,
  () => ({
    eventId: eventId.value,
    userId: getCurrentUserId(),
  }),
  {
    fetchPolicy: 'network-only',
    prefetch: false,
    skip: () => !authStore.isAuthenticated || !getCurrentUserId() || !isValidEventId.value,
  }
);

watch(followResult, (result) => {
  isFollowing.value = result?.follows && result.follows.length > 0;
}, { immediate: true });

// 5. Check Bookmark status
const { result: bookmarkResult, refetch: refetchBookmark } = useQuery(
  CHECK_EVENT_BOOKMARK,
  () => ({
    eventId: eventId.value,
    userId: getCurrentUserId(),
  }),
  {
    fetchPolicy: 'network-only',
    prefetch: false,
    skip: () => !authStore.isAuthenticated || !getCurrentUserId() || !isValidEventId.value,
  }
);

watch(bookmarkResult, (result) => {
  isBookmarked.value = result?.bookmarks && result.bookmarks.length > 0;
}, { immediate: true });

// 6. Get Event Stats
const { result: statsResult } = useQuery(
  GET_EVENT_STATS,
  () => ({ eventId: eventId.value }),
  {
    fetchPolicy: 'cache-and-network',
    skip: () => !isValidEventId.value,
  }
);

const followersCount = computed(() => {
  return statsResult.value?.follows_aggregate?.aggregate?.count || 0;
});

const bookmarksCount = computed(() => {
  return statsResult.value?.bookmarks_aggregate?.aggregate?.count || 0;
});

// 7. Check User Ticket
const { result: ticketResult, refetch: refetchTicket } = useQuery(
  GET_USER_TICKET,
  () => ({
    userId: getCurrentUserId(),
    eventId: eventId.value,
  }),
  {
    fetchPolicy: 'network-only',
    prefetch: false,
    skip: () => !authStore.isAuthenticated || !getCurrentUserId() || !isValidEventId.value,
  }
);

watch(ticketResult, (result) => {
  if (result?.tickets) {
    hasTicket.value = result.tickets.some((t: any) => t.event_id === eventId.value);
  }
}, { immediate: true });

// Mutations
const { mutate: followMutation } = useMutation(FOLLOW_EVENT);
const { mutate: unfollowMutation } = useMutation(UNFOLLOW_EVENT);
const { mutate: bookmarkMutation } = useMutation(BOOKMARK_EVENT);
const { mutate: unbookmarkMutation } = useMutation(UNBOOKMARK_EVENT);
const { mutate: createTicketMutation } = useMutation(CREATE_TICKET);

// Status Colors mapping
const statusColors: Record<string, string> = {
  published: 'bg-green-500/80 text-white',
  upcoming: 'bg-blue-500/80 text-white',
  completed: 'bg-gray-500/80 text-white',
  cancelled: 'bg-red-500/80 text-white',
  draft: 'bg-yellow-500/80 text-white',
};

const formatDate = (dateStr: string, full: boolean = false) => {
  if (!dateStr) return 'TBD';
  const date = new Date(dateStr);
  if (full) {
    return date.toLocaleDateString('en-US', {
      weekday: 'long',
      month: 'long',
      day: 'numeric',
      year: 'numeric',
    });
  }
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });
};

// --- FIXED Action Handlers ---

const handleTicketAction = async () => {
  if (!authStore.isAuthenticated) {
    toast.warning('Please sign in to get a ticket');
    router.push('/login');
    return;
  }

  const userId = getCurrentUserId();
  if (!userId) {
    toast.error('User ID not found. Please log in again.');
    return;
  }

  const targetEventId = getCleanEventId();
  if (!targetEventId || targetEventId.length < 10) {
    toast.error('Invalid event configuration context.');
    return;
  }

  if (hasTicket.value) {
    toast.info('You already have a ticket for this event!');
    return;
  }

  try {
    if (event.value?.is_free) {
      ticketLoading.value = true;

      await createTicketMutation({
        eventId: targetEventId,
        userId: userId,
        quantity: 1,
        totalPrice: 0,
        status: 'confirmed',
      });

      await refetchTicket();
      toast.success('🎉 Free ticket claimed successfully!');
    } else {
      router.push({
        path: `/payment/${targetEventId}`,
        query: {
          title: event.value?.title,
          price: event.value?.price,
          eventId: targetEventId,
        }
      });
    }
  } catch (error: any) {
    console.error('Ticket operation critical error:', error);
    toast.error(error.message || 'Failed to process ticket request.');
  } finally {
    ticketLoading.value = false;
  }
};

// FIXED: toggleFollow with user_id
const toggleFollow = async () => {
  if (!authStore.isAuthenticated) {
    toast.warning('Please sign in to follow events');
    router.push('/login');
    return;
  }

  const userId = getCurrentUserId();
  if (!userId) {
    toast.error('User ID not found. Please log in again.');
    return;
  }

  const targetEventId = getCleanEventId();
  if (!targetEventId || targetEventId.length < 10) {
    toast.error('Action aborted: Invalid identifier.');
    return;
  }

  followLoading.value = true;
  try {
    if (isFollowing.value) {
      await unfollowMutation({
        eventId: targetEventId
      });
      isFollowing.value = false;
      toast.success('Unfollowed event');
    } else {
      await followMutation({
        eventId: targetEventId,
        userId: userId
      });
      isFollowing.value = true;
      toast.success('Following event!');
    }
    await refetchFollow();
  } catch (error: any) {
    console.error('Follow request error details:', error);
    toast.error(error.message || 'Failed to update follow status');
  } finally {
    followLoading.value = false;
  }
};

// FIXED: toggleBookmark with user_id
const toggleBookmark = async () => {
  if (!authStore.isAuthenticated) {
    toast.warning('Please sign in to bookmark events');
    router.push('/login');
    return;
  }

  const userId = getCurrentUserId();
  if (!userId) {
    toast.error('User ID not found. Please log in again.');
    return;
  }

  const targetEventId = getCleanEventId();
  if (!targetEventId || targetEventId.length < 10) {
    toast.error('Action aborted: Invalid identifier.');
    return;
  }

  bookmarkLoading.value = true;
  try {
    if (isBookmarked.value) {
      const result = await unbookmarkMutation({
        eventId: targetEventId
      });
      console.log('Unbookmark result:', result);
      isBookmarked.value = false;
      toast.success('Bookmark removed');
    } else {
      const result = await bookmarkMutation({
        eventId: targetEventId,
        userId: userId
      });
      console.log('Bookmark result:', result);
      isBookmarked.value = true;
      toast.success('Event bookmarked!');
    }
    await refetchBookmark();
  } catch (error: any) {
    console.error('Bookmark request error details:', error);
    toast.error(error.message || 'Failed to update bookmark status');
  } finally {
    bookmarkLoading.value = false;
  }
};

const shareEvent = (platform: string) => {
  const url = window.location.href;
  const title = encodeURIComponent(event.value?.title || '');

  switch (platform) {
    case 'facebook':
      window.open(`https://www.facebook.com/sharer/sharer.php?u=${url}`, '_blank');
      break;
    case 'twitter':
      window.open(`https://twitter.com/intent/tweet?text=${title}&url=${url}`, '_blank');
      break;
    case 'linkedin':
      window.open(`https://www.linkedin.com/sharing/share-offsite/?url=${url}`, '_blank');
      break;
    case 'copy':
      navigator.clipboard.writeText(url);
      toast.success('Link copied to clipboard!');
      break;
  }
};

const addToCalendar = () => {
  if (!event.value) return;

  const startDate = new Date(event.value.event_date);
  if (event.value.start_time) {
    const [hours, minutes] = event.value.start_time.split(':');
    startDate.setHours(parseInt(hours), parseInt(minutes));
  }

  const endDate = new Date(startDate);
  if (event.value.end_time) {
    const [hours, minutes] = event.value.end_time.split(':');
    endDate.setHours(parseInt(hours), parseInt(minutes));
  } else {
    endDate.setHours(startDate.getHours() + 2);
  }

  const formatDateForCalendar = (date: Date) => {
    return date.toISOString().replace(/[-:]/g, '').split('.')[0];
  };

  const calendarUrl = `https://calendar.google.com/calendar/render?action=TEMPLATE&text=${encodeURIComponent(event.value.title)}&details=${encodeURIComponent(event.value.description)}&location=${encodeURIComponent(event.value.venue || event.value.address)}&dates=${formatDateForCalendar(startDate)}/${formatDateForCalendar(endDate)}`;

  window.open(calendarUrl, '_blank');
};

const loading = computed(() => eventLoading.value);
const error = computed(() => eventError.value);
</script>

<style scoped>
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 8px;
}

::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

.dark ::-webkit-scrollbar-thumb {
  background: #475569;
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: #64748b;
}
</style>