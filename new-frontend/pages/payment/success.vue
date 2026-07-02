<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 flex items-center justify-center py-12 px-4">
    <div class="max-w-md w-full bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-8 text-center">
      <!-- Loading State -->
      <div v-if="loading" class="py-12">
        <div class="inline-block animate-spin rounded-full h-16 w-16 border-4 border-green-500 border-t-transparent"></div>
        <p class="mt-4 text-gray-600 dark:text-gray-400">Verifying your payment...</p>
      </div>

      <!-- Success State -->
      <div v-else-if="success">
        <div class="w-20 h-20 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-6">
          <i class="fas fa-check text-4xl text-green-600 dark:text-green-400"></i>
        </div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">Payment Successful! 🎉</h1>
        <p class="text-gray-600 dark:text-gray-400 mb-6">
          Your ticket has been confirmed. You'll receive a confirmation email shortly.
        </p>

        <div v-if="ticket" class="bg-gray-50 dark:bg-gray-700/50 rounded-xl p-4 mb-6 text-left">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-2">Ticket Details</h3>
          <div class="space-y-2 text-sm">
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">Event</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ ticket.event?.title || 'Event' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">Quantity</span>
              <span class="font-medium text-gray-900 dark:text-white">{{ ticket.quantity }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">Total Paid</span>
              <span class="font-medium text-gray-900 dark:text-white">ETB {{ ticket.total_price }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600 dark:text-gray-400">Status</span>
              <span class="px-2 py-1 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded-full text-xs font-medium">
                {{ ticket.status }}
              </span>
            </div>
          </div>
        </div>

        <div class="space-y-3">
          <NuxtLink
            :to="`/events/${eventId}`"
            class="block w-full py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-semibold rounded-lg transition-all transform hover:scale-[1.02]"
          >
            <i class="fas fa-ticket-alt mr-2"></i>
            View My Ticket
          </NuxtLink>
          <NuxtLink
            to="/events"
            class="block w-full py-3 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 font-semibold rounded-lg transition-all"
          >
            <i class="fas fa-arrow-left mr-2"></i>
            Browse More Events
          </NuxtLink>
        </div>
      </div>

      <!-- Error State -->
      <div v-else>
        <div class="w-20 h-20 bg-red-100 dark:bg-red-900/30 rounded-full flex items-center justify-center mx-auto mb-6">
          <i class="fas fa-times text-4xl text-red-600 dark:text-red-400"></i>
        </div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">Payment Failed</h1>
        <p class="text-gray-600 dark:text-gray-400 mb-6">
          {{ errorMessage || 'Something went wrong with your payment. Please try again.' }}
        </p>

        <div class="space-y-3">
          <button
            @click="retryPayment"
            class="block w-full py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-lg transition-all"
          >
            <i class="fas fa-redo mr-2"></i>
            Retry Payment
          </button>
          <NuxtLink
            :to="`/events/${eventId}`"
            class="block w-full py-3 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 font-semibold rounded-lg transition-all"
          >
            <i class="fas fa-arrow-left mr-2"></i>
            Back to Event
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery } from '@vue/apollo-composable'
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

const { mutate: verifyPayment } = useQuery(VERIFY_CHAPA_PAYMENT)

const verifyPaymentStatus = async () => {
  try {
    loading.value = true
    const txRef = route.query.tx_ref as string
    
    if (!txRef) {
      throw new Error('No transaction reference found')
    }

    transactionRef.value = txRef

    // You can add a query to get the ticket details
    const { data, error } = await useQuery(
      VERIFY_CHAPA_PAYMENT,
      { transactionRef: txRef }
    )

    if (error) {
      throw new Error(error.message)
    }

    const result = data?.verifyPayment
    
    if (result?.success) {
      success.value = true
      ticket.value = result.ticket
      eventId.value = result.ticket?.event_id
      toast.success('Payment verified successfully!')
    } else {
      errorMessage.value = result?.message || 'Payment verification failed'
      success.value = false
    }
  } catch (error: any) {
    console.error('Verification error:', error)
    errorMessage.value = error.message || 'Failed to verify payment'
    success.value = false
  } finally {
    loading.value = false
  }
}

const retryPayment = () => {
  if (eventId.value) {
    router.push(`/payment/${eventId.value}`)
  } else {
    router.push('/events')
  }
}

onMounted(() => {
  verifyPaymentStatus()
})
</script>