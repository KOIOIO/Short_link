<template>
  <div class="auth-container">
    <h2>注册</h2>
    <form @submit.prevent="onSubmit">
      <label>
        用户名
        <input v-model.trim="username" type="text" required placeholder="请输入用户名" />
      </label>
      <label>
        密码
        <input v-model="password" type="password" required placeholder="请输入密码" />
      </label>
      <label>
        确认密码
        <input v-model="confirmPassword" type="password" required placeholder="请再次输入密码" />
      </label>
      <button type="submit" :disabled="loading">{{ loading ? '注册中...' : '注册' }}</button>
      <p class="alt">
        已有账号？<router-link to="/login">去登录</router-link>
      </p>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import http from '../api/http'

const router = useRouter()
const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)

const onSubmit = async () => {
  if (password.value !== confirmPassword.value) {
    alert('两次输入的密码不一致')
    return
  }
  loading.value = true
  try {
    const { data } = await http.post('/users/register', {
      username: username.value,
      password: password.value,
      confirm_password: confirmPassword.value,
    })
    if (data && data.code === 200) {
      alert('注册成功，请登录')
      router.replace('/login')
    } else {
      alert(data?.message || '注册失败')
    }
  } catch (e) {
    alert('注册失败，请重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.auth-container {
  max-width: 420px;
  margin: 60px auto;
  padding: 24px;
  border-radius: 14px;
  border: 1px solid rgba(255,255,255,0.14);
  background: linear-gradient(135deg, rgba(17,25,40,.88), rgba(31,41,55,.88));
  color: #e6edf7;
  box-shadow: 0 18px 50px rgba(0,0,0,0.25);
  animation: fadeIn .5s ease;
}
label { display: block; margin-bottom: 14px; color: #c9d4e8; }
input { width: 100%; padding: 12px; margin-top: 6px; border: 1px solid rgba(255,255,255,0.18); border-radius: 10px; background: rgba(255,255,255,.08); color: #e6edf7; }
button { width: 100%; padding: 12px; background: linear-gradient(90deg, #4776E6 0%, #8E54E9 100%); border: none; color: #fff; border-radius: 10px; cursor: pointer; box-shadow: 0 12px 24px rgba(71,118,230,0.35); transition: transform .2s, box-shadow .2s; }
.button:hover { transform: translateY(-2px); box-shadow: 0 16px 30px rgba(71,118,230,0.45); }
.alt { text-align: center; margin-top: 12px; }

@keyframes fadeIn { from { opacity: 0; transform: translateY(6px); } to { opacity: 1; transform: translateY(0); } }
</style>