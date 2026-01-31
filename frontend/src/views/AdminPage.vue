<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <h1 class="text-3xl font-bold text-gray-900 mb-8">Admin Dashboard</h1>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
      <div class="bg-white rounded-xl shadow-md p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-500">Total Bicycles</p>
            <p class="text-2xl font-bold">{{ stats.total_bicycles || 0 }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl shadow-md p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-green-100 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-500">Total Orders</p>
            <p class="text-2xl font-bold">{{ stats.total_orders || 0 }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl shadow-md p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-yellow-100 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-500">Total Revenue</p>
            <p class="text-2xl font-bold">{{ formatPrice(stats.total_revenue || 0) }}</p>
          </div>
        </div>
      </div>
      <div class="bg-white rounded-xl shadow-md p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-purple-100 rounded-lg flex items-center justify-center">
            <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-500">Customers</p>
            <p class="text-2xl font-bold">{{ stats.total_customers || 0 }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabs -->
    <div class="border-b mb-6">
      <nav class="flex space-x-8">
        <button 
          v-for="tab in tabs" 
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="activeTab === tab.id ? 'border-primary-500 text-primary-600' : 'border-transparent text-gray-500 hover:text-gray-700'"
          class="py-4 px-1 border-b-2 font-medium"
        >
          {{ tab.label }}
        </button>
      </nav>
    </div>

    <!-- Categories Tab -->
    <div v-if="activeTab === 'categories'">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold">Categories</h2>
        <button @click="openCategoryModal()" class="btn btn-primary">Add Category</button>
      </div>

      <div class="bg-white rounded-xl shadow-md overflow-hidden">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="cat in categories" :key="cat.id">
              <td class="px-6 py-4 font-medium">{{ cat.category_name }}</td>
              <td class="px-6 py-4 text-gray-500">{{ cat.description }}</td>
              <td class="px-6 py-4">
                <button @click="openCategoryModal(cat)" class="text-primary-600 hover:underline mr-4">Edit</button>
                <button @click="deleteCategory(cat.id)" class="text-red-600 hover:underline">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Bicycles Tab -->
    <div v-if="activeTab === 'bicycles'">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-xl font-bold">Bicycles</h2>
        <button @click="openBicycleModal()" class="btn btn-primary">Add Bicycle</button>
      </div>

      <div class="bg-white rounded-xl shadow-md overflow-hidden">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Model</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Brand</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Price</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Stock</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="bike in bicycles" :key="bike.id">
              <td class="px-6 py-4 font-medium">{{ bike.model_name }}</td>
              <td class="px-6 py-4">{{ bike.brand }}</td>
              <td class="px-6 py-4">{{ formatPrice(bike.price) }}</td>
              <td class="px-6 py-4">
                <span :class="bike.stock_quantity > 0 ? 'text-green-600' : 'text-red-600'">
                  {{ bike.stock_quantity }}
                </span>
              </td>
              <td class="px-6 py-4">
                <button @click="openBicycleModal(bike)" class="text-primary-600 hover:underline mr-4">Edit</button>
                <button @click="deleteBicycle(bike.id)" class="text-red-600 hover:underline">Delete</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Orders Tab -->
    <div v-if="activeTab === 'orders'">
      <h2 class="text-xl font-bold mb-6">All Orders</h2>
      
      <div class="bg-white rounded-xl shadow-md overflow-hidden">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Order ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Customer</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Total</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-200">
            <tr v-for="order in orders" :key="order.id">
              <td class="px-6 py-4 font-mono text-sm">{{ order.id.slice(0, 8) }}...</td>
              <td class="px-6 py-4">{{ order.customer_name || 'N/A' }}</td>
              <td class="px-6 py-4 font-medium">{{ formatPrice(order.total_amount) }}</td>
              <td class="px-6 py-4">
                <span :class="getStatusClass(order.status)" class="badge">{{ order.status }}</span>
              </td>
              <td class="px-6 py-4 text-gray-500">{{ formatDate(order.order_date) }}</td>
              <td class="px-6 py-4">
                <select 
                  :value="order.status" 
                  @change="updateOrderStatus(order.id, $event.target.value)"
                  class="input py-1 px-2 text-sm"
                >
                  <option value="pending">Pending</option>
                  <option value="confirmed">Confirmed</option>
                  <option value="shipped">Shipped</option>
                  <option value="delivered">Delivered</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Reports Tab -->
    <div v-if="activeTab === 'reports'">
      <h2 class="text-xl font-bold mb-6">Sales Reports</h2>

      <!-- Sales by Category -->
      <div class="bg-white rounded-xl shadow-md p-6 mb-6">
        <h3 class="font-semibold mb-4">Sales by Category</h3>
        <div class="space-y-3">
          <div v-for="item in salesByCategory" :key="item.category_id" class="flex items-center">
            <span class="w-40 font-medium">{{ item.category_name || 'Unknown' }}</span>
            <div class="flex-1 h-4 bg-gray-200 rounded-full overflow-hidden">
              <div 
                class="h-full bg-primary-500" 
                :style="{ width: `${(item.total_sales / maxRevenue) * 100}%` }"
              ></div>
            </div>
            <span class="w-32 text-right font-medium">{{ formatPrice(item.total_sales) }}</span>
            <span class="w-20 text-right text-gray-500">{{ item.total_items }} sold</span>
          </div>
        </div>
      </div>

      <!-- Top Selling Bicycles -->
      <div class="bg-white rounded-xl shadow-md p-6">
        <h3 class="font-semibold mb-4">Top Selling Bicycles</h3>
        <table class="w-full">
          <thead>
            <tr class="text-left text-gray-500 text-sm">
              <th class="pb-3">#</th>
              <th class="pb-3">Model</th>
              <th class="pb-3">Units Sold</th>
              <th class="pb-3">Revenue</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in topSelling" :key="item._id">
              <td class="py-2">{{ index + 1 }}</td>
              <td class="py-2 font-medium">{{ item.model_name }}</td>
              <td class="py-2">{{ item.total_sold }}</td>
              <td class="py-2 font-medium text-primary-600">{{ formatPrice(item.total_sales) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Category Modal -->
    <div v-if="showCategoryModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-xl p-6 w-full max-w-md">
        <h3 class="text-xl font-bold mb-4">{{ editingCategory ? 'Edit' : 'Add' }} Category</h3>
        <form @submit.prevent="saveCategory" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input v-model="categoryForm.category_name" type="text" required class="input" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea v-model="categoryForm.description" rows="3" class="input"></textarea>
          </div>
          <div class="flex justify-end space-x-3">
            <button type="button" @click="showCategoryModal = false" class="btn btn-secondary">Cancel</button>
            <button type="submit" class="btn btn-primary">Save</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Bicycle Modal -->
    <div v-if="showBicycleModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 overflow-y-auto">
      <div class="bg-white rounded-xl p-6 w-full max-w-2xl m-4">
        <h3 class="text-xl font-bold mb-4">{{ editingBicycle ? 'Edit' : 'Add' }} Bicycle</h3>
        <form @submit.prevent="saveBicycle" class="space-y-4 max-h-[70vh] overflow-y-auto">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Model Name</label>
              <input v-model="bicycleForm.model_name" type="text" required class="input" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Brand</label>
              <input v-model="bicycleForm.brand" type="text" required class="input" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Price (KZT)</label>
              <input v-model.number="bicycleForm.price" type="number" required class="input" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Stock Quantity</label>
              <input v-model.number="bicycleForm.stock_quantity" type="number" required class="input" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
              <select v-model="bicycleForm.category_id" required class="input">
                <option value="">Select category</option>
                <option v-for="cat in categories" :key="cat.id" :value="cat.id">{{ cat.category_name }}</option>
              </select>
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
            <textarea v-model="bicycleForm.description" rows="3" class="input"></textarea>
          </div>
          <div class="border-t pt-4">
            <h4 class="font-medium mb-3">Specifications</h4>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm text-gray-700 mb-1">Frame Material</label>
                <input v-model="bicycleForm.specifications.frame_material" class="input" />
              </div>
              <div>
                <label class="block text-sm text-gray-700 mb-1">Wheel Size</label>
                <input v-model="bicycleForm.specifications.wheel_size" class="input" />
              </div>
              <div>
                <label class="block text-sm text-gray-700 mb-1">Gear Count</label>
                <input v-model.number="bicycleForm.specifications.gear_count" type="number" class="input" />
              </div>
              <div>
                <label class="block text-sm text-gray-700 mb-1">Brake Type</label>
                <input v-model="bicycleForm.specifications.brake_type" class="input" />
              </div>
            </div>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="showBicycleModal = false" class="btn btn-secondary">Cancel</button>
            <button type="submit" class="btn btn-primary">Save</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { categoryApi, bicycleApi, orderApi, reportApi, customerApi } from '../api/endpoints'
import { useToastStore } from '../stores/toast'

const toastStore = useToastStore()

const tabs = [
  { id: 'categories', label: 'Categories' },
  { id: 'bicycles', label: 'Bicycles' },
  { id: 'orders', label: 'Orders' },
  { id: 'reports', label: 'Reports' }
]

const activeTab = ref('categories')
const categories = ref([])
const bicycles = ref([])
const orders = ref([])
const stats = ref({})
const salesByCategory = ref([])
const topSelling = ref([])

// Category Modal
const showCategoryModal = ref(false)
const editingCategory = ref(null)
const categoryForm = reactive({
  category_name: '',
  description: ''
})

// Bicycle Modal
const showBicycleModal = ref(false)
const editingBicycle = ref(null)
const bicycleForm = reactive({
  model_name: '',
  brand: '',
  price: 0,
  stock_quantity: 0,
  category_id: '',
  description: '',
  specifications: {
    frame_material: '',
    wheel_size: '',
    gear_count: 0,
    brake_type: ''
  }
})

const maxRevenue = computed(() => {
  return Math.max(...salesByCategory.value.map(s => s.total_sales || 0), 1)
})

async function fetchData() {
  try {
    const [catRes, bikeRes, orderRes, statsRes, salesRes, topRes, custRes] = await Promise.all([
      categoryApi.getAll(),
      bicycleApi.getAll({ limit: 100 }),
      orderApi.getAll({ limit: 100 }),
      reportApi.getSalesSummary(),
      reportApi.getSalesByCategory(),
      reportApi.getTopSelling(5),
      customerApi.getAll()
    ])
    
    categories.value = catRes.data.data || []
    bicycles.value = bikeRes.data.data || []
    orders.value = orderRes.data.data || []
    
    stats.value = {
      ...statsRes.data.data,
      total_bicycles: bicycles.value.length,
      total_customers: custRes.data.data?.length || 0
    }
    
    salesByCategory.value = salesRes.data.data || []
    topSelling.value = topRes.data.data || []
  } catch (error) {
    console.error('Failed to fetch data:', error)
  }
}

// Category functions
function openCategoryModal(cat = null) {
  editingCategory.value = cat
  if (cat) {
    categoryForm.category_name = cat.category_name
    categoryForm.description = cat.description
  } else {
    categoryForm.category_name = ''
    categoryForm.description = ''
  }
  showCategoryModal.value = true
}

async function saveCategory() {
  try {
    if (editingCategory.value) {
      await categoryApi.update(editingCategory.value.id, categoryForm)
      toastStore.success('Category updated')
    } else {
      await categoryApi.create(categoryForm)
      toastStore.success('Category created')
    }
    showCategoryModal.value = false
    fetchData()
  } catch (error) {
    toastStore.error('Failed to save category')
  }
}

async function deleteCategory(id) {
  if (!confirm('Delete this category?')) return
  try {
    await categoryApi.delete(id)
    toastStore.success('Category deleted')
    fetchData()
  } catch (error) {
    toastStore.error('Failed to delete category')
  }
}

// Bicycle functions
function openBicycleModal(bike = null) {
  editingBicycle.value = bike
  if (bike) {
    Object.assign(bicycleForm, {
      model_name: bike.model_name,
      brand: bike.brand,
      price: bike.price,
      stock_quantity: bike.stock_quantity,
      category_id: bike.category_id,
      description: bike.description,
      specifications: { ...bike.specifications }
    })
  } else {
    Object.assign(bicycleForm, {
      model_name: '',
      brand: '',
      price: 0,
      stock_quantity: 0,
      category_id: '',
      description: '',
      specifications: {
        frame_material: '',
        wheel_size: '',
        gear_count: 0,
        brake_type: ''
      }
    })
  }
  showBicycleModal.value = true
}

async function saveBicycle() {
  try {
    if (editingBicycle.value) {
      await bicycleApi.update(editingBicycle.value.id, bicycleForm)
      toastStore.success('Bicycle updated')
    } else {
      await bicycleApi.create(bicycleForm)
      toastStore.success('Bicycle created')
    }
    showBicycleModal.value = false
    fetchData()
  } catch (error) {
    toastStore.error('Failed to save bicycle')
  }
}

async function deleteBicycle(id) {
  if (!confirm('Delete this bicycle?')) return
  try {
    await bicycleApi.delete(id)
    toastStore.success('Bicycle deleted')
    fetchData()
  } catch (error) {
    toastStore.error('Failed to delete bicycle')
  }
}

// Order functions
async function updateOrderStatus(orderId, status) {
  try {
    await orderApi.updateStatus(orderId, status)
    toastStore.success('Order status updated')
    fetchData()
  } catch (error) {
    toastStore.error('Failed to update order')
  }
}

function formatPrice(price) {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'KZT',
    minimumFractionDigits: 0
  }).format(price)
}

function formatDate(date) {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

function getStatusClass(status) {
  const classes = {
    pending: 'badge-warning',
    confirmed: 'badge-info',
    shipped: 'badge-info',
    delivered: 'badge-success',
    cancelled: 'badge-danger'
  }
  return classes[status] || 'badge-info'
}

onMounted(() => {
  fetchData()
})
</script>
