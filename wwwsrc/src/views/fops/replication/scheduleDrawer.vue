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
                    <div style="margin-top: 10px;">
                        <el-tag v-for="t,index in baseData" :key="index" :type="ck_t==t?'':'info'" style="cursor: pointer;margin-left: 5px;" @click="ck_ts(t)">{{ t }}</el-tag>
                    </div>
                    <el-table :data="dataList"  :max-height="mhs">
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
                baseData:[],
                ck_t:'',
                dialogVisible:false,
                mhs:'600px',
                direction: 'rtl',
                title:'备份详细',
                dataList:[]
            }
        },
        methods:{
        ck_ts(t){
            this.ck_t = t;
            this.search()
        },
        del(row){
           
            const str = "确定删除["+row.FileName+"]?";
            var par = {
                    "backupId": row.BackupId,   
                    "FileName": row.FileName   
                }
            this.$confirm(str, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                serverApi.backupData_deleteHistory(par).then(d => {
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
           
            const str = "确定恢复["+row.FileName+"]?";
            var par = {
                    "backupId": row.BackupId,   
                    "FileName": row.FileName,   
                    "database": row.Database   
                }
            this.$confirm(str, '提示', {
                confirmButtonText: '确定',
                cancelButtonText: '取消',
                type: 'warning'
            }).then(() => {
                serverApi.backupData_recoverBackupFile(par).then(d => {
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
        handleNav(t,row){ //线路列表
                this.dataList = [];
                this.Id = '';
                this.ck_t = t;
                if(row){
                    this.Id = row.Id;
                    this.baseData = row.Database
                    this.search()
                }
            },
            search(){
               
                serverApi.backupData_backupList({
                    backupId:this.Id,
                    database:this.ck_t
                    }).then(d => {
                        let { Status, StatusMessage,Data } = d;
                        if (Status) {
                            this.dataList = Data;
                            this.dialogVisible = true
                            this.$nextTick(() => {
                                const divHeight = this.$refs.navHe.clientHeight - 70;
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