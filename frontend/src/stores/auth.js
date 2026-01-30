import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api/endpoints'

export const useAuthStore = defineStore('auth', () => {
    const user = ref(null)
    const token = ref(localStorage.getItem('token'))

    const isAuthenticated = computed(() => !!token.value)
    const isAdmin = computed(() => user.value?.role === 'admin')

    async function login(credentials) {
        const response = await authApi.login(credentials)
        if (response.data.success) {
            token.value = response.data.token
            user.value = response.data.user
            localStorage.setItem('token', response.data.token)
        }
        return response.data
    }

    async function register(data) {
        const response = await authApi.register(data)
        if (response.data.success) {
            token.value = response.data.token
            user.value = response.data.user
            localStorage.setItem('token', response.data.token)
        }
        return response.data
    }

    async function fetchUser() {
        if (!token.value) return
        try {
            const response = await authApi.getMe()
            if (response.data.success) {
                user.value = response.data.data
            }
        } catch (error) {
            logout()
        }
    }

    function logout() {
        user.value = null
        token.value = null
        localStorage.removeItem('token')
    }

    // Initialize user data if token exists
    if (token.value) {
        fetchUser()
    }

    return {
        user,
        token,
        isAuthenticated,
        isAdmin,
        login,
        register,
        fetchUser,
        logout
    }
})
