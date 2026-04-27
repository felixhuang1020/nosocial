import { post, get } from '../../utils/request';

Page({
  data: {
    step: 'question', // question, shuffle, draw
    question: '',
    shuffling: false,
    canDraw: false,
    shuffleText: '正在洗牌...',
    deckCards: Array.from({ length: 5 }, (_, i) => i),
    drawnCard: {},
    recommend: {},
    isReversed: false,
    revealing: false,
    readingId: 0
  },

  onLoad() {
    // 页面加载
  },

  onShow() {
    // 重置状态
    this.setData({
      step: 'question',
      question: '',
      shuffling: false,
      canDraw: false,
      drawnCard: {},
      recommend: {},
      isReversed: false,
      revealing: false
    });
  },

  // 输入问题
  onQuestionInput(e) {
    this.setData({ question: e.detail.value });
  },

  // 开始占卜
  startDivination() {
    this.setData({ step: 'shuffle', shuffleText: '正在洗牌...' }, () => {
      // 开始洗牌动画
      this.setData({ shuffling: true });

      // 3秒后允许抽牌
      setTimeout(() => {
        this.setData({
          shuffling: false,
          canDraw: true,
          shuffleText: '请点击抽牌'
        });
      }, 3000);
    });
  },

  // 抽牌
  drawCard() {
    if (!this.data.canDraw) return;
    
    this.setData({ canDraw: false });
    
    // 调用后端占卜接口
    post('/wx/tarot/divine', {
      question: this.data.question || '今天适合喝什么？',
      spread_type: 1
    }).then(res => {
      this.showCardReveal(res);
    }).catch(err => {
      console.error('占卜请求失败:', err);
      wx.showToast({ title: '占卜服务暂时不可用', icon: 'none' });
      this.setData({ canDraw: true });
    });
  },

  // 展示抽牌动画
  showCardReveal(readingData) {
    const card = readingData.cards ? readingData.cards[0] : readingData.card;
    const isReversed = card ? card.is_reversed : false;
    
    this.setData({
      step: 'draw',
      drawnCard: card || readingData,
      recommend: readingData.recommend || {},
      isReversed: isReversed,
      readingId: readingData.reading_id || 0
    }, () => {
      // 延迟展示翻牌动画
      setTimeout(() => {
        this.setData({ revealing: true });
      }, 500);
    });
  },

  // 查看结果
  viewResult() {
    const readingId = this.data.readingId;
    const drawnCard = this.data.drawnCard || {};
    const resultData = {
      card: drawnCard,
      recommend: this.data.recommend || {},
      reading_id: readingId
    };
    wx.navigateTo({
      url: `/pages/tarot/result?id=${readingId}`,
      success: (res) => {
        res.eventChannel.emit('readingData', resultData);
      }
    });
  },

  // 导航到历史记录
  navigateToHistory() {
    wx.navigateTo({ url: '/pages/tarot/history' });
  },
});
