import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { postProblem } from '../../api/admin'
import { ProblemStatus, CodeTemplate } from '../../api/types' // Ensure CodeTemplate is imported

export function useNewProblem() {
    const router = useRouter()
    const loading = ref(false)

    // 1. Dynamic language and file state
    const supportedLangs = ['cpp', 'java', 'python', 'go']
    const activeLang = ref('')
    const activeFileName = ref('')

    const form = reactive({
        title: '',
        timeLimitMs: 1000,
        memoryLimitMb: 256,
        description: '', 
        templates: [] as CodeTemplate[] // [REVIEW FIX]: Empty by default to support ACM mode
    })

    // 2. Computed properties for dynamic tabs
    const existingLangs = computed(() => {
        const langs = new Set(form.templates.map(t => t.language))
        return Array.from(langs)
    })

    const currentLangFiles = computed(() => {
        return form.templates.filter(t => t.language === activeLang.value)
    })

    const activeTemplate = computed(() => {
        return form.templates.find(
            t => t.language === activeLang.value && t.name === activeFileName.value
        )
    })

    // 3. Dynamic template management functions
    const switchLang = (lang: string) => {
        activeLang.value = lang
        const files = form.templates.filter(t => t.language === lang)
        activeFileName.value = files.length > 0 ? files[0].name : ''
    }

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
            if (activeFileName.value === name) {
                const remaining = form.templates.filter(t => t.language === activeLang.value)
                activeFileName.value = remaining.length > 0 ? remaining[0].name : ''
            }
        }
    }

    // Save problem
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
                    status: ProblemStatus.PROBLEM_STATUS_UNAVAILABLE,
                    description: form.description // Ensure metadata description is aligned
                },
                description: form.description,
                templates: form.templates // Send dynamic templates
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
        router,
        activeLang,
        supportedLangs,
        activeFileName,
        existingLangs,
        currentLangFiles,
        activeTemplate,
        switchLang,
        addFile,
        renameFile,
        deleteFile
    }
}