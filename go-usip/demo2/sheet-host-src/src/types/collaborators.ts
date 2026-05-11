export type CollaboratorSubject = {
  id: string
  name: string
  avatar: string
}

export type CollaboratorItem = {
  subject: CollaboratorSubject
  role: 'owner' | 'editor' | 'reader' | string
}

export type CollaboratorsResp = {
  collaborators: Array<{
    unitID: string
    subjects: CollaboratorItem[]
  }>
}
