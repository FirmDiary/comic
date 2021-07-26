//api
const API = {
	'login': {
		label: '登录',
		value: 'user/login',
		method: 'post',
		no_auth: true,
	},
	'transfer_u2': {
		label: '转换图片',
		value: 'upload/transfer',
		method: 'post',
	},
	'transfer_old_fix': {
		label: '修复老照片',
		value: 'upload/transferOldFix',
		method: 'post',
	},
}

export default API;