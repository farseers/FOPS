<template>
    <div>
        <el-drawer
            :title="title"
            v-model="dialogVisible"
            :direction="direction"
             size="85%"
            :before-close="handleClose">
            <div style="display: flex;flex-flow: column;height: 100%;">
                <div style="flex: 1;"  ref="navHe">
                    <el-table :data="dataList" size="mini" :max-height="mhs">
                        <el-table-column type="index" width="50" label="#"></el-table-column>
                        <el-table-column property="FileName" label="文件名" min-width="280"></el-table-column>
                        <el-table-column prop="StoreType" label="存储类型" min-width="100px">
                        <template #default="scope">
                           <span v-show="scope.row.StoreType == 0">OSS</span>
                           <span v-show="scope.row.StoreType == 1">本地目录</span>
                        </template>
                    </el-table-column>
                    <el-table-column prop="CreateAt" label="备份时间" min-width="180px"></el-table-column>
                    <el-table-column prop="Size" label="文件大小" min-width="150px">
                        <template #default="scope">
                          <span>{{scope.row.Size}}（KB）</span>
                        </template>
                    </el-table-column>
                    <el-table-column label="操作" width="100px" fixed="right" align="center">
                        <template #default="scope">
                            <el-button @click="del(scope.row)" type="primary" text size="small">删除</el-button>
                            <el-button @click="rest(scope.row)" type="primary" text size="small">恢复</el-button>
                        </template>
                    </el-table-column>
                    </el-table>
                </div>
                <span slot="footer" class="dialog-footer" style="padding: 10px 0;background-color: #f1f1f1;text-align:center">
                    <el-button @click="handleClose()" type="primary">关 闭</el-button>
                </span>
            </div>
        </el-drawer>
    </div>
</template>
<script>
import { fopsApi } from "/@/api/fops";
import { ElMessage } from 'element-plus';
const serverApi = fopsApi();
    export default{
        components: {  },
        name:'lineNav',
        data(){
            return {
                dialogVisible:false,
                mhs:'600px',
                direction: 'rtl',
                title:'备份详细',
                dataList:[]
            }
        },
        methods:{
            handleChange(page, size) { //分页变化
            this.pagRow.currentPage = page;
            this.pagRow.pageSize = size;
            this.search(page)
        },
        del(row){
            const str = "确定删除["+row.FileName+"]?"
            this.$confirm(str, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                serverApi.backupData_deleteHistory({
                    "backupId": row.backupId,   
                    "FileName": row.FileName   
                }).then(d => {
            let { Status, StatusMessage } = d;
            if (Status) {
                this.$message({
                    type: 'success',
                    message: '删除成功'
                });
              this.search()
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
        rest(row){
            const str = "确定恢复["+row.FileName+"]?"
            this.$confirm(str, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                serverApi.backupData_recoverBackupFile({
                    "backupId": row.backupId,   
                    "FileName": row.FileName   
                }).then(d => {
            let { Status, StatusMessage } = d;
            if (Status) {
                this.$message({
                    type: 'success',
                    message: '恢复成功'
                });
              this.search()
            } else {
              ElMessage.error(StatusMessage)
            }
          })
            }).catch(() => {
                this.$message({
                    type: 'info',
                    message: '已取消恢复'
                });
            });
        },
        handleNav(row){ //线路列表
                this.dataList = [];
                this.Id = '';
                if(row){
                    this.Id = row.Id;
                    this.search()
                }
            },
            search(){
               
                serverApi.backupData_backupList({
                    backupId:this.Id
                    }).then(d => {
                        let { Status, StatusMessage,Data } = d;
                        if (Status) {
                            this.dataList = Data;
                            this.dialogVisible = true
                            this.$nextTick(() => {
                                const divHeight = this.$refs.navHe.clientHeight - 40;
                                this.mhs = divHeight + 'px'
                            });
                        } else {
                        ElMessage.error(StatusMessage)
                        }
                    })
            },
            handleClose(){
                this.dialogVisible = false
            },
        },
        filters:{
            formattedAmounts(v){
                return formattedAmount(v)
            }
        }
    }
</script>
<style scoped>
.lineDialog{
    z-index: 4099 !important;
}
.pages{
    margin-bottom: 10px;
}

</style>