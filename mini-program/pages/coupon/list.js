import { get } from '../../utils/request';
import { COUPON_TYPES, COUPON_STATUS } from '../../utils/constants';

Page({
  data: {
    couponList: [],
    loading: false
  },

  _loaded: false,

  onLoad() {
    this._loaded = true;
    this.fetchCoupons();
  },

  onShow() {
    // onLoad 已加载，onShow 仅在从其他页面返回时刷新
    if (this._loaded && this.data.couponList.length > 0) {
      this.fetchCoupons();
    }
  },

  // 获取优惠券列表
  fetchCoupons() {
    this.setData({ loading: true });
    
    get('/wx/coupons').then(res => {
      const list = res && res.list ? res.list : [];
      
      // 处理数据
      const processedList = list.map(item => ({
        ...item,
        typeName: COUPON_TYPES[item.type] ? COUPON_TYPES[item.type].name : '优惠券',
        statusName: COUPON_STATUS[item.status] ? COUPON_STATUS[item.status].name : '未知'
      }));
      
      this.setData({
        couponList: processedList,
        loading: false
      });
    }).catch(err => {
      console.error('获取优惠券失败:', err);
      this.setData({
        couponList: [],
        loading: false
      });
    });
  },
});
