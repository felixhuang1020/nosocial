import { get } from '../../utils/request';

Page({
  data: {
    galleryList: [],
    leftColumn: [],
    rightColumn: [],
    page: 1,
    size: 10,
    loading: false,
    hasMore: true,
    leftHeight: 0,
    rightHeight: 0
  },

  onLoad() {
    this.fetchGallery();
  },

  onShow() {
    // 页面显示
  },

  onReady() {
    // 页面首次渲染完成
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.setData({
      galleryList: [],
      leftColumn: [],
      rightColumn: [],
      leftHeight: 0,
      rightHeight: 0,
      page: 1,
      hasMore: true
    }, () => {
      this.fetchGallery().finally(() => {
        wx.stopPullDownRefresh();
      });
    });
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.fetchGallery();
    }
  },

  // 获取画廊数据
  fetchGallery() {
    if (this.data.loading) return Promise.resolve();

    this.setData({ loading: true });

    return get('/public/banners', { position: 2 }).then(res => {
      const list = Array.isArray(res) ? res : (res && res.list ? res.list : []);
      // 过滤出 position=2 的展示图，只取未加载的
      const existingIds = new Set(this.data.galleryList.map(i => i.id));
      const newItems = list
        .filter(item => !existingIds.has(item.id))
        .map(item => ({
          ...item,
          loaded: false
        }));
      const hasMore = false; // banners 接口不分页

      this.setData({
        galleryList: [...this.data.galleryList, ...newItems],
        hasMore: hasMore,
        loading: false
      }, () => {
        this.distributeColumns(newItems);
      });
    }).catch(err => {
      console.error('获取展示数据失败:', err);
      this.setData({
        hasMore: false,
        loading: false
      });
    });
  },

  // 分配到左右列
  distributeColumns(newItems) {
    const { leftHeight, rightHeight } = this.data;
    const leftColumn = [...this.data.leftColumn];
    const rightColumn = [...this.data.rightColumn];
    
    // 这里简化处理，实际应根据图片高度计算
    // 奇数索引放左列，偶数索引放右列
    newItems.forEach((item, index) => {
      if (index % 2 === 0) {
        leftColumn.push(item);
      } else {
        rightColumn.push(item);
      }
    });
    
    this.setData({ leftColumn, rightColumn });
  },

  // 图片加载完成
  onImageLoad(e) {
    const id = e.currentTarget.dataset.id;
    const updateField = (column) => {
      const idx = column.findIndex(item => item.id === id);
      if (idx !== -1) {
        column[idx].loaded = true;
        return true;
      }
      return false;
    };
    
    const leftColumn = [...this.data.leftColumn];
    const rightColumn = [...this.data.rightColumn];
    
    if (updateField(leftColumn)) {
      this.setData({ leftColumn });
    } else if (updateField(rightColumn)) {
      this.setData({ rightColumn });
    }
  },

  // 预览图片
  previewImage(e) {
    const url = e.currentTarget.dataset.url;
    const urls = this.data.galleryList.map(item => item.image_url).filter(Boolean);

    wx.previewImage({
      current: url,
      urls: urls
    });
  }
});
