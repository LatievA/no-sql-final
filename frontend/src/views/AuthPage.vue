<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100 py-12 px-4">
    <div class="max-w-md w-full">
      <div class="text-center mb-8">
        <svg class="w-16 h-16 text-primary-600 mx-auto" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
          <circle cx="25" cy="65" r="18" stroke="currentColor" stroke-width="4"/>
          <circle cx="75" cy="65" r="18" stroke="currentColor" stroke-width="4"/>
          <path d="M25 65 L40 35 L60 35 L75 65" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M40 35 L50 65 L75 65" stroke="currentColor" stroke-width="4" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <h2 class="mt-4 text-3xl font-bold text-gray-900">
          {{ isLogin ? 'Welcome Back' : 'Create Account' }}
        </h2>
        <p class="text-gray-600">
          {{ isLogin ? 'Sign in to your account' : 'Join our bicycle community' }}
        </p>
      </div>

      <div class="bg-white rounded-xl shadow-md p-8">
        <!-- Toggle -->
        <div class="flex rounded-lg bg-gray-100 p-1 mb-6">
          <button 
            @click="isLogin = true" 
            :class="isLogin ? 'bg-white shadow' : 'text-gray-500'"
            class="flex-1 py-2 rounded-lg font-medium transition-all"
          >
            Login
          </button>
          <button 
            @click="isLogin = false" 
            :class="!isLogin ? 'bg-white shadow' : 'text-gray-500'"
            class="flex-1 py-2 rounded-lg font-medium transition-all"
          >
            Register
          </button>
        </div>

        <!-- Login Form -->
        <form v-if="isLogin" @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input 
              v-model="loginForm.email" 
              type="email" 
              required
              class="input"
              placeholder="you@example.com"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input 
              v-model="loginForm.password" 
              type="password" 
              required
              class="input"
              placeholder="••••••••"
            />
          </div>
          <button 
            type="submit" 
            :disabled="loading"
            class="w-full btn btn-primary py-3 disabled:opacity-50"
          >
            {{ loading ? 'Signing in...' : 'Sign In' }}
          </button>
        </form>

        <!-- Register Form -->
        <form v-else @submit.prevent="handleRegister" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
            <input 
              v-model="registerForm.name" 
              type="text" 
              required
              class="input"
              placeholder="John Doe"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input 
              v-model="registerForm.email" 
              type="email" 
              required
              class="input"
              placeholder="you@example.com"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
            <input 
              v-model="registerForm.phone" 
              type="tel" 
              class="input"
              placeholder="+7 777 123 4567"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
            <input 
              v-model="registerForm.password" 
              type="password" 
              required
              minlength="6"
              class="input"
              placeholder="••••••••"
            />
          </div>
          <button 
            type="submit" 
            :disabled="loading"
            class="w-full btn btn-primary py-3 disabled:opacity-50"
          >
            {{ loading ? 'Creating account...' : 'Create Account' }}
          </button>
        </form>

        <!-- Demo Accounts -->
        <div class="mt-6 pt-6 border-t">
          <p class="text-sm text-gray-500 text-center mb-3">Demo Accounts</p>
          <div class="grid grid-cols-2 gap-3 text-sm">
            <button 
              @click="fillDemo('admin')" 
              class="p-2 bg-gray-50 rounded-lg hover:bg-gray-100"
            >
              <span class="font-medium">Admin</span><br>
              <span class="text-gray-500 text-xs">admin@store.com</span>
            </button>
            <button 
              @click="fillDemo('customer')" 
              class="p-2 bg-gray-50 rounded-lg hover:bg-gray-100"
            >
              <span class="font-medium">Customer</span><br>
              <span class="text-gray-500 text-xs">customer@store.com</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useToastStore } from '../stores/toast'

const router = useRouter()
const authStore = useAuthStore()
const toastStore = useToastStore()

const isLogin = ref(true)
const loading = ref(false)

const loginForm = reactive({
  email: '',
  password: ''
})

const registerForm = reactive({
  name: '',
  email: '',
  phone: '',
  password: ''
})

async function handleLogin() {
  loading.value = true
  try {
    await authStore.login(loginForm)
    toastStore.success('Welcome back!')
    router.push('/')
  } catch (error) {
    const message = error.response?.data?.error || 'Login failed'
    toastStore.error(message)
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  loading.value = true
  try {
    await authStore.register(registerForm)
    toastStore.success('Account created successfully!')
    router.push('/')
  } catch (error) {
    const message = error.response?.data?.error || 'Registration failed'
    toastStore.error(message)
  } finally {
    loading.value = false
  }
}

function fillDemo(type) {
  isLogin.value = true
  if (type === 'admin') {
    loginForm.email = 'admin@store.com'
    loginForm.password = 'admin123'
  } else {
    loginForm.email = 'customer@store.com'
    loginForm.password = 'password123'
  }
}
</script>
