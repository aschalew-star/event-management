<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 py-12">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center min-h-[60vh]">
        <div class="text-center">
          <div class="inline-block animate-spin rounded-full h-16 w-16 border-4 border-indigo-500 border-t-transparent"></div>
          <p class="mt-4 text-gray-600 dark:text-gray-400">Loading payment details...</p>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="flex items-center justify-center min-h-[60vh]">
        <div class="text-center max-w-md">
          <div class="text-6xl mb-4">😅</div>
          <h3 class="text-2xl font-bold text-gray-800 dark:text-gray-200 mb-2">Payment Error</h3>
          <p class="text-gray-600 dark:text-gray-400 mb-6">{{ error }}</p>
          <NuxtLink
            to="/events"
            class="inline-flex items-center gap-2 px-6 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
          >
            <i class="fas fa-arrow-left"></i>
            Back to Events
          </NuxtLink>
        </div>
      </div>

      <!-- Payment Form -->
      <div v-else-if="event" class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column: Payment Form -->
        <div class="lg:col-span-2">
          <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6">
            <div class="flex items-center gap-3 mb-6">
              <button
                @click="router.back()"
                class="p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
              >
                <i class="fas fa-arrow-left text-gray-600 dark:text-gray-400"></i>
              </button>
              <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Complete Payment</h1>
            </div>

            <!-- Order Summary -->
            <div class="mb-6 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-xl">
              <h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-2">Order Summary</h3>
              <div class="space-y-2 text-sm">
                <div class="flex justify-between text-gray-600 dark:text-gray-400">
                  <span>{{ event.title }}</span>
                  <span>ETB {{ Number(event.price).toFixed(2) }}</span>
                </div>
                <div class="flex justify-between text-gray-600 dark:text-gray-400">
                  <span>Quantity</span>
                  <span>{{ quantity }}</span>
                </div>
                <div class="border-t border-gray-200 dark:border-gray-600 pt-2 mt-2">
                  <div class="flex justify-between font-semibold text-gray-900 dark:text-white">
                    <span>Total</span>
                    <span>ETB {{ totalPrice.toFixed(2) }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Customer Details -->
            <div class="space-y-4 mb-6">
              <h3 class="text-sm font-semibold text-gray-900 dark:text-white">Customer Information</h3>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Email Address *
                </label>
                <input
                  v-model="customer.email"
                  type="email"
                  required
                  class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="your@email.com"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                  Phone Number *
                </label>
                <input
                  v-model="customer.phone"
                  type="tel"
                  required
                  class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                  placeholder="09XXXXXXXX"
                />
              </div>

              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    First Name *
                  </label>
                  <input
                    v-model="customer.firstName"
                    type="text"
                    required
                    class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    placeholder="John"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                    Last Name *
                  </label>
                  <input
                    v-model="customer.lastName"
                    type="text"
                    required
                    class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent dark:bg-gray-700 dark:text-white"
                    placeholder="Doe"
                  />
                </div>
              </div>
            </div>

            <!-- Payment Method Info -->
            <div class="mb-6 p-4 bg-indigo-50 dark:bg-indigo-900/20 rounded-xl">
              <div class="flex items-center gap-3">
                <img src="" alt="Chapa" class="h-8" />
                <div>
                  <p class="text-sm font-medium text-gray-900 dark:text-white">Pay with Chapa</p>
                  <p class="text-xs text-gray-600 dark:text-gray-400">Secure payment via Chapa</p>
                </div>
              </div>
            </div>
    
            <!-- Pay Button -->
            <button
              @click="processPayment"
              :disabled="processing || !isFormValid"
              class="w-full py-3.5 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 disabled:from-gray-400 disabled:to-gray-400 text-white font-semibold rounded-lg shadow-lg transition-all transform hover:scale-[1.02] disabled:scale-100 disabled:cursor-not-allowed text-base"
            >
              <i v-if="processing" class="fas fa-spinner fa-spin mr-2"></i>
              <i v-else class="fas fa-credit-card mr-2"></i>
              {{ processing ? 'Processing...' : `Pay ETB ${totalPrice.toFixed(2)} with Chapa` }}
            </button>

            <p class="mt-4 text-xs text-center text-gray-500 dark:text-gray-400">
              <i class="fas fa-lock mr-1"></i>
              Your payment is secure and encrypted via Chapa
            </p>
          </div>
        </div>

        <!-- Right Column: Event Info -->
        <div class="lg:col-span-1 space-y-6">
          <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6 sticky top-24">
            <h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-4">Event Details</h3>
            
            <div class="relative rounded-lg overflow-hidden mb-4">
              <img
                :src="event.featured_image || '/images/event-placeholder.jpg'"
                :alt="event.title"
                class="w-full h-40 object-cover"
              />
            </div>

            <h4 class="font-semibold text-gray-900 dark:text-white mb-2">{{ event.title }}</h4>
            
            <div class="space-y-2 text-sm text-gray-600 dark:text-gray-400">
              <div class="flex items-center gap-2">
                <i class="fas fa-calendar-day text-indigo-500"></i>
                <span>{{ formatDate(event.event_date, true) }}</span>
              </div>
              <div v-if="event.venue" class="flex items-center gap-2">
                <i class="fas fa-map-pin text-indigo-500"></i>
                <span>{{ event.venue }}</span>
              </div>
              <div class="flex items-center gap-2">
                <i class="fas fa-tag text-indigo-500"></i>
                <span>ETB {{ Number(event.price).toFixed(2) }}</span>
              </div>
            </div>

            <!-- Organizer -->
            <div v-if="event.user" class="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700">
              <p class="text-xs text-gray-500 dark:text-gray-400">Organized by</p>
              <div class="flex items-center gap-2 mt-1">
                <div class="w-8 h-8 rounded-full bg-gradient-to-r from-indigo-500 to-purple-500 flex items-center justify-center text-white font-bold text-sm">
                  {{ event.user.name?.charAt(0) || 'U' }}
                </div>
                <span class="text-sm font-medium text-gray-900 dark:text-white">{{ event.user.name }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery, useMutation } from '@vue/apollo-composable'
import { useAuthStore } from '~/stores/auth'
import { useToast } from 'vue-toastification'
import { GET_EVENT_BY_ID } from '~/graphql/eventQueries'
import { PROCESS_CHAPA_PAYMENT } from '~/graphql/paymentMutations'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const eventId = computed(() => route.params.id as string)
const quantity = ref(1)
const processing = ref(false)

const customer = ref({
  email: '',
  phone: '',
  firstName: '',
  lastName: ''
})

const isFormValid = computed(() => {
  return customer.value.email && 
         customer.value.phone && 
         customer.value.firstName && 
         customer.value.lastName
})

// Fetch event details
const {
  result: eventResult,
  loading: eventLoading,
  error: eventError,
} = useQuery(
  GET_EVENT_BY_ID,
  () => ({ id: eventId.value }),
  {
    fetchPolicy: 'cache-and-network',
  }
)

const event = computed(() => eventResult.value?.events_by_pk)

const totalPrice = computed(() => {
  if (!event.value) return 0
  return event.value.is_free ? 0 : Number(event.value.price) * quantity.value
})

const loading = computed(() => eventLoading.value)
const error = computed(() => {
  if (eventError.value) return eventError.value.message
  if (!authStore.isAuthenticated) return 'Please sign in to continue'
  return null
})

// Chapa payment mutation
const { mutate: processChapaPayment } = useMutation(PROCESS_CHAPA_PAYMENT)

const formatDate = (dateStr: string, full: boolean = false) => {
  if (!dateStr) return 'TBD'
  const date = new Date(dateStr)
  if (full) {
    return date.toLocaleDateString('en-US', {
      weekday: 'long',
      month: 'long',
      day: 'numeric',
      year: 'numeric',
    })
  }
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

const processPayment = async () => {
  if (!authStore.isAuthenticated) {
    toast.warning('Please sign in to continue')
    router.push('/login')
    return
  }

  if (!isFormValid.value) {
    toast.warning('Please fill in all required fields')
    return
  }

  if (event.value?.is_free) {
    // Handle free ticket directly
    try {
      processing.value = true
      // Redirect to event page with success
      toast.success('🎉 Free ticket claimed successfully!')
      router.push(`/events/${eventId.value}`)
    } catch (error: any) {
      toast.error(error.message || 'Failed to claim free ticket')
    } finally {
      processing.value = false
    }
    return
  }

  processing.value = true
  try {
    const result = await processChapaPayment({
      eventId: eventId.value,
      quantity: quantity.value,
      totalPrice: totalPrice.value,
      email: customer.value.email,
      phone: customer.value.phone,
      firstName: customer.value.firstName,
      lastName: customer.value.lastName
    })

    console.log('Chapa payment result:', result)

    if (result?.data?.processChapaPayment?.success) {
      const checkoutUrl = result.data.processChapaPayment.data.checkout_url
      const transactionRef = result.data.processChapaPayment.data.transaction_ref
      
      if (checkoutUrl) {
        // Redirect to Chapa checkout page
        window.location.href = checkoutUrl
      } else {
        toast.success('Payment initialized! Please check your email for the payment link.')
        router.push(`/events/${eventId.value}`)
      }
    } else {
      toast.error(result?.data?.processChapaPayment?.message || 'Payment initialization failed')
    }
  } catch (error: any) {
    console.error('Payment error:', error)
    toast.error(error.message || 'Payment processing failed')
  } finally {
    processing.value = false
  }
}

// Pre-fill customer data if available
if (authStore.isAuthenticated && authStore.user) {
  customer.value.email = authStore.user.email || ''
  customer.value.firstName = authStore.user.name?.split(' ')[0] || ''
  customer.value.lastName = authStore.user.name?.split(' ').slice(1).join(' ') || ''
}
</script>

<style scoped>
input {
  transition: all 0.2s;
}

input:focus {
  outline: none;
}
</style>