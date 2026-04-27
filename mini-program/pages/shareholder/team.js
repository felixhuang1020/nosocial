import { createStoreBindings } from 'mobx-miniprogram-bindings';
import { appStore } from '../../stores/app';
import { get } from '../../utils/request';

Page({
  data: {
    totalMembers: 0,
    activeMembers: 0,
    totalCommission: '0.00',
    teamList: [],
    loading: false
  },

  storeBindings: null,

  onLoad() {
    this.storeBindings = createStoreBindings(this, {
      store: appStore,
      fields: ['inviteCode'],
      actions: []
    });

    this.fetchTeam();
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

  // 获取团队数据
  fetchTeam() {
    this.setData({ loading: true });

    get('/wx/shareholder/team').then(res => {
      const list = Array.isArray(res) ? res : (res && res.list ? res.list : []);
      this.setData({
        teamList: list,
        totalMembers: list.length,
        activeMembers: list.filter(m => m.is_shareholder === 1).length,
        totalCommission: '0.00',
        loading: false
      });
    }).catch(err => {
      console.error('获取团队数据失败:', err);
      this.setData({
        teamList: [],
        totalMembers: 0,
        activeMembers: 0,
        totalCommission: '0.00',
        loading: false
      });
    });
  },

  // 分享邀请
  shareInvite() {
    wx.showShareMenu({
      withShareTicket: true,
      menus: ['shareAppMessage']
    });
  },
});
