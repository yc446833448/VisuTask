import type { User, WalletRecord } from '@/types'

export const mockUser: User = {
  id: 'user_001',
  vipLevel: 1,
  maxConcurrent: 5,
  balance: 128.50,
  avatar: undefined,
}

export const mockWalletRecords: WalletRecord[] = [
  { id: 'w1', description: 'VIP 1 月卡', amount: -29.90, date: '6/15' },
  { id: 'w2', description: '充值', amount: 50.00, date: '6/10' },
  { id: 'w3', description: '任务执行消耗 ×48', amount: -4.80, date: '6/09' },
  { id: 'w4', description: '充值', amount: 100.00, date: '6/01' },
  { id: 'w5', description: '任务执行消耗 ×32', amount: -3.20, date: '5/28' },
]
