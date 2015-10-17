package dbstruct

var SqlLine string = `
SELECT 
	products.id as products_id, products.title as products_title,
	COALESCE(offers.id, 0) as offers_id, COALESCE(offers.title, '') as offers_title 
FROM products LEFT JOIN offers ON offers.products_id = products.id
WHERE products.id > $1 OR ( products.id = $2 AND offers.id > $3 )
ORDER BY products.id, offers.id
limit $4
`

type Row struct {
	PoductsID     int    `db:"products_id"`
	OffersID      int    `db:"offers_id"`
	OffersTitle   string `db:"offers_title"`
	ProductsTitle string `db:"products_title"`
}
