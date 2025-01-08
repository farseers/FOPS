<template>
  <div class="system-user-container layout-padding">
    <el-card shadow="hover" class="layout-padding-auto">
      <div class="system-user-search mb15">
        <label class="ml5">应用</label>
        <el-select class="ml5" style="max-width: 110px;" size="small" v-model="state.appName">
          <el-option label="全部" value=""></el-option>
          <el-option v-for="item in state.appData" :label="item.AppName" :value="item.AppName" ></el-option>
        </el-select>
        <label class="ml10">访问节点</label>
        <el-input class="ml5" size="default" v-model="state.visitsNode" clearable style="max-width: 400px;"> </el-input>
        <label class="ml10">开始时间</label>
        <el-date-picker class="ml5" size="default" v-model="state.startAt" type="datetime" :shortcuts="shortcuts" placeholder="开始时间" clearable style="width: 200px"></el-date-picker>
        <label class="ml10">结束时间</label>
        <el-date-picker class="ml5" size="default" v-model="state.endAt" type="datetime" :shortcuts="shortcuts" placeholder="结束时间" clearable style="width: 200px"></el-date-picker>
        <el-button size="default" type="primary" class="ml5" @click="onQuery">
          <el-icon><ele-Search /></el-icon>
          查询
        </el-button>
        <el-button size="default" type="info" class="ml5" @click="onBack">
          回退
        </el-button>
      </div>
      <el-table :data="state.tableData.data" v-loading="state.tableData.loading" style="width: 100%" :cell-style="{padding:'2px 0'}">
        <el-table-column label="访问节点" show-overflow-tooltip>
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.VisitsNode.endsWith('/')" @click="onVisitsNode(scope.row.VisitsNode)" style="font-size: 14px;cursor:pointer">{{scope.row.VisitsNode}}</el-tag>
            <el-tag size="small" type="info" v-else  style="font-size: 14px;">{{scope.row.VisitsNode}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column width="140px" label="min" show-overflow-tooltip align="right">
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.MinMs >= 1000" type="danger" style="font-size: 14px;">{{scope.row.MinMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.MinMs >= 500" type="warning" style="font-size: 14px;">{{scope.row.MinMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.MinMs >= 100" style="font-size: 14px;">{{scope.row.MinMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.MinMs >= 50" type="success" style="font-size: 14px;">{{scope.row.MinMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else type="info" style="font-size: 14px;">{{scope.row.MinMs.toFixed(1)}}</el-tag>
            ms
          </template>
        </el-table-column>
        <el-table-column width="140px" label="max" show-overflow-tooltip align="right">
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.MaxMs >= 1000" type="danger" style="font-size: 14px;">{{scope.row.MaxMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.MaxMs >= 500" type="warning" style="font-size: 14px;">{{scope.row.MaxMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.MaxMs >= 100" style="font-size: 14px;">{{scope.row.MaxMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.MaxMs >= 50" type="success" style="font-size: 14px;">{{scope.row.MaxMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else type="info" style="font-size: 14px;">{{scope.row.MaxMs.toFixed(1)}}</el-tag>
            ms
          </template>
        </el-table-column>
        <el-table-column width="140px" label="avg" show-overflow-tooltip align="right">
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.AvgMs >= 1000" type="danger" style="font-size: 14px;">{{scope.row.AvgMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.AvgMs >= 500" type="warning" style="font-size: 14px;">{{scope.row.AvgMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.AvgMs >= 100" style="font-size: 14px;">{{scope.row.AvgMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.AvgMs >= 50" type="success" style="font-size: 14px;">{{scope.row.AvgMs.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else type="info" style="font-size: 14px;">{{scope.row.AvgMs.toFixed(1)}}</el-tag>
            ms
          </template>
        </el-table-column>
        <el-table-column width="140px" label="95line" show-overflow-tooltip align="right">
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.Line95Ms >= 1000" type="danger" style="font-size: 14px;">{{scope.row.Line95Ms.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.Line95Ms >= 500" type="warning" style="font-size: 14px;">{{scope.row.Line95Ms.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.Line95Ms >= 100" style="font-size: 14px;">{{scope.row.Line95Ms.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.Line95Ms >= 50" type="success" style="font-size: 14px;">{{scope.row.Line95Ms.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else type="info" style="font-size: 14px;">{{scope.row.Line95Ms.toFixed(1)}}</el-tag>
            ms
          </template>
        </el-table-column>
        <el-table-column width="140px" label="99line" show-overflow-tooltip align="right">
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.Line99Ms >= 1000" type="danger" style="font-size: 14px;">{{scope.row.Line99Ms.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.Line99Ms >= 500" type="warning" style="font-size: 14px;">{{scope.row.Line99Ms.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.Line99Ms >= 100" style="font-size: 14px;">{{scope.row.Line99Ms.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else-if="scope.row.Line99Ms >= 50" type="success" style="font-size: 14px;">{{scope.row.Line99Ms.toFixed(1)}}</el-tag>
            <el-tag size="small" v-else type="info" style="font-size: 14px;">{{scope.row.Line99Ms.toFixed(1)}}</el-tag>
            ms
          </template>
        </el-table-column>
        <el-table-column width="140px" label="错误量" show-overflow-tooltip align="right">
          <template #default="scope">
            <el-tag size="small" v-if="scope.row.ErrorCount >= 1" type="danger" style="font-size: 14px;">{{scope.row.ErrorCount}}</el-tag>
            <el-tag size="small" v-else type="info" style="font-size: 14px;">{{scope.row.ErrorCount}}</el-tag>
          </template>
        </el-table-column>
        <el-table-column width="100px" prop="TotalCount" label="访问次数" show-overflow-tooltip style="font-size: 14px;" align="right"></el-table-column>
        <el-table-column width="140px" label="QPS" show-overflow-tooltip style="font-size: 14px;" align="right">
          <template #default="scope">
            {{scope.row.QPS.toFixed(1)}} /s
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <detailDialog ref="detailDialogRef" @refresh="getTableData()" />
  </div>
</template>

<script setup lang="ts" name="indexVisits">

import {reactive, onMounted, ref, watch} from 'vue';
import { ElMessageBox, ElMessage } from 'element-plus';
import {fopsApi} from "/@/api/fops";
import {friendlyJSONstringify} from "@intlify/shared";

// 引入 api 请求接口
const serverApi = fopsApi();

// 定义变量内容
const detailDialogRef = ref();

let date = new Date()
date.setHours(0,0,0)

const state = reactive({
  appName:'',
  visitsNode:'',
  startAt: date.toLocaleString(),
  endAt: '',
  tableData: {
    data: [],
    total: 0,
    loading: false,
    param: {
    },
  },
  appData:[],
});

const shortcuts = [
  {
    text: '15分钟前',
    value: () => {
      const date = new Date()
      date.setTime(date.getTime() - 900 * 1000)
      return date
    },
  },
  {
    text: '半小时前',
    value: () => {
      const date = new Date()
      date.setTime(date.getTime() - 1800 * 1000)
      return date
    },
  },
  {
    text: '1小时前',
    value: () => {
      const date = new Date()
      date.setTime(date.getTime() - 3600 * 1000 * 1)
      return date
    },
  },
  {
    text: '2小时前',
    value: () => {
      const date = new Date()
      date.setTime(date.getTime() - 3600 * 1000 * 2)
      return date
    },
  },
  {
    text: '今天',
    value: () => {
      const date = new Date()
      date.setHours(0,0,0)
      return date
    },
  },
]

watch(() => state.appName, (newValue, oldValue) => {
  getTableData()
});

watch(() => state.visitsNode, (newValue, oldValue) => {
  getTableData()
});

watch(() => state.startAt, (newValue, oldValue) => {
  getTableData()
});

watch(() => state.endAt, (newValue, oldValue) => {
  getTableData()
});

// 初始化表格数据
const getTableData = () => {
  state.tableData.loading = true;

  var data={
    appName: state.appName,
    visitsNode: state.visitsNode,
    startAt: formatDate(state.startAt),
    endAt: formatDate(state.endAt),
  }
  const params = new URLSearchParams(data).toString();
  // 请求接口
  serverApi.visitsApi(params).then(function (res) {
    if (res.Status){
      state.tableData.data = res.Data;
      state.tableData.loading = false;
    }else{
      state.tableData.data=[]
      state.tableData.loading = false;
    }
  })
};

const getAppData=()=>{
  serverApi.dropDownList({}).then(function (res){
    if (res.Status){
      state.appData=res.Data
    }else{
      state.appData=[]
    }
  })
}

const onVisitsNode=(row: any)=>{
  state.visitsNode = row
}
const onBack=()=>{
  if (state.visitsNode == '') {
    return
  }
  if (state.visitsNode.endsWith('/')) {
    state.visitsNode = state.visitsNode.substring(0,state.visitsNode.length - 1)
  }
  state.visitsNode = state.visitsNode.substring(0, state.visitsNode.lastIndexOf('/')+1)

  if (state.visitsNode.endsWith("//")) {
    state.visitsNode = ""
  }
}

const onQuery=()=>{
  getTableData();
}
// 页面加载时
onMounted(() => {
  getAppData();
  getTableData();

  // 跳转到登录界面之后，不让其回退。就直接添加下面这段代码即可实现
  history.pushState(null, null, document.URL);
  window.addEventListener("popstate", function () {
    onBack()
    history.pushState(null, null, document.URL);
  });
});

const formatDate =(date, fmt) => {
  if (typeof date == 'string') {
    date = new Date(date)
  }

  if (!fmt) fmt = "yyyy-MM-dd hh:mm:ss";

  if (!date || date == null) return null;
  var o = {
    'M+': date.getMonth() + 1, // 月份
    'd+': date.getDate(), // 日
    'h+': date.getHours(), // 小时
    'm+': date.getMinutes(), // 分
    's+': date.getSeconds(), // 秒
    'q+': Math.floor((date.getMonth() + 3) / 3), // 季度
    'S': date.getMilliseconds() // 毫秒
  }
  if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length))
  for (var k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? (o[k]) : (('00' + o[k]).substr(('' + o[k]).length)))
  }
  return fmt
}
</script>

<style scoped lang="scss">
.system-user-container {
  :deep(.el-card__body) {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: auto;
    .el-table {
      flex: 1;
    }
  }
}
.cell {
  padding: 0 0;
}
.el-table--large .el-table__cell {
  padding: 0 0;
}
</style>
