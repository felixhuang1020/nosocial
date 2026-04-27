import { get } from '../../utils/request';
import { COUPON_TYPES, COUPON_STATUS } from '../../utils/constants';

Page({
  data: {
    couponList: [],
    loading: false
  },

  onLoad() {
    this.fetchCoupons();
  },

  onShow() {
    this.fetchCoupons();
  },

  // 获取优惠券列表
  fetchCoupons() {
    this.setData({ loading: true });
    
    get('/wx/coupons').then(res => {
      const list = res && res.list ? res.list : [];
      
      // 处理数据
      const processedList = list.map(item => ({
        ...item,
        typeName: COUPON_TYPES[item.type] || '优惠券',
        statusName: COUPON_STATUS[item.status] || '未知'
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
