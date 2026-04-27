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
    giftHistory: [],
    showDatePicker: false,
    pickerValue: '',
    currentDate: ''
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
    // 计算当前日期作为选择器上限
    const now = new Date();
    const y = now.getFullYear();
    const m = String(now.getMonth() + 1).padStart(2, '0');
    const d = String(now.getDate()).padStart(2, '0');
    this.setData({
      currentDate: `${y}-${m}-${d}`,
      pickerValue: this.data.birthday || `${y - 25}-01-01`
    });
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

  // 打开生日选择器
  openDatePicker() {
    const now = new Date();
    const y = now.getFullYear();
    this.setData({
      showDatePicker: true,
      pickerValue: this.data.birthday || `${y - 25}-01-01`
    });
  },

  // 日期选择器确认
  onDateConfirm(e) {
    const birthday = e.detail.value;
    this.setData({ showDatePicker: false });
    put('/wx/user/birthday', { birthday }).then(() => {
      appStore.updateUserField('birthday', birthday);
      wx.showToast({ title: '生日设置成功', icon: 'success' });
      this.checkBirthday();
    }).catch(err => {
      wx.showToast({ title: '设置失败', icon: 'none' });
    });
  },

  // 日期选择器取消
  onDateCancel() {
    this.setData({ showDatePicker: false });
  }
});
