import Vue from 'vue'
import axios from 'axios'
import VueAxios from 'vue-axios'
import VueContext from 'vue-context';
import App from './App.vue'
import router from './router'
import 'bootstrap'; import 'bootstrap/dist/css/bootstrap.min.css'
import 'colors.css/css/colors.min.css';
import {BootstrapVue} from 'bootstrap-vue'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import {store} from '@/store';

Vue.use(VueAxios, axios);
Vue.use(BootstrapVue)
Vue.config.productionTip = false

new Vue({
  router,
  components: {VueContext},
  render: h => h(App),
  store
}).$mount('#app')
