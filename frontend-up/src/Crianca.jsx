import { useState, useRef, useEffect } from 'react'
import {
  ArrowLeft,
  Cake,
  GraduationCap,
  IdCard,
  Droplet,
  AlertCircle,
  Utensils,
  Pill,
  Phone,
  ShieldCheck,
  Heart,
  Camera,
  Plus,
  Pencil,
  Trash2,
} from 'lucide-react'
import defaultAvatar from '../assets/default.png'
import SearchableSelect from './SearchableSelect'
import { updateStudent, createGuardian, updateGuardian, deleteGuardian, fileToBase64, listTurmas, listUsersByRole } from './api'

const MAX_PHOTO_BYTES = 4 * 1024 * 1024

// mini card de estatistica (idade, turma, matricula) no topo do perfil
function Stat({ icon: Icon, label, value }) {
  return (
    <div className="rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100">
      <span className="grid h-9 w-9 place-items-center rounded-lg bg-emerald-50 text-emerald-600">
        <Icon className="h-4 w-4" />
      </span>
      <p className="mt-2 text-xs text-slate-400">{label}</p>
      <p className="text-sm font-bold text-slate-800">{value || '—'}</p>
    </div>
  )
}

// linha da secao de saude (tipo sanguineo, alergias, restricoes, medicamentos)
function HealthRow({ icon: Icon, label, children, tone = 'slate' }) {
  const tones = {
    slate: 'bg-slate-50 text-slate-500',
    rose: 'bg-rose-50 text-rose-600',
    amber: 'bg-amber-50 text-amber-600',
  }
  return (
    <div className="flex items-start gap-3 py-3">
      <span className={`grid h-9 w-9 shrink-0 place-items-center rounded-lg ${tones[tone]}`}>
        <Icon className="h-4 w-4" />
      </span>
      <div className="min-w-0 flex-1">
        <p className="text-xs text-slate-400">{label}</p>
        <div className="mt-0.5 text-sm font-medium text-slate-700">{children}</div>
      </div>
    </div>
  )
}

// calcula a idade da crianca em anos (ou em meses, se tiver menos de 1 ano)
function calcAge(birthDate) {
  if (!birthDate) return '—'
  const birth = new Date(birthDate)
  const now = new Date()
  let years = now.getFullYear() - birth.getFullYear()
  let months = now.getMonth() - birth.getMonth()
  if (months < 0) {
    years -= 1
    months += 12
  }
  if (years > 0) return `${years} ano${years > 1 ? 's' : ''}`
  return `${months} ${months === 1 ? 'mês' : 'meses'}`
}

// formulario de criar/editar um responsavel/autorizado a buscar a crianca
function GuardianForm({ initial, onSave, onCancel }) {
  const [name, setName] = useState(initial?.name ?? '')
  const [relation, setRelation] = useState(initial?.relation ?? '')
  const [phone, setPhone] = useState(initial?.phone ?? '')
  const [authorized, setAuthorized] = useState(initial?.authorized ?? false)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  async function submit(e) {
    e.preventDefault()
    if (!name.trim() || !relation.trim()) {
      setError('Preencha nome e parentesco.')
      return
    }
    setSaving(true)
    setError('')
    try {
      await onSave({ name, relation, phone, authorized, avatar_url: initial?.avatar_url ?? '' })
    } catch (err) {
      setError(err.message)
      setSaving(false)
    }
  }

  return (
    <form onSubmit={submit} className="animate-slide-up space-y-2 rounded-2xl bg-white p-4 shadow-sm ring-1 ring-slate-100">
      <input
        value={name}
        onChange={(e) => setName(e.target.value)}
        placeholder="Nome completo"
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />
      <input
        value={relation}
        onChange={(e) => setRelation(e.target.value)}
        placeholder="Parentesco (ex: Pai, Mãe, Avó)"
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />
      <input
        value={phone}
        onChange={(e) => setPhone(e.target.value)}
        placeholder="Telefone"
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />
      <label className="flex items-center gap-2 text-sm text-slate-600">
        <input type="checkbox" checked={authorized} onChange={(e) => setAuthorized(e.target.checked)} />
        Autorizado a buscar a criança
      </label>
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
          {saving ? 'Salvando...' : 'Salvar'}
        </button>
      </div>
    </form>
  )
}

// tela de perfil completo da crianca: dados escolares, saude e responsaveis (tudo editavel se canEdit)
export default function Crianca({ student, canEdit, onBack, onUpdated }) {
  const [uploading, setUploading] = useState(false)
  const [photoError, setPhotoError] = useState('')
  const [addingGuardian, setAddingGuardian] = useState(false)
  const [editingGuardianId, setEditingGuardianId] = useState(null)
  const [editingHealth, setEditingHealth] = useState(false)
  const [healthForm, setHealthForm] = useState({ blood_type: '', allergies: '', restrictions: '', medications: '' })
  const [savingHealth, setSavingHealth] = useState(false)
  const [healthError, setHealthError] = useState('')
  const fileInputRef = useRef(null)

  const [editingSchool, setEditingSchool] = useState(false)
  const [schoolForm, setSchoolForm] = useState({ turma_id: null, teacher_user_id: null })
  const [savingSchool, setSavingSchool] = useState(false)
  const [schoolError, setSchoolError] = useState('')
  const [turmas, setTurmas] = useState([])
  const [professionals, setProfessionals] = useState([])

  useEffect(() => {
    // ao entrar no modo de edicao escolar, busca as opcoes de turma e professores
    if (editingSchool) {
      async function loadOptions() {
        try {
          const [tList, pList] = await Promise.all([
            listTurmas(),
            listUsersByRole('profissional'),
          ])
          setTurmas(tList || [])
          setProfessionals(pList || [])
        } catch (err) {
          setSchoolError('Erro ao carregar opções: ' + err.message)
        }
      }
      loadOptions()
    }
  }, [editingSchool])

  // preenche o formulario escolar com os dados atuais e entra no modo de edicao
  function startEditingSchool() {
    setSchoolForm({
      turma_id: student.turma_id || null,
      teacher_user_id: student.teacher_user_id || null,
    })
    setSchoolError('')
    setEditingSchool(true)
  }

  // salva a turma e o professor escolhidos pra crianca
  async function handleSaveSchool(e) {
    e.preventDefault()
    setSavingSchool(true)
    setSchoolError('')
    try {
      await updateStudent(student.id, buildStudentPayload({
        turma_id: schoolForm.turma_id,
        teacher_user_id: schoolForm.teacher_user_id,
      }))
      setEditingSchool(false)
      onUpdated?.()
    } catch (err) {
      setSchoolError(err.message || 'Não foi possível salvar as informações escolares.')
    } finally {
      setSavingSchool(false)
    }
  }

  const turmaOptions = turmas.map((t) => ({ value: t.id, label: t.name }))
  const teacherOptions = professionals.map((p) => ({ value: p.id, label: p.name, hint: p.email }))

  // preenche o formulario de saude com os dados atuais e entra no modo de edicao
  function startEditingHealth() {
    setHealthForm({
      blood_type: student.blood_type || '',
      allergies: student.allergies?.join(', ') || '',
      restrictions: student.restrictions || '',
      medications: student.medications || '',
    })
    setHealthError('')
    setEditingHealth(true)
  }

  // converte as alergias digitadas (separadas por virgula) em lista e salva os dados de saude
  async function handleSaveHealth(e) {
    e.preventDefault()
    setSavingHealth(true)
    setHealthError('')
    try {
      const allergiesArray = healthForm.allergies
        ? healthForm.allergies.split(',').map((s) => s.trim()).filter(Boolean)
        : []
      await updateStudent(student.id, buildStudentPayload({
        blood_type: healthForm.blood_type.trim(),
        allergies: allergiesArray,
        restrictions: healthForm.restrictions.trim(),
        medications: healthForm.medications.trim(),
      }))
      setEditingHealth(false)
      onUpdated?.()
    } catch (err) {
      setHealthError(err.message || 'Não foi possível salvar os dados de saúde.')
    } finally {
      setSavingHealth(false)
    }
  }

  if (!student) {
    return (
      <div className="mx-auto max-w-2xl">
        <p className="rounded-2xl bg-white p-6 text-center text-sm text-slate-400 ring-1 ring-slate-100">
          Carregando dados da criança...
        </p>
      </div>
    )
  }

  const guardians = student.guardians ?? []

  // monta o payload completo de atualizacao do aluno, usando os dados atuais como base e so trocando o que for passado em overrides
  function buildStudentPayload(overrides) {
    return {
      name: student.name,
      guardian_user_id: student.guardian_user_id ?? null,
      photo_url: student.photo_url ?? '',
      group_name: student.group_name ?? '',
      teacher_name: student.teacher_name ?? '',
      birth_date: student.birth_date ?? null,
      enrollment_code: student.enrollment_code ?? '',
      blood_type: student.blood_type ?? '',
      allergies: student.allergies ?? [],
      restrictions: student.restrictions ?? '',
      medications: student.medications ?? '',
      turma_id: student.turma_id ?? null,
      teacher_user_id: student.teacher_user_id ?? null,
      ...overrides,
    }
  }

  // valida tipo/tamanho da foto e ja salva ela no cadastro do aluno
  async function handlePhotoChange(e) {
    const file = e.target.files?.[0]
    if (!file) return
    setPhotoError('')
    if (!file.type.startsWith('image/')) {
      setPhotoError('Selecione um arquivo de imagem.')
      return
    }
    if (file.size > MAX_PHOTO_BYTES) {
      setPhotoError('Imagem muito grande (máximo 4MB).')
      return
    }
    setUploading(true)
    try {
      const base64 = await fileToBase64(file)
      await updateStudent(student.id, buildStudentPayload({ photo_url: base64 }))
      onUpdated?.()
    } catch (err) {
      setPhotoError(err.message || 'Não foi possível atualizar a foto.')
    } finally {
      setUploading(false)
      e.target.value = ''
    }
  }

  // cadastra um responsavel novo pra crianca
  async function handleAddGuardian(data) {
    await createGuardian(student.id, data)
    setAddingGuardian(false)
    onUpdated?.()
  }

  // atualiza os dados de um responsavel existente
  async function handleEditGuardian(id, data) {
    await updateGuardian(id, data)
    setEditingGuardianId(null)
    onUpdated?.()
  }

  // confirma com o usuario e remove o responsavel
  async function handleDeleteGuardian(id) {
    if (!window.confirm('Remover este responsável/autorizado?')) return
    await deleteGuardian(id)
    onUpdated?.()
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
        <h2 className="text-xl font-bold text-slate-800">Perfil da Criança</h2>
      </div>

      <section className="animate-slide-up rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-100">
        <div className="flex flex-col items-center text-center sm:flex-row sm:items-center sm:gap-5 sm:text-left">
          {canEdit ? (
            <button
              type="button"
              onClick={() => fileInputRef.current?.click()}
              disabled={uploading}
              className="group relative shrink-0 rounded-3xl"
              aria-label="Alterar foto da criança"
            >
              <img
                src={student.photo_url || defaultAvatar}
                alt={student.name}
                onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
                className="h-24 w-24 rounded-3xl object-cover ring-4 ring-emerald-50 transition-transform duration-300 group-hover:scale-105"
              />
              <span className="absolute inset-0 grid place-items-center rounded-3xl bg-slate-900/0 text-white opacity-0 transition-all duration-200 group-hover:bg-slate-900/40 group-hover:opacity-100">
                {uploading ? (
                  <span className="h-5 w-5 animate-spin rounded-full border-2 border-white/40 border-t-white" />
                ) : (
                  <Camera className="h-5 w-5" />
                )}
              </span>
              <input ref={fileInputRef} type="file" accept="image/*" onChange={handlePhotoChange} className="hidden" />
            </button>
          ) : (
            <img
              src={student.photo_url || defaultAvatar}
              alt={student.name}
              onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
              className="h-24 w-24 rounded-3xl object-cover ring-4 ring-emerald-50"
            />
          )}
          <div className="mt-4 sm:mt-0">
            <h3 className="text-2xl font-extrabold text-slate-800">{student.name}</h3>
            <p className="text-slate-500">{student.group_name} · {student.teacher_name}</p>
            <div className="mt-2 flex items-center justify-center gap-2 sm:justify-start">
              <span className={`h-2.5 w-2.5 rounded-full ${student.presence_status === 'present' ? 'bg-emerald-500' : 'bg-slate-300'}`} />
              <span className="text-sm font-semibold tracking-wide text-emerald-600">
                {student.presence_status === 'present' ? 'PRESENTE' : 'AUSENTE'}
              </span>
            </div>
          </div>
        </div>
        {photoError && <p className="mt-2 text-xs text-rose-500">{photoError}</p>}
      </section>

      <section className="grid grid-cols-3 gap-4">
        <div className="animate-pop-in stagger-1">
          <Stat icon={Cake} label="Idade" value={calcAge(student.birth_date)} />
        </div>
        <div className="animate-pop-in stagger-2">
          <Stat icon={GraduationCap} label="Turma" value={student.group_name} />
        </div>
        <div className="animate-pop-in stagger-3">
          <Stat icon={IdCard} label="Matrícula" value={student.enrollment_code} />
        </div>
      </section>

      <section className="animate-slide-up rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-100 space-y-4">
        <div className="flex items-center justify-between border-b border-slate-100 pb-3">
          <h3 className="flex items-center gap-2 text-sm font-bold uppercase tracking-wide text-slate-400">
            <GraduationCap className="h-4 w-4 text-emerald-500" />
            Informações Escolares
          </h3>
          {canEdit && !editingSchool && (
            <button
              onClick={startEditingSchool}
              className="flex items-center gap-1 text-sm font-medium text-emerald-600 hover:underline"
            >
              <Pencil className="h-3.5 w-3.5" />
              Editar
            </button>
          )}
        </div>

        {editingSchool ? (
          <form onSubmit={handleSaveSchool} className="space-y-4">
            <div>
              <label className="mb-1 block text-xs font-medium text-slate-500">Turma</label>
              <SearchableSelect
                options={turmaOptions}
                value={schoolForm.turma_id}
                onChange={(val) => setSchoolForm((prev) => ({ ...prev, turma_id: val }))}
                placeholder="Selecionar Turma"
                emptyMessage="Nenhuma turma encontrada"
              />
            </div>

            <div>
              <label className="mb-1 block text-xs font-medium text-slate-500">Professor(a) Responsável</label>
              <SearchableSelect
                options={teacherOptions}
                value={schoolForm.teacher_user_id}
                onChange={(val) => setSchoolForm((prev) => ({ ...prev, teacher_user_id: val }))}
                placeholder="Selecionar Professor(a)"
                emptyMessage="Nenhum profissional encontrado"
              />
            </div>

            {schoolError && <p className="text-xs text-rose-500">{schoolError}</p>}

            <div className="flex justify-end gap-2 pt-2">
              <button
                type="button"
                onClick={() => setEditingSchool(false)}
                className="rounded-full px-3 py-1.5 text-sm text-slate-500 hover:bg-slate-100"
              >
                Cancelar
              </button>
              <button
                type="submit"
                disabled={savingSchool}
                className="rounded-full bg-emerald-500 px-4 py-1.5 text-sm font-semibold text-white hover:bg-emerald-600 disabled:opacity-60"
              >
                {savingSchool ? 'Salvando...' : 'Salvar'}
              </button>
            </div>
          </form>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p className="text-xs text-slate-400">Turma / Grupo</p>
              <p className="text-sm font-semibold text-slate-800">{student.group_name || 'Nenhuma turma vinculada'}</p>
            </div>
            <div>
              <p className="text-xs text-slate-400">Professor(a) Responsável</p>
              <p className="text-sm font-semibold text-slate-800">{student.teacher_name || 'Nenhum professor responsável'}</p>
            </div>
          </div>
        )}
      </section>

      <section>
        <div className="mb-3 flex items-center justify-between">
          <h3 className="flex items-center gap-2 text-sm font-bold uppercase tracking-wide text-slate-400">
            <Heart className="h-4 w-4 text-rose-400" />
            Saúde
          </h3>
          {canEdit && !editingHealth && (
            <button
              onClick={startEditingHealth}
              className="flex items-center gap-1 text-sm font-medium text-emerald-600 hover:underline"
            >
              <Pencil className="h-3.5 w-3.5" />
              Editar
            </button>
          )}
        </div>

        {editingHealth ? (
          <form onSubmit={handleSaveHealth} className="animate-slide-up space-y-3 rounded-2xl bg-white p-5 shadow-sm ring-1 ring-slate-100">
            <div>
              <label className="mb-1 block text-xs font-medium text-slate-500">Tipo Sanguíneo</label>
              <input
                value={healthForm.blood_type}
                onChange={(e) => setHealthForm((prev) => ({ ...prev, blood_type: e.target.value }))}
                placeholder="Ex: A+, O-"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
              />
            </div>

            <div>
              <label className="mb-1 block text-xs font-medium text-slate-500">Alergias (separadas por vírgula)</label>
              <input
                value={healthForm.allergies}
                onChange={(e) => setHealthForm((prev) => ({ ...prev, allergies: e.target.value }))}
                placeholder="Ex: Glúten, Lactose, Amendoim"
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
              />
            </div>

            <div>
              <label className="mb-1 block text-xs font-medium text-slate-500">Restrições Alimentares</label>
              <textarea
                value={healthForm.restrictions}
                onChange={(e) => setHealthForm((prev) => ({ ...prev, restrictions: e.target.value }))}
                placeholder="Ex: Não consome açúcar industrializado"
                rows={2}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
              />
            </div>

            <div>
              <label className="mb-1 block text-xs font-medium text-slate-500">Medicamentos de uso contínuo</label>
              <textarea
                value={healthForm.medications}
                onChange={(e) => setHealthForm((prev) => ({ ...prev, medications: e.target.value }))}
                placeholder="Ex: Nootropil 2x ao dia"
                rows={2}
                className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
              />
            </div>

            {healthError && <p className="text-xs text-rose-500">{healthError}</p>}

            <div className="flex justify-end gap-2">
              <button
                type="button"
                onClick={() => setEditingHealth(false)}
                className="rounded-full px-3 py-1.5 text-sm text-slate-500 hover:bg-slate-100"
              >
                Cancelar
              </button>
              <button
                type="submit"
                disabled={savingHealth}
                className="rounded-full bg-emerald-500 px-4 py-1.5 text-sm font-semibold text-white hover:bg-emerald-600 disabled:opacity-60"
              >
                {savingHealth ? 'Salvando...' : 'Salvar'}
              </button>
            </div>
          </form>
        ) : (
          <div className="divide-y divide-slate-100 rounded-2xl bg-white px-4 shadow-sm ring-1 ring-slate-100">
            <HealthRow icon={Droplet} label="Tipo sanguíneo">{student.blood_type || '—'}</HealthRow>
            <HealthRow icon={AlertCircle} label="Alergias" tone="rose">
              {student.allergies?.length ? (
                <div className="flex flex-wrap gap-1.5">
                  {student.allergies.map((a) => (
                    <span key={a} className="rounded-full bg-rose-50 px-2.5 py-0.5 text-xs font-semibold text-rose-600">
                      {a}
                    </span>
                  ))}
                </div>
              ) : (
                'Nenhuma registrada'
              )}
            </HealthRow>
            <HealthRow icon={Utensils} label="Restrições alimentares" tone="amber">
              {student.restrictions || 'Nenhuma'}
            </HealthRow>
            <HealthRow icon={Pill} label="Medicamentos">{student.medications || 'Nenhum de uso contínuo.'}</HealthRow>
          </div>
        )}
      </section>

      <section>
        <div className="mb-3 flex items-center justify-between">
          <h3 className="flex items-center gap-2 text-sm font-bold uppercase tracking-wide text-slate-400">
            <ShieldCheck className="h-4 w-4 text-emerald-500" />
            Responsáveis e autorizados
          </h3>
          {canEdit && !addingGuardian && (
            <button
              onClick={() => setAddingGuardian(true)}
              className="flex items-center gap-1 text-sm font-medium text-emerald-600 hover:underline"
            >
              <Plus className="h-4 w-4" />
              Adicionar
            </button>
          )}
        </div>

        {addingGuardian && (
          <div className="mb-3">
            <GuardianForm onSave={handleAddGuardian} onCancel={() => setAddingGuardian(false)} />
          </div>
        )}

        {guardians.length === 0 && !addingGuardian ? (
          <p className="rounded-2xl bg-white p-4 text-sm text-slate-400 shadow-sm ring-1 ring-slate-100">
            Nenhum responsável cadastrado.
          </p>
        ) : (
          <div className="space-y-3">
            {guardians.map((g, i) =>
              editingGuardianId === g.id ? (
                <GuardianForm
                  key={g.id}
                  initial={g}
                  onSave={(data) => handleEditGuardian(g.id, data)}
                  onCancel={() => setEditingGuardianId(null)}
                />
              ) : (
                <div
                  key={g.id}
                  className={`animate-pop-in stagger-${Math.min(i + 1, 6)} flex items-center gap-3 rounded-2xl bg-white p-3 shadow-sm ring-1 ring-slate-100 transition-all duration-200 hover:-translate-y-0.5 hover:shadow-md`}
                >
                  <img
                    src={g.avatar_url || defaultAvatar}
                    alt={g.name}
                    onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
                    className="h-12 w-12 rounded-full object-cover"
                  />
                  <div className="min-w-0 flex-1">
                    <p className="font-semibold text-slate-800">{g.name}</p>
                    <p className="text-xs text-slate-400">{g.relation}</p>
                  </div>
                  <div className="flex flex-col items-end gap-1">
                    <span className="inline-flex items-center gap-1.5 text-sm text-slate-500">
                      <Phone className="h-3.5 w-3.5 text-slate-400" />
                      {g.phone}
                    </span>
                    <span
                      className={`rounded-full px-2 py-0.5 text-[11px] font-semibold ${
                        g.authorized ? 'bg-emerald-50 text-emerald-600' : 'bg-slate-100 text-slate-400'
                      }`}
                    >
                      {g.authorized ? 'Autorizado a buscar' : 'Não autorizado'}
                    </span>
                  </div>
                  {canEdit && (
                    <div className="flex shrink-0 gap-1">
                      <button
                        onClick={() => setEditingGuardianId(g.id)}
                        aria-label="Editar"
                        className="grid h-8 w-8 place-items-center rounded-full text-slate-400 transition hover:bg-slate-100 hover:text-emerald-600"
                      >
                        <Pencil className="h-4 w-4" />
                      </button>
                      <button
                        onClick={() => handleDeleteGuardian(g.id)}
                        aria-label="Remover"
                        className="grid h-8 w-8 place-items-center rounded-full text-slate-400 transition hover:bg-rose-50 hover:text-rose-600"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  )}
                </div>
              )
            )}
          </div>
        )}
      </section>
    </div>
  )
}
