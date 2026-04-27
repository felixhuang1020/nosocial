import { createStoreBindings } from 'mobx-miniprogram-bindings';
import { appStore } from '../../stores/app';
import { get, post } from '../../utils/request';

Page({
  data: {
    totalEarning: '0.00',
    balance: '0.00',
    inviteCode: '',
    earningsList: [],
    loading: false
  },

  storeBindings: null,

  onLoad() {
    this.storeBindings = createStoreBindings(this, {
      store: appStore,
      fields: ['totalEarning', 'balance', 'inviteCode'],
      actions: []
    });

    this.fetchEarnings();
  },

  onShow() {
    if (this.storeBindings) {
      this.storeBindings.updateStoreBindings();
    }
  },

  onUnload() {
    if (this.storeBindings) {
      this.storeBindings.destroyStoreBindings();
    }
  },

  // 获取收益明细
  fetchEarnings() {
    this.setData({ loading: true });
    
    get('/wx/shareholder/earnings').then(res => {
      const list = res && res.list ? res.list : [];
      this.setData({
        earningsList: list,
        loading: false
      });
    }).catch(err => {
      console.error('获取收益明细失败:', err);
      this.setData({
        earningsList: [],
        loading: false
      });
    });
  },

  // 提现
  withdraw() {
    if (parseFloat(this.data.balance) < 50) {
      wx.showToast({ title: '满50元才可提现', icon: 'none' });
      return;
    }

    wx.showModal({
      title: '申请提现',
      content: `确认提现 ¥${this.data.balance} 到微信零钱？`,
      confirmColor: '#7C9A92',
      success: (res) => {
        if (res.confirm) {
          post('/wx/shareholder/withdraw', {
            amount: parseFloat(this.data.balance)
          }).then(() => {
            wx.showToast({ title: '提现申请已提交', icon: 'success' });
            // 重新获取用户信息以刷新余额
            get('/wx/user').then(res => {
              appStore.setUserInfo(res);
            });
          }).catch(err => {
            wx.showToast({ title: err.message || '提现失败', icon: 'none' });
          });
        }
      }
    });
  },

  // 分享邀请
  shareInvite() {
    wx.showShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage']
    });
  },

  // 复制邀请码
  copyInviteCode() {
    wx.setClipboardData({
      data: this.data.inviteCode,
      success: () => {
        wx.showToast({ title: '邀请码已复制', icon: 'success' });
      }
    });
  },
});
