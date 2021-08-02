/**
 * 初始化Date函数
 * @param {Object} fmt
 */
function DateTimeInit(fmt) {
	var o = {
		'M+': this.getMonth() + 1, //月份 
		'd+': this.getDate(), //日 
		'h+': this.getHours(), //小时 
		'm+': this.getMinutes(), //分 
		's+': this.getSeconds(), //秒 
		'q+': Math.floor((this.getMonth() + 3) / 3), //季度 
		'S': this.getMilliseconds() //毫秒 
	};
	if (/(y+)/.test(fmt)) {
		fmt = fmt.replace(RegExp.$1, (this.getFullYear() + '').substr(4 - RegExp.$1.length));
	}
	for (var k in o) {
		if (new RegExp('(' + k + ')').test(fmt)) {
			fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (('00' + o[k]).substr(('' + o[k]).length)));
		}
	}
	return fmt;
}


/**
 * 时间格式化
 * @param {String} format 时间格式
 * @param {Time} time   
 */
function timeFormat(format, time = '') {
	time = time ? new Date(time.replace(/\-/g, '/')).getTime() : new Date();
	return new Date(time).format(format);
}

/**
 * 获取当天日期
 */
function getNowFormatDate() {
	let date = new Date();
	let seperator1 = '-';
	let year = date.getFullYear();
	let month = date.getMonth() + 1;
	let strDate = date.getDate();
	if (month >= 1 && month <= 9) {
		month = '0' + month;
	}
	if (strDate >= 0 && strDate <= 9) {
		strDate = '0' + strDate;
	}
	let currentdate = year + seperator1 + month + seperator1 + strDate;
	return currentdate;
}


/**
 * 获取这周日期
 * @param {Int} i 周几  
 * @param {Date} 想获取那周的日期
 */
function getWeek(i, date = '') {
	let now = date ? new Date(date) : new Date();
	// let firstDay = new Date(now - (now.getDay()-1) * 86400000); 减去1是从周一开始
	let firstDay = new Date(now - now.getDay() * 86400000);
	firstDay.setDate(firstDay.getDate() + i);
	let mon = Number(firstDay.getMonth()) + 1;
	mon = mon / 10 >= 1 ? mon : ('0' + mon);
	return mon + '-' + (firstDay.getDate() / 10 < 1 ? '0' + firstDay.getDate() :
		firstDay.getDate());
}


function getCurrentTimeStamp() {
	return parseInt(new Date().getTime() / 1000);
}


export {
	DateTimeInit,
	timeFormat,
	getNowFormatDate,
	getWeek,
	getCurrentTimeStamp,
};