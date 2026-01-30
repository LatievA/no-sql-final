import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useToastStore = defineStore('toast', () => {
    const message = ref('')
    const type = ref('info') // success, error, info, warning
    const show = ref(false)

    function showToast(msg, toastType = 'info') {
        message.value = msg
        type.value = toastType
        show.value = true

        setTimeout(() => {
            show.value = false
        }, 3000)
    }

    function success(msg) {
        showToast(msg, 'success')
    }

    function error(msg) {
        showToast(msg, 'error')
    }

    function info(msg) {
        showToast(msg, 'info')
    }

    function warning(msg) {
        showToast(msg, 'warning')
    }

    return {
        message,
        type,
        show,
        showToast,
        success,
        error,
        info,
        warning
    }
})
