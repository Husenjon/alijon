package inkassback

type BranchData struct {
	Данные []Branch `json:"данные" binding:"required"`
	Ошибка string   `json:"ошибка" binding:"required"`
}

type Branch struct {
	Id           int    `db:"id"`
	ИНН          string `json:"ИНН" db:"inn" binding:"required"`
	Код          string `json:"Код" db:"code" binding:"required"`
	Наименование string `json:"наименование" db:"name" binding:"required"`
}

type Bank struct {
	Код          string `json:"Код" db:"code" binding:"required"`
	Наименование string `json:"наименование" db:"name" binding:"required"`
}

type Contract struct {
	Номер            string `json:"номер" db:"number" binding:"required"`
	Дата             string `json:"дата" db:"date" binding:"required"`
	Учетная_карточка string `json:"учетная_карточка" db:"registration_card" binding:"required"`
}

type Client struct {
	ИНН          string `json:"ИНН" db:"inn" binding:"required"`
	Наименование string `json:"наименование" db:"name" binding:"required"`
}
type Bags struct {
	Мешки []string `json:"мешки" db:"name"`
}
type Contracts struct {
	Id                  int      `db:"id"`
	Дополнительный_банк Bank     `json:"дополнительный_банк" db:"additional_bank"`
	Долг                float32  `json:"долг" db:"debt"`
	Маршрут             string   `json:"маршрут" db:"route"`
	Тариф               string   `json:"тариф" db:"rate"`
	Номера_мешков       string   `json:"номера_мешков" db:"bag_numbers"`
	Банк                Bank     `json:"банк" db:"bank_id"`
	Контракт            Contract `json:"контракт" db:"contract"`
	Состояние           string   `json:"Состояние" db:"state"`
	Клиент              Client   `json:"клиент" db:"client"`
	// Мешки               Bags     `json:"мешки" db:"name"`
}
type ContractsBody struct {
	Ошибка       string      `json:"ошибка"`
	ИНН          string      `json:"ИНН"`
	Наименование string      `json:"наименование"`
	Данные       []Contracts `json:"данные"`
}
type Revenue struct {
	Клиент   Client   `json:"клиент" db:"client_id"`
	Сумма    float64  `json:"сумма"`
	Контракт Contract `json:"контракт" db:"contract_id"`
	Дата     string   `json:"дата"`
}
type RevenueBody struct {
	Ошибка       string    `json:"ошибка"`
	ИНН          string    `json:"ИНН"`
	Наименование string    `json:"наименование"`
	Данные       []Revenue `json:"данные"`
}
