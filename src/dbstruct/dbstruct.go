package dbstruct

type Row struct {
	PoductsID     int    `db:"products_id"`
	OffersID      int    `db:"offers_id"`
	OffersTitle   string `db:"offers_title"`
	ProductsTitle string `db:"products_title"`
}
