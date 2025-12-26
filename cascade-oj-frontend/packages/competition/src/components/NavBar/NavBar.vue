<template>
  <nav class="navbar">
    <div class="nav-container">
      <!-- 网站Logo -->
      <div class="nav-logo">
        <router-link class="logo-link" to="/home" aria-label="CasCade Home">
          <!-- inline cascade svg -->
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 164 164">
            <title>CasCade</title>
            <defs/>
            <g>
              <g data-cell-id="0">
                <g data-cell-id="1">
                  <g data-cell-id="-t6uBNjXSVFn4JsbTugf-1">
                    <g transform="translate(0.5,0.5)"/>
                    <g data-cell-id="-t6uBNjXSVFn4JsbTugf-2">
                      <g transform="translate(0.5,0.5)">
                        <path d="M 74 34 L 94 34 Q -26 74 94 114 L 74 114 Q -46 74 74 34 Z" fill="#4bedc4" stroke="none" transform="rotate(-30,54,74)" style="fill: linear-gradient(#4bedc4,#18cb9e);"/>
                      </g>
                    </g>
                    <g data-cell-id="-t6uBNjXSVFn4JsbTugf-3">
                      <g transform="translate(0.5,0.5)">
                        <path d="M 128 48 L 148 48 Q 28 88 148 128 L 128 128 Q 8 88 128 48 Z" fill="#4bedc4" stroke="none" transform="rotate(-210,108,88)" style="fill: linear-gradient(#4bedc4,#18cb9e);"/>
                      </g>
                    </g>
                  </g>
                </g>
              </g>
            </g>
          </svg>
          <span class="logo-text">CasCade</span>
        </router-link>
      </div>
      
      <!-- 导航菜单 -->
      <ul class="nav-menu">
        <li class="nav-item">
          <router-link class="nav-link" to="/home">主页</router-link>
        </li>
        <li class="nav-item">
          <router-link class="nav-link" to="/competition">比赛</router-link>
        </li>
        <li class="nav-item">
          <router-link class="nav-link" to="/about">关于我们</router-link>
        </li>
      </ul>
      
      <!-- 用户功能区 -->
      <div class="nav-user">
        <router-link class="login-btn" to="/login">登录/注册</router-link>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useRouter } from 'vue-router'; 

const router = useRouter(); 

const handleLogout = () => {
  // 1. 清除 LocalStorage
  localStorage.removeItem('cascade_token');
  
  // 2.以此类推，清除用户信息状态
  
  // 3. 跳转回登录页或首页
  router.push('/login');
  
  // 4. (可选) 调用后端 logout 接口让 token 失效
  // request.post('/auth/logout'); 
};
</script>

<style scoped src="./NavBar.css"></style>