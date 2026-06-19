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
