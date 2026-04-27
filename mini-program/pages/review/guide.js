Page({
  data: {
    dianpingUrl: 'https://m.dianping.com/shop/NoSocialBar'
  },

  onLoad() {
    // 页面加载
  },

  // 跳转到大众点评
  goToDianping() {
    // 尝试打开外部网页
    wx.navigateTo({
      url: `/pages/webview/webview?url=${encodeURIComponent(this.data.dianpingUrl)}&title=大众点评`
    });
  },

  // 跳转到提交页面
  goToSubmit() {
    wx.navigateTo({ url: '/pages/review/submit' });
  }
});
