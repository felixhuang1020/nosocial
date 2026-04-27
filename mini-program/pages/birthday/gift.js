import { createStoreBindings } from 'mobx-miniprogram-bindings';
import { appStore } from '../../stores/app';
import { get, post, put } from '../../utils/request';

Page({
  data: {
    isBirthday: false,
    birthday: '',
    daysUntil: 0,
    giftInfo: null,
    giftClaimed: false,
    claiming: false,
    giftHistory: []
  },

  storeBindings: null,

  onLoad() {
    this.storeBindings = createStoreBindings(this, {
      store: appStore,
      fields: ['birthday'],
      actions: []
    });
  },

  onShow() {
    if (this.storeBindings) {
      this.storeBindings.updateStoreBindings();
    }
    this.checkBirthday();
    this.fetchGiftInfo();
    this.fetchGiftHistory();
  },

  onUnload() {
    if (this.storeBindings) {
      this.storeBindings.destroyStoreBindings();
    }
  },

  // 检查是否今天生日
  checkBirthday() {
    const birthday = this.data.birthday;
    if (!birthday) {
      this.setData({ isBirthday: false, birthday: '' });
      return;
    }

    const today = new Date();
    const birthDate = new Date(birthday);
    
    const isToday = today.getMonth() === birthDate.getMonth() && 
                    today.getDate() === birthDate.getDate();
    
    // 计算距离下一个生日的天数
    const nextBirthday = new Date(today.getFullYear(), birthDate.getMonth(), birthDate.getDate());
    if (nextBirthday < today) {
      nextBirthday.setFullYear(today.getFullYear() + 1);
    }
    const diffTime = nextBirthday - today;
    const daysUntil = Math.ceil(diffTime / (1000 * 60 * 60 * 24));

    this.setData({
      isBirthday: isToday,
      birthday: birthday,
      daysUntil: daysUntil
    });
  },

  // 获取礼品信息
  fetchGiftInfo() {
    get('/wx/birthday/gift').then(res => {
      if (res) {
        this.setData({
          giftInfo: res,
          giftClaimed: res.status === 1
        });
      }
    }).catch(err => {
      console.log('获取礼品信息失败');
    });
  },

  // 获取礼品历史
  fetchGiftHistory() {
    get('/wx/birthday/history').then(res => {
      const list = Array.isArray(res) ? res : [];
      this.setData({ giftHistory: list });
    }).catch(err => {
      console.log('获取礼品历史失败');
    });
  },

  // 领取礼品
  claimGift() {
    if (this.data.giftClaimed) return;

    this.setData({ claiming: true });
    post('/wx/birthday/claim', {}).then(res => {
      wx.showToast({ title: '领取成功', icon: 'success' });
      this.setData({ giftClaimed: true });
    }).catch(err => {
      wx.showToast({ title: '领取失败', icon: 'none' });
    }).finally(() => {
      this.setData({ claiming: false });
    });
  },

  // 设置生日
  setBirthday() {
    const now = new Date();
    const currentYear = now.getFullYear();
    
    wx.showModal({
      title: '设置生日',
      editable: true,
      placeholderText: '格式: yyyy-MM-DD',
      content: this.data.birthday || `${currentYear - 25}-01-01`,
      success: (res) => {
        if (!res.confirm || !res.content) return;
        const birthday = res.content.trim();
        const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
        if (!dateRegex.test(birthday)) {
          wx.showToast({ title: '日期格式错误', icon: 'none' });
          return;
        }
        put('/wx/user/birthday', { birthday }).then(() => {
          wx.showToast({ title: '生日设置成功', icon: 'success' });
          this.checkBirthday();
        }).catch(err => {
          wx.showToast({ title: '设置失败', icon: 'none' });
        });
      }
    });
  }
});
