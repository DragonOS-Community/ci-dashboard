<template>
  <footer class="footer">
    <div class="footer-content">
      <div class="footer-left">
        <span class="copyright">
          © Copyright DragonOS-Community 2022-2026
        </span>
      </div>
      <div class="footer-right">
        <a
          href="https://dragonos.org"
          target="_blank"
          rel="noopener noreferrer"
          class="footer-link"
        >
          <t-icon name="link" />
          官网
        </a>
        <a
          href="https://github.com/DragonOS-Community/DragonOS"
          target="_blank"
          rel="noopener noreferrer"
          class="footer-link"
        >
          <t-icon name="logo-github" />
          社区主页
        </a>
        <a
          href="https://github.com/DragonOS-Community/ci-dashboard/issues"
          target="_blank"
          rel="noopener noreferrer"
          class="footer-link"
        >
          <t-icon name="error" />
          提交Bug
        </a>
        <t-button
          theme="default"
          variant="text"
          @click="goToLogin"
          class="footer-login-btn"
        >
          <t-icon name="user" />
          管理员登录
        </t-button>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { useRouter } from "vue-router";

const router = useRouter();

const goToLogin = () => {
  router.push("/admin/login");
};
</script>

<style scoped>
.footer {
  background-color: #ffffff;
  border-top: 1px solid #e5e7eb;
  padding: 24px 32px;
  margin-top: auto;
}

.footer-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;
}

.footer-left {
  display: flex;
  align-items: center;
}

.copyright {
  font-size: 14px;
  color: #6b7280;
  font-weight: 400;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 0;
  flex-wrap: wrap;
}

/* 添加竖线分隔符（仅桌面端） */
.footer-right > *:not(:first-child) {
  position: relative;
  padding-left: 24px;
}

.footer-right > *:not(:first-child)::before {
  content: "";
  position: absolute;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  width: 1px;
  height: 16px;
  background-color: #e5e7eb;
}

.footer-link {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #6b7280;
  text-decoration: none;
  transition: all 0.2s ease;
  padding: 4px 8px;
  border-radius: 6px;
  font-weight: 500;
}

.footer-link :deep(.t-icon) {
  font-size: 16px;
  color: #6b7280;
  transition: color 0.2s ease;
}

.footer-link:hover {
  color: #f59e0b;
  background-color: #fef9f3;
}

.footer-link:hover :deep(.t-icon) {
  color: #f59e0b;
}

.footer-login-btn {
  border-radius: 6px;
  transition: all 0.2s ease;
  height: auto;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px;
  color: #6b7280;
  font-weight: 400;
  padding: 4px 8px;
  font-size: 14px;
}

.footer-login-btn :deep(.t-button__text) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px;
  line-height: 1;
}

.footer-login-btn :deep(.t-icon) {
  font-size: 16px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #6b7280;
}

.footer-login-btn:hover {
  background-color: #f3f4f6;
  color: #1f2937;
}

.footer-login-btn:hover :deep(.t-icon) {
  color: #1f2937;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .footer {
    padding: 20px 16px;
  }

  .footer-content {
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 16px;
  }

  .footer-right {
    justify-content: center;
    gap: 16px;
  }

  /* 移动端移除分隔符和左侧padding */
  .footer-right > *:not(:first-child) {
    padding-left: 0;
  }

  .footer-right > *:not(:first-child)::before {
    display: none;
  }

  .copyright {
    font-size: 13px;
  }

  .footer-link {
    font-size: 13px;
    padding: 6px 12px;
  }

  .footer-login-btn {
    width: 100%;
    max-width: 200px;
  }
}
</style>

