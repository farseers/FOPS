<template>
    <div class="conlyTab">
        <h3 style="padding: 5px;">应用日志</h3>
        <div v-if="!state.empty" class="dsv">
            <div v-for="row,index in state.tableData" :key="index">
                <p class="cl1" >
                    <span style="color: #9caf62;margin-right: 3px;">{{ row.CreateAt.split(" ")[1] }}</span>
                    <el-tag size="small" style="margin: 2px 5px;">{{ row.AppName }} {{ row.AppIp }}</el-tag>
                  <el-tag :type="type_set(row)" size="small">{{ row.LogLevel }}</el-tag>
                </p>
                <p class="cl1">

                </p>
                <p class="cl1">
                   {{ row.Content }}
                </p>
              <hr />
            </div>
        </div>
        <el-empty v-else description="暂无数据"></el-empty>
    </div>
</template>
<script setup name="QueTab">
import { reactive, onMounted, defineExpose, ref } from 'vue';
import { fopsApi } from "/@/api/fops";

// 引入 api 请求接口
const serverApi = fopsApi();
// 定义变量内容
const state = reactive({
    tableData: [],
    empty:false
});

// 初始化表格数据
const getData = () => {
    var data = {
        appName: '',
        appIp: '',
        traceId: '',
        logContent: '',
        minute: 60,
        logLevel: '3',
        pageSize: '19',
        pageIndex: '1',
    }
    const params = new URLSearchParams(data).toString();
    // 请求接口
    serverApi.logList(params).then(function (res) {
        if (res.Status) {
            state.tableData = res.Data.List;
            if(state.tableData && state.tableData.length>0){
                state.empty = false
            }else{
                state.empty = true
            }
        }
    })
};
const type_set = (item)=>{
    let str = '';
    if(item){
        if(item.LogLevel == 'Info'){str = ''}
        else if(item.LogLevel == 'Debug'){str = 'info'}
        else if(item.LogLevel == 'Warn'){str = 'warning'}
        else if(item.LogLevel == 'Error'){str = 'danger'}
    }
    return str
     
}
onMounted(() => {
    getData()
});
defineExpose({
    getData,
});
</script>
<style scoped lang="scss">
.conlyTab {
    height: 100%;
    display: flex;
    flex-flow: column;

    .dsv {
        flex: 1;
        overflow: auto;color: #fff;background-color:#393d49;padding: 5px;font-size: 14px;
        p{
            display: flex;
            align-items: center;
        }
    }
    .cl1{
        color: var(--el-table-text-color);
    }
}
</style>