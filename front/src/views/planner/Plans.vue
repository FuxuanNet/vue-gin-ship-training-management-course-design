<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getPlansList, createPlan, updatePlan, deletePlan } from '../../api/planner'

const loading = ref(false)
const dialogVisible = ref(false)
const dialogMode = ref('add') // 'add' or 'edit'
const currentPlan = ref(null)
const plansList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

const planForm = ref({
  planName: '',
  planStatus: '规划中',
  planStartDatetime: '',
  planEndDatetime: ''
})

const statusOptions = [
  { label: '规划中', value: '规划中' },
  { label: '进行中', value: '进行中' },
  { label: '已完成', value: '已完成' }
]

// 获取培训计划列表
const fetchPlansList = async () => {
  loading.value = true
  try {
    const response = await getPlansList({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    if (response.code === 200) {
      plansList.value = response.data.list || []
      total.value = response.data.total || 0
    }
  } catch (error) {
    ElMessage.error('获取培训计划列表失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 分页变化
const handlePageChange = (page) => {
  currentPage.value = page
  fetchPlansList()
}

// 打开新增对话框
const openAddDialog = () => {
  dialogMode.value = 'add'
  planForm.value = {
    planName: '',
    planStatus: '规划中',
    planStartDatetime: '',
    planEndDatetime: ''
  }
  dialogVisible.value = true
}

// 打开编辑对话框
const openEditDialog = (plan) => {
  dialogMode.value = 'edit'
  currentPlan.value = plan
  planForm.value = {
    planName: plan.planName,
    planStatus: plan.planStatus,
    planStartDatetime: plan.planStartDatetime,
    planEndDatetime: plan.planEndDatetime
  }
  dialogVisible.value = true
}

// 提交表单
const submitForm = async () => {
  if (!planForm.value.planName) {
    ElMessage.warning('请输入计划名称')
    return
  }
  if (!planForm.value.planStartDatetime) {
    ElMessage.warning('请选择开始时间')
    return
  }
  if (!planForm.value.planEndDatetime) {
    ElMessage.warning('请选择结束时间')
    return
  }
  
  loading.value = true
  try {
    if (dialogMode.value === 'add') {
      await createPlan(planForm.value)
      ElMessage.success('培训计划创建成功')
    } else {
      await updatePlan(currentPlan.value.planId, planForm.value)
      ElMessage.success('培训计划更新成功')
    }
    
    dialogVisible.value = false
    fetchPlansList()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

// 删除计划
const handleDeletePlan = (plan) => {
  ElMessageBox.confirm(
    `确定要删除培训计划"${plan.planName}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deletePlan(plan.planId)
      ElMessage.success('删除成功')
      fetchPlansList()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {})
}

// 查看计划详情
const viewPlanDetail = (plan) => {
  console.log('查看计划详情:', plan)
  ElMessage.info('详情页面开发中')
}

// 获取状态标签类型
const getStatusType = (status) => {
  const typeMap = {
    '规划中': 'info',
    '进行中': 'primary',
    '已完成': 'success'
  }
  return typeMap[status] || 'info'
}

onMounted(() => {
  fetchPlansList()
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>培训计划管理</h1>
      <el-button type="primary" @click="openAddDialog">
        <el-icon><Plus /></el-icon>
        新建培训计划
      </el-button>
    </div>
    
    <div class="content-wrapper" v-loading="loading">
      <el-table :data="plansList" border stripe>
        <el-table-column prop="planId" label="计划ID" width="80" />
        <el-table-column prop="planName" label="计划名称" min-width="200" />
        <el-table-column label="计划状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.planStatus)">{{ row.planStatus }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="planStartDatetime" label="开始时间" width="180" />
        <el-table-column prop="planEndDatetime" label="结束时间" width="180" />
        <el-table-column prop="creatorName" label="制定人" width="100" />
        <el-table-column label="参与员工" width="100" align="center">
          <template #default="{ row }">
            {{ row.employeeCount }} 人
          </template>
        </el-table-column>
        <el-table-column label="课程数" width="100" align="center">
          <template #default="{ row }">
            {{ row.courseCount }} 门
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" align="center">
          <template #default="{ row }">
            <el-button size="small" @click="viewPlanDetail(row)">详情</el-button>
            <el-button size="small" type="primary" @click="openEditDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDeletePlan(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div style="margin-top: 20px; display: flex; justify-content: flex-end;">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next, jumper"
          @current-change="handlePageChange"
        />
      </div>
    </div>
    
    <!-- 新增/编辑对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="dialogMode === 'add' ? '新建培训计划' : '编辑培训计划'"
      width="600px"
    >
      <el-form :model="planForm" label-width="120px">
        <el-form-item label="计划名称" required>
          <el-input v-model="planForm.planName" placeholder="请输入培训计划名称" />
        </el-form-item>
        
        <el-form-item label="计划状态">
          <el-select v-model="planForm.planStatus" placeholder="请选择计划状态" style="width: 100%;">
            <el-option 
              v-for="status in statusOptions" 
              :key="status.value"
              :label="status.label"
              :value="status.value"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="开始时间" required>
          <el-date-picker
            v-model="planForm.planStartDatetime"
            type="datetime"
            placeholder="选择开始时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
        
        <el-form-item label="结束时间" required>
          <el-date-picker
            v-model="planForm.planEndDatetime"
            type="datetime"
            placeholder="选择结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm" :loading="loading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-container {
  padding: 30px;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h1 {
  font-size: 28px;
  color: #303133;
}

.content-wrapper {
  background-color: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}
</style>
