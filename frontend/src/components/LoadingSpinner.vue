<template>
  <div v-if="loading" class="flex items-center justify-center py-12">
    <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
  </div>
</template>

<script setup>
defineProps({
  loading: {
    type: Boolean,
    default: true
  }
})
</script>
