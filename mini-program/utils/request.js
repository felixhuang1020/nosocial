/**
 * 网络请求封装
 * 统一处理 baseURL、token、错误处理
 */

// 从配置或环境变量获取 BASE_URL，支持开发和生产环境
const getBaseUrl = () => {
  // 尝试从全局配置读取
  const app = getApp && getApp();
  if (app && app.globalData && app.globalData.apiBaseUrl) {
    return app.globalData.apiBaseUrl;
  }
  // 开发环境默认值（微信模拟器内 localhost 会被代理到宿主机，127.0.0.1 不行）
  return 'http://localhost:8080/api/v1';
};

const request = (options) => {
  return new Promise((resolve, reject) => {
    const token = wx.getStorageSync('token') || '';
    
    wx.request({
      url: `${getBaseUrl()}${options.url}`,
      enableHttp2: false,
      enableQuic: false,
      method: options.method || 'GET',
      data: options.data,
      header: {
        'Authorization': token ? `Bearer ${token}` : '',
        'Content-Type': 'application/json'
      },
      timeout: 30000,
      success: (res) => {
        if (res.statusCode === 401) {
          // Token 过期，清除登录状态并跳转登录
          wx.removeStorageSync('token');
          const app = getApp();
          if (app) {
            app.handleLoginExpired();
          }
          wx.showToast({ 
            title: '登录已过期，请重新登录', 
            icon: 'none',
            duration: 2000
          });
          setTimeout(() => {
            wx.switchTab({ url: '/pages/profile/profile' });
          }, 2000);
          return reject({ code: 401, message: '登录已过期' });
        }
        
        if (res.statusCode >= 500) {
          wx.showToast({ title: '服务器繁忙，请稍后重试', icon: 'none' });
          return reject({ code: res.statusCode, message: '服务器错误' });
        }
        
        if (res.data.code !== 0) {
          const msg = res.data.msg || res.data.message || '请求失败';
          wx.showToast({ title: msg, icon: 'none' });
          return reject({ code: res.data.code, message: msg });
        }
        
        resolve(res.data.data);
      },
      fail: (err) => {
        console.error('请求失败:', err);
        const isRefused = err && err.errMsg && err.errMsg.includes('ECONNREFUSED');
        wx.showToast({ title: isRefused ? '无法连接服务器' : '网络请求失败', icon: 'none' });
        reject({ code: -1, message: isRefused ? '无法连接服务器' : '网络请求失败' });
      }
    });
  });
};

// GET 请求
const get = (url, params = {}) => {
  // 构建查询字符串
  const queryString = Object.keys(params)
    .filter(key => params[key] !== undefined && params[key] !== null)
    .map(key => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
    .join('&');
  
  const fullUrl = queryString ? `${url}?${queryString}` : url;
  return request({ url: fullUrl, method: 'GET' });
};

// POST 请求
const post = (url, data = {}) => {
  return request({ url, method: 'POST', data });
};

// PUT 请求
const put = (url, data = {}) => {
  return request({ url, method: 'PUT', data });
};

// DELETE 请求
const del = (url, data = {}) => {
  return request({ url, method: 'DELETE', data });
};

// 上传文件
const upload = (filePath, options = {}) => {
  return new Promise((resolve, reject) => {
    const token = wx.getStorageSync('token') || '';
    
    wx.uploadFile({
      url: `${getBaseUrl()}${options.url || '/wx/upload'}`,
      filePath: filePath,
      name: options.name || 'file',
      formData: options.formData || {},
      header: {
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (res) => {
        if (res.statusCode === 200) {
          try {
            const data = JSON.parse(res.data);
            if (data.code === 0) {
              resolve(data.data);
            } else {
              reject({ code: data.code, message: data.msg || '上传失败' });
            }
          } catch (e) {
            resolve(res.data);
          }
        } else {
          reject({ code: res.statusCode, message: '上传失败' });
        }
      },
      fail: (err) => {
        reject({ code: -1, message: '上传失败' });
      }
    });
  });
};

// OSS 直传：签名模式
// 1) 向后端 GET /wx/upload/signature 拿签名
// 2) wx.uploadFile 直交 OSS host (multipart/form-data, PostObject)
const inferContentType = (filePath) => {
  const ext = ((filePath || '').split('.').pop() || '').toLowerCase();
  const map = {
    png: 'image/png',
    jpg: 'image/jpeg',
    jpeg: 'image/jpeg',
    gif: 'image/gif',
    webp: 'image/webp',
    bmp: 'image/bmp'
  };
  return map[ext] || 'application/octet-stream';
};

const uploadToOSS = (filePath, dir = 'reviews') => {
  return new Promise((resolve, reject) => {
    const ext = ((filePath || '').split('.').pop() || '').toLowerCase();
    get('/wx/upload/signature', { dir, ext }).then(sig => {
      if (!sig || !sig.host || !sig.key) {
        return reject({ code: -1, message: '签名数据不完整' });
      }
      wx.uploadFile({
        url: sig.host,
        filePath,
        name: 'file',
        // 注意：PostObject 协议要求 Content-Type 在 file 之前（wx.uploadFile 将 formData 置于 file 部分前，天然满足）
        // 不带 Content-Type 会被存为 application/octet-stream，导致前端图片加载失败
        formData: {
          key: sig.key,
          policy: sig.policy,
          OSSAccessKeyId: sig.OSSAccessKeyId,
          signature: sig.signature,
          success_action_status: '200',
          'Content-Type': inferContentType(filePath)
        },
        success: (res) => {
          if (res.statusCode === 200 || res.statusCode === 204) {
            // PostObject 不支持通过 form 设置 Content-Disposition，上传成功后异步修复
            put('/wx/upload/inline', { key: sig.key }).catch(() => {});
            resolve({ url: `${sig.host}/${sig.key}`, key: sig.key });
          } else {
            const body = (res.data || '').toString().slice(0, 120);
            reject({ code: res.statusCode, message: `OSS 上传失败 (${res.statusCode}) ${body}` });
          }
        },
        fail: () => reject({ code: -1, message: 'OSS 上传失败，请检查网络' })
      });
    }).catch(reject);
  });
};

module.exports = {
  request,
  get,
  post,
  put,
  del,
  upload,
  uploadToOSS
};
