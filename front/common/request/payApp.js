import go from './go.js';

import {
	showToast,
} from '../helper/base.js';

import API from '../enum/api.js';


/**
 * 支付类
 * @pay_for 支付对象
 */
class PayApp {
	constructor(pay_for) {
		this.pay_api = API[pay_for].value;
		this.after_pay_api = API[pay_for].after;
	}

	/**
	 * @param {Object} data
	 */
	async pay(request_params) {
		let orderResult = await this._reuqestOrder(request_params);
		if (orderResult.data.status === 1) {
			//已经优惠掉全部金额
			showToast('支付成功');
			return {
				status: orderResult.data.status,
				actuallyAmount: orderResult.payInfo.actuallyAmount
			};
		} else {
			return await this.reuqestMiniPay(orderResult);
		}
	}
	
	/**
	 * 吊起微信支付
	 */
	reuqestMiniPay(orderResult) {
		let order_id = orderResult.data.id;
		return new Promise((resolve, reject) => {
			wx.requestPayment({
				...orderResult.payInfo.gateways.wxMini.jsPayInfo,
				success: (res) => {
					showToast('支付成功');
					if (this.after_pay_api) {
						//需要回调通知后台验证支付状态的
						this._afterPay(order_id).then((res) => {
							return resolve({
								status: res.data.status,
								order_status:res.data.order_status,
								order_id,
							});
						});
					} else {
						return resolve({
							status: 1,
							order_id,
						});
					}
				},
				fail: (err) => {
					console.log('支付fail:' + JSON.stringify(err));
					if ('requestPayment:fail cancel' === err.errMsg) {
						showToast('您已取消支付');
						return resolve({
							status: '-1',
							order_id,
							orderResult: orderResult,
						});
					} else {
						return reject(false);
					}
				},
			});
		});
	}

	/**
	 * 下单并取得支付必须参数
	 * @param {Object} request_params
	 */
	_reuqestOrder(request_params) {
		return new Promise((resolve) => {
			go(
				this.pay_api,
				request_params,
				'post'
			).then((res) => {
				return resolve(res.data);
			});
		});
	}

	_afterPay(order_id) {
		return new Promise((resolve) => {
			go(
				this.after_pay_api, {
					order_id
				},
				'post'
			).then((res) => {
				return resolve(res.data);
			});
		});
	}
}

export default PayApp;
