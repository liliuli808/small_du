const app = getApp()

const SERVINGS_RANGE = [1, 2, 3, 4, 5, 6, 8, 10]
const UNIT_RANGE = ['克', '个', '勺', '毫升', '根', '片', '块', '适量']

Page({
  data: {
    id: null,
    dishName: '',
    servings: 2,
    servingsIndex: 1,
    ingredients: [
      { name: '', amount: null, unit: '克', unitIndex: 0 },
    ],
    steps: [
      { title: '', description: '' },
    ],
    tips: '',
    loading: false,
    isEdit: false,
    servingsRange: SERVINGS_RANGE,
    unitRange: UNIT_RANGE,
  },

  onLoad(options) {
    if (options.id) {
      this.setData({ id: parseInt(options.id), isEdit: true })
      this.loadRecipe(options.id)
    } else if (options.from_recipe_id) {
      this.setData({ fromRecipeId: options.from_recipe_id })
      this.loadAIDerivedRecipe(options.from_recipe_id)
    }
  },

  loadAIDerivedRecipe(recipeId) {
    this.setData({ loading: true })
    wx.request({
      url: `${app.globalData.apiBaseURL}/recipes/${recipeId}/derive`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200) {
          const data = res.data
          const ingredients = (data.ingredients || []).map(ing => {
            const unitIdx = UNIT_RANGE.indexOf(ing.unit)
            return {
              name: ing.name || '',
              amount: ing.amount || ing.grams || null,
              unit: ing.unit || '克',
              unitIndex: unitIdx >= 0 ? unitIdx : 0,
            }
          })
          const steps = (data.steps || []).map(step => ({
            title: step.title || '',
            description: step.description || '',
          }))
          const servingsIdx = SERVINGS_RANGE.indexOf(data.servings)
          this.setData({
            dishName: data.dish_name || '',
            servings: data.servings || 2,
            servingsIndex: servingsIdx >= 0 ? servingsIdx : 1,
            ingredients: ingredients.length > 0 ? ingredients : [{ name: '', amount: null, unit: '克', unitIndex: 0 }],
            steps: steps.length > 0 ? steps : [{ title: '', description: '' }],
            tips: (data.tips || []).join('\n'),
            loading: false,
          })
        } else {
          this.setData({ loading: false })
          wx.showToast({ title: '加载菜谱失败', icon: 'none' })
        }
      },
      fail: () => {
        this.setData({ loading: false })
        wx.showToast({ title: '网络错误', icon: 'none' })
      },
    })
  },

  loadRecipe(id) {
    this.setData({ loading: true })
    wx.request({
      url: `${app.globalData.apiBaseURL}/user/recipes/${id}`,
      method: 'GET',
      success: (res) => {
        if (res.statusCode === 200) {
          const data = res.data
          const ingredients = (data.recipe.ingredients || []).map(ing => {
            const unitIdx = UNIT_RANGE.indexOf(ing.unit)
            return {
              name: ing.name || '',
              amount: ing.amount || ing.grams || null,
              unit: ing.unit || '克',
              unitIndex: unitIdx >= 0 ? unitIdx : 0,
            }
          })
          const steps = (data.recipe.steps || []).map(step => ({
            title: step.title || '',
            description: step.description || '',
          }))
          const servingsIdx = SERVINGS_RANGE.indexOf(data.servings)
          this.setData({
            dishName: data.dish_name || '',
            servings: data.servings || 2,
            servingsIndex: servingsIdx >= 0 ? servingsIdx : 1,
            ingredients: ingredients.length > 0 ? ingredients : [{ name: '', amount: null, unit: '克', unitIndex: 0 }],
            steps: steps.length > 0 ? steps : [{ title: '', description: '' }],
            tips: (data.recipe.tips || []).join('\n'),
            loading: false,
          })
        } else {
          this.setData({ loading: false })
          wx.showToast({ title: '加载失败', icon: 'none' })
        }
      },
      fail: () => {
        this.setData({ loading: false })
        wx.showToast({ title: '加载失败', icon: 'none' })
      },
    })
  },

  onDishNameChange(e) {
    this.setData({ dishName: e.detail.value })
  },

  onServingsChange(e) {
    const idx = e.detail.value
    this.setData({
      servings: SERVINGS_RANGE[idx],
      servingsIndex: idx,
    })
  },

  onIngredientNameChange(e) {
    const index = e.currentTarget.dataset.index
    const key = `ingredients[${index}].name`
    this.setData({ [key]: e.detail.value })
  },

  onIngredientAmountChange(e) {
    const index = e.currentTarget.dataset.index
    const key = `ingredients[${index}].amount`
    const val = e.detail.value
    // 允许空字符串（用户正在输入），转换为数字时处理
    this.setData({ [key]: val === '' ? null : parseFloat(val) })
  },

  onIngredientUnitChange(e) {
    const index = e.currentTarget.dataset.index
    const idx = e.detail.value
    this.setData({
      [`ingredients[${index}].unit`]: UNIT_RANGE[idx],
      [`ingredients[${index}].unitIndex`]: idx,
    })
  },

  addIngredient() {
    const ingredients = this.data.ingredients
    ingredients.push({ name: '', amount: null, unit: '克', unitIndex: 0 })
    this.setData({ ingredients })
  },

  removeIngredient(e) {
    const index = e.currentTarget.dataset.index
    const ingredients = this.data.ingredients
    if (ingredients.length <= 1) {
      wx.showToast({ title: '至少保留一个食材', icon: 'none' })
      return
    }
    ingredients.splice(index, 1)
    this.setData({ ingredients })
  },

  onStepTitleChange(e) {
    const index = e.currentTarget.dataset.index
    const key = `steps[${index}].title`
    this.setData({ [key]: e.detail.value })
  },

  onStepDescChange(e) {
    const index = e.currentTarget.dataset.index
    const key = `steps[${index}].description`
    this.setData({ [key]: e.detail.value })
  },

  addStep() {
    const steps = this.data.steps
    steps.push({ title: '', description: '' })
    this.setData({ steps })
  },

  removeStep(e) {
    const index = e.currentTarget.dataset.index
    const steps = this.data.steps
    if (steps.length <= 1) {
      wx.showToast({ title: '至少保留一个步骤', icon: 'none' })
      return
    }
    steps.splice(index, 1)
    this.setData({ steps })
  },

  onTipsChange(e) {
    this.setData({ tips: e.detail.value })
  },

  validate() {
    if (!this.data.dishName.trim()) {
      wx.showToast({ title: '请输入菜名', icon: 'none' })
      return false
    }
    for (let i = 0; i < this.data.ingredients.length; i++) {
      if (!this.data.ingredients[i].name.trim()) {
        wx.showToast({ title: `食材${i + 1}名称不能为空`, icon: 'none' })
        return false
      }
    }
    for (let i = 0; i < this.data.steps.length; i++) {
      if (!this.data.steps[i].description.trim()) {
        wx.showToast({ title: `步骤${i + 1}描述不能为空`, icon: 'none' })
        return false
      }
    }
    return true
  },

  submit() {
    if (!this.validate()) return

    const tipsArr = this.data.tips.split('\n').filter(t => t.trim())

    const payload = {
      dish_name: this.data.dishName.trim(),
      servings: this.data.servings,
      ingredients: this.data.ingredients.map(ing => ({
        name: ing.name.trim(),
        amount: ing.amount || 0,
        unit: ing.unit,
        grams: ing.amount || 0,
      })),
      steps: this.data.steps.map((step, idx) => ({
        order: idx + 1,
        title: step.title.trim() || `步骤${idx + 1}`,
        description: step.description.trim(),
      })),
      tips: tipsArr,
    }

    this.setData({ loading: true })

    if (this.data.isEdit) {
      wx.request({
        url: `${app.globalData.apiBaseURL}/user/recipes/${this.data.id}`,
        method: 'PUT',
        header: { 'Content-Type': 'application/json' },
        data: payload,
        success: (res) => {
          this.setData({ loading: false })
          if (res.statusCode === 200) {
            wx.showToast({ title: '更新成功', icon: 'success' })
            setTimeout(() => wx.navigateBack(), 800)
          } else {
            wx.showToast({ title: res.data?.message || '更新失败', icon: 'none' })
          }
        },
        fail: () => {
          this.setData({ loading: false })
          wx.showToast({ title: '网络错误', icon: 'none' })
        },
      })
    } else {
      wx.request({
        url: `${app.globalData.apiBaseURL}/user/recipes`,
        method: 'POST',
        header: { 'Content-Type': 'application/json' },
        data: payload,
        success: (res) => {
          this.setData({ loading: false })
          if (res.statusCode === 200) {
            wx.showToast({ title: '创建成功', icon: 'success' })
            setTimeout(() => wx.navigateBack(), 800)
          } else {
            wx.showToast({ title: res.data?.message || '创建失败', icon: 'none' })
          }
        },
        fail: () => {
          this.setData({ loading: false })
          wx.showToast({ title: '网络错误', icon: 'none' })
        },
      })
    }
  },

  goBack() {
    wx.navigateBack()
  },
})
