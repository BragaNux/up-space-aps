import { useState, useEffect } from 'react'
import {
  Users,
  GraduationCap,
  CalendarCheck,
  Search,
  Plus,
  Layers,
  ChevronRight,
  CalendarDays,
  Megaphone,
} from 'lucide-react'
import defaultAvatar from '../assets/default.png'
import { listTurmas, getEvents, getAnnouncements } from './api'

// foto com fallback pro avatar padrao se a imagem der erro ou nao existir
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

// painel administrativo: estatisticas gerais, busca de alunos por turma, proximo evento e ultimo comunicado
export default function DashboardAdmin({
  students,
  currentUser,
  onSelectChild,
  onCreateChildClick,
  onManageTurmasClick,
  onNavigate,
}) {
  const [search, setSearch] = useState('')
  const [selectedTurmaId, setSelectedTurmaId] = useState('all')
  const [turmas, setTurmas] = useState([])
  const [nextEvent, setNextEvent] = useState(null)
  const [latestAnnouncement, setLatestAnnouncement] = useState(null)

  useEffect(() => {
    listTurmas().then(setTurmas).catch(() => {})
    
    getEvents()
      .then((events) => {
        const upcoming = (events ?? [])
          .filter((e) => new Date(e.starts_at) > new Date())
          .sort((a, b) => new Date(a.starts_at) - new Date(b.starts_at))
        setNextEvent(upcoming[0] ?? null)
      })
      .catch(() => {})

    getAnnouncements()
      .then((announcements) => {
        setLatestAnnouncement(announcements?.[0] ?? null)
      })
      .catch(() => {})
  }, [])

  // Stats calculation
  const totalStudents = students.length
  const totalTurmas = turmas.length
  const presentToday = students.filter((s) => s.presence_status === 'present').length
  const absentToday = students.filter((s) => s.presence_status === 'absent').length

  // Filtering students
  const filteredStudents = students.filter((student) => {
    const matchesSearch = student.name.toLowerCase().includes(search.toLowerCase())
    const matchesTurma = selectedTurmaId === 'all' || student.turma_id === Number(selectedTurmaId)
    return matchesSearch && matchesTurma
  })

  return (
    <div className="space-y-6">
      {/* Header Row */}
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <h2 className="text-xl font-bold text-slate-800">Painel de Controle</h2>
          <p className="text-sm text-slate-400">Visão geral do espaço escolar e gestão</p>
        </div>
        <div className="flex gap-2">
          <button
            onClick={onManageTurmasClick}
            className="flex items-center gap-1.5 rounded-full border border-slate-200 bg-white px-4 py-2 text-sm font-semibold text-slate-600 transition hover:bg-slate-50"
          >
            <Layers className="h-4 w-4 text-slate-400" />
            Turmas
          </button>
          <button
            onClick={onCreateChildClick}
            className="flex items-center gap-1.5 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white shadow-md shadow-emerald-200 hover:bg-emerald-600"
          >
            <Plus className="h-4 w-4" />
            Nova criança
          </button>
        </div>
      </div>

      {/* Stats Cards Grid */}
      <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
        <div className="animate-pop-in stagger-1 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-blue-50 text-blue-500">
            <Users className="h-5 w-5" />
          </div>
          <p className="mt-3 text-xs font-semibold uppercase tracking-wider text-slate-400">Crianças</p>
          <p className="mt-1 text-2xl font-extrabold text-slate-800">{totalStudents}</p>
        </div>

        <div className="animate-pop-in stagger-2 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-purple-50 text-purple-500">
            <GraduationCap className="h-5 w-5" />
          </div>
          <p className="mt-3 text-xs font-semibold uppercase tracking-wider text-slate-400">Turmas</p>
          <p className="mt-1 text-2xl font-extrabold text-slate-800">{totalTurmas}</p>
        </div>

        <div className="animate-pop-in stagger-3 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-emerald-50 text-emerald-500">
            <CalendarCheck className="h-5 w-5" />
          </div>
          <p className="mt-3 text-xs font-semibold uppercase tracking-wider text-slate-400">Presentes Hoje</p>
          <p className="mt-1 text-2xl font-extrabold text-emerald-600">{presentToday}</p>
        </div>

        <div className="animate-pop-in stagger-4 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-rose-50 text-rose-500">
            <CalendarCheck className="h-5 w-5" />
          </div>
          <p className="mt-3 text-xs font-semibold uppercase tracking-wider text-slate-400">Ausentes Hoje</p>
          <p className="mt-1 text-2xl font-extrabold text-rose-600">{absentToday}</p>
        </div>
      </div>

      {/* Main Content Layout */}
      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Left 2 Columns: Searchable Students List */}
        <div className="lg:col-span-2 space-y-4">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between bg-white p-4 rounded-2xl shadow-sm ring-1 ring-slate-100">
            <div className="relative flex-1">
              <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-slate-400" />
              <input
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Buscar criança..."
                className="w-full rounded-lg border border-slate-200 py-2 pl-9 pr-4 text-sm outline-none transition focus:border-emerald-400 focus:ring-1 focus:ring-emerald-100"
              />
            </div>
            
            <select
              value={selectedTurmaId}
              onChange={(e) => setSelectedTurmaId(e.target.value)}
              className="rounded-lg border border-slate-200 bg-white px-3 py-2 text-sm outline-none focus:border-emerald-400 cursor-pointer"
            >
              <option value="all">Todas as turmas</option>
              {turmas.map((t) => (
                <option key={t.id} value={t.id}>{t.name}</option>
              ))}
            </select>
          </div>

          {filteredStudents.length === 0 ? (
            <p className="rounded-2xl bg-white p-12 text-center text-sm text-slate-400 ring-1 ring-slate-100">
              Nenhuma criança encontrada com os filtros selecionados.
            </p>
          ) : (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              {filteredStudents.map((s, i) => (
                <button
                  key={s.id}
                  onClick={() => onSelectChild(s.id)}
                  className={`animate-pop-in stagger-${Math.min(i + 1, 6)} group flex items-center gap-3 rounded-2xl bg-white p-4 text-left shadow-sm ring-1 ring-slate-100 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md hover:ring-emerald-200`}
                >
                  <Avatar src={s.photo_url} alt={s.name} className="h-12 w-12 rounded-xl object-cover transition-transform duration-200 group-hover:scale-105" />
                  <div className="min-w-0 flex-1">
                    <p className="font-semibold text-slate-800">{s.name}</p>
                    <p className="text-xs text-slate-400">{s.group_name || 'Sem turma'}</p>
                  </div>
                  <ChevronRight className="h-4 w-4 text-slate-300 transition-transform duration-200 group-hover:translate-x-1 group-hover:text-emerald-500" />
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Right 1 Column: School Info Quick Access */}
        <div className="space-y-6">
          {/* Quick Agenda */}
          <div className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-slate-100">
            <h3 className="flex items-center gap-2 font-bold text-slate-800">
              <CalendarDays className="h-4 w-4 text-emerald-600" />
              Agenda Próxima
            </h3>
            {nextEvent ? (
              <div className="mt-3 rounded-xl bg-slate-50 p-3.5">
                <h4 className="font-semibold text-sm text-slate-800">{nextEvent.title}</h4>
                <p className="mt-1 text-xs text-slate-500 line-clamp-2">{nextEvent.description}</p>
                <p className="mt-2 text-xs font-bold text-emerald-600">
                  {new Date(nextEvent.starts_at).toLocaleDateString('pt-BR', { day: '2-digit', month: 'long', hour: '2-digit', minute: '2-digit' })}
                </p>
                <button onClick={() => onNavigate('Agenda')} className="mt-3 text-xs font-semibold text-emerald-600 hover:underline">
                  Ver agenda completa
                </button>
              </div>
            ) : (
              <p className="mt-3 text-xs text-slate-400">Sem eventos marcados próximos.</p>
            )}
          </div>

          {/* Quick Announcements */}
          <div className="rounded-2xl bg-white p-5 shadow-sm ring-1 ring-slate-100">
            <h3 className="flex items-center gap-2 font-bold text-slate-800">
              <Megaphone className="h-4 w-4 text-emerald-600" />
              Último Comunicado
            </h3>
            {latestAnnouncement ? (
              <div className="mt-3 rounded-xl bg-slate-50 p-3.5">
                <span className="inline-block rounded-full bg-indigo-50 px-2 py-0.5 text-[10px] font-semibold text-indigo-600">
                  {latestAnnouncement.priority}
                </span>
                <h4 className="mt-1.5 font-semibold text-sm text-slate-800">{latestAnnouncement.title}</h4>
                <p className="mt-1 text-xs text-slate-500 line-clamp-2">{latestAnnouncement.preview}</p>
                <button onClick={() => onNavigate('Comunicados')} className="mt-3 text-xs font-semibold text-emerald-600 hover:underline">
                  Ver comunicados
                </button>
              </div>
            ) : (
              <p className="mt-3 text-xs text-slate-400">Nenhum comunicado publicado.</p>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
