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
    const activeLang = ref('') 
    const activeFileName = ref('') 

    const form = reactive({
        title: '',
        timeLimitMs: 1000,
        memoryLimitMb: 256,
        description: '',
        status: ProblemStatus.PROBLEM_STATUS_UNAVAILABLE,
        templates: [] as CodeTemplate[] 
    })

    // 完全解耦 PL 列表。
    // 不再写死支持的语言，而是动态从当前题目的模板数据中提取存在的语言
    const existingLangs = computed(() => {
        const langs = new Set(form.templates.map(t => t.language))
        return Array.from(langs)
    })

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
        activeFileName.value = files.length > 0 ? files[0].name : ''
    }

    // 允许管理员手动输入任何语言，实现高自由度自定义
    const addFile = () => {
        const lang = prompt('请输入编程语言 :')
        if (!lang) return
        
        const newName = prompt(`为 ${lang.toUpperCase()} 添加新文件，请输入文件名 :`)
        if (!newName) return
        
        const lowerLang = lang.toLowerCase()
        if (form.templates.some(t => t.language === lowerLang && t.name === newName)) {
            alert('该语言下已存在同名文件！')
            return
        }
        
        form.templates.push({ language: lowerLang, name: newName, code: '' })
        activeLang.value = lowerLang
        activeFileName.value = newName
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
            const res: any = await getSingleProblem(problemId)
            // 增加安全赋值，防止某一层级为 undefined 导致程序崩溃
            if (res && res.metadata) {
                form.title = res.metadata.title || ''
                form.timeLimitMs = res.metadata.timeLimitMs || 1000
                form.memoryLimitMb = res.metadata.memoryLimitMb || 256
                form.status = res.metadata.status || ProblemStatus.PROBLEM_STATUS_UNAVAILABLE
            }
            form.description = res.description || ''
            
            if (res.templates && res.templates.length > 0) {
                form.templates = res.templates
                switchLang(form.templates[0].language) 
            } else {
                form.templates = []
                activeLang.value = ''
                activeFileName.value = ''
            }
        } catch (err) {
            console.error('Fetch detail failed:', err)
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
                // 删除了硬编码的过滤逻辑，直接原样提交管理员设定的所有文件
                templates: form.templates 
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
        activeLang, existingLangs, 
        activeFileName, currentLangFiles, activeTemplate, 
        switchLang, addFile, renameFile, deleteFile,      
        handleFileUpload, handleUpdate
    }
}