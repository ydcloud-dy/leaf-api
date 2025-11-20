<template>
  <div class="main-layout">
    <Header />
    <main class="main-content">
      <router-view />
    </main>
    <footer class="footer">
      <div class="container">
        <p>&copy; 2024 个人博客. All rights reserved.</p>
        <p>Powered by Vue 3 + Element Plus</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import Header from '@/components/Header.vue'
</script>

<style scoped>
.main-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.main-content {
  flex: 1;
  padding: 20px 0;
}

.footer {
  background-color: #2c3e50;
  color: #fff;
  padding: 30px 0;
  text-align: center;
  margin-top: 60px;
}

.footer p {
  margin: 5px 0;
  opacity: 0.8;
}
</style>
