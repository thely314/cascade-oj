import { ref, reactive, computed, onMounted, onBeforeUnmount, watch } from 'vue';
import { useRoute, useRouter, onBeforeRouteLeave, onBeforeRouteUpdate } from 'vue-router';
import { marked } from 'marked';
import markedKatex from "marked-katex-extension";
import "katex/dist/katex.min.css"; // 防止公式乱码
import { 
  fetchProblemDetail, 
  fetchProblemList,
  submitCode, 
  getSelfTestResult, 
  getSubmissionResult,
  type ProblemDetail,
  type ProblemSimple,
  type SubmissionResult
} from '@/api/problem';

export function useProblemDetail() {
  const route = useRoute();
  const router = useRouter();

  // 获取 URL 参数 (注意：如果用户直接输 URL 进来，这两个可能是 undefined，最好做个兜底)
  const contestId = computed(() => route.params.contestId as string || '1');
  const problemId = computed(() => route.params.id as string || '1');

  // --- 配置 marked 使用 KaTeX ---
  marked.use(markedKatex({ 
    throwOnError: false, 
    output: 'html' 
  } as any));

  // --- 状态管理 ---
  const leftTabs = ['描述', '提交', /*'成绩', '笔记',*/ '排名'/*, '答案'*/];
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
  const supportedLanguages = ['c', 'c++11', 'c++11(O2)'];
  const currentLanguage = ref('c');

  // 测试运行的详细结果状态
  const runStats = reactive({
    time: '-',
    memory: '-',
    stderr: '' 
  });

  // 通用轮询函数
  const pollResult = async (
    uuid: string, 
    apiFunc: (id: string) => Promise<SubmissionResult>,
    onSuccess: (res: SubmissionResult) => void,
    onError: (err: any) => void
  ) => {
    const MAX_RETRIES = 20; // 最大轮询次数 (防止死循环)
    let count = 0;

    const intervalId = setInterval(async () => {
      count++;
      try {
        const res = await apiFunc(uuid);
        
        // 判断是否结束 (后端状态不是 pending 或 judging)
        // 注意：根据后端实际返回的状态字符串调整这里
        if (res.status !== 'pending' && res.status !== 'judging') {
          clearInterval(intervalId); // 停止轮询
          onSuccess(res); // 回调成功
        }
      } catch (e) {
        clearInterval(intervalId);
        onError(e);
      }

      // 超时处理
      if (count >= MAX_RETRIES) {
        clearInterval(intervalId);
        onError(new Error("请求超时，请稍后重试"));
      }
    }, 1000); // 每 1 秒查一次
  };

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
  
  // [修改点1] filesMap 硬编码初始化 (仿洛谷/CF风格)
  // 既然后端暂时不给模板，这里写死默认状态
  const filesMap = reactive<Record<string, CodeFile[]>>({
    'c': [
      { name: 'main.c', lang: 'c', content: '', isReadOnly: false }
    ],
    'c++11': [
      { name: 'main.cpp', lang: 'c++11', content: '', isReadOnly: false }
    ],
    'c++11(O2)': [
      { name: 'main.cpp', lang: 'c++11(O2)', content: '', isReadOnly: false }
    ],
    'Java': [
      { name: 'Main.java', lang: 'java', content: 'public class Main {\n    public static void main(String[] args) {\n\n    }\n}', isReadOnly: false }
    ],
    'Python3': [
      { name: 'solution.py', lang: 'python', content: '', isReadOnly: false }
    ]
  });

  // 代码快照，用于检测未保存的更改
  const originalCodeSnapshot = reactive<Record<string, string>>({});

  // 当前显示的文件列表（根据语言变化）
  const currentFiles = computed(() => filesMap[currentLanguage.value] || []);
  
  // 当前选中的文件名 (切换语言时默认选中第一个)
  const currentFileName = ref('main.cpp'); // 给个默认值防止为空
  
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
  
  // 更新代码快照
  const updateCodeSnapshot = () => {
    for (const lang of Object.keys(filesMap)) {
      for (const file of filesMap[lang]) {
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
        if (!file.isReadOnly && originalContent !== undefined && file.content !== originalContent) {
          return true;
        }
      }
    }
    return false;
  };

  // 抽离统一的检查逻辑
  const checkUnsavedChange = (next: Function, cancel: Function) => {
    if (hasUnsavedChanges()) {
      const confirmLeave = globalThis.confirm('您有未保存的代码更改。确定要切题/离开吗？更改将丢失。');
      if (confirmLeave) next();
      else cancel();
    } else {
      next();
    }
  };

  // 离开页面
  onBeforeRouteLeave((to, from, next) => {
    checkUnsavedChange(() => next(), () => next(false));
  });

  // 同组件路由更新 (切题)
  onBeforeRouteUpdate((to, from, next) => {
    checkUnsavedChange(() => next(), () => next(false));
  });

  // --- 核心方法 ---

  // 加载单题详情
  const loadProblem = async (id: string) => {
    try {
      loading.value = true;
      const data = await fetchProblemDetail(id);
      problemData.value = data;

      // [修改点2] 删除了后端模板覆盖逻辑
      // 因为后端没做这个功能，保留原有的硬编码 filesMap 即可
      // 只需更新快照，视为当前状态就是“初始状态”
      
      // 加载完成后，立即更新快照
      updateCodeSnapshot();
      
    } catch (error) {
      console.error("题目加载失败", error);
      problemData.value.title = "题目加载失败";
    } finally {
      loading.value = false;
    }
  };

  // 新增一个 map 来存储 序号 -> 真实ID 的映射
  const problemIndexMap = reactive<Record<string, string>>({});
  // 反向映射
  const problemIdToIndex = reactive<Record<string, string>>({});

  // 加载题目列表
  const loadProblemList = async () => {
    try {
      const list = await fetchProblemList(contestId.value); 
      problemList.value = list;
      totalProblems.value = list.length;
      
      list.forEach((p, index) => {
        // 假设 URL 里的 1 代表数组第 0 个
        const routeIndex = String(index + 1); 
        problemIndexMap[routeIndex] = p.id;
        problemIdToIndex[p.id] = routeIndex;
      });
      
      // 如果当前路由参数是序号，转化为真实 ID 后再加载详情
      const routeId = route.params.id as string;
      const realId = problemIndexMap[routeId];
      
      if (realId) {
        loadProblem(realId); // 用真实 ID 去查
      } else {
        // 可能是直接传了真实 ID，或者是无效序号
        loadProblem(routeId);
      }
      
    } catch (e) {
      console.error("题目列表加载失败", e);
    }
  };

  // 跳转到指定题目
  const jumpToProblem = (id: string) => {
    showProblemDrawer.value = false; 
    // [修改点3] 修正路由跳转，必须带上 contestId
    // 假设路由名配置为 'ProblemDetail'
    router.push({
      name: 'ProblemDetail',
      params: { 
        contestId: contestId.value, 
        id: id 
      }
    });
  };

  // 总题数：以接口返回为准，初始为 0
  const totalProblems = ref(0); 
  const currentIndex = computed(() => Number(problemIdToIndex[problemData.value.id]) || 1);

  const isFirstProblem = computed(() => currentIndex.value <= 1);
  const isLastProblem = computed(() => currentIndex.value >= totalProblems.value);

  // 上一题
  const handlePrevProblem = () => {
    if (isFirstProblem.value) return;
    // const currentId = Number.parseInt(problemData.value.id);
    // const currentIndex = Number.parseInt(problemIdToIndex[problemData.value.id]);
    // jumpToProblem(String(currentId - 1));
    jumpToProblem(problemIndexMap[String(currentIndex.value - 1)]);
  };

  // 下一题
  const handleNextProblem = () => {
    if (isLastProblem.value) return;
    // const currentId = Number.parseInt(problemData.value.id);
    // const currentIndex = Number.parseInt(problemIdToIndex[problemData.value.id]);
    // jumpToProblem(String(currentId + 1));
    jumpToProblem(problemIndexMap[String(currentIndex.value + 1)]);
  };

  // 测试运行 (集成轮询)
  const handleTestRun = async () => {
    if (isRunning.value) return;
    isRunning.value = true;
    playgroundOutput.value = "正在提交代码...";
    runStats.time = '-';
    runStats.memory = '-';
    runStats.stderr = '';

    try {
      const codeToSend = activeFile.value?.content || '';

      // 发起提交，获取 UUID
      const res = await submitCode({
        contestId: contestId.value, // 使用 computed
        problemId: problemData.value.id,
        code: codeToSend,
        language: currentLanguage.value,
        type: 'test',
        input: playgroundInput.value
      });

      if (res.uuid) {
        playgroundOutput.value = "正在运行中...";
        
        // 开始轮询
        await pollResult(
          res.uuid,
          getSelfTestResult, // 使用查询自测的 API
          (finalRes) => {
            // 轮询结束，更新 UI
            isRunning.value = false;
            runStats.time = finalRes.timeCost || '0ms';
            runStats.memory = finalRes.memoryCost || '0KB';
            runStats.stderr = finalRes.stderr || ''; // 这里假设 errorMsg 存的是 stderr
            if (finalRes.status === 'compile_error') {
              playgroundOutput.value = `=== 编译错误 ===\n${finalRes.stderr}`;
            } else {
              runStats.time = finalRes.timeCost || '0ms';
              runStats.memory = finalRes.memoryCost || '0KB';
              // 填充 stderr 注：api/problem.ts 的 submitCode 适配器里，把 backendData.stderr 映射到了 errorMsg，所以这里取 res.errorMsg
              runStats.stderr = finalRes.stderr || ''; 
              playgroundOutput.value = finalRes.output || '程序无输出';
            }
          },
          (err) => {
            isRunning.value = false;
            playgroundOutput.value = "查询结果失败: " + err.message;
          }
        );
      } else {
        // 如果没有 UUID，可能是 Mock 模式或者出错了
        isRunning.value = false;
        playgroundOutput.value = "未获取到运行 ID";
      }

    } catch (e) {
      console.error("测试运行出错:", e);
      playgroundOutput.value = "系统错误";
      isRunning.value = false;
    }
  };

  // 提交代码(集成轮询)
  const handleSubmit = async () => {
    if (isSubmitting.value) return;
    isSubmitting.value = true;
    
    try {
      const codeToSend = activeFile.value?.content || '';

      // 1. 发起提交
      const res = await submitCode({
        contestId: contestId.value, 
        problemId: problemData.value.id,
        code: codeToSend,
        language: currentLanguage.value,
        type: 'submit'
      });
      
      updateCodeSnapshot();

      if (res.uuid) {
        // 2. 开始轮询
        // 这里我们可以做一个简单的 UI 反馈，比如按钮上显示 "判题中..."
        // 或者直接弹窗显示 "Waiting..."
        
        await pollResult(
          res.uuid,
          getSubmissionResult, // 使用查询提交的 API
          (finalRes) => {
            isSubmitting.value = false;
            submissionResult.value = finalRes;
            showResultModal.value = true;
          },
          (err) => {
            isSubmitting.value = false;
            alert("判题超时或失败: " + err.message);
          }
        );
      } else {
        isSubmitting.value = false;
        alert("提交失败，未获取到 ID");
      }

    } catch (e) {
      console.error("提交代码出错:", e);
      alert("提交失败，网络错误");
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
  
  // 路由拦截
  onBeforeRouteLeave((to, from, next) => {
    checkUnsavedChange(() => next(), () => next(false));
  });

  // 浏览器刷新/关闭拦截
  const handleBeforeUnload = (e: BeforeUnloadEvent) => {
    if (hasUnsavedChanges()) {
      e.preventDefault(); 
      (e as any).returnValue = '';
    }
  };

  // --- 生命周期 ---
  onMounted(() => {
    const id = problemId.value; // 从 computed 取值
    loadProblemList(); 
    globalThis.addEventListener('mousemove', onMouseMove);
    globalThis.addEventListener('beforeunload', handleBeforeUnload);
  });
  
  onBeforeUnmount(() => {
    globalThis.removeEventListener('beforeunload', handleBeforeUnload);
  });

  // 监听题目 ID 变化
  watch(() => route.params.id, (newId) => {
    if (newId) {
      loadProblem(newId as string);
    }
  });

  // [新增] 复制代码
  const copyCode = async () => {
    try {
      await navigator.clipboard.writeText(activeFile.value.content);
      alert('代码已复制到剪贴板'); // 这里可以用个 Toast，暂时用 alert 替代
    } catch (err) {
      console.error('复制失败', err);
    }
  };

  // [新增] 重置代码
  const resetCode = () => {
    if (confirm('确定要重置当前文件代码吗？您的修改将丢失。')) {
      // 简单粗暴：直接清空，或者恢复成特定模板
      // 如果之前 filesMap 初始化时有默认模板，这里可以恢复成那个值
      // 目前我们就置空，或者给个基础框架
      activeFile.value.content = ''; 
    }
  };

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
    currentIndex, isFirstProblem, isLastProblem, totalProblems,
    showResultModal, submissionResult,
    runStats,
    contestId,
    copyCode, 
    resetCode 
  };
}