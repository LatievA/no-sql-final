import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useCartStore = defineStore('cart', () => {
    const items = ref(JSON.parse(localStorage.getItem('cart') || '[]'))

    const totalItems = computed(() => {
        return items.value.reduce((sum, item) => sum + item.quantity, 0)
    })

    const totalAmount = computed(() => {
        return items.value.reduce((sum, item) => sum + (item.price * item.quantity), 0)
    })

    function saveToStorage() {
        localStorage.setItem('cart', JSON.stringify(items.value))
    }

    function addItem(bicycle, quantity = 1, customizations = []) {
        const existingIndex = items.value.findIndex(item =>
            item.bicycle_id === bicycle.id &&
            JSON.stringify(item.selected_customizations) === JSON.stringify(customizations)
        )

        if (existingIndex > -1) {
            items.value[existingIndex].quantity += quantity
        } else {
            items.value.push({
                bicycle_id: bicycle.id,
                model_name: bicycle.model_name,
                brand: bicycle.brand,
                price: bicycle.price,
                image_url: bicycle.image_url,
                quantity,
                selected_customizations: customizations,
                stock_quantity: bicycle.stock_quantity
            })
        }
        saveToStorage()
    }

    function updateQuantity(index, quantity) {
        if (quantity <= 0) {
            removeItem(index)
        } else {
            items.value[index].quantity = quantity
            saveToStorage()
        }
    }

    function removeItem(index) {
        items.value.splice(index, 1)
        saveToStorage()
    }

    function clearCart() {
        items.value = []
        saveToStorage()
    }

    function getOrderItems() {
        return items.value.map(item => ({
            bicycle_id: item.bicycle_id,
            quantity: item.quantity,
            selected_customizations: item.selected_customizations
        }))
    }

    return {
        items,
        totalItems,
        totalAmount,
        addItem,
        updateQuantity,
        removeItem,
        clearCart,
        getOrderItems
    }
})
