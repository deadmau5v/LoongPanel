<script setup>
import Aside from "./components/Aside.vue";
</script>

<template>
  <div class="common-layout">
    <el-container>
    <el-aside width="250px">
      <Aside/>
    </el-aside>
    <el-main>
      <router-view></router-view>
    </el-main>
  </el-container>
  </div>
</template>

<style scoped>
  body {
      margin: 0;
      padding: 0;
      font-family: 'Roboto', sans-serif;
      background-color: #f0f0f0;
      width: 100%;
      height: 100%;
  }



</style>