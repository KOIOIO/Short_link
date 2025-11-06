import { ref, computed } from 'vue'

// 全局认证状态
const tokenRef = ref(localStorage.getItem('token') || '')
export const isAuthed = computed(() => !!tokenRef.value)

export function setToken(t) {
  tokenRef.value = t || ''
  if (t) localStorage.setItem('token', t)
  else localStorage.removeItem('token')
}

export function clearToken() {
  setToken('')
}

// 监听跨标签的登录/退出变化
window.addEventListener('storage', (e) => {
  if (e.key === 'token') {
    tokenRef.value = e.newValue || ''
  }
})