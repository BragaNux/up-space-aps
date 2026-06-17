const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000'
const TOKEN_KEY = 'up_espaco_token'

// pega o token salvo no localStorage
export function getToken() {
  return localStorage.getItem(TOKEN_KEY)
}

// salva o token (ou remove, se vier null/undefined)
export function setToken(token) {
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
}

// faz a chamada fetch pra API, ja colocando o token de auth e tratando erro/204/json
async function request(path, { method = 'GET', body, auth = true } = {}) {
  const headers = { 'Content-Type': 'application/json' }
  const token = auth ? getToken() : null
  if (token) headers.Authorization = `Bearer ${token}`

  const res = await fetch(`${API_URL}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })

  if (res.status === 204) return null

  const data = await res.json().catch(() => null)

  if (!res.ok) {
    throw new Error(data?.error || `Erro ${res.status} ao acessar ${path}`)
  }

  return data
}

// Auth
export const register = (name, email, password, role) =>
  request('/api/auth/register', { method: 'POST', body: { name, email, password, role }, auth: false })

export const login = (email, password) =>
  request('/api/auth/login', { method: 'POST', body: { email, password }, auth: false })

export const forgotPassword = (email) =>
  request('/api/auth/forgot-password', { method: 'POST', body: { email }, auth: false })

export const getMe = () => request('/api/me')

export const updateMe = (profile) => request('/api/me', { method: 'PUT', body: profile })

export const getMyChildren = () => request('/api/me/children')

export const listUsersByRole = (role) => request(`/api/users?role=${encodeURIComponent(role)}`)

// Students
export const listStudents = () => request('/api/students')

export const getStudentById = (id) => request(`/api/students/${id}`)

export const createStudent = (payload) => request('/api/students', { method: 'POST', body: payload })

export const updateStudent = (studentId, payload) =>
  request(`/api/students/${studentId}`, { method: 'PUT', body: payload })

export const updateStudentPresence = (studentId, status) =>
  request(`/api/students/${studentId}/presence`, { method: 'PATCH', body: { status } })

// Guardians (authorized pickup contacts for a child)
export const listGuardians = (studentId) => request(`/api/students/${studentId}/guardians`, { auth: false })

export const createGuardian = (studentId, data) =>
  request(`/api/students/${studentId}/guardians`, { method: 'POST', body: data })

export const updateGuardian = (guardianId, data) =>
  request(`/api/guardians/${guardianId}`, { method: 'PUT', body: data })

export const deleteGuardian = (guardianId) =>
  request(`/api/guardians/${guardianId}`, { method: 'DELETE' })

// Converts an image File into a base64 data URI for inline storage.
export function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result)
    reader.onerror = () => reject(new Error('Não foi possível ler a imagem.'))
    reader.readAsDataURL(file)
  })
}

// Timeline (scoped to a child)
export const getTimeline = (studentId) => request(`/api/timeline?student_id=${studentId}`)

export const createTimelineEvent = (studentId, data) =>
  request(`/api/timeline?student_id=${studentId}`, { method: 'POST', body: data })

// Posts (scoped to a child)
export const getPosts = (studentId) => request(`/api/posts?student_id=${studentId}`)

export const createPost = (studentId, data) =>
  request(`/api/posts?student_id=${studentId}`, { method: 'POST', body: data })

export const likePost = (id) => request(`/api/posts/${id}/like`, { method: 'POST', auth: false })

export const unlikePost = (id) => request(`/api/posts/${id}/unlike`, { method: 'POST', auth: false })

export const bookmarkPost = (id) => request(`/api/posts/${id}/bookmark`, { method: 'POST', auth: false })

export const unbookmarkPost = (id) => request(`/api/posts/${id}/unbookmark`, { method: 'POST', auth: false })

export const getComments = (postId) => request(`/api/posts/${postId}/comments`, { auth: false })

export const addComment = (postId, text) =>
  request(`/api/posts/${postId}/comments`, { method: 'POST', body: { text } })

// Events (school-wide, not per child)
export const getEvents = () => request('/api/events')

export const createEvent = (data) => request('/api/events', { method: 'POST', body: data })

export const updateEvent = (id, data) => request(`/api/events/${id}`, { method: 'PUT', body: data })

export const deleteEvent = (id) => request(`/api/events/${id}`, { method: 'DELETE' })

export const rsvpEvent = (id) => request(`/api/events/${id}/rsvp`, { method: 'POST' })

// Announcements (school-wide, not per child)
export const getAnnouncements = () => request('/api/announcements')

export const createAnnouncement = (data) => request('/api/announcements', { method: 'POST', body: data })

export const updateAnnouncement = (id, data) => request(`/api/announcements/${id}`, { method: 'PUT', body: data })

export const deleteAnnouncement = (id) => request(`/api/announcements/${id}`, { method: 'DELETE' })

export const markAnnouncementRead = (id) => request(`/api/announcements/${id}/read`, { method: 'POST' })

// Milestones (scoped to a child)
export const getMilestones = (studentId) => request(`/api/milestones?student_id=${studentId}`)

export const createMilestone = (studentId, data) =>
  request('/api/milestones', { method: 'POST', body: { ...data, student_id: studentId } })

// Turmas (classes)
export const listTurmas = () => request('/api/turmas', { auth: false })

export const createTurma = (data) => request('/api/turmas', { method: 'POST', body: data })

export const updateTurma = (id, data) => request(`/api/turmas/${id}`, { method: 'PUT', body: data })

export const deleteTurma = (id) => request(`/api/turmas/${id}`, { method: 'DELETE' })

// Attendance
export const saveAttendance = (studentId, data) =>
  request(`/api/students/${studentId}/attendance`, { method: 'POST', body: data })

export const getTurmaAttendance = (turmaId, date) =>
  request(`/api/turmas/${turmaId}/attendance?date=${encodeURIComponent(date)}`)

export const getStudentAttendanceHistory = (studentId) =>
  request(`/api/students/${studentId}/attendance`)
