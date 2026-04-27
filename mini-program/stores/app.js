/**
 * 全局状态管理 - 使用 MobX-miniprogram
 * 管理用户登录状态、股东信息等全局数据
 */

import { observable, action } from 'mobx-miniprogram';

export const appStore = observable({
  // ========== 状态数据 ==========
  token: '',
  userInfo: null,
  isShareholder: false,
  shareholderInfo: null,
  unreadCount: 0,
  
  // ========== 计算属性 ==========
  get isLogin() {
    return !!this.token;
  },
  
  get userId() {
    return this.userInfo ? this.userInfo.id : 0;
  },
  
  get nickname() {
    return this.userInfo ? this.userInfo.nickname : '';
  },
  
  get avatar() {
    return this.userInfo ? this.userInfo.avatar : '';
  },
  
  get inviteCode() {
    return this.userInfo ? this.userInfo.invite_code : '';
  },
  
  get birthday() {
    return this.userInfo ? this.userInfo.birthday : '';
  },
  
  get balance() {
    return this.userInfo ? this.userInfo.balance : 0;
  },
  
  get totalEarning() {
    return this.userInfo ? this.userInfo.total_earning : 0;
  },
  
  // ========== Actions ==========
  
  // 设置 Token
  setToken: action(function(token) {
    this.token = token;
    if (token) {
      wx.setStorageSync('token', token);
    }
  }),
  
  // 设置用户信息
  setUserInfo: action(function(info) {
    this.userInfo = info;
    this.isShareholder = info ? info.is_shareholder : false;
    if (info) {
      wx.setStorageSync('user_info', JSON.stringify(info));
    }
  }),
  
  // 更新股东信息
  setShareholderInfo: action(function(info) {
    this.shareholderInfo = info;
    if (info) {
      this.isShareholder = true;
    }
  }),
  
  // 更新用户信息字段
  updateUserField: action(function(field, value) {
    if (this.userInfo) {
      this.userInfo[field] = value;
      wx.setStorageSync('user_info', JSON.stringify(this.userInfo));
    }
  }),
  
  // 更新余额
  updateBalance: action(function(amount) {
    if (this.userInfo) {
      this.userInfo.balance = (parseFloat(this.userInfo.balance) + parseFloat(amount)).toFixed(2);
      wx.setStorageSync('user_info', JSON.stringify(this.userInfo));
    }
  }),
  
  // 设置未读消息数
  setUnreadCount: action(function(count) {
    this.unreadCount = count;
  }),
  
  // 清除用户数据（登出）
  clearUserInfo: action(function() {
    this.token = '';
    this.userInfo = null;
    this.isShareholder = false;
    this.shareholderInfo = null;
    this.unreadCount = 0;
    wx.removeStorageSync('token');
    wx.removeStorageSync('user_info');
  }),
  
  // 从本地存储恢复状态
  restoreFromStorage: action(function() {
    const token = wx.getStorageSync('token');
    const userInfoStr = wx.getStorageSync('user_info');
    
    if (token) {
      this.token = token;
    }
    
    if (userInfoStr) {
      try {
        const userInfo = JSON.parse(userInfoStr);
        this.userInfo = userInfo;
        this.isShareholder = userInfo.is_shareholder || false;
      } catch (e) {
        console.error('解析用户信息失败:', e);
      }
    }
  })
});
