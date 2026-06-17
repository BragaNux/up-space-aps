import { useState, useEffect, useCallback } from 'react'
import Login from './Login'
import EsqueceuSenha from './EsqueceuSenha'
import CriarConta from './CriarConta'
import Feed from './Feed'
import Jornada from './Jornada'
import Agenda from './Agenda'
import Comunicados from './Comunicados'
import Perfil from './Perfil'
import Crianca from './Crianca'
import Turmas from './Turmas'
import SearchableSelect from './SearchableSelect'
import DashboardAdmin from './DashboardAdmin'
import defaultAvatar from '../assets/default.png'
import logoUpEspaco from '../assets/logo_up_espaco.png'
import {
  getToken,
  setToken,
  getMe,
  getMyChildren,
  listStudents,
  getStudentById,
  createStudent,
  getTimeline,
  getEvents,
  createTimelineEvent,
  listTurmas,
  listUsersByRole,
} from './api'
import {
  LogIn,
  Images,
  CalendarDays,
  Mail,
  Megaphone,
  CheckCircle2,
  GraduationCap,
  Utensils,
  Heart,
  Home,
  Compass,
  Calendar,
  User,
  Search,
  LogOut,
  ChevronRight,
  ChevronLeft,
  TrendingUp,
  Plus,
  Layers,
} from 'lucide-react'

const quickAccess = [
  { label: 'Fotos', desc: 'Veja as novidades', icon: Images, bg: 'bg-pink-100', fg: 'text-pink-500', target: 'Atividades' },
  { label: 'Eventos', desc: 'Próximos compromissos', icon: CalendarDays, bg: 'bg-indigo-100', fg: 'text-indigo-500', target: 'Agenda' },
  { label: 'Comunicados', desc: 'Avisos da escola', icon: Mail, bg: 'bg-emerald-100', fg: 'text-emerald-600', target: 'Comunicados' },
  { label: 'Jornada', desc: 'Marcos de desenvolvimento', icon: Megaphone, bg: 'bg-amber-100', fg: 'text-amber-500', target: 'Jornada' },
]

const navItems = [
  { label: 'Início', icon: Home },
  { label: 'Atividades', icon: Compass },
  { label: 'Jornada', icon: TrendingUp },
  { label: 'Agenda', icon: Calendar },
  { label: 'Perfil', icon: User },
]

// escolhe um icone pra timeline com base em palavras-chave no titulo do evento
const timelineIconFor = (title = '') => {
  const t = title.toLowerCase()
  if (t.includes('lanche') || t.includes('almoço') || t.includes('refeição')) return { icon: Utensils, fg: 'text-amber-600' }
  if (t.includes('terapia') || t.includes('saúde')) return { icon: Heart, fg: 'text-emerald-700' }
  if (t.includes('pedagóg') || t.includes('aula') || t.includes('artes')) return { icon: GraduationCap, fg: 'text-pink-600' }
  return { icon: CheckCircle2, fg: 'text-emerald-500' }
}

// foto de perfil com fallback pro avatar padrao se a imagem der erro ou nao existir
function Avatar({ src, alt, className }) {
  return (
    <img
      src={src || defaultAvatar}
      alt={alt}
      className={className}
      onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
    />
  )
}

// menu lateral fixo (desktop) com navegacao e botao de sair
function Sidebar({ active, onChange, onLogout, showNav }) {
  return (
    <aside className="hidden w-64 shrink-0 flex-col border-r border-slate-100 bg-white px-4 py-6 md:flex">
      <div className="flex items-center gap-2.5 px-2">
        <img src={logoUpEspaco} alt="Logo Up - Espaço" className="h-10 w-10 object-contain" />
        <span className="text-xl font-extrabold tracking-tight text-emerald-600">
          Up - Espaço
        </span>
      </div>

      <nav className="mt-8 flex flex-1 flex-col gap-1">
        {showNav && navItems.map(({ label, icon: Icon }) => {
          const isActive = active === label
          return (
            <button
              key={label}
              onClick={() => onChange(label)}
              className={`flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium transition-all duration-200 ${
                isActive
                  ? 'bg-emerald-50 text-emerald-700 translate-x-0.5'
                  : 'text-slate-500 hover:translate-x-0.5 hover:bg-slate-50 hover:text-slate-800'
              }`}
            >
              <Icon className="h-5 w-5 transition-transform duration-200" strokeWidth={isActive ? 2.4 : 2} />
              {label}
            </button>
          )
        })}
      </nav>

      <button
        onClick={onLogout}
        className="flex items-center gap-3 rounded-xl px-3 py-2.5 text-sm font-medium text-slate-500 transition-all duration-200 hover:bg-rose-50 hover:text-rose-600"
      >
        <LogOut className="h-5 w-5" />
        Sair
      </button>
    </aside>
  )
}

// barra superior com saudacao, troca de crianca e acesso ao perfil
function Topbar({ user, onSwitchChild, onOpenProfile }) {
  const firstName = user?.name?.split(' ')[0] ?? ''
  return (
    <header className="rounded-3xl bg-white px-6 py-4 shadow-sm ring-1 ring-slate-100 flex items-center justify-between gap-4">
      <div>
        <p className="text-xs font-semibold uppercase tracking-wide text-slate-400">Bem-vindo(a) ao Up - Espaço</p>
        <h1 className="text-xl font-extrabold text-emerald-600">Olá, {firstName}!</h1>
      </div>

      <div className="flex items-center gap-3">
        {onSwitchChild && (
          <button
            onClick={onSwitchChild}
            className="hidden items-center gap-1.5 rounded-full border border-slate-200 px-4 py-2 text-xs font-bold text-slate-600 transition hover:bg-slate-50 sm:flex"
          >
            <ChevronLeft className="h-4 w-4" />
            Trocar criança
          </button>
        )}
        <button
          onClick={onOpenProfile}
          className="transition-transform duration-200 hover:scale-105 active:scale-95 flex shrink-0"
          aria-label="Ver Perfil"
        >
          <Avatar
            src={user?.avatar_url}
            alt={`Foto de ${firstName}`}
            className="h-10 w-10 rounded-full object-cover ring-2 ring-emerald-500/20 shadow-sm cursor-pointer"
          />
        </button>
      </div>
    </header>
  )
}

// card grande da tela inicial com foto, status de presenca e atalho pro diario da crianca
function ChildCard({ student, onOpenChild }) {
  if (!student) return null
  const present = student.presence_status === 'present'
  const checkInTime = student.check_in_at
    ? new Date(student.check_in_at).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
    : '--:--'

  return (
    <section className="animate-slide-up rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-100 transition-shadow duration-300 hover:shadow-md">
      <div className="flex flex-col gap-5 sm:flex-row sm:items-center">
        <button onClick={onOpenChild} className="shrink-0" aria-label={`Abrir perfil de ${student.name}`}>
          <Avatar
            src={student.photo_url}
            alt={`Foto de ${student.name}`}
            className="h-28 w-28 rounded-2xl object-cover transition-transform duration-300 hover:scale-105 hover:opacity-90"
          />
        </button>
        <div className="flex-1">
          <button onClick={onOpenChild} className="text-left">
            <h2 className="text-2xl font-extrabold text-slate-800 transition-colors hover:text-emerald-700">{student.name}</h2>
          </button>
          <p className="text-slate-500">{student.group_name}</p>
          <div className="mt-3 flex items-center gap-2">
            <span className={`h-2.5 w-2.5 rounded-full ${present ? 'animate-pulse bg-emerald-500' : 'bg-slate-300'}`} />
            <span className={`text-sm font-semibold tracking-wide ${present ? 'text-emerald-600' : 'text-slate-400'}`}>
              {present ? 'PRESENTE' : 'AUSENTE'}
            </span>
            {present && (
              <span className="ml-3 flex items-center gap-1.5 text-sm text-slate-500">
                <LogIn className="h-4 w-4 text-slate-400" />
                Entrada: <span className="font-bold text-slate-700">{checkInTime}</span>
              </span>
            )}
          </div>
        </div>
        <button onClick={onOpenChild} className="flex items-center justify-center gap-2 rounded-full bg-emerald-500 px-6 py-3 text-sm font-semibold text-white shadow-md shadow-emerald-200 transition-all duration-200 hover:bg-emerald-600 hover:shadow-lg hover:shadow-emerald-200 active:scale-95">
          Ver Diário
          <ChevronRight className="h-4 w-4 transition-transform group-hover:translate-x-0.5" />
        </button>
      </div>
    </section>
  )
}

// linha de cards com status do dia, horario de entrada e proximo evento
function StatsRow({ student, nextEvent }) {
  const present = student?.presence_status === 'present'
  const stats = [
    { label: 'Status hoje', value: present ? 'Presente' : 'Ausente', hint: student?.group_name ?? '' },
    {
      label: 'Entrada',
      value: student?.check_in_at
        ? new Date(student.check_in_at).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
        : '--:--',
      hint: 'registrado hoje',
    },
    {
      label: 'Próximo evento',
      value: nextEvent?.title ?? 'Sem eventos',
      hint: nextEvent ? new Date(nextEvent.starts_at).toLocaleDateString('pt-BR', { day: '2-digit', month: 'long' }) : '',
    },
  ]

  return (
    <section className="grid grid-cols-1 gap-4 sm:grid-cols-3">
      {stats.map(({ label, value, hint }, i) => (
        <div
          key={label}
          className={`animate-pop-in stagger-${i + 1} rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md`}
        >
          <p className="text-xs font-medium uppercase tracking-wide text-slate-400">{label}</p>
          <p className="mt-1 text-xl font-bold text-slate-800">{value}</p>
          <p className="text-xs text-slate-400">{hint}</p>
        </div>
      ))}
    </section>
  )
}

// grade de atalhos rapidos pras principais telas (fotos, eventos, comunicados, jornada)
function QuickAccess({ onOpen }) {
  return (
    <section>
      <div className="mb-4 flex items-center justify-between">
        <h3 className="text-lg font-bold text-slate-800">Acesso Rápido</h3>
      </div>

      <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        {quickAccess.map(({ label, desc, icon: Icon, bg, fg, target }, i) => (
          <button
            key={label}
            onClick={() => onOpen(target)}
            className={`animate-pop-in stagger-${i + 1} group flex flex-col items-start gap-3 rounded-2xl bg-white p-4 text-left shadow-sm ring-1 ring-slate-100 transition-all duration-200 hover:-translate-y-1 hover:shadow-lg`}
          >
            <span className={`grid h-12 w-12 place-items-center rounded-xl ${bg} ${fg} transition-transform duration-200 group-hover:scale-110`}>
              <Icon className="h-6 w-6" strokeWidth={1.8} />
            </span>
            <div>
              <p className="font-semibold text-slate-800">{label}</p>
              <p className="text-xs text-slate-400">{desc}</p>
            </div>
          </button>
        ))}
      </div>
    </section>
  )
}

// formulario rapido pra adicionar um registro novo na timeline do dia
function NewTimelineEntryForm({ onCreate, onClose }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  async function submit(e) {
    e.preventDefault()
    if (!title.trim() || !description.trim()) return
    setSaving(true)
    setError('')
    try {
      await onCreate({ title, description, occurred_at: new Date().toISOString() })
      setTitle('')
      setDescription('')
      onClose()
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <form onSubmit={submit} className="mb-4 space-y-2 rounded-2xl bg-slate-50 p-4">
      <input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Título (ex: Lanche da tarde)"
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />
      <textarea
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        placeholder="Descrição"
        rows={2}
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />
      {error && <p className="text-xs text-rose-500">{error}</p>}
      <div className="flex justify-end gap-2">
        <button type="button" onClick={onClose} className="rounded-full px-3 py-1.5 text-sm text-slate-500 hover:bg-slate-100">
          Cancelar
        </button>
        <button
          type="submit"
          disabled={saving}
          className="rounded-full bg-emerald-500 px-4 py-1.5 text-sm font-semibold text-white hover:bg-emerald-600 disabled:opacity-60"
        >
          {saving ? 'Salvando...' : 'Adicionar'}
        </button>
      </div>
    </form>
  )
}

// lista os eventos da timeline de hoje, com opcao de adicionar um novo (se puder editar)
function DaySummary({ timeline, canEdit, onCreate }) {
  const [showForm, setShowForm] = useState(false)
  return (
    <section className="rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-100">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-bold text-slate-800">Resumo do Dia</h3>
        {canEdit && !showForm && (
          <button onClick={() => setShowForm(true)} className="flex items-center gap-1 text-sm font-medium text-emerald-600 hover:underline">
            <Plus className="h-4 w-4" />
            Adicionar
          </button>
        )}
      </div>

      {showForm && <div className="mt-4"><NewTimelineEntryForm onCreate={onCreate} onClose={() => setShowForm(false)} /></div>}

      {timeline.length === 0 ? (
        <p className="mt-4 text-sm text-slate-400">Nenhum registro por hoje ainda.</p>
      ) : (
        <ol className="mt-5 space-y-1">
          {timeline.map((event, i) => {
            const { icon: Icon, fg } = timelineIconFor(event.title)
            return (
              <li key={event.id} className={`animate-slide-in-left stagger-${Math.min(i + 1, 6)} flex gap-4`}>
                <div className="flex flex-col items-center">
                  <span className="grid h-10 w-10 place-items-center rounded-full bg-slate-50 ring-1 ring-slate-100 transition-transform duration-200 hover:scale-110">
                    <Icon className={`h-5 w-5 ${fg}`} strokeWidth={2} />
                  </span>
                  {i < timeline.length - 1 && (
                    <span className="my-1 w-px flex-1 border-l-2 border-dashed border-slate-200" />
                  )}
                </div>

                <div className="mb-2 flex-1 rounded-2xl p-3.5">
                  <div className="flex items-start justify-between gap-2">
                    <h4 className="font-bold text-slate-800">{event.title}</h4>
                    <span className="text-sm text-slate-400">
                      {new Date(event.occurred_at).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}
                    </span>
                  </div>
                  <p className="mt-1 text-sm leading-relaxed text-slate-500">{event.description}</p>
                </div>
              </li>
            )
          })}
        </ol>
      )}
    </section>
  )
}

// mensagem de estado vazio quando nao tem crianca selecionada/vinculada
function EmptyChildState({ message }) {
  return (
    <div className="animate-slide-up mx-auto max-w-md py-12 text-center">
      <span className="mx-auto mb-4 grid h-16 w-16 place-items-center rounded-full bg-slate-100 text-slate-400">
        <GraduationCap className="h-8 w-8" />
      </span>
      <p className="text-sm text-slate-500">{message}</p>
    </div>
  )
}

// pega o ultimo sobrenome do nome completo, usado pra sugerir parentesco entre crianca e responsavel
function lastName(fullName = '') {
  const parts = fullName.trim().split(/\s+/)
  return parts.length > 1 ? parts[parts.length - 1].toLowerCase() : ''
}

// formulario de cadastro de uma crianca nova: nome, nascimento, turma, professor e responsavel
function CreateChildForm({ currentUser, onCreated, onCancel }) {
  const [form, setForm] = useState({ name: '', birth_date: '', turma_id: null, teacher_user_id: currentUser.id, guardian_user_id: null })
  const [turmas, setTurmas] = useState([])
  const [teachers, setTeachers] = useState([])
  const [responsaveis, setResponsaveis] = useState([])
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    listTurmas().then(setTurmas).catch(() => {})
    listUsersByRole('profissional').then(setTeachers).catch(() => {})
    listUsersByRole('responsavel').then(setResponsaveis).catch(() => {})
  }, [])

  function update(field, value) {
    setForm((prev) => ({ ...prev, [field]: value }))
  }

  const turmaOptions = turmas.map((t) => ({ value: t.id, label: t.name }))

  const teacherOptions = [
    { value: currentUser.id, label: `Eu (${currentUser.name})` },
    ...teachers.filter((t) => t.id !== currentUser.id).map((t) => ({ value: t.id, label: t.name })),
  ]

  const childSurname = lastName(form.name)
  const responsavelOptions = responsaveis
    .map((r) => ({
      value: r.id,
      label: r.name,
      hint: r.email,
      badge: childSurname && lastName(r.name) === childSurname ? 'Possível parente' : undefined,
    }))
    .sort((a, b) => (b.badge ? 1 : 0) - (a.badge ? 1 : 0))

  async function submit(e) {
    e.preventDefault()
    if (!form.name.trim()) {
      setError('Informe o nome da criança.')
      return
    }
    if (!form.birth_date) {
      setError('Informe a data de nascimento.')
      return
    }
    setSaving(true)
    setError('')
    try {
      const student = await createStudent({
        name: form.name,
        birth_date: new Date(`${form.birth_date}T00:00:00Z`).toISOString(),
        turma_id: form.turma_id,
        teacher_user_id: form.teacher_user_id,
        guardian_user_id: form.guardian_user_id,
        allergies: [],
      })
      onCreated(student)
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <form onSubmit={submit} className="animate-slide-up space-y-3 rounded-2xl bg-white p-5 shadow-sm ring-1 ring-slate-100">
      <h3 className="font-bold text-slate-800">Cadastrar nova criança</h3>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Nome completo</label>
        <input
          value={form.name}
          onChange={(e) => update('name', e.target.value)}
          placeholder="Nome completo"
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
      </div>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Data de nascimento</label>
        <input
          type="date"
          value={form.birth_date}
          onChange={(e) => update('birth_date', e.target.value)}
          className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
      </div>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Turma</label>
        <SearchableSelect
          options={turmaOptions}
          value={form.turma_id}
          onChange={(v) => update('turma_id', v)}
          placeholder="Selecione a turma"
          emptyMessage="Nenhuma turma cadastrada — crie em 'Gerenciar turmas'"
        />
      </div>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Professor(a)</label>
        <SearchableSelect
          options={teacherOptions}
          value={form.teacher_user_id}
          onChange={(v) => update('teacher_user_id', v)}
          placeholder="Selecione o(a) professor(a)"
        />
      </div>

      <div>
        <label className="mb-1 block text-xs font-medium text-slate-500">Responsável</label>
        <SearchableSelect
          options={responsavelOptions}
          value={form.guardian_user_id}
          onChange={(v) => update('guardian_user_id', v)}
          placeholder="Busque pelo nome do responsável"
          emptyMessage="Nenhum responsável encontrado"
        />
      </div>

      {error && <p className="text-xs text-rose-500">{error}</p>}
      <div className="flex justify-end gap-2">
        <button type="button" onClick={onCancel} className="rounded-full px-3 py-1.5 text-sm text-slate-500 hover:bg-slate-100">
          Cancelar
        </button>
        <button
          type="submit"
          disabled={saving}
          className="rounded-full bg-emerald-500 px-4 py-1.5 text-sm font-semibold text-white hover:bg-emerald-600 disabled:opacity-60"
        >
          {saving ? 'Salvando...' : 'Cadastrar'}
        </button>
      </div>
    </form>
  )
}

// tela de selecao de crianca pra contas profissionais (lista, criar nova, gerenciar turmas)
function ChildPicker({ students, currentUser, onSelect, onChildCreated, onManageTurmas }) {
  const [creating, setCreating] = useState(false)

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-xl font-bold text-slate-800">Crianças cadastradas</h2>
          <p className="text-sm text-slate-400">Selecione uma criança para ver e postar atualizações.</p>
        </div>
        {!creating && (
          <div className="flex gap-2">
            <button
              onClick={onManageTurmas}
              className="flex items-center gap-1.5 rounded-full border border-slate-200 px-4 py-2 text-sm font-semibold text-slate-600 transition hover:bg-slate-50"
            >
              <Layers className="h-4 w-4" />
              Turmas
            </button>
            <button
              onClick={() => setCreating(true)}
              className="flex items-center gap-1.5 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white shadow-md shadow-emerald-200 hover:bg-emerald-600"
            >
              <Plus className="h-4 w-4" />
              Nova criança
            </button>
          </div>
        )}
      </div>

      {creating && (
        <CreateChildForm
          currentUser={currentUser}
          onCancel={() => setCreating(false)}
          onCreated={(student) => {
            setCreating(false)
            onChildCreated(student)
            onSelect(student.id)
          }}
        />
      )}

      {students.length === 0 ? (
        <EmptyChildState message="Nenhuma criança cadastrada ainda. Use o botão acima para cadastrar a primeira." />
      ) : (
        <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
          {students.map((s, i) => (
            <button
              key={s.id}
              onClick={() => onSelect(s.id)}
              className={`animate-pop-in stagger-${Math.min(i + 1, 6)} group flex items-center gap-3 rounded-2xl bg-white p-4 text-left shadow-sm ring-1 ring-slate-100 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md hover:ring-emerald-200`}
            >
              <Avatar src={s.photo_url} alt={s.name} className="h-12 w-12 rounded-xl object-cover transition-transform duration-200 group-hover:scale-105" />
              <div className="min-w-0 flex-1">
                <p className="font-semibold text-slate-800">{s.name}</p>
                <p className="text-xs text-slate-400">{s.group_name || 'Sem turma definida'}</p>
              </div>
              <ChevronRight className="h-4 w-4 text-slate-300 transition-transform duration-200 group-hover:translate-x-1 group-hover:text-emerald-500" />
            </button>
          ))}
        </div>
      )}
    </div>
  )
}

// menu de navegacao fixo no rodape, usado no mobile
function BottomNav({ active, onChange }) {
  return (
    <nav className="fixed inset-x-0 bottom-0 z-20 border-t border-slate-100 bg-white/95 px-2 pb-[env(safe-area-inset-bottom)] backdrop-blur md:hidden">
      <ul className="flex items-stretch justify-between">
        {navItems.map(({ label, icon: Icon }) => {
          const isActive = active === label
          return (
            <li key={label} className="flex-1">
              <button
                onClick={() => onChange(label)}
                className={`flex w-full flex-col items-center gap-1 py-2.5 transition-all duration-200 ${
                  isActive ? '-translate-y-0.5 text-emerald-600' : 'text-slate-400'
                }`}
              >
                <Icon className={`h-5 w-5 transition-transform duration-200 ${isActive ? 'scale-110' : ''}`} strokeWidth={isActive ? 2.4 : 2} />
                <span className="text-[11px] font-medium">{label}</span>
              </button>
            </li>
          )
        })}
      </ul>
    </nav>
  )
}

// componente raiz: cuida de login/sessao, carrega as criancas do usuario e decide qual tela mostrar
export default function App() {
  const [active, setActive] = useState('Início')
  const [user, setUser] = useState(null)
  const [authScreen, setAuthScreen] = useState('login')
  const [checkingSession, setCheckingSession] = useState(true)

  const [students, setStudents] = useState([])
  const [selectedChildId, setSelectedChildId] = useState(null)
  const [selectedStudent, setSelectedStudent] = useState(null)
  const [timeline, setTimeline] = useState([])
  const [nextEvent, setNextEvent] = useState(null)
  const [showTurmas, setShowTurmas] = useState(false)
  const [creatingStudent, setCreatingStudent] = useState(false)

  const isProfissional = user?.role === 'profissional'

  // ao abrir o app, tenta recuperar a sessao com o token salvo (se nao for valido, desloga)
  useEffect(() => {
    async function restoreSession() {
      if (!getToken()) {
        setCheckingSession(false)
        return
      }
      try {
        const me = await getMe()
        setUser(me)
      } catch {
        setToken(null)
      } finally {
        setCheckingSession(false)
      }
    }
    restoreSession()
  }, [])

  // busca a lista de criancas certa pro papel do usuario (todas, se profissional; so as suas, se responsavel)
  const loadStudentsList = useCallback(async (currentUser) => {
    const loader = currentUser.role === 'profissional' ? listStudents : getMyChildren
    try {
      const list = await loader()
      setStudents(list ?? [])
      if (currentUser.role === 'responsavel' && list?.length) {
        setSelectedChildId((prev) => prev ?? list[0].id)
      }
      return list ?? []
    } catch (err) {
      console.error('Falha ao carregar crianças', err)
      return []
    }
  }, [])

  useEffect(() => {
    if (!user) return
    loadStudentsList(user)
  }, [user, loadStudentsList])

  // carrega os dados da crianca selecionada: perfil, timeline de hoje e proximo evento
  const loadChildData = useCallback(async (studentId) => {
    try {
      const [studentData, timelineData, eventsData] = await Promise.all([
        getStudentById(studentId),
        getTimeline(studentId),
        getEvents().catch(() => []),
      ])
      setSelectedStudent(studentData)
      setTimeline(timelineData ?? [])
      const upcoming = (eventsData ?? [])
        .filter((e) => new Date(e.starts_at) > new Date())
        .sort((a, b) => new Date(a.starts_at) - new Date(b.starts_at))
      setNextEvent(upcoming[0] ?? null)
    } catch (err) {
      console.error('Falha ao carregar dados da criança', err)
    }
  }, [])

  useEffect(() => {
    if (selectedChildId) loadChildData(selectedChildId)
    else setSelectedStudent(null)
  }, [selectedChildId, loadChildData])

  // limpa token e todo o estado da sessao, voltando pra tela de login
  function handleLogout() {
    setToken(null)
    setUser(null)
    setStudents([])
    setSelectedChildId(null)
    setSelectedStudent(null)
    setShowTurmas(false)
    setActive('Início')
  }

  // cria um evento novo na timeline e recarrega os dados da crianca pra atualizar a tela
  async function handleCreateTimelineEntry(data) {
    await createTimelineEvent(selectedChildId, data)
    await loadChildData(selectedChildId)
  }

  if (checkingSession) {
    return <div className="grid min-h-screen w-full place-items-center bg-white text-slate-400">Carregando...</div>
  }

  if (!user) {
    if (authScreen === 'forgot') {
      return <EsqueceuSenha onBack={() => setAuthScreen('login')} />
    }
    if (authScreen === 'register') {
      return <CriarConta onBack={() => setAuthScreen('login')} />
    }
    return (
      <Login
        onLogin={(loggedUser) => setUser(loggedUser)}
        onForgot={() => setAuthScreen('forgot')}
        onCreateAccount={() => setAuthScreen('register')}
      />
    )
  }

  const showingPicker = isProfissional && !selectedChildId
  const noChildLinked = !isProfissional && students.length === 0

  // decide qual tela central mostrar com base na aba ativa e no estado de selecao de crianca
  function renderMain() {
    if (showingPicker) {
      if (showTurmas) {
        return <Turmas onBack={() => setShowTurmas(false)} />
      }
      if (creatingStudent) {
        return (
          <CreateChildForm
            currentUser={user}
            onCancel={() => setCreatingStudent(false)}
            onCreated={(student) => {
              setCreatingStudent(false)
              setStudents((prev) => [...prev, student])
              setSelectedChildId(student.id)
            }}
          />
        )
      }
      return (
        <DashboardAdmin
          students={students}
          currentUser={user}
          onSelectChild={setSelectedChildId}
          onCreateChildClick={() => setCreatingStudent(true)}
          onManageTurmasClick={() => setShowTurmas(true)}
          onNavigate={setActive}
        />
      )
    }

    if (noChildLinked) {
      return (
        <EmptyChildState message="Nenhuma criança vinculada à sua conta ainda. Peça para a escola vincular seu e-mail ao cadastro do seu filho(a)." />
      )
    }

    if (active === 'Atividades') {
      return <Feed studentId={selectedChildId} canEdit={isProfissional} currentUser={user} />
    }
    if (active === 'Jornada') {
      return <Jornada studentId={selectedChildId} canEdit={isProfissional} />
    }
    if (active === 'Agenda') {
      return <Agenda canEdit={isProfissional} />
    }
    if (active === 'Comunicados') {
      return <Comunicados canEdit={isProfissional} />
    }
    if (active === 'Perfil') {
      return <Perfil user={user} onLogout={handleLogout} onOpenChild={() => setActive('Crianca')} onUserUpdated={setUser} />
    }
    if (active === 'Crianca') {
      return (
        <Crianca
          student={selectedStudent}
          canEdit={isProfissional}
          onBack={() => setActive('Perfil')}
          onUpdated={() => loadChildData(selectedChildId)}
        />
      )
    }

    return (
      <div className="grid grid-cols-1 gap-6 xl:grid-cols-3">
        <div className="space-y-6 xl:col-span-2">
          <ChildCard student={selectedStudent} onOpenChild={() => setActive('Atividades')} />
          <StatsRow student={selectedStudent} nextEvent={nextEvent} />
          <QuickAccess onOpen={setActive} />
        </div>

        <div className="xl:col-span-1">
          <DaySummary timeline={timeline} canEdit={isProfissional} onCreate={handleCreateTimelineEntry} />
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen w-full bg-slate-50">
      <Sidebar active={active} onChange={setActive} onLogout={handleLogout} showNav={!showingPicker} />

      <main className="flex-1 overflow-y-auto px-5 py-6 pb-24 sm:px-8 lg:px-10 md:pb-6">
        <Topbar
          user={user}
          onSwitchChild={isProfissional && selectedChildId ? () => { setSelectedChildId(null); loadStudentsList(user) } : null}
          onOpenProfile={() => setActive('Perfil')}
        />

        <div key={`${active}-${selectedChildId ?? 'none'}`} className="animate-fade-in mt-8">
          {renderMain()}
        </div>
      </main>

      {!showingPicker && <BottomNav active={active} onChange={setActive} />}
    </div>
  )
}
