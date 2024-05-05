<template>
  <div id="content" style="width: 100%">
    <el-form @submit.native.prevent="submitForm" id="filePwd">
      <el-form-item>
        <el-input v-model="path" id="pwd"></el-input>
      </el-form-item>
    </el-form>
    <!--文件名-->
    <el-table :data="files || []" style="width: 100%">
      <el-table-column prop="name" label="文件名">
        <template #default="{ row }">
          <!--   文件夹     -->
          <el-button v-if="row['isDir']" @click="dir(row['path'])" class="fa fa-folder">
            <el-text style="margin-left: 10px;">
              {{ row['name'] }}
            </el-text>
          </el-button>
          <!--    文件      -->
          <el-button v-else class="fa fa-file">
            <el-text style="margin-left: 10px;">
              {{ row['name'] }}
            </el-text>
          </el-button>

        </template>
      </el-table-column>
      <!--  文件大小    -->
      <el-table-column prop="size" label="大小" >
        <template #default="{ row }" >
          <el-text v-if="row['showSize']">{{ formatSize(row['size']) }}</el-text>
        </template>
      </el-table-column>
      <!--  修改时间    -->
      <el-table-column prop="time" label="修改时间">
        <template #default="{ row }" >
          <el-text  >{{ row['time'] }}</el-text>
        </template>
      </el-table-column>
      <!--  文件操作   -->
      <el-table-column label="操作">
        <template #default="{ row }" >
          <el-button v-if="row['showEdit']" type="primary" @click="download(row['path'])">下载</el-button>
          <el-button v-if="row['showEdit']" type="danger" @click="deleteFile(row['path'])">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import axios from 'axios';
import {ElForm, ElFormItem, ElInput, ElTable, ElTableColumn, ElButton} from 'element-plus';

export default {
  components: {
    ElForm,
    ElFormItem,
    ElInput,
    ElTable,
    ElTableColumn,
    ElButton
  },
  data() {
    return {
      path: '',
      "files": [],
    };
  },
  methods: {
    // 获取文件夹
    async dir(path) {
      const response = await axios.get("/api/v1/files/dir?path=" + path);
      this.files = response.data.files.sort((a, b) => {
        if (a["isDir"] && a["name"] === "..") {
          return -1;
        } else if (b["isDir"] && b["name"] === "..") {
          return 1;
        }
        if (a["isDir"] && !b["isDir"]) {
          return -1;
        }
        if (!a["isDir"] && b["isDir"]) {
          return 1;
        }
        return a["name"].localeCompare(b["name"]);
      });
      this.path = decodeURIComponent(path);
    },
    // 提交表单
    submitForm() {
      this.dir(this.path);
    },
    // 格式化文件大小
    formatSize(size) {
      if (size > 1073741824) {
        return (size / 1024 / 1024 / 1024).toFixed(2) + "GB";
      } else if (size > 10488832) {
        return (size / 1024 / 1024).toFixed(2) + "MB";
      } else if (size > 1024) {
        return (size / 1024).toFixed(2) + "KB";
      } else {
        return size.toFixed(2) + "B";
      }
    },
    // 下载文件
    download(path) {
      console.log(this)
      // window.open("/api/v1/files/download?path=" + path);
    },
    // 删除文件
    async deleteFile(path) {
      await axios.delete("/api/v1/files/delete?path=" + path);
      await this.dir(this.path);
    },
  },
  mounted() {
    this.dir(this.path);
  },
};
</script>

<style scoped>
</style>