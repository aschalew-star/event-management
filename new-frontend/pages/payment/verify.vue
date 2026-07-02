<template>
  <div class="min-h-screen bg-slate-50 flex items-center justify-center p-4">
    <div class="max-w-md w-full bg-white rounded-xl shadow-lg p-8 text-center border border-slate-100">
      
      <div v-if="loading" class="py-6 space-y-4">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-4 border-indigo-600 border-t-transparent"></div>
        <p class="text-slate-600 font-medium text-sm">Verifying transaction data status records...</p>
      </div>

      <div v-else-if="success && ticket" class="space-y-6">
        <div class="w-16 h-16 bg-emerald-100 rounded-full flex items-center justify-center mx-auto text-emerald-600 text-3xl font-bold">✓</div>
        <div>
          <h1 class="text-2xl font-bold text-slate-900">Payment Confirmed!</h1>
          <p class="text-slate-500 text-sm mt-1">Your allocation indices have been verified internally.</p>
        </div>
        
        <div class="bg-slate-50 rounded-lg p-4 text-left font-mono text-xs space-y-2 text-slate-700">
          <div><span class="text-slate-400">ALLOCATION_ID:</span> {{ ticket.id }}</div>
          <div><span class="text-slate-400">QUANTITY:</span> {{ ticket.quantity }}</div>
          <div><span class="text-slate-400">TOTAL_PRICE:</span> ETB {{ Number(ticket.total_price).toFixed(2) }}</div>
        </div>

        <button @click="navigateTo('/')" class="w-full py-3 bg-indigo-600 hover:bg-indigo-700 text-white font-semibold rounded-lg transition-colors">
          Return to Dashboard
        </button>
      </div>

      <div v-else class="space-y-6">
        <div class="w-16 h-16 bg-rose-100 rounded-full flex items-center justify-center mx-auto text-rose-600 text-3xl font-bold">×</div>
        <div>
          <h1 class="text-2xl font-bold text-slate-900">Verification Failure</h1>
          <p class="text-rose-500 text-sm mt-2 font-medium">{{ errorMessage }}</p>
        </div>
        <button @click="navigateTo('/')" class="w-full py-3 bg-slate-100 hover:bg-slate-200 text-slate-700 font-semibold rounded-lg transition-colors">
          Back to Catalog Browser
        </button>
      </div>

    </div>
  </div>
</template>


<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMutation } from '@vue/apollo-composable' // Actions act like Mutations/Queries
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

// Use useMutation because Actions that modify or verify data 
// are typically registered as mutations in your frontend definition
const { mutate: verifyPayment } = useMutation(VERIFY_CHAPA_PAYMENT)

const verifyPaymentStatus = async () => {
  try {
    loading.value = true
    const txRef = route.query.tx_ref as string
    
    if (!txRef) {
      throw new Error('No transaction reference found')
    }
    transactionRef.value = txRef

    // Execute through your Hasura GraphQL Engine seamlessly!
    const { data } = await verifyPayment({
      transactionRef: txRef
    })

    // Match your action response schema name: "verifyChapaPayment"
    const result = data?.verifyChapaPayment
    
    if (result?.success) {
      success.value = true
      ticket.value = result.ticket
      eventId.value = result.ticket?.event_id
      toast.success('Payment verified successfully via GraphQL!')
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
    router.push(`/events/${eventId.value}`)
  } else {
    router.push('/events')
  }
}

onMounted(() => {
  verifyPaymentStatus()
})
</script>