import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { postProblem } from '../../api/admin'
import { ProblemStatus } from '../../api/types'

export function useNewProblem() {
    const router = useRouter()
    const loading = ref(false)

    const form = reactive({
        title: '',
        timeLimitMs: 1000,
        memoryLimitMb: 256,
        description: '', 
        codeTemplate: '// Write your code template here...' 
    })

    // 保存题目
    const handleSave = async () => {
    if (!form.title) return alert('Please enter a title')
    
    loading.value = true
    try {
        await postProblem({
            metadata: {
                id: 0,
                title: form.title,
                timeLimitMs: form.timeLimitMs,
                memoryLimitMb: form.memoryLimitMb,
                status: ProblemStatus.PROBLEM_STATUS_UNAVAILABLE 
            },
            description: form.description,
            codeTemplate: form.codeTemplate
        })
        alert('Problem created successfully!')
        router.push('/problems')
    } catch (err) {
        console.error(err)
        alert('Failed to create problem')
    } finally {
        loading.value = false
    }
}
    return {
        form,
        loading,
        handleSave,
        router
    }
}