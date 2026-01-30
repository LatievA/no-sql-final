<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">Shopping Cart</h1>

    <div v-if="cartStore.items.length === 0" class="text-center py-12 bg-white rounded-xl shadow-md">
      <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
      <h3 class="text-lg font-medium text-gray-900">Your cart is empty</h3>
      <p class="text-gray-500 mb-4">Add some bicycles to get started!</p>
      <router-link to="/" class="btn btn-primary">Browse Catalog</router-link>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Cart Items -->
      <div class="lg:col-span-2 space-y-4">
        <div 
          v-for="(item, index) in cartStore.items" 
          :key="index"
          class="bg-white rounded-xl shadow-md p-4 flex gap-4"
        >
          <!-- Image -->
          <div class="w-24 h-24 bg-gradient-to-br from-gray-100 to-gray-200 rounded-lg flex-shrink-0 flex items-center justify-center">
            <svg class="w-12 h-12 text-gray-400" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="25" cy="65" r="18" stroke="currentColor" stroke-width="3"/>
              <circle cx="75" cy="65" r="18" stroke="currentColor" stroke-width="3"/>
              <path d="M25 65 L40 35 L60 35 L75 65" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>

          <!-- Details -->
          <div class="flex-1">
            <div class="flex justify-between">
              <div>
                <h3 class="font-semibold text-gray-900">{{ item.model_name }}</h3>
                <p class="text-sm text-gray-500">{{ item.brand }}</p>
              </div>
              <button @click="cartStore.removeItem(index)" class="text-red-500 hover:text-red-700">
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>

            <!-- Customizations -->
            <div v-if="item.selected_customizations?.length > 0" class="mt-2">
              <div v-for="custom in item.selected_customizations" :key="custom.name" class="text-sm text-gray-600">
                <span class="capitalize">{{ custom.name.replace('_', ' ') }}:</span> {{ custom.value }}
              </div>
            </div>

            <div class="flex justify-between items-center mt-3">
              <div class="flex items-center border rounded-lg">
                <button 
                  @click="cartStore.updateQuantity(index, item.quantity - 1)" 
                  class="px-2 py-1 hover:bg-gray-100"
                >-</button>
                <span class="px-3 py-1 border-x">{{ item.quantity }}</span>
                <button 
                  @click="cartStore.updateQuantity(index, item.quantity + 1)"
                  :disabled="item.quantity >= item.stock_quantity"
                  class="px-2 py-1 hover:bg-gray-100 disabled:opacity-50"
                >+</button>
              </div>
              <span class="font-bold text-primary-600">{{ formatPrice(item.price * item.quantity) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Order Summary -->
      <div class="lg:col-span-1">
        <div class="bg-white rounded-xl shadow-md p-6 sticky top-24">
          <h2 class="text-xl font-bold mb-4">Order Summary</h2>
          
          <div class="space-y-3 mb-4">
            <div class="flex justify-between text-gray-600">
              <span>Items ({{ cartStore.totalItems }})</span>
              <span>{{ formatPrice(cartStore.totalAmount) }}</span>
            </div>
            <div class="flex justify-between text-gray-600">
              <span>Shipping</span>
              <span>Free</span>
            </div>
          </div>

          <div class="border-t pt-4 mb-6">
            <div class="flex justify-between text-xl font-bold">
              <span>Total</span>
              <span class="text-primary-600">{{ formatPrice(cartStore.totalAmount) }}</span>
            </div>
          </div>

          <!-- Checkout Form -->
          <div v-if="authStore.isAuthenticated">
            <h3 class="font-semibold mb-3">Delivery Address</h3>
            <div class="space-y-3 mb-4">
              <input v-model="deliveryAddress.street" type="text" placeholder="Street Address" class="input" />
              <div class="grid grid-cols-2 gap-3">
                <input v-model="deliveryAddress.city" type="text" placeholder="City" class="input" />
                <input v-model="deliveryAddress.postal_code" type="text" placeholder="Postal Code" class="input" />
              </div>
              <input v-model="deliveryAddress.phone" type="text" placeholder="Phone Number" class="input" />
            </div>

            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">Payment Method</label>
              <select v-model="paymentMethod" class="input">
                <option value="card">Credit Card</option>
                <option value="cash">Cash on Delivery</option>
              </select>
            </div>

            <button 
              @click="placeOrder" 
              :disabled="ordering"
              class="w-full btn btn-primary py-3 disabled:opacity-50"
            >
              {{ ordering ? 'Processing...' : 'Place Order' }}
            </button>
          </div>

          <div v-else class="text-center">
            <p class="text-gray-600 mb-4">Please login to checkout</p>
            <router-link to="/auth" class="btn btn-primary w-full">Login to Continue</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import { orderApi } from '../api/endpoints'

const router = useRouter()
const cartStore = useCartStore()
const authStore = useAuthStore()
const toastStore = useToastStore()

const ordering = ref(false)
const paymentMethod = ref('card')
const deliveryAddress = reactive({
  street: '',
  city: '',
  postal_code: '',
  phone: ''
})

async function placeOrder() {
  // Validate
  if (!deliveryAddress.street || !deliveryAddress.city || !deliveryAddress.postal_code || !deliveryAddress.phone) {
    toastStore.error('Please fill in all delivery details')
    return
  }

  ordering.value = true
  try {
    const orderData = {
      items: cartStore.getOrderItems(),
      delivery_address: deliveryAddress,
      payment_method: paymentMethod.value
    }

    await orderApi.create(orderData)
    
    cartStore.clearCart()
    toastStore.success('Order placed successfully!')
    router.push('/orders')
  } catch (error) {
    const message = error.response?.data?.error || 'Failed to place order'
    toastStore.error(message)
  } finally {
    ordering.value = false
  }
}

function formatPrice(price) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'KZT',
    minimumFractionDigits: 0
  }).format(price)
}
</script>
