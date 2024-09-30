<template>

  <div class="system-role-dialog-container">
    <el-dialog :title="title" v-model="isShowDialog" width="800px">
      
      <el-form ref="ruleFormRef" :model="infoRow" :rules="rules">
          <el-form-item label="应用名称" prop="AppNames">
          <el-select v-model="infoRow.AppNames" filterable placeholder="请选择" multiple style="flex: 1;">
            <el-option v-for="item in m_list" :key="item.AppName" :label="item.AppName" :value="item.AppName" />
          </el-select>
        </el-form-item>
        <el-form-item label="生效时间">
          <el-form-item  prop="TimeType" style="margin-right: 10px;">
          <el-radio-group v-model="infoRow.TimeType" @change="radioTimeTypeChange()">
            <el-radio :label="0">小时</el-radio>
            <el-radio :label="1">天</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item  prop="daterange" v-show="infoRow.TimeType == 1"  style="flex: 1;">
          <el-date-picker v-model="infoRow.daterange" type="daterange" range-separator="To"
            start-placeholder="开始时间" end-placeholder="结束时间" value-format="YYYY-MM-DD" format="YYYY-MM-DD"/>
        </el-form-item>
        <el-form-item prop="daterange" v-show="infoRow.TimeType == 0"  style="flex: 1;" >
          <el-time-picker
            v-model="infoRow.daterange"
            is-range
            format="HH:mm:ss"
            value-format="HH:mm:ss"
            range-separator="To"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>
        </el-form-item>
        
        
        
        <el-form-item style="display: flex;">
          
          <div style="flex: 1;">
            <el-form-item label="监控键值" required>
            <el-form-item prop="KeyName" style="flex: 1;">
              <el-input v-model="infoRow.KeyName" placeholder="键"/>
            </el-form-item>
            <div style="width: 120px;padding:0 5px;">
            <el-form-item prop="Comparison">
            <el-select v-model="infoRow.Comparison" filterable placeholder="比较方式" style="flex: 1;">
              <el-option v-for="item in t_list" :key="item.Key" :label="item.Value" :value="item.Value" />
            </el-select>
          </el-form-item>
          </div>
            <el-form-item prop="KeyValue" style="flex: 1;">
              <el-input v-model="infoRow.KeyValue" placeholder="值"/>
            </el-form-item>
          </el-form-item>
          </div>
        </el-form-item>

        <el-form-item label="是否启用">
          <el-switch v-model="infoRow.Enable" active-text="启用" inactive-text="停用"
            style="--el-switch-on-color: #13ce66; --el-switch-off-color: #ff4949" />
        </el-form-item>
        <el-form-item label="关联人ID">
          <div style="display: flex;">
            <div>
              <el-button size="small" type="success" @click="setPers()">
                <el-icon>
                  <ele-Edit />
                </el-icon>
                设置</el-button>
            </div>
            <div style="flex: 1;">
              <span style="margin: 3px;" v-for="item, index in infoRow.NoticeIds" :key="index"
                v-text="set_name(item, index)"></span>
            </div>
          </div>
        </el-form-item>
        
        <el-form-item label="模版"><el-input v-model="infoRow.TipTemplate" /></el-form-item>
        <el-form-item label="备注"><el-input v-model="infoRow.Remark" /></el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="onCancel" size="default">取 消</el-button>
          <el-button type="primary" @click="onSubmit" size="default">保 存</el-button>
        </span>
      </template>
    </el-dialog>
    <el-dialog title="设置关联人" v-model="isTransfer" width="650px">
      <el-transfer style="text-align: center;" :titles="['用户列表', '选中列表']" filterable :filter-method="filterMethod"
        v-model="ck_list" :props="{
          key: 'Id',
          label: 'Name',
        }" :data="p_list" />
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
const validators = (e, s) => {
  if (s && s.length == 2) {
    if(s[0] && s[1]){
      return true
    }else{
      return false
    }
   
  } else {
    return false
  }
}
const validatorName = (e,s)=>{
  if (s && s.length > 0) {
    return true
  } else {
    return false
  }
}
const defaultRow = {
       'AppNames':[],
        "daterange": [],
        "timerange": [],
        "Id": null,
        "AppName": "",
        "TimeType": 1,
        "Comparison": "",
        "KeyName": "",
        "KeyValue": "",
        "TipTemplate":'',
        "Remark": "",
        "NoticeIds": [],
        "NoticeList": [],
        "StartDate":'00:00:00',
        "EndDate":'23:59:59',
        "StartDay":'',
        "EndDay":'',
        "Enable":true
}
export default {
  data() {
    return {
      rules: {
        daterange: [{ required: true, trigger: 'change', type: 'date', message: '请选择时间', validator: validators }],
        AppNames: [{ required: true, trigger: 'change', message: '请选择应用名称',validator: validatorName }],
        TimeType: [{ required: true, trigger: 'change', }],
        Comparison: [{ required: true, trigger: 'blur', message: '请选择比较方式' }],
        KeyName: [{ required: true, trigger: 'blur', message: '请输入键值' }],
        KeyValue: [{ required: true, trigger: 'blur', message: '请输入键值' }],
      },
      ck_list: [],//穿梭框选中
      p_list: [],//关联人列表
      m_list: [],//项目列表
      t_list: [],//比较方式
      title: '编辑规则',
      isShowDialog: false,
      isTransfer: false,//设置关联人
      infoRow: {
        ...defaultRow
      },
    }
  },
  methods: {
    getNow(){
      const now = new Date();
      const dateTime = now.toLocaleString().split(' ')[0].split('/');
      const year = parseInt(dateTime[0]);
      const month = parseInt(dateTime[1]);
      const day = parseInt(dateTime[2]);
      return year + '-' + month + '-' + day
    },
    daterangeChange(){
      const daterange = this.infoRow.daterange
      if(daterange && daterange.length > 0){
        if(this.infoRow.TimeType == 0){ //小时
          this.infoRow.StartDate = daterange[0];
          this.infoRow.EndDate = daterange[1];
        }else{
          this.infoRow.StartDay = daterange[0];
          this.infoRow.EndDay = daterange[1];
        }
        }
      
    },
    
    radioTimeTypeChange(){
      const { StartDate, EndDate,StartDay,EndDay} = this.infoRow
        if(this.infoRow.TimeType == 0){ //小时
          this.infoRow.daterange = [StartDate,EndDate]
        }else{
          this.infoRow.daterange = [StartDay,EndDay]
        }
    },
    set_name(id, i) {
      const row = this.p_list.find(d => {
        return d.Id == id
      })

      if (row) {
        let str = row.Name + ' /';
        if (this.infoRow.NoticeIds.length - 1 == i) { str = row.Name }
        return str
      }
    },
    filterMethod(query, item) {
      return item.Name.toLowerCase().includes(query.toLowerCase())
    },
    tranSave() { //关联人保存 
      this.infoRow.NoticeIds = this.ck_list;
      this.isTransfer = false;
    },
    tranCancel() {

      this.isTransfer = false;
    },
    setPers() {//设置关联人
      this.isTransfer = true;
    },
    init() {
      this.ck_list = [];
      this.p_list = [];
      this.m_list = [];
      this.t_list = [];
      this.infoRow = {
        ...defaultRow
      }
      this.radioTimeTypeChange()
      this.$refs.ruleFormRef && this.$refs.ruleFormRef.resetFields()
    },
    onCancel() {
      this.isShowDialog = false;
      this.init()
      this.$emit('search')
    },
    get_param() {
      let param = { ...this.infoRow };
      const daterange = param.daterange;
      if (daterange && daterange.length > 0) {
        if(param.TimeType == 0){ //小时
          param.StartDate = daterange[0];
          param.EndDate = daterange[1];
        }else{
          param.StartDay = daterange[0];
          param.EndDay = daterange[1];
        }
      }
      const AppNames = param.AppNames;
      if(AppNames && AppNames.length>0){
        param.AppName = AppNames.join(',')
      }
      return param
    },
    onSubmit() {
      this.$refs.ruleFormRef && this.$refs.ruleFormRef.validate((valid) => {
        if (valid) {
          const param = this.get_param()
          serverApi.monitorSaveRule(param).then(d => {
            let { Status, StatusMessage } = d;
            if (Status) {
              this.onCancel()
            } else {
              ElMessage.error(StatusMessage)
            }
          })
        }
      })
    },
    info(id, list1, list2, list3) {
      this.init()
      if (list1) {
        this.p_list = [...list1];
      }
      if (list2) {
        this.m_list = [...list2]
      }
      if (list3) {
        this.t_list = [...list3]
      }
       this.title = '新增规则'
      if (id) {
        this.title = '编辑规则'
        serverApi.monitorInfoRule({ id: id }).then(d => {
          let { Data, Status, StatusMessage } = d;
          if (Status) {
            this.infoRow = { ...Data, daterange: [],AppNames:[] }
            if(this.infoRow.TimeType == 0){ //小时
                this.infoRow.daterange = [Data.StartDate||'00:00:00',Data.EndDate||'23:59:59']
            }else{
                 this.infoRow.daterange = [Data.StartDay,Data.EndDay]
            }
            const AppName = Data.AppName;
            if(AppName.indexOf(',')==-1){
              this.infoRow.AppNames = [AppName]
            }else{
              this.infoRow.AppNames = AppName.split(',')
            }
            this.ck_list = [...this.infoRow.NoticeIds]
            this.isShowDialog = true;
          } else {
            ElMessage.error(StatusMessage)
          }
        })
      } else {
        this.isShowDialog = true;
      }

    }
  }
}
</script>
