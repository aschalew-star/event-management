<!-- pages/my-tickets.vue -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 py-8 transition-colors duration-300">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">My Tickets</h1>
        <button 
          @click="fetchTickets" 
          :disabled="loading"
          class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          <i :class="loading ? 'fas fa-spinner fa-spin' : 'fas fa-sync-alt'"></i>
          Refresh
        </button>
      </div>
      
      <ClientOnly>
        <!-- Loading State -->
        <div v-if="loading" class="flex justify-center py-16">
          <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
        </div>

        <!-- Not Logged In State -->
        <div v-else-if="!authStore.isAuthenticated" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700">
          <div class="text-5xl mb-4 text-gray-300 dark:text-gray-600">
            <i class="fas fa-lock"></i>
          </div>
          <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-200 mb-2">Please Log In</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Log in to view and manage your tickets</p>
          <NuxtLink 
            to="/login"
            class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium rounded-xl hover:shadow-lg transition-all transform hover:scale-[1.02]"
          >
            <i class="fas fa-sign-in-alt mr-1"></i>
            Log In
          </NuxtLink>
        </div>

        <!-- Logged In but No Tickets -->
        <div v-else-if="tickets.length === 0 && !loading" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700">
          <div class="text-5xl mb-4 text-gray-300 dark:text-gray-600">
            <i class="fas fa-ticket-alt"></i>
          </div>
          <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-200 mb-2">No Tickets Yet</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Purchase tickets to events you're interested in!</p>
          <NuxtLink 
            to="/events"
            class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium rounded-xl hover:shadow-lg transition-all transform hover:scale-[1.02]"
          >
            <i class="fas fa-search mr-1"></i>
            Browse Events
          </NuxtLink>
        </div>

        <!-- Display Tickets -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div 
            v-for="ticket in tickets" 
            :key="ticket.id"
            class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg overflow-hidden border border-gray-100 dark:border-gray-700 hover:shadow-2xl transition-all hover:-translate-y-1 group cursor-pointer"
            @click="navigateToEvent(ticket.event.id)"
          >
            <div class="relative h-48 bg-gray-200 dark:bg-gray-700">
              <img 
                :src="ticket.event.featured_image || '/images/event-placeholder.jpg'" 
                :alt="ticket.event.title"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                @error="handleImageError"
              />
              <!-- Status Badge -->
              <div class="absolute top-3 right-3">
                <span 
                  class="px-2.5 py-1 rounded-full text-xs font-semibold backdrop-blur-md shadow-sm"
                  :class="getStatusBadgeClass(ticket.status)"
                >
                  {{ formatStatus(ticket.status) }}
                </span>
              </div>
              <!-- Quantity Badge -->
              <div class="absolute bottom-3 left-3">
                <span 
                  class="px-2.5 py-1 rounded-full text-xs font-semibold backdrop-blur-md shadow-sm bg-black/50 text-white"
                >
                  <i class="fas fa-ticket-alt mr-1"></i>
                  {{ ticket.quantity }} {{ ticket.quantity === 1 ? 'Ticket' : 'Tickets' }}
                </span>
              </div>
              <!-- Price Badge -->
              <div class="absolute bottom-3 right-3">
                <span 
                  class="px-2.5 py-1 rounded-full text-xs font-semibold backdrop-blur-md shadow-sm"
                  :class="ticket.event.is_free ? 'bg-green-500/90 text-white' : 'bg-indigo-600/90 text-white'"
                >
                  {{ ticket.event.is_free ? 'Free' : `$${Number(ticket.total_price).toFixed(2)}` }}
                </span>
              </div>
            </div>
            
            <div class="p-5">
              <h3 class="font-bold text-gray-900 dark:text-white text-lg mb-2 line-clamp-1 group-hover:text-indigo-500 dark:group-hover:text-indigo-400 transition-colors">
                {{ ticket.event.title }}
              </h3>
              
              <div class="space-y-2 text-sm text-gray-600 dark:text-gray-300 mb-4">
                <div class="flex items-center gap-2.5">
                  <i class="fas fa-calendar-day text-indigo-500 w-4 text-center"></i>
                  <span>{{ formatDate(ticket.event.event_date) }}</span>
                </div>
                <div class="flex items-center gap-2.5">
                  <i class="fas fa-map-pin text-indigo-500 w-4 text-center"></i>
                  <span class="line-clamp-1">{{ ticket.event.venue || 'Venue TBD' }}</span>
                </div>
                <div class="flex items-center gap-2.5" v-if="ticket.transaction_ref">
                  <i class="fas fa-receipt text-indigo-500 w-4 text-center"></i>
                  <span class="text-xs">Ref: {{ ticket.transaction_ref }}</span>
                </div>
              </div>
              
              <!-- Ticket Footer -->
              <div class="pt-3 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between">
                <div class="flex flex-col">
                  <span class="text-xs text-gray-400 dark:text-gray-500">
                    <i class="far fa-clock mr-1"></i>
                    Purchased {{ formatRelativeTime(ticket.created_at) }}
                  </span>
                  <span class="text-xs text-gray-400 dark:text-gray-500" v-if="ticket.payment_id">
                    <i class="fas fa-credit-card mr-1"></i>
                    Payment: {{ ticket.payment_id.substring(0, 8) }}...
                  </span>
                </div>
                <div class="flex gap-2">
                  <!-- View Ticket Details Button -->
                  <button 
                    @click.stop="viewTicketDetails(ticket)"
                    class="px-3 py-1.5 bg-indigo-50 dark:bg-indigo-950/30 text-indigo-600 dark:text-indigo-400 hover:bg-indigo-100 dark:hover:bg-indigo-900/40 rounded-lg text-xs font-medium transition-colors"
                    title="View Ticket Details"
                  >
                    <i class="fas fa-eye mr-1"></i>
                    Details
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <template #fallback>
          <div class="flex justify-center py-16">
            <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
          </div>
        </template>
      </ClientOnly>
    </div>

    <!-- Ticket Details Modal -->
    <div v-if="showDetailsModal" class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeDetailsModal">
      <div class="bg-white dark:bg-gray-800 rounded-2xl max-w-md w-full p-6 shadow-2xl">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Ticket Details</h3>
          <button @click="closeDetailsModal" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200">
            <i class="fas fa-times"></i>
          </button>
        </div>
        
        <div v-if="selectedTicket" class="space-y-4">
          <div class="flex items-center gap-3 pb-3 border-b border-gray-100 dark:border-gray-700">
            <div class="w-12 h-12 rounded-lg bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center text-indigo-600 dark:text-indigo-400">
              <i class="fas fa-ticket-alt text-xl"></i>
            </div>
            <div>
              <h4 class="font-semibold text-gray-900 dark:text-white">{{ selectedTicket.event.title }}</h4>
              <p class="text-xs text-gray-500 dark:text-gray-400">{{ selectedTicket.event.venue }}</p>
            </div>
          </div>
          
          <div class="grid grid-cols-2 gap-3 text-sm">
            <div>
              <p class="text-gray-500 dark:text-gray-400">Status</p>
              <p class="font-medium text-gray-900 dark:text-white">
                <span class="inline-block px-2 py-0.5 rounded-full text-xs" :class="getStatusBadgeClass(selectedTicket.status)">
                  {{ formatStatus(selectedTicket.status) }}
                </span>
              </p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400">Quantity</p>
              <p class="font-medium text-gray-900 dark:text-white">{{ selectedTicket.quantity }}</p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400">Total Price</p>
              <p class="font-medium text-gray-900 dark:text-white">${{ Number(selectedTicket.total_price).toFixed(2) }}</p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400">Event Date</p>
              <p class="font-medium text-gray-900 dark:text-white">{{ formatDate(selectedTicket.event.event_date) }}</p>
            </div>
            <div v-if="selectedTicket.payment_id" class="col-span-2">
              <p class="text-gray-500 dark:text-gray-400">Payment ID</p>
              <p class="font-mono text-xs text-gray-900 dark:text-white break-all">{{ selectedTicket.payment_id }}</p>
            </div>
            <div v-if="selectedTicket.transaction_ref" class="col-span-2">
              <p class="text-gray-500 dark:text-gray-400">Transaction Reference</p>
              <p class="font-mono text-xs text-gray-900 dark:text-white break-all">{{ selectedTicket.transaction_ref }}</p>
            </div>
          </div>
          
          <button 
            @click="closeDetailsModal"
            class="w-full mt-4 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onActivated } from 'vue'
import { useRouter } from 'vue-router'
import { useQuery } from '@vue/apollo-composable'
import { useAuthStore } from '~/stores/auth'
import { useToast } from 'vue-toastification'
import { GET_USER_TICKETS } from '~/graphql/ticketQueries'

definePageMeta({
  middleware: 'auth'
})

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const showDetailsModal = ref(false)
const selectedTicket = ref<any>(null)

const userId = computed(() => {
  if (authStore.isAuthenticated && authStore.user?.id) {
    return authStore.user.id
  }
  return null
})

// Fetch user tickets
const { 
  result: ticketsResult, 
  loading, 
  refetch,
  onError
} = useQuery(
  GET_USER_TICKETS,
  () => ({ 
    userId: userId.value || '00000000-0000-0000-0000-000000000000' 
  }),
  {
    fetchPolicy: 'network-only',
    skip: () => {
      return !process.client || !authStore.isAuthenticated || !userId.value
    },
    notifyOnNetworkStatusChange: true,
  }
)

onError((error) => {
  console.error('Error fetching tickets:', error)
  toast.error('Failed to load tickets')
})

const tickets = computed(() => {
  return ticketsResult.value?.tickets || []
})

const navigateToEvent = (eventId: string) => {
  router.push(`/events/${eventId}`)
}

const handleImageError = (event: Event) => {
  const img = event.target as HTMLImageElement
  img.src = '/images/event-placeholder.jpg'
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return 'TBD'
  try {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    })
  } catch {
    return 'Invalid date'
  }
}

const formatRelativeTime = (dateStr: string) => {
  if (!dateStr) return ''
  
  try {
    const now = new Date()
    const past = new Date(dateStr)
    const diff = Math.floor((now.getTime() - past.getTime()) / (1000 * 60 * 60 * 24))
    
    if (diff <= 0) return 'today'
    if (diff === 1) return 'yesterday'
    if (diff < 7) return `${diff} days ago`
    if (diff < 30) return `${Math.floor(diff / 7)} weeks ago`
    
    return formatDate(dateStr)
  } catch {
    return ''
  }
}

const formatStatus = (status: string) => {
  const statusMap: Record<string, string> = {
    'confirmed': 'Confirmed',
    'pending': 'Pending',
    'cancelled': 'Cancelled',
    'used': 'Used',
    'expired': 'Expired'
  }
  return statusMap[status] || status
}

const getStatusBadgeClass = (status: string) => {
  const classes: Record<string, string> = {
    'confirmed': 'bg-green-500/90 text-white',
    'pending': 'bg-yellow-500/90 text-white',
    'cancelled': 'bg-red-500/90 text-white',
    'used': 'bg-gray-500/90 text-white',
    'expired': 'bg-gray-500/90 text-white'
  }
  return classes[status] || 'bg-indigo-600/90 text-white'
}

const viewTicketDetails = (ticket: any) => {
  selectedTicket.value = ticket
  showDetailsModal.value = true
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedTicket.value = null
}

const fetchTickets = async () => {
  if (authStore.isAuthenticated && userId.value) {
    try {
      await refetch()
    } catch (error) {
      console.error('Error refreshing tickets:', error)
    }
  }
}

// Watch for authentication changes
watch(
  () => authStore.isAuthenticated,
  async (authenticated) => {
    if (!authenticated && process.client) {
      // Optionally redirect or handle logout
    } else if (authenticated && userId.value) {
      await refetch()
    }
  }
)

// Refetch on mount
onMounted(async () => {
  if (authStore.isAuthenticated && userId.value) {
    await refetch()
  }
})

// Refetch when component is activated (coming back from event details)
onActivated(async () => {
  if (authStore.isAuthenticated && userId.value) {
    await refetch()
  }
})
</script>

<style scoped>
.line-clamp-1 {
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>