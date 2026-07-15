package pgrepo

const createUserQuery = `
	insert into users (
		id, username, password
	) values ($1, $2, $3);
`

const getUserByUsernameQuery = `
	select 
		id, username, password, 
		created_at, updated_at
	from users
	where username = $1;
`

const createCryptoPositionQuery = `
	insert into crypto_positions (
		id, user_id, ticker, coin_id,
		amount, avg_price_usd_cents
	) values ($1, $2, $3, $4, $5);
`

const deleteCryptoPositionQuery = `
	delete from crypto_positions
	where id = $1;
`

const updateCryptoPositionQuery = `
	update crypto_positions
	set amount = $2, avg_price_usd_cents = $3
	where id = $1;
`

const getUserCryptoPositionsQuery = `
	select
		id, user_id, ticker, coin_id,
		amount, avg_price_usd_cents
	from crypto_positions
	where user_id = $1;
`

const getCryptoPositionQuery = `
	select
		id, user_id, ticker,
		amount, avg_price_usd_cents
	from crypto_positions
	where id = $1;
`

const createBankDepositQuery = `
	insert into bank_deposits (
		id, user_id, name,
		currency, amount_cents,
		interest_rate, opened_at, closed_at
	)
	values ($1, $2, $3, $4, $5, $6, $7, $8)
`

const deleteBankDepositQuery = `
	delete from bank_deposits
	where id = $1
`

const getBankDepositByIDQuery = `
	select
		id, user_id, name,
		currency, amount_cents, interest_rate,
		opened_at, closed_at
	from bank_deposits
	where id = $1
`

const getUserBankDepositsQuery = `
	select
		id, user_id,
		name, currency,
		amount_cents, interest_rate,
		opened_at, closed_at
	from bank_deposits
	where user_id = $1
	order by opened_at desc
`

const createRealEstateQuery = `
	insert into real_estates (
		id, user_id,
		name, currency,
		purchase_price_cents,
		monthly_income_cents,
		purchased
	)
	values ($1, $2, $3, $4, $5, $6, $7)
`

const getUserRealEstatesQuery = `
	select
		id, user_id,
		name, currency,
		purchase_price_cents,
		monthly_income_cents,
		purchased
	from real_estates
	where user_id = $1
	order by purchased desc
`

const deleteRealEstateQuery = `
	delete from real_estates
	where id = $1
`

const updateRealEstateQuery = `
	update real_estates
	set
		name = $2, currency = $3,
		purchase_price_cents = $4,
		monthly_income_cents = $5,
		purchased = $6
	where id = $1
`

const getRealEstateQuery = `
	select
		id, user_id, 
		name, currency, 
		purchase_price_cents,
		monthly_income_cents, 
		purchased
	from real_estates
	where id = $1
`

const createStockQuery = `
	insert into stocks (
		id, user_id, ticker,
		amount, avg_price_cents,
		currency
	) values ($1, $2, $3, $4, $5, $6);
`

const updateStockQuery = `
	update stocks
	set 
		user_id = $2, ticker = $3,
		amount = $4, avg_price_cents = $5,
		currency = $6
	where id = $1;
`

const getStockQuery = `
	select 
		id, user_id, ticker,
		amount, avg_price_cents,
		currency
	from stocks
	where id = $1;
`

const deleteStockQuery = `
	delete from stocks
	where id = $1;
`

const getUserStocksQuery = `
	select 
		id, user_id, ticker,
		amount, avg_price_cents,
		currency
	from stocks
	where user_id = $1;
`

const createBondQuery = `
	insert into bonds (
		id, user_id, ticker,
		currency, amount,
		avg_price_cents, coupon_cents,
		coupon_period_months
	) values ($1, $2, $3, $4, $5, $6, $7, $8);
`

const getBondQuery = `
	select 
		id, user_id, ticker,
		currency, amount,
		avg_price_cents, coupon_cents,
		coupon_period_months
	from bonds
	where id = $1;
`

const getUserBondsQuery = `
	select 
		id, user_id, ticker,
		currency, amount,
		avg_price_cents, coupon_cents,
		coupon_period_months
	from bonds
	where user_id = $1;
`

const deleteBondQuery = `
	delete from bonds
	where id = $1;
`

const getUniqueCryptoCoinsIDs = `
	select distinct coin_id
	from crypto_positions;
`
