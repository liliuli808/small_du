const app = getApp()

Page({
  data: {
    videoUrl: '',
    loading: false,
    exampleUrls: [
      'https://www.bilibili.com/video/BV1xx411c7mD',
      'https://b23.tv/xxxxxx',
    ],
  },

  onLoad() {
    // 页面加载
  },

  onInputChange(e) {
    this.setData({
      videoUrl: e.detail.value,
    })
  },

  pasteUrl() {
    wx.getClipboardData({
      success: (res) => {
        const text = res.data.trim()
        if (this.isValidBilibiliUrl(text)) {
          this.setData({ videoUrl: text })
          wx.showToast({
            title: '已粘贴',
            icon: 'success',
          })
        } else {
          wx.showToast({
            title: '剪贴板无有效链接',
            icon: 'none',
          })
        }
      },
    })
  },

  clearUrl() {
    this.setData({ videoUrl: '' })
  },

  fillExample(e) {
    const url = e.currentTarget.dataset.url
    this.setData({ videoUrl: url })
  },

  submitAnalyze() {
    const { videoUrl } = this.data
    const trimmed = videoUrl.trim()

    if (!trimmed) {
      wx.showToast({
        title: '请先粘贴B站视频链接',
        icon: 'none',
      })
      return
    }

    if (!this.isValidBilibiliUrl(trimmed)) {
      wx.showToast({
        title: '暂未识别到有效的B站视频链接',
        icon: 'none',
      })
      return
    }

    this.setData({ loading: true })

    wx.request({
      url: `${app.globalData.apiBaseURL}/analyze/bilibili`,
      method: 'POST',
      header: {
        'Content-Type': 'application/json',
      },
      data: { url: trimmed },
      success: (res) => {
        this.setData({ loading: false })

        if (res.statusCode === 200 && res.data.task_id) {
          // 同视频去重：如果已有结果，直接跳转到结果页
          if (res.data.is_duplicate && res.data.recipe_id) {
            wx.showToast({
              title: '该视频已解析过',
              icon: 'success',
            })
            wx.redirectTo({
              url: `/pages/result/result?recipe_id=${res.data.recipe_id}`,
            })
          } else {
            wx.navigateTo({
              url: `/pages/task/task?task_id=${res.data.task_id}`,
            })
          }
        } else {
          const msg = res.data?.message || '任务创建失败'
          wx.showToast({ title: msg, icon: 'none' })
        }
      },
      fail: () => {
        this.setData({ loading: false })
        wx.showToast({
          title: '网络异常，请稍后重试',
          icon: 'none',
        })
      },
    })
  },

  isValidBilibiliUrl(url) {
    const lower = url.toLowerCase()
    return lower.includes('bilibili.com') ||
           lower.includes('b23.tv') ||
           lower.includes('bili.im') ||
           lower.includes('bv')
  },
})
