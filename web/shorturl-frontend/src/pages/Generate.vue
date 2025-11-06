<template>
  <div class="container">
    <div class="nav">
      <router-link to="/" class="brand">短链接生成器</router-link>
      <div class="spacer" />
      <button class="link" @click="logout">退出登录</button>
    </div>

    <form @submit.prevent="generateShortURL" class="form">
      <div class="form-group">
        <label>原始链接:</label>
        <input v-model="url" type="text" required placeholder="请输入原始链接..." />
      </div>
      <div class="form-group">
        <label>过期时间:</label>
        <select v-model="expiration">
          <option value="30m">30 分钟</option>
          <option value="1h">1 小时</option>
          <option value="1d">1 天</option>
        </select>
      </div>
      <button type="submit" :disabled="loading">
        {{ loading ? '生成中...' : '✨ 生成短链接' }}
      </button>
      <div class="alt-actions">
        <button type="button" class="small" @click="postTo('generatebymysnowflake')" :disabled="loading">生成（雪花ID）</button>
        <button type="button" class="small" @click="postTo('filterbymybloomfilter')" :disabled="loading">生成（我的布隆）</button>
        <button type="button" class="small" @click="postTo('generatewithiplimiter')" :disabled="loading">生成（IP 限流）</button>
      </div>
    </form>

    <div v-if="shortUrl" class="result">
      <p>生成成功：<a :href="shortUrl" target="_blank">{{ shortUrl }}</a></p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import http from '../api/http'
import { clearToken } from '../api/auth'

const url = ref('')
const expiration = ref('1h')
const shortUrl = ref('')
const loading = ref(false)

const logout = () => {
  clearToken()
  location.href = '/login'
}

const generateShortURL = async () => {
  loading.value = true
  try {
    const { data } = await http.post('/shorturls/generate', {
      url: url.value,
      expiration: expiration.value,
    })
    shortUrl.value = `${http.defaults.baseURL}/shorturls/${data.short_url}`
  } catch (error) {
    handleError(error)
  } finally {
    loading.value = false
  }
}

const postTo = async (path) => {
  loading.value = true
  try {
    const { data } = await http.post(`/shorturls/${path}`, {
      url: url.value,
      expiration: expiration.value,
    })
    shortUrl.value = `${http.defaults.baseURL}/shorturls/${data.short_url}`
  } catch (error) {
    handleError(error)
  } finally {
    loading.value = false
  }
}

const handleError = (error) => {
  let isRateLimited = false
  if (error.response) {
    const data = error.response.data
    if (data && String(data.code) === '5001') {
      alert('为防止数据库崩溃，请半个小时后再生成')
      isRateLimited = true
    } else if (error.response.status === 500 && error.response.statusText === 'Internal Server Error') {
      alert('为防止数据库崩溃，请半个小时后再生成')
      isRateLimited = true
    } else if (typeof data === 'string' && data.includes('rate limit')) {
      alert('为防止数据库崩溃，请半个小时后再生成')
      isRateLimited = true
    } else if (error.response.status === 401) {
      alert('登录状态已过期，请重新登录')
      location.href = '/login'
      return
    }
  }
  if (!isRateLimited) {
    alert('生成失败，请检查输入！')
  }
}
</script>

<style scoped>
.container { max-width: 560px; margin: 40px auto; padding: 24px; background: linear-gradient(135deg, rgba(17,25,40,.88), rgba(31,41,55,.88)); color: #e6edf7; border-radius: 14px; box-shadow: 0 18px 50px rgba(0,0,0,0.25); border: 1px solid rgba(255,255,255,.14); animation: fadeIn .5s ease; }
.nav { display: flex; align-items: center; margin-bottom: 12px; }
.brand { font-weight: bold; background: linear-gradient(90deg, #63a4ff 0%, #83eaf1 100%); -webkit-background-clip: text; background-clip: text; -webkit-text-fill-color: transparent; text-decoration: none; }
.spacer { flex: 1; }
.link { background: rgba(255,255,255,.08); border: 1px solid rgba(255,255,255,.18); color: #e6edf7; cursor: pointer; padding: 8px 12px; border-radius: 10px; }
.form-group { margin-bottom: 14px; }
label { display: block; margin-bottom: 6px; color: #c9d4e8; }
input, select { width: 100%; padding: 12px; border: 1px solid rgba(255,255,255,.18); border-radius: 10px; background: rgba(255,255,255,.08); color: #e6edf7; }
button { padding: 12px 16px; background: linear-gradient(90deg, #4776E6 0%, #8E54E9 100%); border: none; color: #fff; border-radius: 10px; cursor: pointer; box-shadow: 0 12px 24px rgba(71,118,230,0.35); transition: transform .2s, box-shadow .2s; }
button:hover { transform: translateY(-2px); box-shadow: 0 16px 30px rgba(71,118,230,0.45); }
.alt-actions { display: flex; gap: 8px; margin-top: 10px; }
.small { background: linear-gradient(90deg, #f59f00 0%, #f9d423 100%); border: none; color: #222; }
.result { margin-top: 16px; }

@keyframes fadeIn { from { opacity: 0; transform: translateY(6px); } to { opacity: 1; transform: translateY(0); } }
</style>