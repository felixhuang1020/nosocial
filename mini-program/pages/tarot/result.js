import { get } from '../../utils/request';

Page({
  data: {
    readingId: 0,
    card: {},
    recommend: {},
    isReversed: false,
    meaning: ''
  },

  onLoad(options) {
    const id = options.id || 0;
    this.setData({ readingId: id });

    // 尝试从 eventChannel 获取数据
    const eventChannel = this.getOpenerEventChannel();
    if (eventChannel && eventChannel.on) {
      eventChannel.on('readingData', (data) => {
        if (data) {
          this.processResult(data);
        }
      });
    }

    // 如果没有 eventChannel 数据，尝试从 API 获取
    if (id) {
      this.fetchReadingDetail(id);
    }
  },

  // 获取占卜详情
  fetchReadingDetail(id) {
    get('/wx/tarot/history/' + id).then(res => {
      if (res) {
        this.processHistoryItem(res);
      }
    }).catch(err => {
      console.error('获取占卜详情失败:', err);
    });
  },

  // 处理 eventChannel 传入的完整数据
  processResult(data) {
    const card = data.cards ? data.cards[0] : (data.card || {});
    const recommend = data.recommend || {};
    const isReversed = card.is_reversed || false;

    // 根据正逆位选择含义
    let meaning = card.meaning || '暂无解读';

    this.setData({
      card: card,
      recommend: recommend,
      isReversed: isReversed,
      meaning: meaning
    });
  },

  // 处理历史记录数据（从历史列表 API 返回）
  processHistoryItem(item) {
    let cardInfo = {};
    try {
      const cards = JSON.parse(item.cards || '[]');
      cardInfo = cards[0] || {};
    } catch (e) {
      cardInfo = {};
    }

    const recommend = {
      drink_id: item.recommended_drink_id,
      reason: item.recommended_reason || ''
    };

    this.setData({
      card: {
        name: `塔罗牌 #${cardInfo.card_no !== undefined ? cardInfo.card_no + 1 : '?'}`,
        image_url: '',
        is_reversed: cardInfo.is_reversed || false
      },
      recommend: recommend,
      isReversed: cardInfo.is_reversed || false,
      meaning: item.recommended_reason || '暂无解读'
    });
  },

  // 导航到酒水详情
  navigateToDrink() {
    const drinkId = this.data.recommend && this.data.recommend.drink_id;
    if (drinkId) {
      wx.navigateTo({ url: `/pages/drink/detail?id=${drinkId}` });
    } else {
      wx.showToast({ title: '暂无推荐酒水', icon: 'none' });
    }
  },

  // 再占一次
  divineAgain() {
    wx.redirectTo({ url: '/pages/tarot/divine' });
  },
});
