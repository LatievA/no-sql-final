<template>
  <Transition
    enter-active-class="transition ease-out duration-300"
    enter-from-class="opacity-0 translate-y-2"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition ease-in duration-200"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-2"
  >
    <div v-if="toastStore.show" class="fixed bottom-4 right-4 z-50">
      <div :class="toastClasses" class="px-6 py-4 rounded-lg shadow-lg flex items-center space-x-3">
        <!-- Icon -->
        <svg v-if="toastStore.type === 'success'" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
        <svg v-else-if="toastStore.type === 'error'" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
        <svg v-else class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        
        <span class="font-medium">{{ toastStore.message }}</span>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { computed } from 'vue'
import { useToastStore } from '../stores/toast'

const toastStore = useToastStore()

const toastClasses = computed(() => {
  const base = 'text-white'
  switch (toastStore.type) {
    case 'success':
      return `${base} bg-green-500`
    case 'error':
      return `${base} bg-red-500`
    case 'warning':
      return `${base} bg-yellow-500`
    default:
      return `${base} bg-blue-500`
  }
})
</script>
