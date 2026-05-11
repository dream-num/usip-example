export type FileItem = {
  id: number
  name: string
  unitId: string
  unitType: number
  role: 'owner' | 'editor' | 'reader'
  updatedAt: string
  openUrl: string
  exportUrl: string
}

export type FilesResp = {
  userId: string
  files: FileItem[]
}

export type UserListResp = {
  users: Array<{
    user_id: string
    nickname: string
  }>
  next: number
}
