<template>
    <div class="system-user-container layout-padding">
        <el-card>
            <div class="system-user-search mb15">
                <el-button size="default" type="warning" class="ml10" @click="set_add()">
                    <el-icon>
						<ele-Plus />
					</el-icon>
                    新增用户</el-button>
                <el-button size="default" type="success" class="ml10" @click="onSearch()">
                    <el-icon>
						<ele-Search />
					</el-icon>
                    查询</el-button>
                    
            </div>
            <el-table :data="tableData" v-loading="loading" style="width: 100%">
                <el-table-column type="index" label="序号" width="60" />
                <el-table-column prop="NoticeType" label="通知类型" show-overflow-tooltip>
                    <template #default="scope">
                            <span v-text="set_type(scope.row.NoticeType)"></span>
                    </template>
                </el-table-column>
                <el-table-column prop="Name" label="姓名" show-overflow-tooltip></el-table-column>
                <el-table-column prop="Email" label="邮箱" show-overflow-tooltip></el-table-column>
                <el-table-column prop="Phone" label="号码" show-overflow-tooltip></el-table-column>
                <el-table-column prop="ApiKey" label="接口Key" show-overflow-tooltip></el-table-column>
                <el-table-column label="是否启用" show-overflow-tooltip>
                    <template #default="scope">
                        <el-tag size="small" v-if="scope.row.Enable" type="success">是</el-tag>
                        <el-tag size="small" v-else type="danger">否</el-tag>
                </template>
                </el-table-column>
                <el-table-column prop="Remark" label="备注" show-overflow-tooltip></el-table-column>
                <el-table-column label="操作" width="140px" fixed="right" align="center">
                    <template #default="scope">
                        <el-button @click="set_edit(scope.row)" type="primary" plain size="small">编辑</el-button>
                        <el-button @click="set_del(scope.row)" type="danger" plain size="small">删除</el-button>
                    </template>
            </el-table-column>
            </el-table>
            <InitPagination @sizeChange="onHandleSizeChange" @currentChange="onHandleCurrentChange" :pages="pages" />
        </el-card>
        <noticeDialog ref="editInfo" :typeList="typeList" @search="getTableData"/>
    </div>
</template>
<script>

import { fopsApi } from "/@/api/fops";
import InitPagination from '/src/views/components/InitPagination.vue';
import noticeDialog from './noticeDialog.vue';
import { ElMessageBox, ElMessage } from 'element-plus';
// Id         int64           // 主键
//     NoticeType noticeType.Enum // 通知类型
//     Email      string          // 邮箱
//     Phone      string          // 号码
//     ApiKey     string          // 接口Key
//     Remark     string          // 备注
//     Enable     bool            // 是否启用

const serverApi = fopsApi();
export default {
    components: { InitPagination,noticeDialog },
    data() {
        return {
            tableData: [],
            typeList:[],//通知类型
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
        info(){
            serverApi.monitorNoticeTypeList({}).then(d=>{
                const { Data,Status } = d;
                if(Status){
                    this.typeList = [...Data];
                }
            })
        },
        set_add(){
            this.$refs.editInfo && this.$refs.editInfo.info()
        },
        set_edit(row){
            if(row.Id){
                this.$refs.editInfo && this.$refs.editInfo.info(row.Id)
            }
        },
        set_del(row){ //删除
            let str = '确定删除此数据?'
            const _this = this;
            ElMessageBox.confirm(`${str}`, '提示', {
                confirmButtonText: '确认',
                cancelButtonText: '取消',
                type: 'warning',
            })
                .then(() => {
                    // 删除逻辑
                    serverApi.monitorDelNotice({"id":row.Id}).then(function (res){
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
            serverApi.monitorNoticeList({
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