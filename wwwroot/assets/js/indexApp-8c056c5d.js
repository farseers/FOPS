import{_ as ot}from"./index-74aeefad.js";import{f as lt}from"./index-2149a754.js";import{E as B,a as k}from"./element-plus-7da03efa.js";import{k as $,aN as nt,e as st,$ as it,f as rt,i as pt,_ as dt,ah as r,ar as ut,o as i,c as b,W as s,Q as l,a as c,F as v,a9 as M,U as u,P as d,R as ct,T as w,u as gt,V as h}from"./@vue-d9027515.js";import"./pinia-be7e992d.js";import"./vue-demi-71ba0ef2.js";import"./vue-router-a2fac16f.js";import"./vue-i18n-e0b9a81d.js";import"./@intlify-9e8a497c.js";import"./source-map-7d7e1c08.js";import"./vue-b46b26c7.js";import"./js-cookie-cf83ad76.js";import"./@element-plus-96a31696.js";import"./nprogress-8d2808ea.js";import"./axios-707ed124.js";import"./qs-bbfcbdf3.js";import"./side-channel-eafc5c70.js";import"./get-intrinsic-b9397c9a.js";import"./has-symbols-e8f3ca0e.js";import"./function-bind-22e7ee79.js";import"./has-26d28e02.js";import"./call-bind-9ceb8f5b.js";import"./object-inspect-7b08a46e.js";import"./mitt-f7ef348c.js";import"./vue-grid-layout-0ee01be3.js";import"./lodash-es-36eb724a.js";import"./@vueuse-2d9216d7.js";import"./@ctrl-f8748455.js";import"./dayjs-6aac7dfa.js";import"./async-validator-dee29e8b.js";import"./memoize-one-297ddbcb.js";import"./escape-html-1d60d822.js";import"./normalize-wheel-es-ed76fb12.js";import"./@floating-ui-463e90e0.js";import"./requestFSGet-61715c66.js";const mt={class:"system-user-container layout-padding"},ft={class:"system-user-search mb15"},_t={class:"flex-warp",style:{background:"#e0e0e0"}},Dt=["onClick"],bt={class:"appItem"},ht={class:"appItem"},yt={class:"appItem"},kt={class:"appItem"},vt={class:"appItem"},wt={class:"appItem"},It=c("h3",{style:{padding:"5px"}},"构建队列",-1),Lt=["innerHTML"],St=$({name:"fopsApp"}),pe=$({...St,setup(Ct){const g=lt(),E=nt(()=>ot(()=>import("./dialog-cbbb38f9.js"),["assets/js/dialog-cbbb38f9.js","assets/js/index-2149a754.js","assets/js/index-74aeefad.js","assets/js/@vue-d9027515.js","assets/js/pinia-be7e992d.js","assets/js/vue-demi-71ba0ef2.js","assets/js/vue-router-a2fac16f.js","assets/js/vue-i18n-e0b9a81d.js","assets/js/@intlify-9e8a497c.js","assets/js/source-map-7d7e1c08.js","assets/js/vue-b46b26c7.js","assets/js/js-cookie-cf83ad76.js","assets/js/@element-plus-96a31696.js","assets/js/nprogress-8d2808ea.js","assets/css/nprogress-8b89e2e0.css","assets/js/axios-707ed124.js","assets/js/qs-bbfcbdf3.js","assets/js/side-channel-eafc5c70.js","assets/js/get-intrinsic-b9397c9a.js","assets/js/has-symbols-e8f3ca0e.js","assets/js/function-bind-22e7ee79.js","assets/js/has-26d28e02.js","assets/js/call-bind-9ceb8f5b.js","assets/js/object-inspect-7b08a46e.js","assets/js/element-plus-7da03efa.js","assets/js/lodash-es-36eb724a.js","assets/js/@vueuse-2d9216d7.js","assets/js/@ctrl-f8748455.js","assets/js/dayjs-6aac7dfa.js","assets/js/async-validator-dee29e8b.js","assets/js/memoize-one-297ddbcb.js","assets/js/escape-html-1d60d822.js","assets/js/normalize-wheel-es-ed76fb12.js","assets/js/@floating-ui-463e90e0.js","assets/js/mitt-f7ef348c.js","assets/js/vue-grid-layout-0ee01be3.js","assets/css/index-c2520709.css","assets/js/requestFSGet-61715c66.js","assets/css/detailDialog-2d838ab0.css"])),y=st(),t=it({logDialogIsShow:!1,logContent:"",tableData:{data:[],total:0,loading:!1,param:{pageNum:1,pageSize:12}},tableLogData:{data:[],total:0,loading:!1,param:{pageNum:1,pageSize:12}},appName:"",logId:0,clusterId:0,clusterData:[]}),I=()=>{t.tableData.loading=!0;const e=[];g.appsList({}).then(function(a){if(a.Status){for(let p=0;p<a.Data.length;p++){var n=a.Data[p];n.FrameworkGitsStr=J(n.FrameworkGits),n.AppGitStr=K(n.AppGit),e.push(n)}setTimeout(()=>{t.tableData.data=e,t.tableData.total=e.length,t.tableData.loading=!1},500)}else t.tableData.data=[],setTimeout(()=>{t.tableData.loading=!1},500)})},m=()=>{t.tableLogData.loading=!0;const e={appName:"",pageIndex:t.tableLogData.param.pageNum,pageSize:t.tableLogData.param.pageSize};g.buildList(e).then(function(a){a.Status?(t.tableLogData.data=a.Data.List,t.tableLogData.total=a.Data.RecordCount,t.tableLogData.loading=!1):(t.tableLogData.data=[],t.tableLogData.loading=!1)})},F=()=>{t.tableData.loading=!0,g.clusterList({}).then(function(e){if(e.Status){var a=[];for(let p=0;p<e.Data.length;p++){var n=e.Data[p];p==0&&(t.clusterId=n.Id),n.Name=n.Name+" - "+n.DockerName,a.push(n)}t.clusterData=a}else t.tableData.data=[]})},G=e=>{t.clusterId=e},U=e=>{y.value.openDialog(e)},H=(e,a)=>{y.value.openDialog(e,a)},R=()=>{B.confirm("此操作将永久清除：“None镜像”，是否继续?","提示",{confirmButtonText:"确认",cancelButtonText:"取消",type:"warning"}).then(()=>{g.dockerClearImage().then(function(e){e.Status?k.success("清除成功"):B.alert(e.StatusMessage,"Warning",{type:"warning",dangerouslyUseHTMLString:!0})})}).catch(()=>{})},P=e=>{t.tableLogData.param.pageSize=e,m()},O=e=>{t.appName=e.AppName,t.tableLogData.param.pageNum=1,t.tableLogData.param.pageSize=10,m()},W=e=>{t.tableLogData.param.pageNum=e,m()};let L=null;rt(()=>t.logDialogIsShow,(e,a)=>{e?L=setInterval(Q,1e3):clearInterval(L)});const j=e=>{t.logId=e.Id,g.buildLog(t.logId.toString()).then(function(a){t.logContent=a,t.logDialogIsShow=!0})},Q=()=>{g.buildLog(t.logId.toString()).then(function(e){console.log(e),t.logContent=e})},q=e=>{var a={AppName:e.AppName,ClusterId:t.clusterId};g.buildAdd(a).then(async function(n){n.Status?(k.success("添加成功"),m()):k.error(n.StatusMessage)})},J=e=>{var a=[];for(let n=0;n<e.length;n++)g.gitInfo({gitId:e[n]}).then(function(p){p.Status&&a.push(p.Data.Name)});return a},K=e=>{var a=[];return g.gitInfo({gitId:e}).then(function(n){n.Status&&a.push(n.Data.Name)}),a};let S=null;return pt(()=>{I(),m(),F(),S=setInterval(m,1e3)}),dt(()=>{clearInterval(S)}),(e,a)=>{const n=r("el-option"),p=r("el-select"),X=r("ele-FolderAdd"),C=r("el-icon"),_=r("el-button"),Y=r("ele-Delete"),f=r("el-tag"),N=r("el-col"),x=r("el-empty"),z=r("el-row"),D=r("el-table-column"),Z=r("el-table"),tt=r("el-pagination"),A=r("el-card"),et=r("el-dialog"),at=ut("loading");return i(),b("div",mt,[s(A,{shadow:"hover",class:"layout-padding-auto"},{default:l(()=>[c("div",ft,[s(p,{modelValue:t.clusterId,"onUpdate:modelValue":a[0]||(a[0]=o=>t.clusterId=o),placeholder:"请选择集群",class:"ml10",onChange:G},{default:l(()=>[(i(!0),b(v,null,M(t.clusterData,o=>(i(),d(n,{key:o.Id,label:o.Name,value:o.Id},null,8,["label","value"]))),128))]),_:1},8,["modelValue"]),s(_,{size:"default",type:"success",class:"ml10",onClick:a[1]||(a[1]=o=>U("add"))},{default:l(()=>[s(C,null,{default:l(()=>[s(X)]),_:1}),u(" 新增应用 ")]),_:1}),s(_,{size:"default",type:"danger",class:"ml10",onClick:a[2]||(a[2]=o=>R("add"))},{default:l(()=>[s(C,null,{default:l(()=>[s(Y)]),_:1}),u(" 清除None镜像 ")]),_:1})]),c("div",_t,[s(z,{style:{float:"left",width:"65%"}},{default:l(()=>[t.tableData.data.length>0?(i(),d(N,{key:0,style:{float:"left",margin:"10px"},xs:24,sm:24,md:24,lg:24,xl:24,class:"mb15"},{default:l(()=>[(i(!0),b(v,null,M(t.tableData.data,(o,V)=>(i(),b("div",{style:{background:"#ffffff",width:"24%"},class:"flex-warp-item",key:V},[c("div",{class:"flex-warp-item-box",onClick:T=>O(o)},[c("div",bt,[s(f,{size:"mini"},{default:l(()=>[u(h(o.AppName),1)]),_:2},1024),o.IsHealth?(i(),d(f,{key:0,size:"mini",type:"success"},{default:l(()=>[u("健康")]),_:1})):(i(),d(f,{key:1,size:"mini",type:"warning"},{default:l(()=>[u("不健康")]),_:1}))]),c("div",ht,"容器版本："+h(o.DockerVer),1),c("div",yt,"集群版本："+h(o.ClusterVer),1),c("div",kt,"仓库："+h(o.AppGitName),1),c("div",vt,"容器文件路径："+h(o.DockerfilePath),1),c("div",wt,[s(_,{size:"default",onClick:T=>H("edit",o),type:"warning"},{default:l(()=>[u("修改")]),_:2},1032,["onClick"]),s(_,{onClick:T=>q(o),size:"default",type:"danger"},{default:l(()=>[u("构建")]),_:2},1032,["onClick"])])],8,Dt)]))),128))]),_:1})):(i(),d(x,{key:1,description:"暂无数据"}))]),_:1}),s(z,{style:{width:"35%",float:"left"}},{default:l(()=>[t.tableLogData.data.length>0?(i(),d(N,{key:0,style:{float:"left",background:"#ffffff",margin:"15px 10px 10px 0",padding:"5px",width:"98%"},xs:24,sm:24,md:24,lg:24,xl:24,class:"mb15"},{default:l(()=>[It,t.tableLogData.data.length>0?(i(),b(v,{key:0},[ct((i(),d(Z,{data:t.tableLogData.data,style:{width:"100%",background:"#ffffff"}},{default:l(()=>[s(D,{prop:"Id",label:"编号",width:"70"}),s(D,{prop:"AppName",label:"应用名称"}),s(D,{label:"状态",width:"90","show-overflow-tooltip":""},{default:l(o=>[o.row.Status==0?(i(),d(f,{key:0,size:"mini",type:"info"},{default:l(()=>[u("未开始")]),_:1})):o.row.Status==1?(i(),d(f,{key:1,size:"mini",type:"success"},{default:l(()=>[u("构建中")]),_:1})):o.row.Status==2?(i(),d(f,{key:2,size:"mini"},{default:l(()=>[u("完成")]),_:1})):w("",!0)]),_:1}),s(D,{prop:"FinishAt",width:"170",label:"完成时间"}),s(D,{label:"操作",width:"80"},{default:l(o=>[o.row.Status!=0?(i(),d(_,{key:0,size:"small",type:"success",onClick:V=>j(o.row)},{default:l(()=>[u("日志")]),_:2},1032,["onClick"])):w("",!0)]),_:1})]),_:1},8,["data"])),[[at,t.tableLogData.loading]]),s(tt,{onSizeChange:P,onCurrentChange:W,class:"mt15","pager-count":5,"page-sizes":[10,20,30],"current-page":t.tableLogData.param.pageNum,"onUpdate:currentPage":a[3]||(a[3]=o=>t.tableLogData.param.pageNum=o),background:"","page-size":t.tableLogData.param.pageSize,"onUpdate:pageSize":a[4]||(a[4]=o=>t.tableLogData.param.pageSize=o),layout:"total, sizes, prev, pager, next, jumper",total:t.tableLogData.total},null,8,["current-page","page-size","total"])],64)):w("",!0)]),_:1})):(i(),d(x,{key:1,description:"暂无数据"}))]),_:1})])]),_:1}),s(gt(E),{ref_key:"appDialogRef",ref:y,onRefresh:a[5]||(a[5]=o=>I())},null,512),s(et,{title:"构建日志",modelValue:t.logDialogIsShow,"onUpdate:modelValue":a[6]||(a[6]=o=>t.logDialogIsShow=o),style:{width:"80%",height:"85%",top:"20px","margin-bottom":"50px"}},{default:l(()=>[s(A,{shadow:"hover",class:"layout-padding-auto",style:{"background-color":"#393d49",overflow:"auto"}},{default:l(()=>[c("pre",{style:{color:"#fff","background-color":"#393d49",height:"100%"},innerHTML:t.logContent},null,8,Lt)]),_:1})]),_:1},8,["modelValue"])])}}});export{pe as default};