Page({
  data: {
    url: ''
  },

  onLoad(options) {
    const url = options.url || '';
    const title = options.title || '网页';
    
    this.setData({ url: decodeURIComponent(url) });
    
    wx.setNavigationBarTitle({ title });
  },

  onLoadSuccess() {
    console.log('WebView 加载成功');
  },

  onLoadError() {
    wx.showToast({ title: '页面加载失败', icon: 'none' });
  }
});
