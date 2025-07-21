package storage

import (
	"context"
	//"time"

	"marketplace/internal/models"
)

func CreateAd(ad *models.Ad) error {
	query := `
		INSERT INTO ads (title, text, image_url, price, author_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	return db.QueryRow(
		context.Background(),
		query,
		ad.Title, ad.Text, ad.ImageURL, ad.Price, ad.AuthorID,
	).Scan(&ad.ID, &ad.CreatedAt)
}

func GetAds(page, limit int, orderBy string, minPrice, maxPrice float64) ([]struct {
	models.Ad
	Author string `json:"author"`
}, error) {
	offset := (page - 1) * limit
	query := `
		SELECT 
			a.id, a.title, a.text, a.image_url, a.price, a.author_id, a.created_at,
			u.login as author
		FROM ads a
		JOIN users u ON a.author_id = u.id
		WHERE ($1 = 0 OR a.price >= $1)
		AND ($2 = 0 OR a.price <= $2)
		ORDER BY ` + orderBy + `
		LIMIT $3 OFFSET $4
	`

	rows, err := db.Query(context.Background(), query, minPrice, maxPrice, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ads []struct {
		models.Ad
		Author string `json:"author"`
	}

	for rows.Next() {
		var ad struct {
			models.Ad
			Author string `json:"author"`
		}
		err := rows.Scan(
			&ad.ID, &ad.Title, &ad.Text, &ad.ImageURL, &ad.Price,
			&ad.AuthorID, &ad.CreatedAt, &ad.Author,
		)
		if err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}

	return ads, nil
}
