<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 bg-gray-50 dark:bg-gray-900">
    <div class="max-w-md w-full space-y-8">

      <!-- Header -->
      <div>
        <h2 class="text-center text-3xl font-bold text-gray-900 dark:text-white">
          Welcome Back
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
          Sign in to your account to continue
        </p>
      </div>

      <!-- Form -->
      <VeeForm @submit="handleLogin" class="mt-8 space-y-6" v-slot="{ errors, meta }">

        <div class="space-y-4">

          <!-- Email -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Email Address
            </label>

            <VeeField
              name="email"
              v-model="form.email"
              type="email"
              placeholder="Enter your email"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-800 dark:border-gray-600 dark:text-white"
              :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.email }"
              rules="required|email"
            />

            <VeeErrorMessage name="email" class="text-red-500 text-sm mt-1" />
          </div>

          <!-- Password -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Password
            </label>

            <div class="relative">
              <VeeField
                name="password"
                v-model="form.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="Enter your password"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-800 dark:border-gray-600 dark:text-white"
                :class="{ 'border-red-500 focus:ring-red-500 focus:border-red-500': errors.password }"
                rules="required|min:6"
              />

              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300"
              >
                <span v-if="showPassword">👁️</span>
                <span v-else>👁️‍🗨️</span>
              </button>
            </div>

            <VeeErrorMessage name="password" class="text-red-500 text-sm mt-1" />
          </div>

        </div>

        <!-- Remember + forgot -->
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              type="checkbox"
              v-model="form.remember"
              @change="handleRememberMe"
              class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <label class="ml-2 text-sm text-gray-900 dark:text-gray-300">
              Remember me
            </label>
          </div>

          <a href="#" class="text-sm text-blue-600 hover:text-blue-500 dark:text-blue-400">
            Forgot password?
          </a>
        </div>

        <!-- Submit -->
        <button
          type="submit"
          :disabled="loading || !meta.valid"
          class="w-full py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="loading" class="inline-flex items-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Signing in...
          </span>
          <span v-else>Sign in</span>
        </button>

        <!-- Register -->
        <div class="text-center text-sm">
          <span class="text-gray-600 dark:text-gray-400">
            Don't have an account?
          </span>

          <NuxtLink
            to="/register"
            class="ml-1 font-medium text-blue-600 hover:text-blue-500 dark:text-blue-400"
          >
            Sign up
          </NuxtLink>
        </div>

      </VeeForm>

    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useToast } from 'vue-toastification'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const loading = ref(false)
const showPassword = ref(false)

const form = reactive({
  email: '',
  password: '',
  remember: false
})

// ---------------- SUBMIT ----------------
const handleLogin = async (values, { resetForm }) => {
  loading.value = true

  try {
    const result = await authStore.login(form.email, form.password)

    if (result.success) {
      toast.success('Welcome back!')
      // Redirect to the dashboard or landing home context path
      router.push('/events')
    } else {
      toast.error(result.error || 'Invalid credentials. Please try again.')
    }
  } catch (error) {
    console.error('Login action template error:', error)
    toast.error('An unexpected connection error occurred.')
  } finally {
    loading.value = false
  }
}

// ---------------- REMEMBER ME ----------------
const handleRememberMe = () => {
  if (process.client) {
    if (form.remember && form.email) {
      localStorage.setItem('remembered_email', form.email)
    } else if (!form.remember) {
      localStorage.removeItem('remembered_email')
    }
  }
}

onMounted(() => {
  // Safe validation fallback for server-side environments (Nuxt 3)
  if (process.client) {
    const savedEmail = localStorage.getItem('remembered_email')
    if (savedEmail) {
      form.email = savedEmail
      form.remember = true
    }
  }
})
</script>