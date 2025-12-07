import { ref, reactive, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useRoute, useRouter, onBeforeRouteLeave } from 'vue-router';
import { marked } from 'marked';
import markedKatex from "marked-katex-extension";
import "katex/dist/katex.min.css"; // 防止公式乱码
import { 
  fetchProblemDetail, 
  fetchProblemList,
  submitCode, 
  type ProblemDetail,
  type ProblemSimple,
  type SubmissionResult
} from '@/api/problem';

export function useProblemDetail() {
  const route = useRoute();
  const router = useRouter();

  // --- 配置 marked 使用 KaTeX ---
  marked.use(markedKatex({ 
    throwOnError: false, // 如果公式写错，不要报错崩掉页面，而是显示源码
    output: 'html' 
  } as any));

  // --- 状态管理 ---
  const leftTabs = ['描述', '提交', '成绩', '笔记', '排名', '答案'];
  const currentLeftTab = ref(0);
  
  // 界面状态
  const isPlaygroundOpen = ref(false);
  const loading = ref(true);
  const isSubmitting = ref(false);
  const isRunning = ref(false);
  
  // 题目列表抽屉状态
  const showProblemDrawer = ref(false);
  const problemList = ref<ProblemSimple[]>([]);

  // 提交结果模态框状态
  const showResultModal = ref(false);
  const submissionResult = ref<SubmissionResult | null>(null);

  // 支持的语言列表
  const supportedLanguages = ['C++', 'Java', 'Python3'];
  const currentLanguage = ref('C++');

  // 测试运行的详细结果状态
  const runStats = reactive({
    time: '-',
    memory: '-'
  });

  // 题目数据
  const problemData = ref<ProblemDetail>({
    id: '',
    title: '加载中...',
    timeLimit: '-',
    memoryLimit: '-',
    description: '',
    codeTemplates: {}
  });

  // --- 动态代码文件结构 ---
  // 定义代码文件接口
  interface CodeFile {
    name: string;
    lang: string; 
    content: string;
    isReadOnly: boolean;
  }
  
  // filesMap 定义结构
  const filesMap = reactive<Record<string, CodeFile[]>>({});

  // [新增] 代码快照，用于检测未保存的更改
  const originalCodeSnapshot = reactive<Record<string, string>>({});

  // 当前显示的文件列表（根据语言变化）
  const currentFiles = computed(() => filesMap[currentLanguage.value] || []);
  
  // 当前选中的文件名 (切换语言时默认选中第一个)
  const currentFileName = ref(''); 
  
  // 监听语言变化，自动选中第一个文件
  watch(currentLanguage, (newLang) => {
    if (filesMap[newLang] && filesMap[newLang].length > 0) {
      currentFileName.value = filesMap[newLang][0].name;
    }
  }, { immediate: true });

  // 计算当前文件对象
  const activeFile = computed(() => {
    return currentFiles.value.find(f => f.name === currentFileName.value) || currentFiles.value[0];
  });

  // 编辑器双向绑定代码内容
  const currentCode = computed({
    get: () => activeFile.value?.content || '',
    set: (val) => { 
      if (activeFile.value && !activeFile.value.isReadOnly) {
        activeFile.value.content = val; 
      }
    }
  });

  // 是否只读
  const isCurrentReadOnly = computed(() => activeFile.value?.isReadOnly || false);

  // Playground 输入输出
  const playgroundInput = ref('');
  const playgroundOutput = ref('等待测试运行...');

  // --- 计算 Markdown 渲染后的 HTML ---
  const descriptionHtml = computed(() => marked.parse(problemData.value.description));

  // --- 辅助方法：快照管理 ---
  
  // 更新代码快照 (在加载完成或提交成功后调用)
  const updateCodeSnapshot = () => {
    for (const lang of Object.keys(filesMap)) {
      for (const file of filesMap[lang]) {
        // 记录唯一 Key: 语言_文件名
        originalCodeSnapshot[`${lang}_${file.name}`] = file.content;
      }
    }
  };
  // 检查是否有未保存的更改
  const hasUnsavedChanges = () => {
    for (const lang of Object.keys(filesMap)) {
      for (const file of filesMap[lang]) {
        const key = `${lang}_${file.name}`;
        const originalContent = originalCodeSnapshot[key];
        // 如果快照存在且内容不一致 (排除只读文件)
        if (!file.isReadOnly && originalContent !== undefined && file.content !== originalContent) {
          return true;
        }
      }
    }
    return false;
  };

  // --- 核心方法 ---

  // 加载单题详情
  const loadProblem = async (id: string) => {
    try {
      loading.value = true;
      const data = await fetchProblemDetail(id);
      problemData.value = data;

      // 清空旧文件映射
      for (const key in filesMap) delete filesMap[key];

      // 将后端返回的文件列表映射到前端结构
      if (data.codeTemplates) {
        for (const lang of Object.keys(data.codeTemplates)) {
          const apiFiles = data.codeTemplates[lang];
          
          filesMap[lang] = apiFiles.map(f => ({
            name: f.name,
            // 简单的语言映射逻辑
            lang: lang.toLowerCase().includes('python') ? 'python' : lang.toLowerCase(),
            content: f.content,
            isReadOnly: f.isReadOnly
          }));
        }
        
        // 刷新当前视图：确保有文件显示
        if (!filesMap[currentLanguage.value]) {
           // 如果当前语言没有模板，切换到第一个可用的语言
           const firstLang = Object.keys(filesMap)[0];
           if (firstLang) currentLanguage.value = firstLang;
        }
        // 选中第一个文件
        if (filesMap[currentLanguage.value] && filesMap[currentLanguage.value].length > 0) {
           currentFileName.value = filesMap[currentLanguage.value][0].name;
        }

        // 加载完成后，立即更新快照
        updateCodeSnapshot();
      }
    } catch (error) {
      console.error("题目加载失败", error);
      problemData.value.title = "题目加载失败";
    } finally {
      loading.value = false;
    }
  };

  // 加载题目列表 (用于左下角抽屉)
  const loadProblemList = async () => {
    try {
      const list = await fetchProblemList();
      problemList.value = list;
    } catch (e) {
      console.error("题目列表加载失败", e);
    }
  };

  // 跳转到指定题目
  const jumpToProblem = (id: string) => {
    showProblemDrawer.value = false; 
    router.push(`/problem/${id}`);
  };

  // 上一题 
  const handlePrevProblem = () => {
    const currentId = Number.parseInt(problemData.value.id);
    if (!Number.isNaN(currentId) && currentId > 1) {
      jumpToProblem(String(currentId - 1));
    } else {
      alert("已经是第一题了");
    }
  };

  // 下一题
  const handleNextProblem = () => {
    const currentId = Number.parseInt(problemData.value.id);
    if (!Number.isNaN(currentId)) {
      jumpToProblem(String(currentId + 1));
    }
  };

  // 测试运行
  const handleTestRun = async () => {
    if (isRunning.value) return;
    isRunning.value = true;
    playgroundOutput.value = "正在运行...";
    // 重置统计数据
    runStats.time = '-';
    runStats.memory = '-';

    try {
      const filesToSend = filesMap[currentLanguage.value].map(f => ({
         name: f.name,
         content: f.content
      }));

      const res = await submitCode({
        problemId: problemData.value.id,
        language: currentLanguage.value,
        files: filesToSend, // 发送所有文件
        type: 'test',
        input: playgroundInput.value
      });
      
      if (res.status === 'Compile Error') {
         playgroundOutput.value = `=== 编译错误 ===\n${res.errorMsg}`;
      } else {
         // 填充时间和内存
         runStats.time = res.time || '0ms';
         runStats.memory = res.memory || '0KB';
         playgroundOutput.value = res.output || '程序无输出';
      }
    } catch (e) {
      console.error("测试运行出错:", e);
      playgroundOutput.value = "系统错误";
    } finally {
      isRunning.value = false;
    }
  };

  // 提交代码
  const handleSubmit = async () => {
    if (isSubmitting.value) return;
    isSubmitting.value = true; 
    try {
      const filesToSend = currentFiles.value.map(f => ({
         name: f.name,
         content: f.content
      }));

      const res = await submitCode({
        problemId: problemData.value.id,
        language: currentLanguage.value,
        files: filesToSend,
        type: 'submit'
      });
      
      // 提交成功后，视为已保存，更新快照
      updateCodeSnapshot();

      submissionResult.value = res;
      showResultModal.value = true;
    } catch (e) {
      console.error("提交代码出错:", e);
      alert("提交失败，网络错误");
    } finally {
      isSubmitting.value = false;
    }
  };

  // --- 拖拽调整宽度逻辑 ---
  const leftWidth = ref(50); 
  const isDragging = ref(false);

  const startDrag = () => {
    isDragging.value = true;
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
  };

  const stopDrag = () => {
    isDragging.value = false;
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
  };

  const onMouseMove = (e: MouseEvent) => {
    if (!isDragging.value) return;
    const containerWidth = globalThis.innerWidth;
    let newLeftWidth = (e.clientX / containerWidth) * 100;
    if (newLeftWidth < 20) newLeftWidth = 20;
    if (newLeftWidth > 80) newLeftWidth = 80;
    leftWidth.value = newLeftWidth;
  };

  // --- [新增] 路由守卫与事件监听 ---
  
  // 1. 路由拦截 (Vue Router)
  onBeforeRouteLeave((to, from, next) => {
    if (hasUnsavedChanges()) {
      const confirmLeave = globalThis.confirm('您有未保存的代码更改。确定要离开吗？您的更改将会丢失。');
      if (confirmLeave) {
        next(); 
      } else {
        next(false); 
      }
    } else {
      next(); 
    }
  });

  // 2. 浏览器刷新/关闭拦截
  const handleBeforeUnload = (e: BeforeUnloadEvent) => {
    if (hasUnsavedChanges()) {
      e.preventDefault(); 
      (e as any).returnValue = '';
    }
  };

  // --- 生命周期 ---
  onMounted(() => {
    const id = route.params.id as string || '1';
    loadProblem(id);
    loadProblemList(); // 获取列表数据
    globalThis.addEventListener('mousemove', onMouseMove);
    // [新增] 监听浏览器关闭/刷新
    globalThis.addEventListener('beforeunload', handleBeforeUnload);
  });
  
  onBeforeUnmount(() => {
    // [新增] 销毁监听器
    globalThis.removeEventListener('beforeunload', handleBeforeUnload);
  });

  // 监听路由变化 (题目 ID 变化)
  watch(() => route.params.id, (newId) => {
    if (newId) {
      // 路由切换时，onBeforeRouteLeave 已经处理了确认逻辑
      // 这里只需要直接加载新题目即可
      loadProblem(newId as string);
    }
  });

  // 导出给模板使用
  return {
    leftTabs, currentLeftTab, isPlaygroundOpen, loading, isSubmitting, isRunning,
    problemData, descriptionHtml,
    currentFiles, currentFileName, currentCode, isCurrentReadOnly,
    supportedLanguages, currentLanguage, 
    playgroundInput, playgroundOutput,
    leftWidth, startDrag, stopDrag,
    handleTestRun, handleSubmit,
    showProblemDrawer, problemList,
    jumpToProblem, handlePrevProblem, handleNextProblem,
    showResultModal, submissionResult,
    runStats // 导出统计数据
  };
}