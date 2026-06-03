const app = getApp()

Page({
  data: {
    recipeCount: 0,
    favoriteCount: 0,
    loading: false,
  },

  onShow() {
    this.loadStats()
  },

  loadStats() {
    this.setData({ loading: true })

    // 获取用户菜谱数量
    wx.request({
      url: `${app.globalData.apiBaseURL}/user/recipes?limit=1`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200 && res.data.recipes) {
          // 接口没有返回总数，需要通过列表长度估算或使用单独统计接口
          // 这里简化处理：实际数量需要后端支持 count 接口
          // 暂时用列表长度
        }
      },
      complete: () => {
        this.setData({ loading: false })
      },
    })

    // 获取收藏数量
    wx.request({
      url: `${app.globalData.apiBaseURL}/user/favorites?limit=1`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200 && res.data.favorites) {
          this.setData({ favoriteCount: res.data.favorites.length })
        }
      },
    })

    // 获取菜谱数量（通过列表）
    wx.request({
      url: `${app.globalData.apiBaseURL}/user/recipes?limit=100`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200 && res.data.recipes) {
          this.setData({ recipeCount: res.data.recipes.length })
        }
      },
    })
  },

  goToUserRecipes() {
    wx.navigateTo({
      url: '/pages/user-recipes/user-recipes',
    })
  },

  goToFavorites() {
    wx.navigateTo({
      url: '/pages/favorites/favorites',
    })
  },

  goBack() {
    wx.navigateBack()
  },
})
