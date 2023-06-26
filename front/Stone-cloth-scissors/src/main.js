import Vue from 'vue'
import App from './App.vue'
import {Button,Input} from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import router from './router'
// 国际化(internationalization)
import VueI18n from 'vue-i18n'
import {zh} from './language/zh/zh.js'
import {en} from './language/en/en'

import store from './store'

import { getDefaultLang } from './utils/utils'


Vue.use(Button);
Vue.use(Input);
Vue.use(VueI18n);

Vue.config.productionTip = false

const i18n = new VueI18n({
  locale:getDefaultLang(),
  messages:{
    en:en,
    zh:zh
  }
})

new Vue({
  router,
  store,
  i18n,
  render: h => h(App)
}).$mount('#app')
