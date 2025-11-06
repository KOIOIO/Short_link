<template>
  <div>
    <div v-if="!isAuthPage" class="topbar glass">
      <router-link to="/" class="brand">短链接生成器</router-link>
      <div class="menu" v-if="!isAuthed">
        <router-link to="/login">登录</router-link>
        <router-link to="/register">注册</router-link>
      </div>
    </div>
    <router-view />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { isAuthed } from './api/auth'

const route = useRoute()
const isAuthPage = computed(() => route.path === '/login' || route.path === '/register')
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Montserrat:wght@700&display=swap');

@keyframes fadeIn {
  0% { opacity: 0; transform: translateY(30px) scale(0.95); }
  60% { opacity: 1; transform: translateY(-8px) scale(1.03); }
  100% { opacity: 1; transform: translateY(0) scale(1); }
}

@keyframes bounceIn {
  0% { opacity: 0; transform: scale(0.7); }
  60% { opacity: 1; transform: scale(1.1); }
  100% { opacity: 1; transform: scale(1); }
}

.fade-in {
  animation: fadeIn 1s cubic-bezier(.68,-0.55,.27,1.55);
}

.bounce-enter-active {
  animation: bounceIn 0.7s cubic-bezier(.68,-0.55,.27,1.55);
}
.bounce-leave-active {
  animation: fadeIn 0.3s reverse;
}

.bg-animated {
  position: relative;
  min-height: 100vh;
  width: 100vw;
  overflow: hidden;
  background: linear-gradient(120deg, #f9d423 0%, #ff4e50 100%);
}

.topbar { display:flex; align-items:center; justify-content:space-between; padding: 12px 18px; background: #fff; border-bottom: 1px solid #eee; }
.brand { font-weight: bold; color: #ff4e50; text-decoration: none; }
.menu a { margin-left: 10px; color: #333; text-decoration: none; }

h1 {
  text-align: center;
  color: #fff;
  margin-bottom: 32px;
  letter-spacing: 2px;
  font-size: 2.2rem;
  font-family: 'Montserrat', sans-serif;
}

.highlight {
  background: linear-gradient(90deg, #f9d423 30%, #ff4e50 70%);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: bold;
}

.alt-actions {
  display: flex;
  gap: 8px;
  margin-top: 12px;
}

.small {
  flex: 1;
  padding: 8px 10px;
  font-size: 13px;
  border-radius: 8px;
  background: rgba(255,255,255,0.06);
  color: #fff;
  border: 1px solid rgba(255,255,255,0.06);
}

.form-group {
  margin-bottom: 22px;
}

label {
  color: #ff4e50;
  font-weight: bold;
  margin-bottom: 6px;
  display: block;
  letter-spacing: 1px;
}

input, select {
  width: 100%;
  padding: 12px;
  border: 2px solid #f9d423;
  border-radius: 8px;
  box-sizing: border-box;
  font-size: 15px;
  background: #fffbe6;
  transition: border-color 0.3s, box-shadow 0.3s;
  outline: none;
}

input[type="text"] {
  background: #222;
  color: #fff;
  border: 2px solid #ff4e50;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

input[type="text"]:focus {
  border-color: #f9d423;
  box-shadow: 0 0 0 2px #ffb19955;
}

select {
  background: #222;
  color: #fff;
  border: 2px solid #ff4e50;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
}

select:focus {
  border-color: #f9d423;
  box-shadow: 0 0 0 2px #ffb19955;
}

button {
  background: linear-gradient(90deg, #ff4e50 0%, #f9d423 100%);
  color: #fff;
  border: none;
  padding: 14px;
  width: 100%;
  font-size: 18px;
  border-radius: 10px;
  cursor: pointer;
  font-weight: bold;
  letter-spacing: 1px;
  box-shadow: 0 4px 16px rgba(255, 78, 80, 0.12);
  transition: background 0.3s, transform 0.2s;
}

button:disabled {
  background: #ffd6d6;
  color: #fff;
  cursor: not-allowed;
  opacity: 0.7;
}

button:hover:enabled {
  background: linear-gradient(90deg, #f9d423 0%, #ff4e50 100%);
  transform: translateY(-2px) scale(1.03);
}

.result {
  margin-top: 28px;
  background: linear-gradient(90deg, #43e97b 0%, #38f9d7 100%);
  padding: 18px;
  border-left: 6px solid #43e97b;
  color: #22543d;
  border-radius: 10px;
  font-size: 1.1rem;
  box-shadow: 0 2px 8px rgba(67, 233, 123, 0.08);
  animation: bounceIn 0.7s;
}

.result a {
  color: #ff4e50;
  font-weight: bold;
  word-break: break-all;
  text-decoration: underline;
  transition: color 0.2s;
}

.result a:hover {
  color: #f9d423;
}

/* 科技感增强 */
.glass { background: rgba(17, 25, 40, 0.45); backdrop-filter: blur(12px); }
.brand { background: linear-gradient(90deg, #63a4ff 0%, #83eaf1 100%); -webkit-background-clip: text; background-clip: text; -webkit-text-fill-color: transparent; }
.menu a { margin-left: 12px; color: #e6edf7; text-decoration: none; padding: 6px 10px; border-radius: 8px; transition: transform .2s, background .3s; }
.menu a:hover { transform: translateY(-2px); background: rgba(255,255,255,0.08); }
</style>