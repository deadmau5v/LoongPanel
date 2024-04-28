<template>
  <div id="terminal"></div>
  <el-button-group>
    <el-button @click="createScreen">创建</el-button>
    <el-button @click="connect">连接</el-button>
    <el-button @click="term.clear()">清屏</el-button>
  </el-button-group>
  <el-button-group id="screens">

  </el-button-group>
  <link rel="stylesheet" href="https://cdn1.d5v.cc/CDN/Project/LoongPanel/static/xterm.css"/>
</template>


<script setup>
import {Terminal} from "@xterm/xterm";
import {h} from "vue";
import {FitAddon} from 'xterm-addon-fit';
import {AttachAddon} from 'xterm-addon-attach';
import axios from "axios";


let term;
let socket;
let screenId = 1;
let screenName;
let screens = [];

function init() {
  term = new Terminal();
  term.open(document.getElementById('terminal'));
}

function connect() {
  socket = new WebSocket('ws://127.0.0.1:8080/api/ws/screen');
  socket.onopen = function () {
    const attachAddon = new AttachAddon(socket);
    term.loadAddon(attachAddon);
  };
  socket.onmessage = function (event) {
    term.write(event.data);
  };
  term.onData(data => {
    socket.send(data);
  });
  const fitAddon = new FitAddon();
  term.loadAddon(fitAddon);
  fitAddon.fit();
}

function createScreen() {
  screenName = "Terminal " + screenId;
  let api = "/api/v1/screen/create?id=" + screenId + "&name=" + encodeURIComponent(screenName);
  axios.get(api).then(res => {
    console.log(res)
    screenId++;
    getScreens()
  });
}

function getScreens() {
  let api = "/api/v1/screen/get_screens"
    axios.get(api).then(res => {
    const screensElement = document.getElementById('screens');
    screens = res.data;
    res.data.forEach(datum => {
      const button = h('el-button', {
        onClick: () => {
          // 鼠标点击事件
          console.log(`Button ${datum} clicked`);
        }
      }, datum);
      screensElement.appendChild(button.el);
    });
  });
}

getScreens()
// onMounted(init)

</script>


<style scoped>

</style>