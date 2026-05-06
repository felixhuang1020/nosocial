const { DIANPING_URL } = require('../../utils/constants');

Page({
  data: {
    dianpingUrl: DIANPING_URL
  },

  onLoad() {
    // 页面加载
  },

  // 跳转到大众点评
  // 微信小程序 web-view 只能打开已在公众平台配置的业务域名，
  // 第三方网址（如 m.dianping.com）无法配置，直接 navigateTo 会提示“不在合法域名列表”。
  // 正确交互：复制链接 → 提示用户在大众点评 APP / 浏览器中粘贴打开。
  goToDianping() {
    wx.setClipboardData({
      data: this.data.dianpingUrl,
      success: () => {
        wx.showModal({
          title: '链接已复制',
          content: '已复制大众点评店铺链接。\n\n推荐方式：\n1. 打开大众点评 APP，在搜索框粘贴访问\n2. 或将链接粘贴到浏览器中打开\n\n完成点评后请返回小程序上传截图领取优惠券。',
          confirmText: '我知道了',
          showCancel: false
        });
      },
      fail: (err) => {
        console.error('[review.guide] 复制链接失败:', err);
        wx.showToast({ title: '复制失败，请手动输入大众点评搜索 NoSocial', icon: 'none', duration: 3000 });
      }
    });
  },

  // 跳转到提交页面
  goToSubmit() {
    wx.navigateTo({ url: '/pages/review/submit' });
  }
});
