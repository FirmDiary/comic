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
	'transfer_carton': {
		label: '转化卡通',
		value: 'upload/transferCarton',
		method: 'post',
	},
	'carton_etc': {
		label: '获取案例',
		value: 'common/carton/etc',
		method: 'get',
		no_auth: true,
	},
	'transfer_old_fix': {
		label: '修复老照片',
		value: 'upload/transferOldFix',
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