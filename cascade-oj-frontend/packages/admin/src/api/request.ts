const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

interface RequestOptions extends RequestInit {
    params?: Record<string, string | number | boolean | undefined>;
}

async function request<T>(url: string, options: RequestOptions = {}): Promise<T> {
    const { params, ...init } = options;

    let fullUrl = `${BASE_URL}${url}`;
    if (params) {
        const searchParams = new URLSearchParams();
        Object.entries(params).forEach(([key, value]) => {
            if (value !== undefined) {
                searchParams.append(key, String(value));
            }
        });
        const queryString = searchParams.toString();
        if (queryString) {
            fullUrl += `?${queryString}`;
        }
    }

    const headers = new Headers(init.headers);
    if (!headers.has('Content-Type') && !(init.body instanceof FormData)) {
        headers.set('Content-Type', 'application/json');
    }

    const config: RequestInit = {
        ...init,
        headers,
    };

    try {
        const response = await fetch(fullUrl, config);
        if (!response.ok) {
            const errorBody = await response.text();
            throw new Error(`Request failed with status ${response.status}: ${errorBody}`);
        }
        // Assuming JSON response
        const data = await response.json();
        return data as T;
    } catch (error) {
        console.error('API Request Error:', error);
        throw error;
    }
}

export const get = <T>(url: string, params?: Record<string, any>) =>
    request<T>(url, { method: 'GET', params });

export const post = <T>(url: string, body?: any) =>
    request<T>(url, { method: 'POST', body: JSON.stringify(body) });

export const put = <T>(url: string, body?: any) =>
    request<T>(url, { method: 'PUT', body: JSON.stringify(body) });

export const del = <T>(url: string) =>
    request<T>(url, { method: 'DELETE' });
