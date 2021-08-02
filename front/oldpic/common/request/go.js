import config from '@/common/conifg.js';
import {
	showToast,
	showLoading,
	reLaunch,
} from "@/common/helper/base.js";
import API from "@/common/enum/api.js";
import store from '@/store/index.js';

class Go {
	_goSetting = {}
	_api = {}
	_header = {
		'content-type': 'application/json',
		'app_id': config.app_id,
	}

	constructor() {
		this.host = config.host
		this.domain = config.domain
		this._resetGoSetting()
	}

	set(setting) {
		this._goSetting = setting
		return this
	}

	_resetGoSetting() {
		this._goSetting = {
			errToast: true, //是否将错误信息toast
			loading: config.openLoading, //是否显示loading
		}
	}

	/**
	 * http 请求
	 */
	to(api, data = {}) {
		if (!API.hasOwnProperty(api)) {
			console.error('api不存在!')
			return;
		}
		this._api = API[api]

		let url = `${this.host}/${this._api.value}`;
		// console.log(this._goSetting);
		// if (this._goSetting.loading) {
		// 	showLoading()
		// }
		return new Promise((resolve, reject) => {
			this.request({
				url,
				data,
				method: this._api.method,
			}).then((res) => {
				if (res.statusCode === 401) {
					//token过期
					reLaunch('/pages/index/index')
				}
				if (res.data.code === 0) {
					return resolve(res.data)
				} else {
					if (this._goSetting.errToast) {
						showToast(res.data.msg)
					}
					return reject(res)
				}
			}).catch((err) => {
				console.debug('请求错误' + JSON.stringify(err))
				if (this._goSetting.errToast) {
					showToast("请将错误页面截屏发送给我们，以便进行问题追踪。")
				}
				return reject(err)
			}).finally(() => {
				this._resetGoSetting()
				wx.hideLoading()
			})
		});
	}

	request(params = {
		url: "",
		data: "",
		method: "",
	}) {
		return new Promise((resolve, reject) => {
			if (!this._api.no_auth) {
				let token = store.state.login.token
				if (token) Object.assign(this._header, {
					Authorization: "Bearer " + token
				})
			}
			params.method = params.method ? params.method.toUpperCase() : 'GET'
			uni.request({
				...params,
				header: this._header,
				success: (res) => {
					return resolve(res)
				},
				fail: (err) => {
					return reject(err)
				},
				complete: () => {}
			})
		})
	}

}
export default Go;
