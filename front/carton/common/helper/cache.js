import config from '../conifg.js';
import CACHE_CONFIG from '../enum/cache_config.js'
import {
	getCurrentTimeStamp
} from './time.js';


/**
 * @author Azal
 * @date 2020/4/29 09:56
 * @description 缓存
 */
class Cache {
	constructor() {}

	_buildCacheKey(key) {
		return `${key}${config.brand_id}${config.environment}${config.domain.substr(-6)}`;
	}


	/**
	 * 检查版本
	 */
	checkVersion() {
		let version = config.version;
		let storageVersion = uni.getStorageSync('version');
		if (storageVersion !== version) {
			// uni.clearStorageSync(); 版本不一致更新全部缓存
			uni.setStorageSync('version', version);
		}
	}


	/**
	 * 读取缓存
	 */
	get(key) {
		let data = uni.getStorageSync(this._buildCacheKey(key));
		if (data && data instanceof Object) {
			if (data.hasOwnProperty('expired_time')) {
				if (getCurrentTimeStamp() > data.expired_time) {
					this.remove(key);
					return null;
				}
				return data.data;
			}
			return data;
		}
		return data || null;
	}

	/**
	 * 设置缓存
	 * @param {String} key
	 * @param {Object} data
	 */
	set(key, data) {
		let config = CACHE_CONFIG[key.toUpperCase()]
		if (config && config.ttl) {
			data = {
				data,
				expired_time: config.ttl + getCurrentTimeStamp()
			};
		}
		uni.setStorageSync(this._buildCacheKey(key), data);
	}

	/**
	 * 异步设置缓存
	 * @param {Object} key
	 * @param {Object} data
	 * @param {Int} _ttl 时效 秒
	 */
	setSync(key, data, _ttl = 0) {
		if (_ttl) {
			data = {
				data,
				expired_time: _ttl + getCurrentTimeStamp()
			};
		}
		key = this._buildCacheKey(key);
		return new Promise((resolve, reject) => {
			uni.setStorage({
				key,
				data,
				success: () => {
					return resolve(!0);
				},
				fail: () => {
					return reject(!1);
				}
			});
		});
	}
	
	batchRemove(keys) {
		keys.forEach((key) => {
			this.remove(key)
		})
	}

	/**
	 * 清除缓存
	 * @param {Object} key
	 */
	remove(key) {
		uni.removeStorageSync(this._buildCacheKey(key));
	}
}

export default Cache;
