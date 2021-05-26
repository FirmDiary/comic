import config from '../conifg.js';
import go from './go.js';

/**
 * 获取某页的小程序码
 */
export const getPageQr = (scene, page = '', width = '') => {
	let data = {
		brand_id: config.brand_id,
		scene,
	};
	if (page) data.page = page;
	if (width) data.width = width;
	return new Promise((resolve) => {
		go(`${config.domain}/api/mini/qr/page`, data, 'post', {
			full_url: true
		}).then((res) => {
			import('@/common/helper/utils.js').then((utils) => {
				let {
					bufferToFile,
					getBase64ImageUrl
				} = utils
				return resolve(bufferToFile(getBase64ImageUrl(res.data.data)));
			})
		});
	});
};



/**
 * 获取二维码
 * @param {String} code
 * @return {Base64}
 */
export const code2Qr = (code) => {
	return new Promise((resolve) => {
		go(`${config.domain}/api/mini/qr/${code}`, [], 'get', {
			full_url: true
		}).then((res) => {
			import('@/common/helper/utils.js').then((utils) => {
				let {
					getBase64ImageUrl,
				} = utils
				return resolve(getBase64ImageUrl(res.data.data.qr));
			})
		});
	});
};



/**
 * 获取小程序自定义配置
 */
export const getCustomInfo = () => {
	return new Promise((resolve, reject) => {
		go(`${config.domain}/api/mini/custom/config`, {
			app: 1, //代表会员端
			brand_id: config.brand_id,
		}, 'get', {
			full_url: true,
			loading: false,
			special: true,
		}).then((res) => {
			return resolve(res.data.data);
		});
	});
};


export const getMiniCode = () => {
	return new Promise((resolve, reject) => {
		uni.login({
			success: res => {
				return resolve(res.code);
			},
			fail: err => {
				return reject(err);
			}
		});
	});
};

/**
 * 尝试登陆
 * 成功返回用户信息
 * 未注册用户返回false
 */
export const login = () => {
	return new Promise((resolve) => {
		getMiniCode().then((code) => {
			go('authorizations', {
				code
			}, 'post', {
				'special': true
			}).then((res) => {
				let userData = res.data.data
				import("@/common/helper/cache.js").then((cache) => {
					let Cache = new cache.default()
					if (res.data.code == 0) {
						let user = Object.assign(userData.user || {}, {
							is_member: userData.is_member,
						})
						Cache.set('user', user, userData.expires_in)
						Cache.set('token', userData.access_token, res.data.data.expires_in)
						return resolve(user);
					}
					return resolve(false);
				})
			}, (err) => {
				return resolve(false);
			});
		})
	});
};

/**
 * 获取门店信息
 * @param {Object} Cache对象
 */
export const getShopInfo = (Cache) => {
	const key = 'shop_info';
	return new Promise((resolve) => {
		import("@/common/helper/cache.js").then((cache) => {
			let Cache = new cache.default()
			let info = Cache.get(key);
			if (!info) {
				go('shop').then((res) => {
					Cache.set(key, res.data.data, 60 * 60 * 24 * 7);
					return resolve(res.data.data);
				});
			} else {
				return resolve(info);
			}
		})
	});
};



/**
 * 获取商城基础配置
 * @param {Object} Cache对象
 */
export const getMarketBasicSetting = (Cache) => {
	const key = 'marketBasicSetting';
	return new Promise((resolve) => {
		let setting = Cache.get(key);
		if (!setting) {
			go('plugin/market/setting/detail').then((res) => {
				Cache.set(key, res.data.data, 60 * 60 * 10);
				return resolve(res.data.data);
			});
		} else {
			return resolve(setting);
		}
	});
};

/**
 * 获取商城商品
 */
export const getMarketProducts = (params) => {
	return new Promise((resolve) => {
		go('plugin/market/product/list', params).then((res) => {
			return resolve(res.data);
		});
	});
};


/**
 * 商城取消订单
 */
export const cancelOrder = (id) => {
	return new Promise((resolve) => {
		go(`plugin/market/order/cancel/${id}`, {}, 'post').then((res) => {
			return resolve(res.data);
		});
	});
};

/**
 * 商城继续支付
 */
export const payOrder = (order_id) => {
	return new Promise((resolve) => {
		go(`plugin/market/order/pay`, {
			order_id
		}, 'post').then((res) => {
			return resolve(res.data);
		});
	});
};

/**
 * 商城确认收货
 */
export const confirmOrder = (order_id) => {
	return new Promise((resolve) => {
		go(`plugin/market/order/confirm_receive/${order_id}`, {}, 'post').then((res) => {
			return resolve(res.data);
		});
	});
};


/**
 * 商城取消售后订单
 */
export const cancelAfterSaleOrder = (id) => {
	return new Promise((resolve) => {
		go(`plugin/market/order_return/cancel/${id}`, {}, 'post').then((res) => {
			return resolve(res.data);
		});
	});
};

