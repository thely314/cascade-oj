<template>
  <form class="auth-form" @submit.prevent="handleRegister">
    <!-- 用户名 -->
    <div class="form-group">
      <label>用户名</label>
      <input 
        type="text" 
        v-model="form.username" 
        placeholder="请输入用户名" 
        class="input-field"
        required
      />
    </div>

    <!-- 邮箱 -->
    <div class="form-group">
      <label>邮箱</label>
      <input 
        type="email" 
        v-model="form.email" 
        placeholder="请输入邮箱地址" 
        class="input-field"
        required
      />
    </div>

    <!-- 密码 -->
    <div class="form-group">
      <label>密码</label>
      <div class="password-field">
        <input 
          :type="showPassword ? 'text' : 'password'" 
          v-model="form.password" 
          placeholder="请输入密码（至少8位）" 
          class="input-field"
          required
          minlength="8"
        />
        <button type="button" class="toggle-password" @click="showPassword = !showPassword">
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" v-if="showPassword">
            <path d="M12 15c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0-2c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm6-6.88c.76 0 1.41.36 1.87.91l.8 1.04c.27.35.41.8.41.1.47-.14.92-.41 1.27l-.8 1.04c-.46.55-1.11.91-1.87.91-.76 0-1.41-.36-1.87-.91l-.8-1.04c-.27-.35-.41-.8-.41-1.27v-.08c.01-.47.14-.92.41-1.27l.8-1.04c.46-.55 1.11-.91 1.87-.91zm-18 0c.76 0 1.41.36 1.87.91l.8 1.04c.27.35.41.8.41 1.27v.08c-.01.47-.14.92-.41 1.27l-.8 1.04c-.46.55-1.11.91-1.87.91-.76 0-1.41-.36-1.87-.91l-.8-1.04c-.27-.35-.41-.8-.41-1.27v-.08c.01-.47.14-.92.41-1.27l.8-1.04c.46-.55 1.11-.91 1.87-.91z" fill="currentColor"/>
          </svg>
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" v-else>
            <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zm0 14c-3.87 0-7-3.13-7-7s3.13-7 7-7 7 3.13 7 7-3.13 7-7 7zm0-11c-2.76 0-5 2.24-5 5s2.24 5 5 5 5-2.24 5-5-2.24-5-5-5zm0 8c-1.66 0-3-1.34-3-3s1.34-3 3-3 3 1.34 3 3-1.34 3-3 3z" fill="currentColor"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- 确认密码 (这里增加了眼睛图标) -->
    <div class="form-group">
      <label>确认密码</label>
      <div class="password-field">
        <input 
          :type="showConfirmPassword ? 'text' : 'password'" 
          v-model="form.confirmPassword" 
          placeholder="请再次输入密码" 
          class="input-field"
          required
          minlength="8"
        />
        <button type="button" class="toggle-password" @click="showConfirmPassword = !showConfirmPassword">
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" v-if="showConfirmPassword">
            <path d="M12 15c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0-2c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm6-6.88c.76 0 1.41.36 1.87.91l.8 1.04c.27.35.41.8.41.1.47-.14.92-.41 1.27l-.8 1.04c-.46.55-1.11.91-1.87.91-.76 0-1.41-.36-1.87-.91l-.8-1.04c-.27-.35-.41-.8-.41-1.27v-.08c.01-.47.14-.92.41-1.27l.8-1.04c.46-.55 1.11-.91 1.87-.91zm-18 0c.76 0 1.41.36 1.87.91l.8 1.04c.27.35.41.8.41 1.27v.08c-.01.47-.14.92-.41 1.27l-.8 1.04c-.46.55-1.11.91-1.87.91-.76 0-1.41-.36-1.87-.91l-.8-1.04c-.27-.35-.41-.8-.41-1.27v-.08c.01-.47.14-.92.41-1.27l.8-1.04c.46-.55 1.11-.91 1.87-.91z" fill="currentColor"/>
          </svg>
          <svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" v-else>
            <path d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zm0 14c-3.87 0-7-3.13-7-7s3.13-7 7-7 7 3.13 7 7-3.13 7-7 7zm0-11c-2.76 0-5 2.24-5 5s2.24 5 5 5 5-2.24 5-5-2.24-5-5-5zm0 8c-1.66 0-3-1.34-3-3s1.34-3 3-3 3 1.34 3 3-1.34 3-3 3z" fill="currentColor"/>
          </svg>
        </button>
      </div>
    </div>
    
    <!-- 协议 -->
    <div class="form-check">
      <label class="checkbox-label">
        <input type="checkbox" v-model="form.agree" required />
        <span class="checkmark"></span>
        <span>我同意 <a href="#">用户协议</a> 和 <a href="#">隐私政策</a></span>
      </label>
    </div>
    
    <button 
      type="submit" 
      class="btn-primary" 
      :disabled="!form.agree"
    >
      注册
    </button>
  </form>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import md5 from 'js-md5';

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  agree: false,
});

const showPassword = ref(false);
const showConfirmPassword = ref(false); // 新增状态

const handleRegister = () => {
  if (form.password !== form.confirmPassword) {
    alert('两次密码输入不一致');
    return;
  }
  console.log('注册提交：', form);
  const encryptedPwd = md5.md5(form.password);
  console.log('加密后的密码：', encryptedPwd);
};
</script>

<style scoped>
.auth-form {
  width: 100%;
  max-width: 360px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 16px;
  box-sizing: border-box;
  animation: fadeIn 0.4s ease;
}

.form-group label {
  display: block;
  font-size: 14px;
  margin-bottom: 8px;
  color: #bababa;
}

.input-field {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #2a3331;
  border-radius: 6px;
  font-size: 14px;
  background: #0f1413;
  color: #cbd5c0;
  outline: none;
  transition: all 0.3s ease;
}

.input-field:focus {
  border-color: #23aa8f;
  background: #131a18;
  box-shadow: 0 0 0 2px rgba(35, 170, 143, 0.1);
}

.input-field::placeholder {
  color: #5f6f6b;
}

/* 密码输入框容器 */
.password-field {
  position: relative;
}

.password-field .input-field {
  padding-right: 40px; /* 留出眼睛图标的位置 */
}

.toggle-password {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  color: #5f6f6b;
  display: flex;
  align-items: center;
}

.toggle-password:hover {
  color: #23aa8f;
}

.toggle-password svg {
  width: 20px;
  height: 20px;
}

/* 协议勾选 */
.form-check {
  font-size: 13px;
  margin-bottom: 8px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  color: #8c9e9a;
  user-select: none;
}

.checkbox-label input {
  display: none;
}

.checkmark {
  width: 14px;
  height: 14px;
  border: 1px solid #4a5755;
  border-radius: 3px;
  margin-right: 8px;
  position: relative;
  transition: all 0.2s;
}

.checkbox-label input:checked ~ .checkmark {
  background: #23aa8f;
  border-color: #23aa8f;
}

.form-check a {
  color: #1dad80; /* 这里也可以统一用 Logo 绿 */
  text-decoration: none;
  margin: 0 4px;
}

/* 按钮 */
.btn-primary {
  background-color: #23aa8f;
  color: #01210f;
  border: none;
  padding: 12px;
  border-radius: 6px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background-color: #1eb897;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background-color: #184a41;
  color: #6c7c7a;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>