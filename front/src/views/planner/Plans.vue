<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Delete } from '@element-plus/icons-vue'
import { 
  getPlansList, 
  createPlan, 
  updatePlan, 
  deletePlan,
  getPlanDetail,
  addEmployeesToPlan,
  removeEmployeeFromPlan,
  getEmployeesList,
  getCoursesList,
  getCourseItemsList,
  createCourseItem,
  updateCourseItem,
  deleteCourseItem
} from '../../api/planner'

const loading = ref(false)
const dialogVisible = ref(false)
const detailDialogVisible = ref(false)
const addEmployeeDialogVisible = ref(false)
const courseItemDialogVisible = ref(false)

const dialogMode = ref('add') // 'add' or 'edit'
const courseItemMode = ref('add') // 'add' or 'edit'
const currentPlan = ref(null)
const plansList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 详情数据
const planDetail = ref(null)
const detailEmployees = ref([])
const detailCourseItems = ref([])

// 可选数据
const allEmployees = ref([])
const allCourses = ref([])
const selectedEmployeeIds = ref([])

const planForm = ref({
  planName: '',
  planStatus: '规划中',
  planStartDatetime: '',
  planEndDatetime: ''
})

const courseItemForm = ref({
  courseId: null,
  classDate: '',
  classBeginTime: '',
  classEndTime: '',
  location: ''
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

// 获取计划详情
const fetchPlanDetail = async (planId) => {
  try {
    const response = await getPlanDetail(planId)
    if (response.code === 200) {
      planDetail.value = response.data
      detailEmployees.value = response.data.employees || []
      detailCourseItems.value = response.data.courseItems || []
    }
  } catch (error) {
    ElMessage.error('获取计划详情失败：' + (error.message || '未知错误'))
  }
}

// 获取所有员工列表
const fetchAllEmployees = async () => {
  try {
    const response = await getEmployeesList()
    if (response.code === 200) {
      allEmployees.value = response.data || []
    }
  } catch (error) {
    ElMessage.error('获取员工列表失败')
  }
}

// 获取所有课程列表
const fetchAllCourses = async () => {
  try {
    const response = await getCoursesList({ page: 1, pageSize: 1000 })
    if (response.code === 200) {
      allCourses.value = response.data.list || []
    }
  } catch (error) {
    ElMessage.error('获取课程列表失败')
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
const viewPlanDetail = async (plan) => {
  currentPlan.value = plan
  await fetchPlanDetail(plan.planId)
  detailDialogVisible.value = true
}

// 打开添加员工对话框
const openAddEmployeeDialog = async () => {
  await fetchAllEmployees()
  selectedEmployeeIds.value = []
  addEmployeeDialogVisible.value = true
}

// 添加员工到计划
const submitAddEmployees = async () => {
  if (selectedEmployeeIds.value.length === 0) {
    ElMessage.warning('请至少选择一个员工')
    return
  }
  
  try {
    await addEmployeesToPlan(currentPlan.value.planId, {
      employeeIds: selectedEmployeeIds.value
    })
    ElMessage.success('添加员工成功')
    addEmployeeDialogVisible.value = false
    await fetchPlanDetail(currentPlan.value.planId)
  } catch (error) {
    ElMessage.error(error.message || '添加失败')
  }
}

// 移除员工
const handleRemoveEmployee = (employee) => {
  ElMessageBox.confirm(
    `确定要将员工"${employee.name}"从该培训计划中移除吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await removeEmployeeFromPlan(currentPlan.value.planId, employee.personId)
      ElMessage.success('移除成功')
      await fetchPlanDetail(currentPlan.value.planId)
    } catch (error) {
      ElMessage.error(error.message || '移除失败')
    }
  }).catch(() => {})
}

// 打开课程安排对话框
const openCourseItemDialog = async (mode, item = null) => {
  courseItemMode.value = mode
  await fetchAllCourses()
  
  if (mode === 'add') {
    courseItemForm.value = {
      courseId: null,
      classDate: '',
      classBeginTime: '',
      classEndTime: '',
      location: ''
    }
  } else {
    courseItemForm.value = {
      itemId: item.itemId,
      courseId: item.courseId,
      classDate: item.classDate,
      classBeginTime: item.classBeginTime,
      classEndTime: item.classEndTime,
      location: item.location
    }
  }
  
  courseItemDialogVisible.value = true
}

// 提交课程安排
const submitCourseItem = async () => {
  if (!courseItemForm.value.courseId) {
    ElMessage.warning('请选择课程')
    return
  }
  if (!courseItemForm.value.classDate) {
    ElMessage.warning('请选择上课日期')
    return
  }
  if (!courseItemForm.value.classBeginTime) {
    ElMessage.warning('请选择开始时间')
    return
  }
  if (!courseItemForm.value.classEndTime) {
    ElMessage.warning('请选择结束时间')
    return
  }
  if (!courseItemForm.value.location) {
    ElMessage.warning('请输入上课地点')
    return
  }
  
  try {
    if (courseItemMode.value === 'add') {
      await createCourseItem({
        planId: currentPlan.value.planId,
        courseId: courseItemForm.value.courseId,
        classDate: courseItemForm.value.classDate,
        classBeginTime: courseItemForm.value.classBeginTime,
        classEndTime: courseItemForm.value.classEndTime,
        location: courseItemForm.value.location
      })
      ElMessage.success('课程安排创建成功')
    } else {
      await updateCourseItem(courseItemForm.value.itemId, {
        classDate: courseItemForm.value.classDate,
        classBeginTime: courseItemForm.value.classBeginTime,
        classEndTime: courseItemForm.value.classEndTime,
        location: courseItemForm.value.location
      })
      ElMessage.success('课程安排更新成功')
    }
    
    courseItemDialogVisible.value = false
    await fetchPlanDetail(currentPlan.value.planId)
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  }
}

// 删除课程安排
const handleDeleteCourseItem = (item) => {
  ElMessageBox.confirm(
    `确定要删除该课程安排吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteCourseItem(item.itemId)
      ElMessage.success('删除成功')
      await fetchPlanDetail(currentPlan.value.planId)
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {})
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
        <el-table-column label="操作" width="250" align="center" fixed="right">
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

    <!-- 培训计划详情对话框 -->
    <el-dialog 
      v-model="detailDialogVisible" 
      title="培训计划详情"
      width="1200px"
      top="5vh"
    >
      <div v-if="planDetail" class="detail-container">
        <!-- 基本信息 -->
        <el-card class="detail-card">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
            </div>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="计划名称">{{ planDetail.planName }}</el-descriptions-item>
            <el-descriptions-item label="计划状态">
              <el-tag :type="getStatusType(planDetail.planStatus)">{{ planDetail.planStatus }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="开始时间">{{ planDetail.planStartDatetime }}</el-descriptions-item>
            <el-descriptions-item label="结束时间">{{ planDetail.planEndDatetime }}</el-descriptions-item>
            <el-descriptions-item label="制定人">{{ planDetail.creatorName }}</el-descriptions-item>
            <el-descriptions-item label="参与员工">{{ detailEmployees.length }} 人</el-descriptions-item>
          </el-descriptions>
        </el-card>

        <!-- 参与员工 -->
        <el-card class="detail-card">
          <template #header>
            <div class="card-header">
              <span>参与员工</span>
              <el-button type="primary" size="small" @click="openAddEmployeeDialog">
                <el-icon><Plus /></el-icon>
                添加员工
              </el-button>
            </div>
          </template>
          <el-table :data="detailEmployees" border max-height="300">
            <el-table-column prop="personId" label="员工ID" width="100" />
            <el-table-column prop="name" label="姓名" />
            <el-table-column label="操作" width="100" align="center">
              <template #default="{ row }">
                <el-button size="small" type="danger" :icon="Delete" @click="handleRemoveEmployee(row)" />
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <!-- 课程安排 -->
        <el-card class="detail-card">
          <template #header>
            <div class="card-header">
              <span>课程安排</span>
              <el-button type="primary" size="small" @click="openCourseItemDialog('add')">
                <el-icon><Plus /></el-icon>
                添加课程安排
              </el-button>
            </div>
          </template>
          <el-table :data="detailCourseItems" border max-height="400">
            <el-table-column prop="courseName" label="课程名称" min-width="150" />
            <el-table-column prop="courseClass" label="课程类型" width="100" />
            <el-table-column prop="classDate" label="上课日期" width="120" />
            <el-table-column label="上课时间" width="160">
              <template #default="{ row }">
                {{ row.classBeginTime }} - {{ row.classEndTime }}
              </template>
            </el-table-column>
            <el-table-column prop="location" label="上课地点" min-width="120" />
            <el-table-column prop="teacherName" label="讲师" width="100" />
            <el-table-column label="操作" width="150" align="center">
              <template #default="{ row }">
                <el-button size="small" type="primary" @click="openCourseItemDialog('edit', row)">编辑</el-button>
                <el-button size="small" type="danger" @click="handleDeleteCourseItem(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </div>
    </el-dialog>

    <!-- 添加员工对话框 -->
    <el-dialog 
      v-model="addEmployeeDialogVisible" 
      title="添加员工"
      width="600px"
    >
      <el-select 
        v-model="selectedEmployeeIds" 
        multiple
        filterable
        placeholder="请选择员工"
        style="width: 100%;"
      >
        <el-option
          v-for="emp in allEmployees"
          :key="emp.personId"
          :label="`${emp.personName} (ID: ${emp.personId})`"
          :value="emp.personId"
        />
      </el-select>
      
      <template #footer>
        <el-button @click="addEmployeeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitAddEmployees">确定</el-button>
      </template>
    </el-dialog>

    <!-- 课程安排对话框 -->
    <el-dialog 
      v-model="courseItemDialogVisible" 
      :title="courseItemMode === 'add' ? '添加课程安排' : '编辑课程安排'"
      width="600px"
    >
      <el-form :model="courseItemForm" label-width="120px">
        <el-form-item label="课程" required>
          <el-select 
            v-model="courseItemForm.courseId" 
            placeholder="请选择课程"
            filterable
            style="width: 100%;"
            :disabled="courseItemMode === 'edit'"
          >
            <el-option
              v-for="course in allCourses"
              :key="course.courseId"
              :label="`${course.courseName} (${course.courseClass})`"
              :value="course.courseId"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="上课日期" required>
          <el-date-picker
            v-model="courseItemForm.classDate"
            type="date"
            placeholder="选择上课日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            style="width: 100%;"
          />
        </el-form-item>
        
        <el-form-item label="开始时间" required>
          <el-time-picker
            v-model="courseItemForm.classBeginTime"
            placeholder="选择开始时间"
            format="HH:mm:ss"
            value-format="HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
        
        <el-form-item label="结束时间" required>
          <el-time-picker
            v-model="courseItemForm.classEndTime"
            placeholder="选择结束时间"
            format="HH:mm:ss"
            value-format="HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
        
        <el-form-item label="上课地点" required>
          <el-input 
            v-model="courseItemForm.location" 
            placeholder="请输入上课地点"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="courseItemDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCourseItem">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.page-container {
  padding: 30px;
  max-width: 1600px;
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

.detail-container {
  max-height: 70vh;
  overflow-y: auto;
}

.detail-card {
  margin-bottom: 20px;
}

.detail-card:last-child {
  margin-bottom: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: bold;
}
</style>
