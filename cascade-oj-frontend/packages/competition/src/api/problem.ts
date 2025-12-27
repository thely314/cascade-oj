import request from '@/utils/request';

// --- 1. 定义数据类型 ---

export interface ProblemFile {
  name: string;
  content: string;
  isReadOnly: boolean; 
}

export interface ProblemDetail {
  id: string;
  title: string;
  creator?: string;
  timeLimit: string;
  memoryLimit: string;
  description: string;
  codeTemplates: {
    [lang: string]: any; 
  };
}

export interface SubmissionResult {
  uuid?: string;
  status: string; // 原本是'accepted' | 'wrong_answer' | 'time_limit_exceeded' | 'compile_error' | 'running' | 'system_error';
  timeCost?: string;
  memoryCost?: string;
  output?: string; 
  stderr?: string;
}

export interface ProblemSimple {
  id: string;
  title: string;
}

export interface SubmitRequest {
  contestId: string;
  problemId: string;
  language: string;
  code: string;      // 原：files: { name: string; content: string }[]; 
  type: 'submit' | 'test';
  input?: string;
}

export interface BackendSubmissionMetadata {
  submissionUuid: string;
  problemId: number;
  userId: number;
  status: string; 
  score: number;
  submitTime: any; 
}

export interface BackendSubmissionReply {
  metadata: BackendSubmissionMetadata;
  code: string;
  language: string;
  // 修正：proto里是int32，所以这里是number
  timeCost: number;   
  memoryCost: number;

  // TODO: 等待后端添加这些字段
  // 队友说 stdout/stderr 被遗漏了，所以这里暂时没有
  // 我们不需要定义它们，直接在 adapter 里处理
  stdout?: string; 
  stderr?: string; 
}

// --- 2. 模拟真实请求的延时函数 ---
const delay = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

// 开关：是否使用模拟数据
const IS_MOCK = false;

// --- 3. API 方法 ---

// A. 获取题目详情
export const fetchProblemDetail = async (id: string): Promise<ProblemDetail> => {
  if (IS_MOCK) {
    await delay(300);
    
    // 模拟一个复杂的 C++ 题目：需要补充头文件，只读主文件
    return {
      id: id,
      title: '1. 长方体类 (Cuboid)', // 对应你截图的例子
      timeLimit: '1000ms',
      memoryLimit: '128MB',
      description: `### 题目描述\n请实现一个长方体类 Cuboid...`,
      codeTemplates: {
        'C++': [
          { 
            name: 'main.cpp', 
            content: `#include <iostream>\n#include "cuboid.hpp"\nusing namespace std;\n\nint main() {\n    Cuboid c(3, 4, 5);\n    cout << "Area: " << c.get_area() << endl;\n    return 0;\n}`, 
            isReadOnly: true 
          },
          { 
            name: 'cuboid.hpp', 
            content: `#pragma once\n\nclass Cuboid {\nprivate:\n    int length, width, height;\npublic:\n    Cuboid(int l, int w, int h);\n    double get_area();\n    double get_volume();\n};`, 
            isReadOnly: false 
          }
        ],
        'Java': [
           { name: 'Main.java', content: 'public class Main { ... }', isReadOnly: false }
        ]
      }
    };
  }
  // 真实接口
  // 注意：确保 request.ts 里的 baseURL 已经包含了 /api (或根据后端代理配置)
  const res = await request.get(`/user/problems/${id}`); 
  const backendData = res.data; // 假设对应 GetSingleProblemReply

  return {
    id: String(backendData.metadata.id),
    title: backendData.metadata.title,
    creator: backendData.creator || 'Admin', 
    timeLimit: `${backendData.metadata.timeLimitMs}ms`,
    memoryLimit: `${backendData.metadata.memoryLimitMb}MB`,
    description: backendData.description,
    codeTemplates: {} // 后端暂无模板，留空
  };
};

// B. 获取测试运行结果
export const getSelfTestResult = async (uuid: string): Promise<SubmissionResult> => {
  if (IS_MOCK) {
    // 模拟：随机返回 Running 或 完成
    // 实际 Mock 时，为了效果，我们可以第一次调用返回 Running，第二次返回成功
    // 这里简单处理：直接返回成功
    return {
      status: 'Accepted',
      timeCost: '2ms',
      memoryCost: '9216KB',
      output: 'Mock Output: Hello World',
      stderr: ''
    };
  }

  const res = await request.get(`/user/selftests/${uuid}`);
  const data = res.data; // 对应 GetSelfTestResultReply

  // 适配器 logic
  // 注意：proto 里 SelfTest 只有一个 is_compiled 字段
  // 我们假设：如果 !is_compiled 且无报错，算 Running；否则算结束
  // *具体状态判断可能需要跟后端确认，这里暂且认为只要有返回就是结束*
  
  return {
    status: 'Finished', 
    timeCost: data.timeCost ? `${data.timeCost}ms` : '0ms',
    memoryCost: data.memoryCost ? `${data.memoryCost}KB` : '0KB',
    output: data.stdout || '',
    stderr: data.stderr || ''
  };
};


// C. 查询正式提交结果
export const getSubmissionResult = async (uuid: string): Promise<SubmissionResult> => {
  if (IS_MOCK) {
    return { status: 'Accepted', timeCost: '12ms', memoryCost: '1.2MB' };
  }

  const res = await request.get(`/user/submissions/${uuid}`);
  const data = res.data; // 对应 GetSingleSubmissionReply

  return {
    status: data.metadata?.status || 'Unknown',
    timeCost: data.timeCost ? `${data.timeCost}ms` : '0ms',
    memoryCost: data.memoryCost ? `${data.memoryCost}KB` : '0KB',
    output: '', // 正式提交通常不看 output
    stderr: '' // 如果有编译错误，可能在 metadata.status 或其他字段
  };
};

// D. 提交代码 (或测试运行)
export const submitCode = async (data: SubmitRequest): Promise<SubmissionResult> => {
  
  if (IS_MOCK) {
    await delay(500);
    return { 
      uuid: 'mock-uuid-12345', 
      status: 'Running', 
      output: '正在提交...' 
    };
  }

  let res;
  if (data.type === 'test') {
    const payload = {
      problemId: Number(data.problemId),
      code: data.code,
      language: data.language,
      input: data.input ? data.input : ""
    };
    res = await request.post('/user/selftests', payload);
  } else {
    const payload = {
      contestId: Number(data.contestId),
      problemId: Number(data.problemId),
      code: data.code,
      language: data.language,
    };
    res = await request.post('/user/submissions', payload);
  }
  // const payload = {
  //   // contest_id: Number(data.contestId),
  //   problemId: Number(data.problemId),
  //   code: data.code,
  //   language: data.language,
  //   input: data.input ? data.input : ""
  // };

  // const res = await request.post(
  //   data.type === 'test' ? '/user/selftests' : '/user/submissions', 
  //   payload
  // );
  
  return {
    uuid: res ? res.data.uuid : '',
    status: 'pending', 
    output: '请求已发送，等待结果...'
  };
};

// E. 获取简易题目列表
export const fetchProblemList = async (contestId: string): Promise<ProblemSimple[]> => {
  if (IS_MOCK) {
    return [
      { id: '1', title: '两数之和' },
      { id: '2', title: '两数相加' },
      { id: '3', title: '无重复字符的最长子串' },
      { id: '4', title: '寻找两个正序数组的中位数' },
      { id: '5', title: '最长回文子串' },
      { id: '6', title: 'N 字形变换' },
      { id: '7', title: '整数反转' },
      { id: '8', title: '字符串转换整数 (atoi)' },
    ];
  }
  // 真实接口
  // 注意：GetProblems 需要 contest_id。这里如果做公共题库，需确认 contest_id 传什么
  // 假设暂时获取 ID=1 的比赛题目列表
  const res = await request.get(`/user/contests/${contestId}/problems`); 
  return res.data.problems; // 假设返回结构里有 problems 数组
};


// F. 提交记录相关 

// 提交记录列表项 (对应后端 SubmissionMetadata)
export interface SubmissionItem {
  submissionUuid: string;
  status: string;      
  submitTime: string; // ISO 时间字符串
  score: number; 
  //TODO:      
  // 后端 proto 里 GetSubmissionsReply 列表项似乎没有 time/memory/language？
  // 如果没有，暂时只能展示状态。如果有扩展，这里补上。
  // 通常列表页也需要展示语言、耗时、内存，假设后端之后会补，我们先 Mock 出来
  language?: string;
  timeCost?: number;
  memoryCost?: number;
}

// 获取提交记录列表
export const fetchSubmissions = async (contestId: string, problemId: string): Promise<SubmissionItem[]> => {
  if (IS_MOCK) {
    await delay(300);
    // Mock 数据
    return [
      { submissionUuid: 'sub-001', status: 'Accepted', score: 100, submitTime: '2023-10-01T12:00:00Z', language: 'C++', timeCost: 12, memoryCost: 1024 },
      { submissionUuid: 'sub-002', status: 'Wrong Answer', score: 0, submitTime: '2023-10-01T11:55:00Z', language: 'Python3', timeCost: 20, memoryCost: 2048 },
      { submissionUuid: 'sub-003', status: 'Time Limit Exceeded', score: 0, submitTime: '2023-10-01T11:50:00Z', language: 'Java', timeCost: 1000, memoryCost: 5000 },
    ];
  }

  // 真实请求
  // query 参数: problem_id, contest_id
  const res = await request.get('/user/submissions', {
    params: {
      problemId: problemId,
      contestId: contestId,
      page: 1,      // 暂时写死第一页
      pageSize: 20
    }
  });
  
  // 适配器：如果有字段不一致，在这里转换
  return res.data.submissions || [];
};