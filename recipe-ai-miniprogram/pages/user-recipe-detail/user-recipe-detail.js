const app = getApp()

Page({
  data: {
    id: null,
    loading: true,
    recipe: {},
  },

  onLoad(options) {
    const id = options.id
    if (!id) {
      wx.showToast({ title: 'ID无效', icon: 'none' })
      wx.navigateBack()
      return
    }
    this.setData({ id })
    this.loadRecipe(id)
  },

  loadRecipe(id) {
    this.setData({ loading: true })
    wx.request({
      url: `${app.globalData.apiBaseURL}/user/recipes/${id}`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200) {
          this.setData({
            recipe: res.data,
            loading: false,
          })
        } else {
          this.setData({ loading: false })
          wx.showToast({ title: '加载失败', icon: 'none' })
        }
      },
      fail: () => {
        this.setData({ loading: false })
        wx.showToast({ title: '网络错误', icon: 'none' })
      },
    })
  },

  goToEdit() {
    wx.navigateTo({
      url: `/pages/user-recipe-edit/user-recipe-edit?id=${this.data.id}`,
    })
  },

  deleteRecipe() {
    wx.showModal({
      title: '确认删除',
      content: '删除后不可恢复，确定要删除这个菜谱吗？',
      confirmColor: '#a83836',
      success: (res) => {
        if (res.confirm) {
          wx.request({
            url: `${app.globalData.apiBaseURL}/user/recipes/${this.data.id}`,
            method: 'DELETE',
            success: (res2) => {
              if (res2.statusCode === 200) {
                wx.showToast({ title: '已删除', icon: 'success' })
                setTimeout(() => wx.navigateBack(), 800)
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