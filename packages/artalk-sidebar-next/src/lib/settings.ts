export function transSettingKey(key: string) {
  const transMap: { [key: string]: string } = {
    host: '后端地址',
    port: '后端端口',
    app_key: '安全密钥',
    debug: '调试模式',
    timezone: '时间区域',
    db: '数据库',
    log: '系统日志',
    cache: '系统缓存',
    trusted_domains: '安全域名',
    ssl: 'SSL',
    site_default: '默认站点',
    admin_users: '管理员账户',
    login_timeout: '登陆时长',
    moderator: '评论审核',
    captcha: '验证码',
    email: '邮箱通知',
    img_upload: '图片上传',
    admin_notify: '多元推送',
  }

  return transMap[key] || key
}

export default { transSettingKey }
