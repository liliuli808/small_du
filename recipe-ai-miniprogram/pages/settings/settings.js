const app = getApp()

Page({
  data: {
    apiBaseURL: app.globalData.apiBaseURL,
    envName: 'development',
  },

  onApiURLChange(e) {
    this.setData({ apiBaseURL: e.detail.value })
  },

  saveApiURL() {
    const url = this.data.apiBaseURL.trim()
    if (!url) {
      wx.showToast({ title: '请输入API地址', icon: 'none' })
      return
    }
    app.globalData.apiBaseURL = url
    wx.setStorageSync('api_base_url', url)
    wx.showToast({ title: '已保存', icon: 'success' })
  },

  clearCache() {
    wx.showModal({
      title: '清除缓存',
      content: '将清除本地历史记录和登录状态，确定继续？',
      success: (res) => {
        if (res.confirm) {
          wx.clearStorageSync()
          app.globalData.userOpenID = ''
          wx.showToast({ title: '已清除', icon: 'success' })
        }
      },
    })
  },
})
