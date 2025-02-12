import { RouteRecordRaw } from 'vue-router';

/**
 * 建议：路由 path 路径与文件夹名称相同，找文件可浏览器地址找，方便定位文件位置
 *
 * 路由meta对象参数说明
 * meta: {
 *      title:          菜单栏及 tagsView 栏、菜单搜索名称（国际化）
 *      isLink：        是否超链接菜单，开启外链条件，`1、isLink: 链接地址不为空 2、isIframe:false`
 *      isHide：        是否隐藏此路由
 *      isKeepAlive：   是否缓存组件状态
 *      isAffix：       是否固定在 tagsView 栏上
 *      isIframe：      是否内嵌窗口，开启条件，`1、isIframe:true 2、isLink：链接地址不为空`
 *      roles：         当前路由权限标识，取角色管理。控制路由显示、隐藏。超级管理员：admin 普通角色：common
 *      icon：          菜单、tagsView 图标，阿里：加 `iconfont xxx`，fontawesome：加 `fa xxx`
 * }
 */

// 扩展 RouteMeta 接口
declare module 'vue-router' {
	interface RouteMeta {
		title?: string;
		isLink?: string;
		isHide?: boolean;
		isKeepAlive?: boolean;
		isAffix?: boolean;
		isIframe?: boolean;
		roles?: string[];
		icon?: string;
	}
}

/**
 * 定义动态路由
 * 前端添加路由，请在顶级节点的 `children 数组` 里添加
 * @description 未开启 isRequestRoutes 为 true 时使用（前端控制路由），开启时第一个顶级 children 的路由将被替换成接口请求回来的路由数据
 * @description 各字段请查看 `/@/views/system/menu/component/addMenu.vue 下的 ruleForm`
 * @returns 返回路由菜单数据
 */
export const dynamicRoutes: Array<RouteRecordRaw> = [
	{
		path: '/',
		name: '/',
		component: () => import('/@/layout/index.vue'),
		redirect: '/market',
		meta: {
			isKeepAlive: true,
		},
		children: [
			{
				path: '/market',
				name: 'market',
				component: () => import('/@/views/fops/task/indexColonyNode.vue'),
				meta: {
					title: 'message.router.Market',
					isLink: '',
					isHide: false,
					isKeepAlive: true,
					isAffix: true,
					isIframe: false,
					roles: ['admin', 'common'],
					icon: 'iconfont icon-putong',
				},children: [
					// {
					// 	path: '/market/monitorCenter',
					// 	name: 'marketMonitorCenter',
					// 	component: () => import('/@/views/fops/market/index.vue'),
					// 	meta: {
					// 		title: 'message.router.MonitoringCenter',
					// 		isLink: '',
					// 		isHide: false,
					// 		isKeepAlive: true,
					// 		isAffix: false,
					// 		isIframe: false,
					// 		roles: ['admin', 'common'],
					// 		icon: 'iconfont icon-ico_shuju',
					// 	},
					// },
					{
						path: '/market/colonyNode',
						name: 'marketColonyNode',
						component: () => import('/@/views/fops/task/indexColonyNode.vue'),
						meta: {
							title: 'message.router.ColonyNode',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: true,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-ico_shuju',
						},
					},
					// {
					// 	path: '/market/serverNode',
					// 	name: 'marketServerNode',
					// 	component: () => import('/@/views/fops/task/indexServerNode.vue'),
					// 	meta: {
					// 		title: 'message.router.ServerNode',
					// 		isLink: '',
					// 		isHide: false,
					// 		isKeepAlive: true,
					// 		isAffix: false,
					// 		isIframe: false,
					// 		roles: ['admin', 'common'],
					// 		icon: 'iconfont icon-ico_shuju',
					// 	},
					// },
					{
						path: '/terminal/index',
						name: 'terminalIndex',
						component: () => import('/@/views/fops/terminal/index.vue'),
						meta: {
							title: 'message.router.terminal',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					}
				]
			},{
				path: '/app',
				name: 'app',
				component: () => import('/@/views/fops/build/indexApp.vue'),
				meta: {
					title: 'message.router.AppCenter',
					isLink: '',
					isHide: false,
					isKeepAlive: true,
					isAffix: true,
					isIframe: false,
					roles: ['admin', 'common'],
					icon: 'iconfont icon-shouye',
				},children: [
					{
						path: '/app',
						name: 'appBuildDeployment',
						component: () => import('/@/views/fops/build/indexApp.vue'),
						meta: {
							title: 'message.router.BuildDeployment',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-step',
						},
					},{
						path: '/autoBuild',
						name: 'autoBuildDeployment',
						component: () => import('/@/views/fops/autoBuild/indexApp.vue'),
						meta: {
							title: 'message.router.AutoBuildDeployment',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-step',
						},
					},{
						path: '/log/list',
						name: 'logList',
						component: () => import('/@/views/fops/log/indexLogV2.vue'),
						meta: {
							title: 'message.router.LogList',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: true,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/cluster/clusterAdd',
						name: 'clusterAdd',
						component: () => import('/@/views/fops/cluster/index.vue'),
						meta: {
							title: 'message.router.ClusterAdd',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-wenducanshu-05',
						},
					},{
						path: '/app/gitAdd',
						name: 'appGitAdd',
						component: () => import('/@/views/fops/git/index.vue'),
						meta: {
							title: 'message.router.GitAdd',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zujian',
						},
					},
					{
						path: '/app/Configuration',
						name: 'Configuration',
						component: () => import('/@/views/fops/configu/index.vue'),
						meta: {
							title: 'message.router.Configuration',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin'],
							icon: 'iconfont icon-caidan',
						},
					},
					//
				]
			},{
				path: '/dispatch',
				name: 'dispatch',
				component: () => import('/@/views/home/index.vue'),
				meta: {
					title: 'message.router.DispatchCenter',
					isLink: '',
					isHide: false,
					isKeepAlive: true,
					isAffix: false,
					isIframe: false,
					roles: ['admin', 'common'],
					icon: 'iconfont icon-crew_feature',
				},children: [
					{
						path: '/dispatch/taskGroup',
						name: 'taskGroup',
						component: () => import('/@/views/fops/task/indexGroupList.vue'),
						meta: {
							title: 'message.router.TaskGroup',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: true,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-shuxingtu',
						},
					},{
						path: '/dispatch/taskGroup2',
						name: 'taskGroup2',
						component: () => import('/@/views/fops/task/taskHistory.vue'),
						meta: {
							title: 'message.router.TaskHistory',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-shuxingtu',
						},
					},{
						path: '/dispatch/taskPlan',
						name: 'taskPlan',
						component: () => import('/@/views/fops/task/taskPlan.vue'),
						meta: {
							title: 'message.router.TaskPlan',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: true,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-jinridaiban',
						},
					},{
						path: '/dispatch/taskLogList',
						name: 'taskLogList',
						component: () => import('/@/views/fops/task/indexLog.vue'),
						meta: {
							title: 'message.router.taskLogList',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-jinridaiban',
						},
					},{
						path: '/dispatch/node',
						name: 'taskNode',
						component: () => import('/@/views/fops/task/indexServerNode.vue'),
						meta: {
							title: 'message.router.DispatchNode',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-webicon318',
						},
					},{
						path: '/dispatch/client',
						name: 'client',
						component: () => import('/@/views/fops/task/indexClient.vue'),
						meta: {
							title: 'message.router.Client',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-caidan',
						},
					},
				]
			},{
				path: '/linkTrace',
				name: 'linkTrace',
				component: () => import('/@/views/home/index.vue'),
				meta: {
					title: 'message.router.linkTrace',
					isLink: '',
					isHide: false,
					isKeepAlive: true,
					isAffix: false,
					isIframe: false,
					roles: ['admin', 'common'],
					icon: 'iconfont icon-shuju',
				},children: [
					{
						path: '/linkTrace/visits',
						name: 'visits',
						component: () => import('/@/views/fops/linkTrace/indexVisits.vue'),
						meta: {
							title: 'message.router.visits',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: true,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/linkTrace/webApi',
						name: 'webApi',
						component: () => import('/@/views/fops/linkTrace/indexWebApi.vue'),
						meta: {
							title: 'message.router.webApi',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/linkTrace/webSocket',
						name: 'webSocket',
						component: () => import('/@/views/fops/linkTrace/indexWebSocket.vue'),
						meta: {
							title: 'message.router.webSocket',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/linkTrace/taskApi',
						name: 'taskApi',
						component: () => import('/@/views/fops/linkTrace/indexTaskApi.vue'),
						meta: {
							title: 'message.router.taskApi',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/linkTrace/customerApi',
						name: 'customerApi',
						component: () => import('/@/views/fops/linkTrace/indexCustomerApi.vue'),
						meta: {
							title: 'message.router.customerApi',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/linkTrace/queueApi',
						name: 'queueApi',
						component: () => import('/@/views/fops/linkTrace/indexQueueApi.vue'),
						meta: {
							title: 'message.router.queueApi',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/linkTrace/eventApi',
						name: 'eventApi',
						component: () => import('/@/views/fops/linkTrace/indexEventApi.vue'),
						meta: {
							title: 'message.router.eventApi',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/linkTrace/fscheduleApi',
						name: 'fscheduleApi',
						component: () => import('/@/views/fops/linkTrace/indexFsScheduleApi.vue'),
						meta: {
							title: 'message.router.fscheduleApi',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},]
			},{
				path: '/slowQuery',
				name: 'slowQuery',
				component: () => import('/@/views/home/index.vue'),
				meta: {
					title: 'message.router.SlowQuery',
					isLink: '',
					isHide: false,
					isKeepAlive: true,
					isAffix: false,
					isIframe: false,
					roles: ['admin', 'common'],
					icon: 'iconfont icon-shuju',
				},children: [
					{
						path: '/slowQuery/dataBase',
						name: 'dataBase',
						component: () => import('/@/views/fops/sql/indexSql.vue'),
						meta: {
							title: 'message.router.DataBase',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: true,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/slowQuery/redis',
						name: 'redis',
						component: () => import('/@/views/fops/sql/indexRedis.vue'),
						meta: {
							title: 'message.router.Redis',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/slowQuery/mq',
						name: 'mq',
						component: () => import('/@/views/fops/sql/indexMq.vue'),
						meta: {
							title: 'message.router.MQ',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/slowQuery/http',
						name: 'http',
						component: () => import('/@/views/fops/sql/indexHttp.vue'),
						meta: {
							title: 'message.router.Http',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/slowQuery/etcd',
						name: 'etcd',
						component: () => import('/@/views/fops/sql/indexEtcd.vue'),
						meta: {
							title: 'message.router.Etcd',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/slowQuery/es',
						name: 'es',
						component: () => import('/@/views/fops/sql/indexEs.vue'),
						meta: {
							title: 'message.router.ES',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					},{
						path: '/slowQuery/hand',
						name: 'hand',
						component: () => import('/@/views/fops/sql/indexHand.vue'),
						meta: {
							title: 'message.router.Hand',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
						},
					}
				]
			},
			{
				path:'/monitor',
				name:'monitor',
				component: () => import('/@/views/home/index.vue'),
				meta: {
					title: 'message.router.Monitor',
					isLink: '',
					isHide: false,
					isKeepAlive: true,
					isAffix: false,
					isIframe: false,
					roles: ['admin', 'common'],
					icon: 'iconfont icon-shuju',
				},
				children:[
					   {
						path: '/monitor/magule',
						name: 'magule',
						component: () => import('/@/views/fops/monitor/gule.vue'),
						meta: {
							title: 'message.router.magule',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
					   }
					},
					{
						path: '/monitor/manotice',
						name: 'manotice',
						component: () => import('/@/views/fops/monitor/notice.vue'),
						meta: {
							title: 'message.router.manotice',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
					   }
					},
					{
						path: '/monitor/madata',
						name: 'madata',
						component: () => import('/@/views/fops/monitor/data.vue'),
						meta: {
							title: 'message.router.madata',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
					   }
					},
					{
						path: '/monitor/malog',
						name: 'malog',
						component: () => import('/@/views/fops/monitor/log.vue'),
						meta: {
							title: 'message.router.malog',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
					   }
					}
				]
			},
			{
				path:'/Replication',
				name:'Replication',
				component: () => import('/@/views/home/index.vue'),
				meta: {
					title: 'message.router.Replication',
					isLink: '',
					isHide: false,
					isKeepAlive: true,
					isAffix: false,
					isIframe: false,
					roles: ['admin', 'common'],
					icon: 'iconfont icon-shuju',
				},
				children:[
					   {
						path: '/replication/schedule',
						name: 'schedule',
						component: () => import('/@/views/fops/replication/schedule.vue'),
						meta: {
							title: 'message.router.schedule',
							isLink: '',
							isHide: false,
							isKeepAlive: true,
							isAffix: false,
							isIframe: false,
							roles: ['admin', 'common'],
							icon: 'iconfont icon-zidingyibuju',
					   }
					},
					
				]
			},
		],
	},
];

/**
 * 定义404、401界面
 * @link 参考：https://next.router.vuejs.org/zh/guide/essentials/history-mode.html#netlify
 */
export const notFoundAndNoPower = [
	{
		path: '/:path(.*)*',
		name: 'notFound',
		component: () => import('/@/views/error/404.vue'),
		meta: {
			title: 'message.staticRoutes.notFound',
			isHide: true,
		},
	},
	{
		path: '/401',
		name: 'noPower',
		component: () => import('/@/views/error/401.vue'),
		meta: {
			title: 'message.staticRoutes.noPower',
			isHide: true,
		},
	},
];

/**
 * 定义静态路由（默认路由）
 * 此路由不要动，前端添加路由的话，请在 `dynamicRoutes 数组` 中添加
 * @description 前端控制直接改 dynamicRoutes 中的路由，后端控制不需要修改，请求接口路由数据时，会覆盖 dynamicRoutes 第一个顶级 children 的内容（全屏，不包含 layout 中的路由出口）
 * @returns 返回路由菜单数据
 */
export const staticRoutes: Array<RouteRecordRaw> = [
	{
		path: '/login',
		name: 'login',
		component: () => import('/@/views/login/index.vue'),
		meta: {
			title: '登录',
		},
	},
	/**
	 * 提示：写在这里的为全屏界面，不建议写在这里
	 * 请写在 `dynamicRoutes` 路由数组中
	 */
	{
		path: '/visualizingDemo1',
		name: 'visualizingDemo1',
		component: () => import('/@/views/visualizing/demo1.vue'),
		meta: {
			title: 'message.router.visualizingLinkDemo1',
		},
	},
	{
		path: '/visualizingDemo2',
		name: 'visualizingDemo2',
		component: () => import('/@/views/visualizing/demo2.vue'),
		meta: {
			title: 'message.router.visualizingLinkDemo2',
		},
	},
];
