<template>
	<div class="layout-navbars-breadcrumb-user-news" style="max-height: 500px;display: flex;flex-flow: column;">
		<div class="head-box">
			<div class="head-box-title">{{ $t('message.user.newTitle') }}</div>
			<div class="head-box-btn" v-if="NoReadList.length > 0" @click="allRead">{{ $t('message.user.newBtn') }}</div>
		</div>
		<div class="content-box" style="flex: 1;overflow: auto;">
			<template v-if="NoReadList.length > 0">
				<div class="content-box-item" v-for="(v, k) in NoReadList" :key="k">
					<div>{{ v.AppName }}</div>
					<div class="content-box-msg">
						{{ v.NoticeMsg }}
					</div>
					<div class="content-box-time" v-html="formattedTime(v.NoticeAt)"></div>
				</div>
			</template>
			<el-empty :description="$t('message.user.newDesc')" v-else></el-empty>
		</div>
		<div class="foot-box" v-if="NoReadList.length > 0" @click="goTOPage">{{ $t('message.user.newGo') }}</div>
	</div>
</template>

<script>
import { fopsApi } from "/@/api/fops";
const serverApi = fopsApi();
export default {
	name:'layoutBreadcrumbUserNews',
	props:{
		NoReadList:{
			type:Array,
			default:()=>{
				return []
			}
		}
	},
	data(){
	return {

	}
},
methods:{
	formattedTime(d){
		const date = new Date(d);
		return date.toLocaleString(); // 使用本地时间格式
	},
	allRead(){
		const Ids = this.NoReadList.map(obj => obj.Id);
		serverApi.monitorAllRead({
			"Ids":Ids
		}).then((d)=>{
                if(d.Status){
                    this.$emit('allRead')
                }else{
                    ElMessage.error(d.StatusMessage);
                }
            }).catch(e=>{
                ElMessage.error('网络错误');
            })
	},
	goTOPage(){
		this.$emit('newsCK')
	},
	
}
}


</script>

<style scoped lang="scss">
.layout-navbars-breadcrumb-user-news {
	.head-box {
		display: flex;
		border-bottom: 1px solid var(--el-border-color-lighter);
		box-sizing: border-box;
		color: var(--el-text-color-primary);
		justify-content: space-between;
		height: 35px;
		align-items: center;
		.head-box-btn {
			color: var(--el-color-primary);
			font-size: 13px;
			cursor: pointer;
			opacity: 0.8;
			&:hover {
				opacity: 1;
			}
		}
	}
	.content-box {
		font-size: 13px;
		.content-box-item {
			padding-top: 12px;
			&:last-of-type {
				padding-bottom: 12px;
			}
			.content-box-msg {
				color: var(--el-text-color-secondary);
				margin-top: 5px;
				margin-bottom: 5px;
			}
			.content-box-time {
				color: var(--el-text-color-secondary);
			}
		}
	}
	.foot-box {
		height: 35px;
		color: var(--el-color-primary);
		font-size: 13px;
		cursor: pointer;
		opacity: 0.8;
		display: flex;
		align-items: center;
		justify-content: center;
		border-top: 1px solid var(--el-border-color-lighter);
		&:hover {
			opacity: 1;
		}
	}
	:deep(.el-empty__description p) {
		font-size: 13px;
	}
}
</style>
