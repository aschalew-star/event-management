<!-- components/Navigation.vue -->
<template>
  <nav class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4">
      <div class="flex items-center justify-between h-16">
        <!-- Logo -->
        <NuxtLink to="/" class="flex items-center gap-2">
          <Icon name="lucide:calendar" class="w-8 h-8 text-purple-600" />
          <span class="text-xl font-bold text-gray-900">EventHub</span>
        </NuxtLink>

        <!-- Navigation Links -->
        <div class="hidden md:flex items-center gap-6">
          <NuxtLink
            to="/events"
            class="text-gray-600 hover:text-purple-600 transition-colors"
          >
            Events
          </NuxtLink>
          <NuxtLink
            v-if="authStore.isAuthenticated"
            to="/my-events"
            class="text-gray-600 hover:text-purple-600 transition-colors"
          >
            My Events
          </NuxtLink>
          <NuxtLink
            v-if="authStore.isAuthenticated"
            to="/my-tickets"
            class="text-gray-600 hover:text-purple-600 transition-colors"
          >
            My Tickets
          </NuxtLink>
        </div>

        <!-- Right Section -->
        <div class="flex items-center gap-4">
          <template v-if="authStore.isAuthenticated">
            <!-- Follower Dropdown -->
            <div class="relative" @mouseenter="showFollowerMenu = true" @mouseleave="showFollowerMenu = false">
              <button class="flex items-center gap-2 px-3 py-2 text-gray-600 hover:text-purple-600 transition-colors">
                <Icon name="lucide:users" class="w-5 h-5" />
                <span class="hidden sm:inline">Social</span>
              </button>
              
              <!-- Dropdown Menu -->
              <div
                v-show="showFollowerMenu"
                class="absolute right-0 mt-2 w-56 bg-white rounded-xl shadow-xl border border-gray-100 py-2"
              >
                <NuxtLink
                  to="/my-followers"
                  class="flex items-center gap-3 px-4 py-2 hover:bg-purple-50 transition-colors"
                >
                  <Icon name="lucide:users" class="w-4 h-4 text-purple-600" />
                  <span>My Followers</span>
                  <span class="ml-auto text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded-full">
                    {{ followerStore.stats.total_followers }}
                  </span>
                </NuxtLink>
                <NuxtLink
                  to="/my-following"
                  class="flex items-center gap-3 px-4 py-2 hover:bg-purple-50 transition-colors"
                >
                  <Icon name="lucide:heart" class="w-4 h-4 text-pink-600" />
                  <span>Following</span>
                  <span class="ml-auto text-xs bg-pink-100 text-pink-700 px-2 py-0.5 rounded-full">
                    {{ followerStore.stats.total_following }}
                  </span>
                </NuxtLink>
              </div>
            </div>

            <!-- User Menu -->
            <div class="relative" @mouseenter="showUserMenu = true" @mouseleave="showUserMenu = false">
              <div class="flex items-center gap-2 cursor-pointer">
                <img
                  :src="authStore.user?.avatar_url || '/images/default-avatar.jpg'"
                  :alt="authStore.user?.name"
                  class="w-8 h-8 rounded-full object-cover"
                />
                <span class="hidden sm:inline text-sm text-gray-700">{{ authStore.user?.name }}</span>
              </div>

              <div
                v-show="showUserMenu"
                class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-xl border border-gray-100 py-2"
              >
                <NuxtLink
                  to="/profile"
                  class="flex items-center gap-3 px-4 py-2 hover:bg-purple-50 transition-colors"
                >
                  <Icon name="lucide:user" class="w-4 h-4" />
                  Profile
                </NuxtLink>
                <NuxtLink
                  to="/my-bookmarks"
                  class="flex items-center gap-3 px-4 py-2 hover:bg-purple-50 transition-colors"
                >
                  <Icon name="lucide:bookmark" class="w-4 h-4" />
                  Bookmarks
                </NuxtLink>
                <hr class="my-1 border-gray-100" />
                <button
                  @click="authStore.logout"
                  class="flex items-center gap-3 px-4 py-2 w-full hover:bg-red-50 text-red-600 transition-colors"
                >
                  <Icon name="lucide:log-out" class="w-4 h-4" />
                  Logout
                </button>
              </div>
            </div>
          </template>

          <NuxtLink
            v-else
            to="/login"
            class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
          >
            Login
          </NuxtLink>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { useFollowerStore } from '~/stores/follower'

const authStore = useAuthStore()
const followerStore = useFollowerStore()

const showFollowerMenu = ref(false)
const showUserMenu = ref(false)

onMounted(async () => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    await followerStore.fetchFollowers(authStore.user.id)
    await followerStore.fetchFollowing(authStore.user.id)
  }
})
</script>