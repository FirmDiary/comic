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
	_api={}
	_header = {
		'content-type': 'application/json',
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
	to(api, data = {}, ...query) {
		if (!API.hasOwnProperty(api)) {
			console.error('api不存在!')
			return;
		}
		this._api = API[api]

		let url = `${this.host}/${this._api.value}`;
		if (query && query.length) {
			query.forEach((que) => {
				url += '/' + que
			})
		}
		if (this._goSetting.loading) {
			showLoading()
		}
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
						showToast(res.data.message)
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
				console.log(token);
				if (token) Object.assign(this._header, {
					Authorization: "Bearer " + token
				})
			}
			params.method = params.method ? params.method.toUpperCase() : 'GET'
			uni.request({
				...params,
				header:this._header,
				success: (res) => {
					return resolve(res)
				},
				fail: (err) => {
					return reject(err)
				},
				complete:() => {
				}
			})
		})
	}

}
export default Go;





// import config from '../conifg.js';

// import {
// 	showLoading,
// 	showToast,
// } from '../helper/base.js';

// import Cache from '../helper/cache.js';
// const cache = new Cache();

// const HEADER = {
// 	'content-type': 'application/json',
// };

// /**
//  * @author Azal
//  * @date 2020/4/29 10:02
//  * @description http请求
//  */
// const go = (path, data = {}, type = 'get', {
// 	special = false, //特殊不需要验证是否选择门店的
// 	loading = true, //是否要loading
// 	full_url = false, //是否是完整路径，否会自动组装
// 	show_err = true, //如果有错误信息返回 是否弹出
// } = {}) => {

// 	let _loading = true;
// 	if (loading && type !== 'post') {
// 		_loading = config.openLoading;
// 	} else {
// 		_loading = loading;
// 	}
// 	if (_loading) {
// 		showLoading();
// 	}
// 	let shop = cache.get('shop_info');
// 	if (!shop && !special) {
// 		uni.redirectTo({
// 			url: '/pages/store/select'
// 		});
// 		return;
// 	}
// 	return new Promise((resolve, reject) => {
// 		let token = cache.get('token');
// 		let header = HEADER;
// 		if (token) {
// 			Object.assign(header, {
// 				Authorization: token
// 			});
// 		}
// 		let url = full_url ? path : `${config.host}/${config.brand_id}/${shop ? shop.id : 0}/${path}`;
// 		uni.request({
// 			url,
// 			data,
// 			header,
// 			method: type ? type.toUpperCase() : 'GET',
// 			success: (res) => {
// 				if (res.statusCode === 401) {
// 					//token过期
// 					uni.clearStorage();
// 					uni.reLaunch({
// 						url: '/pages/index/index'
// 					});
// 					return reject(res);
// 				}
// 				if (res.data.code === 0) {
// 					return resolve(res);
// 				} else {
// 					console.log('Fail→');
// 					console.log(JSON.stringify(res));
// 					if (show_err) {
// 						showToast(res.data.message);
// 					}
// 					return reject(res);
// 				}
// 			},
// 			fail: (err) => {
// 				console.log('请求错误' + JSON.stringify(err));
// 				showToast('请将错误页面截屏发送给我们，以便进行问题追踪。');
// 				return reject(err);
// 			},
// 		});
// 	});
// };



// export default go;
