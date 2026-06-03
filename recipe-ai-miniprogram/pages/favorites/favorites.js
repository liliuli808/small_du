const app = getApp()

Page({
  data: {
    favorites: [],
    loading: false,
  },

  onShow() {
    this.loadFavorites()
  },

  loadFavorites() {
    this.setData({ loading: true })
    wx.request({
      url: `${app.globalData.apiBaseURL}/user/favorites?limit=50`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200 && res.data.favorites) {
          this.setData({
            favorites: res.data.favorites,
            loading: false,
          })
        } else {
          this.setData({ loading: false })
        }
      },
      fail: () => {
        this.setData({ loading: false })
        wx.showToast({ title: '加载失败', icon: 'none' })
      },
    })
  },

  goToRecipe(e) {
    const recipeId = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/result/result?recipe_id=${recipeId}`,
    })
  },

  unfavorite(e) {
    const recipeId = e.currentTarget.dataset.id
    wx.request({
      url: `${app.globalData.apiBaseURL}/recipes/${recipeId}/favorite`,
      method: 'POST',
      success: (res) => {
        if (res.statusCode === 200) {
          wx.showToast({ title: '已取消收藏', icon: 'none' })
          this.loadFavorites()
        }
      },
    })
  },

  goBack() {
    wx.navigateBack()
  },
})