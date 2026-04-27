import { appStore } from './stores/app';

App({
  globalData: {
    userInfo: null,
    token: '',
    isShareholder: false,
    systemInfo: null,
    statusBarHeight: 0,
    navBarHeight: 44,
    screenWidth: 375,
    screenHeight: 812,
    safeAreaBottom: 0,
    // API 基础地址
    // 开发环境：http://127.0.0.1:8080/api/v1
    // 生产环境：在 project.config.json / 微信后台配置合法域名后替换
    apiBaseUrl: ''
  },

  store: appStore,

  onLaunch(options) {
    console.log('App Launch', options);
    this.initSystemInfo();
    this.checkLoginStatus();
    
    // 处理扫码场景值
    const scene = options.scene;
    if (scene === 1047 || scene === 1048 || scene === 1049) {
      // 扫描小程序码进入
      const query = options.query;
      if (query.scene) {
        const sceneStr = decodeURIComponent(query.scene);
        console.log('扫码参数:', sceneStr);
        wx.setStorageSync('invite_code', sceneStr);
      }
    }
  },

  onShow(options) {
    console.log('App Show', options);
  },

  onHide() {
    console.log('App Hide');
  },

  onError(msg) {
    console.error('App Error:', msg);
  },

  // 初始化系统信息
  initSystemInfo() {
    const systemInfo = wx.getSystemInfoSync();
    this.globalData.systemInfo = systemInfo;
    this.globalData.statusBarHeight = systemInfo.statusBarHeight;
    this.globalData.screenWidth = systemInfo.screenWidth;
    this.globalData.screenHeight = systemInfo.screenHeight;
    this.globalData.safeAreaBottom = systemInfo.safeArea ? 
      (systemInfo.screenHeight - systemInfo.safeArea.bottom) : 0;
    
    // 计算导航栏高度
    const menuButtonInfo = wx.getMenuButtonBoundingClientRect();
    this.globalData.navBarHeight = (menuButtonInfo.top - systemInfo.statusBarHeight) * 2 + menuButtonInfo.height;
  },

  // 检查登录状态
  checkLoginStatus() {
    const token = wx.getStorageSync('token');
    if (token) {
      this.globalData.token = token;
      this.store.setToken(token);
      this.fetchUserInfo();
    }
  },

  // 获取用户信息
  fetchUserInfo() {
    const { get } = require('./utils/request');
    get('/wx/user').then(res => {
      this.globalData.userInfo = res;
      this.globalData.isShareholder = res.is_shareholder;
      this.store.setUserInfo(res);
    }).catch(err => {
      console.error('获取用户信息失败:', err);
      if (err.code === 401) {
        this.handleLoginExpired();
      }
    });
  },

  // 处理登录过期
  handleLoginExpired() {
    wx.removeStorageSync('token');
    this.globalData.token = '';
    this.globalData.userInfo = null;
    this.globalData.isShareholder = false;
    this.store.clearUserInfo();
  },

  // 微信登录
  wxLogin(inviteCode = '') {
    return new Promise((resolve, reject) => {
      wx.login({
        success: (res) => {
          if (res.code) {
            const { post } = require('./utils/request');
            post('/wx/login', {
              code: res.code,
              invite_code: inviteCode
            }).then(data => {
              wx.setStorageSync('token', data.token);
              this.globalData.token = data.token;
              this.globalData.userInfo = data.user_info;
              this.globalData.isShareholder = data.user_info.is_shareholder;
              this.store.setToken(data.token);
              this.store.setUserInfo(data.user_info);
              resolve(data);
            }).catch(reject);
          } else {
            reject(new Error('微信登录失败'));
          }
        },
        fail: reject
      });
    });
  }
});
