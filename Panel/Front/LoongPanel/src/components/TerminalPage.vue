<template>
  <div id="terminal"></div>
  <el-button-group>
    <el-button @click="createScreen">创建</el-button>
    <el-button @click="term.clear()">清屏</el-button>
  </el-button-group>
  <el-button-group id="screens">

  </el-button-group>
  <link rel="stylesheet" href="https://cdn1.d5v.cc/CDN/Project/LoongPanel/static/xterm.css"/>
</template>


<script setup>
import {Terminal} from "@xterm/xterm"
import {onMounted} from "vue"
import axios from "axios"


let term
let socket
let screenId = 1
let screenName
let screens = []

function init() {
  term = new Terminal()
  term.open(document.getElementById('terminal'))
}

function connect(_id) {
  socket = new WebSocket('ws://127.0.0.1:8080/api/ws/screen?id=' + _id)
  socket.onmessage = function (event) {
    console.log(event)
    term.write(event.data)
  }
  term.onData(data => {
    socket.send(data)
  })
}

function createScreen() {
  screenName = "Terminal " + screenId
  let api = "/api/v1/screen/create?id=" + screenId + "&name=" + encodeURIComponent(screenName);
  axios.get(api).then(res => {
    console.log(res)
    screenId++
    getScreens()
  })
}

function getScreens() {
  let api = "/api/v1/screen/get_screens"
  axios.get(api).then(res => {
    const screensElement = document.getElementById('screens');
    screens = res.data
    res.data.forEach(datum => {
      const button = document.createElement("button")
      button.classList.add("el-button")
      button.type = "button"
      button.innerText = datum.name

      function _click() {
        screenId = datum.id
        try {
          socket.close()
        } catch {
        }
        connect(screenId)
      }

      button.onclick = _click
      screensElement.appendChild(button)
    })
  })
}

getScreens()
onMounted(init)

</script>


<style scoped>

</style>