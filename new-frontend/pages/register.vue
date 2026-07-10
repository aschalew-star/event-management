<template>
  <div class="min-h-screen flex items-center justify-center py-12 px-4 bg-gray-50 dark:bg-gray-900">
    <div class="max-w-md w-full space-y-8">

      <div>
        <h2 class="text-center text-3xl font-bold text-gray-900 dark:text-white">
          Create Account
        </h2>
      </div>

      <VeeForm @submit="handleRegister" class="space-y-6" v-slot="{ errors }">

        <!-- NAME -->
        <div>
          <VeeField
            name="name"
            v-model="form.name"
            placeholder="Full Name"
            class="w-full px-3 py-2 border rounded"
            rules="required"
          />
          <VeeErrorMessage name="name" class="text-red-500 text-sm" />
        </div>

        <!-- EMAIL -->
        <div>
          <VeeField
            name="email"
            v-model="form.email"
            placeholder="Email"
            type="email"
            class="w-full px-3 py-2 border rounded"
            rules="required|email"
          />
          <VeeErrorMessage name="email" class="text-red-500 text-sm" />
        </div>

        <!-- PASSWORD -->
        <div>
          <VeeField
            name="password"
            v-model="form.password"
            placeholder="Password"
            type="password"
            class="w-full px-3 py-2 border rounded"
            rules="required|min:6"
          />
          <VeeErrorMessage name="password" class="text-red-500 text-sm" />
        </div>

        <!-- CONFIRM PASSWORD -->
        <div>
          <VeeField
            name="password_confirmation"
            v-model="form.password_confirmation"
            placeholder="Confirm Password"
            type="password"
            class="w-full px-3 py-2 border rounded"
            rules="required|confirmed:@password"
          />
          <VeeErrorMessage name="password_confirmation" class="text-red-500 text-sm" />
        </div>

        <!-- TERMS -->
        <div>
          <label class="flex items-center gap-2">
            <VeeField
              name="terms"
              v-model="form.terms"
              type="checkbox"
              rules="required"
              :value="true"
              :unchecked-value="false"
            />
            <span>I agree to terms</span>
          </label>
          <VeeErrorMessage name="terms" class="text-red-500 text-sm" />
        </div>

        <!-- BUTTON (ALWAYS VISIBLE) -->
        <button
          type="submit"
          :disabled="loading"
          class="w-full py-2 bg-blue-600 text-white rounded"
        >
          {{ loading ? 'Creating...' : 'Create Account' }}
        </button>

      </VeeForm>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
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

// ---------------- SUBMIT ----------------
const handleRegister = async (values, { resetForm }) => {
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
</script>