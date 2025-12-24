<script setup>
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useMockDataStore } from '../../stores/mockData'

const mockDataStore = useMockDataStore()
const dialogVisible = ref(false)
const dialogMode = ref('add') // 'add' or 'edit'
const currentPlan = ref(null)

const planForm = ref({
  plan_name: '',
  plan_status: '规划中',
  plan_start_datetime: '',
  plan_end_datetime: ''
})

const statusOptions = [
  { label: '规划中', value: '规划中' },
  { label: '进行中', value: '进行中' },
  { label: '已完成', value: '已完成' }
]

// 打开新增对话框
const openAddDialog = () => {
  dialogMode.value = 'add'
  planForm.value = {
    plan_name: '',
    plan_status: '规划中',
    plan_start_datetime: '',
    plan_end_datetime: ''
  }
  dialogVisible.value = true
}

// 打开编辑对话框
const openEditDialog = (plan) => {
  dialogMode.value = 'edit'
  currentPlan.value = plan
  planForm.value = {
    plan_name: plan.plan_name,
    plan_status: plan.plan_status,
    plan_start_datetime: plan.plan_start_datetime,
    plan_end_datetime: plan.plan_end_datetime
  }
  dialogVisible.value = true
}

// 提交表单
const submitForm = () => {
  if (!planForm.value.plan_name) {
    ElMessage.warning('请输入计划名称')
    return
  }
  
  if (dialogMode.value === 'add') {
    ElMessage.success('培训计划创建成功')
  } else {
    ElMessage.success('培训计划更新成功')
  }
  
  dialogVisible.value = false
  
  console.log('提交表单:', planForm.value)
}

// 删除计划
const deletePlan = (plan) => {
  ElMessageBox.confirm(
    `确定要删除培训计划"${plan.plan_name}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    ElMessage.success('删除成功')
    console.log('删除计划:', plan.plan_id)
  }).catch(() => {})
}

// 查看计划详情
const viewPlanDetail = (plan) => {
  console.log('查看计划详情:', plan)
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
    
    <div class="content-wrapper">
      <el-table :data="mockDataStore.trainingPlans" border>
        <el-table-column prop="plan_id" label="计划ID" width="80" />
        <el-table-column prop="plan_name" label="计划名称" min-width="200" />
        <el-table-column label="计划状态" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.plan_status)">{{ row.plan_status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="plan_start_datetime" label="开始时间" width="180" />
        <el-table-column prop="plan_end_datetime" label="结束时间" width="180" />
        <el-table-column label="课程数" width="100" align="center">
          <template #default="{ row }">
            {{ mockDataStore.getCourseItemsByPlanId(row.plan_id).length }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" align="center">
          <template #default="{ row }">
            <el-button size="small" @click="viewPlanDetail(row)">详情</el-button>
            <el-button size="small" type="primary" @click="openEditDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deletePlan(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 新增/编辑对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="dialogMode === 'add' ? '新建培训计划' : '编辑培训计划'"
      width="600px"
    >
      <el-form :model="planForm" label-width="120px">
        <el-form-item label="计划名称" required>
          <el-input v-model="planForm.plan_name" placeholder="请输入培训计划名称" />
        </el-form-item>
        
        <el-form-item label="计划状态">
          <el-select v-model="planForm.plan_status" placeholder="请选择计划状态">
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
            v-model="planForm.plan_start_datetime"
            type="datetime"
            placeholder="选择开始时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
        
        <el-form-item label="结束时间" required>
          <el-date-picker
            v-model="planForm.plan_end_datetime"
            type="datetime"
            placeholder="选择结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
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
