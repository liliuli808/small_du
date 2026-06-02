const app = getApp()

Page({
  data: {
    recipeId: null,
    servings: 1,
    ingredients: [],
    originalServings: 1,
  },

  onLoad(options) {
    const recipeId = options.recipe_id
    const servings = parseInt(options.servings) || 1
    let ingredients = []

    try {
      ingredients = JSON.parse(decodeURIComponent(options.ingredients || '[]'))
    } catch (e) {
      console.error('解析材料失败', e)
    }

    // 转换格式
    const formattedIngredients = ingredients.map(item => ({
      name: item.name,
      grams: item.grams || 0,
      originalGrams: item.grams || 0,
      changeText: '',
    }))

    this.setData({
      recipeId,
      servings,
      originalServings: servings,
      ingredients: formattedIngredients,
    })
  },

  // 增加份数
  increaseServings() {
    this._updateServings(this.data.servings + 1)
  },

  // 减少份数
  decreaseServings() {
    if (this.data.servings <= 1) return
    this._updateServings(this.data.servings - 1)
  },

  // 内部更新份数方法
  _updateServings(newServings) {
    if (newServings < 1) return

    // 按比例调整所有材料的克重
    const ratio = newServings / this.data.originalServings
    const ingredients = this.data.ingredients.map(item => {
      const newGrams = Math.round(item.originalGrams * ratio * 10) / 10
      const diff = Math.round((newGrams - item.originalGrams) * 10) / 10
      return {
        ...item,
        grams: newGrams,
        changeText: diff > 0 ? `+${diff}g` : diff < 0 ? `${diff}g` : '',
      }
    })

    this.setData({
      servings: newServings,
      ingredients,
    })
  },

  onServingsChange(e) {
    const value = e.detail ? e.detail.value : (e.currentTarget ? e.currentTarget.dataset.servings : null)
    const newServings = parseInt(value) || 1
    this._updateServings(newServings)
  },

  onGramsChange(e) {
    const index = e.currentTarget.dataset.index
    const grams = parseFloat(e.detail.value)
    if (grams < 0 || isNaN(grams)) return

    const ingredients = this.data.ingredients.slice()
    ingredients[index].grams = grams

    // 更新变化量显示
    const diff = Math.round((grams - ingredients[index].originalGrams) * 10) / 10
    ingredients[index].changeText = diff > 0 ? `+${diff}g` : diff < 0 ? `${diff}g` : ''

    this.setData({ ingredients })
  },

  submitRecalculate() {
    const { recipeId, servings, ingredients } = this.data

    const payload = {
      servings,
      ingredients: ingredients.map(item => ({
        name: item.name,
        grams: item.grams,
      })),
    }

    wx.showLoading({ title: '计算中...' })

    wx.request({
      url: `${app.globalData.apiBaseURL}/recipes/${recipeId}/recalculate`,
      method: 'POST',
      header: {
        'Content-Type': 'application/json',
      },
      data: payload,
      success: (res) => {
        wx.hideLoading()

        if (res.statusCode === 200) {
          // 将新的营养数据传回结果页
          const pages = getCurrentPages()
          const resultPage = pages[pages.length - 2]
          if (resultPage) {
            resultPage.setData({
              'recipe.servings': servings,
              'nutrition': res.data.nutrition,
            })
            // 更新显示的克重
            const updatedIngredients = this.data.ingredients.map((item, idx) => ({
              ...resultPage.data.recipe.ingredients[idx],
              grams: item.grams,
            }))
            resultPage.setData({
              'recipe.ingredients': updatedIngredients,
            })
          }

          wx.navigateBack()
          wx.showToast({
            title: '已重新计算',
            icon: 'success',
          })
        } else {
          wx.showToast({
            title: res.data?.message || '计算失败',
            icon: 'none',
          })
        }
      },
      fail: () => {
        wx.hideLoading()
        wx.showToast({
          title: '网络错误',
          icon: 'none',
        })
      },
    })
  },

  resetToOriginal() {
    const ingredients = this.data.ingredients.map(item => ({
      ...item,
      grams: item.originalGrams,
      changeText: '',
    }))
    this.setData({
      servings: this.data.originalServings,
      ingredients,
    })
  },
})
