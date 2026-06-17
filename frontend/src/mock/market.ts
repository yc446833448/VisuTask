import type { MarketScript } from '@/types'

export const mockFreeScripts: MarketScript[] = [
  { id: 'm1', name: 'Excel 批量录入', description: '将 Excel 数据逐条录入到 Web 系统，支持自定义字段映射', version: 'v1.2', versions: ['v1.2', 'v1.1', 'v1.0'], downloads: 128, author: 'user_a', isVip: false },
  { id: 'm2', name: '通用自动填表', description: '通用表单自动填写，适配多种表单格式，智能识别输入框', version: 'v2.0', versions: ['v2.0', 'v1.5', 'v1.0'], downloads: 256, author: 'user_b', isVip: false },
  { id: 'm3', name: '邮件自动归档', description: '自动读取邮件附件并归档到指定文件夹，支持过滤规则', version: 'v1.0', versions: ['v1.0'], downloads: 89, author: 'user_c', isVip: false },
  { id: 'm4', name: '批量文件重命名', description: '按规则批量重命名文件，支持正则表达式和序号', version: 'v1.3', versions: ['v1.3', 'v1.2', 'v1.1'], downloads: 175, author: 'user_d', isVip: false },
]

export const mockVipScripts: MarketScript[] = [
  { id: 'v1', name: 'ERP 全自动同步', description: '企业级 ERP 系统全自动数据同步，支持多实例并行', version: 'v3.0', versions: ['v3.0', 'v2.1'], downloads: 64, author: 'pro_user', isVip: true, vipLevel: 2 },
  { id: 'v2', name: '跨平台数据迁移', description: '跨多个平台的数据迁移和格式转换，支持 SAP/Oracle/金蝶', version: 'v1.0', versions: ['v1.0'], downloads: 32, author: 'pro_user', isVip: true, vipLevel: 3 },
  { id: 'v3', name: '智能回归测试套件', description: '桌面应用智能回归测试，自动生成测试用例并执行', version: 'v2.5', versions: ['v2.5', 'v2.0', 'v1.0'], downloads: 48, author: 'pro_user', isVip: true, vipLevel: 1 },
]
