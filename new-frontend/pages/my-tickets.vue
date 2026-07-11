<template>
  <div class="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-900 dark:to-gray-800 py-8 transition-colors duration-300">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">My Tickets</h1>
          <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {{ tickets.length }} {{ tickets.length === 1 ? 'ticket' : 'tickets' }} purchased
          </p>
        </div>
        <div class="flex items-center gap-3">
          <button 
            @click="fetchTickets" 
            :disabled="loading || refreshing"
            class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
          >
            <Icon :name="refreshing ? 'lucide:loader-2' : 'lucide:refresh-cw'" class="w-4 h-4" :class="refreshing ? 'animate-spin' : ''" />
            Refresh
          </button>
          <NuxtLink
            to="/events"
            class="px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg transition-colors flex items-center gap-2"
          >
            <Icon name="lucide:search" class="w-4 h-4" />
            Browse Events
          </NuxtLink>
        </div>
      </div>

      <div v-if="tickets.length > 0" class="grid grid-cols-1 sm:grid-cols-4 gap-4 mb-6">
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Total Tickets</p>
              <p class="text-2xl font-bold text-gray-900 dark:text-white">{{ tickets.length }}</p>
            </div>
            <div class="w-10 h-10 bg-indigo-100 dark:bg-indigo-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:ticket" class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
            </div>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Confirmed</p>
              <p class="text-2xl font-bold text-green-600 dark:text-green-400">{{ confirmedCount }}</p>
            </div>
            <div class="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:check-circle" class="w-5 h-5 text-green-600 dark:text-green-400" />
            </div>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Pending</p>
              <p class="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{{ pendingCount }}</p>
            </div>
            <div class="w-10 h-10 bg-yellow-100 dark:bg-yellow-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:clock" class="w-5 h-5 text-yellow-600 dark:text-yellow-400" />
            </div>
          </div>
        </div>
        <div class="bg-white dark:bg-gray-800 rounded-xl shadow p-4">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm text-gray-500 dark:text-gray-400">Total Spent</p>
              <p class="text-2xl font-bold text-purple-600 dark:text-purple-400">${{ totalSpent.toFixed(2) }}</p>
            </div>
            <div class="w-10 h-10 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center">
              <Icon name="lucide:credit-card" class="w-5 h-5 text-purple-600 dark:text-purple-400" />
            </div>
          </div>
        </div>
      </div>
      
      <ClientOnly>
        <div v-if="loading" class="flex justify-center py-16">
          <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
        </div>

        <div v-else-if="!authStore.isAuthenticated" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700">
          <div class="text-5xl mb-4 text-gray-300 dark:text-gray-600">
            <Icon name="lucide:lock" class="w-16 h-16" />
          </div>
          <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-200 mb-2">Please Log In</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Log in to view and manage your tickets</p>
          <NuxtLink 
            to="/login"
            class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium rounded-xl hover:shadow-lg transition-all transform hover:scale-[1.02]"
          >
            <Icon name="lucide:log-in" class="w-4 h-4" />
            Log In
          </NuxtLink>
        </div>

        <div v-else-if="tickets.length === 0 && !loading" class="text-center py-16 bg-white dark:bg-gray-800 rounded-2xl shadow-xl border border-gray-100 dark:border-gray-700">
          <div class="text-5xl mb-4 text-gray-300 dark:text-gray-600">
            <Icon name="lucide:ticket" class="w-16 h-16" />
          </div>
          <h3 class="text-xl font-semibold text-gray-700 dark:text-gray-200 mb-2">No Tickets Yet</h3>
          <p class="text-gray-500 dark:text-gray-400 mb-6">Purchase tickets to events you're interested in!</p>
          <NuxtLink 
            to="/events"
            class="inline-flex items-center gap-2 px-6 py-3 bg-gradient-to-r from-indigo-600 to-purple-600 hover:from-indigo-700 hover:to-purple-700 text-white font-medium rounded-xl hover:shadow-lg transition-all transform hover:scale-[1.02]"
          >
            <Icon name="lucide:search" class="w-4 h-4" />
            Browse Events
          </NuxtLink>
        </div>

        <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <div 
            v-for="ticket in tickets" 
            :key="ticket.id"
            class="relative group"
          >
            <div class="absolute -top-2 -right-2 z-10 flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold shadow-lg backdrop-blur-sm"
                 :class="getStatusBadgeClass(ticket.status)">
              <Icon :name="getStatusIcon(ticket.status)" class="w-3 h-3" />
              {{ formatStatus(ticket.status) }}
            </div>

            <div class="absolute -top-2 -left-2 z-10 flex items-center gap-1.5 px-2.5 py-1 bg-indigo-600/90 backdrop-blur-sm text-white rounded-full text-xs font-semibold shadow-lg">
              <Icon name="lucide:ticket" class="w-3 h-3" />
              {{ ticket.quantity }} {{ ticket.quantity === 1 ? 'Ticket' : 'Tickets' }}
            </div>

            <button 
              @click.stop="viewTicketDetails(ticket)"
              class="absolute bottom-4 right-4 z-20 p-2 bg-white/95 dark:bg-gray-800/95 backdrop-blur-sm rounded-lg shadow-lg hover:scale-110 transition-all hover:bg-indigo-50 dark:hover:bg-indigo-950/50 md:opacity-0 md:group-hover:opacity-100 opacity-100"
              title="View Ticket Details"
            >
              <Icon name="lucide:eye" class="w-5 h-5 text-indigo-600 dark:text-indigo-400" />
            </button>

            <EventCard
              :event="enrichTicketWithEventData(ticket)"
              @click="navigateToEvent(ticket.event.id)"
            />
          </div>
        </div>

        <template #fallback>
          <div class="flex justify-center py-16">
            <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-500 border-t-transparent"></div>
          </div>
        </template>
      </ClientOnly>
    </div>

    <div v-if="showDetailsModal" class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4" @click.self="closeDetailsModal">
      <div class="bg-white dark:bg-gray-800 rounded-2xl max-w-lg w-full p-6 shadow-2xl animate-in fade-in zoom-in duration-200 max-h-[90vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-4 sticky top-0 bg-white dark:bg-gray-800 z-10 pb-2">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center gap-2">
            <Icon name="lucide:ticket" class="w-5 h-5 text-indigo-600" />
            Ticket & Payment Details
          </h3>
          <button @click="closeDetailsModal" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors">
            <Icon name="lucide:x" class="w-5 h-5" />
          </button>
        </div>
        
        <div v-if="selectedTicket" class="space-y-4">
          <div class="flex items-center gap-3 pb-3 border-b border-gray-100 dark:border-gray-700">
            <div class="w-12 h-12 rounded-lg bg-indigo-100 dark:bg-indigo-900/30 flex items-center justify-center text-indigo-600 dark:text-indigo-400">
              <Icon name="lucide:calendar" class="w-6 h-6" />
            </div>
            <div class="flex-1 min-w-0">
              <h4 class="font-semibold text-gray-900 dark:text-white line-clamp-1">{{ selectedTicket.event.title }}</h4>
              <p class="text-xs text-gray-500 dark:text-gray-400 line-clamp-1">{{ selectedTicket.event.venue || 'Venue TBD' }}</p>
            </div>
          </div>
          
          <div class="grid grid-cols-2 gap-3 text-sm">
            <div>
              <p class="text-gray-500 dark:text-gray-400">Ticket Status</p>
              <p class="font-medium text-gray-900 dark:text-white mt-0.5">
                <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs" :class="getStatusBadgeClass(selectedTicket.status)">
                  <Icon :name="getStatusIcon(selectedTicket.status)" class="w-3 h-3" />
                  {{ formatStatus(selectedTicket.status) }}
                </span>
              </p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400">Quantity</p>
              <p class="font-medium text-gray-900 dark:text-white mt-0.5">{{ selectedTicket.quantity }}</p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400">Total Price</p>
              <p class="font-medium text-gray-900 dark:text-white mt-0.5">${{ Number(selectedTicket.total_price).toFixed(2) }}</p>
            </div>
            <div>
              <p class="text-gray-500 dark:text-gray-400">Event Date</p>
              <p class="font-medium text-gray-900 dark:text-white mt-0.5">{{ formatDate(selectedTicket.event.event_date) }}</p>
            </div>
          </div>

          <div class="border-t border-gray-200 dark:border-gray-700 pt-3 mt-2">
            <h4 class="text-sm font-semibold text-gray-900 dark:text-white flex items-center gap-2 mb-3">
              <Icon name="lucide:credit-card" class="w-4 h-4 text-indigo-600" />
              Payment Information
            </h4>
            
            <div v-if="selectedTicket.payment" class="grid grid-cols-2 gap-3 text-sm bg-gray-50 dark:bg-gray-700/50 rounded-lg p-3">
              <div>
                <p class="text-gray-500 dark:text-gray-400">Payment ID</p>
                <p class="font-mono text-xs text-gray-900 dark:text-white break-all mt-0.5">{{ selectedTicket.payment.id }}</p>
              </div>
              <div>
                <p class="text-gray-500 dark:text-gray-400">Amount</p>
                <p class="font-medium text-gray-900 dark:text-white mt-0.5">
                  {{ selectedTicket.payment.currency || 'USD' }} ${{ Number(selectedTicket.payment.amount).toFixed(2) }}
                </p>
              </div>
              <div>
                <p class="text-gray-500 dark:text-gray-400">Payment Status</p>
                <p class="font-medium text-gray-900 dark:text-white mt-0.5">
                  <span class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs" :class="getPaymentStatusBadgeClass(selectedTicket.payment.status)">
                    <Icon :name="getPaymentStatusIcon(selectedTicket.payment.status)" class="w-3 h-3" />
                    {{ formatPaymentStatus(selectedTicket.payment.status) }}
                  </span>
                </p>
              </div>
              <div>
                <p class="text-gray-500 dark:text-gray-400">Payment Method</p>
                <p class="font-medium text-gray-900 dark:text-white mt-0.5">{{ selectedTicket.payment.payment_method || 'N/A' }}</p>
              </div>
              <div v-if="selectedTicket.payment.transaction_ref" class="col-span-2">
                <p class="text-gray-500 dark:text-gray-400">Transaction Reference</p>
                <p class="font-mono text-xs text-gray-900 dark:text-white break-all bg-white dark:bg-gray-800 p-2 rounded-lg mt-0.5 border border-gray-200 dark:border-gray-600">
                  {{ selectedTicket.payment.transaction_ref }}
                </p>
              </div>
              <div class="col-span-2">
                <p class="text-gray-500 dark:text-gray-400">Payment Date</p>
                <p class="font-medium text-gray-900 dark:text-white mt-0.5">{{ formatDate(selectedTicket.payment.created_at) }}</p>
              </div>
            </div>

            <div v-else-if="selectedTicket.payment_id" class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-3">
              <div class="grid grid-cols-2 gap-3 text-sm">
                <div class="col-span-2">
                  <p class="text-gray-500 dark:text-gray-400">Payment ID</p>
                  <p class="font-mono text-xs text-gray-900 dark:text-white break-all mt-0.5">{{ selectedTicket.payment_id }}</p>
                </div>
                <div class="col-span-2">
                  <p class="text-gray-500 dark:text-gray-400 text-xs italic">Full payment details are not available for this ticket.</p>
                </div>
              </div>
            </div>

            <div v-else class="bg-gray-50 dark:bg-gray-700/50 rounded-lg p-3 text-center text-sm text-gray-500 dark:text-gray-400">
              <Icon name="lucide:credit-card" class="w-5 h-5 mx-auto mb-1 text-gray-400" />
              <p>No payment information available for this ticket.</p>
            </div>
          </div>
          
          <div class="flex gap-3 pt-3 border-t border-gray-100 dark:border-gray-700">
            <button 
              @click="closeDetailsModal"
              class="flex-1 px-4 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg transition-colors text-sm font-medium"
            >
              Close
            </button>
            <button 
              @click="navigateToEvent(selectedTicket.event.id); closeDetailsModal()"
              class="flex-1 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors text-sm font-medium flex items-center justify-center gap-2"
            >
              <Icon name="lucide:eye" class="w-4 h-4" />
              View Event
            </button>
          </div>
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
import EventCard from '~/components/events/EventCard.vue'

definePageMeta({
  middleware: 'auth'
})

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const refreshing = ref(false)
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

// Computed stats
const confirmedCount = computed(() => {
  return tickets.value.filter((t: any) => t.status === 'confirmed').length
})

const pendingCount = computed(() => {
  return tickets.value.filter((t: any) => t.status === 'pending').length
})

const totalSpent = computed(() => {
  return tickets.value
    .filter((t: any) => t.status === 'confirmed')
    .reduce((sum: number, t: any) => sum + Number(t.total_price), 0)
})

const navigateToEvent = (eventId: string) => {
  router.push(`/events/${eventId}`)
}

// Enrich ticket with event data for EventCard
const enrichTicketWithEventData = (ticket: any) => {
  const event = ticket.event
  return {
    ...event,
    // Add ticket info to the event object
    ticket_id: ticket.id,
    ticket_status: ticket.status,
    ticket_quantity: ticket.quantity,
    ticket_total_price: ticket.total_price,
    ticket_created_at: ticket.created_at,
    // Ensure images array is properly structured
    images: event.event_images || [],
    // Set featured image from event_images
    featured_image: event.event_images?.find((img: any) => img.is_featured)?.image_url || 
                     event.event_images?.[0]?.image_url || 
                     null,
    // Add bookmark count
    bookmarks_count: event.bookmarks_aggregate?.aggregate?.count || 0
  }
}

// Helper functions for status badges
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

const getStatusIcon = (status: string) => {
  const icons: Record<string, string> = {
    'confirmed': 'lucide:check-circle',
    'pending': 'lucide:clock',
    'cancelled': 'lucide:x-circle',
    'used': 'lucide:check',
    'expired': 'lucide:alert-circle'
  }
  return icons[status] || 'lucide:circle'
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

// Payment status helpers
const getPaymentStatusBadgeClass = (status: string) => {
  const classes: Record<string, string> = {
    'confirmed': 'bg-green-500/90 text-white',
    'pending': 'bg-yellow-500/90 text-white',
    'failed': 'bg-red-500/90 text-white',
    'refunded': 'bg-gray-500/90 text-white',
    'completed': 'bg-green-500/90 text-white'
  }
  return classes[status] || 'bg-gray-500/90 text-white'
}

const getPaymentStatusIcon = (status: string) => {
  const icons: Record<string, string> = {
    'confirmed': 'lucide:check-circle',
    'pending': 'lucide:clock',
    'failed': 'lucide:x-circle',
    'refunded': 'lucide:undo-2',
    'completed': 'lucide:check'
  }
  return icons[status] || 'lucide:circle'
}

const formatPaymentStatus = (status: string) => {
  const statusMap: Record<string, string> = {
    'confirmed': 'Confirmed',
    'pending': 'Pending',
    'failed': 'Failed',
    'refunded': 'Refunded',
    'completed': 'Completed'
  }
  return statusMap[status] || status
}

const formatDate = (dateStr: string) => {
  if (!dateStr) return 'TBD'
  try {
    return new Date(dateStr).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch {
    return 'Invalid date'
  }
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
    refreshing.value = true
    try {
      await refetch()
      toast.success('Tickets refreshed!')
    } catch (error) {
      console.error('Error refreshing tickets:', error)
      toast.error('Failed to refresh tickets')
    } finally {
      refreshing.value = false
    }
  }
}

// Watch for authentication changes
watch(
  () => authStore.isAuthenticated,
  async (authenticated) => {
    if (authenticated && userId.value) {
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

// Refetch when component is activated
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

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes zoomIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-in {
  animation: fadeIn 0.2s ease-out;
}

.fade-in {
  animation: fadeIn 0.2s ease-out;
}

.zoom-in {
  animation: zoomIn 0.2s ease-out;
}

.relative {
  transition: all 0.2s ease;
}

.absolute {
  transition: all 0.2s ease;
}

/* Custom scrollbar for modal */
.max-h-[90vh]::-webkit-scrollbar {
  width: 4px;
}

.max-h-[90vh]::-webkit-scrollbar-track {
  background: transparent;
}

.max-h-[90vh]::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.5);
  border-radius: 9999px;
}

.max-h-[90vh]::-webkit-scrollbar-thumb:hover {
  background: rgba(156, 163, 175, 0.8);
}
</style>