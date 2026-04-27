import { createStoreBindings } from 'mobx-miniprogram-bindings';
import { appStore } from '../../stores/app';
import { get, post } from '../../utils/request';

const STATUS_MAP = {
  0: { text: '待支付', color: '#ff9800', type: 'warning' },
  1: { text: '已支付', color: '#4caf50', type: 'success' },
  2: { text: '制作中', color: '#2196f3', type: 'primary' },
  3: { text: '已完成', color: '#9e9e9e', type: 'default' },
  4: { text: '已取消', color: '#f44336', type: 'danger' }
};

Page({
  data: {
    orders: [],
    page: 1,
    size: 10,
    loading: false,
    hasMore: true,
    total: 0
  },

  storeBindings: null,

  onLoad() {
    this.storeBindings = createStoreBindings(this, {
      store: appStore,
      fields: ['isLogin'],
      actions: []
    });

    if (this.data.isLogin) {
      this.fetchOrders(true);
    }
  },

  onShow() {
    if (this.storeBindings) {
      this.storeBindings.updateStoreBindings();
    }
    if (this.data.isLogin && this.data.orders.length === 0) {
      this.fetchOrders(true);
    }
  },

  onUnload() {
    if (this.storeBindings) {
      this.storeBindings.destroyStoreBindings();
    }
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.fetchOrders(true).finally(() => {
      wx.stopPullDownRefresh();
    });
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.fetchOrders();
    }
  },

  // 获取订单列表
  fetchOrders(reset = false) {
    if (!this.data.isLogin) {
      return Promise.resolve();
    }

    if (this.data.loading) return Promise.resolve();

    this.setData({ loading: true });

    const page = reset ? 1 : this.data.page;

    return get('/wx/orders', { page, size: this.data.size }).then(res => {
      const list = res && res.list ? res.list : [];
      const total = res && res.total !== undefined ? res.total : 0;
      const hasMore = list.length === this.data.size;

      // 处理订单数据
      const processedList = list.map(item => ({
        ...item,
        statusText: STATUS_MAP[item.status] ? STATUS_MAP[item.status].text : '未知',
        statusColor: STATUS_MAP[item.status] ? STATUS_MAP[item.status].color : '#999'
      }));

      this.setData({
        orders: reset ? processedList : [...this.data.orders, ...processedList],
        page: page + 1,
        hasMore,
        total,
        loading: false
      });
    }).catch(err => {
      console.error('获取订单失败:', err);
      this.setData({ loading: false });
    });
  },

  // 支付订单
  payOrder(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({
      title: '确认支付',
      content: '确认支付该订单？',
      confirmColor: '#7C9A92',
      success: (res) => {
        if (res.confirm) {
          post(`/wx/orders/${id}/pay`).then(res => {
            if (res && res.pay_params) {
              wx.requestPayment({
                ...res.pay_params,
                success: () => {
                  wx.showToast({ title: '支付成功', icon: 'success' });
                  this.fetchOrders(true);
                },
                fail: () => {
                  wx.showToast({ title: '支付取消', icon: 'none' });
                }
              });
            }
          }).catch(err => {
            wx.showToast({ title: '发起支付失败', icon: 'none' });
          });
        }
      }
    });
  },

  // 取消订单
  cancelOrder(e) {
    const id = e.currentTarget.dataset.id;
    wx.showModal({
      title: '确认取消',
      content: '确定要取消该订单吗？',
      confirmColor: '#ff4d4f',
      success: (res) => {
        if (res.confirm) {
          post(`/wx/orders/${id}/cancel`).then(() => {
            wx.showToast({ title: '取消成功', icon: 'success' });
            this.fetchOrders(true);
          }).catch(err => {
            wx.showToast({ title: '取消失败', icon: 'none' });
          });
        }
      }
    });
  },

  // 去登录
  goLogin() {
    wx.switchTab({ url: '/pages/profile/profile' });
  },

  // 去点单
  goOrder() {
    wx.switchTab({ url: '/pages/menu/menu' });
  }
});
