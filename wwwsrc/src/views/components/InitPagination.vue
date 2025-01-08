<template>
    <div>
        <el-pagination
				@size-change="onHandleSizeChange"
				@current-change="onHandleCurrentChange"
				class="mt15"
				:pager-count="5"
				:page-sizes="[10, 20, 30]"
				v-model:current-page="pages.pageNum"
				background
				v-model:page-size="pages.pageSize"
				layout="total, sizes, prev, pager, next, jumper"
				:total="pages.total"
			>
			</el-pagination>
    </div>
</template>
<script>
export default {
    props:{
        pages:{
            type:Object,
            default:()=>{
                return{
                    total:0,
                    pageNum:1,
                    pageSize:10
                }
            }
        }
    },
    methods:{
         onHandleSizeChange(val){
            this.$emit('sizeChange',val)
        },
        // 分页改变
        onHandleCurrentChange(val){
            this.$emit('currentChange',val)
        },
    }
}
</script>