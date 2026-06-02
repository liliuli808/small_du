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
