<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4 bg-gray-50 dark:bg-gray-900">
    <div class="max-w-md w-full space-y-8">

      <div>
        <h2 class="text-center text-3xl font-bold text-gray-900 dark:text-white">
          Create Account
        </h2>
      </div>

      <form @submit.prevent="handleRegister" class="space-y-6">

        <!-- NAME -->
        <div>
          <input
            v-model="form.name"
            placeholder="Full Name"
            class="w-full px-3 py-2 border rounded"
            @blur="validateField('name')"
          />
          <p v-if="errors.name" class="text-red-500 text-sm">{{ errors.name }}</p>
        </div>

        <!-- EMAIL -->
        <div>
          <input
            v-model="form.email"
            placeholder="Email"
            type="email"
            class="w-full px-3 py-2 border rounded"
            @blur="validateField('email')"
          />
          <p v-if="errors.email" class="text-red-500 text-sm">{{ errors.email }}</p>
        </div>

        <!-- PASSWORD -->
        <div>
          <input
            v-model="form.password"
            placeholder="Password"
            type="password"
            class="w-full px-3 py-2 border rounded"
            @blur="validateField('password')"
          />
          <p v-if="errors.password" class="text-red-500 text-sm">{{ errors.password }}</p>
        </div>

        <!-- CONFIRM PASSWORD -->
        <div>
          <input
            v-model="form.password_confirmation"
            placeholder="Confirm Password"
            type="password"
            class="w-full px-3 py-2 border rounded"
            @blur="validateField('password_confirmation')"
          />
          <p v-if="errors.password_confirmation" class="text-red-500 text-sm">
            {{ errors.password_confirmation }}
          </p>
        </div>

        <!-- TERMS -->
        <label class="flex items-center gap-2">
          <input type="checkbox" v-model="form.terms" />
          <span>I agree to terms</span>
        </label>
        <p v-if="errors.terms" class="text-red-500 text-sm">{{ errors.terms }}</p>

        <!-- BUTTON (ALWAYS VISIBLE) -->
        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2 bg-blue-600 text-white rounded"
        >
          {{ loading ? 'Creating...' : 'Create Account' }}
        </button>

      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '~/stores/auth'
import { useToast } from 'vue-toastification'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const loading = ref(false)

const form = reactive({
  name: '',
  email: '',
  password: '',
  password_confirmation: '',
  terms: false
})

const errors = reactive({
  name: '',
  email: '',
  password: '',
  password_confirmation: '',
  terms: ''
})

// ---------------- VALIDATION ----------------
const validateField = (field) => {
  if (field === 'name') {
    errors.name = form.name ? '' : 'Name required'
  }

  if (field === 'email') {
    errors.email = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)
      ? ''
      : 'Valid email required'
  }

  if (field === 'password') {
    errors.password = form.password.length >= 6
      ? ''
      : 'Min 6 characters'
  }

  if (field === 'password_confirmation') {
    errors.password_confirmation =
      form.password === form.password_confirmation
        ? ''
        : 'Passwords do not match'
  }
}

const validateAll = () => {
  validateField('name')
  validateField('email')
  validateField('password')
  validateField('password_confirmation')

  errors.terms = form.terms ? '' : 'Accept terms required'
}

// ---------------- SUBMIT ----------------
const handleRegister = async () => {
  validateAll()

  if (
    errors.name ||
    errors.email ||
    errors.password ||
    errors.password_confirmation ||
    errors.terms
  ) {
    toast.error('Fix errors first')
    return
  }

  loading.value = true

  try {
    const res = await authStore.register(
      form.email,
      form.password,
      form.name
    )

    console.log(res)

    if (res.success) {
      toast.success('Account created!')
      await router.push('/login')
    } else {
      toast.error(res.error || 'Failed')
    }
  } catch (e) {
    toast.error('Server error')
  } finally {
    loading.value = false
  }
}

// auto recheck password match
watch(() => form.password, validateField)
watch(() => form.password_confirmation, validateField)
</script>