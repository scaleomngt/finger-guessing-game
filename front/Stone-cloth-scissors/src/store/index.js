import Vue from "vue";
import Vuex from 'vuex'
import socket from './modules/socket'
import getters from './getters'


Vue.use(Vuex)


const store = new Vuex.Store({
    modules: {
        socket,
    },
    getters
  })
  
  export default store