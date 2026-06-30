<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800">
    <!-- Hero Section with Animated Background -->
    <div class="relative overflow-hidden bg-gradient-to-r from-indigo-600 via-purple-600 to-pink-600">
      <div class="absolute inset-0 opacity-20">
        <div class="absolute top-0 left-0 w-64 h-64 bg-white rounded-full filter blur-3xl animate-float"></div>
        <div class="absolute bottom-0 right-0 w-96 h-96 bg-white rounded-full filter blur-3xl animate-float-delayed"></div>
      </div>
      
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16 md:py-24">
        <div class="text-center text-white">
          <div class="flex items-center justify-center gap-3 mb-4">
            <i class="fas fa-calendar-alt text-4xl md:text-5xl animate-pulse"></i>
          </div>
          <h1 class="text-4xl md:text-6xl font-bold mb-4 bg-clip-text text-transparent bg-gradient-to-r from-white to-gray-200">
            Discover Events
          </h1>
          <p class="text-lg md:text-xl text-indigo-100 max-w-2xl mx-auto">
            Find and join amazing events happening around you
          </p>
          
          <!-- Quick Stats -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-8 max-w-3xl mx-auto">
            <div class="bg-white/10 backdrop-blur-sm rounded-xl p-4 border border-white/20">
              <div class="text-2xl font-bold">{{ stats.totalEvents }}</div>
              <div class="text-sm text-indigo-200">Total Events</div>
            </div>
            <div class="bg-white/10 backdrop-blur-sm rounded-xl p-4 border border-white/20">
              <div class="text-2xl font-bold">{{ stats.upcomingEvents }}</div>
              <div class="text-sm text-indigo-200">Upcoming</div>
            </div>
            <div class="bg-white/10 backdrop-blur-sm rounded-xl p-4 border border-white/20">
              <div class="text-2xl font-bold">{{ stats.categories }}</div>
              <div class="text-sm text-indigo-200">Categories</div>
            </div>
            <div class="bg-white/10 backdrop-blur-sm rounded-xl p-4 border border-white/20">
              <div class="text-2xl font-bold">{{ stats.attendees }}+</div>
              <div class="text-sm text-indigo-200">Attendees</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Search and Filter Section -->
    <div class="sticky top-0 z-40 bg-white/80 dark:bg-gray-800/80 backdrop-blur-md border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex flex-col lg:flex-row lg:items-center gap-4">
          <!-- Search Bar -->
          <div class="flex-1 relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search events by title, venue, or location..."
              class="w-full pl-12 pr-4 py-3 bg-gray-100 dark:bg-gray-700 border-2 border-transparent rounded-xl focus:border-indigo-500 focus:bg-white dark:focus:bg-gray-600 transition-all duration-300 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400"
              @input="debounceSearch"
            />
            <i class="fas fa-search absolute left-4 top-1/2 -translate-y-1/2 text-gray-400 dark:text-gray-500"></i>
          </div>

          <!-- Filter Buttons -->
          <div class="flex flex-wrap gap-2">
            <button
              @click="showFilters = !showFilters"
              class="flex items-center gap-2 px-4 py-2.5 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-xl text-gray-700 dark:text-gray-200 transition-colors"
            >
              <i class="fas fa-sliders-h"></i>
              <span class="hidden sm:inline">Filters</span>
              <span class="ml-1 text-xs bg-indigo-500 text-white rounded-full px-2 py-0.5">{{ activeFilters }}</span>
            </button>
            
            <button
              @click="toggleView"
              class="px-4 py-2.5 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-xl text-gray-700 dark:text-gray-200 transition-colors"
            >
              <i :class="[viewMode === 'grid' ? 'fas fa-list' : 'fas fa-th']"></i>
            </button>
          </div>
        </div>

        <!-- Expanded Filters -->
        <div v-if="showFilters" class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-600">
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
            <!-- Category Filter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Category</label>
              <select
                v-model="filters.category_id"
                class="w-full bg-gray-100 dark:bg-gray-700 border-2 border-transparent rounded-lg px-3 py-2 focus:border-indigo-500 focus:bg-white dark:focus:bg-gray-600 transition-colors text-gray-900 dark:text-white"
              >
                <option value="">All Categories</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                  {{ cat.name }}
                </option>
              </select>
            </div>

            <!-- Status Filter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Status</label>
              <select
                v-model="filters.status"
                class="w-full bg-gray-100 dark:bg-gray-700 border-2 border-transparent rounded-lg px-3 py-2 focus:border-indigo-500 focus:bg-white dark:focus:bg-gray-600 transition-colors text-gray-900 dark:text-white"
              >
                <option value="">All Status</option>
                <option value="published">Published</option>
                <option value="upcoming">Upcoming</option>
                <option value="completed">Completed</option>
                <option value="cancelled">Cancelled</option>
              </select>
            </div>

            <!-- Price Filter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Price</label>
              <select
                v-model="filters.price"
                class="w-full bg-gray-100 dark:bg-gray-700 border-2 border-transparent rounded-lg px-3 py-2 focus:border-indigo-500 focus:bg-white dark:focus:bg-gray-600 transition-colors text-gray-900 dark:text-white"
              >
                <option value="all">All Events</option>
                <option value="free">Free Only</option>
                <option value="paid">Paid Only</option>
              </select>
            </div>

            <!-- Sort By -->
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Sort By</label>
              <select
                v-model="filters.sortBy"
                class="w-full bg-gray-100 dark:bg-gray-700 border-2 border-transparent rounded-lg px-3 py-2 focus:border-indigo-500 focus:bg-white dark:focus:bg-gray-600 transition-colors text-gray-900 dark:text-white"
              >
                <option value="date_asc">Date (Earliest)</option>
                <option value="date_desc">Date (Latest)</option>
                <option value="title">Title</option>
                <option value="created_desc">Newest</option>
              </select>
            </div>
          </div>

          <!-- Apply Filters Button -->
          <div class="mt-4 flex justify-end gap-2">
            <button
              @click="resetFilters"
              class="px-4 py-2 text-gray-600 dark:text-gray-400 hover:text-gray-800 dark:hover:text-gray-200 transition-colors"
            >
              Reset All
            </button>
            <button
              @click="applyFilters"
              class="px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
            >
              Apply Filters
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-16">
        <div class="text-center">
          <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
          <p class="mt-4 text-gray-600 dark:text-gray-400">Loading events...</p>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-12">
        <div class="text-6xl mb-4">😅</div>
        <h3 class="text-xl font-semibold text-gray-800 dark:text-gray-200 mb-2">Oops! Something went wrong</h3>
        <p class="text-gray-600 dark:text-gray-400">{{ error }}</p>
        <button
          @click="fetchEvents"
          class="mt-4 px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
        >
          Try Again
        </button>
      </div>

      <!-- Events Grid -->
      <div v-else>
        <!-- Results Count -->
        <div class="flex items-center justify-between mb-6">
          <div class="text-sm text-gray-600 dark:text-gray-400">
            Showing {{ events.length }} events
            <span v-if="totalEvents > events.length" class="text-gray-500 dark:text-gray-500">
              ({{ totalEvents }} total)
            </span>
          </div>
          <button
            v-if="filtersApplied"
            @click="clearAllFilters"
            class="text-sm text-indigo-600 dark:text-indigo-400 hover:text-indigo-800 dark:hover:text-indigo-300"
          >
            <i class="fas fa-times-circle mr-1"></i> Clear all filters
          </button>
        </div>

        <!-- Grid/List View -->
        <div :class="[viewMode === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6' : 'space-y-4']">
          <EventCard
            v-for="event in events"
            :key="event.id"
            :event="event"
            :view-mode="viewMode"
            @click="navigateToEvent(event.id)"
          />
        </div>

        <!-- Empty State -->
        <div v-if="events.length === 0" class="text-center py-16">
          <div class="text-6xl mb-4">🎯</div>
          <h3 class="text-xl font-semibold text-gray-800 dark:text-gray-200 mb-2">No events found</h3>
          <p class="text-gray-600 dark:text-gray-400">Try adjusting your filters or search terms</p>
          <button
            @click="resetFilters"
            class="mt-4 px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
          >
            Reset Filters
          </button>
        </div>

        <!-- Load More -->
        <div v-if="hasMore" class="mt-8 text-center">
          <button
            @click="loadMore"
            :disabled="loadingMore"
            class="px-6 py-3 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg text-gray-700 dark:text-gray-200 transition-colors disabled:opacity-50"
          >
            <i v-if="loadingMore" class="fas fa-spinner fa-spin mr-2"></i>
            <span v-else>Load More Events</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Quick Action FAB -->
    <NuxtLink
      to="/events/create"
      class="fixed bottom-6 right-6 z-50 p-4 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white rounded-full shadow-lg hover:shadow-xl transition-all duration-300 transform hover:scale-110"
    >
      <i class="fas fa-plus text-xl"></i>
    </NuxtLink>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue';
import { useQuery } from '@vue/apollo-composable';
import { GET_ALL_EVENTS, GET_EVENT_CATEGORIES } from '~/graphql/eventQueries';
import { useDebounce } from '~/composables/useDebounce';
import { useRouter } from 'vue-router';
import EventCard from '~/components/events/EventCard.vue';

const router = useRouter();
const searchQuery = ref('');
const showFilters = ref(false);
const viewMode = ref<'grid' | 'list'>('grid');
const currentPage = ref(1);
const pageSize = ref(12);
const loadingMore = ref(false);

// Filter State
const filters = reactive({
  category_id: '',
  status: '',
  price: 'all',
  sortBy: 'date_asc',
});

// Stats
const stats = reactive({
  totalEvents: 0,
  upcomingEvents: 0,
  categories: 0,
  attendees: 2543,
});

// Build GraphQL Where Clause
const buildWhereClause = () => {
  const where: any = {};
  
  // Status filter
  if (filters.status) {
    where.status = { _eq: filters.status };
  } else {
    // Default: show published and upcoming
    where._or = [
      { status: { _eq: 'published' } },
      { status: { _eq: 'upcoming' } }
    ];
  }

  // Category filter
  if (filters.category_id) {
    where.category_id = { _eq: filters.category_id };
  }

  // Price filter
  if (filters.price === 'free') {
    where.is_free = { _eq: true };
  } else if (filters.price === 'paid') {
    where.is_free = { _eq: false };
  }

  // Search query
  if (searchQuery.value) {
    where._or = [
      { title: { _ilike: `%${searchQuery.value}%` } },
      { description: { _ilike: `%${searchQuery.value}%` } },
      { venue: { _ilike: `%${searchQuery.value}%` } },
      { address: { _ilike: `%${searchQuery.value}%` } },
    ];
  }

  return where;
};

// Build Order By
const buildOrderBy = () => {
  switch (filters.sortBy) {
    case 'date_asc':
      return { event_date: 'asc' };
    case 'date_desc':
      return { event_date: 'desc' };
    case 'title':
      return { title: 'asc' };
    case 'created_desc':
      return { created_at: 'desc' };
    default:
      return { event_date: 'asc' };
  }
};

// Fetch Events
const { result, loading, error, fetchMore, refetch } = useQuery(
  GET_ALL_EVENTS,
  () => ({
    filters: buildWhereClause(),
    limit: pageSize.value,
    offset: (currentPage.value - 1) * pageSize.value,
  }),
  {
    fetchPolicy: 'cache-and-network',
    context: {
      headers: {
        'x-hasura-role': 'anonymous',
      },
    },
  }
);

// Fetch Categories
const { result: categoriesResult } = useQuery(GET_EVENT_CATEGORIES);

const events = computed(() => result.value?.events || []);
const totalEvents = computed(() => result.value?.events_aggregate?.aggregate?.count || 0);
const categories = computed(() => categoriesResult.value?.categories || []);
const hasMore = computed(() => events.value.length < totalEvents.value);
const activeFilters = computed(() => {
  let count = 0;
  if (filters.category_id) count++;
  if (filters.status) count++;
  if (filters.price !== 'all') count++;
  if (searchQuery.value) count++;
  return count;
});
const filtersApplied = computed(() => activeFilters.value > 0);

// Debounced Search
const debounceSearch = useDebounce(() => {
  refetch({
    filters: buildWhereClause(),
    limit: pageSize.value,
    offset: 0,
  });
}, 500);

// Methods
const applyFilters = () => {
  currentPage.value = 1;
  refetch({
    filters: buildWhereClause(),
    limit: pageSize.value,
    offset: 0,
  });
  showFilters.value = false;
};

const resetFilters = () => {
  filters.category_id = '';
  filters.status = '';
  filters.price = 'all';
  filters.sortBy = 'date_asc';
  searchQuery.value = '';
  currentPage.value = 1;
  applyFilters();
};

const clearAllFilters = resetFilters;

const loadMore = async () => {
  if (loadingMore.value || !hasMore.value) return;
  loadingMore.value = true;
  currentPage.value++;
  
  await fetchMore({
    variables: {
      offset: (currentPage.value - 1) * pageSize.value,
    },
    updateQuery: (prev, { fetchMoreResult }) => {
      if (!fetchMoreResult) return prev;
      return {
        ...fetchMoreResult,
        events: [...prev.events, ...fetchMoreResult.events],
      };
    },
  });
  loadingMore.value = false;
};

const toggleView = () => {
  viewMode.value = viewMode.value === 'grid' ? 'list' : 'grid';
  // Save preference
  localStorage.setItem('eventViewMode', viewMode.value);
};

const navigateToEvent = (id: string) => {
  router.push(`/events/${id}`);
};

const fetchEvents = () => {
  refetch();
};

// Watch for filter changes
watch(
  () => ({ ...filters }),
  () => {
    if (!showFilters.value) {
      applyFilters();
    }
  },
  { deep: true }
);

// Lifecycle
onMounted(() => {
  // Load view preference
  const savedView = localStorage.getItem('eventViewMode');
  if (savedView === 'grid' || savedView === 'list') {
    viewMode.value = savedView;
  }
  
  // Update stats
  stats.categories = categories.value.length;
});
</script>

<style scoped>
@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(5deg); }
}

@keyframes float-delayed {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-30px) rotate(-5deg); }
}

.animate-float {
  animation: float 6s ease-in-out infinite;
}

.animate-float-delayed {
  animation: float-delayed 8s ease-in-out infinite;
}
</style>