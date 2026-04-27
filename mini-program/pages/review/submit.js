import { post, uploadToOSS } from '../../utils/request';

Page({
  data: {
    screenshots: [],
    reviewContent: '',
    submitting: false
  },

  onLoad() {
    // 页面加载
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
        const newImages = res.tempFiles.map(f => ({
          url: '',
          tempFilePath: f.tempFilePath,
          uploading: true
        }));

        const startIdx = this.data.screenshots.length;
        this.setData({
          screenshots: [...this.data.screenshots, ...newImages]
        });

        // 逐个上传（OSS 直传）
        newImages.forEach((img, idx) => {
          const realIdx = startIdx + idx;
          uploadToOSS(img.tempFilePath, 'reviews').then(uploadRes => {
            const screenshots = [...this.data.screenshots];
            if (screenshots[realIdx]) {
              screenshots[realIdx] = {
                url: uploadRes.url,
                tempFilePath: '',
                uploading: false
              };
              this.setData({ screenshots });
            }
          }).catch(err => {
            const screenshots = [...this.data.screenshots];
            if (screenshots[realIdx]) {
              screenshots[realIdx].uploading = false;
              screenshots[realIdx].error = true;
              this.setData({ screenshots });
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
      this.setData({ submitting: false });
    });
  }
});
