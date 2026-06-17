import { useState, useEffect } from 'react'
import { Layers, Plus, Trash2, ArrowLeft, CalendarCheck, X } from 'lucide-react'
import { listTurmas, createTurma, deleteTurma, listStudents, saveAttendance, getTurmaAttendance, getStudentAttendanceHistory } from './api'
import defaultAvatar from '../assets/default.png'

// tela de chamada de uma turma: marca presenca por aluno num dia e mostra o historico de cada um
function ChamadaView({ turma, onBack }) {
  const [students, setStudents] = useState([])
  const [attendanceMap, setAttendanceMap] = useState({})
  const [date, setDate] = useState(() => new Date().toISOString().slice(0, 10))
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  // State for attendance history details
  const [selectedStudentHistory, setSelectedStudentHistory] = useState(null)
  const [historyRecords, setHistoryRecords] = useState([])
  const [loadingHistory, setLoadingHistory] = useState(false)

  useEffect(() => {
    let active = true
    setLoading(true)
    setError('')
    Promise.all([
      listStudents().catch(() => []),
      getTurmaAttendance(turma.id, date).catch(() => []),
    ])
      .then(([allStudents, records]) => {
        if (!active) return
        const filtered = allStudents.filter((s) => s.turma_id === turma.id)
        setStudents(filtered)

        const map = {}
        records.forEach((r) => {
          map[r.student_id] = r.status
        })
        setAttendanceMap(map)
      })
      .catch((err) => setError(err.message))
      .finally(() => {
        if (active) setLoading(false)
      })

    return () => {
      active = false
    }
  }, [turma.id, date])

  // busca o historico de presenca de um aluno especifico pra mostrar no painel lateral
  async function loadHistory(student) {
    setSelectedStudentHistory(student)
    setLoadingHistory(true)
    try {
      const data = await getStudentAttendanceHistory(student.id)
      setHistoryRecords(data ?? [])
    } catch (err) {
      console.error(err)
    } finally {
      setLoadingHistory(false)
    }
  }

  // inverte o status de presenca do aluno (otimista: atualiza a tela antes de confirmar com a API)
  async function toggleStatus(studentId) {
    const current = attendanceMap[studentId] ?? 'absent'
    const next = current === 'present' ? 'absent' : 'present'

    setAttendanceMap((prev) => ({ ...prev, [studentId]: next }))

    try {
      await saveAttendance(studentId, { date, status: next })
    } catch (err) {
      setError(err.message || 'Falha ao salvar presença.')
      setAttendanceMap((prev) => ({ ...prev, [studentId]: current }))
    }
  }

  return (
    <div className="space-y-6 animate-slide-up">
      <div className="flex items-center gap-3">
        <button
          onClick={onBack}
          className="grid h-9 w-9 place-items-center rounded-full border border-slate-200 text-slate-500 transition-all duration-200 hover:-translate-x-0.5 hover:bg-slate-50"
          aria-label="Voltar"
        >
          <ArrowLeft className="h-4 w-4" />
        </button>
        <div>
          <h2 className="text-xl font-bold text-slate-800">Chamada — {turma.name}</h2>
          <p className="text-sm text-slate-400">Registre a presença dos alunos por data</p>
        </div>
      </div>

      <div className="flex flex-wrap items-center gap-3 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100">
        <div className="flex items-center gap-2">
          <label className="text-sm font-semibold text-slate-600">Data:</label>
          <input
            type="date"
            value={date}
            onChange={(e) => setDate(e.target.value)}
            className="rounded-lg border border-slate-200 px-3 py-1.5 text-sm outline-none focus:border-emerald-400"
          />
        </div>
      </div>

      {error && <p className="rounded-2xl bg-rose-50 p-4 text-sm text-rose-600">{error}</p>}

      {loading ? (
        <p className="text-sm text-slate-400">Carregando chamada...</p>
      ) : students.length === 0 ? (
        <p className="rounded-2xl bg-white p-6 text-center text-sm text-slate-400 ring-1 ring-slate-100">
          Nenhum aluno cadastrado nesta turma.
        </p>
      ) : (
        <div className="divide-y divide-slate-100 rounded-3xl bg-white p-4 shadow-sm ring-1 ring-slate-100">
          {students.map((student) => {
            const isPresent = attendanceMap[student.id] === 'present'
            return (
              <div key={student.id} className="flex items-center justify-between py-3.5 first:pt-0 last:pb-0">
                <button
                  onClick={() => loadHistory(student)}
                  className="flex items-center gap-3 text-left hover:opacity-80 transition group"
                >
                  <img
                    src={student.photo_url || defaultAvatar}
                    alt={student.name}
                    className="h-10 w-10 rounded-xl object-cover ring-1 ring-slate-100 transition-transform duration-200 group-hover:scale-105"
                    onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
                  />
                  <div>
                    <p className="font-semibold text-slate-800 group-hover:text-emerald-600 transition-colors">{student.name}</p>
                    <p className="text-xs text-slate-400">Ver histórico de presenças</p>
                  </div>
                </button>

                <button
                  onClick={() => toggleStatus(student.id)}
                  className={`rounded-full px-4 py-1.5 text-xs font-bold transition-all duration-200 active:scale-95 ${
                    isPresent
                      ? 'bg-emerald-500 text-white shadow-md shadow-emerald-100'
                      : 'bg-rose-500 text-white shadow-md shadow-rose-100'
                  }`}
                >
                  {isPresent ? 'PRESENTE' : 'AUSENTE'}
                </button>
              </div>
            )
          })}
        </div>
      )}

      {selectedStudentHistory && (
        <div className="animate-fade-in fixed inset-0 z-50 flex items-center justify-end bg-slate-900/50 p-0 backdrop-blur-sm sm:p-4">
          <div className="animate-slide-in-right flex h-full w-full max-w-md flex-col bg-white shadow-2xl sm:rounded-3xl sm:h-[90vh]">
            <div className="flex items-center justify-between border-b border-slate-100 p-5">
              <div className="flex items-center gap-3">
                <img
                  src={selectedStudentHistory.photo_url || defaultAvatar}
                  alt={selectedStudentHistory.name}
                  className="h-10 w-10 rounded-xl object-cover"
                />
                <div>
                  <h3 className="font-bold text-slate-800">{selectedStudentHistory.name}</h3>
                  <p className="text-xs text-slate-400">Histórico de Presença</p>
                </div>
              </div>
              <button
                onClick={() => setSelectedStudentHistory(null)}
                className="rounded-full p-1 text-slate-400 hover:bg-slate-100 hover:text-slate-600"
              >
                <X className="h-5 w-5" />
              </button>
            </div>

            <div className="flex-1 overflow-y-auto p-5 space-y-6">
              {loadingHistory ? (
                <p className="text-sm text-slate-400">Carregando histórico...</p>
              ) : (
                <>
                  {/* Resumo */}
                  <div className="grid grid-cols-3 gap-3">
                    <div className="rounded-2xl bg-emerald-50 p-3 text-center ring-1 ring-emerald-100">
                      <p className="text-[10px] font-bold text-emerald-700 uppercase">Presenças</p>
                      <p className="mt-1 text-xl font-extrabold text-emerald-800">
                        {historyRecords.filter(r => r.status === 'present').length}
                      </p>
                    </div>
                    <div className="rounded-2xl bg-rose-50 p-3 text-center ring-1 ring-rose-100">
                      <p className="text-[10px] font-bold text-rose-700 uppercase">Faltas</p>
                      <p className="mt-1 text-xl font-extrabold text-rose-800">
                        {historyRecords.filter(r => r.status === 'absent').length}
                      </p>
                    </div>
                    <div className="rounded-2xl bg-indigo-50 p-3 text-center ring-1 ring-indigo-100">
                      <p className="text-[10px] font-bold text-indigo-700 uppercase">Frequência</p>
                      <p className="mt-1 text-lg font-extrabold text-indigo-800">
                        {historyRecords.length > 0 
                          ? `${Math.round((historyRecords.filter(r => r.status === 'present').length / historyRecords.length) * 100)}%` 
                          : '100%'}
                      </p>
                    </div>
                  </div>

                  {/* Lista detalhada */}
                  <div className="space-y-3">
                    <p className="text-sm font-bold text-slate-700">Histórico de Chamadas ({historyRecords.length})</p>
                    {historyRecords.length === 0 ? (
                      <p className="text-xs text-slate-400">Nenhum registro de chamada encontrado.</p>
                    ) : (
                      <div className="divide-y divide-slate-100 rounded-2xl border border-slate-100 bg-slate-50/50 p-4 max-h-[50vh] overflow-y-auto">
                        {historyRecords.map((record) => {
                          const isPresent = record.status === 'present'
                          const dateObj = new Date(record.date)
                          const localDate = new Date(dateObj.getTime() + dateObj.getTimezoneOffset() * 60000)
                          
                          return (
                            <div key={record.id || record.date} className="flex items-center justify-between py-2.5 first:pt-0 last:pb-0">
                              <span className="text-sm font-medium text-slate-600">
                                {localDate.toLocaleDateString('pt-BR', { weekday: 'short', day: '2-digit', month: '2-digit', year: 'numeric' })}
                              </span>
                              <span className={`rounded-full px-2.5 py-1 text-[10px] font-bold ${
                                isPresent ? 'bg-emerald-100 text-emerald-800' : 'bg-rose-100 text-rose-800'
                              }`}>
                                {isPresent ? 'PRESENTE' : 'FALTA'}
                              </span>
                            </div>
                          )
                        })}
                      </div>
                    )}
                  </div>
                </>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

// tela de gerenciamento de turmas: lista, cria, remove e abre a chamada de cada turma
export default function Turmas({ onBack }) {
  const [turmas, setTurmas] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [name, setName] = useState('')
  const [saving, setSaving] = useState(false)
  const [activeChamadaTurma, setActiveChamadaTurma] = useState(null)

  // busca a lista de turmas na API
  function load() {
    setLoading(true)
    listTurmas()
      .then(setTurmas)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  // cria a turma com o nome digitado e recarrega a lista
  async function submit(e) {
    e.preventDefault()
    if (!name.trim()) return
    setSaving(true)
    setError('')
    try {
      await createTurma({ name })
      setName('')
      load()
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  // confirma com o usuario e remove a turma
  async function handleDelete(id) {
    if (!window.confirm('Remover esta turma? As crianças vinculadas ficarão sem turma definida.')) return
    try {
      await deleteTurma(id)
      load()
    } catch (err) {
      setError(err.message)
    }
  }

  if (activeChamadaTurma) {
    return <ChamadaView turma={activeChamadaTurma} onBack={() => setActiveChamadaTurma(null)} />
  }

  return (
    <div className="mx-auto max-w-2xl space-y-6">
      <div className="flex items-center gap-3">
        {onBack && (
          <button
            onClick={onBack}
            className="grid h-9 w-9 place-items-center rounded-full border border-slate-200 text-slate-500 transition-all duration-200 hover:-translate-x-0.5 hover:bg-slate-50"
            aria-label="Voltar"
          >
            <ArrowLeft className="h-4 w-4" />
          </button>
        )}
        <div>
          <h2 className="flex items-center gap-2 text-xl font-bold text-slate-800">
            <Layers className="h-5 w-5 text-emerald-600" />
            Turmas
          </h2>
          <p className="text-sm text-slate-400">Gerencie as turmas da escola</p>
        </div>
      </div>

      <form onSubmit={submit} className="animate-slide-up flex gap-2 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100">
        <input
          value={name}
          onChange={(e) => setName(e.target.value)}
          placeholder="Nome da turma (ex: Berçário II)"
          className="flex-1 rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
        />
        <button
          type="submit"
          disabled={saving || !name.trim()}
          className="flex items-center gap-1.5 rounded-lg bg-emerald-500 px-4 py-2 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:opacity-60"
        >
          <Plus className="h-4 w-4" />
          Criar
        </button>
      </form>

      {error && <p className="rounded-2xl bg-rose-50 p-4 text-sm text-rose-600">{error}</p>}

      {loading ? (
        <p className="text-sm text-slate-400">Carregando...</p>
      ) : turmas.length === 0 ? (
        <p className="rounded-2xl bg-white p-6 text-center text-sm text-slate-400 ring-1 ring-slate-100">
          Nenhuma turma cadastrada ainda.
        </p>
      ) : (
        <div className="space-y-3">
          {turmas.map((t, i) => (
            <div
              key={t.id}
              className={`animate-pop-in stagger-${Math.min(i + 1, 6)} flex items-center justify-between gap-3 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md`}
            >
              <p className="font-semibold text-slate-800">{t.name}</p>
              <div className="flex items-center gap-1">
                <button
                  onClick={() => setActiveChamadaTurma(t)}
                  aria-label="Fazer chamada"
                  className="flex items-center gap-1.5 rounded-full border border-slate-200 px-3.5 py-1.5 text-xs font-semibold text-slate-600 transition hover:border-emerald-300 hover:bg-emerald-50 hover:text-emerald-600"
                >
                  <CalendarCheck className="h-3.5 w-3.5" />
                  Chamada
                </button>
                <button
                  onClick={() => handleDelete(t.id)}
                  aria-label="Remover turma"
                  className="grid h-8 w-8 place-items-center rounded-full text-slate-400 transition hover:bg-rose-50 hover:text-rose-600"
                >
                  <Trash2 className="h-4 w-4" />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
