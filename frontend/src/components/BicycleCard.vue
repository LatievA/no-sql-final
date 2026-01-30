<template>
  <router-link :to="`/bicycles/${bicycle.id}`" class="card hover:shadow-lg transition-shadow duration-300">
    <div class="aspect-w-16 aspect-h-12 bg-gray-100">
      <div class="w-full h-48 flex items-center justify-center bg-gradient-to-br from-gray-100 to-gray-200">
        <svg class="w-20 h-20 text-gray-400" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="25" cy="65" r="18" stroke="currentColor" stroke-width="3"/>
          <circle cx="75" cy="65" r="18" stroke="currentColor" stroke-width="3"/>
          <path d="M25 65 L40 35 L60 35 L75 65" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M40 35 L50 65 L75 65" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </div>
    </div>
    <div class="p-4">
      <div class="flex justify-between items-start mb-2">
        <div>
          <h3 class="font-semibold text-gray-900 line-clamp-1">{{ bicycle.model_name }}</h3>
          <p class="text-sm text-gray-500">{{ bicycle.brand }}</p>
        </div>
        <span v-if="bicycle.stock_quantity > 0" class="badge badge-success">In Stock</span>
        <span v-else class="badge badge-danger">Out of Stock</span>
      </div>
      
      <div class="flex items-center justify-between mt-4">
        <span class="text-lg font-bold text-primary-600">{{ formatPrice(bicycle.price) }}</span>
        <div class="flex items-center text-yellow-400">
          <svg v-for="i in 5" :key="i" class="w-4 h-4" :class="i <= averageRating ? 'text-yellow-400' : 'text-gray-300'" fill="currentColor" viewBox="0 0 20 20">
            <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
          </svg>
          <span class="text-xs text-gray-500 ml-1">({{ bicycle.reviews?.length || 0 }})</span>
        </div>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  bicycle: {
    type: Object,
    required: true
  }
})

const averageRating = computed(() => {
  if (!props.bicycle.reviews || props.bicycle.reviews.length === 0) return 0
  const sum = props.bicycle.reviews.reduce((acc, r) => acc + r.rating, 0)
  return Math.round(sum / props.bicycle.reviews.length)
})

function formatPrice(price) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'KZT',
    minimumFractionDigits: 0
  }).format(price)
}
</script>
