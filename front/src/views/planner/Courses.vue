<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCoursesList, createCourse, updateCourse, deleteCourse } from '../../api/planner'
import axios from 'axios'

const loading = ref(false)
const activeTab = ref('courses')
const dialogVisible = ref(false)
const dialogMode = ref('add')

// 课程列表数据
const coursesList = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 讲师列表
const teachers = ref([])

// 课程表单
const courseForm = ref({
  courseName: '',
  courseDesc: '',
  courseRequire: '',
  courseClass: '',
  teacherId: null
})

// 获取讲师列表
const fetchTeachersList = async () => {
  try {
    const sessionId = localStorage.getItem('sessionId')
    const response = await axios.get('http://localhost:8080/api/planner/teachers', {
      headers: { 'Session-ID': sessionId }
    })
    if (response.data.code === 200) {
      teachers.value = response.data.data.list || []
    }
  } catch (error) {
    console.error('获取讲师列表失败：', error)
  }
}

// 获取课程列表
const fetchCoursesList = async () => {
  loading.value = true
  try {
    const response = await getCoursesList({
      page: currentPage.value,
      pageSize: pageSize.value
    })
    if (response.code === 200) {
      coursesList.value = response.data.list
      total.value = response.data.total
    }
  } catch (error) {
    ElMessage.error('获取课程列表失败：' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 分页变化
const handlePageChange = (page) => {
  currentPage.value = page
  fetchCoursesList()
}

// 打开新增课程对话框
const openAddCourseDialog = () => {
  dialogMode.value = 'add-course'
  courseForm.value = {
    courseName: '',
    courseDesc: '',
    courseRequire: '',
    courseClass: '',
    teacherId: null
  }
  dialogVisible.value = true
}

// 打开编辑课程对话框
const openEditCourseDialog = (course) => {
  dialogMode.value = 'edit-course'
  courseForm.value = {
    courseId: course.courseId,
    courseName: course.courseName,
    courseDesc: course.courseDesc,
    courseRequire: course.courseRequire,
    courseClass: course.courseClass,
    teacherId: course.teacherId
  }
  dialogVisible.value = true
}

// 提交表单
const submitForm = async () => {
  if (!courseForm.value.courseName) {
    ElMessage.warning('请输入课程名称')
    return
  }
  if (!courseForm.value.courseClass) {
    ElMessage.warning('请输入课程类型')
    return
  }
  if (!courseForm.value.teacherId) {
    ElMessage.warning('请选择讲师')
    return
  }

  loading.value = true
  try {
    if (dialogMode.value === 'add-course') {
      // 确保 teacherId 是数字类型
      const courseData = {
        ...courseForm.value,
        teacherId: Number(courseForm.value.teacherId)
      }
      await createCourse(courseData)
      ElMessage.success('创建课程成功')
    } else {
      const courseId = courseForm.value.courseId
      const updateData = {
        courseName: courseForm.value.courseName,
        courseDesc: courseForm.value.courseDesc,
        courseRequire: courseForm.value.courseRequire,
        courseClass: courseForm.value.courseClass,
        teacherId: Number(courseForm.value.teacherId)
      }
      await updateCourse(courseId, updateData)
      ElMessage.success('修改课程成功')
    }
    dialogVisible.value = false
    fetchCoursesList()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

// 删除课程
const deleteCourseHandler = (course) => {
  ElMessageBox.confirm(
    `确定要删除课程"${course.courseName}"吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    try {
      await deleteCourse(course.courseId)
      ElMessage.success('删除成功')
      fetchCoursesList()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  }).catch(() => {})
}

onMounted(() => {
  fetchTeachersList()
  fetchCoursesList()
})
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1>课程管理</h1>
    </div>
    
    <div class="content-wrapper" v-loading="loading">
      <div class="tab-header">
        <el-button type="primary" @click="openAddCourseDialog">
          新建课程
        </el-button>
      </div>
      
      <el-table :data="coursesList" border>
        <el-table-column prop="courseId" label="课程ID" width="80" />
        <el-table-column prop="courseName" label="课程名称" min-width="150" />
        <el-table-column prop="courseDesc" label="课程描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="courseClass" label="课程类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.courseClass }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="teacherName" label="主讲讲师" width="120" />
        <el-table-column prop="scheduledCount" label="安排次数" width="100" />
        <el-table-column label="操作" width="180" align="center">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="openEditCourseDialog(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteCourseHandler(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        v-if="total > 0"
        v-model:current-page="currentPage"
        :page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
        style="margin-top: 20px; justify-content: center;"
      />
    </div>
    
    <!-- 课程对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="dialogMode === 'add-course' ? '新建课程' : '编辑课程'"
      width="600px"
    >
      <el-form :model="courseForm" label-width="100px">
        <el-form-item label="课程名称" required>
          <el-input v-model="courseForm.courseName" placeholder="请输入课程名称" />
        </el-form-item>
        
        <el-form-item label="课程描述">
          <el-input v-model="courseForm.courseDesc" type="textarea" :rows="3" placeholder="请输入课程描述" />
        </el-form-item>
        
        <el-form-item label="课程要求">
          <el-input v-model="courseForm.courseRequire" type="textarea" :rows="2" placeholder="请输入课程要求" />
        </el-form-item>
        
        <el-form-item label="课程类型" required>
          <el-input v-model="courseForm.courseClass" placeholder="如：船舶结构、动力系统等" />
        </el-form-item>
        
        <el-form-item label="主讲讲师" required>
          <el-select v-model="courseForm.teacherId" placeholder="请选择讲师" style="width: 100%;">
            <el-option
              v-for="teacher in teachers"
              :key="teacher.personId"
              :label="teacher.name"
              :value="teacher.personId"
            />
          </el-select>
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
