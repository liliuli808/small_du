App({
  globalData: {
    apiBaseURL: 'http://localhost:8080/api/v1',
    userOpenID: '',
  },

  onLaunch() {
    // 从缓存恢复API地址
    const savedURL = wx.getStorageSync('api_base_url')
    if (savedURL) {
      this.globalData.apiBaseURL = savedURL
    }

    // 重写 wx.request，自动带上用户标识
    this.wrapRequest()

    // 尝试微信登录，失败则使用匿名ID
    this.doWxLogin()
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

  // ===== 微信登录 =====
  loginKey: 'wx_login_openid',

  doWxLogin() {
    return new Promise((resolve) => {
      const cached = wx.getStorageSync(this.loginKey)
      if (cached) {
        this.globalData.userOpenID = cached
        resolve(cached)
        return
      }

      wx.login({
        success: (res) => {
          if (res.code) {
            wx.request({
              url: `${this.globalData.apiBaseURL}/auth/wx-login`,
              method: 'POST',
              data: { code: res.code },
              success: (r) => {
                if (r.statusCode === 200 && r.data.openid) {
                  const openid = r.data.openid
                  wx.setStorageSync(this.loginKey, openid)
                  this.globalData.userOpenID = openid
                  resolve(openid)
                  return
                }
              },
              complete: () => {
                // 登录接口失败时使用本地匿名ID
                if (!this.globalData.userOpenID) {
                  this.ensureAnonymousID()
                }
                resolve(this.globalData.userOpenID)
              },
            })
          } else {
            this.ensureAnonymousID()
            resolve(this.globalData.userOpenID)
          }
        },
        fail: () => {
          this.ensureAnonymousID()
          resolve(this.globalData.userOpenID)
        },
      })
    })
  },

  ensureAnonymousID() {
    if (this.globalData.userOpenID) return
    let openid = wx.getStorageSync('user_openid')
    if (!openid) {
      openid = this.generateOpenID()
      wx.setStorageSync('user_openid', openid)
    }
    this.globalData.userOpenID = openid
  },

  // ===== 解析历史（本地存储） =====
  historyKey: 'analyze_history',
  historyMax: 50,

  getAnalyzeHistory() {
    return wx.getStorageSync(this.historyKey) || []
  },

  addAnalyzeHistory(item) {
    if (!item || !item.recipe_id) return
    const list = this.getAnalyzeHistory()
    // 去重：同一菜谱只保留最新一条
    const filtered = list.filter((it) => it.recipe_id !== item.recipe_id)
    filtered.unshift({
      recipe_id: item.recipe_id,
      dish_name: item.dish_name || '未命名菜谱',
      video_title: item.video_title || '',
      total_kcal: item.total_kcal || 0,
      created_at: Date.now(),
    })
    wx.setStorageSync(this.historyKey, filtered.slice(0, this.historyMax))
  },

  removeAnalyzeHistory(recipeId) {
    const list = this.getAnalyzeHistory().filter((it) => it.recipe_id !== recipeId)
    wx.setStorageSync(this.historyKey, list)
    return list
  },

  clearAnalyzeHistory() {
    wx.removeStorageSync(this.historyKey)
  },
})
