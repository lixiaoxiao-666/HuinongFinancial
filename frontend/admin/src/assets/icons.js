// 创建登录页面所需的SVG图标
(function() {
  const iconSvg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
  iconSvg.setAttribute('aria-hidden', 'true');
  iconSvg.style.position = 'absolute';
  iconSvg.style.width = '0';
  iconSvg.style.height = '0';
  iconSvg.style.overflow = 'hidden';
  iconSvg.setAttribute('version', '1.1');
  iconSvg.setAttribute('xmlns', 'http://www.w3.org/2000/svg');
  iconSvg.setAttribute('xmlns:xlink', 'http://www.w3.org/1999/xlink');
  
  // 用户图标
  const userIcon = document.createElementNS('http://www.w3.org/2000/svg', 'symbol');
  userIcon.setAttribute('id', 'icon-user');
  userIcon.setAttribute('viewBox', '0 0 24 24');
  userIcon.innerHTML = `
    <path fill="currentColor" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"></path>
  `;
  
  // 锁图标
  const lockIcon = document.createElementNS('http://www.w3.org/2000/svg', 'symbol');
  lockIcon.setAttribute('id', 'icon-lock');
  lockIcon.setAttribute('viewBox', '0 0 24 24');
  lockIcon.innerHTML = `
    <path fill="currentColor" d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z"></path>
  `;
  
  // 眼睛图标
  const eyeIcon = document.createElementNS('http://www.w3.org/2000/svg', 'symbol');
  eyeIcon.setAttribute('id', 'icon-eye');
  eyeIcon.setAttribute('viewBox', '0 0 24 24');
  eyeIcon.innerHTML = `
    <path fill="currentColor" d="M12 4.5C7 4.5 2.73 7.61 1 12c1.73 4.39 6 7.5 11 7.5s9.27-3.11 11-7.5c-1.73-4.39-6-7.5-11-7.5zM12 17c-2.76 0-5-2.24-5-5s2.24-5 5-5 5 2.24 5 5-2.24 5-5 5zm0-8c-1.66 0-3 1.34-3 3s1.34 3 3 3 3-1.34 3-3-1.34-3-3-3z"></path>
  `;
  
  // 盾牌图标（用于验证码）
  const shieldIcon = document.createElementNS('http://www.w3.org/2000/svg', 'symbol');
  shieldIcon.setAttribute('id', 'icon-shield');
  shieldIcon.setAttribute('viewBox', '0 0 24 24');
  shieldIcon.innerHTML = `
    <path fill="currentColor" d="M12 1L3 5v6c0 5.55 3.84 10.74 9 12 5.16-1.26 9-6.45 9-12V5l-9-4zm0 10.99h7c-.53 4.12-3.28 7.79-7 8.94V12H5V6.3l7-3.11v8.8z"></path>
  `;
  
  // 将所有图标添加到SVG中
  iconSvg.appendChild(userIcon);
  iconSvg.appendChild(lockIcon);
  iconSvg.appendChild(eyeIcon);
  iconSvg.appendChild(shieldIcon);
  
  // 将SVG添加到文档中
  document.body.appendChild(iconSvg);
})(); 