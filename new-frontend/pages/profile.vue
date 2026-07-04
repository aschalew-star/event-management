<template>
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-16">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl">
        <Icon name="lucide:alert-circle" class="w-20 h-20 text-red-400 mx-auto mb-4" />
        <h3 class="text-2xl font-semibold text-gray-700 dark:text-gray-200 mb-2">Failed to Load Profile</h3>
        <p class="text-gray-500 dark:text-gray-400 mb-6">{{ error.message || 'Something went wrong' }}</p>
        <button
          @click="refetch"
          class="inline-flex items-center gap-2 px-6 py-3 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition-colors"
        >
          <Icon name="lucide:refresh-cw" class="w-5 h-5" />
          Try Again
        </button>
      </div>

      <!-- Profile Content -->
      <div v-else-if="profile" class="space-y-6">
        <!-- Profile Header -->
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl overflow-hidden">
          <!-- Cover Image -->
          <div class="relative h-48 md:h-64 bg-gradient-to-r from-blue-500 to-purple-600">
            <div class="absolute inset-0 bg-black/20"></div>
            <div class="absolute bottom-0 left-0 right-0 p-6 bg-gradient-to-t from-black/60 to-transparent">
              <div class="flex flex-col md:flex-row items-start md:items-end gap-4">
                <!-- Avatar -->
                <div class="relative -mb-12 md:-mb-16">
                  <img
                    :src="profile.avatar_url || '/images/default-avatar.jpg'"
                    :alt="profile.name"
                    class="w-24 h-24 md:w-32 md:h-32 rounded-full object-cover ring-4 ring-white dark:ring-gray-800 shadow-xl"
                    @error="handleAvatarError"
                  />
                  <button
                    @click="triggerAvatarUpload"
                    class="absolute bottom-0 right-0 p-1.5 bg-blue-600 hover:bg-blue-700 rounded-full text-white shadow-lg transition-colors"
                    title="Change avatar"
                  >
                    <Icon name="lucide:camera" class="w-4 h-4" />
                  </button>
                  <input
                    ref="avatarInput"
                    type="file"
                    accept="image/*"
                    class="hidden"
                    @change="handleAvatarUpload"
                  />
                </div>

                <!-- User Info -->
                <div class="flex-1 text-white md:pb-2">
                  <h1 class="text-2xl md:text-3xl font-bold">{{ profile.name }}</h1>
                  <p class="text-sm md:text-base text-white/80">{{ profile.email }}</p>
                  <div class="flex flex-wrap items-center gap-4 mt-2 text-sm">
                    <span class="flex items-center gap-1">
                      <Icon name="lucide:calendar" class="w-4 h-4" />
                      Joined {{ formatDate(profile.created_at) }}
                    </span>
                    <span class="flex items-center gap-1">
                      <Icon name="lucide:map-pin" class="w-4 h-4" />
                      {{ profile.location || 'Location not set' }}
                    </span>
                  </div>
                </div>

                <!-- Edit Button -->
                <button
                  @click="isEditing = !isEditing"
                  class="px-4 py-2 bg-white/20 hover:bg-white/30 backdrop-blur-sm text-white rounded-lg transition-colors text-sm font-medium"
                >
                  <Icon :name="isEditing ? 'lucide:x' : 'lucide:edit-2'" class="w-4 h-4 inline mr-2" />
                  {{ isEditing ? 'Cancel' : 'Edit Profile' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Stats Cards -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center justify-between">
              <div>
                <span class="text-sm text-gray-500 dark:text-gray-400">Events Created</span>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.totalEvents }}</p>
              </div>
              <div class="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
                <Icon name="lucide:calendar" class="w-5 h-5 text-blue-600 dark:text-blue-400" />
              </div>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center justify-between">
              <div>
                <span class="text-sm text-gray-500 dark:text-gray-400">Bookmarks</span>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.totalBookmarks }}</p>
              </div>
              <div class="w-10 h-10 bg-amber-100 dark:bg-amber-900/30 rounded-lg flex items-center justify-center">
                <Icon name="lucide:bookmark" class="w-5 h-5 text-amber-600 dark:text-amber-400" />
              </div>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center justify-between">
              <div>
                <span class="text-sm text-gray-500 dark:text-gray-400">Following</span>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.totalFollowing }}</p>
              </div>
              <div class="w-10 h-10 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center">
                <Icon name="lucide:heart" class="w-5 h-5 text-purple-600 dark:text-purple-400" />
              </div>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center justify-between">
              <div>
                <span class="text-sm text-gray-500 dark:text-gray-400">Followers</span>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ stats.totalFollowers }}</p>
              </div>
              <div class="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
                <Icon name="lucide:users" class="w-5 h-5 text-green-600 dark:text-green-400" />
              </div>
            </div>
          </div>
        </div>

        <!-- Edit Profile Form -->
        <div v-if="isEditing" class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Edit Profile</h3>
          <form @submit.prevent="updateProfile" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Full Name</label>
              <input
                v-model="editForm.name"
                type="text"
                class="w-full bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Bio</label>
              <textarea
                v-model="editForm.bio"
                rows="3"
                class="w-full bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Tell us about yourself..."
              ></textarea>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Location</label>
              <input
                v-model="editForm.location"
                type="text"
                class="w-full bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="City, Country"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Website</label>
              <input
                v-model="editForm.website"
                type="url"
                class="w-full bg-gray-50 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg px-4 py-2 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="https://yourwebsite.com"
              />
            </div>
            <div class="flex gap-3 pt-2">
              <button
                type="submit"
                :disabled="updating"
                class="px-6 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors disabled:opacity-50"
              >
                <Icon v-if="updating" name="lucide:loader-2" class="w-4 h-4 animate-spin inline mr-2" />
                {{ updating ? 'Saving...' : 'Save Changes' }}
              </button>
              <button
                type="button"
                @click="isEditing = false"
                class="px-6 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg transition-colors"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>

        <!-- Recent Activity -->
        <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Recent Activity</h3>
          <div v-if="recentActivity.length === 0" class="text-center py-8">
            <Icon name="lucide:activity" class="w-12 h-12 text-gray-300 dark:text-gray-600 mx-auto mb-3" />
            <p class="text-gray-500 dark:text-gray-400">No recent activity</p>
          </div>
          <div v-else class="space-y-3">
            <div
              v-for="activity in recentActivity"
              :key="activity.id"
              class="flex items-center gap-4 p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
              <div class="w-10 h-10 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center text-white">
                <Icon :name="getActivityIcon(activity.type)" class="w-5 h-5" />
              </div>
              <div class="flex-1">
                <p class="text-sm text-gray-800 dark:text-gray-200">{{ activity.message }}</p>
                <p class="text-xs text-gray-400">{{ formatRelativeTime(activity.created_at) }}</p>
              </div>
              <NuxtLink
                v-if="activity.link"
                :to="activity.link"
                class="text-blue-600 dark:text-blue-400 hover:underline text-sm"
              >
                View
              </NuxtLink>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <NuxtLink
            to="/my-events"
            class="flex items-center justify-center gap-2 p-4 bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-md transition-all hover:-translate-y-0.5"
          >
            <Icon name="lucide:calendar" class="w-5 h-5 text-blue-600" />
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">My Events</span>
          </NuxtLink>
          <NuxtLink
            to="/my-bookmarks"
            class="flex items-center justify-center gap-2 p-4 bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-md transition-all hover:-translate-y-0.5"
          >
            <Icon name="lucide:bookmark" class="w-5 h-5 text-amber-600" />
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Bookmarks</span>
          </NuxtLink>
          <NuxtLink
            to="/my-following"
            class="flex items-center justify-center gap-2 p-4 bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-md transition-all hover:-translate-y-0.5"
          >
            <Icon name="lucide:heart" class="w-5 h-5 text-purple-600" />
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Following</span>
          </NuxtLink>
          <NuxtLink
            to="/my-followers"
            class="flex items-center justify-center gap-2 p-4 bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-md transition-all hover:-translate-y-0.5"
          >
            <Icon name="lucide:users" class="w-5 h-5 text-green-600" />
            <span class="text-sm font-medium text-gray-700 dark:text-gray-300">Followers</span>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '~/stores/auth';
import { useQuery, useMutation } from '@vue/apollo-composable';
import { useToast } from 'vue-toastification';
import { GET_USER_PROFILE } from '~/graphql/userQueries';
import { UPDATE_USER_PROFILE, UPDATE_USER_AVATAR } from '~/graphql/userMutations';

definePageMeta({
  middleware: 'auth'
});

const authStore = useAuthStore();
const toast = useToast();

const isEditing = ref(false);
const updating = ref(false);
const avatarInput = ref<HTMLInputElement | null>(null);

// Edit form
const editForm = ref({
  name: '',
  bio: '',
  location: '',
  website: '',
});

// Fetch user profile
const { 
  result: profileResult, 
  loading, 
  error,
  refetch,
  onError 
} = useQuery(
  GET_USER_PROFILE,
  () => ({ 
    userId: authStore.user?.id || '00000000-0000-0000-0000-000000000000' 
  }),
  {
    fetchPolicy: 'network-only',
    skip: () => !authStore.isAuthenticated || !authStore.user?.id,
    notifyOnNetworkStatusChange: true,
  }
);

onError((error) => {
  console.error('Error fetching profile:', error);
  toast.error('Failed to load profile');
});

const profile = computed(() => {
  const user = profileResult.value?.users_by_pk;
  if (user) {
    // Populate edit form when profile loads
    if (!isEditing.value) {
      editForm.value = {
        name: user.name || '',
        bio: user.bio || '',
        location: user.location || '',
        website: user.website || '',
      };
    }
  }
  return user;
});

// Stats
const stats = computed(() => ({
  totalEvents: profile.value?.events_aggregate?.aggregate?.count || 0,
  totalBookmarks: profile.value?.bookmarks_aggregate?.aggregate?.count || 0,
  totalFollowing: profile.value?.follows_aggregate?.aggregate?.count || 0,
  totalFollowers: 0, // Will be fetched separately
}));

// Recent activity
const recentActivity = computed(() => {
  const activities = [];
  const user = profile.value;
  
  if (user?.events && user.events.length > 0) {
    user.events.slice(0, 3).forEach((event: any) => {
      activities.push({
        id: event.id,
        type: 'event',
        message: `Created event "${event.title}"`,
        created_at: event.created_at,
        link: `/events/${event.id}`,
      });
    });
  }
  
  if (user?.bookmarks && user.bookmarks.length > 0) {
    user.bookmarks.slice(0, 2).forEach((bookmark: any) => {
      activities.push({
        id: bookmark.id,
        type: 'bookmark',
        message: `Bookmarked "${bookmark.event?.title || 'an event'}"`,
        created_at: bookmark.created_at,
        link: `/events/${bookmark.event?.id}`,
      });
    });
  }
  
  // Sort by date and limit to 5
  return activities
    .sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
    .slice(0, 5);
});

// Mutations
const { mutate: updateProfileMutation } = useMutation(UPDATE_USER_PROFILE);
const { mutate: updateAvatarMutation } = useMutation(UPDATE_USER_AVATAR);

// Update profile
const updateProfile = async () => {
  try {
    updating.value = true;
    
    await updateProfileMutation({
      id: authStore.user?.id,
      input: editForm.value,
    });
    
    await refetch();
    isEditing.value = false;
    toast.success('Profile updated successfully!');
  } catch (error: any) {
    console.error('Error updating profile:', error);
    toast.error(error.message || 'Failed to update profile');
  } finally {
    updating.value = false;
  }
};

// Avatar upload
const triggerAvatarUpload = () => {
  avatarInput.value?.click();
};

const handleAvatarUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement;
  const file = target.files?.[0];
  
  if (!file) return;
  
  // Validate file
  if (file.size > 5 * 1024 * 1024) {
    toast.error('Image size must be less than 5MB');
    return;
  }
  
  const validTypes = ['image/jpeg', 'image/png', 'image/webp'];
  if (!validTypes.includes(file.type)) {
    toast.error('Only JPEG, PNG, and WebP images are allowed');
    return;
  }
  
  try {
    // Convert to base64
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = async () => {
      const base64 = reader.result as string;
      
      await updateAvatarMutation({
        id: authStore.user?.id,
        avatar_url: base64,
      });
      
      await refetch();
      toast.success('Avatar updated successfully!');
    };
  } catch (error: any) {
    console.error('Error uploading avatar:', error);
    toast.error(error.message || 'Failed to upload avatar');
  }
};

// Helper functions
const handleAvatarError = (event: Event) => {
  const img = event.target as HTMLImageElement;
  img.src = '/images/default-avatar.jpg';
};

const getActivityIcon = (type: string) => {
  const icons: Record<string, string> = {
    event: 'lucide:calendar',
    bookmark: 'lucide:bookmark',
    follow: 'lucide:heart',
    ticket: 'lucide:ticket',
  };
  return icons[type] || 'lucide:activity';
};

const formatDate = (date: string) => {
  if (!date) return 'Unknown';
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
};

const formatRelativeTime = (date: string) => {
  if (!date) return '';
  const now = new Date();
  const past = new Date(date);
  const diffMinutes = Math.floor((now.getTime() - past.getTime()) / 1000 / 60);
  const diffHours = Math.floor(diffMinutes / 60);
  const diffDays = Math.floor(diffHours / 24);
  
  if (diffMinutes < 1) return 'Just now';
  if (diffMinutes < 60) return `${diffMinutes}m ago`;
  if (diffHours < 24) return `${diffHours}h ago`;
  if (diffDays < 7) return `${diffDays}d ago`;
  return formatDate(date);
};

// Refetch on mount
onMounted(async () => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    await refetch();
  }
});
</script>

<style scoped>
/* Add any custom styles here */
</style>