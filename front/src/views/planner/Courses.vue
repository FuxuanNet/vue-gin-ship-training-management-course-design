<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useMockDataStore } from '../../stores/mockData'

const mockDataStore = useMockDataStore()
const activeTab = ref('courses')
const dialogVisible = ref(false)
const dialogMode = ref('add')

// 课程表单
const courseForm = ref({
  course_name: '',
  course_desc: '',
  course_require: '',
  course_class: '',
  teacher_id: null
})

// 课程安排表单
const itemForm = ref({
  plan_id: null,
  course_id: null,
  class_date: '',
  class_begin_time: '',
  class_end_time: '',
  location: ''
})

// 获取讲师列表
const teachers = mockDataStore.persons.filter(p => p.role === '讲师')

// 打开新增课程对话框
const openAddCourseDialog = () => {
  dialogMode.value = 'add-course'
  courseForm.value = {
    course_name: '',
    course_desc: '',
    course_require: '',
    course_class: '',
    teacher_id: null
  }
  dialogVisible.value = true
}

// 打开编辑课程对话框
const openEditCourseDialog = (course) => {
  dialogMode.value = 'edit-course'
  courseForm.value = { ...course }
  dialogVisible.value = true
}

// 打开新增课程安排对话框
const openAddItemDialog = () => {
  dialogMode.value = 'add-item'
  itemForm.value = {
    plan_id: null,
    course_id: null,
    class_date: '',
    class_begin_time: '',
    class_end_time: '',
    location: ''
  }
  dialogVisible.value = true
}

// 提交表单
const submitForm = () => {
  if (dialogMode.value.includes('course')) {
    if (!courseForm.value.course_name) {
      ElMessage.warning('请输入课程名称')
      return
    }
    ElMessage.success(dialogMode.value === 'add-course' ? '课程创建成功' : '课程更新成功')
    console.log('提交课程:', courseForm.value)
  } else {
    if (!itemForm.value.plan_id || !itemForm.value.course_id) {
      ElMessage.warning('请完整填写表单')
      return
    }
    ElMessage.success('课程安排添加成功')
    console.log('提交课程安排:', itemForm.value)
  }
  
  dialogVisible.value = false
}

// 删除课程
const deleteCourse = (course) => {
  ElMessageBox.confirm(
    `确定要删除课程"${course.course_name}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    ElMessage.success('删除成功')
    console.log('删除课程:', course.course_id)
  }).catch(() => {})
}

// 删除课程安排
const deleteItem = (item) => {
  ElMessageBox.confirm(
    '确定要删除这个课程安排吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    ElMessage.success('删除成功')
    console.log('删除课程安排:', item.item_id)
  }).catch(() => {})
}

// 获取课程名称
const getCourseName = (courseId) => {
  const course = mockDataStore.courses.find(c => c.course_id === courseId)
  return course?.course_name || '-'
}

// 获取计划名称
const getPlanName = (planId) => {
  const plan = mockDataStore.trainingPlans.find(p => p.plan_id === planId)
  return plan?.plan_name || '-'
}

// 获取讲师名称
const getTeacherName = (teacherId) => {
  const teacher = mockDataStore.persons.find(p => p.person_id === teacherId)
  return teacher?.name || '-'
}
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>课程管理</h1>
    </div>
    
    <div class="content-wrapper">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="课程列表" name="courses">
          <div class="tab-header">
            <el-button type="primary" @click="openAddCourseDialog">
              <el-icon><Plus /></el-icon>
              新建课程
            </el-button>
          </div>
          
          <el-table :data="mockDataStore.courses" border>
            <el-table-column prop="course_id" label="课程ID" width="80" />
            <el-table-column prop="course_name" label="课程名称" min-width="150" />
            <el-table-column prop="course_desc" label="课程描述" min-width="200" show-overflow-tooltip />
            <el-table-column prop="course_class" label="课程类型" width="120">
              <template #default="{ row }">
                <el-tag size="small">{{ row.course_class }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="主讲讲师" width="120">
              <template #default="{ row }">
                {{ getTeacherName(row.teacher_id) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="180" align="center">
              <template #default="{ row }">
                <el-button size="small" type="primary" @click="openEditCourseDialog(row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteCourse(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
        
        <el-tab-pane label="课程安排" name="items">
          <div class="tab-header">
            <el-button type="primary" @click="openAddItemDialog">
              <el-icon><Plus /></el-icon>
              添加课程安排
            </el-button>
          </div>
          
          <el-table :data="mockDataStore.courseItems" border>
            <el-table-column prop="item_id" label="安排ID" width="80" />
            <el-table-column label="培训计划" width="180">
              <template #default="{ row }">
                {{ getPlanName(row.plan_id) }}
              </template>
            </el-table-column>
            <el-table-column label="课程名称" min-width="150">
              <template #default="{ row }">
                {{ getCourseName(row.course_id) }}
              </template>
            </el-table-column>
            <el-table-column prop="class_date" label="上课日期" width="120" />
            <el-table-column label="上课时间" width="150">
              <template #default="{ row }">
                {{ row.class_begin_time.slice(0, 5) }} - {{ row.class_end_time.slice(0, 5) }}
              </template>
            </el-table-column>
            <el-table-column prop="location" label="上课地点" width="150" />
            <el-table-column label="操作" width="120" align="center">
              <template #default="{ row }">
                <el-button size="small" type="danger" @click="deleteItem(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>
    
    <!-- 课程对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="dialogMode === 'add-course' ? '新建课程' : dialogMode === 'edit-course' ? '编辑课程' : '添加课程安排'"
      width="600px"
    >
      <el-form v-if="dialogMode.includes('course')" :model="courseForm" label-width="100px">
        <el-form-item label="课程名称" required>
          <el-input v-model="courseForm.course_name" placeholder="请输入课程名称" />
        </el-form-item>
        
        <el-form-item label="课程描述">
          <el-input v-model="courseForm.course_desc" type="textarea" :rows="3" placeholder="请输入课程描述" />
        </el-form-item>
        
        <el-form-item label="课程要求">
          <el-input v-model="courseForm.course_require" type="textarea" :rows="2" placeholder="请输入课程要求" />
        </el-form-item>
        
        <el-form-item label="课程类型" required>
          <el-input v-model="courseForm.course_class" placeholder="如：船舶结构、动力系统等" />
        </el-form-item>
        
        <el-form-item label="主讲讲师" required>
          <el-select v-model="courseForm.teacher_id" placeholder="请选择讲师" style="width: 100%;">
            <el-option 
              v-for="teacher in teachers" 
              :key="teacher.person_id"
              :label="teacher.name"
              :value="teacher.person_id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      
      <el-form v-else :model="itemForm" label-width="100px">
        <el-form-item label="培训计划" required>
          <el-select v-model="itemForm.plan_id" placeholder="请选择培训计划" style="width: 100%;">
            <el-option 
              v-for="plan in mockDataStore.trainingPlans" 
              :key="plan.plan_id"
              :label="plan.plan_name"
              :value="plan.plan_id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="课程" required>
          <el-select v-model="itemForm.course_id" placeholder="请选择课程" style="width: 100%;">
            <el-option 
              v-for="course in mockDataStore.courses" 
              :key="course.course_id"
              :label="course.course_name"
              :value="course.course_id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="上课日期" required>
          <el-date-picker
            v-model="itemForm.class_date"
            type="date"
            placeholder="选择上课日期"
            value-format="YYYY-MM-DD"
            style="width: 100%;"
          />
        </el-form-item>
        
        <el-form-item label="开始时间" required>
          <el-time-picker
            v-model="itemForm.class_begin_time"
            placeholder="选择开始时间"
            value-format="HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
        
        <el-form-item label="结束时间" required>
          <el-time-picker
            v-model="itemForm.class_end_time"
            placeholder="选择结束时间"
            value-format="HH:mm:ss"
            style="width: 100%;"
          />
        </el-form-item>
        
        <el-form-item label="上课地点" required>
          <el-input v-model="itemForm.location" placeholder="请输入上课地点" />
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

.tab-header {
  margin-bottom: 20px;
}
</style>
