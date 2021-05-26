import Vue from 'vue'
import App from '@/App'


import store from '@/store/index.js';
import * as base from '@/common/helper/base.js';
import Go from '@/common/request/go.js';
import Cache from '@/common/helper/cache.js';

import {
	DateTimeInit
} from "@/common/helper/time.js"

Vue.prototype.$base = base
Vue.prototype.$store = store
Vue.prototype.$go = new Go()
Vue.prototype.$cache = new Cache()

Date.prototype.format = DateTimeInit
Vue.config.productionTip = false

App.mpType = 'app'

const app = new Vue({
    ...App
})
app.$mount()

 



