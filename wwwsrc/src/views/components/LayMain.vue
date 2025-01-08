<template>
    <div class="system-user-container layout-padding w100">
        <el-card class="l_main">
            <header class="lm_header">
                <slot name="header"></slot>
            </header>
            <main class="lm_main">
                <slot name="main"></slot>
            </main>
            <footer class="lm_footer">
                <slot name="footer"></slot>
            </footer>
        </el-card>
    </div>
</template>
<script>
export default {
    name: 'LayMain'
}

</script>
<style>
.l_main {
    width: 100%;
    height: 100%;
}
.l_main .el-card__body{
    width: 100%;
    height: 100%;
    display: flex;
    flex-flow: column;
}
.lm_header{
    display: flex;
    align-items: center;
    margin-bottom: 5px;
}
.lm_header .el-form-item--large{
    margin-bottom: 0;
}
.lm_main{
    flex: 1;
    overflow: auto;
}
</style>