//api
const API = {
	'login': {
		label: '登录',
		value: 'user/login',
		method: 'post',
		no_auth: true,
	},
	'user': {
		label: '获取用户信息',
		value: 'user/info',
		method: 'get',
	},
	'transfer_old_fix': {
		label: '修复老照片',
		value: 'upload/transferOldFixMT',
		method: 'post',
	},
	'transfer_waifu_2x': {
		label: '放大高清图片',
		value: 'upload/transfer2x',
		method: 'post',
	},
	'old_etc': {
		label: '获取案例',
		value: 'common/old/etc',
		method: 'get',
		no_auth: true,
	},
	'invite_success': {
		label: '邀请成功',
		value: 'invite/success',
		method: 'get',
	},
}

export default API;