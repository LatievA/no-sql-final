<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Hero Section -->
    <div class="bg-gradient-to-r from-primary-600 to-primary-800 rounded-2xl p-8 mb-8 text-white">
      <h1 class="text-4xl font-bold mb-4">Sports Bicycle Store</h1>
      <p class="text-xl opacity-90 mb-6">Discover our collection of customizable bicycles for every adventure</p>
      <div class="flex flex-wrap gap-4">
        <div class="bg-white/20 rounded-lg px-4 py-2">
          <span class="font-bold">{{ totalBicycles }}</span> Bicycles
        </div>
        <div class="bg-white/20 rounded-lg px-4 py-2">
          <span class="font-bold">{{ categories.length }}</span> Categories
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-xl shadow-md p-6 mb-8">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <!-- Search -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <input 
            v-model="filters.search" 
            @input="debouncedFetch"
            type="text" 
            placeholder="Search bicycles..." 
            class="input"
          />
        </div>

        <!-- Category -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
          <select v-model="filters.category_id" @change="fetchBicycles" class="input">
            <option value="">All Categories</option>
            <option v-for="cat in categories" :key="cat.id" :value="cat.id">
              {{ cat.category_name }}
            </option>
          </select>
        </div>

        <!-- Price Range -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Min Price</label>
          <input 
            v-model.number="filters.min_price" 
            @input="debouncedFetch"
            type="number" 
            placeholder="0" 
            class="input"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Max Price</label>
          <input 
            v-model.number="filters.max_price" 
            @input="debouncedFetch"
            type="number" 
            placeholder="Any" 
            class="input"
          />
        </div>
      </div>

      <!-- Sort & Clear -->
      <div class="flex flex-wrap items-center justify-between mt-4 pt-4 border-t">
        <div class="flex items-center space-x-4">
          <select v-model="filters.sort" @change="fetchBicycles" class="input w-auto">
            <option value="created_at">Newest First</option>
            <option value="price">Price</option>
            <option value="model_name">Name</option>
          </select>
          <select v-model="filters.order" @change="fetchBicycles" class="input w-auto">
            <option value="desc">Descending</option>
            <option value="asc">Ascending</option>
          </select>
        </div>
        <button @click="clearFilters" class="text-primary-600 hover:text-primary-700 font-medium">
          Clear Filters
        </button>
      </div>
    </div>

    <!-- Loading -->
    <LoadingSpinner v-if="loading" />

    <!-- Bicycles Grid -->
    <div v-else>
      <div v-if="bicycles.length === 0" class="text-center py-12">
        <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900">No bicycles found</h3>
        <p class="text-gray-500">Try adjusting your filters</p>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        <BicycleCard v-for="bicycle in bicycles" :key="bicycle.id" :bicycle="bicycle" />
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="mt-8">
        <Pagination 
          :current-page="filters.page" 
          :total-pages="totalPages"
          @prev="prevPage"
          @next="nextPage"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { bicycleApi, categoryApi } from '../api/endpoints'
import BicycleCard from '../components/BicycleCard.vue'
import Pagination from '../components/Pagination.vue'
import LoadingSpinner from '../components/LoadingSpinner.vue'

const bicycles = ref([])
const categories = ref([])
const loading = ref(true)
const totalBicycles = ref(0)
const totalPages = ref(1)

const filters = ref({
  search: '',
  category_id: '',
  min_price: null,
  max_price: null,
  sort: 'created_at',
  order: 'desc',
  page: 1,
  limit: 12
})

let debounceTimeout = null

function debouncedFetch() {
  clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => {
    filters.value.page = 1
    fetchBicycles()
  }, 300)
}

async function fetchBicycles() {
  loading.value = true
  try {
    const params = { ...filters.value }
    // Remove empty values
    Object.keys(params).forEach(key => {
      if (params[key] === '' || params[key] === null) {
        delete params[key]
      }
    })
    
    const response = await bicycleApi.getAll(params)
    bicycles.value = response.data.data || []
    totalBicycles.value = response.data.total
    totalPages.value = response.data.total_pages
  } catch (error) {
    console.error('Failed to fetch bicycles:', error)
  } finally {
    loading.value = false
  }
}

async function fetchCategories() {
  try {
    const response = await categoryApi.getAll()
    categories.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  }
}

function clearFilters() {
  filters.value = {
    search: '',
    category_id: '',
    min_price: null,
    max_price: null,
    sort: 'created_at',
    order: 'desc',
    page: 1,
    limit: 12
  }
  fetchBicycles()
}

function prevPage() {
  if (filters.value.page > 1) {
    filters.value.page--
    fetchBicycles()
  }
}

function nextPage() {
  if (filters.value.page < totalPages.value) {
    filters.value.page++
    fetchBicycles()
  }
}

onMounted(() => {
  fetchCategories()
  fetchBicycles()
})
</script>
