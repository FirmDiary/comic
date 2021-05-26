import config from '../conifg.js';

import {
	showToast,
} from '../helper/base.js';

import API from "../enum/api.js"

/**
 * 上传类
 */
class Upload {
	constructor() {}

	/**
	 * 批量上传图片
	 * @param {Object} paths
	 */
	uploadMutipleImg(paths) {
		return new Promise((resolve) => {
			let img_url_ok = [];
			paths.forEach((path) => {
				this.uploadImg(path).then((res) => {
					img_url_ok.push(res.data.path);
					if (img_url_ok.length == paths.length) {
						return resolve(img_url_ok);
					}
				});
			});
		});
	}

	/**
	 * 上传单图
	 * @param {Object} path
	 */
	uploadImg(path, header = {}) {
		return new Promise((resolve, reject) => {
			wx.uploadFile({
				url: `${config.domain}/${API.transfer.value}`,
				filePath: path,
				name: 'file',
				header,
				success: (res) => {
					return resolve(JSON.parse(res.data));
				},
				fail: (err) => {
					return reject(err);
				},
			});
		});
	}
}

export default Upload;
