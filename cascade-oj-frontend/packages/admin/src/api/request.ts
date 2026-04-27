import axios from 'axios';
import router from '@/router';

// 1. 创建 axios 实例
const service = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
    timeout: 10000
});

// 2. 请求拦截器 (Request Interceptor)
service.interceptors.request.use(
    (config) => {
        // 从 localStorage 获取 token
        const token = localStorage.getItem('cascade_token');

        if (token) {
            // 遵循队友要求的自定义 'token' 字段
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
        // Axios 的响应对象中，后端返回的实际数据在 .data 中
        // 之前使用 Fetch 时，.json() 返回的就是后端数据对象
        // 为了保持 admin.ts 的泛型正常工作及 Users.vue 能拿到 response.users，
        // 我们应该在这里直接返回 response.data
        return response.data;
    },
    (error) => {
        // 处理 HTTP 错误状态码
        if (error.response) {
            const status = error.response.status;

            switch (status) {
                case 401:
                    // 401 未授权处理
                    localStorage.removeItem('cascade_token');
                    // 只有当不在登录页时，才跳转，防止死循环
                    if (router.currentRoute.value.path !== '/login') {
                        alert('登录已过期或权限不足，请重新登录');
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
            // 断网或请求超时
            alert('网络连接异常，请检查网络');
        }

        return Promise.reject(error);
    }
);

// 定义泛型请求方法，保持与原 admin 逻辑的兼容性
export const get = <T>(url: string, params?: Record<string, any>) =>
    service.get<any, T>(url, { params });

export const post = <T>(url: string, body?: any) =>
    service.post<any, T>(url, body);

export const put = <T>(url: string, body?: any) =>
    service.put<any, T>(url, body);

export const del = <T>(url: string) =>
    service.delete<any, T>(url);

export default service;
