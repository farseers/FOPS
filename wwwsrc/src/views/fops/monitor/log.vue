<template>
    <div class="system-user-container layout-padding">
        <el-card>
            <div class="system-user-search mb15">
                <el-button size="default" type="primary" class="ml10" @click="onSearch()">
                    <el-icon>
						<ele-Search />
					</el-icon>
                    查询</el-button>
                    <el-button size="default" type="warning" class="ml5" @click="set_del">
					<el-icon><ele-Delete /></el-icon>
					删除七天前数据
				</el-button>
            </div>
            <el-table :data="tableData" v-loading="loading" style="width: 100%" size="default">
                <el-table-column type="index" label="序号" width="60" />
                <el-table-column prop="AppName" label="项目名称" width="120px"></el-table-column>
                <el-table-column prop="NoticeName" label="通知人" width="110px"></el-table-column>
                <el-table-column prop="NoticeType" label="通知类型" width="110px">
                    <template #default="scope">
                            <span v-text="set_type(scope.row.NoticeType)"></span>
                    </template>
                </el-table-column>
                <el-table-column prop="NoticeMsg" label="通知消息"></el-table-column>
                <el-table-column prop="NoticeAt" label="通知时间" width="160px"></el-table-column>
            </el-table>
            <InitPagination @sizeChange="onHandleSizeChange" @currentChange="onHandleCurrentChange" :pages="pages" />
        </el-card>
        <noticeDialog ref="editInfo" @search="getTableData"/>
    </div>
</template>
<script>

import { fopsApi } from "/@/api/fops";
import InitPagination from '/src/views/components/InitPagination.vue';
import noticeDialog from './noticeDialog.vue';
import { ElMessageBox, ElMessage } from 'element-plus';
// Id         int64           // 主键
//     AppName    string          // 项目名称
//     NoticeId   int64           // 通知Id
//     NoticeType noticeType.Enum // 0 whatsapp
//     NoticeMsg  string          // 通知消息
//     NoticeAt   time.Time       // 通知时间

const serverApi = fopsApi();
export default {
    components: { InitPagination,noticeDialog },
    data() {
        return {
            typeList:[],//通知类型
            tableData: [],
            loading: false,
            pages: {
                pageNum: 1,
                pageSize: 10,
                total: 0
            }

        }
    },
    created(){
        this.info()
    },
    mounted(){
        this.getTableData()
    },
    methods: {
        info(){
            serverApi.monitorNoticeTypeList({}).then(d=>{
                const { Data,Status } = d;
                if(Status){
                    this.typeList = [...Data];
                }
            })
        },
        set_type(type){
           const row =  this.typeList.find(item=>{
                return item.NoticeType == type
            })
            if(row){
                return row.NoticeTypeName
            }else{
                return ''
            }
        },
        set_del(){ //删除
            let str = '确定删除七天前数据?'
            const _this = this;
            ElMessageBox.confirm(`${str}`, '提示', {
                confirmButtonText: '确认',
                cancelButtonText: '取消',
                type: 'warning',
            })
                .then(() => {
                    // 删除逻辑
                    serverApi.monitorDelNoticeLog({}).then(function (res){
                    if (res.Status){
                        ElMessage.success('删除成功');
                        _this.getTableData()
                    }else{
                        ElMessage.error(res.StatusMessage)
                    }
                    })
                })
                .catch(() => {});
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
        onSearch(){
            this.pageNum = 1;
            this.getTableData();
        },
        getTableData() {
            this.loading = true;
            serverApi.monitorNoticeLogList({
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