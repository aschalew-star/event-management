<!-- components/events/EventCard.vue -->
<template>
  <div
    class="group bg-white dark:bg-gray-800 rounded-xl shadow-lg hover:shadow-2xl transition-all duration-300 overflow-hidden cursor-pointer border border-gray-200 dark:border-gray-700 hover:border-indigo-400 dark:hover:border-indigo-500 transform hover:-translate-y-1"
    @click="$emit('click')"
  >
    <!-- Image Section -->
    <div class="relative h-48 overflow-hidden bg-gray-200 dark:bg-gray-700">
      <img
        :src="eventImage || '/images/event-placeholder.jpg'"
        :alt="event.title"
        class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
        loading="lazy"
        @error="handleImageError"
      />
      
      <!-- Gradient Overlay for better text visibility -->
      <div class="absolute inset-0 bg-gradient-to-t from-black/50 via-transparent to-transparent" />
      
      <!-- Image counter badge -->
      <div v-if="event.images && event.images.length > 1" 
           class="absolute bottom-3 right-3 px-2 py-1 bg-black/60 backdrop-blur-sm rounded-full text-white text-xs font-medium flex items-center gap-1">
        <Icon name="lucide:image" class="w-3 h-3" />
        {{ event.images.length }}
      </div>
      
      <!-- Top Badges -->
      <div class="absolute top-3 left-3 flex flex-wrap gap-2">
        <!-- Status Badge -->
        <span
          v-if="event.status"
          class="px-2.5 py-1 text-xs font-medium rounded-full backdrop-blur-sm shadow-sm"
          :class="getStatusColor(event.status)"
        >
          <Icon :name="getStatusIcon(event.status)" class="w-3 h-3 inline mr-1" />
          {{ formatStatus(event.status) }}
        </span>
        
        <!-- Price Badge -->
        <span
          class="px-2.5 py-1 text-xs font-medium rounded-full backdrop-blur-sm shadow-sm"
          :class="event.is_free ? 'bg-green-500/80 text-white' : 'bg-indigo-600/80 text-white'"
        >
          <Icon :name="event.is_free ? 'lucide:gift' : 'lucide:dollar-sign'" class="w-3 h-3 inline mr-1" />
          {{ event.is_free ? 'Free' : `$${Number(event.price || 0).toFixed(2)}` }}
        </span>
      </div>

      <!-- Category Badge (Bottom Left) -->
      <div
        v-if="event.category"
        class="absolute bottom-3 left-3 px-2.5 py-1 bg-black/60 backdrop-blur-sm rounded-full text-white text-xs font-medium flex items-center gap-1.5"
        :style="categoryStyle"
      >
        <Icon v-if="event.category.icon" :name="event.category.icon" class="w-3 h-3" />
        <span>{{ event.category.name }}</span>
      </div>

      <!-- Date Badge (Bottom Right) -->
      <div class="absolute bottom-3 right-3 px-2.5 py-1 bg-black/60 backdrop-blur-sm rounded-full text-white text-xs font-medium flex items-center gap-1.5">
        <Icon name="lucide:calendar" class="w-3 h-3" />
        {{ formatDate(event.event_date) }}
      </div>
    </div>

    <!-- Content Section -->
    <div class="p-4">
      <!-- Title -->
      <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-1 line-clamp-1 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors">
        {{ event.title }}
      </h3>
      
      <!-- Description -->
      <p class="text-sm text-gray-600 dark:text-gray-400 mb-3 line-clamp-2">
        {{ event.description }}
      </p>

      <!-- Event Details -->
      <div class="space-y-1.5 text-sm">
        <!-- Venue/Address -->
        <div v-if="event.venue || event.address" class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
          <Icon name="lucide:map-pin" class="w-4 h-4 text-indigo-500 flex-shrink-0" />
          <span class="truncate">{{ event.venue || event.address }}</span>
        </div>
        
        <!-- Date & Time -->
        <div class="flex items-center gap-2 text-gray-600 dark:text-gray-400">
          <Icon name="lucide:calendar" class="w-4 h-4 text-indigo-500 flex-shrink-0" />
          <span>{{ formatDate(event.event_date) }}</span>
          <span v-if="event.start_time" class="text-xs text-gray-400">
            at {{ formatTime(event.start_time) }}
          </span>
        </div>
        
        <!-- Stats -->
        <div class="flex items-center gap-4 text-xs text-gray-500 dark:text-gray-400 mt-1">
          <div v-if="event.bookmarks_count !== undefined" class="flex items-center gap-1">
            <Icon name="lucide:bookmark" class="w-3.5 h-3.5" />
            <span>{{ event.bookmarks_count }}</span>
          </div>
          <div v-if="event.view_count !== undefined" class="flex items-center gap-1">
            <Icon name="lucide:eye" class="w-3.5 h-3.5" />
            <span>{{ event.view_count }}</span>
          </div>
          <div v-if="event.tickets_count !== undefined" class="flex items-center gap-1">
            <Icon name="lucide:ticket" class="w-3.5 h-3.5" />
            <span>{{ event.tickets_count }}</span>
          </div>
        </div>
      </div>

      <!-- Footer - Organizer & Actions -->
      <div class="mt-3 pt-3 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between">
        <!-- Organizer -->
        <div v-if="event.user" class="flex items-center gap-2">
          <div class="w-6 h-6 rounded-full bg-indigo-100 dark:bg-indigo-900 flex items-center justify-center text-indigo-600 dark:text-indigo-400 text-xs font-bold">
            {{ getInitials(event.user.name) }}
          </div>
          <span class="text-xs text-gray-600 dark:text-gray-400 truncate max-w-[100px]">
            {{ event.user.name }}
          </span>
        </div>
        
        <!-- Bookmark Button (if in bookmarks context) -->
        <button
          v-if="showRemoveBookmark"
          @click.stop="$emit('remove-bookmark', event.id)"
          class="p-1.5 rounded-full hover:bg-red-100 dark:hover:bg-red-900/30 text-gray-400 hover:text-red-500 transition-colors"
          title="Remove bookmark"
        >
          <Icon name="lucide:bookmark" class="w-4 h-4 fill-current text-red-500" />
        </button>
        
        <!-- View Details Link -->
        <span class="text-xs text-indigo-600 dark:text-indigo-400 font-medium flex items-center gap-1">
          View Details
          <Icon name="lucide:arrow-right" class="w-3 h-3" />
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

export interface EventCardProps {
  event: {
    id: string;
    title: string;
    description: string;
    price: number;
    is_free: boolean;
    venue: string;
    address: string;
    event_date: string;
    start_time?: string;
    end_time?: string;
    status: string;
    view_count?: number;
    featured_image?: string | null;
    images?: Array<{ id: string; image_url: string; is_featured: boolean }>;
    category?: {
      id: string;
      name: string;
      icon: string;
      color: string;
    } | null;
    user?: {
      id: string;
      name: string;
      email: string;
      avatar_url: string;
    } | null;
    bookmarks_count?: number;
    tickets_count?: number;
  };
  showRemoveBookmark?: boolean;
}

const props = withDefaults(defineProps<EventCardProps>(), {
  showRemoveBookmark: false,
});

const emit = defineEmits<{
  (e: 'click'): void;
  (e: 'remove-bookmark', eventId: string): void;
}>();

// Computed image with proper fallback
const eventImage = computed(() => {
  // Check for featured image
  if (props.event.featured_image) {
    return props.event.featured_image;
  }
  // Check for images array
  if (props.event.images && props.event.images.length > 0) {
    const featured = props.event.images.find(img => img.is_featured);
    return featured?.image_url || props.event.images[0]?.image_url;
  }
  return null;
});

// Category style with color
const categoryStyle = computed(() => {
  if (props.event.category?.color) {
    return { backgroundColor: props.event.category.color + 'CC' };
  }
  return {};
});

// Helper functions
const getStatusColor = (status: string): string => {
  const colors: Record<string, string> = {
    published: 'bg-green-500/80 text-white',
    upcoming: 'bg-blue-500/80 text-white',
    completed: 'bg-gray-500/80 text-white',
    cancelled: 'bg-red-500/80 text-white',
    draft: 'bg-yellow-500/80 text-white',
  };
  return colors[status] || 'bg-gray-500/80 text-white';
};

const getStatusIcon = (status: string): string => {
  const icons: Record<string, string> = {
    published: 'lucide:check-circle',
    upcoming: 'lucide:calendar-clock',
    completed: 'lucide:check',
    cancelled: 'lucide:x-circle',
    draft: 'lucide:file-text',
  };
  return icons[status] || 'lucide:circle';
};

const formatStatus = (status: string): string => {
  if (!status) return 'Draft';
  return status.charAt(0).toUpperCase() + status.slice(1);
};

const formatDate = (dateStr: string): string => {
  if (!dateStr) return 'TBD';
  try {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  } catch {
    return 'Invalid date';
  }
};

const formatTime = (timeStr: string): string => {
  if (!timeStr) return '';
  try {
    // Handle time in HH:mm format
    const [hours, minutes] = timeStr.split(':');
    const date = new Date();
    date.setHours(parseInt(hours), parseInt(minutes));
    return date.toLocaleTimeString('en-US', {
      hour: 'numeric',
      minute: '2-digit',
      hour12: true
    });
  } catch {
    return timeStr;
  }
};

const getInitials = (name?: string): string => {
  if (!name) return 'U';
  return name
    .split(' ')
    .map(word => word[0])
    .join('')
    .toUpperCase()
    .slice(0, 2);
};

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement;
  img.src = '/images/event-placeholder.jpg';
};
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