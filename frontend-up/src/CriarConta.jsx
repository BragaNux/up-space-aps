import { useState } from 'react'
import { ArrowLeft, Eye, EyeOff, GraduationCap, Heart, CheckCircle2 } from 'lucide-react'
import logoUpEspaco from '../assets/logo_up_espaco.png'
import { register as registerRequest } from './api'

const roles = [
  {
    id: 'profissional',
    label: 'Profissional',
    desc: 'Professoras e equipe da escola',
    icon: GraduationCap,
  },
  {
    id: 'responsavel',
    label: 'Responsável',
    desc: 'Mãe, pai ou responsável pelo aluno',
    icon: Heart,
  },
]

// tela de cadastro: escolhe o papel, preenche os dados e cria a conta na API
export default function CriarConta({ onBack, onCreated }) {
  const [role, setRole] = useState(null)
  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [created, setCreated] = useState(false)

  // valida os campos (senha confirmada etc) e registra a conta na API
  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    if (!name || !email || !password || !confirmPassword) {
      setError('Preencha todos os campos.')
      return
    }
    if (password !== confirmPassword) {
      setError('As senhas não coincidem.')
      return
    }
    setLoading(true)
    try {
      await registerRequest(name, email, password, role)
      setCreated(true)
    } catch (err) {
      setError(err.message || 'Não foi possível criar a conta.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen w-full items-center justify-center overflow-y-auto bg-white p-4 sm:p-6">
      <div className="relative grid w-full max-w-4xl overflow-hidden rounded-[28px] bg-white shadow-2xl shadow-slate-200/60 ring-1 ring-slate-100 md:grid-cols-2">
        {/* Painel marca */}
        <div className="relative flex flex-col bg-slate-50 p-6 md:p-8">
          <div className="flex flex-1 items-center justify-center py-2">
            <div className="h-64 w-64 sm:h-72 sm:w-72 md:h-72 md:w-72">
              <img src={logoUpEspaco} alt="Up Espaço" className="h-full w-full object-contain" />
            </div>
          </div>
        </div>

        {/* Painel formulário */}
        <div className="flex flex-col justify-center px-7 py-9 sm:px-10">
          {created ? (
            <div className="flex flex-col items-center text-center">
              <span className="grid h-16 w-16 place-items-center rounded-full bg-slate-100 text-slate-900">
                <CheckCircle2 className="h-8 w-8" />
              </span>
              <h1 className="mt-5 text-2xl font-extrabold text-slate-800">
                Conta criada!
              </h1>
              <p className="mt-2 text-sm leading-relaxed text-slate-500">
                Sua conta de{' '}
                <span className="font-semibold text-slate-700">
                  {roles.find((r) => r.id === role)?.label}
                </span>{' '}
                foi criada com sucesso.
              </p>

              <button
                onClick={onBack}
                className="mt-7 w-full rounded-lg bg-slate-900 py-3.5 text-sm font-semibold text-white shadow-md shadow-slate-300 transition hover:bg-slate-800 active:scale-[0.98]"
              >
                Ir para o login
              </button>
            </div>
          ) : !role ? (
            <>
              <button
                onClick={onBack}
                className="mb-6 inline-flex items-center gap-1.5 self-start text-sm font-medium text-slate-500 transition hover:text-slate-800"
              >
                <ArrowLeft className="h-4 w-4" />
                Voltar
              </button>

              <h1 className="text-3xl font-extrabold text-slate-800">
                Criar conta
              </h1>
              <p className="mt-2 text-sm leading-relaxed text-slate-500">
                Selecione como você vai entrar na Up Espaço.
              </p>

              <div className="mt-7 space-y-3">
                {roles.map(({ id, label, desc, icon: Icon }) => (
                  <button
                    key={id}
                    onClick={() => setRole(id)}
                    className="flex w-full items-center gap-4 rounded-xl border border-slate-200 bg-white p-4 text-left shadow-sm transition hover:border-emerald-400 hover:bg-emerald-50"
                  >
                    <span className="grid h-12 w-12 shrink-0 place-items-center rounded-xl bg-slate-100 text-slate-700">
                      <Icon className="h-6 w-6" strokeWidth={1.8} />
                    </span>
                    <div>
                      <p className="font-semibold text-slate-800">{label}</p>
                      <p className="text-xs text-slate-500">{desc}</p>
                    </div>
                  </button>
                ))}
              </div>
            </>
          ) : (
            <>
              <button
                onClick={() => { setRole(null); setError('') }}
                className="mb-6 inline-flex items-center gap-1.5 self-start text-sm font-medium text-slate-500 transition hover:text-slate-800"
              >
                <ArrowLeft className="h-4 w-4" />
                Voltar
              </button>

              <h1 className="text-3xl font-extrabold text-slate-800">
                Criar conta
              </h1>
              <p className="mt-2 text-sm leading-relaxed text-slate-500">
                Cadastro de {roles.find((r) => r.id === role)?.label.toLowerCase()}.
              </p>

              <form onSubmit={handleSubmit} noValidate className="mt-7 space-y-5">
                <div>
                  <label className="mb-1.5 block text-sm font-semibold text-slate-700">
                    Nome completo
                  </label>
                  <input
                    type="text"
                    placeholder="Seu nome"
                    value={name}
                    onChange={(e) => { setName(e.target.value); setError('') }}
                    className="w-full rounded-lg border border-black bg-white px-4 py-3 text-sm text-slate-800 placeholder-slate-300 shadow-sm outline-none transition focus:border-black focus:ring-2 focus:ring-slate-200"
                  />
                </div>

                <div>
                  <label className="mb-1.5 block text-sm font-semibold text-slate-700">
                    E-mail
                  </label>
                  <input
                    type="email"
                    placeholder="seu@email.com"
                    value={email}
                    onChange={(e) => { setEmail(e.target.value); setError('') }}
                    className="w-full rounded-lg border border-black bg-white px-4 py-3 text-sm text-slate-800 placeholder-slate-300 shadow-sm outline-none transition focus:border-black focus:ring-2 focus:ring-slate-200"
                  />
                </div>

                <div>
                  <label className="mb-1.5 block text-sm font-semibold text-slate-700">
                    Senha
                  </label>
                  <div className="relative">
                    <button
                      type="button"
                      onClick={() => setShowPassword((v) => !v)}
                      className="absolute right-3.5 top-1/2 -translate-y-1/2 text-slate-400 transition hover:text-[#0F5CC0]"
                      aria-label={showPassword ? 'Ocultar senha' : 'Mostrar senha'}
                    >
                      {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                    </button>
                    <input
                      type={showPassword ? 'text' : 'password'}
                      placeholder="••••••••"
                      value={password}
                      onChange={(e) => { setPassword(e.target.value); setError('') }}
                      className="w-full rounded-lg border border-black bg-white py-3 pl-4 pr-10 text-sm text-slate-800 placeholder-slate-300 shadow-sm outline-none transition focus:border-black focus:ring-2 focus:ring-slate-200"
                    />
                  </div>
                </div>

                <div>
                  <label className="mb-1.5 block text-sm font-semibold text-slate-700">
                    Confirmar senha
                  </label>
                  <input
                    type={showPassword ? 'text' : 'password'}
                    placeholder="••••••••"
                    value={confirmPassword}
                    onChange={(e) => { setConfirmPassword(e.target.value); setError('') }}
                    className="w-full rounded-lg border border-black bg-white px-4 py-3 text-sm text-slate-800 placeholder-slate-300 shadow-sm outline-none transition focus:border-black focus:ring-2 focus:ring-slate-200"
                  />
                </div>

                {error && <p className="text-xs text-red-500">{error}</p>}

                <button
                  type="submit"
                  disabled={loading}
                  className="w-full rounded-lg bg-emerald-500 py-3.5 text-sm font-semibold text-white shadow-md shadow-emerald-200 transition hover:bg-emerald-600 active:scale-[0.98] disabled:opacity-70"
                >
                  {loading ? 'Criando...' : 'Criar conta'}
                </button>
              </form>
            </>
          )}
        </div>
      </div>
    </div>
  )
}
