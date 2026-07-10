<template>
  <header 
    class="sticky top-0 z-50 w-full transition-all duration-300"
    :class="scrolled ? 'bg-white/95 dark:bg-gray-900/95 shadow-lg backdrop-blur-md' : 'bg-white dark:bg-gray-900 shadow-sm'"
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16 md:h-20">
        <!-- Logo -->
        <NuxtLink 
          to="/" 
          class="flex items-center gap-3 group"
        >
          <div class="relative">
            <div class="w-10 h-10 bg-gradient-to-br from-blue-600 to-purple-600 rounded-xl flex items-center justify-center shadow-lg group-hover:scale-105 transition-transform duration-300">
              <Icon name="lucide:calendar" class="w-6 h-6 text-white" />
            </div>
            <div class="absolute -top-1 -right-1 w-3 h-3 bg-green-400 rounded-full animate-pulse"></div>
          </div>
          <div>
            <span class="text-xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              EventHub
            </span>
            <span class="text-xs text-gray-400 block -mt-1">Discover & Connect</span>
          </div>
        </NuxtLink>

        <!-- Desktop Navigation -->
        <nav class="hidden md:flex items-center gap-1">
          <NuxtLink
            to="/"
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-gray-100 dark:hover:bg-gray-800"
            :class="$route.path === '/' ? 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20' : 'text-gray-700 dark:text-gray-300'"
          >
            <Icon name="lucide:home" class="w-4 h-4 inline mr-2" />
            Home
          </NuxtLink>
          
          <NuxtLink
            to="/events"
            class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-gray-100 dark:hover:bg-gray-800"
            :class="$route.path.startsWith('/events') ? 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20' : 'text-gray-700 dark:text-gray-300'"
          >
            <Icon name="lucide:calendar" class="w-4 h-4 inline mr-2" />
            Events
          </NuxtLink>

          <!-- Protected Links -->
          <template v-if="authStore.isAuthenticated">
            <NuxtLink
              to="/my-events"
              class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-gray-100 dark:hover:bg-gray-800"
              :class="$route.path === '/my-events' ? 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20' : 'text-gray-700 dark:text-gray-300'"
            >
              <Icon name="lucide:folder" class="w-4 h-4 inline mr-2" />
              My Events
            </NuxtLink>
            
            <NuxtLink
              to="/my-bookmarks"
              class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 hover:bg-gray-100 dark:hover:bg-gray-800"
              :class="$route.path === '/my-bookmarks' ? 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/20' : 'text-gray-700 dark:text-gray-300'"
            >
              <Icon name="lucide:bookmark" class="w-4 h-4 inline mr-2" />
              Bookmarks
            </NuxtLink>
          </template>
        </nav>

        <!-- Right Section -->
        <div class="flex items-center gap-3">
          <!-- Dark Mode Toggle -->
          <button
            @click="toggleDarkMode"
            class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors relative"
            aria-label="Toggle dark mode"
          >
            <Icon 
              :name="isDark ? 'lucide:sun' : 'lucide:moon'" 
              class="w-5 h-5 text-gray-600 dark:text-gray-300 transition-all duration-300"
              :class="isDark ? 'rotate-0' : 'rotate-0'"
            />
            <span class="absolute -top-1 -right-1 w-2 h-2 bg-blue-500 rounded-full animate-pulse" v-if="isDark"></span>
          </button>

          <!-- Desktop User Menu -->
          <template v-if="authStore.isAuthenticated">
            <div class="relative" @click="toggleDropdown" @mouseleave="closeDropdown">
              <button
                class="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-gray-100 dark:hover:bg-gray-800 transition-all duration-200 group"
              >
                <div class="relative">
                  <img
                    :src="authStore.user?.avatar_url || '/images/default-avatar.jpg'"
                    alt="User avatar"
                    class="w-9 h-9 rounded-full object-cover ring-2 ring-gray-200 dark:ring-gray-700 group-hover:ring-blue-400 transition-all"
                    @error="handleAvatarError"
                  />
                  <div class="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-green-500 rounded-full border-2 border-white dark:border-gray-800"></div>
                </div>
                <div class="hidden lg:block text-left">
                  <p class="text-sm font-medium text-gray-700 dark:text-gray-200">
                    {{ authStore.user?.name || 'User' }}
                  </p>
                  <p class="text-xs text-gray-400">{{ authStore.user?.email }}</p>
                </div>
                <Icon name="lucide:chevron-down" class="w-4 h-4 text-gray-400 transition-transform duration-200" :class="isDropdownOpen ? 'rotate-180' : ''" />
              </button>

              <!-- Dropdown Menu -->
              <div
                v-if="isDropdownOpen"
                class="absolute right-0 mt-2 w-72 bg-white dark:bg-gray-800 rounded-2xl shadow-2xl border border-gray-100 dark:border-gray-700 overflow-hidden animate-slide-down"
              >
                <!-- User Info -->
                <div class="p-4 bg-gradient-to-r from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20">
                  <div class="flex items-center gap-3">
                    <img
                      :src="authStore.user?.avatar_url || '/images/default-avatar.jpg'"
                      alt="User avatar"
                      class="w-12 h-12 rounded-full object-cover ring-2 ring-white dark:ring-gray-800"
                      @error="handleAvatarError"
                    />
                    <div>
                      <p class="font-semibold text-gray-800 dark:text-white">{{ authStore.user?.name || 'User' }}</p>
                      <p class="text-sm text-gray-500 dark:text-gray-400">{{ authStore.user?.email }}</p>
                    </div>
                  </div>
                </div>

                <div class="p-2">
                  <NuxtLink
                    to="/profile"
                    class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-200"
                    @click="closeDropdown"
                  >
                    <Icon name="lucide:user" class="w-4 h-4" />
                    <span class="text-sm">Profile</span>
                  </NuxtLink>
                  
                  <NuxtLink
                    to="/my-followers"
                    class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-200"
                    @click="closeDropdown"
                  >
                    <Icon name="lucide:calendar" class="w-4 h-4" />
                    <span class="text-sm">My Followers</span>
                  </NuxtLink>
                  
                  <NuxtLink
                    to="/my-tickets"
                    class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-200"
                    @click="closeDropdown"
                  >
                    <Icon name="lucide:bookmark" class="w-4 h-4" />
                    <span class="text-sm">My Tickets</span>
                  </NuxtLink>
                  
                  <NuxtLink
                    to="/my-following"
                    class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors text-gray-700 dark:text-gray-200"
                    @click="closeDropdown"
                  >
                    <Icon name="lucide:heart" class="w-4 h-4" />
                    <span class="text-sm">Following</span>
                  </NuxtLink>

                  <div class="border-t border-gray-200 dark:border-gray-700 my-2"></div>

                  <button
                    @click="handleLogout"
                    class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors text-red-600 dark:text-red-400"
                  >
                    <Icon name="lucide:log-out" class="w-4 h-4" />
                    <span class="text-sm font-medium">Logout</span>
                  </button>
                </div>
              </div>
            </div>
          </template>

          <!-- Auth Buttons -->
          <template v-else>
            <NuxtLink
              to="/login"
              class="hidden sm:flex px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
            >
              Sign In
            </NuxtLink>
            <NuxtLink
              to="/register"
              class="px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 text-white text-sm font-medium rounded-xl shadow-lg hover:shadow-xl transition-all duration-300 transform hover:-translate-y-0.5"
            >
              <span class="hidden sm:inline">Get Started</span>
              <Icon name="lucide:user-plus" class="sm:hidden w-5 h-5" />
            </NuxtLink>
          </template>

          <!-- Mobile Menu Button -->
          <button
            @click="isMobileMenuOpen = !isMobileMenuOpen"
            class="md:hidden p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
            aria-label="Toggle menu"
          >
            <Icon :name="isMobileMenuOpen ? 'lucide:x' : 'lucide:menu'" class="w-6 h-6 text-gray-600 dark:text-gray-300" />
          </button>
        </div>
      </div>
    </div>

    <!-- Mobile Menu -->
    <div
      v-if="isMobileMenuOpen"
      class="md:hidden bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-700 animate-slide-down"
    >
      <div class="max-w-7xl mx-auto px-4 py-3 space-y-1">
        <NuxtLink
          to="/"
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-700 dark:text-gray-300"
          @click="isMobileMenuOpen = false"
        >
          <Icon name="lucide:home" class="w-5 h-5" />
          <span>Home</span>
        </NuxtLink>
        
        <NuxtLink
          to="/events"
          class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-700 dark:text-gray-300"
          @click="isMobileMenuOpen = false"
        >
          <Icon name="lucide:calendar" class="w-5 h-5" />
          <span>Events</span>
        </NuxtLink>

        <template v-if="authStore.isAuthenticated">
          <NuxtLink
            to="/my-events"
            class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-700 dark:text-gray-300"
            @click="isMobileMenuOpen = false"
          >
            <Icon name="lucide:folder" class="w-5 h-5" />
            <span>My Events</span>
          </NuxtLink>
          
          <NuxtLink
            to="/my-tickets"
            class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-700 dark:text-gray-300"
            @click="isMobileMenuOpen = false"
          >
            <Icon name="lucide:bookmark" class="w-5 h-5" />
            <span>My Tickets</span>
          </NuxtLink>

          <NuxtLink
            to="/my-following"
            class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-700 dark:text-gray-300"
            @click="isMobileMenuOpen = false"
          >
            <Icon name="lucide:heart" class="w-5 h-5" />
            <span>Following</span>
          </NuxtLink>

          <button
            @click="handleLogout"
            class="w-full flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors text-red-600 dark:text-red-400 mt-2"
          >
            <Icon name="lucide:log-out" class="w-5 h-5" />
            <span>Logout</span>
          </button>
        </template>

        <template v-else>
          <NuxtLink
            to="/login"
            class="flex items-center gap-3 px-3 py-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors text-gray-700 dark:text-gray-300"
            @click="isMobileMenuOpen = false"
          >
            <Icon name="lucide:log-in" class="w-5 h-5" />
            <span>Sign In</span>
          </NuxtLink>
          <NuxtLink
            to="/register"
            class="flex items-center gap-3 px-3 py-2 rounded-lg bg-gradient-to-r from-blue-600 to-purple-600 text-white hover:shadow-lg transition-all"
            @click="isMobileMenuOpen = false"
          >
            <Icon name="lucide:user-plus" class="w-5 h-5" />
            <span>Get Started</span>
          </NuxtLink>
        </template>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue';
import { useAuthStore } from '~/stores/auth';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';

const authStore = useAuthStore();
const router = useRouter();
const toast = useToast();

const isDropdownOpen = ref(false);
const isMobileMenuOpen = ref(false);
const scrolled = ref(false);
const isDark = ref(false);

// Toggle dropdown
const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value;
};

const closeDropdown = () => {
  isDropdownOpen.value = false;
};

// Handle scroll effect
const handleScroll = () => {
  scrolled.value = window.scrollY > 20;
};

// Toggle dark mode
const toggleDarkMode = () => {
  isDark.value = !isDark.value;
  document.documentElement.classList.toggle('dark', isDark.value);
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light');
};

// Handle logout
const handleLogout = async () => {
  try {
    await authStore.logout();
    closeDropdown();
    isMobileMenuOpen.value = false;
    toast.success('Logged out successfully');
    router.push('/');
  } catch (error) {
    console.error('Logout error:', error);
    toast.error('Failed to logout');
  }
};

// Handle avatar error
const handleAvatarError = (event: Event) => {
  const img = event.target as HTMLImageElement;
  img.src = '/images/default-avatar.jpg';
};

// Close dropdown on click outside
const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  if (!target.closest('.relative')) {
    isDropdownOpen.value = false;
  }
};

// Initialize theme
onMounted(() => {
  window.addEventListener('scroll', handleScroll);
  document.addEventListener('click', handleClickOutside);
  
  // Load theme preference
  const savedTheme = localStorage.getItem('theme');
  if (savedTheme) {
    isDark.value = savedTheme === 'dark';
    document.documentElement.classList.toggle('dark', isDark.value);
  } else if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
    isDark.value = true;
    document.documentElement.classList.add('dark');
  }
});

onBeforeUnmount(() => {
  window.removeEventListener('scroll', handleScroll);
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped>
@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-slide-down {
  animation: slideDown 0.2s ease-out;
}
</style>