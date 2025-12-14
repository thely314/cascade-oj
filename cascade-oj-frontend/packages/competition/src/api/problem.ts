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
  timeLimit: string;
  memoryLimit: string;
  description: string;
  codeTemplates: {
    [lang: string]: ProblemFile[]; 
  };
}

export interface SubmissionResult {
  status: 'Accepted' | 'Wrong Answer' | 'Time Limit Exceeded' | 'Compile Error' | 'Running' | 'System Error';
  time?: string;
  memory?: string;
  output?: string; 
  errorMsg?: string;
}

export interface ProblemSimple {
  id: string;
  title: string;
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
      title: '3. 长方体类 (Cuboid)', // 对应你截图的例子
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
  const res = await request.get(`/problems/${id}`);
  return res.data;
};

// B. 提交代码 (或测试运行)
export const submitCode = async (data: {
  problemId: string;
  language: string;
  files: { name: string; content: string }[]; 
  type: 'submit' | 'test'; 
  input?: string;
}): Promise<SubmissionResult> => {
  
  if (IS_MOCK) {
    await delay(800);
    
    // 模拟测试运行成功
    if (data.type === 'test') {
      return {
        status: 'Accepted',
        time: '2ms',
        memory: '9216KB',
        output: 'Cuboid Information:\nlength: 3\nwidth: 4\nheight: 5\narea: 94\nvolume: 60'
      };
    } 
    
    // 模拟提交
    return { status: 'Accepted', time: '12ms', memory: '1.2MB' };
  }
  const res = await request.post('/submissions', data);
  return res.data;
};

// C.获取简易题目列表
export const fetchProblemList = async (): Promise<ProblemSimple[]> => {
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
  const res = await request.get('/problems/simple-list');
  return res.data;
};