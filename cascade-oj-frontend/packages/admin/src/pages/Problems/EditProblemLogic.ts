import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSingleProblem, putProblem } from '../../api/admin'
import { ProblemStatus } from '../../api/types'

export function useEditProblem() {
    const route = useRoute()
    const router = useRouter()
    const loading = ref(false)
    const problemId = Number(route.params.id) 

    const form = reactive({
        title: '',
        timeLimitMs: 1000,
        memoryLimitMb: 256,
        description: '',
        codeTemplate: '',
        status: ProblemStatus.PROBLEM_STATUS_UNAVAILABLE
    })

    // 1. 初始化：从后端拉取现有数据
    const fetchDetail = async () => {
        try {
            const res = await getSingleProblem(problemId)
            form.title = res.metadata.title
            form.timeLimitMs = res.metadata.timeLimitMs
            form.memoryLimitMb = res.metadata.memoryLimitMb
            form.description = res.description
            form.codeTemplate = res.codeTemplate
            form.status = res.metadata.status
        } catch (err) {
            console.error('Failed to fetch detail', err)
            alert('题目不存在或加载失败')
        }
    }

    // 2. 保存修改
    const handleUpdate = async () => {
        loading.value = true
        try {
            await putProblem(problemId, {
                problemId: problemId, 
                metadata: {
                    id: problemId,
                    title: form.title,
                    timeLimitMs: form.timeLimitMs,
                    memoryLimitMb: form.memoryLimitMb,
                    status: form.status
                },
                description: form.description,
                codeTemplate: form.codeTemplate
            })
            alert('Update successful!')
            router.push('/problems')
        } catch (err: any) {
            alert(err.response?.data?.message || 'Update failed')
        } finally {
            loading.value = false
        }
    }

    onMounted(fetchDetail)

    return { form, loading, handleUpdate, router }
}