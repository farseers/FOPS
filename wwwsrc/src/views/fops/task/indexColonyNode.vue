<template>
	<div class="system-user-container layout-padding mtor_main">
		<div class="content">
			<div class="header">
				<AsyncCount ref="childCount" />
			</div>
			<div class="main">
				<div class="aside">
					<AsyncEcfy ref="childEcfy" />
					<AsyncConly ref="childConly" />
				</div>
				<div class="argid">
					<AsyncQueTab ref="childQueTab" />
				</div>
			</div>
		</div>
	</div>
</template>

<script setup name="fopsTaskTimeOut">
import { ref, defineAsyncComponent, reactive, onMounted, onUnmounted } from 'vue';
// 引入 api 请求接口
const AsyncConly = defineAsyncComponent(() => import('./components/Conly.vue'))
const AsyncCount = defineAsyncComponent(() => import('./components/Counts.vue'))
const AsyncQueTab = defineAsyncComponent(() => import('./components/QueTab.vue'))
const AsyncEcfy = defineAsyncComponent(() => import('./components/Ecfy.vue'))
// 定义变量内容
const state = reactive({
	timer: null, //
	
});

const childConly = ref(null);
const childCount = ref(null);
const childQueTab = ref(null);
const childEcfy = ref(null);
const init = () => {
	let time = 0;
	let time1 = 0;
	let time2 = 0;
	state.timer = setInterval(() => {
		time++
		if (childQueTab.value) {
			childQueTab.value.getData()
		}
		if (time1 >= 3) { //调用3秒一次的
			time1 = 0;
			if (childConly.value) {
				childConly.value.getData()
			}
			if (childEcfy.value) {
				childEcfy.value.getData()
			}
		}
		if (time2 >= 10) {
			if (childCount.value) {
				childCount.value.getData()
			}
			time2 = 0
		}

	}, 1000)
}
onUnmounted(() => {
	clearInterval(state.timer)
})
// 页面加载时
onMounted(() => {
	// init();
});
</script>

<style lang="scss">
.mtor_main {
	.content {
		padding: 0;
		flex: 1;
		display: flex;
		flex-flow: column;
		overflow: hidden;
	}

	.header {
		height: auto;
		padding: 5px
	}

	.main {
		padding: 0;
		flex: 1;
		display: flex;
		overflow: hidden;
	}

	.argid {
		flex: 1;
		height: 100%;
		width: 40%;
	}

	.aside {
		width: 60%;
		display: flex;
		flex-flow: column;
		height: 100%;
		overflow: auto
	}
	
}

	
</style>
