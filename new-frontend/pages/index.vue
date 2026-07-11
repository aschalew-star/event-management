<template>
  <div class="min-h-screen">
    <!-- Hero Section -->
    <section class="relative overflow-hidden bg-gradient-to-br from-blue-600 via-purple-600 to-pink-600">
      <div class="absolute inset-0">
        <div class="absolute inset-0 bg-[url('/images/hero-pattern.svg')] opacity-10"></div>
        <div class="absolute top-0 -left-4 w-72 h-72 bg-purple-300 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob"></div>
        <div class="absolute top-0 -right-4 w-72 h-72 bg-yellow-300 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-2000"></div>
        <div class="absolute -bottom-8 left-20 w-72 h-72 bg-pink-300 rounded-full mix-blend-multiply filter blur-xl opacity-70 animate-blob animation-delay-4000"></div>
      </div>
      
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-32">
        <div class="text-center">
          <div class="inline-flex items-center gap-2 px-4 py-2 bg-white/20 backdrop-blur-sm rounded-full text-white text-sm mb-6">
            <Icon name="lucide:sparkles" class="w-4 h-4" />
            <span>Discover Amazing Events</span>
          </div>
          
          <h1 class="text-4xl md:text-6xl font-bold text-white mb-6 leading-tight">
            Find Your Next
            <span class="bg-gradient-to-r from-yellow-200 to-pink-200 bg-clip-text text-transparent">Adventure</span>
          </h1>
          
          <p class="text-xl text-white/90 max-w-2xl mx-auto mb-10">
            Discover, connect, and experience the best events in your area. 
            From concerts to workshops, find what makes you feel alive.
          </p>
          
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <NuxtLink
              to="/events"
              class="px-8 py-4 bg-white text-gray-900 font-semibold rounded-xl hover:shadow-2xl transition-all transform hover:-translate-y-1 flex items-center justify-center gap-2"
            >
              <Icon name="lucide:search" class="w-5 h-5" />
              Browse Events
            </NuxtLink>
            <NuxtLink
              v-if="!authStore.isAuthenticated"
              to="/register"
              class="px-8 py-4 bg-white/20 backdrop-blur-sm text-white font-semibold rounded-xl hover:bg-white/30 transition-all transform hover:-translate-y-1 border border-white/30 flex items-center justify-center gap-2"
            >
              <Icon name="lucide:user-plus" class="w-5 h-5" />
              Get Started
            </NuxtLink>
            <NuxtLink
              v-else
              to="/my-events"
              class="px-8 py-4 bg-white/20 backdrop-blur-sm text-white font-semibold rounded-xl hover:bg-white/30 transition-all transform hover:-translate-y-1 border border-white/30 flex items-center justify-center gap-2"
            >
              <Icon name="lucide:calendar" class="w-5 h-5" />
              My Events
            </NuxtLink>
          </div>
        </div>
      </div>
      
      <!-- Wave Divider -->
      <div class="absolute bottom-0 left-0 right-0">
        <svg viewBox="0 0 1440 120" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M0 120L60 110C120 100 240 80 360 70C480 60 600 60 720 65C840 70 960 80 1080 75C1200 70 1320 50 1380 40L1440 30V120H1380C1320 120 1200 120 1080 120H360C240 120 120 120 60 120H0Z" fill="currentColor" class="text-gray-50 dark:text-gray-900"></path>
        </svg>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-16 bg-gray-50 dark:bg-gray-900">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-12">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white mb-4">
            Why <span class="bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">EventHub</span>?
          </h2>
          <p class="text-gray-600 dark:text-gray-400 max-w-2xl mx-auto">
            Everything you need to discover, create, and manage amazing events
          </p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-lg hover:shadow-xl transition-all hover:-translate-y-1">
            <div class="w-14 h-14 bg-blue-100 dark:bg-blue-900/30 rounded-xl flex items-center justify-center mb-4">
              <Icon name="lucide:compass" class="w-7 h-7 text-blue-600 dark:text-blue-400" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">Discover Events</h3>
            <p class="text-gray-600 dark:text-gray-400 text-sm">
              Find the perfect events near you. Filter by category, date, or location to discover your next adventure.
            </p>
          </div>
          
          <div class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-lg hover:shadow-xl transition-all hover:-translate-y-1">
            <div class="w-14 h-14 bg-purple-100 dark:bg-purple-900/30 rounded-xl flex items-center justify-center mb-4">
              <Icon name="lucide:calendar-plus" class="w-7 h-7 text-purple-600 dark:text-purple-400" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">Create Events</h3>
            <p class="text-gray-600 dark:text-gray-400 text-sm">
              Create and manage your own events. Reach thousands of attendees and build your community.
            </p>
          </div>
          
          <div class="bg-white dark:bg-gray-800 rounded-2xl p-6 shadow-lg hover:shadow-xl transition-all hover:-translate-y-1">
            <div class="w-14 h-14 bg-pink-100 dark:bg-pink-900/30 rounded-xl flex items-center justify-center mb-4">
              <Icon name="lucide:users" class="w-7 h-7 text-pink-600 dark:text-pink-400" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 dark:text-white mb-2">Connect & Follow</h3>
            <p class="text-gray-600 dark:text-gray-400 text-sm">
              Follow your favorite creators and events. Never miss out on updates and new opportunities.
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Events Preview Section -->
    <section class="py-16 bg-white dark:bg-gray-800">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-8">
          <div>
            <h2 class="text-2xl md:text-3xl font-bold text-gray-900 dark:text-white">
              Popular Events
            </h2>
            <p class="text-gray-600 dark:text-gray-400 mt-1">Discover trending events near you</p>
          </div>
          <NuxtLink
            to="/events"
            class="text-blue-600 dark:text-blue-400 hover:text-blue-700 font-medium flex items-center gap-1"
          >
            View All
            <Icon name="lucide:arrow-right" class="w-4 h-4" />
          </NuxtLink>
        </div>

        <!-- Loading -->
        <div v-if="loadingEvents" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>

        <!-- Events Grid using EventCard -->
        <div v-else-if="enrichedEvents.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
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
          <p class="text-gray-500 dark:text-gray-400">No events available yet</p>
        </div>
      </div>
    </section>

    <!-- Testimonials / Call to Action -->
    <section class="py-16 bg-gradient-to-r from-blue-600 to-purple-600">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h2 class="text-3xl md:text-4xl font-bold text-white mb-4">
          Ready to Get Started?
        </h2>
        <p class="text-xl text-white/90 mb-8">
          Join thousands of people discovering amazing events every day
        </p>
        <NuxtLink
          v-if="!authStore.isAuthenticated"
          to="/register"
          class="inline-flex items-center gap-2 px-8 py-4 bg-white text-gray-900 font-semibold rounded-xl hover:shadow-2xl transition-all transform hover:-translate-y-1"
        >
          <Icon name="lucide:user-plus" class="w-5 h-5" />
          Create Free Account
        </NuxtLink>
        <NuxtLink
          v-else
          to="/events"
          class="inline-flex items-center gap-2 px-8 py-4 bg-white text-gray-900 font-semibold rounded-xl hover:shadow-2xl transition-all transform hover:-translate-y-1"
        >
          <Icon name="lucide:search" class="w-5 h-5" />
          Explore Events
        </NuxtLink>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '~/stores/auth';
import { useQuery } from '@vue/apollo-composable';
import { useToast } from 'vue-toastification';
import { GET_POPULAR_EVENTS } from '~/graphql/eventQueries';
import EventCard from '~/components/events/EventCard.vue';

const router = useRouter();
const authStore = useAuthStore();
const toast = useToast();

// Fetch popular events
const { 
  result: eventsResult, 
  loading: loadingEvents, 
  error: eventsError,
  refetch: refetchEvents,
  onError 
} = useQuery(
  GET_POPULAR_EVENTS,
  null,
  {
    fetchPolicy: 'network-only',
    notifyOnNetworkStatusChange: true,
  }
);

onError((error) => {
  console.error('Error fetching events:', error);
  toast.error('Failed to load popular events');
});

const popularEvents = computed(() => {
  return eventsResult.value?.events || [];
});

// Enrich events with proper data structure for EventCard
const enrichedEvents = computed(() => {
  return popularEvents.value.map((event: any) => {
    // Get images from event_images relationship
    const images = event.event_images || [];
    const featuredImage = images.find((img: any) => img.is_featured)?.image_url || 
                          images[0]?.image_url || 
                          null;
    
    return {
      ...event,
      featured_image: featuredImage,
      images: images,
      // Add counts from aggregates
      bookmarks_count: event.bookmarks_aggregate?.aggregate?.count || 0,
      tickets_count: event.tickets_aggregate?.aggregate?.count || 0
    };
  });
});

const navigateToEvent = (eventId: string) => {
  router.push(`/events/${eventId}`);
};

// Refetch on mount
onMounted(async () => {
  await refetchEvents();
});
</script>

<style scoped>
@keyframes blob {
  0% {
    transform: translate(0px, 0px) scale(1);
  }
  33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
  100% {
    transform: translate(0px, 0px) scale(1);
  }
}

.animate-blob {
  animation: blob 7s infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.animation-delay-4000 {
  animation-delay: 4s;
}
</style>