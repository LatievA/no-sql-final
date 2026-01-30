<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">My Orders</h1>

    <LoadingSpinner v-if="loading" />

    <div v-else-if="orders.length === 0" class="text-center py-12 bg-white rounded-xl shadow-md">
      <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
      </svg>
      <h3 class="text-lg font-medium text-gray-900">No orders yet</h3>
      <p class="text-gray-500 mb-4">Start shopping to see your orders here!</p>
      <router-link to="/" class="btn btn-primary">Browse Catalog</router-link>
    </div>

    <div v-else class="space-y-6">
      <div 
        v-for="order in orders" 
        :key="order.id"
        class="bg-white rounded-xl shadow-md overflow-hidden"
      >
        <!-- Order Header -->
        <div class="bg-gray-50 px-6 py-4 flex flex-wrap justify-between items-center gap-4">
          <div>
            <p class="text-sm text-gray-500">Order ID</p>
            <p class="font-mono text-sm">{{ order.id }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Date</p>
            <p class="font-medium">{{ formatDate(order.order_date) }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Total</p>
            <p class="font-bold text-primary-600">{{ formatPrice(order.total_amount) }}</p>
          </div>
          <div>
            <span :class="getStatusClass(order.status)" class="badge">
              {{ order.status }}
            </span>
          </div>
        </div>

        <!-- Order Items -->
        <div class="p-6">
          <h3 class="font-semibold mb-4">Items</h3>
          <div class="space-y-4">
            <div 
              v-for="(item, index) in order.items" 
              :key="index"
              class="flex gap-4 pb-4 border-b last:border-0"
            >
              <div class="w-16 h-16 bg-gradient-to-br from-gray-100 to-gray-200 rounded-lg flex-shrink-0 flex items-center justify-center">
                <svg class="w-8 h-8 text-gray-400" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <circle cx="25" cy="65" r="18" stroke="currentColor" stroke-width="3"/>
                  <circle cx="75" cy="65" r="18" stroke="currentColor" stroke-width="3"/>
                  <path d="M25 65 L40 35 L60 35 L75 65" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </div>
              <div class="flex-1">
                <h4 class="font-medium">{{ item.model_name }}</h4>
                <p class="text-sm text-gray-500">{{ item.brand }}</p>
                <div v-if="item.selected_customizations?.length > 0" class="mt-1">
                  <span 
                    v-for="custom in item.selected_customizations" 
                    :key="custom.name"
                    class="text-xs bg-gray-100 rounded px-2 py-1 mr-2"
                  >
                    {{ custom.name.replace('_', ' ') }}: {{ custom.value }}
                  </span>
                </div>
              </div>
              <div class="text-right">
                <p class="font-medium">{{ formatPrice(item.price_at_purchase) }}</p>
                <p class="text-sm text-gray-500">Qty: {{ item.quantity }}</p>
              </div>
            </div>
          </div>

          <!-- Delivery Address -->
          <div class="mt-6 pt-4 border-t">
            <h3 class="font-semibold mb-2">Delivery Address</h3>
            <p class="text-gray-600">
              {{ order.delivery_address.street }}<br>
              {{ order.delivery_address.city }}, {{ order.delivery_address.postal_code }}<br>
              Phone: {{ order.delivery_address.phone }}
            </p>
          </div>

          <!-- Payment Info -->
          <div class="mt-4 flex gap-4 text-sm">
            <span class="text-gray-500">Payment: <span class="capitalize font-medium text-gray-900">{{ order.payment_method }}</span></span>
            <span :class="order.payment_status === 'paid' ? 'text-green-600' : 'text-yellow-600'" class="font-medium capitalize">
              {{ order.payment_status }}
            </span>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <Pagination 
        v-if="totalPages > 1"
        :current-page="page" 
        :total-pages="totalPages"
        @prev="page > 1 && fetchOrders(--page)"
        @next="page < totalPages && fetchOrders(++page)"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { orderApi } from '../api/endpoints'
import LoadingSpinner from '../components/LoadingSpinner.vue'
import Pagination from '../components/Pagination.vue'

const orders = ref([])
const loading = ref(true)
const page = ref(1)
const totalPages = ref(1)

async function fetchOrders(pageNum = 1) {
  loading.value = true
  try {
    const response = await orderApi.getMyOrders({ page: pageNum, limit: 10 })
    orders.value = response.data.data || []
    totalPages.value = response.data.total_pages
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  } finally {
    loading.value = false
  }
}

function formatPrice(price) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'KZT',
    minimumFractionDigits: 0
  }).format(price)
}

function formatDate(date) {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function getStatusClass(status) {
  const classes = {
    pending: 'badge-warning',
    confirmed: 'badge-info',
    shipped: 'badge-info',
    delivered: 'badge-success',
    cancelled: 'badge-danger'
  }
  return classes[status] || 'badge-info'
}

onMounted(() => {
  fetchOrders()
})
</script>
