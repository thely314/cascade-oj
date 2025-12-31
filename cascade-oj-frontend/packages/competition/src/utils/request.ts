import axios from 'axios';
import router from '@/router'; 

// 1. 创建 axios 实例
const service = axios.create({
  baseURL: '/api', 
  timeout: 10000 
});

// 2. 请求拦截器 (Request Interceptor)
service.interceptors.request.use(
  (config) => {
    // 从 localStorage 获取 token
    // (假设我们在登录成功后，把 token 存到了 'cascade_token' 这个键里)
    const token = localStorage.getItem('cascade_token');

    if (token) {
      // --- 关键点：根据你队友的要求修改这里 ---
      // 不使用 Authorization，而是使用自定义字段 'token'
      // 也不需要 'Bearer ' 前缀
      config.headers['token'] = token;
    }

    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// 3. 响应拦截器 (Response Interceptor)
service.interceptors.response.use(
  (response) => {
    // 如果后端返回状态码是 2xx，直接返回数据
    return response;
  },
  (error) => {
    // TODO 目前的错误信息对用户不友好
    // 如错误将未登录等信息作为 internal error 处理
    // 处理 HTTP 错误状态码
    if (error.response) {
      const status = error.response.status;

      switch (status) {
        case 401:
          // --- 关键点：401 未授权处理 ---
          // 1. 清除本地过期的 token
          localStorage.removeItem('cascade_token');
          // 2. 只有当不在登录页时，才跳转，防止死循环
          if (router.currentRoute.value.path !== '/login') {
            alert('登录已过期，请重新登录'); // 或者使用更优雅的 Toast
            router.push('/login');
          }
          break;
          
        case 403:
          alert('您没有权限执行此操作');
          break;
          
        case 404:
          console.error('请求的资源不存在');
          break;
          
        case 500:
          alert('服务器内部错误，请稍后重试');
          break;
          
        default:
          console.error('网络请求错误:', error.message);
      }
    } else {
      //断网或请求超时
      alert('网络连接异常，请检查网络');
    }
    
    return Promise.reject(error);
  }
);

// 导出这个实例
export default service;