import { defineStore } from 'pinia'

export const useMockDataStore = defineStore('mockData', {
  state: () => ({
    // 培训计划数据
    trainingPlans: [
      {
        plan_id: 1,
        plan_name: '2024年春季船舶技术培训',
        plan_status: '进行中',
        plan_start_datetime: '2024-03-01 09:00:00',
        plan_end_datetime: '2024-06-30 18:00:00',
        creator_id: 1
      },
      {
        plan_id: 2,
        plan_name: '2024年安全管理培训计划',
        plan_status: '规划中',
        plan_start_datetime: '2024-07-01 09:00:00',
        plan_end_datetime: '2024-09-30 18:00:00',
        creator_id: 1
      }
    ],
    
    // 课程数据
    courses: [
      {
        course_id: 1,
        course_name: '船舶结构基础',
        course_desc: '学习船舶基本结构和组成部分',
        course_require: '无基础要求',
        course_class: '船舶结构',
        teacher_id: 2
      },
      {
        course_id: 2,
        course_name: '船舶动力系统',
        course_desc: '掌握船舶动力系统原理与维护',
        course_require: '需要机械基础',
        course_class: '动力系统',
        teacher_id: 2
      },
      {
        course_id: 3,
        course_name: '航海安全管理',
        course_desc: '学习航海安全规范与应急处理',
        course_require: '需要基本航海知识',
        course_class: '安全管理',
        teacher_id: 3
      },
      {
        course_id: 4,
        course_name: '船舶电气系统',
        course_desc: '船舶电气系统的原理、维护与故障排除',
        course_require: '需要电气基础知识',
        course_class: '电气系统',
        teacher_id: 2
      }
    ],
    
    // 课程安排数据
    courseItems: [
      {
        item_id: 1,
        plan_id: 1,
        course_id: 1,
        class_date: '2024-03-15',
        class_begin_time: '09:00:00',
        class_end_time: '11:00:00',
        location: '培训中心A201'
      },
      {
        item_id: 2,
        plan_id: 1,
        course_id: 2,
        class_date: '2024-03-18',
        class_begin_time: '14:00:00',
        class_end_time: '16:00:00',
        location: '培训中心A202'
      },
      {
        item_id: 3,
        plan_id: 1,
        course_id: 3,
        class_date: '2024-12-20',
        class_begin_time: '09:00:00',
        class_end_time: '11:00:00',
        location: '培训中心B101'
      },
      {
        item_id: 4,
        plan_id: 1,
        course_id: 4,
        class_date: '2024-12-20',
        class_begin_time: '14:00:00',
        class_end_time: '17:00:00',
        location: '培训中心A203'
      },
      {
        item_id: 5,
        plan_id: 1,
        course_id: 1,
        class_date: '2024-12-23',
        class_begin_time: '10:00:00',
        class_end_time: '12:00:00',
        location: '培训中心A201'
      }
    ],
    
    // 人员数据
    persons: [
      {
        person_id: 1,
        name: '张主任',
        role: '课程大纲制定者'
      },
      {
        person_id: 2,
        name: '李老师',
        role: '讲师'
      },
      {
        person_id: 3,
        name: '王老师',
        role: '讲师'
      },
      {
        person_id: 4,
        name: '刘员工',
        role: '员工'
      },
      {
        person_id: 5,
        name: '陈员工',
        role: '员工'
      }
    ],
    
    // 账号数据
    accounts: [
      {
        account_id: 1,
        person_id: 1,
        login_name: 'planner',
        password_hash: '123456'
      },
      {
        account_id: 2,
        person_id: 2,
        login_name: 'teacher',
        password_hash: '123456'
      },
      {
        account_id: 3,
        person_id: 4,
        login_name: 'employee',
        password_hash: '123456'
      }
    ],
    
    // 参与和评价数据
    attendanceEvaluations: [
      {
        person_id: 4,
        item_id: 1,
        self_score: 85.5,
        self_comment: '课程内容很实用，基本掌握了船舶结构的基础知识',
        teacher_score: 88.0,
        teacher_comment: '学习态度认真，课堂表现积极',
        score_ratio: 0.7
      },
      {
        person_id: 4,
        item_id: 2,
        self_score: 78.0,
        self_comment: '动力系统较为复杂，还需要进一步学习',
        teacher_score: 82.0,
        teacher_comment: '理解能力较好，需加强实践操作',
        score_ratio: 0.7
      },
      {
        person_id: 5,
        item_id: 1,
        self_score: 90.0,
        self_comment: '对船舶结构有了全面的认识',
        teacher_score: 92.0,
        teacher_comment: '优秀学员，掌握扎实',
        score_ratio: 0.7
      }
    ]
  }),
  
  getters: {
    // 获取特定培训计划的课程安排
    getCourseItemsByPlanId: (state) => (planId) => {
      return state.courseItems.filter(item => item.plan_id === planId)
    },
    
    // 获取今日课程（员工）
    getTodayCoursesForEmployee: (state) => (personId) => {
      const today = new Date().toISOString().split('T')[0]
      return state.courseItems
        .filter(item => item.class_date === today)
        .map(item => ({
          ...item,
          course: state.courses.find(c => c.course_id === item.course_id),
          plan: state.trainingPlans.find(p => p.plan_id === item.plan_id)
        }))
    },
    
    // 获取今日课程（讲师）
    getTodayCoursesForTeacher: (state) => (teacherId) => {
      const today = new Date().toISOString().split('T')[0]
      return state.courseItems
        .filter(item => {
          const course = state.courses.find(c => c.course_id === item.course_id)
          return course && course.teacher_id === teacherId && item.class_date === today
        })
        .map(item => ({
          ...item,
          course: state.courses.find(c => c.course_id === item.course_id),
          plan: state.trainingPlans.find(p => p.plan_id === item.plan_id)
        }))
    },
    
    // 获取员工成绩
    getScoresByPersonId: (state) => (personId) => {
      return state.attendanceEvaluations
        .filter(ae => ae.person_id === personId)
        .map(ae => {
          const item = state.courseItems.find(i => i.item_id === ae.item_id)
          const course = state.courses.find(c => c.course_id === item?.course_id)
          const weightedScore = ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
          return {
            ...ae,
            item,
            course,
            weighted_score: weightedScore.toFixed(2)
          }
        })
    },
    
    // 获取课程类型统计
    getCourseTypeScores: (state) => (personId) => {
      const scores = state.attendanceEvaluations
        .filter(ae => ae.person_id === personId)
        .map(ae => {
          const item = state.courseItems.find(i => i.item_id === ae.item_id)
          const course = state.courses.find(c => c.course_id === item?.course_id)
          const weightedScore = ae.self_score * (1 - ae.score_ratio) + ae.teacher_score * ae.score_ratio
          return {
            course_class: course?.course_class,
            weighted_score: weightedScore
          }
        })
      
      // 按课程类型分组计算平均分
      const typeScores = {}
      scores.forEach(s => {
        if (!typeScores[s.course_class]) {
          typeScores[s.course_class] = []
        }
        typeScores[s.course_class].push(s.weighted_score)
      })
      
      return Object.entries(typeScores).map(([type, scores]) => ({
        course_class: type,
        avg_score: (scores.reduce((a, b) => a + b, 0) / scores.length).toFixed(2)
      }))
    }
  }
})
