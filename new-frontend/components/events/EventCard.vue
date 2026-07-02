<template>
  <div
    class="group bg-white dark:bg-gray-800 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden cursor-pointer border border-gray-200 dark:border-gray-700 hover:border-indigo-400 dark:hover:border-indigo-500 transform hover:-translate-y-1"
    @click="$emit('click')"
  >
    <!-- Image -->
    <div class="relative h-48 overflow-hidden">
      <img
        :src="event.featured_image || '/images/event-placeholder.jpg'"
        :alt="event.title"
        class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
        loading="lazy"
      />
      
      <!-- Status Badge -->
      <div class="absolute top-3 left-3 flex gap-2">
        <span
          v-if="event.status"
          class="px-2 py-1 text-xs font-medium rounded-full backdrop-blur-sm"
          :class="statusColors[event.status] || 'bg-gray-500/80 text-white'"
        >
          {{ event.status }}
        </span>
        <span
          v-if="event.is_free"
          class="px-2 py-1 text-xs font-medium bg-green-500/80 backdrop-blur-sm text-white rounded-full"
        >
          Free
        </span>
      </div>

      <!-- Category Badge -->
      <div
        v-if="event.category"
        class="absolute bottom-3 left-3 px-2 py-1 bg-black/60 backdrop-blur-sm rounded-full text-white text-xs font-medium flex items-center gap-1"
      >
        <i :class="[event.category.icon || 'fas fa-tag']"></i>
        {{ event.category.name }}
      </div>
    </div>

    <!-- Content -->
    <div class="p-4">
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-1 line-clamp-1">
        {{ event.title }}
      </h3>
      
      <p class="text-sm text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">
        {{ event.description }}
      </p>

      <!-- Venue & Date -->
      <div class="space-y-2 text-sm">
        <div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
          <i class="fas fa-map-marker-alt text-indigo-500 w-4"></i>
          <span class="truncate">{{ event.venue || event.address || 'TBD' }}</span>
        </div>
        <div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
          <i class="fas fa-calendar-alt text-indigo-500 w-4"></i>
          <span>{{ formatDate(event.event_date) }}</span>
          <span v-if="event.start_time" class="text-xs text-gray-400">
            at {{ event.start_time }}
          </span>
        </div>
        <div v-if="!event.is_free && event.price" class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
          <i class="fas fa-tag text-indigo-500 w-4"></i>
          <span class="font-semibold text-indigo-600 dark:text-indigo-400">
            ${{ Number(event.price).toFixed(2) }}
          </span>
        </div>
      </div>

      <!-- Organizer -->
      <div v-if="event.user" class="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700 flex items-center gap-2">
        <div class="w-6 h-6 rounded-full bg-indigo-100 dark:bg-indigo-900 flex items-center justify-center text-indigo-600 dark:text-indigo-400 text-xs font-bold">
          {{ event.user.name?.charAt(0) || 'U' }}
        </div>
        <span class="text-xs text-gray-600 dark:text-gray-400">{{ event.user.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue';

const props = defineProps<{
  event: {
    id: string;
    title: string;
    description: string;
    price: number;
    is_free: boolean;
    venue: string;
    address: string;
    event_date: string;
    start_time: string;
    end_time: string;
    status: string;
    featured_image: string;
    category: {
      id: string;
      name: string;
      icon: string;
      color: string;
    };
    user: {
      id: string;
      name: string;
      email: string;
      avatar_url: string;
    };
  };
  viewMode: 'grid' | 'list';
}>();

const emit = defineEmits<{
  (e: 'click'): void;
}>();

const statusColors: Record<string, string> = {
  published: 'bg-green-500/80 text-white',
  upcoming: 'bg-blue-500/80 text-white',
  completed: 'bg-gray-500/80 text-white',
  cancelled: 'bg-red-500/80 text-white',
  draft: 'bg-yellow-500/80 text-white',
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return 'TBD';
  const date = new Date(dateStr);
  return date.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric', 
    year: 'numeric' 
  });
};
</script>