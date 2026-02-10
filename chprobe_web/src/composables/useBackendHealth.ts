import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

export const useBackendHealth = () => {
  const router = useRouter()
  const isBackendDown = ref(false)
  const healthCheckInterval = ref<number | null>(null)
  const HEALTH_CHECK_INTERVAL = 10000

  const checkHealth = async () => {
    try {
      const baseURL = import.meta.env.VITE_API_BASE_URL || '/api'
      console.log('Checking health at:', `${baseURL}/health`)
      const response = await axios.get(`${baseURL}/health`, {
        timeout: 2000
      })
      
      console.log('Health check response:', response.status)
      
      if (response.status === 200) {
        console.log('Backend is up, clearing interval and redirecting to login')
        isBackendDown.value = false
        if (healthCheckInterval.value) {
          clearInterval(healthCheckInterval.value)
          healthCheckInterval.value = null
        }
        
        localStorage.removeItem('token')
        localStorage.removeItem('session')
        router.push('/login')
      }
    } catch (error) {
      console.log('Health check failed:', error)
      isBackendDown.value = true
    }
  }

  const startHealthCheck = () => {
    if (healthCheckInterval.value) {
      clearInterval(healthCheckInterval.value)
    }
    
    console.log('Starting health check')
    isBackendDown.value = true
    
    healthCheckInterval.value = window.setInterval(() => {
      checkHealth()
    }, HEALTH_CHECK_INTERVAL)
    
    checkHealth()
  }

  const stopHealthCheck = () => {
    if (healthCheckInterval.value) {
      clearInterval(healthCheckInterval.value)
      healthCheckInterval.value = null
    }
    isBackendDown.value = false
  }

  onMounted(() => {
    window.addEventListener('backend-down', startHealthCheck)
    window.addEventListener('backend-up', stopHealthCheck)
  })

  onUnmounted(() => {
    stopHealthCheck()
    window.removeEventListener('backend-down', startHealthCheck)
    window.removeEventListener('backend-up', stopHealthCheck)
  })

  return {
    isBackendDown,
    startHealthCheck,
    stopHealthCheck
  }
}
