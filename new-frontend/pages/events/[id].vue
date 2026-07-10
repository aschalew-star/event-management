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
          <!-- Hero Image Carousel -->
          <div class="relative rounded-2xl overflow-hidden bg-gray-200 dark:bg-gray-700 shadow-xl">
            <!-- Main Image -->
            <img
              :src="currentImage || '/images/event-placeholder.jpg'"
              :alt="event.title"
              class="w-full h-64 sm:h-96 object-cover"
            />
            
            <!-- Thumbnail Strip for Multiple Images -->
            <div v-if="eventImages.length > 1" class="absolute bottom-0 left-0 right-0 p-3 bg-gradient-to-t from-black/60 to-transparent">
              <div class="flex gap-2 overflow-x-auto pb-2 scrollbar-thin scrollbar-thumb-gray-400 scrollbar-track-transparent">
                <div
                  v-for="(img, index) in eventImages"
                  :key="index"
                  class="flex-shrink-0 w-16 h-16 rounded-lg overflow-hidden cursor-pointer border-2 transition-all relative"
                  :class="index === currentImageIndex ? 'border-indigo-500' : 'border-transparent hover:border-gray-300'"
                  @click="currentImageIndex = index"
                >
                  <img
                    :src="img.image_url"
                    :alt="`${event.title} - Image ${index + 1}`"
                    class="w-full h-full object-cover"
                  />
                  <div v-if="img.is_featured" class="absolute bottom-0 left-0 right-0 bg-indigo-600 text-white text-[8px] font-bold px-1 py-0.5 text-center">
                    FEATURED
                  </div>
                </div>
              </div>
            </div>

            <div class="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent"></div>
            
            <!-- Image Navigation Arrows -->
            <button
              v-if="eventImages.length > 1"
              @click="prevImage"
              class="absolute left-3 top-1/2 -translate-y-1/2 p-2 rounded-full bg-black/40 backdrop-blur-sm text-white hover:bg-black/60 transition-colors"
            >
              <i class="fas fa-chevron-left"></i>
            </button>
            <button
              v-if="eventImages.length > 1"
              @click="nextImage"
              class="absolute right-3 top-1/2 -translate-y-1/2 p-2 rounded-full bg-black/40 backdrop-blur-sm text-white hover:bg-black/60 transition-colors"
            >
              <i class="fas fa-chevron-right"></i>
            </button>

            <!-- Image Counter -->
            <div v-if="eventImages.length > 1" class="absolute bottom-20 left-1/2 -translate-x-1/2 px-3 py-1 bg-black/50 backdrop-blur-sm rounded-full text-white text-xs font-medium">
              {{ currentImageIndex + 1 }} / {{ eventImages.length }}
            </div>

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
                v-if="authStore.isAuthenticated"
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
            
            <!-- Organizer Info with Follow Button -->
            <div v-if="event.user" class="flex items-center gap-3 mb-4 pb-4 border-b border-gray-200 dark:border-gray-700">
              <div class="w-10 h-10 rounded-full bg-gradient-to-r from-indigo-500 to-purple-500 flex items-center justify-center text-white font-bold">
                {{ event.user.name?.charAt(0) || 'U' }}
              </div>
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  Organized by {{ event.user.name }}
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                  {{ formatDate(event.created_at) }}
                </p>
              </div>
              <button
                v-if="authStore.isAuthenticated && authStore.user?.id !== event.user.id"
                @click="toggleFollow"
                :disabled="followLoading"
                class="py-1.5 px-4 rounded-lg text-sm font-medium transition-all flex items-center justify-center gap-2"
                :class="isFollowing 
                  ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-200' 
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200'"
              >
                <i :class="[isFollowing ? 'fas fa-check-circle' : 'fas fa-plus-circle']"></i>
                {{ isFollowing ? 'Following' : 'Follow' }}
              </button>
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
                <p class="text-lg font-bold text-gray-900 dark:text-white">{{ attendeesCount }}</p>
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
              :disabled="ticketLoading || hasTicket || !authStore.isAuthenticated"
              class="w-full py-3.5 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 disabled:from-gray-400 disabled:to-gray-400 text-white font-semibold rounded-lg shadow-lg transition-all transform hover:scale-[1.02] disabled:scale-100 disabled:cursor-not-allowed text-base"
            >
              <i v-if="ticketLoading" class="fas fa-spinner fa-spin mr-2"></i>
              <i v-else-if="hasTicket" class="fas fa-check-circle mr-2"></i>
              <template v-else-if="!authStore.isAuthenticated">
                <i class="fas fa-sign-in-alt mr-2"></i>
                Sign in to Get Ticket
              </template>
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
            <div class="mt-4">
              <button
                v-if="authStore.isAuthenticated"
                @click="toggleBookmark"
                :disabled="bookmarkLoading"
                class="w-full py-2 px-3 rounded-lg text-sm font-medium transition-all flex items-center justify-center gap-2"
                :class="isBookmarked 
                  ? 'bg-amber-100 dark:bg-amber-900/30 text-amber-600 dark:text-amber-400 hover:bg-amber-200' 
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200'"
              >
                <i :class="[isBookmarked ? 'fas fa-bookmark' : 'far fa-bookmark']"></i>
                {{ isBookmarked ? 'Saved' : 'Save Event' }}
              </button>
              <div v-else class="text-center text-sm text-gray-500 dark:text-gray-400 py-2">
                <i class="fas fa-lock mr-1"></i> Sign in to save events
              </div>
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
            <div v-if="attendeesList.length > 0" class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
              <div class="flex items-center gap-2">
                <div class="flex -space-x-2">
                  <div
                    v-for="attendee in attendeesList.slice(0, 6)"
                    :key="attendee.id"
                    class="w-8 h-8 rounded-full bg-gradient-to-r from-indigo-500 to-purple-500 flex items-center justify-center text-white text-xs font-bold ring-2 ring-white dark:ring-gray-800"
                    :title="attendee.user?.name"
                  >
                    {{ attendee.user?.name?.charAt(0) || 'U' }}
                  </div>
                  <div
                    v-if="attendeesList.length > 6"
                    class="w-8 h-8 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center text-xs font-medium text-gray-600 dark:text-gray-400 ring-2 ring-white dark:ring-gray-800"
                  >
                    +{{ attendeesList.length - 6 }}
                  </div>
                </div>
                <span class="text-xs text-gray-500 dark:text-gray-400">
                  {{ attendeesList.length }} people going
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
                  :src="getRelatedEventImage(related.id) || '/images/event-placeholder.jpg'"
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
import { ref, computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useQuery, useMutation } from '@vue/apollo-composable';
import { useAuthStore } from '~/stores/auth';
import { useToast } from 'vue-toastification';
import EventMapSingle from '~/components/EventMapSingle.vue';

// Import queries from eventQueries
import {
  GET_EVENT_BY_ID,
  GET_RELATED_EVENTS,
  GET_EVENT_ATTENDEES,
  CHECK_EVENT_FOLLOW,
  CHECK_EVENT_BOOKMARK,
  GET_EVENT_STATS,
  GET_USER_TICKET,
  GET_EVENT_IMAGES,
  GET_EVENT_IMAGES_BY_EVENT_IDS,
} from '~/graphql/eventQueries';

// Import mutations from eventMutations
import {
  FOLLOW_EVENT,
  UNFOLLOW_EVENT,
  BOOKMARK_EVENT,
  UNBOOKMARK_EVENT,
  CREATE_TICKET,
} from '~/graphql/eventMutations';

// REMOVED auth middleware - anyone can view event details
// definePageMeta({
//   middleware: 'auth'
// });

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

// Image carousel state
const currentImageIndex = ref(0);

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

// console.log(authStore.isAuthenticated && authStore.user?.id !== event.user.id)
// 1. Fetch Event Details
const {
  result: eventResult,
  loading: eventLoading,
  error: eventError,
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

console.log("events from du id ",event)

// 2. Fetch Event Images
const {
  result: imagesResult,
  loading: imagesLoading,
} = useQuery(
  GET_EVENT_IMAGES,
  () => ({ eventId: eventId.value }),
  {
    fetchPolicy: 'cache-and-network',
    skip: () => !isValidEventId.value,
  }
);

const eventImages = computed(() => {
  return imagesResult.value?.event_images || [];
});

const currentImage = computed(() => {
  const images = eventImages.value;
  if (images.length > 0 && currentImageIndex.value < images.length) {
    return images[currentImageIndex.value]?.image_url;
  }
  return null;
});

// Image carousel controls
const nextImage = () => {
  if (eventImages.value.length > 0) {
    currentImageIndex.value = (currentImageIndex.value + 1) % eventImages.value.length;
  }
};

const prevImage = () => {
  if (eventImages.value.length > 0) {
    currentImageIndex.value = currentImageIndex.value === 0 
      ? eventImages.value.length - 1 
      : currentImageIndex.value - 1;
  }
};

// 3. Fetch Related Events
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

// 4. Fetch images for related events
const relatedEventIds = computed(() => {
  return relatedEvents.value.map((e: any) => e.id);
});

const { result: relatedImagesResult } = useQuery(
  GET_EVENT_IMAGES_BY_EVENT_IDS,
  () => ({ eventIds: relatedEventIds.value }),
  {
    fetchPolicy: 'cache-and-network',
    prefetch: false,
    skip: () => relatedEventIds.value.length === 0,
  }
);

// Build image map for related events
const eventImageMap = ref<Record<string, string>>({});

watch(relatedImagesResult, (result) => {
  if (result?.event_images) {
    const map: Record<string, string> = {};
    result.event_images.forEach((img: any) => {
      if (!map[img.event_id]) {
        map[img.event_id] = img.image_url;
      }
    });
    eventImageMap.value = map;
  }
}, { immediate: true });

const getRelatedEventImage = (eventId: string): string | null => {
  return eventImageMap.value[eventId] || null;
};

// 5. Fetch Attendees (people who have tickets)
const { result: attendeesResult } = useQuery(
  GET_EVENT_ATTENDEES,
  () => ({ eventId: eventId.value }),
  {
    fetchPolicy: 'cache-and-network',
    skip: () => !isValidEventId.value,
  }
);

const attendeesList = computed(() => attendeesResult.value?.tickets || []);
const attendeesCount = computed(() => attendeesList.value.length);

// 6. Check Follow status - Following the event organizer
const { result: followResult, refetch: refetchFollow } = useQuery(
  CHECK_EVENT_FOLLOW,
  () => {
    const userId = getCurrentUserId();
    const organizerId = event.value?.user?.id;
    if (!userId || !organizerId) return null;
    return {
      userId: userId,
      followedUserId: organizerId
    };
  },
  {
    fetchPolicy: 'network-only',
    prefetch: false,
    skip: () => !authStore.isAuthenticated || !getCurrentUserId() || !isValidEventId.value || !event.value?.user?.id,
  }
);

watch(followResult, (result) => {
  if (result?.follows && result.follows.length > 0) {
    const userId = getCurrentUserId();
    isFollowing.value = result.follows.some((follow: any) => 
      follow.follower_id === userId && 
      follow.followed_user_id === event.value?.user?.id
    );
  } else {
    isFollowing.value = false;
  }
}, { immediate: true });

// 7. Check Bookmark status
const { result: bookmarkResult, refetch: refetchBookmark } = useQuery(
  CHECK_EVENT_BOOKMARK,
  () => {
    const userId = getCurrentUserId();
    if (!userId) return null;
    return {
      eventId: eventId.value,
      userId: userId,
    };
  },
  {
    fetchPolicy: 'network-only',
    prefetch: false,
    skip: () => !authStore.isAuthenticated || !getCurrentUserId() || !isValidEventId.value,
  }
);

watch(bookmarkResult, (result) => {
  isBookmarked.value = result?.bookmarks && result.bookmarks.length > 0;
}, { immediate: true });

// 8. Get Event Stats
const { result: statsResult } = useQuery(
  GET_EVENT_STATS,
  () => ({ 
    eventId: eventId.value
  }),
  {
    fetchPolicy: 'cache-and-network',
    skip: () => !isValidEventId.value,
  }
);

const followersCount = computed(() => {
  if (statsResult.value?.follows_aggregate?.aggregate?.count) {
    return statsResult.value.follows_aggregate.aggregate.count;
  }
  return 0;
});

const bookmarksCount = computed(() => {
  if (statsResult.value?.bookmarks_aggregate?.aggregate?.count) {
    return statsResult.value.bookmarks_aggregate.aggregate.count;
  }
  return 0;
});

// 9. Check User Ticket
const { result: ticketResult, refetch: refetchTicket } = useQuery(
  GET_USER_TICKET,
  () => {
    const userId = getCurrentUserId();
    if (!userId) return null;
    return {
      userId: userId,
      eventId: eventId.value,
    };
  },
  {
    fetchPolicy: 'network-only',
    prefetch: false,
    skip: () => !authStore.isAuthenticated || !getCurrentUserId() || !isValidEventId.value,
  }
);

watch(ticketResult, (result) => {
  if (result?.tickets) {
    hasTicket.value = result.tickets.some((t: any) => t.event_id === eventId.value && t.status === 'confirmed');
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

// --- Action Handlers ---

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

      const result = await createTicketMutation({
        eventId: targetEventId,
        // userId: userId,
        quantity: 1,
        totalPrice: 0,
        status: 'confirmed',
      });

      if (result?.data?.insert_tickets_one) {
        await refetchTicket();
        toast.success('🎉 Free ticket claimed successfully!');
      }
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

const toggleFollow = async () => {
  if (!authStore.isAuthenticated) {
    toast.warning('Please sign in to follow organizers');
    router.push('/login');
    return;
  }

  const userId = getCurrentUserId();
  if (!userId) {
    toast.error('User ID not found. Please log in again.');
    return;
  }

  const organizerId = event.value?.user?.id;
  if (!organizerId) {
    toast.error('Event organizer not found.');
    return;
  }

  // Don't allow following yourself
  if (userId === organizerId) {
    toast.info('You cannot follow yourself');
    return;
  }

  followLoading.value = true;
  try {
    if (isFollowing.value) {
      // Unfollow the organizer - need both followedUserId and followerId
      const result = await unfollowMutation({
        followedUserId: organizerId,
        // followerId: userId
      });
      
      if (result?.data?.delete_follows?.affected_rows > 0) {
        isFollowing.value = false;
        toast.success('Unfollowed organizer');
        await refetchFollow();
      }
    } else {
      // Follow the organizer
      const result = await followMutation({
        followedUserId: organizerId,
        // followerId: userId
      });
      
      if (result?.data?.insert_follows_one) {
        isFollowing.value = true;
        toast.success(`Following ${event.value?.user?.name || 'organizer'}!`);
        await refetchFollow();
      }
    }
  } catch (error: any) {
    console.error('Follow request error details:', error);
    toast.error(error.message || 'Failed to update follow status');
  } finally {
    followLoading.value = false;
  }
};

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
      // Unbookmark - need both eventId and userId
      const result = await unbookmarkMutation({
        eventId: targetEventId,
        userId: userId
      });
      
      if (result?.data?.delete_bookmarks?.affected_rows > 0) {
        isBookmarked.value = false;
        toast.success('Bookmark removed');
        await refetchBookmark();
      } else {
        toast.warning('Bookmark not found');
      }
    } else {
      // Bookmark
      const result = await bookmarkMutation({
        eventId: targetEventId,
        // userId: userId
      });
      
      if (result?.data?.insert_bookmarks_one) {
        isBookmarked.value = true;
        toast.success('Event bookmarked!');
        await refetchBookmark();
      }
    }
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

const loading = computed(() => eventLoading.value || imagesLoading.value);
const error = computed(() => eventError.value);
</script>

<style scoped>
/* Custom scrollbar for thumbnail strip */
.scrollbar-thin::-webkit-scrollbar {
  height: 4px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.5);
  border-radius: 9999px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background: rgba(156, 163, 175, 0.8);
}

/* Smooth image transitions */
img {
  transition: opacity 0.3s ease;
}
</style>