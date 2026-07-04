<template>
  <footer class="bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
        <!-- Brand -->
        <div class="col-span-1 md:col-span-2">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-10 h-10 bg-gradient-to-br from-blue-600 to-purple-600 rounded-xl flex items-center justify-center">
              <Icon name="lucide:calendar" class="w-6 h-6 text-white" />
            </div>
            <div>
              <span class="text-xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                EventHub
              </span>
              <p class="text-sm text-gray-500 dark:text-gray-400">Discover & Connect</p>
            </div>
          </div>
          <p class="text-gray-600 dark:text-gray-400 max-w-md">
            The ultimate platform for discovering, creating, and attending amazing events. 
            Connect with like-minded people and make memories that last.
          </p>
          <div class="flex gap-4 mt-4">
            <a href="#" class="text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
              <Icon name="lucide:facebook" class="w-5 h-5" />
            </a>
            <a href="#" class="text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
              <Icon name="lucide:twitter" class="w-5 h-5" />
            </a>
            <a href="#" class="text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
              <Icon name="lucide:instagram" class="w-5 h-5" />
            </a>
            <a href="#" class="text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
              <Icon name="lucide:youtube" class="w-5 h-5" />
            </a>
            <a href="#" class="text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors">
              <Icon name="lucide:linkedin" class="w-5 h-5" />
            </a>
          </div>
        </div>

        <!-- Quick Links -->
        <div>
          <h3 class="font-semibold text-gray-900 dark:text-white mb-4">Quick Links</h3>
          <ul class="space-y-2">
            <li>
              <NuxtLink to="/events" class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors text-sm">
                Browse Events
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/my-events" class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors text-sm">
                My Events
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/my-bookmarks" class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors text-sm">
                Bookmarks
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/my-following" class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors text-sm">
                Following
              </NuxtLink>
            </li>
          </ul>
        </div>

        <!-- Support -->
        <div>
          <h3 class="font-semibold text-gray-900 dark:text-white mb-4">Support</h3>
          <ul class="space-y-2">
            <li>
              <NuxtLink to="/help" class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors text-sm">
                Help Center
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/contact" class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors text-sm">
                Contact Us
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/privacy" class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors text-sm">
                Privacy Policy
              </NuxtLink>
            </li>
            <li>
              <NuxtLink to="/terms" class="text-gray-600 dark:text-gray-400 hover:text-blue-600 dark:hover:text-blue-400 transition-colors text-sm">
                Terms of Service
              </NuxtLink>
            </li>
          </ul>
        </div>
      </div>

      <!-- Bottom Bar -->
      <div class="border-t border-gray-200 dark:border-gray-800 mt-8 pt-8 flex flex-col sm:flex-row justify-between items-center gap-4">
        <p class="text-sm text-gray-500 dark:text-gray-400">
          &copy; {{ new Date().getFullYear() }} EventHub. All rights reserved.
        </p>
        <div class="flex items-center gap-6 text-sm text-gray-500 dark:text-gray-400">
          <span>Made with</span>
          <Icon name="lucide:heart" class="w-4 h-4 text-red-500 animate-pulse" />
          <span>by EventHub Team</span>
        </div>
      </div>
    </div>
  </footer>
</template>