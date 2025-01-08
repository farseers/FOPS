<template>
    <div>
        <LayMain>
            <template #header>
                <el-form-item label="项目名称">
                    <el-select v-model="AppName" clearable style="width: 200px;" @change="onSearch()" filterable
                        placeholder="请选择项目名称">
                        <el-option v-for="item in m_list" :key="item.AppName" :label="item.AppName"
                            :value="item.AppName" />
                    </el-select>
                </el-form-item>
                <el-button size="default" type="primary" class="ml10" @click="onSearch()">
                    <el-icon>
                        <ele-Search />
                    </el-icon>
                    查询</el-button>
            </template>
            <template #main>
                <el-table :data="tableData" v-loading="loading" style="width: 100%;height: 100%;" size="default">
                    <el-table-column type="index" label="序号" width="60" />
                    <el-table-column prop="AppName" label="项目名称" show-overflow-tooltip></el-table-column>
                    <el-table-column prop="Key" label="监控key" show-overflow-tooltip></el-table-column>
                    <el-table-column prop="Value" label="监控value" show-overflow-tooltip></el-table-column>
                    <el-table-column prop="CreateAt" label="发生时间" show-overflow-tooltip></el-table-column>
                </el-table>
            </template>
            <template #footer>
                <InitPagination @sizeChange="onHandleSizeChange" @currentChange="onHandleCurrentChange"
                    :pages="pages" />
            </template>
        </LayMain>
    </div>
</template>
<script>
import LayMain from '/src/views/components/LayMain.vue';
import { fopsApi } from "/@/api/fops";
import InitPagination from '/src/views/components/InitPagination.vue';
import { ElMessage } from 'element-plus';
// AppName  string            // 项目名称
//     Key      string            // 监控key
//     Value    string            // 监控value
//     CreateAt dateTime.DateTime // 发生时间


const serverApi = fopsApi();
export default {
    components: { InitPagination, LayMain },
    data() {
        return {
            tableData: [],
            m_list: [],
            loading: false,
            AppName: '',//项目名称
            pages: {
                pageNum: 1,
                pageSize: 10,
                total: 0
            }

        }
    },
    created() {
        serverApi.dropDownList({ IsAll: true }).then(d => {
            const { Status, Data } = d;
            if (Status) {
                this.m_list = [...Data]
            }
        })
    },
    mounted() {
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
        onSearch() {
            this.pageNum = 1;
            this.getTableData();
        },
        getTableData() {
            this.loading = true;
            serverApi.monitorDataList({
                'appName': this.AppName,
                "pageSize": this.pages.pageSize,
                "pageIndex": this.pages.pageNum
            }).then((d) => {
                this.loading = false;
                if (d.Status) {
                    const { List, RecordCount } = d.Data;
                    this.tableData = [...List];
                    this.pages.total = RecordCount;
                } else {
                    ElMessage.error(d.StatusMessage);
                }
            }).catch(e => {
                ElMessage.error('网络错误');
            })
        }
    }
}
</script>