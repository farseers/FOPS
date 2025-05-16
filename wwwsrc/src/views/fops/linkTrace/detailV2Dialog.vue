<template>
	<div class="system-user-dialog-container">
		<el-dialog :title="state.dialog.title" v-model="state.dialog.isShowDialog" width="70%">
			<el-form ref="gitDialogFormRef" size="default" label-width="100px">
				<el-row :gutter="35">
					<el-col :xs="24" :sm="24" :md="24" :lg="24" :xl="24" class="mb20" v-if="state.tableData.length>0">
            <div>
              应用名称：<el-tag size="small">{{state.AppName}}</el-tag>，
              应用ID：{{state.AppId}}，
              应用IP：{{state.AppIp}}，
              请求时间：{{state.CreateAt}}，
              整体耗时：
                <el-tag size="small" v-if="state.UseTs > 100000000" type="danger">{{state.UseDesc}}</el-tag>
                <el-tag size="small" v-else-if="state.UseTs > 50000000" type="warning">{{state.UseDesc}}</el-tag>
                <el-tag size="small" v-else-if="state.UseTs > 1000000">{{state.UseDesc}}</el-tag>
                <el-tag size="small" v-else type="success">{{state.UseDesc}}</el-tag>
            </div>
            <!--webapi-->
            <div class="mt10">
              <span v-if="state.TraceType == 0">
                <el-tag size="small">{{state.WebStatusCode}}</el-tag> {{state.WebRequestIp}} <el-tag type="success" size="small">{{state.WebMethod}}</el-tag>
                <el-tag v-if="state.WebContentType!=''" type="info" size="small">{{state.WebContentType}}</el-tag>{{state.WebPath}}
                <el-button style="margin-left: 20px" type="primary" @click="onShow()" size="small">查看报文</el-button>
              </span>
              <!--MqConsumer--> <!--QueueConsumer-->
              <span v-else-if="state.TraceType == 1 || state.TraceType == 2">
                {{state.ConsumerServer}}
                <el-tag v-if="state.ConsumerRoutingKey !=''" size="small">{{state.ConsumerRoutingKey}}</el-tag>
                <br v-if="state.ConsumerRoutingKey !=''" />
                {{state.ConsumerQueueName}}
              </span>
              <!--FSchedule--> <!--Task-->
              <span v-else-if="state.TraceType == 3 || state.TraceType == 4">
                <el-tag v-if="state.tgn >0" size="small">任务组Id：{{state.tgn}}</el-tag>
                <el-tag v-if="state.TaskId >0" size="small" type="success">任务Id：{{state.TaskId}}</el-tag>
                {{state.TaskName}}
              </span>
              <!--webSocket-->
              <span v-else-if="state.TraceType == 7">
                {{state.WebRequestIp}} <el-tag type="success" size="small">ws</el-tag>
                {{state.WebPath}}
                <el-button style="margin-left: 20px" type="primary" @click="onShow()" size="small">查看报文</el-button>
              </span>
              <el-button style="margin-left: 20px" size="small" type="success" @click="showLog()">查看日志</el-button>
              <div v-if="state.Exception!=null" class="mt5">
              <el-tag type="danger">
                异常：{{state.Exception.ExceptionDetails[0].ExceptionCallFile}}:{{state.Exception.ExceptionDetails[0].ExceptionCallLine}} {{state.Exception.ExceptionDetails[0].ExceptionCallFuncName}}
                {{state.Exception.ExceptionMessage}}
              </el-tag>
              </div>
            </div>
            <div :style="{'width':'95%','white-space': 'nowrap'}">
            <ul class="custom-list mt10">
              <li style="height: 35px;padding: 10px 0;position: relative;display: flex;">
                <!-- <el-tooltip v-for="(info, index) in state.tableData" :key="index"  effect="dark" content="Top Left 提示文字" placement="top-start">
                  <span style="height: 5px;flex: 1;background-color: red;"></span>
          </el-tooltip> -->
          <span v-for="(info, index) in state.timeDatas" :key="index" 
                :style="{'left':info.StartRate+'%','position':'absolute'}">
                  <el-tag  size="small" style="max-width: 65px;" type="success">{{info.StartTs / 1000}}ms</el-tag>
                </span> 
                  <!-- <span style="float:right;margin-right: 10px">
                  <el-tag size="small" type="success">{{state.UseTs / 1000000}}ms</el-tag>
                </span> -->
              </li>
                <!-- <span v-for="(info, index) in state.timeDatas" :key="index" 
                :style="{'left':info.StartRate+'%','position':'absolute'}">
                  <el-tag  size="small" style="max-width: 65px;" type="success">{{info.StartTs / 1000}}ms</el-tag>
                </span> -->

              <!--详情-->
              <li style="clear: both;padding:2px 0;height:21px" v-for="(info, index) in state.tableData" :key="index">
                <div>
                  <span style="float:left;position:absolute;margin: 0 5px;color:#8c8b8b">{{index + 1}}</span>
                  <span :style="{'margin-left':info.StartRate+'%','float':'left','position':'relative','width':'100%','cursor':'pointer'}"  @click="copyText(info.Desc)"  :title="info.Desc">
                    <div class="el-progress el-progress--line is-exception el-progress--text-inside" role="progressbar" aria-valuenow="50" aria-valuemin="0" aria-valuemax="100" :style="{'width':info.UseRate+'%'}">
                      <div class="el-progress-bar">
                          <div class="el-progress-bar__inner" :style="{'height': '21px','width': '100%', 'animation-duration': '3s','text-align': 'left','background-color':'rgb('+info.Rgba+')'}">
                            <div class="el-progress-bar__innerText" style="color:#181818">
                              <el-tag size="small" style="margin-right: 5px;">{{info.AppName}}</el-tag>
                              <el-tag size="small" style="margin-right: 5px;" v-if="info.Exception!=null" :title="info.Exception.ExceptionMessage" type="danger">【异常】</el-tag>
                              <span v-html="info.Caption"></span>
                              <span v-if="index > 0 && info.UseTs > 0">，耗时：
                                <el-tag size="small" v-if="info.UseTs > 100000000" type="danger">{{info.UseDesc}}</el-tag>
                                <el-tag size="small" v-else-if="info.UseTs > 50000000" type="warning">{{info.UseDesc}}</el-tag>
                                <el-tag size="small" v-else-if="info.UseTs > 1000000">{{info.UseDesc}}</el-tag>
                                <el-tag size="small" v-else type="success">{{info.UseDesc}}</el-tag>
                              </span>
                            </div>
                        </div>
                      </div>
                    </div>
                  </span>
<!--                  <span v-if="info.Exception!=null">异常：{{friendlyJSONstringify(info.Exception)}}</span>-->
<!--                  <span v-else></span>-->
                </div>
              </li>
            </ul>
            </div>
					</el-col>
          <el-col style="text-align: center" v-else>
            <span style="width: 100%;">暂无数据</span>
          </el-col>
				</el-row>
			</el-form>
<!--      <el-table :data="state.tableData" v-loading="state.loading" style="width: 100%">-->
<!--        <el-table-column prop="UseDesc" label="时间" show-overflow-tooltip></el-table-column>-->
<!--        <el-table-column width="250px" label="内容" show-overflow-tooltip>-->
<!--          <template #default="scope">-->
<!--            <el-tag size="small">{{scope.row.AppName}}</el-tag><span>{{scope.row.Desc}} {{scope.row.Caption}}</span>-->
<!--          </template>-->
<!--        </el-table-column>-->
<!--        <el-table-column width="120px" prop="AppIp" label="应用IP" show-overflow-tooltip></el-table-column>-->
<!--        <el-table-column width="200px" label="异常" show-overflow-tooltip>-->
<!--          <template #default="scope">-->
<!--            <el-tag size="small" v-if="scope.row.Exception!=null" type="danger">{{scope.row.Exception.ExceptionCallFile}}:{{scope.row.Exception.ExceptionCallLine}} {{scope.row.Exception.ExceptionCallFuncName}}</el-tag><br  v-if="scope.row.Exception!=null">-->
<!--            <el-tag size="small" v-if="scope.row.Exception!=null" type="danger">{{scope.row.Exception.ExceptionMessage}}</el-tag>-->
<!--            <el-tag size="small" v-else type="info">无</el-tag>-->
<!--          </template>-->
<!--        </el-table-column>-->
<!--      </el-table>-->
		</el-dialog>

    <logDialog ref="logDialogRef"  />
    <showDialog ref="showDialogRef"  />

	</div>
</template>

<script setup lang="ts" name="fopsLinkDetailV2Dialog">
import {defineAsyncComponent, reactive, ref} from 'vue';
import {fopsApi} from "/@/api/fops";
import { ElMessageBox, ElMessage } from 'element-plus';
import {friendlyJSONstringify} from "@intlify/shared";
// 定义变量内容
import commonFunction from '/@/utils/commonFunction';
const { copyText } = commonFunction();
// 引入 api 请求接口
const serverApi = fopsApi();
// 定义子组件向父组件传值/事件
const emit = defineEmits(['refresh']);
const showDialog = defineAsyncComponent(() => import('/src/views/fops/sql/httpDialog.vue'));
const logDialog = defineAsyncComponent(() => import('/src/views/fops/log/logV2Dialog.vue'));
// 定义变量内容
const gitDialogFormRef = ref();
const showDialogRef = ref();
const logDialogRef = ref();
const state = reactive({
	ruleForm: {},
  loading:false,
  timeDatas:[
  {
    Rgba:'',
    AppId:0,
    AppIp:'',
    AppName:'',
    StartTs:0,
    StartRate:0,
    UseTs:0,
    UseRate:0,
    Caption:'',
    Desc:'',
    UseDesc:'',
    Exception:'',
  }
  ],
  tableData:[{
    Rgba:'',
    AppId:0,
    AppIp:'',
    AppName:'',
    StartTs:0,
    StartRate:0,
    UseTs:0,
    UseRate:0,
    Caption:'',
    Desc:'',
    UseDesc:'',
    Exception:'',
  }],
  TraceId:'',
  Rgba:'',
  AppId:0,
  AppIp:'',
  AppName:'',
  Desc:'',
  Caption:'',
  UseDesc:'',
  UseTs:0,
  TraceType:0,
  WebStatusCode:0,
  WebRequestIp:'',
  WebMethod:'',
  WebContentType:'',
  WebPath:'',
  WebHeaders:'',
  WebRequestBody:'',
  WebResponseBody:'',
  tgn:0,
  TaskId:0,
  TaskName:'',
  ConsumerServer:'',
  ConsumerRoutingKey:'',
  ConsumerQueueName:'',
  CreateAt:'',
  Exception:{},
	dialog: {
		isShowDialog: false,
		type: '',
		title: '',
		submitTxt: '',
	},
  traceInfo:{},
  spacePx:100,
});

// 打开弹窗
const openDialog = (row2: any) => {
  state.loading=true
  //state.ruleForm = row;
  state.dialog.title = '链路追踪详情(TraceId：'+row2.tid+')';
  state.traceInfo=row2

  // 详情
  serverApi.linkTraceInfo(row2.tid).then(function (res){
    state.loading=false
    if (res.Status){
      // 计算宽度
      if (res.Data.List.length <= 5) {
        state.spacePx = 200
      } else if (res.Data.List.length <= 10) {
        state.spacePx = 150
      } else{
        state.spacePx = 60
      }
      let originalArray = [...res.Data.List]
      const set_arr = function(StartRate,index){
        let row = null
          for(var i=0;i<originalArray.length;i++){
            if (i>=index){
              var n_StartRate = originalArray[i].StartRate
              if(n_StartRate >= StartRate+6){ 
                row = {
                  StartRate: n_StartRate,
                  StartTs: originalArray[i].StartTs,
                  index:i
                }
                break
              }
            }
          }
          return row
      }
      let crr = [];
      const items = function(item,i){
          const StartRate = item.StartRate;
          crr.push({
              StartRate:item.StartRate,
              StartTs:item.StartTs,
            })
          const row = set_arr(StartRate,i);
          if(row){
            var index = row.index + 1;
            if(originalArray[index]){
              items(originalArray[index],index)
            }
          
          }
      }
      if(originalArray && originalArray.length>0){
        items(originalArray[1],1)
      }
      
      // 绑定数据
      state.timeDatas = crr
      state.tableData=res.Data.List
      state.AppId=res.Data.Entry.aid
      state.AppIp=res.Data.Entry.aip
      state.AppName=res.Data.Entry.an
      state.Desc=res.Data.Entry.Desc
      state.UseDesc=res.Data.Entry.ud
      state.TraceId=res.Data.Entry.tid
      state.UseTs=res.Data.Entry.ut
      state.TraceType=res.Data.Entry.tt
      state.WebStatusCode=res.Data.Entry.wsc
      state.WebRequestIp=res.Data.Entry.wip
      state.WebMethod=res.Data.Entry.wm
      state.WebContentType=res.Data.Entry.wct
      state.WebPath=res.Data.Entry.wp
      state.WebHeaders=res.Data.Entry.wh
      state.WebRequestBody=res.Data.Entry.wrb
      state.WebResponseBody=res.Data.Entry.wpb
      state.tgn=res.Data.Entry.tgn
      state.TaskId=res.Data.Entry.tid
      state.TaskName=res.Data.Entry.tn
      state.ConsumerServer=res.Data.Entry.cs
      state.ConsumerRoutingKey=res.Data.Entry.cr
      state.ConsumerQueueName=res.Data.Entry.cq
      state.CreateAt=res.Data.Entry.ca
      state.Exception=res.Data.Entry.e
    }
  })
	state.dialog.isShowDialog = true;
};
// 关闭弹窗
const closeDialog = () => {
	state.dialog.isShowDialog = false;
};
const onShow=()=>{
  showDialogRef.value.openDialog(2,state);
}
const showLog=()=>{
  logDialogRef.value.openDialog(state.traceInfo);
}

// 取消
const onCancel = () => {
	closeDialog();
};

// 暴露变量
defineExpose({
	openDialog,
});
</script>

<style>
textarea{
  height: 220px;
}

/* 基本样式 */
.custom-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

/* 每个列表项的样式 */
.custom-list li {
  padding: 2px 15px;
  margin-bottom: 8px;
  border-radius: 5px;
  background-color: #f2f2f2;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: background-color 0.3s ease;
}

/* 悬停时的样式 */
.custom-list li:hover {
  background-color: #e0e0e0;
}
/* 基本样式，根据需要自定义样式 */
.timeline {
  position: relative;
  padding: 20px 0;
}

.timeline-event {
  display: flex;
  align-items: flex-start;
  margin-bottom: 20px;
}

.timeline-event-marker {
  width: 20px;
  height: 20px;
  background-color: #333;
  border-radius: 50%;
  margin-right: 10px;
}

.timeline-event-content {
  border: 1px solid #ccc;
  padding: 10px;
  border-radius: 5px;
  background-color: #f9f9f9;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}
.el-dialog__body{
  overflow:auto!important;
}
</style>
