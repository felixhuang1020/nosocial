import { createStoreBindings } from 'mobx-miniprogram-bindings';
import { appStore } from '../../stores/app';
import { get, post } from '../../utils/request';

Page({
  data: {
    registerFee: '99.00',
    paying: false
  },

  storeBindings: null,

  onLoad() {
    this.storeBindings = createStoreBindings(this, {
      store: appStore,
      fields: ['isLogin', 'isShareholder'],
      actions: []
    });
    this.fetchConfig();
  },

  // 获取业务配置
  fetchConfig() {
    get('/public/config').then(res => {
      if (res && res.shareholder_fee) {
        this.setData({
          registerFee: res.shareholder_fee.toFixed(2)
        });
      }
    }).catch(err => {
      console.log('使用默认注册费配置');
    });
  },

  onShow() {
    if (this.storeBindings) {
      this.storeBindings.updateStoreBindings();
    }
    // 检查是否已经是股东
    if (this.data.isShareholder) {
      wx.redirectTo({ url: '/pages/shareholder/earnings' });
      return;
    }
  },

  onUnload() {
    if (this.storeBindings) {
      this.storeBindings.destroyStoreBindings();
    }
  },

  // 申请成为股东
  applyShareholder() {
    if (!this.data.isLogin) {
      wx.showModal({
        title: '提示',
        content: '请先登录后再申请',
        confirmText: '去登录',
        confirmColor: '#7C9A92',
        success: (res) => {
          if (res.confirm) {
            wx.switchTab({ url: '/pages/profile/profile' });
          }
        }
      });
      return;
    }

    wx.showModal({
      title: '确认支付',
      content: `确认支付 ¥${this.data.registerFee} 成为共享股东？`,
      confirmColor: '#7C9A92',
      success: (res) => {
        if (res.confirm) {
          this.processPayment();
        }
      }
    });
  },

  // 处理支付
  processPayment() {
    this.setData({ paying: true });

    post('/wx/shareholder/pay', {
      amount: parseFloat(this.data.registerFee)
    }).then(res => {
      // 调起微信支付
      wx.requestPayment({
        ...res.pay_params,
        success: (payRes) => {
          wx.showToast({ title: '支付成功', icon: 'success' });
          // 延迟后跳转到股东中心
          setTimeout(() => {
            wx.redirectTo({ url: '/pages/shareholder/earnings' });
          }, 1500);
        },
        fail: (err) => {
          if (err.errMsg && err.errMsg.includes('cancel')) {
            wx.showToast({ title: '支付已取消', icon: 'none' });
          } else {
            wx.showToast({ title: '支付失败', icon: 'none' });
          }
        },
        complete: () => {
          this.setData({ paying: false });
        }
      });
    }).catch(err => {
      wx.showToast({ title: '发起支付失败', icon: 'none' });
      this.setData({ paying: false });
    });
  }
});
