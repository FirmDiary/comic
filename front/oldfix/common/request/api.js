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