const app = getApp()

Page({
  data: {
    recipes: [],
    loading: false,
  },

  onShow() {
    this.loadRecipes()
  },

  loadRecipes() {
    this.setData({ loading: true })
    wx.request({
      url: `${app.globalData.apiBaseURL}/user/recipes?limit=100`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200 && res.data.recipes) {
          this.setData({
            recipes: res.data.recipes,
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

  goToCreate() {
    wx.navigateTo({
      url: '/pages/user-recipe-edit/user-recipe-edit',
    })
  },

  goToDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/user-recipe-detail/user-recipe-detail?id=${id}`,
    })
  },

  goToEdit(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/user-recipe-edit/user-recipe-edit?id=${id}`,
    })
  },

  deleteRecipe(e) {
    const id = e.currentTarget.dataset.id
    wx.showModal({
      title: '确认删除',
      content: '删除后不可恢复，确定要删除这个菜谱吗？',
      confirmColor: '#a83836',
      success: (res) => {
        if (res.confirm) {
          wx.request({
            url: `${app.globalData.apiBaseURL}/user/recipes/${id}`,
            method: 'DELETE',
            success: (res2) => {
              if (res2.statusCode === 200) {
                wx.showToast({ title: '已删除', icon: 'success' })
                this.loadRecipes()
              } else {
                wx.showToast({ title: '删除失败', icon: 'none' })
              }
            },
            fail: () => {
              wx.showToast({ title: '删除失败', icon: 'none' })
            },
          })
        }
      },
    })
  },

  goBack() {
    wx.navigateBack()
  },
})