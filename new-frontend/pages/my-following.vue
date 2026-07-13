<!-- pages/my-following.vue -->
<template>
  <div
    class="min-h-screen bg-gradient-to-br from-purple-50 via-white to-pink-50 dark:from-gray-900 dark:via-gray-800 dark:to-gray-900"
  >
    <!-- Header -->
    <header
      class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700 sticky top-0 z-40"
    >
      <div class="max-w-7xl mx-auto px-4 py-4">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
              My Following
            </h1>
            <p class="text-sm text-gray-500 dark:text-gray-400">
              People you follow
            </p>
          </div>
          <div class="flex items-center gap-4">
            <NuxtLink
              to="/my-followers"
              class="flex items-center gap-2 px-4 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
            >
              <Icon name="lucide:users" class="w-4 h-4" />
              Followers
            </NuxtLink>
          </div>
        </div>
      </div>
    </header>

    <!-- Stats -->
    <div class="max-w-7xl mx-auto px-4 pt-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div
          class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow"
        >
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500 dark:text-gray-400"
              >Total Following</span
            >
            <Icon
              name="lucide:heart"
              class="w-5 h-5 text-purple-600 dark:text-purple-400"
            />
          </div>
          <p class="text-2xl font-bold text-gray-900 dark:text-white mt-2">
            {{ totalFollowing }}
          </p>
        </div>
        <div
          class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow"
        >
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500 dark:text-gray-400"
              >Total Hosted Events</span
            >
            <Icon
              name="lucide:calendar-check"
              class="w-5 h-5 text-blue-600 dark:text-blue-400"
            />
          </div>
          <p class="text-2xl font-bold text-gray-900 dark:text-white mt-2">
            {{ totalHostedEvents }}
          </p>
        </div>
        <div
          class="bg-white dark:bg-gray-800 rounded-xl shadow-sm p-4 hover:shadow-md transition-shadow"
        >
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500 dark:text-gray-400"
              >Active Profiles</span
            >
            <Icon
              name="lucide:activity"
              class="w-5 h-5 text-green-600 dark:text-green-400"
            />
          </div>
          <p class="text-2xl font-bold text-gray-900 dark:text-white mt-2">
            {{ activeFollowing }}
          </p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 pb-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-16">
        <div
          class="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"
        ></div>
      </div>

      <!-- Error State -->
      <div
        v-else-if="error"
        class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl"
      >
        <Icon
          name="lucide:alert-circle"
          class="w-20 h-20 text-red-400 mx-auto mb-4"
        />
        <h3
          class="text-2xl font-semibold text-gray-700 dark:text-gray-200 mb-2"
        >
          Failed to Load Following
        </h3>
        <p class="text-gray-500 dark:text-gray-400 mb-6">
          {{ error.message || "Something went wrong" }}
        </p>
        <button
          @click="refetch"
          class="inline-flex items-center gap-2 px-6 py-3 bg-purple-600 text-white rounded-xl hover:bg-purple-700 transition-colors"
        >
          <Icon name="lucide:refresh-cw" class="w-5 h-5" />
          Try Again
        </button>
      </div>

      <!-- Not Logged In State -->
      <div
        v-else-if="!authStore.isAuthenticated"
        class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl"
      >
        <Icon
          name="lucide:lock"
          class="w-20 h-20 text-gray-300 dark:text-gray-600 mx-auto mb-4"
        />
        <h3
          class="text-2xl font-semibold text-gray-700 dark:text-gray-200 mb-2"
        >
          Please Log In
        </h3>
        <p class="text-gray-500 dark:text-gray-400 mb-6 max-w-md mx-auto">
          Log in to see who you're following.
        </p>
        <NuxtLink
          to="/login"
          class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-xl hover:shadow-lg transition-all"
        >
          <Icon name="lucide:log-in" class="w-5 h-5" />
          Log In
        </NuxtLink>
      </div>

      <!-- Following Grid -->
      <div
        v-else-if="following.length > 0"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
      >
        <div
          v-for="follow in following"
          :key="follow.id"
          class="group bg-white dark:bg-gray-800 rounded-2xl shadow-xl overflow-hidden hover:shadow-2xl transition-all duration-300 hover:-translate-y-1 border border-gray-100 dark:border-gray-700"
        >
          <div class="p-6">
            <div class="flex items-start gap-4">
              <!-- Avatar -->
              <div class="relative">
                <img
                  :src="
                    follow.user?.avatar_url ||
                    '/images/default-avatar.jpg'
                  "
                  :alt="follow.user?.name"
                  class="w-16 h-16 rounded-full object-cover border-2 border-purple-100 dark:border-purple-900 group-hover:border-purple-400 transition-colors"
                  @error="handleAvatarError"
                />
              </div>

              <!-- User Info -->
              <div class="flex-1 min-w-0">
                <h3
                  class="font-semibold text-gray-900 dark:text-white truncate"
                >
                  {{ follow.user?.name || "Anonymous User" }}
                </h3>
                <p class="text-sm text-gray-500 dark:text-gray-400 truncate">
                  {{ follow.user?.email }}
                </p>
                <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">
                  Following since {{ formatDate(follow.created_at) }}
                </p>
              </div>
            </div>

            <!-- Bio -->
            <p
              v-if="follow.user?.bio"
              class="text-sm text-gray-600 dark:text-gray-400 mt-3 line-clamp-2"
            >
              {{ follow.user.bio }}
            </p>

            <!-- Actions -->
            <div
              class="mt-4 pt-3 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between"
            >
              <div class="flex items-center gap-2">
                <NuxtLink
                  :to="`/profile/${follow.followed_user_id}`"
                  class="px-3 py-1.5 text-sm bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-400 rounded-lg hover:bg-purple-200 dark:hover:bg-purple-900/50 transition-colors flex items-center gap-1"
                >
                  <Icon name="lucide:user" class="w-3 h-3" />
                  Profile
                </NuxtLink>
                
                <button
                  @click="handleUnfollow(follow.followed_user_id)"
                  class="px-3 py-1.5 text-sm bg-red-50 dark:bg-red-950/20 text-red-600 dark:text-red-400 rounded-lg hover:bg-red-100 dark:hover:bg-red-950/50 transition-colors flex items-center gap-1"
                >
                  <Icon name="lucide:user-minus" class="w-3 h-3" />
                  Unfollow
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div
        v-else
        class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl"
      >
        <Icon
          name="lucide:heart"
          class="w-20 h-20 text-gray-300 dark:text-gray-600 mx-auto mb-4"
        />
        <h3
          class="text-2xl font-semibold text-gray-700 dark:text-gray-200 mb-2"
        >
          Not Following Anyone Yet
        </h3>
        <p class="text-gray-500 dark:text-gray-400 mb-6 max-w-md mx-auto">
          You're not following anyone yet. Discover amazing events and creators
          to follow!
        </p>
        <NuxtLink
          to="/events"
          class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-xl hover:shadow-lg transition-all hover:-translate-y-0.5"
        >
          <Icon name="lucide:compass" class="w-5 h-5" />
          Discover Events
        </NuxtLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onActivated } from "vue";
import { useQuery, useMutation } from "@vue/apollo-composable";
import { useAuthStore } from "~/stores/auth";
import { useToast } from "vue-toastification";
import { GET_USER_FOLLOWINGs, UNFOLLOW_USER } from "~/graphql/eventQueries";

definePageMeta({
  middleware: "auth",
});

const authStore = useAuthStore();
const toast = useToast();

const userId = computed(() => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    return authStore.user.id;
  }
  return null;
});

// Fetch following for the current user
const {
  result: followingResult,
  loading,
  error,
  refetch,
  onError,
} = useQuery(
  GET_USER_FOLLOWINGs,
  () => ({
    userId: userId.value || "00000000-0000-0000-0000-000000000000",
  }),
  {
    fetchPolicy: "network-only",
    skip: () => !process.client || !authStore.isAuthenticated || !userId.value,
    notifyOnNetworkStatusChange: true,
  },
);

onError((error) => {
  console.error("Error fetching following:", error);
  toast.error("Failed to load following list");
});

const following = computed(() => {
  return followingResult.value?.follows || [];
});

// Unfollow Mutation setup
const { mutate: unfollowMutate } = useMutation(UNFOLLOW_USER);

const handleUnfollow = async (followedUserId: string) => {
  if (!userId.value) return;
  
  try {
    await unfollowMutate({
      followerId: userId.value,
      followedUserId: followedUserId
    });
    toast.success("Unfollowed successfully");
    refetch();
  } catch (err) {
    console.error("Unfollow mutation error:", err);
    toast.error("Failed to unfollow user");
  }
};

// Stats Calculations
const totalFollowing = computed(() => following.value.length);

const totalHostedEvents = computed(() => {
  let count = 0;
  following.value.forEach((f) => {
    if (f.user?.events) {
      count += f.user.events.length;
    }
  });
  return count;
});

const activeFollowing = computed(() => {
  return following.value.filter((f) => f.user).length;
});

const handleAvatarError = (event: Event) => {
  const img = event.target as HTMLImageElement;
  img.src = "/images/default-avatar.jpg";
};

const formatDate = (date: string) => {
  if (!date) return "";
  try {
    return new Date(date).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  } catch {
    return "";
  }
};

// Watch for authentication changes
watch(
  () => authStore.isAuthenticated,
  async (authenticated) => {
    if (authenticated && userId.value) {
      await refetch();
    }
  },
);

// Refetch hooks
onMounted(async () => {
  if (authStore.isAuthenticated && userId.value) {
    await refetch();
  }
});

onActivated(async () => {
  if (authStore.isAuthenticated && userId.value) {
    await refetch();
  }
});
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>