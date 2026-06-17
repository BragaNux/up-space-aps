package postgres

import (
	"context"
	"database/sql"

	"up-espaco/backend/internal/domain/entities"
	"up-espaco/backend/internal/domain/repositories"
)

// AnnouncementRepository e a implementacao em postgres do repositorio de comunicados
type AnnouncementRepository struct {
	db *DB
}

func NewAnnouncementRepository(db *DB) *AnnouncementRepository {
	return &AnnouncementRepository{db: db}
}

// scanAnnouncement le uma linha do banco pra dentro de um Announcement; withRead controla se a coluna "read" foi pedida na query
func scanAnnouncement(scan func(dest ...any) error, withRead bool) (*entities.Announcement, error) {
	a := &entities.Announcement{}
	var attachment sql.NullString
	var read sql.NullBool

	dest := []any{&a.ID, &a.Title, &a.Sender, &a.Priority, &a.Preview, &a.Body, &attachment, &a.CreatedAt}
	if withRead {
		dest = append(dest, &read)
	}

	if err := scan(dest...); err != nil {
		return nil, err
	}
	if attachment.Valid {
		a.AttachmentName = &attachment.String
	}
	a.Read = read.Valid && read.Bool
	return a, nil
}

// List busca os comunicados; se vier userID, faz join pra marcar quais ja foram lidos por ele
func (r *AnnouncementRepository) List(ctx context.Context, userID *int64) ([]*entities.Announcement, error) {
	var rows *sql.Rows
	var err error

	if userID != nil {
		rows, err = r.db.Conn().QueryContext(ctx, `
			SELECT a.id, a.title, a.sender, a.priority, a.preview, a.body, a.attachment_name, a.created_at,
			       (ar.id IS NOT NULL) AS read
			FROM announcements a
			LEFT JOIN announcement_reads ar ON ar.announcement_id = a.id AND ar.user_id = $1
			ORDER BY a.created_at DESC
		`, *userID)
	} else {
		rows, err = r.db.Conn().QueryContext(ctx, `
			SELECT id, title, sender, priority, preview, body, attachment_name, created_at
			FROM announcements
			ORDER BY created_at DESC
		`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	announcements := make([]*entities.Announcement, 0)
	for rows.Next() {
		a, err := scanAnnouncement(rows.Scan, userID != nil)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, a)
	}
	return announcements, rows.Err()
}

// GetByID busca um comunicado pelo id, marcando se o usuario informado ja leu
func (r *AnnouncementRepository) GetByID(ctx context.Context, id int64, userID *int64) (*entities.Announcement, error) {
	if userID != nil {
		row := r.db.Conn().QueryRowContext(ctx, `
			SELECT a.id, a.title, a.sender, a.priority, a.preview, a.body, a.attachment_name, a.created_at,
			       (ar.id IS NOT NULL) AS read
			FROM announcements a
			LEFT JOIN announcement_reads ar ON ar.announcement_id = a.id AND ar.user_id = $1
			WHERE a.id = $2
		`, *userID, id)
		return scanAnnouncement(row.Scan, true)
	}

	row := r.db.Conn().QueryRowContext(ctx, `
		SELECT id, title, sender, priority, preview, body, attachment_name, created_at
		FROM announcements WHERE id = $1
	`, id)
	return scanAnnouncement(row.Scan, false)
}

// Create insere um comunicado novo
func (r *AnnouncementRepository) Create(ctx context.Context, a *entities.Announcement) error {
	row := r.db.Conn().QueryRowContext(ctx, `
		INSERT INTO announcements (title, sender, priority, preview, body, attachment_name)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`, a.Title, a.Sender, a.Priority, a.Preview, a.Body, a.AttachmentName)
	return row.Scan(&a.ID, &a.CreatedAt)
}

// Update atualiza um comunicado existente
func (r *AnnouncementRepository) Update(ctx context.Context, a *entities.Announcement) error {
	_, err := r.db.Conn().ExecContext(ctx, `
		UPDATE announcements
		SET title = $1, sender = $2, priority = $3, preview = $4, body = $5, attachment_name = $6
		WHERE id = $7
	`, a.Title, a.Sender, a.Priority, a.Preview, a.Body, a.AttachmentName, a.ID)
	return err
}

// Delete apaga o comunicado pelo id
func (r *AnnouncementRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `DELETE FROM announcements WHERE id = $1`, id)
	return err
}

// MarkRead registra que o usuario leu o comunicado (ignora se ja tinha marcado antes)
func (r *AnnouncementRepository) MarkRead(ctx context.Context, announcementID int64, userID int64) error {
	_, err := r.db.Conn().ExecContext(ctx, `
		INSERT INTO announcement_reads (announcement_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (announcement_id, user_id) DO NOTHING
	`, announcementID, userID)
	return err
}

var _ repositories.AnnouncementRepository = (*AnnouncementRepository)(nil)
