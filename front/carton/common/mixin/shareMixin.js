const shareMixin = {
	onShareAppMessage() {
		return {
			title: '神奇的自拍风格化!',
			path: '/page/index/index',
		}
	}
};
export default shareMixin;
