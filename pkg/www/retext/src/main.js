import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import App from './App.vue'
import router from './router'
import { BootstrapVue } from 'bootstrap-vue'

Vue.use(VueAxios, axios);
Vue.use(BootstrapVue)
Vue.config.productionTip = false

new Vue({
  router,
  render: h => h(App)
}).$mount('#app')

import 'bootstrap-vue/dist/bootstrap-vue.css'