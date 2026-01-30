import api from './index'

export const authApi = {
    login(credentials) {
        return api.post('/auth/login', credentials)
    },

    register(data) {
        return api.post('/auth/register', data)
    },

    getMe() {
        return api.get('/auth/me')
    }
}

export const categoryApi = {
    getAll() {
        return api.get('/categories')
    },

    getById(id) {
        return api.get(`/categories/${id}`)
    },

    create(data) {
        return api.post('/categories', data)
    },

    update(id, data) {
        return api.put(`/categories/${id}`, data)
    },

    delete(id) {
        return api.delete(`/categories/${id}`)
    }
}

export const bicycleApi = {
    getAll(params = {}) {
        return api.get('/bicycles', { params })
    },

    getById(id) {
        return api.get(`/bicycles/${id}`)
    },

    create(data) {
        return api.post('/bicycles', data)
    },

    update(id, data) {
        return api.put(`/bicycles/${id}`, data)
    },

    delete(id) {
        return api.delete(`/bicycles/${id}`)
    },

    addReview(id, data) {
        return api.post(`/bicycles/${id}/reviews`, data)
    },

    updateStock(id, quantity) {
        return api.patch(`/bicycles/${id}/stock`, { quantity })
    }
}

export const orderApi = {
    getAll(params = {}) {
        return api.get('/orders', { params })
    },

    getMyOrders(params = {}) {
        return api.get('/orders/my', { params })
    },

    getById(id) {
        return api.get(`/orders/${id}`)
    },

    create(data) {
        return api.post('/orders', data)
    },

    updateStatus(id, status) {
        return api.patch(`/orders/${id}/status`, { status })
    }
}

export const customerApi = {
    getAll() {
        return api.get('/customers')
    },

    getById(id) {
        return api.get(`/customers/${id}`)
    },

    updateProfile(data) {
        return api.put('/customers/profile', data)
    },

    addAddress(data) {
        return api.post('/customers/addresses', data)
    },

    removeAddress(type) {
        return api.delete(`/customers/addresses/${type}`)
    }
}

export const reportApi = {
    getSalesByCategory() {
        return api.get('/reports/sales-by-category')
    },

    getTopSelling(limit = 10) {
        return api.get('/reports/top-selling', { params: { limit } })
    }
}
