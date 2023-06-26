<template>
  <div id="app">
    <router-view />
  </div>
</template>
<script>
import store from "@/store";
import {  createWebSocket } from "./utils/socket";
export default {
  mounted(){
    createWebSocket(this.global_callback);
  },
  methods:{
    global_callback(msg) {
      if (msg.Cmd == 'loginsucess') {
        this.$router.replace({path:'/game'});
      }else if(msg.Cmd == 'full'){
        store.dispatch('socket/setMsg', msg)
      }else if(msg.Cmd == 'result'){
        store.dispatch('socket/setMsg', msg)
      }
    },   
  },
}
</script>
<style>
* {
  padding: 0;
  margin: 0;
}
</style>
