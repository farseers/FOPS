<template>
	<div class="layout-navbars-breadcrumb-user pr15" :style="{ flex: layoutUserFlexNum }">
		<el-dropdown :show-timeout="70" :hide-timeout="50" trigger="click" @command="onComponentSizeChange">
			<div class="layout-navbars-breadcrumb-user-icon">
				<i class="iconfont icon-ziti" :title="$t('message.user.title0')"></i>
			</div>
			<template #dropdown>
				<el-dropdown-menu>
					<el-dropdown-item command="large" :disabled="state.disabledSize === 'large'">{{ $t('message.user.dropdownLarge') }}</el-dropdown-item>
					<el-dropdown-item command="default" :disabled="state.disabledSize === 'default'">{{ $t('message.user.dropdownDefault') }}</el-dropdown-item>
					<el-dropdown-item command="small" :disabled="state.disabledSize === 'small'">{{ $t('message.user.dropdownSmall') }}</el-dropdown-item>
				</el-dropdown-menu>
			</template>
		</el-dropdown>
		<el-dropdown :show-timeout="70" :hide-timeout="50" trigger="click" @command="onLanguageChange">
			<div class="layout-navbars-breadcrumb-user-icon">
				<i
					class="iconfont"
					:class="state.disabledI18n === 'en' ? 'icon-fuhao-yingwen' : 'icon-fuhao-zhongwen'"
					:title="$t('message.user.title1')"
				></i>
			</div>
			<template #dropdown>
				<el-dropdown-menu>
					<el-dropdown-item command="zh-cn" :disabled="state.disabledI18n === 'zh-cn'">简体中文</el-dropdown-item>
					<el-dropdown-item command="en" :disabled="state.disabledI18n === 'en'">English</el-dropdown-item>
					<el-dropdown-item command="zh-tw" :disabled="state.disabledI18n === 'zh-tw'">繁體中文</el-dropdown-item>
				</el-dropdown-menu>
			</template>
		</el-dropdown>
		<div class="layout-navbars-breadcrumb-user-icon" @click="onSearchClick">
			<el-icon :title="$t('message.user.title2')">
				<ele-Search />
			</el-icon>
		</div>
		<div class="layout-navbars-breadcrumb-user-icon" @click="onLayoutSetingClick">
			<i class="icon-skin iconfont" :title="$t('message.user.title3')"></i>
		</div>
		<div class="layout-navbars-breadcrumb-user-icon" ref="userNewsBadgeRef" @click="onUserNewsClick">
			<el-badge :is-dot="state.isDot">
				<el-icon :title="$t('message.user.title4')">
					<ele-Bell />
				</el-icon>
			</el-badge>
		</div>
		<el-popover
			ref="userNewsRef"
			:visible="state.manualVisible"
			:virtual-ref="userNewsBadgeRef"
			placement="bottom"
			trigger="click"
			transition="el-zoom-in-top"
			virtual-triggering
			:width="300"
			:persistent="false"
		>
			<UserNews :NoReadList="state.NoReadList" @newsCK="newsCK" @allRead="allRead"/>
		</el-popover>
		<div class="layout-navbars-breadcrumb-user-icon mr10" @click="onScreenfullClick">
			<i
				class="iconfont"
				:title="state.isScreenfull ? $t('message.user.title6') : $t('message.user.title5')"
				:class="!state.isScreenfull ? 'icon-fullscreen' : 'icon-tuichuquanping'"
			></i>
		</div>
		<el-dropdown :show-timeout="70" :hide-timeout="50" @command="onHandleCommandClick">
			<span class="layout-navbars-breadcrumb-user-link">
				<img :src="userInfos.photo" class="layout-navbars-breadcrumb-user-link-photo mr5" />
				{{ userInfos.userName === '' ? 'common' : userInfos.userName }}
				<el-icon class="el-icon--right">
					<ele-ArrowDown />
				</el-icon>
			</span>
			<template #dropdown>
				<el-dropdown-menu>
					<el-dropdown-item command="/market/monitorCenter">{{ $t('message.user.dropdown1') }}</el-dropdown-item>
					<el-dropdown-item command="changePsd">{{ $t('message.user.dropdown7') }}</el-dropdown-item>
					<!-- <el-dropdown-item command="wareHouse">{{ $t('message.user.dropdown6') }}</el-dropdown-item> -->
					<!-- <el-dropdown-item command="/personal">{{ $t('message.user.dropdown2') }}</el-dropdown-item> -->
					<!-- <el-dropdown-item command="/404">{{ $t('message.user.dropdown3') }}</el-dropdown-item> -->
					<!-- <el-dropdown-item command="/401">{{ $t('message.user.dropdown4') }}</el-dropdown-item> -->
					<el-dropdown-item divided command="logOut">{{ $t('message.user.dropdown5') }}</el-dropdown-item>
				</el-dropdown-menu>
			</template>
		</el-dropdown>
		<Search ref="searchRef" />
		<el-dialog v-model="state.psdShow" title="修改密码" width="400px">
		<div>
			<el-form-item label="新密码"><el-input clearable placeholder="请输入密码" maxlength="20"
                        v-model="state.LoginPwd"/></el-form-item>
                <el-form-item label="再次输入"><el-input clearable placeholder="请再次输入密码" maxlength="20"
                        v-model="state.ConfirmPwd"/></el-form-item>
                <el-form-item style="text-align: center;">
					<div style="text-align: center;width: 100%;">
						<el-button type="success" @click="submitForm()">确认</el-button>
					</div>
				</el-form-item>
		</div>
	</el-dialog>
	</div>
</template>

<script setup lang="ts" name="layoutBreadcrumbUser">
import { defineAsyncComponent, ref, unref, computed, reactive, onMounted, onUnmounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessageBox, ElMessage,ElNotification, ClickOutside as vClickOutside } from 'element-plus';
import screenfull from 'screenfull';
import { useI18n } from 'vue-i18n';
import { storeToRefs } from 'pinia';
import { useUserInfo } from '/@/stores/userInfo';
import { useThemeConfig } from '/@/stores/themeConfig';
import other from '/@/utils/other';
import mittBus from '/@/utils/mitt';
import { Session, Local } from '/@/utils/storage';
import WebSocketService from '/@/utils/socket.js';
import { fopsApi } from '/@/api/fops';
const serverApi = fopsApi();
// 引入组件
const UserNews = defineAsyncComponent(() => import('/@/layout/navBars/topBar/userNews.vue'));
const Search = defineAsyncComponent(() => import('/@/layout/navBars/topBar/search.vue'));

// 定义变量内容
const userNewsRef = ref();
const userNewsBadgeRef = ref();
const { locale, t } = useI18n();
const router = useRouter();
const stores = useUserInfo();
const storesThemeConfig = useThemeConfig();
const { userInfos } = storeToRefs(stores);
const { themeConfig } = storeToRefs(storesThemeConfig);
const searchRef = ref();
const wsInfo =  ref();
const state = reactive({
	isScreenfull: false,
	disabledI18n: 'zh-cn',
	disabledSize: 'large',
	NoReadList:[],
	isDot:false,
	manualVisible:false,//未读消息
	psdShow:false,//修改密码
	LoginPwd:'',
	ConfirmPwd:'',
});

// 设置分割样式
const layoutUserFlexNum = computed(() => {
	let num: string | number = '';
	const { layout, isClassicSplitMenu } = themeConfig.value;
	const layoutArr: string[] = ['defaults', 'columns'];
	if (layoutArr.includes(layout) || (layout === 'classic' && !isClassicSplitMenu)) num = '1';
	else num = '';
	return num;
});
// 全屏点击时
const onScreenfullClick = () => {
	if (!screenfull.isEnabled) {
		ElMessage.warning('暂不不支持全屏');
		return false;
	}
	screenfull.toggle();
	screenfull.on('change', () => {
		if (screenfull.isFullscreen) state.isScreenfull = true;
		else state.isScreenfull = false;
	});
};
// 消息通知点击时
const onUserNewsClick = () => {
	state.manualVisible = !state.manualVisible
	// unref(userNewsRef).popperRef?.delayHide?.();
};
// 布局配置 icon 点击时
const onLayoutSetingClick = () => {
	mittBus.emit('openSetingsDrawer');
};
// 下拉菜单点击时
const onHandleCommandClick = (path: string) => {
	if (path === 'logOut') {
		ElMessageBox({
			closeOnClickModal: false,
			closeOnPressEscape: false,
			title: t('message.user.logOutTitle'),
			message: t('message.user.logOutMessage'),
			showCancelButton: true,
			confirmButtonText: t('message.user.logOutConfirm'),
			cancelButtonText: t('message.user.logOutCancel'),
			buttonSize: 'default',
			beforeClose: (action, instance, done) => {
				if (action === 'confirm') {
					instance.confirmButtonLoading = true;
					instance.confirmButtonText = t('message.user.logOutExit');
					setTimeout(() => {
						done();
						setTimeout(() => {
							instance.confirmButtonLoading = false;
						}, 300);
					}, 700);
				} else {
					done();
				}
			},
		})
			.then(async () => {
				// 清除缓存/token等
				Session.clear();
				// 使用 reload 时，不需要调用 resetRoute() 重置路由
				window.location.reload();
			})
			.catch(() => {});
	} else if (path === 'wareHouse') {
		window.open('https://gitee.com/lyt-top/vue-next-admin');
	} else if(path == 'changePsd') { //修改密码
		state.psdShow = true
	}else  {
		router.push(path);
	}
};
// 菜单搜索点击
const onSearchClick = () => {
	searchRef.value.openSearch();
};
// 组件大小改变
const onComponentSizeChange = (size: string) => {
	Local.remove('themeConfig');
	themeConfig.value.globalComponentSize = size;
	Local.set('themeConfig', themeConfig.value);
	initI18nOrSize('globalComponentSize', 'disabledSize');
	window.location.reload();
};
// 语言切换
const onLanguageChange = (lang: string) => {
	Local.remove('themeConfig');
	themeConfig.value.globalI18n = lang;
	Local.set('themeConfig', themeConfig.value);
	locale.value = lang;
	other.useTitle();
	initI18nOrSize('globalI18n', 'disabledI18n');
};
// 初始化组件大小/i18n
const initI18nOrSize = (value: string, attr: string) => {
	(<any>state)[attr] = Local.get('themeConfig')[value];
};
// 页面加载时
onMounted(() => {
	if (Local.get('themeConfig')) {
		initI18nOrSize('globalComponentSize', 'disabledSize');
		initI18nOrSize('globalI18n', 'disabledI18n');
	}
	infoConnect()
});
onUnmounted(() => {
	if (wsInfo.value) {
        wsInfo.value.close();
      }
})
const submitForm = ()=>{
	if(state.LoginPwd && state.ConfirmPwd ){
		if(state.LoginPwd != state.ConfirmPwd){
			ElMessage.warning('两次输入密码不相同');
		}else{
			serverApi.changePsd({
				ConfirmPwd:state.ConfirmPwd,
				LoginPwd:state.LoginPwd
			}).then((d)=>{
                if(d.Status){
                    ElMessage.success('修改成功');
					state.psdShow = false;
					Session.clear();
					window.location.reload();
                }else{
                    ElMessage.error(d.StatusMessage);
                }
            }).catch(e=>{
                ElMessage.error('网络错误');
            })
		}
	}else{
		ElMessage.warning('请输入完整');
	}
}
const allRead = ()=>{ //全部已读
	ElNotification.closeAll()
	state.manualVisible = false;
	wsInfo.value && wsInfo.value.sendMessage(JSON.stringify('1'))
}
const elNotClick =()=>{ //点击消息
	ElNotification.closeAll()
	state.manualVisible = true;
}
const newsCK = ()=>{ //跳到通知中心
	ElNotification.closeAll()
	state.manualVisible = false;
	router.push('/monitor/malog');
}
const infoConnect =() =>{ //
	wsInfo.value = new WebSocketService(`ws/notice`, JSON.stringify('1')); 
	wsInfo.value.onMessage((data:any) => {
        if (data) {
          state.NoReadList = data.NoReadList ? data.NoReadList : [];
        }else{
		  state.NoReadList = []
		}
		if(state.NoReadList.length > 0){
			state.isDot = true;
			ElNotification.closeAll()
			// if(!state.manualVisible){
			// 	ElNotification({
			// 	title: '提示',
			// 	message: '有新的日志信息',
			// 	type: 'warning',
			// 	duration:0,
			// 	position: 'bottom-right',
			// 	onClick:elNotClick
			// });
			// }
			
		}else{
			state.isDot = false;
		}
      });
    }
</script>

<style scoped lang="scss">
.layout-navbars-breadcrumb-user {
	display: flex;
	align-items: center;
	justify-content: flex-end;
	&-link {
		height: 100%;
		display: flex;
		align-items: center;
		white-space: nowrap;
		&-photo {
			width: 25px;
			height: 25px;
			border-radius: 100%;
		}
	}
	&-icon {
		padding: 0 10px;
		cursor: pointer;
		color: var(--next-bg-topBarColor);
		height: 50px;
		line-height: 50px;
		display: flex;
		align-items: center;
		&:hover {
			background: var(--next-color-user-hover);
			i {
				display: inline-block;
				animation: logoAnimation 0.3s ease-in-out;
			}
		}
	}
	:deep(.el-dropdown) {
		color: var(--next-bg-topBarColor);
	}
	:deep(.el-badge) {
		height: 40px;
		line-height: 40px;
		display: flex;
		align-items: center;
	}
	:deep(.el-badge__content.is-fixed) {
		top: 12px;
	}
}
</style>
