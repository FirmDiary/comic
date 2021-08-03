<template>
	<view
		class="container"
		v-if="is_load_over"
		:style="{
			'--status_bar_height': system_info.statusBarHeight + 'px',
		}"
	>
		<view class="quota">额度:{{ user.quota }}</view>
		<swiper
			v-if="etcs.length"
			class="screen-swiper etcs"
			:class="dotStyle ? 'square-dot' : 'round-dot'"
			:indicator-dots="true"
			:circular="true"
			:autoplay="true"
			interval="5000"
			duration="500"
			@change="cardSwiper"
			indicator-color="#8799a3"
			indicator-active-color="#fbbd08"
		>
			<swiper-item v-for="(item, index) in etcs" :key="index" :class="cardCur == index ? 'cur' : ''">
				<view class="swiper-item">
					<view class="etcs-imgs" :class="item.direction == IMG_DIRECTION_COLUMN ? 'h-50' : 'w-50'">
						<image @tap="previewImgs(0)" :src="item.origin" mode="aspectFill"></image>
						<image @tap="previewImgs(0)" :src="item.res" mode="aspectFill"></image>
					</view>
					<view class="etcs-desc" v-if="item.desc">{{ item.desc }}</view>
				</view>
			</swiper-item>
		</swiper>

		<view class="transfer" @tap="upload"><view class="btn cu-btn bg-yellow lg shadow">开始</view></view>

		<view class="preview" v-if="has_transfer">
			<view class="preview-imgs" :class="img_direction == IMG_DIRECTION_COLUMN ? 'h-50' : 'w-50'">
				<image @tap="previewImgs(0)" :src="img_origin" mode="aspectFill"></image>
				<image @tap="previewImgs(1)" v-if="img_result" :src="img_result" mode="aspectFill"></image>
			</view>

			<view class="box-btns" v-if="img_result">
				<view class="cu-btn bg-yellow shadow radius" @tap="save">保存</view>
				<button open-type="share"><view class="cu-btn bg-yellow shadow radius" @tap="share">分享</view></button>
			</view>
		</view>

		<view class="cu-modal invite" :class="{ show: showTip }">
			<view class="cu-dialog">
				<view class="cu-bar bg-yellow justify-end">
					<view class="content">额度不足</view>
					<view class="action" @tap="hideTip"><text class="cuIcon-close text-red"></text></view>
				</view>
				<view class="padding-xl invite-desc">
					<p>
						每邀请一名新用户可增加
						<span>2</span>
						次额度
					</p>
					<p>新用户仅需点开小程序即可完成邀请</p>
					<p>邀请成功后下拉刷新额度</p>
				</view>
				<button open-type="share" class="invite-share">
					<view class="cu-btn bg-yellow shadow radius lg" @tap="share">邀请</view>
				</button>
			</view>
		</view>

		<view class="copyright">v1.1 | by deepai.org</view>
	</view>
</template>

<script>
import cuCustom from '@/static/colorui/components/cu-custom.vue';

import Upload from '@/common/request/upload.js';
let upload = new Upload();

import config from '@/common/conifg.js';
import { downloadFile, getImgInfo } from '@/common/helper/utils.js';

const IMG_OUT_URL = config.img_prefix;

const IMG_DIRECTION_ROW = 'row';
const IMG_DIRECTION_COLUMN = 'column';

export default {
	components: {
		cuCustom,
	},
	data() {
		return {
			IMG_DIRECTION_ROW,
			IMG_DIRECTION_COLUMN,

			user: {},
			auth: {},
			system_info: getApp().globalData.system_info,

			invite_id: 0,

			etcs: [],
			etc_directions: [],

			image_support: ['png', 'jpg'],
			image_size: 10,

			img_origin: '', //转换前
			img_direction: IMG_DIRECTION_COLUMN, //转换前
			img_result: '',

			is_transfering: false,
			has_transfer: false,

			cardCur: 0,

			showTip: false,
		};
	},

	onLoad(options) {
		this.invite_id = options.invite_id || 0;
		this.checkLogin();
		this.loadEtcs();
	},
	computed: {
		is_load_over() {
			return this.etcs.length;
		},
	},
	onPullDownRefresh() {
		this.$base.showLoading();
		this.getUser();
	},
	onShareAppMessage(res) {
		if (res.from === 'button') {
			// console.log(res.target);
		}
		return {
			title: '使用人工智能将你变成卡通!',
			imageUrl: this.img_result || '',
			path: `/pages/index/index?invite_id=${this.user.id}`,
		};
	},
	onShareTimeline() {
		return {
			title: '使用人工智能将你变成卡通!',
			imageUrl: this.img_result || '',
			path: `/pages/index/index?invite_id=${this.user.id}`,
		};
	},
	methods: {
		loadEtcs() {
			this.$go.to('old_etc').then(res => {
				this.etcs = res.data;
			});
		},
		getUser() {
			this.$go.to('user').then(res => {
				this.user = res.data;
				uni.hideLoading();
				uni.stopPullDownRefresh();
			});
		},
		inviteSuccess() {
			if (this.invite_id) {
				//受邀请
				this.$go.to('invite_success', { invite_id: this.invite_id }).then(res => {
					console.log('受邀请:' + this.invite_id);
					console.log(res);
				});
			}
		},
		checkLogin() {
			let auth = this.$cache.get('auth');
			if (auth) {
				this.auth = auth;
				this.$store.commit('login/login', auth);
				this.getUser();
				return;
			}
			wx.login().then(res => {
				this.$go
					.to('login', {
						code: res.code,
					})
					.then(res => {
						let auth = {
							token: res.data.token,
						};
						this.auth = auth;
						this.$cache.set('auth', auth);
						this.$store.commit('login/login', auth);
						this.getUser();
						this.inviteSuccess();
					});
			});
		},

		async upload() {
			if (this.is_transfering) {
				this.$base.showToast('拼命绘制中...');
				return;
			}
			if (this.user.quota == 0) {
				//额度不足
				this.showTip = true;
				return;
			}
			wx.chooseImage({
				count: 1,
				// sizeType:['compressed'],
				success: res => {
					let file = res.tempFiles[0];
					let file_name = file.name || file.path;

					this.img_origin = file_name;
					this.img_result = '';
					let file_type = file_name.substr(file_name.lastIndexOf('.') + 1).toLowerCase();
					if (!this.image_support.includes(file_type)) {
						this.$base.showToast(`只支持${this.image_support.join("、")}格式的图片`);
						return true;
					}
					if (file.size > 1024 * 1024 * this.image_size) {
						this.$base.showToast(`图片大小请控制在${this.image_size}兆以内`);
						return true;
					}
					
					this.has_transfer = true;
					setTimeout(() => {
						uni.pageScrollTo({
							scrollTop: this.system_info.screenHeight,
							duration: 1000,
						});
					}, 100);
					this.$base.showLoading('绘制中...');
					this.is_transfering = true;
					
					let head = {
						Authorization: 'Bearer ' + this.auth.token,
					};
					upload
						.uploadImg(file_name, head)
						.then(res => {
							uni.hideLoading();
							this.is_transfering = false;
							if (res.quota == -1) {
								//额度不足
								this.showTip = true;
								return;
							}
							this.img_result = IMG_OUT_URL + res.data.filename;
							this.img_direction = res.data.direction;
							this.user.quota--;
							this.$base.showToast('绘制成功!');
						})
						.catch(err => {
							console.error(err);
							this.is_transfering = false;
							this.$base.showToast('绘制失败~换张图片试试吧');
						});
				},
			});
		},

		save() {
			downloadFile(this.img_result).then(filePath => {
				uni.saveImageToPhotosAlbum({
					filePath,
				});
			});
		},

		previewImgs(index) {
			let images = [this.img_origin];
			if (this.img_result) {
				images.push(this.img_result);
			}
			console.log(images);
			uni.previewImage({
				current: index,
				urls: images,
			});
		},

		cardSwiper(e) {
			this.cardCur = e.detail.current;
		},

		hideTip() {
			this.showTip = false;
		},
	},
};
</script>

<style lang="scss">
$light: #fbbd08;
$dark: #333;
$common: #ddd;

page {
	color: $dark;
}

.container {
	min-height: 100vh;
	background-image: linear-gradient(45deg, #1cbbb4, #0081ff);
	padding-bottom: 80rpx;
}

.logo {
	display: block;
	width: 100%;
	// margin: 80rpx auto 0;
}

.swiper-item {
	width: 100%;
	height: 100%;
}

.swiper-item image {
	display: block;
	margin: 0;
	pointer-events: none;
}
.h-50 image {
	height: 50% !important;
	border-radius: 10rpx;
	width: 100% !important;
}
.h-50 image:first-child {
	margin-bottom: 10rpx;
}
.w-50 {
	flex-wrap: inherit !important;
}
.w-50 image {
	height: 50% !important;
	border-radius: 10rpx;
	width: 50% !important;
}
.w-50 image:first-child {
	margin-right: 10rpx;
}

.etcs {
	height: 88vh;
	padding: 0rpx 20rpx;
	padding-top: calc(var(--status_bar_height) + 100rpx);

	&-imgs {
		width: 100%;
		height: 80%;
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: center;
	}

	&-desc {
		word-break: break-word;
		white-space: pre-line;
		font-size: 26rpx;
		color: #ffffff;
		padding: 40rpx 80rpx;
		line-height: 40rpx;
		text-align: center;
	}
}

.preview {
	padding: 0 20rpx;
	width: 100%;
	height: 100vh;
	padding-top: calc(var(--status_bar_height) + 180rpx);

	&-imgs {
		height: 89%;
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: center;
	}

	image {
		display: block;
		height: 100%;
	}
}

.transfer {
	display: flex;
	justify-content: center;
	.btn {
		color: #333;
		width: 388rpx;
		border-radius: 10rpx;
		height: 90rpx;
		box-shadow: 0rpx 6rpx 20rpx 6rpx rgba(224, 170, 7, 0.5);
	}
}

.box {
	padding: 0 50rpx;
	&-btns {
		display: flex;
		flex-direction: row;
		align-items: center;
		justify-content: space-evenly;
		width: 100%;
		margin-top: 46rpx;
		margin-bottom: 56rpx;
		.cu-btn {
			width: 168rpx;
		}
	}
}

.select_transfer {
	color: $dark;

	&-next {
		position: fixed;
		right: 4%;
		bottom: 4%;
	}
}

.copyright {
	position: fixed;
	text-align: center;
	width: 100%;
	bottom: 10rpx;
	font-size: 20rpx;
}

.quota {
	color: #fff;
	position: fixed;
	text-align: center;
	font-size: 24rpx;
	top: calc(var(--status_bar_height) + 26rpx);
	left: 26rpx;
}
.invite {
	.cuIcon-close {
		color: #fff;
	}
	&-desc {
		color: $dark;
		font-size: 28rpx;
		line-height: 50rpx;
		span {
			color: $light;
			font-size: 32rpx;
		}
	}
	&-share {
		margin: 0rpx 0 40rpx;
		background-color: #f8f8f8;
	}

	.cu-btn {
		width: 272rpx;
		font-size: 30rpx;
		border-radius: 10rpx;
	}
}
</style>
