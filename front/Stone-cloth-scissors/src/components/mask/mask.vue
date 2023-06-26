<template>
  <div class="mask" v-if="show">
    <slot></slot>
  </div>
</template>

<script>
export default {
    props:{
        show:{
            type:Boolean,
            default:false,
        } 
    }
}
</script>

<style lang="less" scoped>
*{
    position: relative;
    .mask{
        width: 100vw;
        height: 100vh;
        background-color: rgba(161, 159, 159, 0.7);
        position: absolute;
        top: 0;
        z-index: 1 !important;
        display: flex;
        justify-content: center;
        align-items: center;
        // backdrop-filter: blur(3px);
    }
}
</style>