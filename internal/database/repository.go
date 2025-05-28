package database

import (
	"bitespeed-identity-reconciliation/internal/models"
	"database/sql"
	"time"
)

type ContactRepository struct {
	db *sql.DB
}

func NewContactRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

func (r *ContactRepository) FindByEmailOrPhone(email, phoneNumber *string) ([]models.Contact, error) {
	query := `
		SELECT id, phone_number, email, linked_id, link_precedence, created_at, updated_at, deleted_at
        FROM contacts 
        WHERE deleted_at IS NULL 
        AND (email = ? OR phone_number = ?)
        ORDER BY created_at ASC
	`
	
	rows, err := r.db.Query(query, email, phoneNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.Contact
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
			&contact.ID,
			&contact.PhoneNumber,
			&contact.Email,
            &contact.LinkedID,
            &contact.LinkPrecedence,
            &contact.CreatedAt,
            &contact.UpdatedAt,
            &contact.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (r *ContactRepository) FindByLinkedID(linkedID int) ([]models.Contact, error) {
	query := `
        SELECT id, phone_number, email, linked_id, link_precedence, created_at, updated_at, deleted_at
        FROM contacts 
        WHERE deleted_at IS NULL AND linked_id = ?
        ORDER BY created_at ASC
    `

	rows, err := r.db.Query(query, linkedID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []models.Contact
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
			&contact.ID,
            &contact.PhoneNumber,
            &contact.Email,
            &contact.LinkedID,
            &contact.LinkPrecedence,
            &contact.CreatedAt,
            &contact.UpdatedAt,
            &contact.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func (r *ContactRepository) Create(contact *models.Contact) error {
	query := `
        INSERT INTO contacts (phone_number, email, linked_id, link_precedence, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    
    now := time.Now()
    contact.CreatedAt = now
    contact.UpdatedAt = now

    result, err := r.db.Exec(query, 
        contact.PhoneNumber, 
        contact.Email, 
        contact.LinkedID, 
        contact.LinkPrecedence,
        contact.CreatedAt,
        contact.UpdatedAt,
    )
    if err != nil {
        return err
    }

	id, err := result.LastInsertId()
    if err != nil {
        return err
    }

    contact.ID = int(id)
    return nil
}

func (r *ContactRepository) UpdateLinkPrecedence(id int, linkedID int, linkPrecedence string) error {
    query := `
        UPDATE contacts 
        SET linked_id = ?, link_precedence = ?, updated_at = ?
        WHERE id = ? AND deleted_at IS NULL
    `
    
    _, err := r.db.Exec(query, linkedID, linkPrecedence, time.Now(), id)
    return err
}

func (r *ContactRepository) FindByID(id int) (*models.Contact, error) {
    query := `
        SELECT id, phone_number, email, linked_id, link_precedence, created_at, updated_at, deleted_at
        FROM contacts 
        WHERE id = ? AND deleted_at IS NULL
    `
    
    var contact models.Contact
    err := r.db.QueryRow(query, id).Scan(
        &contact.ID,
        &contact.PhoneNumber,
        &contact.Email,
        &contact.LinkedID,
        &contact.LinkPrecedence,
        &contact.CreatedAt,
        &contact.UpdatedAt,
        &contact.DeletedAt,
    )

	if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, err
    }

    return &contact, nil
}