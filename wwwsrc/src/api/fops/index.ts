import request from '/@/utils/request';
import requestGet from '/@/utils/requestGet';
import requestFS from '/@/utils/requestFS';
import requestFSGet from '/@/utils/requestFSGet';

/**
 * （不建议写成 request.post(xxx)，因为这样 post 时，无法 params 与 data 同时传参）
 *
 * 登录api接口集合
 * @method gitList git列表
 * @method gitAdd git添加
 * @method gitEdit git修改
 * @method gitDel git删除
 */
export function fopsApi() {
	return {
		gitList: (param: object) => {
			return request({
				url: '/git/list',
				method: 'post',
				data:param,
			});
		},
		gitAdd: (param: object) => {
			return request({
				url: '/git/add',
				method: 'post',
				data:param,
			});
		},
		gitInfo: (param: object) => {
			return request({
				url: '/git/info',
				method: 'post',
				data:param,
			});
		},
		gitEdit: (param: object) => {
			return request({
				url: '/git/update',
				method: 'post',
				data:param,
			});
		},
		gitDel: (param: object) => {
			return request({
				url: '/git/delete',
				method: 'post',
				data:param,
			});
		},appsList: (param: object) => {
			return request({
				url: '/apps/list',
				method: 'post',
				data:param,
			});
		},
		appsAdd: (param: object) => {
			return request({
				url: '/apps/add',
				method: 'post',
				data:param,
			});
		},
		appsEdit: (param: object) => {
			return request({
				url: '/apps/update',
				method: 'post',
				data:param,
			});
		},
		appsDel: (param: object) => {
			return request({
				url: '/apps/delete',
				method: 'post',
				data:param,
			});
		},appsDetail: (param: object) => {
			return request({
				url: '/apps/info',
				method: 'post',
				data:param,
			});
		},clusterAdd: (param: object) => {
			return request({
				url: '/cluster/add',
				method: 'post',
				data:param,
			});
		},clusterEdit: (param: object) => {
			return request({
				url: '/cluster/update',
				method: 'post',
				data:param,
			});
		},clusterList: (param: object) => {
			return request({
				url: '/cluster/list',
				method: 'post',
				data:param,
			});
		},clusterDel: (param: object) => {
			return request({
				url: '/cluster/delete',
				method: 'post',
				data:param,
			});
		},buildList: (param: object) => {
			return request({
				url: '/apps/build/list',
				method: 'post',
				data:param,
			});
		},buildAdd: (param: object) => {
			return request({
				url: '/apps/build/add',
				method: 'post',
				data:param,
			});
		},buildStop: (param: object) => {
			return request({
				url: '/apps/build/stop',
				method: 'post',
				data:param,
			});
		},restartDocker: (param: object) => {
			return request({
				url: '/apps/build/restartDocker',
				method: 'post',
				data:param,
			});
		},syncDockerImage: (param: object) => {
			return request({
				url: '/apps/build/syncDockerImage',
				method: 'post',
				data:param,
			});
		},buildLog: (param: string) => {
			return requestGet({
				url: '/apps/build/view-'+param,
				method: 'get',
				data:{},
			});
		},dockerClearImage: () => {
			return request({
				url: '/apps/build/clearDockerImage',
				method: 'post'
			});
		},taskGroupList: (param: string) => {
			return requestFSGet({
				url: '/basicapi/taskGroup/list?'+param,
				method: 'get',
			});
		},taskList: (param: string) => {
			return requestFSGet({
				url: '/basicapi/task/list?'+param,
				method: 'get',
			});
		},taskGroupInfo: (url:string) => {
			return requestFSGet({
				url: url,
				method: 'get',
			});
		},taskGroupSetEnable: (param: object) => {
			return requestFS({
				url: '/basicapi/taskGroup/setEnable',
				method: 'post',
				data:param,
			});
		},taskUpdate: (param: object) => {
			return requestFS({
				url: '/basicapi/taskGroup/update',
				method: 'post',
				data:param,
			});
		},taskDel: (param: object) => {
			return requestFS({
				url: '/basicapi/taskGroup/delete',
				method: 'post',
				data:param,
			});
		},taskCount: (param: object) => {
			return requestFS({
				url: '/basicapi/taskGroup/count',
				method: 'get',
				data:param,
			});
		},taskNoRunCount: (param: object) => {
			return requestFS({
				url: '/basicapi/taskGroup/unRunCount',
				method: 'get',
				data:param,
			});
		},taskNoRunList: (param: string) => {
			return requestFSGet({
				url: '/basicapi/taskGroup/unRunList?'+param,
				method: 'get',
			});
		},taskRunningList: (param: string) => {
			return requestFSGet({
				url: '/basicapi/taskGroup/schedulerWorkingList?'+param,
				method: 'get',
			});
		},taskLogList: (param: string) => {
			return requestFSGet({
				url: '/basicapi/log/list?'+param,
				method: 'get',
			});
		},taskLogListClientName: (param: string) => {
			return requestFSGet({
				url: '/basicapi/log/listByClientName?'+param,
				method: 'get',
			});
		},taskStatList: (param: string) => {
			return requestFSGet({
				url: '/basicapi/task/statList?'+param,
				method: 'get',
			});
		},taskFinishedList: (param: object) => {
			return requestFSGet({
				url: '/basicapi/task/list',
				method: 'get',
				data:param,
			});
		},taskTodayFailCount: (param: object) => {
			return requestFS({
				url: '/basicapi/task/todayFailCount',
				method: 'get',
				data:param,
			});
		},serverNodeList: (param: object) => {
			return requestFS({
				url: '/basicapi/server/list',
				method: 'get',
				data:param,
			});
		},clientList: (param: object) => {
			return requestFS({
				url: '/basicapi/client/list',
				method: 'get',
				data:param,
			});
		},linkTraceWebApi:(param: string) => {
			return requestGet({
				url: '/linkTrace/webApiList?'+param,
				method: 'get',
				//data:param,
			});
		},linkTraceTask:(param: string) => {
			return requestGet({
				url: '/linkTrace/taskList?'+param,
				method: 'get',
				//data:param,
			});
		},linkTraceConsumerList:(param: string) => {
			return requestGet({
				url: '/linkTrace/consumerList?'+param,
				method: 'get',
				//data:param,
			});
		},linkTraceFScheduleList:(param: string) => {
			return requestGet({
				url: '/linkTrace/fScheduleList?'+param,
				method: 'get',
				//data:param,
			});
		},linkTraceInfo:(traceId: object) => {
			return requestGet({
				url: '/linkTrace/info/'+traceId,
				method: 'get',
			});
		},slowSql:(param: string) => {
			return requestGet({
				url: '/linkTrace/slowDbList?'+param,
				method: 'get',
				//data:param,
			});
		},slowEs:(param: string) => {
			return requestGet({
				url: '/linkTrace/slowEsList?'+param,
				method: 'get',
				//data:param,
			});
		},slowEtcd:(param: string) => {
			return requestGet({
				url: '/linkTrace/slowEtcdList?'+param,
				method: 'get',
				//data:param,
			});
		},slowHand:(param: string) => {
			return requestGet({
				url: '/linkTrace/slowHandList?'+param,
				method: 'get',
				//data:param,
			});
		},slowHttp:(param: string) => {
			return requestGet({
				url: '/linkTrace/slowHttpList?'+param,
				method: 'get',
				//data:param,
			});
		},slowMq:(param: string) => {
			return requestGet({
				url: '/linkTrace/slowMqList?'+param,
				method: 'get',
				//data:param,
			});
		},slowRedis:(param: string) => {
			return requestGet({
				url: '/linkTrace/slowRedisList?'+param,
				method: 'get',
				//data:param,
			});
		},logList:(param: string) => {
			return requestGet({
				url: '/flog/list?'+param,
				method: 'get',
				//data:param,
			});
		},logInfo:(param: string) => {
			return requestGet({
				url: '/flog/info-'+param,
				method: 'get',
				//data:param,
			});
		},logStatCount:(param: string) => {
			return requestGet({
				url: '/flog/StatCount?'+param,
				method: 'get',
				//data:param,
			});
		}
	};
}
