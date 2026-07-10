<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 flex items-center justify-center p-4">
    <!-- Loading State -->
    <div v-if="loading" class="text-center">
      <div class="relative">
        <div class="w-20 h-20 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin mx-auto"></div>
        <div class="absolute inset-0 flex items-center justify-center">
          <svg class="w-8 h-8 text-indigo-600 animate-pulse" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
      </div>
      <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-300 mt-6">Verifying Your Payment</h3>
      <p class="text-gray-500 dark:text-gray-400 mt-2">Please wait while we confirm your transaction...</p>
    </div>

    <!-- Success State -->
    <div v-else-if="success" class="max-w-md w-full bg-white dark:bg-gray-800 rounded-2xl shadow-2xl overflow-hidden transform transition-all">
      <div class="bg-gradient-to-r from-green-500 to-emerald-600 p-6">
        <div class="flex justify-center">
          <div class="w-20 h-20 bg-white/20 rounded-full flex items-center justify-center animate-bounce">
            <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7"/>
            </svg>
          </div>
        </div>
        <h2 class="text-2xl font-bold text-white text-center mt-4">Payment Successful!</h2>
        <p class="text-green-100 text-center mt-1">Your ticket has been confirmed</p>
      </div>

      <div class="p-6 space-y-4">
        <!-- Ticket Info -->
        <div v-if="ticket" class="bg-gray-50 dark:bg-gray-700/50 rounded-xl p-4 space-y-3">
          <h4 class="text-sm font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Ticket Details</h4>
          
          <div class="flex justify-between items-center border-b border-gray-200 dark:border-gray-600 pb-2">
            <span class="text-sm text-gray-600 dark:text-gray-300">Transaction ID</span>
            <span class="text-sm font-mono font-semibold text-gray-800 dark:text-gray-200">{{ transactionRef }}</span>
          </div>
          
          <div v-if="ticket.event_id" class="flex justify-between items-center border-b border-gray-200 dark:border-gray-600 pb-2">
            <span class="text-sm text-gray-600 dark:text-gray-300">Event ID</span>
            <span class="text-sm font-medium text-gray-800 dark:text-gray-200">{{ ticket.event_id }}</span>
          </div>
          
          <div v-if="ticket.status" class="flex justify-between items-center">
            <span class="text-sm text-gray-600 dark:text-gray-300">Status</span>
            <span class="px-3 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 text-xs font-semibold rounded-full">Active</span>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex flex-col gap-3">
          <button 
            @click="navigateTo('/my-tickets')"
            class="w-full bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-700 hover:to-indigo-800 text-white font-semibold py-3 px-4 rounded-xl transition-all duration-200 transform hover:scale-[1.02] shadow-md hover:shadow-lg flex items-center justify-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z"/>
            </svg>
            View My Tickets
          </button>
          
          <button 
            @click="navigateTo('/events')"
            class="w-full bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 font-semibold py-3 px-4 rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
            </svg>
            Browse More Events
          </button>
        </div>
      </div>
    </div>

    <!-- Error State -->
    <div v-else class="max-w-md w-full bg-white dark:bg-gray-800 rounded-2xl shadow-2xl overflow-hidden transform transition-all">
      <div class="bg-gradient-to-r from-red-500 to-rose-600 p-6">
        <div class="flex justify-center">
          <div class="w-20 h-20 bg-white/20 rounded-full flex items-center justify-center animate-shake">
            <svg class="w-12 h-12 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </div>
        </div>
        <h2 class="text-2xl font-bold text-white text-center mt-4">Payment Failed</h2>
        <p class="text-red-100 text-center mt-1">We couldn't verify your transaction</p>
      </div>

      <div class="p-6 space-y-4">
        <!-- Error Message -->
        <div class="bg-red-50 dark:bg-red-900/20 border-2 border-red-200 dark:border-red-800/50 rounded-xl p-4">
          <div class="flex items-start gap-3">
            <svg class="w-5 h-5 text-red-600 dark:text-red-400 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            <div>
              <p class="text-sm font-medium text-red-800 dark:text-red-300">Error Details</p>
              <p class="text-sm text-red-700 dark:text-red-400 mt-1">{{ errorMessage || 'An unexpected error occurred' }}</p>
            </div>
          </div>
        </div>

        <!-- Transaction Reference -->
        <div v-if="transactionRef" class="bg-gray-50 dark:bg-gray-700/50 rounded-xl p-4">
          <p class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider">Transaction Reference</p>
          <p class="text-sm font-mono font-semibold text-gray-800 dark:text-gray-200 mt-1">{{ transactionRef }}</p>
        </div>

        <!-- Action Buttons -->
        <div class="flex flex-col gap-3">
          <button 
            @click="retryPayment"
            class="w-full bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-700 hover:to-indigo-800 text-white font-semibold py-3 px-4 rounded-xl transition-all duration-200 transform hover:scale-[1.02] shadow-md hover:shadow-lg flex items-center justify-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
            </svg>
            Try Again
          </button>
          
          <button 
            @click="navigateTo('/support')"
            class="w-full bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 font-semibold py-3 px-4 rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 5.636a9 9 0 010 12.728m0 0a9 9 0 01-12.728 0m12.728 0L21 21M4.243 4.243a9 9 0 00-2.828 9.9m2.828-9.9L3 3"/>
            </svg>
            Contact Support
          </button>
        </div>
      </div>
    </div>
  </div>
</template>


<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMutation } from '@vue/apollo-composable'
import { useToast } from 'vue-toastification'
import { VERIFY_CHAPA_PAYMENT } from '~/graphql/paymentMutations'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const loading = ref(true)
const success = ref(false)
const errorMessage = ref('')
const ticket = ref<any>(null)
const eventId = ref('')
const transactionRef = ref('')

const { mutate: verifyPayment, onError } = useMutation(VERIFY_CHAPA_PAYMENT)

// Handle GraphQL errors
onError((error) => {
  console.error('GraphQL error:', error)
  errorMessage.value = error.message || 'Payment verification failed'
  loading.value = false
  success.value = false
})

const verifyPaymentStatus = async () => {
  try {
    loading.value = true
    success.value = false
    errorMessage.value = ''

    // Get transaction reference from URL query params
    const txRef = route.query.tx_ref as string || route.query.transaction_ref as string
    
    if (!txRef) {
      throw new Error('No transaction reference found in URL')
    }
    
    transactionRef.value = txRef
    console.log('[DEBUG] Verifying payment with tx_ref:', txRef)

    // FIXED: Wrap the arguments in an "input" object
    const { data } = await verifyPayment({
      input: {
        transactionRef: txRef
      }
    })

    console.log('[DEBUG] Verification response:', data)

    // Check response
    const result = data?.verifyChapaPayment
    
    if (!result) {
      throw new Error('No response from verification service')
    }

    if (result.success) {
      success.value = true
      ticket.value = result.ticket
      eventId.value = result.ticket?.event_id
      
      toast.success('Payment verified successfully!')
    } else {
      errorMessage.value = result.message || 'Payment verification failed'
      success.value = false
      toast.error(errorMessage.value)
    }
  } catch (error: any) {
    console.error('Verification error:', error)
    
    if (error.networkError) {
      errorMessage.value = 'Network error. Please check your connection.'
    } else if (error.graphQLErrors && error.graphQLErrors.length > 0) {
      errorMessage.value = error.graphQLErrors[0].message || 'GraphQL error occurred'
    } else {
      errorMessage.value = error.message || 'Failed to verify payment'
    }
    
    success.value = false
    toast.error(errorMessage.value)
  } finally {
    loading.value = false
  }
}

const retryPayment = () => {
  if (eventId.value) {
    router.push(`/events/${eventId.value}`)
  } else {
    router.push('/events')
  }
}

const navigateTo = (path: string) => {
  router.push(path)
}

onMounted(() => {
  verifyPaymentStatus()
})
</script>

<style scoped>
@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-8px); }
  75% { transform: translateX(8px); }
}
.animate-shake {
  animation: shake 0.5s ease-in-out;
}
</style>