<template>
    <div>
        <LayMain>
            <template #header>
                <el-button size="default" type="primary" class="ml10" @click="getTableData()">
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
                <el-table default-expand-all :data="tableData" v-loading="loading" style="width: 100%;max-height: 100%;" size="default">
                    <el-table-column type="index" label="序号" width="60" />
                    <el-table-column type="expand">
                        <template #default="scope">
                           <div class="reps_expand">
                            <div style="display: flex;align-items: center;margin-bottom: 5px;">
                                <span style="margin-right: 10px;">用户名：{{ scope.row.Username }}</span>
                                <span style="margin-right: 10px;">上次备份时间：{{ scope.row.LastBackupAt }}</span>
                                <span style="margin-right: 10px;">下次执行时间：{{ scope.row.NextBackupAt }}</span>
                            </div>
                            <div style="display: flex;align-items: center;">
                                数据库：
                                <div style="flex: 1;" v-if="scope.row.Database && scope.row.Database.length > 0">
                                        <el-button size="small" style="margin-right: 5px;" 
                                        type="primary" 
                                        v-for="r,index in scope.row.Database" :key="index" @click="get_base_data_dr(r, scope.row)">{{ r }}</el-button>
                                    </div>
                            </div>
                           </div>
                        </template>
                    </el-table-column>
                    <el-table-column prop="BackupDataType" label="数据库类型" min-width="100px">
                        <template #default="scope">
                           <span v-show="scope.row.BackupDataType == 0">Mysql</span>
                           <span v-show="scope.row.BackupDataType == 1">Clickhouse</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="主机" prop="Host" min-width="100">
                        <template #default="scope">
                          <span>{{ scope.row.Host }}:{{ scope.row.Port }}</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="Cron" prop="Cron" min-width="100"></el-table-column>
                   
                    <el-table-column prop="StoreType" label="存储类型" min-width="100px">
                        <template #default="scope">
                           <span v-show="scope.row.StoreType == 0">OSS</span>
                           <span v-show="scope.row.StoreType == 1">本地目录</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="180px" fixed="right" align="center">
                        <template #default="scope">
                            <el-button @click="base_info(scope.row)" type="primary" text size="small">备份</el-button>
                            <el-button @click="set_edit(scope.row)" type="primary" text size="small">编辑</el-button>
                            <el-button @click="em_info(scope.row)" type="primary" text size="small">清空</el-button>
                            <el-button @click="set_del(scope.row)" type="primary" text size="small">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </template>
            <template #footer>
                <!-- <InitPagination @sizeChange="onHandleSizeChange" @currentChange="onHandleCurrentChange"
                    :pages="pages" /> -->
            </template>
        </LayMain>
        <ScheduleDialog ref="editInfo" @search="getTableData" />
        <ScheduleDrawer ref="scheduleDrawer"/>
    </div>
</template>

<script>
import LayMain from '/src/views/components/LayMain.vue';
import { fopsApi } from "/@/api/fops";
import InitPagination from '/src/views/components/InitPagination.vue';
import ScheduleDialog from './scheduleDialog.vue';
import  ScheduleDrawer from './scheduleDrawer.vue';
import { ElMessageBox, ElMessage } from 'element-plus';
const serverApi = fopsApi();
export default {
    components: { InitPagination, ScheduleDialog,ScheduleDrawer, LayMain },
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

        em_info(row){
            const str = "确定清空此计划?"
            this.$confirm(str, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.loading = true;
                serverApi.backupData_clear({
                    "id": row.Id,   
                }).then(d => {
            let { Status, StatusMessage } = d;
            this.loading = false;
            if (Status) {
                this.$message({
                    type: 'success',
                    message: '清空成功'
                });
              this.getTableData()
            } else {
              ElMessage.error(StatusMessage)
            }
          })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消清空'
                });
            });
        },
        base_info(row){
            const str = "确定备份此计划?"
            this.$confirm(str, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.loading = true;
                serverApi.backupData_backup({
                    "id": row.Id,   
                }).then(d => {
            let { Status, StatusMessage } = d;
            this.loading = false;
            if (Status) {
              
                this.$message({
                    type: 'success',
                    message: '备份成功'
                });

              this.getTableData()
            } else {
              ElMessage.error(StatusMessage)
            }
          })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消备份'
                });
            });
        },
        get_base_data_dr(t,row){
            this.$refs.scheduleDrawer.handleNav(t,row)
        },
        set_add() {
            this.$refs.editInfo && this.$refs.editInfo.info(null)
        },
        set_edit(row) {
            if (row.Id) {
                this.$refs.editInfo && this.$refs.editInfo.info(row.Id)
            }
        },
        set_del(row){
            const str = "确定删除此备份计划?"
            this.$confirm(str, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                this.loading = true;
                serverApi.backupData_delete({
                    "id": row.Id,   
                }).then(d => {
            let { Status, StatusMessage } = d;
            this.loading = false;
            if (Status) {
               
                this.$message({
                    type: 'success',
                    message: '删除成功'
                });
              this.getTableData()
            } else {
              ElMessage.error(StatusMessage)
            }
          })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消删除'
                });
            });
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
        getTableData() {
            this.loading = true;
            serverApi.backupData_list({}).then((d) => {
                this.loading = false;
                if (d.Status) {
                    const List = d.Data;
                    this.tableData = [...List];
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
