import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
    headers: {
        'Content-Type': 'application/json'
    }
})

// Request interceptor to add auth token
api.interceptors.request.use(
    (config) => {
        const authStore = useAuthStore()
        if (authStore.token) {
            config.headers.Authorization = `Bearer ${authStore.token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Response interceptor to handle errors
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            const authStore = useAuthStore()
            authStore.logout()
            window.location.href = '/auth'
        }
        return Promise.reject(error)
    }
)

export default api
