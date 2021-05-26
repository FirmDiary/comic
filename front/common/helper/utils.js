/**
 * 保存画布图片至手机相册
 * @param {String} canvas_id
 */
export const saveCanvasToPhone = async (canvas_id) => {
	return new Promise((resolve, reject) => {
		wx.canvasToTempFilePath({
			canvasId: canvas_id,
			quality: 1,
			success: (res) => {
				let path = res.tempFilePath;
				wx.getSetting({
					success: (res) => {
						if (!res.authSetting['scope.writePhotosAlbum']) {
							wx.authorize({
								scope: 'scope.writePhotosAlbum',
								success: () => {
									saveToPhone(path).then((res) => {
										return resolve(res);
									});
								},
								fail: (res) => {
									return reject(res);
								}
							});
						} else {
							saveToPhone(path).then((res) => {
								return resolve(res);
							});
						}
					}
				});
			},
			fail: (res) => {
				return reject(res);
			}
		});
	});
};

/**
 * 保存图片至本地
 * @param {String} path
 */
export const saveToPhone = (path) => {
	return new Promise((resolve) => {
		uni.saveImageToPhotosAlbum({
			filePath: path,
			success: (res) => {
				return resolve(res);
			},
		});
	});
};


/**
 * canvas 文字单行截断
 * @param {Object} ctx canvas对象
 * @param {String} text
 * @param {Int} size
 * @param {Int} width
 */
export const textCut = (ctx, text, size, width) => {
	ctx.setFontSize(size);
	if (ctx.measureText(text).width < width) {
		return text;
	}
	let char_arr = text.split('');

	let res_text = '';
	let i = 0;
	let step_width = 0; //单字长
	let full_width = 0; //全长
	do {
		res_text += char_arr[i++];
		full_width = ctx.measureText(res_text).width;
		step_width = full_width - step_width;
	} while (full_width + step_width < width && i + 1 < char_arr.length);

	res_text += (i + 1 < char_arr.length) ? '...' : '';

	return res_text;
};

//canvas 文字换行
// textWrap() {}

/**
 * 下载网络资源
 * @param {String} url
 */
export const downloadFile = (url) => {
	return new Promise((resolve, reject) => {
		wx.downloadFile({
			url: url,
			success: (res) => {
				return resolve(res.tempFilePath);
			},
			fail: (err) => {
				console.log(err, 'err');
				return reject(!!err);
			}
		});
	});
};

/**
 * Transfer Buffer to File
 * @param {String} url
 */
export const bufferToFile = (base64Data) => {
	const fs = wx.getFileSystemManager();
	const FILE_BASE_NAME = 'tmp_base64imgsrc';
	const [, format, bodyData] = /data:image\/(\w+);base64,(.*)/.exec(base64Data) || [];
	if (!format) {
		return (new Error('ERROR_PARSE'));
	}
	// const filePath = `${wx.env.USER_DATA_PATH}/${FILE_BASE_NAME+Date.parse(new Date())}.${format}`;
	const filePath = `${wx.env.USER_DATA_PATH}/${FILE_BASE_NAME}.${format}`;
	const buffer = wx.base64ToArrayBuffer(bodyData);

	return new Promise((resolve, reject) => {
		// fx.removeSavedFile({
		// })
		fs.writeFile({
			filePath,
			data: buffer,
			encoding: 'binary',
			success: () => {
				return resolve(filePath);
			},
			fail: (res) => {
				console.log(res);
				reject(new Error('ERROR_BASE64SRC_WRITE'));
			},
		});
	});
};




/**
 * 数字转换中文
 */
export const fuckNumberToChinese = (str) => {
	str = ('' + str).trim().replace(/^0*/, ''); //去掉前面修饰的0
	let match = ['', '一', '二', '三', '四', '五', '六', '七', '八', '九', '零'];
	return ('0000' + str).substr((str.length % 4) || 4).replace(/(\d){4}/g, function(_str, endIndex, startIndex) {
		let dot = (((str.length - 1) / 4) >> 0) - ((startIndex / 4) >> 0);
		let prefix = (function getPrfix(dot) {
			return (dot > 2 ? (+_str ? (dot == 3 ? '万' : (getPrfix(dot - 1) + '万')) : '') : (dot == 1 ? (+_str ? '万' : '') :
				(dot == 2 ? '亿' : '')));
		})(dot);
		/0+$/g.test(_str) && (prefix += match[10]); //处理单元内后半部分有零的地方
		return (+_str) ? (_str.replace(/(\d)(\d)(\d)(\d)/g, function($0, $1, $2, $3, $4) {
			!match[$1] && (match[$2] ? ($1 = 10) : (match[$3] ? ($2 = 10) : (match[$4] ? ($3 = 10) : ''))); //处理相邻单元前半部分
			match[$1] && match[$3] && !match[$2] && ($2 = 10), match[$2] && match[$4] && !match[$3] && ($3 = 10), match[
				$1] && match[$4] && !match[$3] && !match[$2] && ($3 = 10); //中间两个连续为0，只是获取最后一个
			return (match[$1] && ($1 < 10 ? (match[$1] + '千') : match[$1])) + (match[$2] && ($2 < 10 ? (match[$2] + '百') :
				match[$2])) + (match[$3] && ($3 < 10 ? ($3 == 1 ? '十' : (match[$3] + '十')) : match[$3])) + (match[$4] &&
				match[$4]);
		}) + prefix) : (prefix);
	}).replace(/^零*/g, '').replace(/零*$/g, '').replace(/(零)*/g, '$1').replace(/零亿/g, '亿') || match[10]; //处理连续零的问题
};


/**
 * 组装base64文件头
 * @param {Object} base64Data
 */
export const getBase64ImageUrl = (base64Data) => {
	let base64ImgUrl = 'data:image/png;base64,' + base64Data;
	return base64ImgUrl;
};


/**
 * 复制
 */
export const copyText = (text) => {
	uni.setClipboardData({
		data:text.toString(),
		success(res) {
			wx.getClipboardData({});
		},
	});
};


//计算两地距离
export const getDistance = (lat1, lng1, lat2, lng2) => {
	let radLat1 = lat1 * Math.PI / 180.0;
	let radLat2 = lat2 * Math.PI / 180.0;
	let a = radLat1 - radLat2;
	let b = lng1 * Math.PI / 180.0 - lng2 * Math.PI / 180.0;
	let s = 2 * Math.asin(Math.sqrt(Math.pow(Math.sin(a / 2), 2) +
		Math.cos(radLat1) * Math.cos(radLat2) * Math.pow(Math.sin(b / 2), 2)));
	s = s * 6378.137; // EARTH_RADIUS;
	s = Math.round(s * 10000) / 10000;
	return s.toFixed(2);
};



/**
 * 颜色16进制转rgb
 */
export const colorRgb = (color, opacity = 1) => {
	color = color.toLowerCase();
	//十六进制颜色值的正则表达式
	var reg = /^#([0-9a-fA-f]{3}|[0-9a-fA-f]{6})$/;
	// 如果是16进制颜色
	if (color && reg.test(color)) {
		if (color.length === 4) {
			var colorNew = '#';
			for (var i = 1; i < 4; i += 1) {
				colorNew += color.slice(i, i + 1).concat(color.slice(i, i + 1));
			}
			color = colorNew;
		}
		//处理六位的颜色值
		var colorChange = [];
		for (var i = 1; i < 7; i += 2) {
			colorChange.push(parseInt('0x' + color.slice(i, i + 2)));
		}
		return 'RGB(' + colorChange.join(',') + ',' + opacity + ')';
	}
	return color;
};


/**
 * 绘制圆角图案
 * @param {CanvasContext} ctx canvas上下文
 * @param {number} x 圆角矩形选区的左上角 x坐标
 * @param {number} y 圆角矩形选区的左上角 y坐标
 * @param {number} w 圆角矩形选区的宽度
 * @param {number} h 圆角矩形选区的高度
 * @param {number} r 圆角的半径
 * @param {number} bgC 周圈的背景颜色
 */
export const roundRect = (ctx, x, y, w, h, r, bgC) => {
	// 开始绘制
	ctx.beginPath();
	// ctx.setFillStyle(c);
	ctx.fillStyle = bgC;
	// 左上角
	ctx.moveTo(x + r, y);
	ctx.arcTo(x, y, x, y + r, r);
	ctx.lineTo(x, y);
	ctx.lineTo(x + r, y);
	ctx.fill();

	// 右上角
	ctx.moveTo(x + w - r, y);
	ctx.arcTo(x + w, y, x + w, y + r, r);
	ctx.lineTo(x + w, y);
	ctx.lineTo(x + w - r, y);
	ctx.fill();

	// 左下角
	ctx.moveTo(x, y + h - r);
	ctx.arcTo(x, y + h, x + r, y + h, r);
	ctx.lineTo(x, y + h);
	ctx.lineTo(x, y + h - r);
	ctx.fill();

	// 右下角
	ctx.moveTo(x + w, y + h - r);
	ctx.arcTo(x + w, y + h, x + w - r, y + h, r);
	ctx.lineTo(x + w, y + h);
	ctx.lineTo(x + w, y + h - r);
	ctx.fill();
};
