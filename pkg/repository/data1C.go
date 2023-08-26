package repository

import (
	"encoding/json"
	"fmt"

	inkassback "github.com/Husenjon/InkassBack"
	"github.com/jmoiron/sqlx"
)

type Data1CPostgres struct {
	db *sqlx.DB
}

func NewData1CPostgres(db *sqlx.DB) *Data1CPostgres {
	return &Data1CPostgres{
		db: db,
	}
}
func (r *Data1CPostgres) CreateBranch(br inkassback.Branch) (inkassback.Branch, error) {
	var branch inkassback.Branch
	query := fmt.Sprintf(
		`insert
			into
			%s (
				inn,
				code,
				name
			)
		values(
			$1,
			$2,
			$3
		) on conflict(code) do
		update
		set
			inn = EXCLUDED.inn,
			name = EXCLUDED.name,
			code = EXCLUDED.code returning *;`,
		branchTable,
	)
	err := r.db.Get(&branch, query, br.ИНН, br.Код, br.Наименование)
	return branch, err
}
func (r *Data1CPostgres) CreateContracts(cs inkassback.Contracts, branchId int) (inkassback.Contracts, error) {
	var contracts inkassback.Contracts
	var id int
	var additional_bank_id int
	var bank_id int
	var contract_id int
	var client_id int

	tx, err := r.db.Begin()
	if err != nil {
		return contracts, err
	}
	createAB := fmt.Sprintf(
		`insert into
			%s (
				code,
				name
			)
		values
		(
			$1,
			$2
		) on conflict(code) do
		update set	name = EXCLUDED.name returning id;`,
		additionalBankTable)
	row := tx.QueryRow(createAB, cs.Дополнительный_банк.Код, cs.Дополнительный_банк.Наименование)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return contracts, err
	}
	additional_bank_id = id

	createB := fmt.Sprintf(
		`insert into
			%s (
				code,
				name
			)
		values
		(
			$1,
			$2
		) on conflict(code) do
		update set	name = EXCLUDED.name returning id;`,
		bankTable)
	row = tx.QueryRow(createB, cs.Банк.Код, cs.Банк.Наименование)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return contracts, err
	}
	bank_id = id

	createContract := fmt.Sprintf(
		`insert into
			%s (
				number,
				registration_card,
				date
			)
		values
		(
			$1,
			$2,
			$3
		) on conflict(number) do
		update
			set
				registration_card = EXCLUDED.registration_card,
				date = EXCLUDED.date
			returning id;`,
		contractTable)
	row = tx.QueryRow(createContract, cs.Контракт.Номер, cs.Контракт.Учетная_карточка, cs.Контракт.Дата)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return contracts, err
	}
	contract_id = id

	createClient := fmt.Sprintf(
		`INSERT INTO
			%s (
				inn,
				name
			)
		VALUES
			(
				$1,
				$2
			) on conflict(inn) do
		update
			set
			name = EXCLUDED.name
			returning id;`,
		clientTable)
	row = tx.QueryRow(createClient, cs.Клиент.ИНН, cs.Клиент.Наименование)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return contracts, err
	}
	client_id = id

	createContracts := fmt.Sprintf(
		`WITH cons as(
			INSERT INTO
				%s (
					contract_number,
					additional_bank_id,
					debt,
					route,
					rate,
					bag_numbers,
					bank_id,
					contract_id,
					state,
					client_id,
					branch_id
				)
			VALUES
				(
					$1,
					$2,
					$3,
					$4,
					$5,
					$6,
					$7,
					$8,
					$9,
					$10,
					$11
				) on conflict(contract_number) do
			update
			set
				additional_bank_id = EXCLUDED.additional_bank_id,
				debt = EXCLUDED.debt,
				route = EXCLUDED.route,
				rate = EXCLUDED.rate,
				bag_numbers = EXCLUDED.bag_numbers,
				bank_id = EXCLUDED.bank_id,
				contract_id = EXCLUDED.contract_id,
				state = EXCLUDED.state,
				client_id = EXCLUDED.client_id,
				branch_id = EXCLUDED.branch_id returning *
		)
		select
			JSON_BUILD_OBJECT(
				'id', cons.ID,
				'contract_number', cons.contract_number,
				'дополнительный_банк', JSON_BUILD_OBJECT(
					'id', additional_bank.id,
					'Код', additional_bank.code,
					'наименование', additional_bank.name
				),
				'долг', cons.debt,
				'маршрут', cons.route,
				'тариф', cons.rate,
				'номера_мешков', cons.bag_numbers,
				'банк', JSON_BUILD_OBJECT(
					'id', bank.id,
					'Код', bank.code,
					'наименование', bank.name
				),
				'контракт', JSON_BUILD_OBJECT(
					'id', contract.id,
					'номер', contract.number,
					'учетная_карточка', contract.registration_card,
					'дата', contract.date
				),
				'Состояние', cons.state, 
				'клиент', JSON_BUILD_OBJECT(
					'id', client.id,
					'ИНН', client.inn,
					'наименование', client.name
				),
				'branch', JSON_BUILD_OBJECT(
					'id', branch.id,
					'ИНН', branch.inn,
					'Код', branch.code,
					'наименование', branch.name
				)
			) as contracts
		from
			cons
			inner join additional_bank on cons.additional_bank_id = additional_bank.id
			inner join bank on bank.id = cons.bank_id
			inner join contract on contract.id = cons.contract_id
			inner join client on client.id = cons.client_id
			inner join branch on branch.id = cons.branch_id`,
		contractsTable)
	row = tx.QueryRow(
		createContracts,
		cs.Контракт.Номер,
		additional_bank_id,
		cs.Долг,
		cs.Маршрут,
		cs.Тариф,
		cs.Номера_мешков,
		bank_id,
		contract_id,
		cs.Состояние,
		client_id,
		branchId,
	)
	var j []byte
	if err := row.Scan(&j); err != nil {
		tx.Rollback()
		return contracts, err
	}
	tx.Commit()
	err = json.Unmarshal(j, &contracts)
	if err != nil {
		return contracts, err
	}
	return contracts, nil
}
func (r *Data1CPostgres) CreateRevenue(revenue inkassback.Revenue) (inkassback.Revenue, error) {
	var re inkassback.Revenue
	var id int
	// var additional_bank_id int
	// var bank_id int
	var contract_id int
	var client_id int
	tx, err := r.db.Begin()
	if err != nil {
		return re, err
	}

	createClient := fmt.Sprintf(
		`INSERT INTO
			%s (
				inn, 
				name
			)
		VALUES
			(
				$1, 
				$2
			) on conflict(inn) do
		update 
			set	
			name = EXCLUDED.name
			returning id;`,
		clientTable)
	row := tx.QueryRow(createClient, revenue.Клиент.ИНН, revenue.Клиент.Наименование)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return re, err
	}
	client_id = id

	createContract := fmt.Sprintf(
		`insert into
			%s (
				number, 
				registration_card, 
				date
			)
		values
		(
			$1, 
			$2, 
			$3
		) on conflict(number) do
		update 
			set	
				registration_card = EXCLUDED.registration_card,
				date = EXCLUDED.date 
			returning id;`,
		contractTable)
	row = tx.QueryRow(createContract, revenue.Контракт.Номер, revenue.Контракт.Учетная_карточка, revenue.Контракт.Дата)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return re, err
	}
	contract_id = id
	createRevenue := fmt.Sprintf(
		`with rev as(
			insert into
				%s (
					client_id,
					sum,
					contract_id,
					date
				)
			values
				(
					$1, 
					$2, 
					$3, 
					$4
				) 
			returning *
		)
		select
			JSON_BUILD_OBJECT(
				'id', rev.ID,
				'клиент', JSON_BUILD_OBJECT(
					'id', client.id,
					'ИНН', client.inn,
					'наименование', client.name
				),
				'сумма', rev.sum,
				'контракт', JSON_BUILD_OBJECT(
					'id', contract.id,
					'номер', contract.number,
					'учетная_карточка', contract.registration_card,
					'дата', contract.date
				),
				'дата', rev.date
			) as revenue
		from
			rev
			inner join client on client.id = rev.client_id
			inner join contract on contract.id = rev.contract_id;`,
		revenueTable)
	row = tx.QueryRow(createRevenue, client_id, revenue.Сумма, contract_id, revenue.Дата)
	var j []byte
	if err := row.Scan(&j); err != nil {
		tx.Rollback()
		return re, err
	}
	tx.Commit()
	err = json.Unmarshal(j, &re)
	if err != nil {
		return re, err
	}
	return re, nil
}
