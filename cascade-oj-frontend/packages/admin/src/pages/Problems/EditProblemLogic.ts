import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSingleProblem, putProblem, uploadTestCases } from '../../api/admin'
import { ProblemStatus, CodeTemplate } from '../../api/types' 

export function useEditProblem() {
    const route = useRoute()
    const router = useRouter()
    const loading = ref(false)
    const problemId = Number(route.params.id)

    // 1. 语言与文件状态管理
    const supportedLangs = ['cpp', 'java', 'python', 'go']
    const activeLang = ref('cpp')
    const activeFileName = ref('main.cpp') // 追踪当前正在编辑的文件名

    const form = reactive({
        title: '',
        timeLimitMs: 1000,
        memoryLimitMb: 256,
        description: '',
        status: ProblemStatus.PROBLEM_STATUS_UNAVAILABLE,
        templates: [] as CodeTemplate[] // 现在的模板是一个完整的数组
    })

    // 初始化默认的单文件结构（防止新题目一片空白）
    const initDefaultFiles = () => {
        if (form.templates.length === 0) {
            form.templates = [
                { language: 'cpp', name: 'main.cpp', code: '' },
                { language: 'java', name: 'Main.java', code: '' },
                { language: 'python', name: 'main.py', code: '' },
                { language: 'go', name: 'main.go', code: '' }
            ]
        }
    }

    // 2. 核心计算属性
    // 过滤出当前语言下的所有文件
    const currentLangFiles = computed(() => {
        return form.templates.filter(t => t.language === activeLang.value)
    })
    // 找出当前编辑器正在绑定的那个文件
    const activeTemplate = computed(() => {
        return form.templates.find(
            t => t.language === activeLang.value && t.name === activeFileName.value
        )
    })

    // 3. 文件操作方法
    const switchLang = (lang: string) => {
        activeLang.value = lang
        const files = form.templates.filter(t => t.language === lang)
        // 切换语言时，默认选中该语言的第一个文件
        activeFileName.value = files.length > 0 ? files[0].name : ''
    }

    const addFile = () => {
        const newName = prompt(`为 ${activeLang.value.toUpperCase()} 添加新文件，请输入文件名 :`)
        if (!newName) return
        
        // 查重
        if (form.templates.some(t => t.language === activeLang.value && t.name === newName)) {
            alert('该文件名已存在！')
            return
        }
        form.templates.push({ language: activeLang.value, name: newName, code: '' })
        activeFileName.value = newName // 焦点切到新文件
    }

    const renameFile = (oldName: string) => {
        const newName = prompt('重命名文件:', oldName)
        if (!newName || newName === oldName) return
        
        if (form.templates.some(t => t.language === activeLang.value && t.name === newName)) {
            alert('该文件名已存在！')
            return
        }
        const target = form.templates.find(t => t.language === activeLang.value && t.name === oldName)
        if (target) {
            target.name = newName
            if (activeFileName.value === oldName) activeFileName.value = newName
        }
    }

    const deleteFile = (name: string) => {
        if (!confirm(`确定要删除 ${name} 吗？`)) return
        
        const idx = form.templates.findIndex(t => t.language === activeLang.value && t.name === name)
        if (idx !== -1) {
            form.templates.splice(idx, 1)
            // 如果删掉的是当前正在看的文件，随便切到另一个文件
            if (activeFileName.value === name) {
                const remaining = form.templates.filter(t => t.language === activeLang.value)
                activeFileName.value = remaining.length > 0 ? remaining[0].name : ''
            }
        }
    }

    // 4. API 交互逻辑
    const fetchDetail = async () => {
        try {
            const res = await getSingleProblem(problemId)
            form.title = res.metadata.title
            form.timeLimitMs = res.metadata.timeLimitMs
            form.memoryLimitMb = res.metadata.memoryLimitMb
            form.description = res.description
            form.status = res.metadata.status
            
            if (res.templates && res.templates.length > 0) {
                form.templates = res.templates
            } else {
                initDefaultFiles()
            }
            switchLang(activeLang.value) 
        } catch (err) {
            console.error('Fetch detail failed', err)
        }
    }

    const handleFileUpload = async (event: Event) => {
        const file = (event.target as HTMLInputElement).files?.[0]
        if (!file) return
        loading.value = true
        try {
            await uploadTestCases(problemId, file)
            alert('Test cases uploaded successfully!')
        } catch (err: any) {
            alert(err.response?.data?.message || 'Upload failed')
        } finally {
            loading.value = false
            ;(event.target as HTMLInputElement).value = ''
        }
    }

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
                    status: form.status,
                    description: form.description
                },
                description: form.description,
                // 这里过滤掉没有任何代码，且是默认生成的 main 文件（保留用户特意创建的空文件）
                templates: form.templates.filter(t => t.code.trim() !== '' || t.name !== 'main.cpp' && t.name !== 'Main.java')
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

    return { 
        form, loading, router,
        activeLang, supportedLangs, 
        activeFileName, currentLangFiles, activeTemplate, 
        switchLang, addFile, renameFile, deleteFile,      
        handleFileUpload, handleUpdate
    }
}