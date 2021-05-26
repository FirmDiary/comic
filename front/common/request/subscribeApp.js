import config from '../conifg.js';
import go from './go.js';

import Cache from '../helper/cache.js';
const cache = new Cache();


/**
 * 订阅消息
 * @subscribe_for 发送场景  
 * personal_train 私教课预约成功  
 * group_train 团课预约成功
 */
class SubscribeApp {
	constructor(subscribe_for) {
		this.templateIds = [];
		this.subscribe_for = subscribe_for;
		this.templates = [];
	}

	async requestSubscribe() {
		let subscribeTemplateValid = await this.getSubscribeTemplates();
		if (!subscribeTemplateValid) {
			// this.errorToast('未配置消息订阅')
			return false;
		}
		await this.initTemplates(subscribeTemplateValid);
		let requestResult = await this.reuqestMiniSubscribe();
		let resultTemplateData = [];
		for (let key in this.templates) {
			if (requestResult[this.templates[key].priTmplId] === 'accept') {
				resultTemplateData.push({
					id: this.templates[key].id,
					page: this.templates[key].page,
					page_params: this.templates[key].page_params,
				});
			}
		}

		if (!resultTemplateData.length) {
			return false;
		}
		this.saveScribeToOutBox(resultTemplateData);
		return true;
	}

	/**
	 * 存储成功订阅的消息至服务端发件箱
	 * @param {array} templateIds
	 */
	saveScribeToOutBox(resultTemplateData) {
		let user = cache.get('user');
		if (user && user.is_member) {
			go(`${config.domain}/api/mini/subscribe-message-template/insert`, {
				receiver_type: 1, //代表会员
				receiver_id: user.member_id,
				templates: resultTemplateData,
			}, 'post', {
				full_url: true,
				loading: false,
				show_err: false
			}).then((res) => {
				console.log(res);
			});
		}
	}

	/**
	 * 初始化即将申请的模板消息
	 */
	initTemplates(templatesExsits) {
		this.templates = templatesExsits[this.subscribe_for];
		return new Promise((resolve, reject) => {
			try {
				for (let sign in this.templates) {
					this.templateIds.push(this.templates[sign].priTmplId);
				}
			} catch (e) {
				return reject(false);
			}
			return resolve(true);
		});
	}

	reuqestMiniSubscribe() {
		return new Promise((resolve, reject) => {
			wx.requestSubscribeMessage({
				tmplIds: this.templateIds,
				success: (res) => {
					return resolve(res);
				},
				fail: (err) => {
					return reject(err);
				}
			});
		});
	}

	// 获取可以使用的订阅消息模板
	getSubscribeTemplates() {
		return new Promise((resolve, reject) => {
			let shop = cache.get('shop_info');
			if (!shop) {
				return reject(!1);
			}
			go(`${config.domain}/api/mini/subscribe-message-template/list`, {
				app: 1, //代表会员端
				brand_id: config.brand_id,
				shop_id: shop.id,
			}, 'get', {
				full_url: true,
				loading: false,
			}).then((res) => {
				return resolve(res.data.data);
			});
		});
	}

}

export default SubscribeApp;
