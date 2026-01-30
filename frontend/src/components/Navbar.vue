<template>
  <nav class="bg-white shadow-md sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <router-link to="/" class="flex items-center space-x-2">
          <svg class="w-8 h-8 text-primary-600" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
            <circle cx="25" cy="65" r="18" stroke="currentColor" stroke-width="4"/>
            <circle cx="75" cy="65" r="18" stroke="currentColor" stroke-width="4"/>
            <path d="M25 65 L40 35 L60 35 L75 65" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
            <path d="M40 35 L50 65 L75 65" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <span class="text-xl font-bold text-gray-900">BikeStore</span>
        </router-link>

        <!-- Navigation -->
        <div class="hidden md:flex items-center space-x-4">
          <router-link to="/" class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md font-medium">
            Catalog
          </router-link>
          
          <router-link v-if="authStore.isAuthenticated" to="/orders" class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md font-medium">
            My Orders
          </router-link>
          
          <router-link v-if="authStore.isAdmin" to="/admin" class="text-gray-700 hover:text-primary-600 px-3 py-2 rounded-md font-medium">
            Admin
          </router-link>
        </div>

        <!-- Right side -->
        <div class="flex items-center space-x-4">
          <!-- Cart -->
          <router-link to="/cart" class="relative p-2 text-gray-700 hover:text-primary-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            <span v-if="cartStore.totalItems > 0" class="absolute -top-1 -right-1 bg-primary-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
              {{ cartStore.totalItems }}
            </span>
          </router-link>

          <!-- Auth -->
          <div v-if="authStore.isAuthenticated" class="flex items-center space-x-3">
            <span class="text-sm text-gray-600">{{ authStore.user?.name }}</span>
            <button @click="logout" class="btn btn-secondary text-sm">
              Logout
            </button>
          </div>
          <router-link v-else to="/auth" class="btn btn-primary">
            Login
          </router-link>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useAuthStore } from '../stores/auth'
import { useCartStore } from '../stores/cart'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const cartStore = useCartStore()
const router = useRouter()

function logout() {
  authStore.logout()
  router.push('/')
}
</script>
