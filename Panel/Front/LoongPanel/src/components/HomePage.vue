<template>
  <div>
    <el-row id="box1">
      <el-col :span="col_w">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>磁盘使用</span>
          </div>
          <div id="diskUsageChart" class="charts"></div>
        </el-card>
      </el-col>
      <el-col :span="col_w">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>平均负荷</span>
          </div>
          <div id="averageLoadChart" class="charts"></div>
        </el-card>
      </el-col>
      <el-col :span="col_w">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>内存使用</span>
          </div>
          <div id="memoryUsageChart" class="charts"></div>
        </el-card>
      </el-col>
      <el-col :span="col_w">
        <el-card class="box-card">
          <div slot="header" class="clearfix">
            <span>CPU使用</span>
          </div>
          <div id="cpuUsageChart" class="charts"></div>
        </el-card>
      </el-col>
    </el-row>
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
}

function getData() {
  drawDiskUsageChart(0);
  drawAverageLoadChart(0);
  drawMemoryUsageChart(0);
  drawCpuUsageChart(0);
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

.el-row {
  justify-content: center
}

.charts {
  width: 200px;
  height: 200px;

}

.clearfix {
  text-align: center;
}

#box1 {
  background-color: white;
  border-radius: 10px;
  box-shadow: 0 0 5px rgba(255, 0, 0, 0.3);
  margin: 20px;
  padding: 20px;
}

html {
  background-color: #f0f0f0;
}
</style>