export const is_phone = (phone) => {
	let regExp = /^((13[0-9])|(14[5,7,9])|(15[0-3,5-9])|(16[2,5-7])|(17[0-3,5-8])|(18[0-9])|(19[1,3,5,8-9]))\d{8}$/;
	return phone && regExp.test(phone);
}

export const is_password = (password) => {
	let regExp = /^(?![0-9]+$)(?![a-zA-Z]+$)[0-9A-Za-z]{6,20}$/;
	return password && regExp.test(password)
}

export const is_email= (email) => {
	let regExp = /^([A-Za-z0-9_\-\.\u4e00-\u9fa5])+\@([A-Za-z0-9_\-\.])+\.([A-Za-z]{2,8})$/;
	return email && regExp.test(email)
}
