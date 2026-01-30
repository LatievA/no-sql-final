import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
    {
        path: '/',
        name: 'Home',
        component: () => import('../views/HomePage.vue')
    },
    {
        path: '/bicycles/:id',
        name: 'BicycleDetail',
        component: () => import('../views/BicycleDetailPage.vue')
    },
    {
        path: '/cart',
        name: 'Cart',
        component: () => import('../views/CartPage.vue')
    },
    {
        path: '/orders',
        name: 'Orders',
        component: () => import('../views/OrdersPage.vue'),
        meta: { requiresAuth: true }
    },
    {
        path: '/admin',
        name: 'Admin',
        component: () => import('../views/AdminPage.vue'),
        meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
        path: '/auth',
        name: 'Auth',
        component: () => import('../views/AuthPage.vue'),
        meta: { guestOnly: true }
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

router.beforeEach((to, from, next) => {
    const authStore = useAuthStore()

    if (to.meta.requiresAuth && !authStore.isAuthenticated) {
        next('/auth')
    } else if (to.meta.requiresAdmin && !authStore.isAdmin) {
        next('/')
    } else if (to.meta.guestOnly && authStore.isAuthenticated) {
        next('/')
    } else {
        next()
    }
})

export default router
