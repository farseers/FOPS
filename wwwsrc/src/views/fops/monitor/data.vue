<template>
    <div class="system-user-container layout-padding">
        <el-card>
            <div class="system-user-search mb15">
                <el-button size="default" type="primary" class="ml10" @click="onSearch()">
                    <el-icon>
						<ele-Search />
					</el-icon>
                    查询</el-button>
                    
            </div>
            <el-table :data="tableData" v-loading="loading" style="width: 100%" size="default">
                <el-table-column type="index" label="序号" width="60" />
                <el-table-column prop="AppName" label="项目名称" show-overflow-tooltip></el-table-column>
                <el-table-column prop="Key" label="监控key" show-overflow-tooltip></el-table-column>
                <el-table-column prop="Value" label="监控value" show-overflow-tooltip></el-table-column>
                <el-table-column prop="CreateAt" label="发生时间" show-overflow-tooltip></el-table-column>
            </el-table>
            <InitPagination @sizeChange="onHandleSizeChange" @currentChange="onHandleCurrentChange" :pages="pages" />
        </el-card>
    </div>
</template>
<script>

import { fopsApi } from "/@/api/fops";
import InitPagination from '/src/views/components/InitPagination.vue';
import { ElMessage } from 'element-plus';
// AppName  string            // 项目名称
//     Key      string            // 监控key
//     Value    string            // 监控value
//     CreateAt dateTime.DateTime // 发生时间


const serverApi = fopsApi();
export default {
    components: { InitPagination },
    data() {
        return {
            tableData: [],
            loading: false,
            pages: {
                pageNum: 1,
                pageSize: 10,
                total: 0
            }

        }
    },
    mounted(){
        this.getTableData()
    },
    methods: {
        onHandleSizeChange(val) {
            this.pageSize = val;
            this.getTableData();
        },
        // 分页改变
        onHandleCurrentChange(val) {
            this.pageNum = val;
            this.getTableData();
        },
        onSearch(){
            this.pageNum = 1;
            this.getTableData();
        },
        getTableData() {
            this.loading = true;
            serverApi.monitorDataList({
                "pageSize": this.pages.pageSize,
                "pageIndex": this.pages.pageNum
            }).then((d)=>{
                this.loading = false;
                if(d.Status){
                    const { List,RecordCount } = d.Data;
                    this.tableData = [...List];
                    this.pages.total = RecordCount;
                }else{
                    ElMessage.error(d.StatusMessage);
                }
            }).catch(e=>{
                ElMessage.error('网络错误');
            })
        }
    }
}
</script>