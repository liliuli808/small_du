const app = getApp()

const POLLING_INTERVAL = 2000
const MAX_POLLING_TIME = 60000

Page({
  data: {
    taskId: '',
    status: 'pending',
    message: '准备解析',
    recipeId: null,
    elapsedTime: 0,
    pollingCount: 0,
    errorMsg: '',
  },

  timer: null,
  startTime: 0,

  onLoad(options) {
    const taskId = options.task_id
    if (!taskId) {
      wx.showToast({ title: '任务ID无效', icon: 'none' })
      wx.navigateBack()
      return
    }

    this.setData({ taskId })
    this.startTime = Date.now()
    this.startPolling()
  },

  onUnload() {
    this.stopPolling()
  },

  startPolling() {
    // 立即查询一次
    this.queryTaskStatus()

    // 定时轮询
    this.timer = setInterval(() => {
      const elapsed = Date.now() - this.startTime
      this.setData({ elapsedTime: Math.floor(elapsed / 1000) })

      if (elapsed > MAX_POLLING_TIME) {
        this.stopPolling()
        this.setData({
          status: 'timeout',
          message: '解析时间较长，请稍后重新尝试',
        })
        return
      }

      this.queryTaskStatus()
    }, POLLING_INTERVAL)
  },

  stopPolling() {
    if (this.timer) {
      clearInterval(this.timer)
      this.timer = null
    }
  },

  queryTaskStatus() {
    wx.request({
      url: `${app.globalData.apiBaseURL}/tasks/${this.data.taskId}`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode !== 200) {
          return
        }

        const data = res.data
        this.setData({
          status: data.status,
          message: data.message,
          recipeId: data.recipe_id,
          pollingCount: this.data.pollingCount + 1,
        })

        if (data.status === 'done') {
          this.stopPolling()
          if (data.recipe_id) {
            wx.redirectTo({
              url: `/pages/result/result?recipe_id=${data.recipe_id}`,
            })
          }
        } else if (data.status === 'failed') {
          this.stopPolling()
          this.setData({ errorMsg: data.message || '解析失败' })
        }
      },
      fail: () => {
        // 网络错误，继续轮询
      },
    })
  },

  goBack() {
    wx.navigateBack()
  },

  retry() {
    this.setData({
      status: 'pending',
      message: '准备解析',
      errorMsg: '',
      elapsedTime: 0,
    })
    this.startTime = Date.now()
    this.startPolling()
  },
})
