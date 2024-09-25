<template>
  <div class="system-role-dialog-container">
    <el-dialog :title="title" v-model="isShowDialog" width="800px" @close="close">
      <el-form-item label="项目名称"><el-input v-model="infoRow.AppName" /></el-form-item>
      <el-form-item label="时间类型">
        <el-radio-group v-model="infoRow.TimeType">
          <el-radio :label="0">小时</el-radio>
          <el-radio :label="1">天</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="开始时间">
        <el-date-picker v-model="infoRow.daterange" type="datetimerange" range-separator="To" start-placeholder="Start date"
          end-placeholder="End date" />
      </el-form-item>
      <el-form-item label="比较方式">
        <el-input v-model="infoRow.Comparison" />
      </el-form-item>
      <el-form-item label="监控键值">
        <el-input v-model="infoRow.KeyName" style="flex: 1;" />
        <span style="width: 15px;text-align: center;"> : </span>
        <el-input style="flex: 1;" v-model="infoRow.KeyValue" />
      </el-form-item>
      <el-form-item label="是否启用">
        <el-switch v-model="infoRow.Enable" active-text="启用" inactive-text="关闭"
          style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949" />
      </el-form-item>
      <el-form-item label="关联人ID">
        <div style="display: flex;">
          <div >
            <el-button size="small" type="success" @click="setPers()">
          <el-icon>
            <ele-Edit />
          </el-icon>
          设置</el-button>
          </div>
         <div style="flex: 1;">
          <span style="margin: 3px;" v-for="item,index in infoRow.NoticeList" :key="index">{{item.Name}}、</span>
         </div>
        </div>
      </el-form-item>
      <el-form-item label="备注"><el-input v-model="infoRow.Remark" /></el-form-item>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="onCancel" size="default">取 消</el-button>
          <el-button type="primary" @click="onSubmit" size="default">保 存</el-button>
        </span>
      </template>
    </el-dialog>
    <el-dialog title="设置关联人" v-model="isTransfer" width="650px">
        <el-transfer style="text-align: center;"
                :titles="['用户列表', '选中列表']"
                filterable
                :filter-method="filterMethod"
                v-model="ck_list"
                :props="{
                key: 'Id',
                label: 'Name',
                }"
                :data="p_list"
            />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="tranCancel" size="default">取 消</el-button>
          <el-button type="primary" @click="tranSave" size="default">保 存</el-button>
        </span>
      </template>
    </el-dialog>

  </div>
</template>


<script>

import { fopsApi } from "/@/api/fops";
import { ElMessage } from 'element-plus';
const serverApi = fopsApi();
export default {
  data() {
    return {
      ck_list: [],//穿梭框选中
      p_list: [],//关联人列表
      title: '编辑规则',
      isShowDialog: false,
      isTransfer: false,//设置关联人
      infoRow: {
        "daterange": [],
        "Id": null,
        "AppName": "",
        "TimeType": null,
        "StartTime": "",
        "EndTime": "",
        "Comparison": "",
        "KeyName": "",
        "KeyValue": "",
        "Remark": "",
        "NoticeIds": [],
        "NoticeList":[]
      },
    }
  },
  methods: {
    filterMethod (query, item){
            return item.Name.toLowerCase().includes(query.toLowerCase())
       } ,
    tranSave(){ //关联人保存 
      this.infoRow.NoticeIds = this.ck_list;
      const param = this.get_param()
      serverApi.monitorSaveRule(param).then(d => {
        let { Status, StatusMessage } = d;
        if (Status) {
          this.isTransfer = false;
          this.info(this.infoRow.Id)
        } else {
          ElMessage.error(StatusMessage)
        }
      })
    
      
    },
    tranCancel(){
     
      this.isTransfer = false;
    },
    setPers() {//设置关联人
      this.isTransfer = true;
    },
    close(){
      this.$emit('search')
    },
    onCancel() {
      this.isShowDialog = false;
      this.ck_list = [];
      this.p_list = [];
      this.infoRow = {
        "daterange": [],
        "Id": null,
        "AppName": "",
        "TimeType": null,
        "StartTime": "",
        "EndTime": "",
        "Comparison": "",
        "KeyName": "",
        "KeyValue": "",
        "Remark": "",
        "NoticeIds": [],
        "NoticeList":[]
      }
      this.$emit('search')
    },
    get_param(){
      let param = { ...this.infoRow };
      const daterange = param.daterange;
      if (daterange && daterange.length > 0) {
        param.StartTime = daterange[0];
        param.EndTime = daterange[1];
      }
      return param
    },
    onSubmit() {
      const param = this.get_param()
      serverApi.monitorSaveRule(param).then(d => {
        let { Status, StatusMessage } = d;
        if (Status) {
          this.onCancel()
        
        } else {
          ElMessage.error(StatusMessage)
        }
      })
    },
    info(id, list) {
      if(list){
        this.p_list = [...list];
      }
     
      serverApi.monitorInfoRule({ id: id }).then(d => {
        let { Data, Status, StatusMessage } = d;
        if (Status) {
          this.infoRow = { ...Data, daterange: [Data.StartTime, Data.EndTime] }
          this.ck_list = [...this.infoRow.NoticeIds]
          this.isShowDialog = true;
        } else {
          ElMessage.error(StatusMessage)
        }
      })
    }
  }
}
</script>
