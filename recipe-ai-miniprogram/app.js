App({
  globalData: {
    apiBaseURL: 'http://localhost:8080/api/v1',
    userOpenID: '',
  },

  onLaunch() {
    // 获取或生成用户标识
    let openid = wx.getStorageSync('user_openid')
    if (!openid) {
      openid = this.generateOpenID()
      wx.setStorageSync('user_openid', openid)
    }
    this.globalData.userOpenID = openid

    // 重写 wx.request，自动带上用户标识
    this.wrapRequest()
  },

  generateOpenID() {
    // 生成简易用户标识：时间戳 + 随机数
    const ts = Date.now().toString(36)
    const rand = Math.random().toString(36).substring(2, 10)
    return `u_${ts}_${rand}`
  },

  wrapRequest() {
    const originalRequest = wx.request
    const app = this
    wx.request = function(options) {
      if (!options.header) {
        options.header = {}
      }
      if (app.globalData.userOpenID) {
        options.header['X-User-OpenID'] = app.globalData.userOpenID
      }
      return originalRequest(options)
    }
  },
})
