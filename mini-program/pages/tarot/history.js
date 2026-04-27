import { get } from '../../utils/request';

Page({
  data: {
    historyList: [],
    page: 1,
    size: 10,
    loading: false,
    hasMore: true
  },

  onLoad() {
    this.fetchHistory();
  },

  onShow() {
    // onLoad 已加载，onShow 仅在返回时刷新（historyList 非空时才触发）
    if (this.data.historyList.length > 0) {
      this.fetchHistory(true);
    }
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.setData({ page: 1, hasMore: true }, () => {
      this.fetchHistory(true).finally(() => {
        wx.stopPullDownRefresh();
      });
    });
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.fetchHistory();
    }
  },

  // 获取历史记录
  fetchHistory(reset = false) {
    if (this.data.loading) return Promise.resolve();

    this.setData({ loading: true });
    const page = reset ? 1 : this.data.page;

    return get('/wx/tarot/history', { page, size: this.data.size }).then(res => {
      const list = res && res.list ? res.list : [];
      const hasMore = list.length === this.data.size;

      // 解析后端数据，适配前端展示
      const processedList = list.map(item => {
        let cardInfo = {};
        try {
          const cards = JSON.parse(item.cards || '[]');
          cardInfo = cards[0] || {};
        } catch (e) {
          cardInfo = {};
        }
        return {
          ...item,
          card_no: cardInfo.card_no,
          is_reversed: cardInfo.is_reversed || false,
          displayDate: item.reading_date || (item.created_at ? item.created_at.split('T')[0] : '')
        };
      });

      this.setData({
        historyList: reset ? processedList : [...this.data.historyList, ...processedList],
        page: page + 1,
        hasMore: hasMore,
        loading: false
      });
    }).catch(err => {
      console.error('获取占卜历史失败:', err);
      this.setData({
        historyList: reset ? [] : this.data.historyList,
        hasMore: false,
        loading: false
      });
    });
  },

  // 查看详情
  viewDetail(e) {
    const id = e.currentTarget.dataset.id;
    wx.navigateTo({ url: `/pages/tarot/result?id=${id}` });
  },

  // 去占卜
  goDivine() {
    wx.navigateTo({ url: '/pages/tarot/divine' });
  },
});
