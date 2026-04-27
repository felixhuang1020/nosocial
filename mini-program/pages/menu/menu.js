import { get } from '../../utils/request';

Page({
  data: {
    categories: [{ id: 0, name: '全部' }],
    currentCategory: 0,
    drinkList: [],
    page: 1,
    size: 10,
    loading: false,
    hasMore: true,
    searchKeyword: ''
  },

  onLoad() {
    this.fetchCategories();
    this.fetchDrinkList();
  },

  onShow() {
    // 页面显示
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.setData({ page: 1, hasMore: true }, () => {
      Promise.all([
        this.fetchCategories(),
        this.fetchDrinkList(true)
      ]).finally(() => {
        wx.stopPullDownRefresh();
      });
    });
  },

  // 上拉加载更多
  onReachBottom() {
    if (this.data.hasMore && !this.data.loading) {
      this.fetchDrinkList();
    }
  },

  // 获取分类列表（从酒水数据中动态提取）
  fetchCategories() {
    return get('/public/drinks', { page: 1, size: 100 }).then(res => {
      const list = res && res.list ? res.list : [];
      if (list.length > 0) {
        // 动态提取分类
        const categoryMap = new Map();
        list.forEach(item => {
          if (item.category && item.category.id && item.category.name) {
            categoryMap.set(item.category.id, { id: item.category.id, name: item.category.name });
          }
        });
        const categories = [{ id: 0, name: '全部' }, ...Array.from(categoryMap.values())];
        this.setData({ categories });
      }
    }).catch(err => {
      console.log('使用默认分类');
    });
  },

  // 获取酒水列表
  fetchDrinkList(reset = false) {
    if (this.data.loading) return Promise.resolve();
    
    this.setData({ loading: true });
    
    const page = reset ? 1 : this.data.page;
    const params = {
      page: page,
      size: this.data.size,
      status: 1
    };
    
    if (this.data.currentCategory > 0) {
      params.category_id = this.data.currentCategory;
    }
    
    if (this.data.searchKeyword) {
      params.keyword = this.data.searchKeyword;
    }
    
    return get('/public/drinks', params).then(res => {
      const list = res && res.list ? res.list : [];
      const hasMore = list.length === this.data.size;
      
      // 处理配料数据
      const processedList = list.map(item => ({
        ...item,
        ingredientList: item.ingredients ? item.ingredients.split(',').slice(0, 3) : []
      }));
      
      this.setData({
        drinkList: reset ? processedList : [...this.data.drinkList, ...processedList],
        page: page + 1,
        hasMore: hasMore,
        loading: false
      });
    }).catch(err => {
      console.error('获取酒水列表失败:', err);
      this.setData({
        drinkList: reset ? [] : this.data.drinkList,
        hasMore: false,
        loading: false
      });
    });
  },

  // 切换分类
  switchCategory(e) {
    const id = parseInt(e.currentTarget.dataset.id);
    if (id === this.data.currentCategory) return;
    
    this.setData({
      currentCategory: id,
      drinkList: [],
      page: 1,
      hasMore: true
    }, () => {
      this.fetchDrinkList(true);
    });
  },

  // 显示搜索
  showSearch() {
    wx.showModal({
      title: '搜索酒水',
      editable: true,
      placeholderText: '输入酒水名称...',
      confirmColor: '#7C9A92',
      success: (res) => {
        if (res.confirm && res.content) {
          this.setData({
            searchKeyword: res.content,
            drinkList: [],
            page: 1,
            hasMore: true
          }, () => {
            this.fetchDrinkList(true);
          });
        }
      }
    });
  },

  // 重置分类
  resetCategory() {
    this.setData({
      currentCategory: 0,
      searchKeyword: '',
      drinkList: [],
      page: 1,
      hasMore: true
    }, () => {
      this.fetchDrinkList(true);
    });
  },

  // 导航到酒水详情
  navigateToDetail(e) {
    const id = e.currentTarget.dataset.id;
    wx.navigateTo({ url: `/pages/drink/detail?id=${id}` });
  },
});
