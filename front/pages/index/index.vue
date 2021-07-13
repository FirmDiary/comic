<template>
	<view class="container">
		<image class="logo" src="/static/images/logo.png" mode="aspectFit"></image>
		<view class="transfer" @tap="showModal"><view class="btn cu-btn bg-yellow lg shadow">开始</view></view>

		<view class="box" v-if="img_origin">
			<view class="box-imgs">
				<view><image @tap="previewImgs(0)" :src="img_origin" mode="aspectFill"></image></view>
				<!-- <view>==</view> -->
				<view>
					<image
						v-if="!img_result"
						style="width: 120rpx;"
						src="/static/images/loading_transfer.gif"
						mode="aspectFit"
					></image>
					<image :src="img_result" @tap="previewImgs(1)" v-else mode="aspectFill"></image>
				</view>
			</view>

			<view class="box-btns" v-if="img_result">
				<view class="cu-btn bg-brown shadow radius" @tap="save">保存</view>
			</view>
		</view>

		<view class="cu-modal select_transfer" :class="modal_show ? 'show' : ''" @tap="hideModal">
			<view class="cu-dialog" @tap.stop="">
				<scroll-view scroll-x class="bg-black nav">
					<view class="flex text-center">
						<view
							class="cu-item flex-sub"
							:class="item.value == transfer_type ? 'text-yellow' : ''"
							v-for="(item, index) in TRANSFER_TYPE"
							:key="index"
							@tap="tabSelect(item.value)"
						>
							{{ item.label }}
						</view>
					</view>
				</scroll-view>

				<view class="etcs">
					<swiper
						class="screen-swiper"
						style="min-height: 900rpx;"
						:class="dotStyle ? 'square-dot' : 'round-dot'"
						:indicator-dots="true"
						:circular="true"
						:autoplay="true"
						interval="5000"
						duration="500"
					>
						<swiper-item v-for="(item, index) in etc_imgs" class="bg-black" :key="index">
							<image :src="item" mode="aspectFill"></image>
						</swiper-item>
					</swiper>
				</view>
			</view>

			<button @tap="upload" class="btn cu-btn bg-yellow lg shadow select_transfer-next">
				<text class="cuIcon-upload"></text>
				上传
			</button>
		</view>
	</view>
</template>

<script>
import Upload from '@/common/request/upload.js';
let upload = new Upload();

import { downloadFile } from '@/common/helper/utils.js';

const TRANSFER_TYPE = [
	{
		label: '脸部',
		value: 1,
	},
	{
		label: '其他',
		value: 2,
	},
];

const ETCS = {
	1: [
		'/static/images/face1.png',
		'/static/images/face2.png',
		'/static/images/face3.png',
		'/static/images/face4.png',
		'/static/images/face5.png',
		'/static/images/face6.png',
	],
	2: ['/static/images/face4.png', '/static/images/face5.png', '/static/images/face6.png'],
};

const IMG_OUT_URL = 'https://comic-img.zwww.cool/out/';

export default {
	data() {
		return {
			TRANSFER_TYPE,
			ETCS,

			auth: {},

			image_support: ['png', 'jpg'],
			image_size: 10,

			img_origin: '', //转换前
			img_result: '',

			modal_show: true,
			transfer_type: 1,
		};
	},
	onLoad() {
		this.checkLogin();
	},
	computed: {
		// 支付宝小程序需写成计算属性,prop定义default仍报错
		etc_imgs() {
			return this.ETCS[this.transfer_type];
		},
	},
	methods: {
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

		upload() {
			wx.chooseImage({
				count: 1,
				success: res => {
					this.$base.showLoading("拼命绘制中...");
					let file = res.tempFiles[0];
					let file_name = file.name || file.path;
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
					upload.uploadImg(file_name, head).then(res => {
						this.img_result = IMG_OUT_URL + res.data.path;
						uni.hideLoading();
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

		showModal(e) {
			this.modal_show = true;
		},
		hideModal(e) {
			this.modal_show = false;
		},
		tabSelect(value) {
			this.transfer_type = value;
		},
	},
};
</script>

<style lang="scss">
page {
	background-image: linear-gradient(15deg, #333333, #fff);
	color: #ffffff;
}

.container {
	// min-height: 100vh;
	// position: relative;
}

.logo {
	display: block;
	margin: 80rpx auto 0;
}

.transfer {
	margin: 0rpx auto 120rpx;
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
		flex-direction: column;
		align-items: center;
		margin-top: 140rpx;
		.cu-btn {
			margin-bottom: 30rpx;
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
