<template>
	<view
		class="container"
		v-if="is_load_over"
		:style="{
			'--status_bar_height': system_info.statusBarHeight + 'px',
		}"
	>
		<!-- <cu-custom bgImage="https://comic-img.zwww.cool/images/banner.png"></cu-custom> -->
		<!-- 		<image class="logo" src="https://comic-img.zwww.cool/images/banner.png" mode="aspectFit"></image>

		 -->

		<swiper
			v-if="show_etcs"
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
					<view class="etcs-imgs" :class="item.direction == IMG_DIRECTION_ROW ? 'h-50' : 'w-50'">
						<image @tap="previewImgs(0)" :src="item.origin" mode="aspectFill"></image>
						<image @tap="previewImgs(0)" :src="item.res" mode="aspectFill"></image>
					</view>
					<view class="etcs-desc" v-if="item.desc">{{ item.desc }}</view>
				</view>
			</swiper-item>
		</swiper>

		<view class="transfer" @tap="upload"><view class="btn cu-btn bg-yellow lg shadow">开始</view></view>

		<view class="preview" v-if="has_transfer">
			<view class="preview-imgs" :class="img_direction == IMG_DIRECTION_ROW ? 'h-50' : 'w-50'">
				<image @tap="previewImgs(0)" :src="img_origin" mode="aspectFill"></image>
				<image @tap="previewImgs(1)" v-if="img_result" :src="img_result" mode="aspectFill"></image>
			</view>

			<view class="box-btns" v-if="img_result">
				<view class="cu-btn bg-yellow shadow radius" @tap="save">保存</view>
				<view class="cu-btn bg-yellow shadow radius" @tap="share">分享</view>
			</view>
		</view>
	</view>
</template>

<script>
import cuCustom from '@/static/colorui/components/cu-custom.vue';

import Upload from '@/common/request/upload.js';
let upload = new Upload();

import config from '@/common/conifg.js';
import { downloadFile, getImgInfo } from '@/common/helper/utils.js';

const IMG_OUT_URL = config.img_prefix;

const IMG_DIRECTION_ROW = 1;
const IMG_DIRECTION_COLUMN = 2;

export default {
	components: {
		cuCustom,
	},
	data() {
		return {
			IMG_DIRECTION_ROW,
			IMG_DIRECTION_COLUMN,

			auth: {},
			system_info: getApp().globalData.system_info,

			etcs: [],
			etc_directions: [],

			image_support: ['png', 'jpg'],
			image_size: 10,

			img_origin: '', //转换前
			img_derection: 1, //转换前
			img_result: '',

			modal_show: true,
			transfer_type: 1,

			is_transfering: false,
			has_transfer: false,

			cardCur: 0,
		};
	},

	onLoad() {
		this.checkLogin();
		this.loadEtcs();
	},
	computed: {
		is_load_over() {
			return this.show_etcs;
		},
		show_etcs() {
			return this.etc_directions.length == this.etcs.length;
		},
	},
	methods: {
		async loadEtcs() {
			let info = {};
			this.$go.to('old_etc').then(res => {
				console.log(res);
				this.etcs = res.data;
				this.etcs.forEach(etc => {
					(async () => {
						info = await getImgInfo(etc.res);
						Object.assign(etc, {
							direction: info.width > info.height ? IMG_DIRECTION_ROW : IMG_DIRECTION_COLUMN,
						});
						this.etc_directions.push(1);
					})();
				});
			});
		},

		checkLogin() {
			let auth = this.$cache.get('auth');
			if (auth) {
				this.auth = auth;
				this.$store.commit('login/login', auth);
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
					});
			});
		},

		async upload() {
			if (this.is_transfering) {
				this.$base.showToast('拼命绘制中...');
				return;
			}
			wx.chooseImage({
				count: 1,
				// sizeType:['compressed'],
				success: res => {
					this.has_transfer = true;
					setTimeout(() => {
						uni.pageScrollTo({
							scrollTop: this.system_info.screenHeight,
							duration: 1000,
						});
					}, 100);
					this.$base.showLoading('上色中...');
					this.is_transfering = true;

					let file = res.tempFiles[0];
					let file_name = file.name || file.path;
					
					let self = this
					uni.getImageInfo({
						src:file_name,
						success(info) {
							self.img_derection = info.width > info.height ? IMG_DIRECTION_ROW : IMG_DIRECTION_COLUMN;
							console.log(res);
						}
					})

					this.img_origin = file_name;
					this.img_result = '';
					let file_type = file_name.substr(file_name.lastIndexOf('.') + 1).toLowerCase();
					if (this.image_support.includes(file_type) && file.size > 1024 * 1024 * this.image_size) {
						this.$base.showToast(`图片大小请控制在${this.image_size}兆以内`);
						has_error = true;
						return true;
					}
					let head = {
						Authorization: 'Bearer ' + this.auth.token,
						transfer_type: this.transfer_type,
					};
					upload.uploadImg(file_name, head)
						.then(res => {
							this.img_result = IMG_OUT_URL + res.data.filename;
							this.is_transfering = false;
							uni.hideLoading();
							this.$base.showToast('上色成功!');
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
	},
};
</script>

<style lang="scss">
page {
	color: #ffffff;
}

.container {
	min-height: 100vh;
	background-image: linear-gradient(15deg, #dddddd, #333333);
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
	width: 100%;
}
.h-50 image:first-child {
	margin-bottom: 10rpx;
}
.w-50 image {
	max-height: 100%;
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
		font-size: 26rpx;
		color: #ffffff;
		padding: 40rpx 80rpx;
		word-break: break-all;
		white-space: pre-line;
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
		// flex-wrap: wrap;
		align-items: center;
		justify-content: center;
	}

	image {
		// width: 100%;
		display: block;
		// height: 50%;
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
	&-imgs {
		display: flex;
		justify-content: space-around;
		align-items: center;
		image {
			width: 280rpx;
			max-height: 800rpx;
			border-radius: 10rpx;
		}
	}
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
	color: #333;

	&-next {
		position: fixed;
		right: 4%;
		bottom: 4%;
	}
}
</style>
