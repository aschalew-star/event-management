<template>
  <div>
    <label class="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-1">
      Category *
    </label>
    
    <div class="relative">
      <VeeField 
        name="category_id" 
        as="select" 
        class="w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-700 rounded-lg px-4 py-2.5 text-sm text-gray-900 dark:text-white focus:ring-2 focus:ring-indigo-500 focus:border-transparent disabled:opacity-60"
        :value="modelValue"
        @input="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
        :disabled="loading"
      >
        <option value="">
          {{ loading ? 'Loading categories...' : 'Select a category' }}
        </option>
        <option v-for="cat in categories" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </option>
      </VeeField>

      <div v-if="loading" class="absolute right-8 top-3">
        <i class="fas fa-spinner fa-spin text-gray-400 text-xs"></i>
      </div>
    </div>

    <VeeErrorMessage name="category_id" class="text-red-500 text-xs mt-1 block" />
    
    <p v-if="errorMsg" class="text-red-500 text-xs mt-1">
      ⚠️ {{ errorMsg }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useQuery } from '@vue/apollo-composable'
import { gql } from 'graphql-tag'

defineProps({
  modelValue: {
    type: String,
    default: ''
  }
})

defineEmits(['update:modelValue'])

const categories = ref<{ id: string; name: string }[]>([])
const errorMsg = ref<string | null>(null)

// Exact working Hasura Query from your console
const GET_CATEGORIES_QUERY = gql`
  query GetCategories {
    categories {
      id
      name
    }
  }
`

const { onResult, onError, loading } = useQuery(GET_CATEGORIES_QUERY)

onResult((result) => {
  if (result.data?.categories) {
    categories.value = result.data.categories
    errorMsg.value = null
  }
})

onError((err) => {
  console.error('Apollo error fetching categories:', err)
  errorMsg.value = 'Failed to load categories from backend.'
})
</script>