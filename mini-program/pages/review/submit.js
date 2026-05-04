import { post, uploadToOSS } from '../../utils/request';

Page({
  data: {
    screenshots: [],
    reviewContent: '',
    submitting: false,
    canSubmit: false
  },

  onLoad() {
    // 页面加载
  },

  onUnload() {
    // 标记页面已销毁，防止异步回调中调用 setData
    this._destroyed = true;
  },

  // 计算提交按钮状态
  updateSubmitState() {
    const { screenshots, submitting } = this.data;
    const hasValidImage = screenshots.some(s => s.url && !s.uploading && !s.error);
    const hasUploading = screenshots.some(s => s.uploading);
    const canSubmit = hasValidImage && !hasUploading && !submitting;
    this.setData({ canSubmit });
  },

  // 选择图片（支持多张）
  chooseImages() {
    const remainCount = 9 - this.data.screenshots.length;
    if (remainCount <= 0) {
      wx.showToast({ title: '最多9张图片', icon: 'none' });
      return;
    }

    wx.chooseMedia({
      count: remainCount,
      mediaType: ['image'],
      sourceType: ['album', 'camera'],
      success: (res) => {
        const newImages = res.tempFiles.map((f, idx) => ({
          id: `img_${Date.now()}_${idx}`,
          url: '',
          tempFilePath: f.tempFilePath,
          uploading: true
        }));

        this.setData({
          screenshots: [...this.data.screenshots, ...newImages]
        });
        this.updateSubmitState();

        // 逐个上传（OSS 直传）
        newImages.forEach((img, idx) => {
          const tempId = newImages[idx].id;
          uploadToOSS(img.tempFilePath, 'reviews').then(uploadRes => {
            if (this._destroyed) return;
            const screenshots = [...this.data.screenshots];
            const targetIdx = screenshots.findIndex(s => s.id === tempId);
            if (targetIdx >= 0) {
              screenshots[targetIdx] = {
                ...screenshots[targetIdx],
                url: uploadRes.url,
                tempFilePath: '',
                uploading: false
              };
              this.setData({ screenshots });
              this.updateSubmitState();
            }
          }).catch(err => {
            if (this._destroyed) return;
            const screenshots = [...this.data.screenshots];
            const targetIdx = screenshots.findIndex(s => s.id === tempId);
            if (targetIdx >= 0) {
              screenshots[targetIdx] = {
                ...screenshots[targetIdx],
                uploading: false,
                error: true
              };
              this.setData({ screenshots });
              this.updateSubmitState();
            }
            wx.showToast({ title: (err && err.message) || '有图片上传失败', icon: 'none' });
          });
        });
      }
    });
  },

  // 删除图片
  deleteImage(e) {
    const index = e.currentTarget.dataset.index;
    const screenshots = [...this.data.screenshots];
    screenshots.splice(index, 1);
    this.setData({ screenshots });
    this.updateSubmitState();
  },

  // 重试上传
  retryUpload(e) {
    const index = e.currentTarget.dataset.index;
    const item = this.data.screenshots[index];
    if (!item || !item.tempFilePath) {
      // 如果没有本地文件路径了，就删除这个错误项
      const screenshots = this.data.screenshots.filter((_, i) => i !== index);
      this.setData({ screenshots });
      this.updateSubmitState();
      return;
    }

    const screenshots = [...this.data.screenshots];
    const tempId = item.id;
    screenshots[index] = { ...item, uploading: true, error: false };
    this.setData({ screenshots });
    this.updateSubmitState();

    uploadToOSS(item.tempFilePath, 'reviews').then(uploadRes => {
      if (this._destroyed) return;
      const screenshots = [...this.data.screenshots];
      const targetIdx = screenshots.findIndex(s => s.id === tempId);
      if (targetIdx >= 0) {
        screenshots[targetIdx] = {
          ...screenshots[targetIdx],
          url: uploadRes.url,
          uploading: false,
          error: false
        };
        this.setData({ screenshots });
        this.updateSubmitState();
      }
    }).catch(err => {
      if (this._destroyed) return;
      const screenshots = [...this.data.screenshots];
      const targetIdx = screenshots.findIndex(s => s.id === tempId);
      if (targetIdx >= 0) {
        screenshots[targetIdx] = { ...screenshots[targetIdx], uploading: false, error: true };
        this.setData({ screenshots });
        this.updateSubmitState();
      }
      wx.showToast({ title: '上传失败，请点击重试', icon: 'none' });
    });
  },

  // 预览图片
  previewImage(e) {
    const index = e.currentTarget.dataset.index;
    const urls = this.data.screenshots.map(s => s.url || s.tempFilePath).filter(Boolean);
    if (urls.length === 0) return;
    wx.previewImage({
      current: urls[index],
      urls
    });
  },

  // 内容变化
  onContentChange(e) {
    this.setData({ reviewContent: e.detail.value });
  },

  // 提交点评
  submitReview() {
    const uploaded = this.data.screenshots.filter(s => s.url && !s.uploading);
    if (uploaded.length === 0) {
      wx.showToast({ title: '请上传至少一张截图', icon: 'none' });
      return;
    }

    // 检查是否还有未上传完成的
    const uploading = this.data.screenshots.some(s => s.uploading);
    if (uploading) {
      wx.showToast({ title: '图片上传中，请稍候', icon: 'none' });
      return;
    }

    this.setData({ submitting: true });
    this.updateSubmitState();

    // 提交点评（支持多图 image_urls）
    const imageURLs = uploaded.map(s => s.url);
    post('/wx/review/submit', {
      image_urls: imageURLs,
      screenshot_url: imageURLs[0], // 兼容旧后端
      content: this.data.reviewContent
    }).then(res => {
      wx.showModal({
        title: '提交成功',
        content: '审核通过后将自动发放20元优惠券到您的账户',
        showCancel: false,
        confirmColor: '#7C9A92',
        success: () => {
          wx.navigateBack();
        }
      });
    }).catch(err => {
      wx.showToast({ title: '提交失败，请重试', icon: 'none' });
    }).finally(() => {
      if (this._destroyed) return;
      this.setData({ submitting: false });
      this.updateSubmitState();
    });
  }
});
