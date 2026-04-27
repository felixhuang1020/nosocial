import { createStoreBindings } from 'mobx-miniprogram-bindings';
import { appStore } from '../../stores/app';
import { get, post } from '../../utils/request';

Page({
  data: {
    drinkId: 0,
    drink: {},
    ingredientList: [],
    isFreeDrink: false,
    ordering: false
  },

  storeBindings: null,

  onLoad(options) {
    this.storeBindings = createStoreBindings(this, {
      store: appStore,
      fields: ['isLogin', 'isShareholder'],
      actions: []
    });

    const id = options.id || 0;
    this.setData({ drinkId: id });
    this.fetchDrinkDetail(id);
  },

  onShow() {
    if (this.storeBindings) {
      this.storeBindings.updateStoreBindings();
    }
    // 股东状态更新后重新计算免费酒水标识
    if (this.data.drink && this.data.drink.id) {
      const isFreeDrink = this.data.drink.is_free_drink && this.data.isShareholder;
      if (isFreeDrink !== this.data.isFreeDrink) {
        this.setData({ isFreeDrink });
      }
    }
  },

  onUnload() {
    if (this.storeBindings) {
      this.storeBindings.destroyStoreBindings();
    }
  },

  // 获取酒水详情
  fetchDrinkDetail(id) {
    get(`/public/drinks/${id}`).then(res => {
      this.processDrinkData(res);
    }).catch(err => {
      console.error('获取酒水详情失败:', err);
      wx.showToast({ title: '获取详情失败', icon: 'none' });
    });
  },

  // 处理酒水数据
  processDrinkData(drink) {
    const ingredientList = drink.ingredients ? drink.ingredients.split(',').map(s => s.trim()) : [];
    const isFreeDrink = drink.is_free_drink && this.data.isShareholder;
    
    this.setData({
      drink: drink,
      ingredientList: ingredientList,
      isFreeDrink: isFreeDrink
    });
  },

  // 下单
  placeOrder() {
    if (!this.data.isLogin) {
      wx.showModal({
        title: '提示',
        content: '请先登录后再下单',
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

    if (this.data.isFreeDrink) {
      this.claimFreeDrink();
    } else {
      this.createOrder();
    }
  },

  // 领取免费酒水
  claimFreeDrink() {
    this.setData({ ordering: true });
    post('/wx/free-drink/claim').then(() => {
      wx.showToast({ title: '领取成功', icon: 'success' });
      setTimeout(() => {
        wx.navigateTo({ url: '/pages/order/list' });
      }, 1500);
    }).catch(err => {
      wx.showToast({ title: err.message || '领取失败', icon: 'none' });
    }).finally(() => {
      this.setData({ ordering: false });
    });
  },

  // 创建订单
  createOrder() {
    wx.showModal({
      title: '确认下单',
      content: `确认购买「${this.data.drink.name}」？\n金额: ¥${this.data.drink.price}`,
      confirmColor: '#7C9A92',
      success: (res) => {
        if (res.confirm) {
          this.setData({ ordering: true });
          const price = parseFloat(this.data.drink.price) || 0;
          post('/wx/orders', {
            total_amount: price,
            discount_amount: 0
          }).then(res => {
            wx.showToast({ title: '下单成功', icon: 'success' });
            setTimeout(() => {
              wx.navigateTo({ url: '/pages/order/list' });
            }, 1500);
          }).catch(err => {
            wx.showToast({ title: '下单失败', icon: 'none' });
            this.setData({ ordering: false });
          });
        }
      }
    });
  },

  // 发起支付
  requestPayment(payParams) {
    wx.requestPayment({
      ...payParams,
      success: (res) => {
        wx.showToast({ title: '支付成功', icon: 'success' });
        setTimeout(() => {
          wx.navigateTo({ url: '/pages/order/list' });
        }, 1500);
      },
      fail: (err) => {
        wx.showToast({ title: '支付取消', icon: 'none' });
      },
      complete: () => {
        this.setData({ ordering: false });
      }
    });
  }
});
      },
      fail: (err) => {
        wx.showToast({ title: '支付取消', icon: 'none' });
      },
      complete: () => {
        this.setData({ ordering: false });
      }
    });
  }
});
