export const warnRelaunch = (url = '/pages/index/index', title = '异常访问!') => {
	uni.showModal({
		showCancel: false,
		title: '提示',
		content: title,
		success: () => {
			reLaunch(url);
		}
	});
};

export const showModal = (title = '发生异常了!', config = {
	showCancel: false
}) => {
	return new Promise((resolve) => {
		uni.showModal({
			showCancel: config.showCancel,
			title: '提示',
			content: title,
			success: (res) => {
				return resolve(res.confirm);
			}
		});
	});
};


export const reLaunch = (url = '/pages/index/index') => {
	uni.reLaunch({
		url
	});
};

export const navigateTo = (url) => {
	uni.navigateTo({
		url,
	});
};


export const navigateBack = () => {
	uni.navigateBack({});
};

export const redirectTo = (url) => {
	uni.redirectTo({
		url,
	});
};


export const showToast = (title, position = 'center', duration = 1500) => {
	uni.showToast({
		title,
		position,
		duration,
		icon: 'none',
	});
};


export const showLoading = (title = '数据加载中', mask = true) => {
	console.log(title);
	uni.showLoading({
		title: title,
		mask: mask
	});
};


export const currentPage = (prefix = '/') => {
	let pages = getCurrentPages();
	return prefix + pages[pages.length - 1].route;
};
