<!-- components/ImageUploadProgress.vue -->
<template>
  <div v-if="uploading" class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center">
    <div class="bg-white dark:bg-gray-800 rounded-2xl p-8 max-w-md w-full shadow-2xl">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
        <h3 class="text-xl font-semibold text-gray-900 dark:text-white mt-4">Uploading Images</h3>
        <p class="text-gray-500 dark:text-gray-400 mt-2">Please wait while we upload your images...</p>
        
        <div class="mt-4 w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
          <div 
            class="bg-gradient-to-r from-indigo-600 to-purple-600 h-2.5 rounded-full transition-all duration-300"
            :style="{ width: `${progress}%` }"
          ></div>
        </div>
        <p class="text-sm text-gray-500 dark:text-gray-400 mt-2">{{ progress }}%</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  uploading: boolean;
  progress: number;
}>();
</script>