//api
const API = {
	'login': {
		label: '登录',
		value: 'user/login',
		method: 'post',
		no_auth: true,
	},
	'transfer': {
		label: '转换图片',
		value: 'upload/transfer',
		method: 'post',
	},
}

export default API;