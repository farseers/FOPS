<template>
    <div class="system-user-container layout-padding">
        <el-card>
            <div class="system-user-search mb15">
                <el-button size="default" type="success" class="ml10" @click="onSearch()">
                    <el-icon>
						<ele-Search />
					</el-icon>
                    查询</el-button>
                    
            </div>
            <el-table :data="tableData" v-loading="loading" style="width: 100%">
                <el-table-column type="index" label="序号" width="60" />
                <el-table-column prop="AppName" label="项目名称" show-overflow-tooltip></el-table-column>
                <el-table-column label="时间类型" show-overflow-tooltip>
                <template #default="scope">
                    <span v-if="scope.row.TimeType==0">小时</span>
                    <span v-if="scope.row.TimeType==1">天</span>
                </template>
                </el-table-column>
                <el-table-column prop="StartTime" label="开始时间" show-overflow-tooltip></el-table-column>
                <el-table-column prop="EndTime" label="结束时间" show-overflow-tooltip></el-table-column>
                <el-table-column prop="Comparison" label="比较方式" show-overflow-tooltip></el-table-column>
                <el-table-column label="监控键值" show-overflow-tooltip>
                <template #default="scope">
                    <span>{{scope.row.KeyName}}</span>
                    <span v-show="scope.row.KeyName">:</span>
                    <span>{{scope.row.KeyValue}}</span>
                </template>
                </el-table-column>
                <el-table-column label="是否启用" show-overflow-tooltip>
                    <template #default="scope">
                        <el-tag size="small" v-if="scope.row.Enable" type="success">是</el-tag>
                        <el-tag size="small" v-else type="danger">否</el-tag>
                </template>
                </el-table-column>
                <el-table-column prop="NoticeIds" label="关联人" show-overflow-tooltip>
                    <template #default="scope">
                        <span style="margin: 3px;" v-for="item,index in scope.row.NoticeList" :key="index">{{item.Name}}、</span>
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
        <GuleDialog ref="editInfo" @search="getTableData" :p_list="p_list"/>
    </div>
</template>
<script>
import { fopsApi } from "/@/api/fops";
import InitPagination from '/src/views/components/InitPagination.vue';
import GuleDialog from './guleDialog.vue';
import { ElMessageBox, ElMessage } from 'element-plus';
const serverApi = fopsApi();
export default {
    components: { InitPagination,GuleDialog },
    data() {
        return {
            tableData: [],
            loading: false,
            p_list:[],//关联人
            pages: {
                pageNum: 1,
                pageSize: 10,
                total: 0
            }

        }
    },
    mounted(){
        this.getTableData()
        this.get_notie()
    },
    methods: {
        get_notie(){ //获取关联人
        serverApi.monitorNoticeList({
                "pageSize": 10000,
                "pageIndex": 1
            }).then((d)=>{
                if(d.Status){
                    const { List } = d.Data;
                    this.p_list = List.filter(item => item.Enable === true)
                }else{
                    ElMessage.error(d.StatusMessage);
                }
            }).catch(e=>{
                ElMessage.error('网络错误');
            })
      },
        set_edit(row){
            if(row.Id){
                this.$refs.editInfo && this.$refs.editInfo.info(row.Id,this.p_list)
            }
        },
        set_del(row){ //删除
            let str = '确定删除此数据?'
            if(row.AppName){
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
                    serverApi.monitorDelRule({"id":row.Id}).then(function (res){
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
            serverApi.monitorRuleList({
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