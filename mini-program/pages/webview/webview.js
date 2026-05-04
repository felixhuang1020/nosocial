Page({
  data: {
    url: ''
  },
  onLoad(options) {
    let url = options.url ? decodeURIComponent(options.url) : '';
    
    // URL 白名单验证
    const allowedDomains = [
      'https://m.dianping.com',
      'https://www.dianping.com',
    ];
    
    let isAllowed = false;
    try {
      if (url) {
        isAllowed = allowedDomains.some(domain => url.startsWith(domain));
      }
    } catch (e) {
      isAllowed = false;
    }
    
    if (!isAllowed) {
      wx.showToast({ title: '链接不安全', icon: 'none' });
      setTimeout(() => { wx.navigateBack(); }, 1500);
      return;
    }
    
    this.setData({ url });
    wx.setNavigationBarTitle({ title: options.title || '网页' });
  },
  onLoadSuccess() {
    console.log('WebView 加载成功');
  },
  onLoadError() {
    wx.showToast({ title: '页面加载失败', icon: 'none' });
  }
});
