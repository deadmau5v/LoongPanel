<template>
  <div id="content" style="width: 100%">
    <el-form @submit.native.prevent="submitForm" id="filePwd">
      <el-form-item>
        <el-input v-model="path" id="pwd"></el-input>
      </el-form-item>
    </el-form>

    <el-table :data="files || []" style="width: 100%">
      <el-table-column prop="name" label="文件名">
        <template #default="{ row }">
          <el-button v-if="row['isDir']" @click="dir(row['path'])" class="fa fa-folder"><el-text style="margin-left: 10px;">{{ row['name'] }}</el-text></el-button>
          <el-text v-else>{{ row['name'] }}</el-text>
        </template>
      </el-table-column>
      <el-table-column prop="size" label="大小">
        <template #default="{ row }">
          <el-text>{{ formatSize(row['size']) }}</el-text>
        </template>
      </el-table-column>
      <el-table-column prop="time" label="修改时间">
        <template #default="{ row }">
          <el-text>{{ row['time'] }}</el-text>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import axios from 'axios';
import { ElForm, ElFormItem, ElInput, ElTable, ElTableColumn, ElButton } from 'element-plus';

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
    submitForm() {
      this.dir(this.path);
    },
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
  },
  mounted() {
    this.dir(this.path);
  },
};
</script>

<style scoped>
</style>