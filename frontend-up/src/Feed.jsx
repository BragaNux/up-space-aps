import { useState, useEffect, useCallback, useRef } from 'react'
import {
  Heart,
  MessageCircle,
  Share2,
  MoreHorizontal,
  Bookmark,
  X,
  Send,
  Plus,
  Camera,
  Users,
  Link,
  Mail,
  Phone,
} from 'lucide-react'
import defaultAvatar from '../assets/default.png'
import { getPosts, likePost, unlikePost, bookmarkPost, unbookmarkPost, getComments, addComment, createPost, fileToBase64 } from './api'

const tagStyles = {
  Pedagógico: 'bg-pink-50 text-pink-600',
  Lanche: 'bg-amber-50 text-amber-600',
  Evento: 'bg-indigo-50 text-indigo-600',
}

// transforma uma data em texto relativo (agora, ha 5 min, ontem, ha 3 dias...)
function relativeTime(dateStr) {
  const date = new Date(dateStr)
  const diffMs = Date.now() - date.getTime()
  const minutes = Math.floor(diffMs / 60000)
  if (minutes < 1) return 'agora'
  if (minutes < 60) return `há ${minutes} min`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `há ${hours} h`
  const days = Math.floor(hours / 24)
  if (days === 1) return 'ontem'
  return `há ${days} dias`
}

// etiqueta colorida pro tipo do post (Pedagogico/Lanche/Evento)
function Tag({ tag }) {
  if (!tag) return null
  return (
    <span className={`rounded-full px-3 py-1 text-xs font-semibold ${tagStyles[tag] ?? 'bg-emerald-50 text-emerald-600'}`}>
      {tag}
    </span>
  )
}

// barra de botoes de curtir, comentar, compartilhar e salvar de um post
function PostActions({ post, liked, commentCount, bookmarked, onLike, onComment, onBookmark, onShare }) {
  const likeCount = post.likes
  return (
    <div className="flex items-center gap-5">
      <button
        onClick={onLike}
        className={`flex items-center gap-1.5 text-sm font-medium transition-all duration-150 active:scale-90 ${
          liked ? 'text-rose-500' : 'text-slate-500 hover:text-rose-500'
        }`}
      >
        <Heart className={`h-5 w-5 transition-transform duration-200 ${liked ? 'fill-rose-500 scale-110' : ''}`} />
        {likeCount}
      </button>
      <button
        onClick={onComment}
        className="flex items-center gap-1.5 text-sm font-medium text-slate-500 transition-colors hover:text-emerald-600"
      >
        <MessageCircle className="h-5 w-5" />
        {commentCount}
      </button>
      <button
        onClick={onShare}
        className="flex items-center gap-1.5 text-sm font-medium text-slate-500 transition-colors hover:text-emerald-600"
      >
        <Share2 className="h-5 w-5" />
      </button>
      <button
        onClick={onBookmark}
        className={`ml-auto transition-colors active:scale-90 ${
          bookmarked ? 'text-amber-500' : 'text-slate-400 hover:text-amber-500'
        }`}
      >
        <Bookmark className={`h-5 w-5 ${bookmarked ? 'fill-amber-500' : ''}`} />
      </button>
    </div>
  )
}

// card de um post no feed, com duplo-clique pra curtir tipo Instagram
function Post({ post, liked, commentCount, bookmarked, onLike, onOpen, onBookmark, onShare, index, isOwnPost }) {
  const [showHeartOverlay, setShowHeartOverlay] = useState(false)

  // curte no duplo clique e mostra o coracao animado por um instante
  function handleDoubleClick() {
    if (!liked) {
      onLike()
    }
    setShowHeartOverlay(true)
    setTimeout(() => setShowHeartOverlay(false), 800)
  }

  return (
    <article
      className={`animate-slide-up stagger-${Math.min(index + 1, 6)} overflow-hidden rounded-3xl bg-white shadow-sm ring-1 ring-slate-100 transition-shadow duration-300 hover:shadow-md`}
    >
      <div className="flex items-center gap-3 p-4">
        <img
          src={post.teacher_avatar || defaultAvatar}
          alt={post.teacher_name}
          onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
          className="h-11 w-11 rounded-full object-cover ring-2 ring-white shadow-sm"
        />
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="font-semibold text-slate-800">
              {post.teacher_name || 'Professora'} <span className="font-normal text-slate-400">postou</span> <span className="font-extrabold text-emerald-600">{post.title}</span>
            </span>
            {post.visibility === 'turma' && (
              <span className="inline-flex items-center gap-1 rounded-full bg-indigo-50 px-2 py-0.5 text-[10px] font-semibold text-indigo-600">
                <Users className="h-3 w-3" />
                Mural da Turma
              </span>
            )}
          </div>
          <p className="text-xs text-slate-400">
            {relativeTime(post.created_at)}
            {!isOwnPost && post.student_name && (
              <> · <span className="font-medium text-emerald-600">Atividade de {post.student_name}</span></>
            )}
          </p>
        </div>
        <button className="text-slate-300 transition hover:text-slate-500">
          <MoreHorizontal className="h-5 w-5" />
        </button>
      </div>

      <p className="px-4 pb-3 text-sm leading-relaxed text-slate-600">{post.description}</p>
      {post.image_url && (
        <div className="relative px-4 pb-3 cursor-pointer select-none" onDoubleClick={handleDoubleClick}>
          <img src={post.image_url} alt={post.title} className="max-h-85 w-full rounded-2xl object-cover ring-1 ring-slate-100" />
          {showHeartOverlay && (
            <div className="absolute inset-0 flex items-center justify-center bg-slate-900/10 rounded-2xl mx-4 mb-3" style={{ animation: 'pop-heart 0.8s ease-out forwards' }}>
              <Heart className="h-20 w-20 text-rose-500 fill-rose-500 drop-shadow-md" />
            </div>
          )}
        </div>
      )}
      {post.pedagogical_note && (
        <p className="px-4 pb-3 text-sm italic leading-relaxed text-slate-400">{post.pedagogical_note}</p>
      )}

      <div className="px-4 py-3">
        <PostActions
          post={post}
          liked={liked}
          commentCount={commentCount}
          bookmarked={bookmarked}
          onLike={onLike}
          onComment={onOpen}
          onBookmark={onBookmark}
          onShare={onShare}
        />
      </div>
    </article>
  )
}

// bolha de um comentario, com avatar e nome do autor
function Comment({ comment }) {
  return (
    <div className="flex gap-3">
      <img
        src={comment.avatar_url || defaultAvatar}
        alt={comment.author_name}
        className="h-9 w-9 shrink-0 rounded-full object-cover"
      />
      <div className="min-w-0 flex-1">
        <div className="rounded-2xl rounded-tl-md bg-slate-50 px-3.5 py-2.5">
          <p className="text-sm font-semibold text-slate-800">{comment.author_name}</p>
          <p className="text-sm leading-relaxed text-slate-600">{comment.text}</p>
        </div>
        <div className="mt-1 flex gap-4 pl-1 text-xs text-slate-400">
          <span>{relativeTime(comment.created_at)}</span>
        </div>
      </div>
    </div>
  )
}

// modal com o post aberto, lista de comentarios e campo pra comentar
function PostDetail({ post, liked, comments, onLike, onClose, onAddComment, isOwnPost, userAvatar, bookmarked, onBookmark, onShare }) {
  const [draft, setDraft] = useState('')
  const [sending, setSending] = useState(false)
  const [error, setError] = useState('')

  // envia o comentario digitado
  async function submit(e) {
    e.preventDefault()
    const text = draft.trim()
    if (!text) return
    setSending(true)
    setError('')
    try {
      await onAddComment(post.id, text)
      setDraft('')
    } catch (err) {
      setError(err.message || 'Faça login para comentar.')
    } finally {
      setSending(false)
    }
  }

  return (
    <div className="animate-fade-in fixed inset-0 z-50 flex items-end justify-center bg-slate-900/50 p-0 backdrop-blur-sm sm:items-center sm:p-4">
      <div className="animate-slide-up flex max-h-[92vh] w-full max-w-lg flex-col overflow-hidden rounded-t-3xl bg-white shadow-2xl sm:rounded-3xl">
        <div className="flex items-center gap-3 border-b border-slate-100 p-4">
          <img
            src={post.teacher_avatar || defaultAvatar}
            alt={post.teacher_name}
            onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
            className="h-10 w-10 rounded-full object-cover"
          />
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2 flex-wrap">
              <span className="font-semibold text-slate-800">
                {post.teacher_name || 'Professora'} <span className="font-normal text-slate-400">postou</span> <span className="font-extrabold text-emerald-600">{post.title}</span>
              </span>
              {post.visibility === 'turma' && (
                <span className="inline-flex items-center gap-1 rounded-full bg-indigo-50 px-2 py-0.5 text-[10px] font-semibold text-indigo-600">
                  <Users className="h-3 w-3" />
                  Mural da Turma
                </span>
              )}
            </div>
            <p className="text-xs text-slate-400">
              {relativeTime(post.created_at)}
              {!isOwnPost && post.student_name && (
                <> · <span className="font-medium text-emerald-600">Atividade de {post.student_name}</span></>
              )}
            </p>
          </div>
          <button
            onClick={onClose}
            aria-label="Fechar"
            className="grid h-9 w-9 place-items-center rounded-full text-slate-400 transition hover:bg-slate-100 hover:text-slate-600"
          >
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto">
          <p className="px-4 pb-3 pt-4 text-sm leading-relaxed text-slate-600">{post.description}</p>
          {post.image_url && (
            <div className="px-4 pb-3">
              <img src={post.image_url} alt={post.title} className="max-h-85 w-full rounded-2xl object-cover ring-1 ring-slate-100" />
            </div>
          )}

          <div className="border-b border-slate-100 px-4 py-3">
            <PostActions
              post={post}
              liked={liked}
              commentCount={comments.length}
              bookmarked={bookmarked}
              onLike={onLike}
              onComment={() => {}}
              onBookmark={onBookmark}
              onShare={onShare}
            />
          </div>

          <div className="space-y-4 p-4">
            <p className="text-sm font-semibold text-slate-700">Comentários ({comments.length})</p>
            {comments.length === 0 ? (
              <p className="text-sm text-slate-400">Seja o primeiro a comentar.</p>
            ) : (
              comments.map((c) => <Comment key={c.id} comment={c} />)
            )}
          </div>
        </div>

        <form onSubmit={submit} className="flex items-center gap-2 border-t border-slate-100 p-3">
          <img
            src={userAvatar || defaultAvatar}
            alt="Você"
            onError={(e) => { e.currentTarget.onerror = null; e.currentTarget.src = defaultAvatar }}
            className="h-9 w-9 shrink-0 rounded-full object-cover"
          />
          <input
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            placeholder="Escreva um comentário..."
            className="flex-1 rounded-full border border-slate-200 bg-slate-50 px-4 py-2.5 text-sm outline-none transition focus:border-emerald-400 focus:bg-white focus:ring-2 focus:ring-emerald-100"
          />
          <button
            type="submit"
            disabled={!draft.trim() || sending}
            aria-label="Enviar comentário"
            className="grid h-10 w-10 shrink-0 place-items-center rounded-full bg-emerald-500 text-white transition hover:bg-emerald-600 active:scale-95 disabled:opacity-40"
          >
            <Send className="h-4 w-4" />
          </button>
        </form>
        {error && <p className="px-4 pb-3 text-xs text-red-500">{error}</p>}
      </div>
    </div>
  )
}

// formulario de criar um post novo (titulo, descricao, nota pedagogica e foto opcional)
function NewPostForm({ onCreate, onClose }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [pedagogicalNote, setPedagogicalNote] = useState('')
  const [imageUrl, setImageUrl] = useState('')
  const [uploading, setUploading] = useState(false)
  const [photoError, setPhotoError] = useState('')
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')
  const [shareWithClass, setShareWithClass] = useState(false)
  const fileInputRef = useRef(null)

  const MAX_PHOTO_BYTES = 4 * 1024 * 1024

  // valida tipo e tamanho da imagem escolhida e converte pra base64 antes de mostrar a preview
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
      setImageUrl(base64)
    } catch (err) {
      setPhotoError(err.message || 'Não foi possível ler a imagem.')
    } finally {
      setUploading(false)
      e.target.value = ''
    }
  }

  // valida os campos obrigatorios e publica o post
  async function submit(e) {
    e.preventDefault()
    if (!title.trim() || !description.trim() || !pedagogicalNote.trim()) {
      setError('Preencha todos os campos.')
      return
    }
    setSaving(true)
    setError('')
    try {
      await onCreate({
        title,
        description,
        pedagogical_note: pedagogicalNote,
        image_url: imageUrl,
        visibility: shareWithClass ? 'turma' : 'private',
      })
      setTitle('')
      setDescription('')
      setPedagogicalNote('')
      setImageUrl('')
      setShareWithClass(false)
      onClose()
    } catch (err) {
      setError(err.message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <form onSubmit={submit} className="animate-slide-up space-y-2 rounded-3xl bg-white p-5 shadow-sm ring-1 ring-slate-100">
      <h3 className="font-bold text-slate-800">Nova atividade</h3>
      <input
        value={title}
        onChange={(e) => setTitle(e.target.value)}
        placeholder="Título"
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />
      <textarea
        value={description}
        onChange={(e) => setDescription(e.target.value)}
        placeholder="Descrição da atividade"
        rows={2}
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />
      <textarea
        value={pedagogicalNote}
        onChange={(e) => setPedagogicalNote(e.target.value)}
        placeholder="Nota pedagógica"
        rows={2}
        className="w-full rounded-lg border border-slate-200 px-3 py-2 text-sm outline-none focus:border-emerald-400"
      />

      {imageUrl && (
        <div className="relative mt-2 inline-block">
          <img src={imageUrl} alt="Preview do post" className="h-24 w-32 rounded-xl object-cover ring-1 ring-slate-100" />
          <button
            type="button"
            onClick={() => setImageUrl('')}
            className="absolute -right-2 -top-2 grid h-6 w-6 place-items-center rounded-full bg-slate-900 text-white hover:bg-slate-700 active:scale-90"
          >
            <X className="h-3.5 w-3.5" />
          </button>
        </div>
      )}

      {photoError && <p className="text-xs text-rose-500">{photoError}</p>}

      <label className="flex items-center gap-2 text-sm text-slate-600 pb-2 pt-1 cursor-pointer select-none">
        <input
          type="checkbox"
          checked={shareWithClass}
          onChange={(e) => setShareWithClass(e.target.checked)}
          className="rounded border-slate-300 text-emerald-600 focus:ring-emerald-500 cursor-pointer"
        />
        Compartilhar com a turma (visível para outros pais)
      </label>

      <div className="flex justify-between items-center pt-2">
        <div>
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            disabled={uploading}
            className="flex items-center gap-1.5 rounded-full border border-slate-200 px-4 py-2 text-xs font-semibold text-slate-600 transition hover:bg-slate-50 disabled:opacity-50"
          >
            {uploading ? (
              <span className="h-3.5 w-3.5 animate-spin rounded-full border-2 border-slate-300 border-t-slate-600" />
            ) : (
              <Camera className="h-3.5 w-3.5 text-emerald-600" />
            )}
            Adicionar foto
          </button>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handlePhotoChange}
            className="hidden"
          />
        </div>
        <div className="flex gap-2">
          <button type="button" onClick={onClose} className="rounded-full px-3 py-1.5 text-sm text-slate-500 hover:bg-slate-100">
            Cancelar
          </button>
          <button
            type="submit"
            disabled={saving || uploading}
            className="rounded-full bg-emerald-500 px-4 py-1.5 text-sm font-semibold text-white hover:bg-emerald-600 disabled:opacity-60"
          >
            {saving ? 'Publicando...' : 'Publicar'}
          </button>
        </div>
      </div>
    </form>
  )
}

// tela do feed de atividades de uma crianca: lista posts, curtidas, comentarios e compartilhar
export default function Feed({ studentId, canEdit, currentUser }) {
  const [posts, setPosts] = useState([])
  const [likedIds, setLikedIds] = useState(() => new Set())
  const [bookmarkedIds, setBookmarkedIds] = useState(() => new Set())
  const [sharingPost, setSharingPost] = useState(null)
  const [openId, setOpenId] = useState(null)
  const [comments, setComments] = useState({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)

  useEffect(() => {
    if (!studentId) return
    setLoading(true)
    getPosts(studentId)
      .then(setPosts)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false))
  }, [studentId])

  // cria o post e recarrega a lista pra mostrar ele logo
  async function handleCreatePost(data) {
    await createPost(studentId, data)
    const updated = await getPosts(studentId)
    setPosts(updated)
  }

  // busca os comentarios de um post e guarda no estado por id
  const loadComments = useCallback((postId) => {
    getComments(postId)
      .then((data) => setComments((prev) => ({ ...prev, [postId]: data })))
      .catch(() => {})
  }, [])

  // curte ou descurte o post, atualizando contagem e o set de ids curtidos
  async function toggleLike(id) {
    const isLiked = likedIds.has(id)
    try {
      if (isLiked) {
        const { likes } = await unlikePost(id)
        setPosts((prev) => prev.map((p) => (p.id === id ? { ...p, likes } : p)))
        setLikedIds((prev) => {
          const next = new Set(prev)
          next.delete(id)
          return next
        })
      } else {
        const { likes } = await likePost(id)
        setPosts((prev) => prev.map((p) => (p.id === id ? { ...p, likes } : p)))
        setLikedIds((prev) => {
          const next = new Set(prev)
          next.add(id)
          return next
        })
      }
    } catch (err) {
      setError(err.message)
    }
  }

  // salva ou remove o post dos favoritos, atualizando contagem e o set de ids salvos
  async function toggleBookmark(id) {
    const isBookmarked = bookmarkedIds.has(id)
    try {
      if (isBookmarked) {
        const { bookmarks } = await unbookmarkPost(id)
        setPosts((prev) => prev.map((p) => (p.id === id ? { ...p, bookmarks } : p)))
        setBookmarkedIds((prev) => {
          const next = new Set(prev)
          next.delete(id)
          return next
        })
      } else {
        const { bookmarks } = await bookmarkPost(id)
        setPosts((prev) => prev.map((p) => (p.id === id ? { ...p, bookmarks } : p)))
        setBookmarkedIds((prev) => {
          const next = new Set(prev)
          next.add(id)
          return next
        })
      }
    } catch (err) {
      setError(err.message)
    }
  }

  // abre o modal de detalhe do post, carregando os comentarios se ainda nao tiver buscado
  function openPost(id) {
    setOpenId(id)
    if (!comments[id]) loadComments(id)
  }

  // adiciona o comentario e recarrega a lista de comentarios daquele post
  async function handleAddComment(postId, text) {
    await addComment(postId, text)
    loadComments(postId)
  }

  const activePost = posts.find((p) => p.id === openId) ?? null

  return (
    <div className="mx-auto max-w-xl space-y-6">
      <div className="flex items-center justify-between rounded-3xl bg-white p-6 shadow-sm ring-1 ring-slate-100">
        <div>
          <h2 className="text-xl font-extrabold text-slate-800">Feed de Atividades</h2>
          <p className="text-sm text-slate-400">Acompanhe o dia a dia em tempo real</p>
        </div>
        {canEdit && !showForm && (
          <button
            onClick={() => setShowForm(true)}
            className="animate-pop-in flex items-center gap-1.5 rounded-full bg-emerald-500 px-4 py-2 text-sm font-semibold text-white shadow-md shadow-emerald-200 transition-all duration-200 hover:bg-emerald-600 hover:shadow-lg active:scale-95"
          >
            <Plus className="h-4 w-4" />
            Nova
          </button>
        )}
      </div>

      {showForm && <NewPostForm onCreate={handleCreatePost} onClose={() => setShowForm(false)} />}

      {error && <p className="rounded-2xl bg-rose-50 p-4 text-sm text-rose-600">{error}</p>}
      {loading && <p className="text-sm text-slate-400">Carregando...</p>}
      {!loading && posts.length === 0 && !error && (
        <p className="rounded-2xl bg-white p-6 text-center text-sm text-slate-400 ring-1 ring-slate-100">
          Nenhuma atividade publicada ainda.
        </p>
      )}

      {posts.map((post, index) => (
        <Post
          key={post.id}
          post={post}
          index={index}
          liked={likedIds.has(post.id)}
          commentCount={comments[post.id]?.length ?? post.comment_count ?? 0}
          bookmarked={bookmarkedIds.has(post.id)}
          onLike={() => toggleLike(post.id)}
          onOpen={() => openPost(post.id)}
          onBookmark={() => toggleBookmark(post.id)}
          onShare={() => setSharingPost(post)}
          isOwnPost={post.student_id === studentId}
        />
      ))}

      {activePost && (
        <PostDetail
          post={activePost}
          liked={likedIds.has(activePost.id)}
          comments={comments[activePost.id] ?? []}
          bookmarked={bookmarkedIds.has(activePost.id)}
          onLike={() => toggleLike(activePost.id)}
          onClose={() => setOpenId(null)}
          onAddComment={handleAddComment}
          onBookmark={() => toggleBookmark(activePost.id)}
          onShare={() => setSharingPost(activePost)}
          isOwnPost={activePost.student_id === studentId}
          userAvatar={currentUser?.avatar_url}
        />
      )}

      {sharingPost && (
        <div className="animate-fade-in fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 p-4 backdrop-blur-sm">
          <div className="animate-pop-in w-full max-w-sm rounded-3xl bg-white p-6 shadow-2xl">
            <div className="flex items-center justify-between border-b border-slate-100 pb-3">
              <h3 className="text-base font-bold text-slate-800">Compartilhar Atividade</h3>
              <button
                onClick={() => setSharingPost(null)}
                className="rounded-full p-1 text-slate-400 hover:bg-slate-100 hover:text-slate-600"
              >
                <X className="h-5 w-5" />
              </button>
            </div>
            <p className="mt-3 text-sm text-slate-500">
              Escolha onde deseja compartilhar a atividade <strong>{sharingPost.title}</strong>:
            </p>
            <div className="mt-4 grid grid-cols-2 gap-3">
              <button
                onClick={() => { alert('Link copiado com sucesso!'); setSharingPost(null) }}
                className="flex flex-col items-center gap-2 rounded-2xl border border-slate-100 bg-slate-50 p-4 transition hover:bg-slate-100"
              >
                <span className="grid h-10 w-10 place-items-center rounded-full bg-emerald-50 text-emerald-600">
                  <Link className="h-5 w-5" />
                </span>
                <span className="text-xs font-semibold text-slate-700">Copiar Link</span>
              </button>
              <button
                onClick={() => { alert('Compartilhando no WhatsApp...'); setSharingPost(null) }}
                className="flex flex-col items-center gap-2 rounded-2xl border border-slate-100 bg-slate-50 p-4 transition hover:bg-slate-100"
              >
                <span className="grid h-10 w-10 place-items-center rounded-full bg-green-50 text-green-600 animate-pulse">
                  <Phone className="h-5 w-5" />
                </span>
                <span className="text-xs font-semibold text-slate-700">WhatsApp</span>
              </button>
              <button
                onClick={() => { alert('Compartilhando no Telegram...'); setSharingPost(null) }}
                className="flex flex-col items-center gap-2 rounded-2xl border border-slate-100 bg-slate-50 p-4 transition hover:bg-slate-100"
              >
                <span className="grid h-10 w-10 place-items-center rounded-full bg-sky-50 text-sky-600">
                  <Send className="h-5 w-5" />
                </span>
                <span className="text-xs font-semibold text-slate-700">Telegram</span>
              </button>
              <button
                onClick={() => { alert('Enviando por E-mail...'); setSharingPost(null) }}
                className="flex flex-col items-center gap-2 rounded-2xl border border-slate-100 bg-slate-50 p-4 transition hover:bg-slate-100"
              >
                <span className="grid h-10 w-10 place-items-center rounded-full bg-indigo-50 text-indigo-600">
                  <Mail className="h-5 w-5" />
                </span>
                <span className="text-xs font-semibold text-slate-700">E-mail</span>
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
