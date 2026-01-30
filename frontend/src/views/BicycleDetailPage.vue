<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <LoadingSpinner v-if="loading" />

    <div v-else-if="!bicycle" class="text-center py-12">
      <h2 class="text-2xl font-bold text-gray-900">Bicycle not found</h2>
      <router-link to="/" class="text-primary-600 hover:underline mt-4 inline-block">
        Back to catalog
      </router-link>
    </div>

    <div v-else>
      <!-- Back button -->
      <router-link to="/" class="inline-flex items-center text-gray-600 hover:text-gray-900 mb-6">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Back to catalog
      </router-link>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- Image -->
        <div class="bg-white rounded-xl shadow-md overflow-hidden">
          <div class="aspect-w-1 aspect-h-1 bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center h-96">
            <svg class="w-40 h-40 text-gray-400" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="25" cy="65" r="18" stroke="currentColor" stroke-width="3"/>
              <circle cx="75" cy="65" r="18" stroke="currentColor" stroke-width="3"/>
              <path d="M25 65 L40 35 L60 35 L75 65" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
              <path d="M40 35 L50 65 L75 65" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </div>
        </div>

        <!-- Details -->
        <div>
          <div class="bg-white rounded-xl shadow-md p-6">
            <div class="flex justify-between items-start mb-4">
              <div>
                <h1 class="text-3xl font-bold text-gray-900">{{ bicycle.model_name }}</h1>
                <p class="text-lg text-gray-500">{{ bicycle.brand }}</p>
              </div>
              <span v-if="bicycle.stock_quantity > 0" class="badge badge-success text-sm">
                {{ bicycle.stock_quantity }} in stock
              </span>
              <span v-else class="badge badge-danger text-sm">Out of Stock</span>
            </div>

            <p class="text-3xl font-bold text-primary-600 mb-6">{{ formatPrice(bicycle.price) }}</p>

            <p class="text-gray-600 mb-6">{{ bicycle.description }}</p>

            <!-- Specifications -->
            <div class="border-t pt-6 mb-6">
              <h3 class="text-lg font-semibold mb-4">Specifications</h3>
              <div class="grid grid-cols-2 gap-4">
                <div class="flex justify-between">
                  <span class="text-gray-500">Frame</span>
                  <span class="font-medium">{{ bicycle.specifications?.frame_material || 'N/A' }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-500">Wheel Size</span>
                  <span class="font-medium">{{ bicycle.specifications?.wheel_size || 'N/A' }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-500">Gears</span>
                  <span class="font-medium">{{ bicycle.specifications?.gear_count || 'N/A' }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-500">Brakes</span>
                  <span class="font-medium">{{ bicycle.specifications?.brake_type || 'N/A' }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-500">Weight</span>
                  <span class="font-medium">{{ bicycle.specifications?.weight || 'N/A' }}</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-500">Max Load</span>
                  <span class="font-medium">{{ bicycle.specifications?.max_load || 'N/A' }}</span>
                </div>
              </div>
            </div>

            <!-- Customization Options -->
            <div v-if="bicycle.customization_options?.length > 0" class="border-t pt-6 mb-6">
              <h3 class="text-lg font-semibold mb-4">Customize Your Bicycle</h3>
              <div class="space-y-4">
                <div v-for="option in bicycle.customization_options" :key="option.name">
                  <label class="block text-sm font-medium text-gray-700 mb-2 capitalize">
                    {{ option.name.replace('_', ' ') }}
                  </label>
                  <select 
                    v-model="selectedCustomizations[option.name]" 
                    class="input"
                  >
                    <option value="">Select {{ option.name.replace('_', ' ') }}</option>
                    <option v-for="opt in option.options" :key="opt" :value="opt">
                      {{ opt }}
                    </option>
                  </select>
                </div>
              </div>
            </div>

            <!-- Quantity & Add to Cart -->
            <div class="border-t pt-6">
              <div class="flex items-center space-x-4 mb-4">
                <label class="text-sm font-medium text-gray-700">Quantity:</label>
                <div class="flex items-center border rounded-lg">
                  <button @click="quantity > 1 && quantity--" class="px-3 py-2 hover:bg-gray-100">-</button>
                  <span class="px-4 py-2 border-x">{{ quantity }}</span>
                  <button @click="quantity < bicycle.stock_quantity && quantity++" class="px-3 py-2 hover:bg-gray-100">+</button>
                </div>
              </div>

              <button 
                @click="addToCart" 
                :disabled="bicycle.stock_quantity === 0"
                class="w-full btn btn-primary py-3 text-lg disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Add to Cart
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Reviews Section -->
      <div class="mt-8 bg-white rounded-xl shadow-md p-6">
        <div class="flex justify-between items-center mb-6">
          <h2 class="text-2xl font-bold">Reviews ({{ bicycle.reviews?.length || 0 }})</h2>
          <button 
            v-if="authStore.isAuthenticated" 
            @click="showReviewForm = !showReviewForm"
            class="btn btn-secondary"
          >
            Write a Review
          </button>
        </div>

        <!-- Review Form -->
        <div v-if="showReviewForm" class="mb-6 p-4 bg-gray-50 rounded-lg">
          <h3 class="font-semibold mb-4">Write Your Review</h3>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Rating</label>
            <div class="flex space-x-2">
              <button 
                v-for="i in 5" 
                :key="i" 
                @click="reviewForm.rating = i"
                class="text-2xl"
              >
                <span :class="i <= reviewForm.rating ? 'text-yellow-400' : 'text-gray-300'">★</span>
              </button>
            </div>
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Comment</label>
            <textarea 
              v-model="reviewForm.comment" 
              rows="3" 
              class="input"
              placeholder="Share your experience..."
            ></textarea>
          </div>
          <div class="flex space-x-3">
            <button @click="submitReview" class="btn btn-primary">Submit Review</button>
            <button @click="showReviewForm = false" class="btn btn-secondary">Cancel</button>
          </div>
        </div>

        <!-- Reviews List -->
        <div v-if="bicycle.reviews?.length > 0" class="space-y-4">
          <div v-for="review in bicycle.reviews" :key="review.review_id" class="border-b pb-4 last:border-0">
            <div class="flex justify-between items-start mb-2">
              <div>
                <span class="font-medium">{{ review.customer_name }}</span>
                <div class="flex text-yellow-400 mt-1">
                  <span v-for="i in 5" :key="i" :class="i <= review.rating ? 'text-yellow-400' : 'text-gray-300'">★</span>
                </div>
              </div>
              <span class="text-sm text-gray-500">{{ formatDate(review.review_date) }}</span>
            </div>
            <p class="text-gray-600">{{ review.comment }}</p>
          </div>
        </div>
        <p v-else class="text-gray-500 text-center py-4">No reviews yet. Be the first to review!</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { bicycleApi } from '../api/endpoints'
import { useCartStore } from '../stores/cart'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'
import LoadingSpinner from '../components/LoadingSpinner.vue'

const route = useRoute()
const cartStore = useCartStore()
const authStore = useAuthStore()
const toastStore = useToastStore()

const bicycle = ref(null)
const loading = ref(true)
const quantity = ref(1)
const selectedCustomizations = reactive({})
const showReviewForm = ref(false)
const reviewForm = reactive({
  rating: 5,
  comment: ''
})

async function fetchBicycle() {
  loading.value = true
  try {
    const response = await bicycleApi.getById(route.params.id)
    bicycle.value = response.data.data
    
    // Initialize customization selections
    if (bicycle.value.customization_options) {
      bicycle.value.customization_options.forEach(opt => {
        selectedCustomizations[opt.name] = ''
      })
    }
  } catch (error) {
    console.error('Failed to fetch bicycle:', error)
  } finally {
    loading.value = false
  }
}

function addToCart() {
  const customizations = Object.entries(selectedCustomizations)
    .filter(([_, value]) => value !== '')
    .map(([name, value]) => ({ name, value }))

  cartStore.addItem(bicycle.value, quantity.value, customizations)
  toastStore.success('Added to cart!')
}

async function submitReview() {
  if (reviewForm.rating < 1 || !reviewForm.comment.trim()) {
    toastStore.error('Please provide a rating and comment')
    return
  }

  try {
    await bicycleApi.addReview(bicycle.value.id, {
      rating: reviewForm.rating,
      comment: reviewForm.comment
    })
    
    toastStore.success('Review submitted!')
    showReviewForm.value = false
    reviewForm.rating = 5
    reviewForm.comment = ''
    
    // Refresh bicycle data
    fetchBicycle()
  } catch (error) {
    toastStore.error('Failed to submit review')
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
    day: 'numeric'
  })
}

onMounted(() => {
  fetchBicycle()
})
</script>
