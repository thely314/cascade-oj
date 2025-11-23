<template>
  <form class="auth-form" @submit.prevent="handleLogin">
    <div class="form-group">
      <label>用户名或邮箱</label>
      <input 
        type="text" 
        v-model="form.usernameOrEmail" 
        placeholder="请输入用户名或邮箱" 
        class="input-field"
        required
      />
    </div>
    
    <div class="form-group">
      <label>密码</label>
      <div class="password-field">
        <input 
          :type="showPassword ? 'text' : 'password'" 
          v-model="form.password" 
          placeholder="请输入密码" 
          class="input-field"
          required
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
    
    <div class="form-actions">
      <label class="checkbox-label">
        <input type="checkbox" v-model="form.remember" />
        记住我
      </label>
      <a href="#" class="forgot-link">忘记密码?</a>
    </div>
    
    <button type="submit" class="btn-primary">登录</button>
    
    <div class="or-separator">或者</div>
    
    <button type="button" class="btn-social google">
      <span>G</span> 使用 Google 登录
    </button>
    <button type="button" class="btn-social github">
      <span>GH</span> 使用 GitHub 登录
    </button>
  </form>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import md5 from 'js-md5';

const form = reactive({
  usernameOrEmail: '',
  password: '',
  remember: false,
});

const showPassword = ref(false);

const handleLogin = () => {
  console.log('登录提交：', form);
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
  gap: 12px;
  box-sizing: border-box;
}

.form-group label {
  display: block;
  font-size: 14px;
  margin-bottom: 8px;
  color: #bababa;
}

.input-field {
  width: 100%;
  padding: 12px;
  border: 1px solid rgba(255, 255, 255, 0.773);
  border-radius: 6px;
  font-size: 14px;
  background: #0f1716;
  color: #cbd5c0;
}

.password-field {
  position: relative;
}

.toggle-password {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  cursor: pointer;
}

.toggle-password svg {
  width: 20px;
  height: 20px;
  color: #9fbebb;
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 14px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.checkbox-label input {
  margin-right: 8px;
}

.forgot-link {
  color: #1dad7f;
  text-decoration: none;
  font-size: 14px;
}

.btn-primary {
  background-color: #23aa8f;
  color: #01210f;
  border: none;
  padding: 12px;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-primary:hover {
  background-color: #159067;
}

.or-separator {
  text-align: center;
  color: #999999;
  font-size: 14px;
  margin: 12px 0;
  position: relative;
}

.or-separator::before,
.or-separator::after {
  content: '';
  position: absolute;
  top: 50%;
  width: 40%;
  height: 1px;
  background-color: #eeeeee;
}

.or-separator::before {
  left: 0;
}

.or-separator::after {
  right: 0;
}

.btn-social {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border: 1px solid rgba(255,255,255,0.04);
  padding: 12px;
  border-radius: 4px;
  background-color: #23aa8f;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.btn-social span {
  font-weight: bold;
  font-size: 16px;
}

.btn-social.google {
  color: #0044ff;
}

.btn-social.github {
  color: #000000;
}
</style>