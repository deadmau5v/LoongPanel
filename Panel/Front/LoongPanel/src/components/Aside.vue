<template>
  <el-aside>
    <el-menu>
      <el-image :src="logo"></el-image>
      <router-link to="/">
        <el-menu-item index="1">首页</el-menu-item>
      </router-link>
      <el-menu-item index="2">状态</el-menu-item>
      <el-menu-item index="3">用户</el-menu-item>
      <el-menu-item index="4">设置</el-menu-item>
      <el-menu-item index="5">应用</el-menu-item>
      <router-link to="/files">
        <el-menu-item index="6">文件</el-menu-item>
      </router-link>
      <el-menu-item index="7">终端</el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script setup>
import logo from '/src/assets/logo.png'
</script>

<style scoped>
/* 侧边栏整体 */
.el-aside {
  position: fixed;
  top: 0;
  bottom: 0;
  left: 0;
  height: 100%;
  background-image: url("/src/assets/aside.png");
  background-size: auto 100%;
  background-position: center;
}

/* 侧边栏背景 */
.el-menu {
  background: none;
}

/* 侧边栏Logo */
.el-image {
  width: 50%;
  height: auto;
  display: block;
  margin: 10% auto;
}

/* 侧边栏选项 */
.el-menu-item {
  width: 100%;
  justify-content: center;
  border-radius: 0;
  color: white;
  font-size: 20px;
  font-weight: bold;
  text-align: center;
}

.el-menu-item:hover {
  color: grey;
}

.el-menu-item.is-active {
  color: red;
}

/* 移除 <router-link> 的下划线 */
a {
  text-decoration: none;
}

</style>