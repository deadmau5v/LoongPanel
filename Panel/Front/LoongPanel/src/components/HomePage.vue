<template>
  <div>
    <!--  性能监视表  -->
    <el-card id="echarts-cards">
      <el-card class="box-card echarts-card">
        <div slot="header" class="clearfix">
          <span>磁盘使用</span>
        </div>
        <div id="diskUsageChart" class="charts"></div>
      </el-card>

      <el-card class="box-card echarts-card">
        <div slot="header" class="clearfix">
          <span>平均负荷</span>
        </div>
        <div id="averageLoadChart" class="charts"></div>
      </el-card>

      <el-card class="box-card echarts-card">
        <div slot="header" class="clearfix">
          <span>内存使用</span>
        </div>
        <div id="memoryUsageChart" class="charts"></div>
      </el-card>

      <el-card class="box-card echarts-card">
        <div slot="header" class="clearfix">
          <span>CPU使用</span>
        </div>
        <div id="cpuUsageChart" class="charts"></div>
      </el-card>
    </el-card>

    <!--  开关机重启  -->
    <span class="fa fa-power-off"><span style="padding-left: 10px"></span>电源</span>
    <el-card class="box-card">
      <el-button type="primary" class="power-button" @click="shutdown()">关闭服务器</el-button>
      <el-button type="primary" class="power-button" @click="reboot()">重启服务器</el-button>
    </el-card>

    <!--  垃圾清理  -->

  </div>
</template>

<script setup>
import {onMounted} from 'vue';
import * as echarts from 'echarts';
import axios from "axios";

let pageName = "首页"
document.title = "LoongPanel - " + pageName

let diskUsageChart = null;
let averageLoadChart = null;
let memoryUsageChart = null;
let cpuUsageChart = null;

let col_w = 1000

function reboot() {
  axios.get('/api/v1/power/reboot').then(res => {
  })
}

function shutdown() {
  axios.get('/api/v1/power/shutdown').then(res => {
  })
}

function drawDiskUsageChart(usedPercentage) {
  let doc = document.getElementById('diskUsageChart')
  let chart = echarts.init(doc)
  let option = {
    tooltip: {
      show: false
    },
    series: [
      {
        name: '',
        type: 'pie',
        radius: ['70%', '80%'],
        avoidLabelOverlap: false,
        label: {
          show: true,
          position: 'center',
          formatter: `${usedPercentage}%\n磁盘`,
          fontSize: '20',
          fontWeight: 'bold',
          color: 'black',

        },
        labelLine: {
          show: false
        },
        data: [
          {value: usedPercentage, itemStyle: {color: 'rgba(255,0,0,0.3)'}},
          {value: 100 - usedPercentage, itemStyle: {color: 'rgba(100,100,100,0.2)'}},
        ],
      }
    ]
  };
  chart.setOption(option)

}

function drawAverageLoadChart(usedPercentage) {
  let doc = document.getElementById('averageLoadChart')
  let chart = echarts.init(doc)
  let option = {
    tooltip: {
      show: false
    },
    series: [
      {
        name: '',
        type: 'pie',
        radius: ['70%', '80%'],
        avoidLabelOverlap: false,
        label: {
          show: true,
          position: 'center',
          formatter: `${usedPercentage}%\n负荷`,
          fontSize: '20',
          fontWeight: 'bold',
          color: 'black',

        },
        labelLine: {
          show: false
        },
        data: [
          {value: usedPercentage, itemStyle: {color: 'rgba(255,0,0,0.3)'}},
          {value: 100 - usedPercentage, itemStyle: {color: 'rgba(100,100,100,0.2)'}},
        ],
      }
    ]
  };
  chart.setOption(option)
}

function drawMemoryUsageChart(usedPercentage) {
  let doc = document.getElementById('memoryUsageChart')
  let chart = echarts.init(doc)
  let option = {
    tooltip: {
      show: false
    },
    series: [
      {
        name: '',
        type: 'pie',
        radius: ['70%', '80%'],
        avoidLabelOverlap: false,
        label: {
          show: true,
          position: 'center',
          formatter: `${usedPercentage}%\n内存`,
          fontSize: '20',
          fontWeight: 'bold',
          color: 'black',

        },
        labelLine: {
          show: false
        },
        data: [
          {value: usedPercentage, itemStyle: {color: 'rgba(255,0,0,0.3)'}},
          {value: 100 - usedPercentage, itemStyle: {color: 'rgba(100,100,100,0.2)'}},
        ],
      }
    ]
  };
  chart.setOption(option)
}

function drawCpuUsageChart(usedPercentage) {
  let doc = document.getElementById('cpuUsageChart')
  let chart = echarts.init(doc)
  let option = {
    tooltip: {
      show: false
    },
    series: [
      {
        name: '',
        type: 'pie',
        radius: ['70%', '80%'],
        avoidLabelOverlap: false,
        label: {
          show: true,
          position: 'center',
          formatter: `${usedPercentage}%\nCPU`,
          fontSize: '20',
          fontWeight: 'bold',
          color: 'black',

        },
        labelLine: {
          show: false
        },
        data: [
          {value: usedPercentage, itemStyle: {color: 'rgba(255,0,0,0.3)'}},
          {value: 100 - usedPercentage, itemStyle: {color: 'rgba(100,100,100,0.2)'}},
        ],
      }
    ]
  };
  chart.setOption(option)
}

function init() {
  diskUsageChart = echarts.init(document.getElementById('diskUsageChart'));
  averageLoadChart = echarts.init(document.getElementById('averageLoadChart'));
  memoryUsageChart = echarts.init(document.getElementById('memoryUsageChart'));
  cpuUsageChart = echarts.init(document.getElementById('cpuUsageChart'));
  drawDiskUsageChart(0)
  drawAverageLoadChart(0)
  drawMemoryUsageChart(0)
  drawCpuUsageChart(0)
}

function getData() {
  // 设置定时器
  setInterval(() => {
    axios.get('/api/v1/status/system_status').then(res => {
      let disk_usage = res.data["disk_usage"]
      let average_load = res.data["average_load"]
      let memory_usage = res.data["memory_usage"]
      let cpu_usage = res.data["cpu_usage"]
      // 保留两位小数
      drawDiskUsageChart(disk_usage.toFixed(2))
      drawAverageLoadChart(average_load.toFixed(2))
      drawMemoryUsageChart(memory_usage.toFixed(2))
      drawCpuUsageChart(cpu_usage.toFixed(2))

    })
  }, 1000)
}

onMounted(init);
onMounted(getData)
</script>

<style scoped>
.box-card {
  margin: 20px;
  box-shadow: none;
  border: none;
}

.echarts-card {
  width: 250px;
  height: 260px;
}

.el-button {
  margin: 10px 0;
}

.charts {
  width: 200px;
  height: 200px;

}

.clearfix {
  text-align: center;
}

html {
  background-color: #f0f0f0;
}

.power-button {
  background-color: rgba(255, 77, 81, 0.6);
  color: white;
  border-color: #ff4d51;
}

.power-button:hover {
  background-color: rgba(255, 77, 81, 0.8);
  color: white;
  border-color: #ff4d51;
}

.power-button:active {
  background-color: rgba(255, 77, 81, 1);
  color: white;
  border-color: #ff4d51;
}
</style>