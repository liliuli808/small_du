const app = getApp()

Page({
  data: {
    recipeId: null,
    loading: true,
    error: '',
    video: {},
    textSources: {},
    recipe: {},
    nutrition: {},
    isFavorited: false,
    favoriteLoading: false,
  },

  onLoad(options) {
    const recipeId = options.recipe_id
    if (!recipeId) {
      wx.showToast({ title: '菜谱ID无效', icon: 'none' })
      wx.redirectTo({ url: '/pages/index/index' })
      return
    }

    this.setData({ recipeId })
    this.loadRecipe(recipeId)
    this.checkFavoriteStatus(recipeId)
  },

  loadRecipe(recipeId) {
    wx.request({
      url: `${app.globalData.apiBaseURL}/recipes/${recipeId}`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200) {
          this.setData({
            video: res.data.video || {},
            textSources: res.data.text_sources || {},
            recipe: res.data.recipe || {},
            nutrition: res.data.nutrition || {},
            loading: false,
          })
          // 记录到本地解析历史
          app.addAnalyzeHistory({
            recipe_id: recipeId,
            dish_name: res.data.recipe?.dish_name,
            video_title: res.data.video?.title,
            total_kcal: res.data.nutrition?.total_kcal,
          })
        } else {
          this.setData({
            error: res.data?.message || '获取菜谱失败',
            loading: false,
          })
        }
      },
      fail: () => {
        this.setData({
          error: '网络错误，请稍后重试',
          loading: false,
        })
      },
    })
  },

  checkFavoriteStatus(recipeId) {
    wx.request({
      url: `${app.globalData.apiBaseURL}/recipes/${recipeId}/favorite`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200) {
          this.setData({ isFavorited: res.data.is_favorited })
        }
      },
    })
  },

  toggleFavorite() {
    if (this.data.favoriteLoading) return

    this.setData({ favoriteLoading: true })
    wx.request({
      url: `${app.globalData.apiBaseURL}/recipes/${this.data.recipeId}/favorite`,
      method: 'POST',
      success: (res) => {
        this.setData({ favoriteLoading: false })
        if (res.statusCode === 200) {
          this.setData({ isFavorited: res.data.is_favorited })
          wx.showToast({
            title: res.data.is_favorited ? '已收藏' : '已取消收藏',
            icon: 'none',
          })
        } else if (res.statusCode === 401) {
          wx.showToast({
            title: '请先登录',
            icon: 'none',
          })
        } else {
          wx.showToast({
            title: res.data?.message || '操作失败',
            icon: 'none',
          })
        }
      },
      fail: () => {
        this.setData({ favoriteLoading: false })
        wx.showToast({
          title: '网络错误',
          icon: 'none',
        })
      },
    })
  },

  goBack() {
    wx.redirectTo({ url: '/pages/index/index' })
  },

  goEditIngredients() {
    const ingredients = JSON.stringify(this.data.recipe.ingredients || [])
    const servings = this.data.recipe.servings || 1
    wx.navigateTo({
      url: `/pages/edit-ingredients/edit-ingredients?recipe_id=${this.data.recipeId}&ingredients=${encodeURIComponent(ingredients)}&servings=${servings}`,
    })
  },

  formatTime(seconds) {
    if (seconds === null || seconds === undefined) return ''
    const m = Math.floor(seconds / 60)
    const s = seconds % 60
    return `${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`
  },

  formatConfidence(conf) {
    if (conf >= 0.8) return { text: '高', class: 'high' }
    if (conf >= 0.5) return { text: '中', class: 'medium' }
    return { text: '低', class: 'low' }
  },
})
