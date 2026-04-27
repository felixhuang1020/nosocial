import { createStoreBindings } from 'mobx-miniprogram-bindings';
import { appStore } from '../../stores/app';
import { get } from '../../utils/request';

Page({
  data: {
    // 轮播图数据
    banners: [],
    swiperImageProps: { mode: 'aspectFill', src: 'image' },
    
    // 推荐酒水
    recommendDrink: null,
    
    // 店铺信息
    shopInfo: {
      address: '北京市朝阳区三里屯太古里北区 N8-20',
      phone: '010-8888-6666',
      business_hours: '周一至周日 18:00 - 04:00',
      wifi: 'NoSocial_Free / nosocial888',
      latitude: 39.934,
      longitude: 116.455
    }
  },

  storeBindings: null,

  onLoad(options) {
    this.storeBindings = createStoreBindings(this, {
      store: appStore,
      fields: ['userInfo', 'isShareholder', 'isLogin'],
      actions: []
    });
    
    // 获取邀请码参数
    if (options.invite_code) {
      wx.setStorageSync('invite_code', options.invite_code);
    }
    
    this.fetchBanners();
    this.fetchRecommendDrink();
    this.fetchShopInfo();
  },

  onShow() {
    if (this.storeBindings) {
      this.storeBindings.updateStoreBindings();
    }
  },

  onReady() {
    // 页面首次渲染完成
  },

  onHide() {
    // 页面隐藏
  },

  onUnload() {
    if (this.storeBindings) {
      this.storeBindings.destroyStoreBindings();
    }
  },

  // 下拉刷新
  onPullDownRefresh() {
    Promise.all([
      this.fetchBanners(),
      this.fetchRecommendDrink()
    ]).finally(() => {
      wx.stopPullDownRefresh();
    });
  },

  // 获取轮播图
  fetchBanners() {
    return get('/public/banners', { position: 1 }).then(res => {
      if (res && res.length > 0) {
        this.setData({
          banners: res.map(item => ({
            image: item.image_url,
            link_type: item.link_type,
            link_value: item.link_value,
            id: item.id
          }))
        });
      }
    }).catch(err => {
      console.log('使用默认轮播图数据');
    });
  },

  // 获取店铺信息
  fetchShopInfo() {
    return get('/public/shop').then(res => {
      if (res) {
        this.setData({
          shopInfo: {
            address: res.address || this.data.shopInfo.address,
            phone: res.phone || this.data.shopInfo.phone,
            business_hours: res.business_hours || this.data.shopInfo.business_hours,
            wifi: `${res.wifi_name || 'NoSocial_Free'} / ${res.wifi_password || 'nosocial888'}`,
            latitude: parseFloat(res.latitude) || this.data.shopInfo.latitude,
            longitude: parseFloat(res.longitude) || this.data.shopInfo.longitude
          }
        });
      }
    }).catch(err => {
      console.log('使用默认店铺信息');
    });
  },

  // 获取推荐酒水
  fetchRecommendDrink() {
    return get('/public/drinks', { page: 1, size: 20 }).then(res => {
      const list = res && res.list ? res.list : [];
      const recommended = list.find(item => item.is_recommended);
      if (recommended) {
        this.setData({ recommendDrink: recommended });
      }
    }).catch(err => {
      console.log('获取推荐酒水失败');
    });
  },

  // 轮播图点击
  onBannerClick(e) {
    const { index } = e.detail;
    const banner = this.data.banners[index];
    if (!banner) return;
    
    if (banner.link_type === 1 && banner.link_value) {
      // 内部页面跳转
      wx.navigateTo({ url: banner.link_value });
    } else if (banner.link_type === 2 && banner.link_value) {
      // 外部链接 - 使用 web-view
      wx.navigateTo({ 
        url: `/pages/webview/webview?url=${encodeURIComponent(banner.link_value)}` 
      });
    }
  },

  // 导航到共享股东
  navigateToShareholder() {
    if (!this.data.isLogin) {
      this.showLoginTip();
      return;
    }
    
    if (this.data.isShareholder) {
      wx.navigateTo({ url: '/pages/shareholder/earnings' });
    } else {
      wx.navigateTo({ url: '/pages/shareholder/apply' });
    }
  },

  // 导航到生日有礼
  navigateToBirthday() {
    if (!this.data.isLogin) {
      this.showLoginTip();
      return;
    }
    
    if (!this.data.isShareholder) {
      wx.showToast({ 
        title: '仅股东可领取生日礼', 
        icon: 'none',
        duration: 2000
      });
      return;
    }
    
    wx.navigateTo({ url: '/pages/birthday/gift' });
  },

  // 导航到塔罗占卜
  navigateToTarot() {
    wx.navigateTo({ url: '/pages/tarot/divine' });
  },

  // 导航到离店点评
  navigateToReview() {
    if (!this.data.isLogin) {
      this.showLoginTip();
      return;
    }
    wx.navigateTo({ url: '/pages/review/guide' });
  },

  // 导航到酒水详情
  navigateToDrinkDetail(e) {
    const id = e.currentTarget.dataset.id;
    wx.navigateTo({ url: `/pages/drink/detail?id=${id}` });
  },

  // 导航到酒水单
  navigateToMenu() {
    wx.switchTab({ url: '/pages/menu/menu' });
  },

  // 打开地图导航
  openAddress() {
    const info = this.data.shopInfo;
    wx.openLocation({
      latitude: info.latitude || 39.934,
      longitude: info.longitude || 116.455,
      name: 'NoSocial 酒吧',
      address: info.address,
      scale: 18
    });
  },

  // 拨打电话
  makePhoneCall() {
    wx.makePhoneCall({
      phoneNumber: this.data.shopInfo.phone.replace(/-/g, '')
    });
  },

  // 复制WiFi密码
  copyWifi() {
    const wifiStr = this.data.shopInfo.wifi;
    wx.setClipboardData({
      data: wifiStr,
      success: () => {
        wx.showToast({ title: 'WiFi信息已复制', icon: 'success' });
      }
    });
  },

  // 显示登录提示
  showLoginTip() {
    wx.showModal({
      title: '提示',
      content: '请先登录后再使用该功能',
      confirmText: '去登录',
      confirmColor: '#7C9A92',
      success: (res) => {
        if (res.confirm) {
          wx.switchTab({ url: '/pages/profile/profile' });
        }
      }
    });
  }
});
