<template>
  <div class="min-h-screen bg-gray-50">
    <Navbar />
    <main>
      <router-view />
    </main>
    <Toast />
  </div>
</template>

<script setup>
import Navbar from './components/Navbar.vue'
import Toast from './components/Toast.vue'
</script>
