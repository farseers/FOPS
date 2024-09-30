<template>
    <div>
        <LayMain>
            <template #header>
                <el-form-item label="应用名称">
                        <el-select v-model="AppName" clearable style="width: 200px;" @change="onSearch()" filterable
                            placeholder="请选择">
                            <el-option v-for="item in m_list" :key="item.AppName" :label="item.AppName"
                                :value="item.AppName" />
                        </el-select>
                    </el-form-item>
                    <el-button size="default" type="primary" class="ml10" @click="onSearch()">
                        <el-icon>
                            <ele-Search />
                        </el-icon>
                        查询</el-button>
                    <el-button size="default" type="warning" class="ml10" @click="set_add()">
                        <el-icon>
                            <ele-Plus />
                        </el-icon>
                        新增</el-button>
            </template>
            <template #main>
                <el-table :data="tableData" v-loading="loading" style="width: 100%;height: 100%;" size="default">
                    <el-table-column type="index" label="序号" width="60" />
                    <el-table-column prop="AppName" label="应用名称" min-width="200px" show-overflow-tooltip>
                        <template #default="scope">
                            <el-tag size="small" style="margin-left: 3px;" v-if="scope.row.Enable"
                                type="success">启用</el-tag>
                            <el-tag size="small" style="margin-left: 3px;" v-else type="danger">停用</el-tag>
                            {{ scope.row.AppName }}
                        </template>
                    </el-table-column>
                    <el-table-column label="起止时间" width="110px">
                        <template #default="scope">
                            <div v-if="scope.row.TimeType == 1">{{ scope.row.StartDay }} </div>
                            <div v-if="scope.row.TimeType == 1">{{ scope.row.EndDay }}</div>
                            <div v-if="scope.row.TimeType == 0">{{ scope.row.StartDate }} </div>
                            <div v-if="scope.row.TimeType == 0">{{ scope.row.EndDate }}</div>
                        </template>
                    </el-table-column>
                    <el-table-column label="监控键值" width="180px">
                        <template #default="scope">
                            <span>{{ scope.row.KeyName }} </span>
                            <span style="margin: 0 10px;">{{ scope.row.Comparison }} </span>
                            <span>{{ scope.row.KeyValue }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="NoticeIds" label="关联人" show-overflow-tooltip>
                        <template #default="scope">
                            <span style="margin: 3px;" v-for="item, index in scope.row.NoticeList"
                                :key="index">{{ item.Name }}、</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="TipTemplate" label="模版" show-overflow-tooltip></el-table-column>
                    <el-table-column prop="Remark" label="备注" show-overflow-tooltip></el-table-column>
                    <el-table-column label="操作" width="100px" fixed="right" align="center">
                        <template #default="scope">
                            <el-button @click="set_edit(scope.row)" type="primary" text size="small">编辑</el-button>
                            <el-button @click="set_del(scope.row)" type="danger" text size="small">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </template>
            <template #footer>
                <InitPagination @sizeChange="onHandleSizeChange" @currentChange="onHandleCurrentChange"
                    :pages="pages" />
            </template>
        </LayMain>
        <GuleDialog ref="editInfo" @search="getTableData" />
    </div>
</template>

<script>
import LayMain from '/src/views/components/LayMain.vue';
import { fopsApi } from "/@/api/fops";
import InitPagination from '/src/views/components/InitPagination.vue';
import GuleDialog from './guleDialog.vue';
import { ElMessageBox, ElMessage } from 'element-plus';
const serverApi = fopsApi();
export default {
    components: { InitPagination, GuleDialog, LayMain },
    data() {
        return {
            tableData: [],
            loading: false,
            p_list: [],//关联人
            AppName: '',//应用名称
            m_list: [],//项目列表
            typeList: [],//比较方式
            pages: {
                pageNum: 1,
                pageSize: 10,
                total: 0
            }

        }
    },
    mounted() {
        this.getTableData()
        this.info()
    },
    methods: {
        info() { //获取关联人 项目名称
            serverApi.monitorNoticeList({
                "pageSize": 10000,
                "pageIndex": 1
            }).then((d) => {
                if (d.Status) {
                    const { List } = d.Data;
                    this.p_list = List.filter(item => item.Enable === true)
                } else {
                    ElMessage.error(d.StatusMessage);
                }
            }).catch(e => {
                ElMessage.error('网络错误');
            })
            serverApi.dropDownList({}).then(d => {
                const { Status, Data } = d;
                if (Status) {
                    this.m_list = [...Data]
                }
            })
            serverApi.drpBaseList({ baseType: '2' }).then(d => {
                const { Data, Status } = d;
                if (Status) {
                    const { CompareList } = Data
                    this.typeList = [...CompareList];
                }
            })
        },
        set_add() {
            this.$refs.editInfo && this.$refs.editInfo.info(null, this.p_list, this.m_list, this.typeList)
        },
        set_edit(row) {
            if (row.Id) {
                this.$refs.editInfo && this.$refs.editInfo.info(row.Id, this.p_list, this.m_list, this.typeList)
            }
        },
        set_del(row) { //删除
            let str = '确定删除此数据?'
            if (row.AppName) {
                str = `删除项目名称：“${row.AppName}”，是否继续?`;
            }
            const _this = this;
            ElMessageBox.confirm(`${str}`, '提示', {
                confirmButtonText: '确认',
                cancelButtonText: '取消',
                type: 'warning',
            })
                .then(() => {
                    // 删除逻辑
                    serverApi.monitorDelRule({ "id": row.Id }).then(function (res) {
                        if (res.Status) {
                            ElMessage.success('删除成功');
                            _this.getTableData()
                        } else {
                            ElMessage.error(res.StatusMessage)
                        }
                    })
                })
                .catch(() => { });
        },
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
            serverApi.monitorRuleList({
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