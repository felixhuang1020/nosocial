import { createStoreBindings } from 'mobx-miniprogram-bindings';
import { appStore } from '../../stores/app';
import { get, post, put } from '../../utils/request';

Page({
  data: {
    loginLoading: false,
    teamCount: 0,
    showBirthdayPicker: false,
    pickerValue: '',
    currentDate: ''
  },

  storeBindings: null,

  onLoad() {
    this.storeBindings = createStoreBindings(this, {
      store: appStore,
      fields: ['userInfo', 'isLogin', 'isShareholder', 'nickname', 'avatar', 'birthday', 'inviteCode', 'balance', 'totalEarning'],
      actions: ['setToken', 'setUserInfo', 'clearUserInfo']
    });
  },

  onShow() {
    if (this.storeBindings) {
      this.storeBindings.updateStoreBindings();
    }
    // 计算当前日期作为选择器上限
    const now = new Date();
    const y = now.getFullYear();
    const m = String(now.getMonth() + 1).padStart(2, '0');
    const d = String(now.getDate()).padStart(2, '0');
    this.setData({ currentDate: `${y}-${m}-${d}` });
    
    if (this.data.isLogin) {
      this.fetchUserProfile();
      if (this.data.isShareholder) {
        this.fetchTeamCount();
      }
    }
  },

  onReady() {
    // 页面首次渲染完成
  },

  onUnload() {
    if (this.storeBindings) {
      this.storeBindings.destroyStoreBindings();
    }
  },

  // 微信登录
  handleLogin() {
    this.setData({ loginLoading: true });
    
    const app = getApp();
    const inviteCode = wx.getStorageSync('invite_code') || '';
    
    app.wxLogin(inviteCode).then(data => {
      wx.showToast({ title: '登录成功', icon: 'success' });
      this.fetchUserProfile();
    }).catch(err => {
      wx.showToast({ title: '登录失败', icon: 'none' });
    }).finally(() => {
      this.setData({ loginLoading: false });
    });
  },

  // 获取用户资料
  fetchUserProfile() {
    get('/wx/user').then(res => {
      this.setUserInfo(res);
      if (res.is_shareholder) {
        this.fetchTeamCount();
      }
    }).catch(err => {
      console.error('获取用户信息失败:', err);
    });
  },

  // 获取团队人数
  fetchTeamCount() {
    get('/wx/shareholder/team').then(res => {
      const count = Array.isArray(res) ? res.length : 0;
      this.setData({ teamCount: count });
    }).catch(err => {
      console.log('获取团队人数失败');
    });
  },

  // 导航到收益明细
  navigateToEarnings() {
    wx.navigateTo({ url: '/pages/shareholder/earnings' });
  },

  // 导航到团队
  navigateToTeam() {
    wx.navigateTo({ url: '/pages/shareholder/team' });
  },

  // 导航到申请股东
  navigateToApply() {
    wx.navigateTo({ url: '/pages/shareholder/apply' });
  },

  // 导航到订单列表
  navigateToOrders() {
    wx.navigateTo({ url: '/pages/order/list' });
  },

  // 导航到优惠券
  navigateToCoupons() {
    wx.navigateTo({ url: '/pages/coupon/list' });
  },

  // 导航到占卜历史
  navigateToTarotHistory() {
    wx.navigateTo({ url: '/pages/tarot/history' });
  },

  // 分享邀请
  shareInvite() {
    wx.showShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage']
    });
  },

  // 设置生日 - 打开日期选择器弹窗
  setBirthday() {
    const now = new Date();
    const y = now.getFullYear();
    this.setData({
      pickerValue: this.data.birthday || `${y - 25}-01-01`,
      showBirthdayPicker: true
    });
  },

  // 日期选择器确认
  onBirthdayConfirm(e) {
    const birthday = e.detail.value;
    this.setData({ showBirthdayPicker: false });
    this.updateBirthday(birthday);
  },

  // 日期选择器取消
  onBirthdayCancel() {
    this.setData({ showBirthdayPicker: false });
  },

  // 日期选择器显隐变化
  onBirthdayPickerVisible(e) {
    this.setData({ showBirthdayPicker: e.detail.visible });
  },

  // 更新生日
  updateBirthday(birthday) {
    put('/wx/user/birthday', { birthday }).then(res => {
      appStore.updateUserField('birthday', birthday);
      wx.showToast({ title: '生日设置成功', icon: 'success' });
      this.fetchUserProfile();
    }).catch(err => {
      wx.showToast({ title: '设置失败', icon: 'none' });
    });
  },

  // 联系客服（corpId 需配置企业微信客服ID，当前回退到显示固定电话）
  contactService() {
    wx.openCustomerServiceChat({
      extInfo: { url: '' },
      corpId: '',
      success: () => {},
      fail: () => {
        wx.showModal({
          title: '联系客服',
          content: '客服电话: 010-8888-6666\n工作时间: 18:00-04:00',
          showCancel: false,
          confirmColor: '#7C9A92'
        });
      }
    });
  },

  // 关于我们
  aboutUs() {
    wx.showModal({
      title: '关于 NoSocial',
      content: 'NoSocial 酒吧小程序\n版本: 1.0.0\n\n没有社交，只有美酒。',
      showCancel: false,
      confirmColor: '#7C9A92'
    });
  },

  // 退出登录
  logout() {
    wx.showModal({
      title: '确认退出',
      content: '确定要退出登录吗？',
      confirmColor: '#ff4d4f',
      success: (res) => {
        if (res.confirm) {
          this.clearUserInfo();
          wx.showToast({ title: '已退出登录', icon: 'success' });
        }
      }
    });
  },

  // 分享
  onShareAppMessage() {
    const inviteCode = this.data.inviteCode || '';
    return {
      title: 'NoSocial 酒吧 - 没有社交，只有美酒',
      path: `/pages/index/index?invite_code=${inviteCode}`,
      imageUrl: ''
    };
  }
});
