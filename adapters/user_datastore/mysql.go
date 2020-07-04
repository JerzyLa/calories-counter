package user_datastore

import (
	"calories-counter/models"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strings"
)

// implements models.Users_Datastore
type MySQLStore struct {
	db *sqlx.DB
}

func NewMySQLStore(dbSourceName string) (*MySQLStore, error) {
	log.Info("connecting to DB")
	driver := "mysql"
	sqldb, err := sql.Open(driver, dbSourceName)
	if err != nil {
		log.WithError(err).Error("couldn't open db")
		return nil, err
	}

	db := sqlx.NewDb(sqldb, driver)

	err = db.Ping()
	if err != nil {
		log.WithError(err).Error("couldn't ping db")
		return nil, err
	}

	return &MySQLStore{db: db}, nil
}

func (d *MySQLStore) Close() error {
	return d.db.Close()
}

func (d *MySQLStore) GetUser(accountID, username string) (*models.User, error) {
	query := d.db.Rebind(`SELECT id, account_id, username, role_id FROM users WHERE account_id=? AND username=?`)
	row := d.db.QueryRowx(query, accountID, username)

	var user models.User
	err := row.StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (d *MySQLStore) GetUserPassword(accountID, username string) (*string, error) {
	query := d.db.Rebind(`SELECT password FROM users WHERE account_id=? AND username=?`)
	row := d.db.QueryRowx(query, accountID, username)

	var pass string
	err := row.Scan(&pass)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &pass, nil
}

func (d *MySQLStore) GetUserById(accountID, userID string) (*models.User, error) {
	query := d.db.Rebind(`SELECT id, account_id, username, role_id FROM users WHERE account_id=? AND id=?`)
	row := d.db.QueryRowx(query, accountID, userID)

	var user models.User
	err := row.StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (d *MySQLStore) SaveUser(accountID, username, pass string, role int) (*models.User, error) {
	id := uuid.New().String()

	query := d.db.Rebind(`INSERT INTO users (id, account_id, username, password, role_id) VALUES (?, ?, ?, ?, ?);`)
	_, err := d.db.Exec(query, id, accountID, username, pass, role)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, models.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &models.User{
		ID:        id,
		AccountID: accountID,
		Username:  username,
		RoleID:    role,
	}, nil
}

func (d *MySQLStore) SaveRootUser(username, pass string) (*models.User, error) {
	id := uuid.New().String()
	accountID := uuid.New().String()

	query := d.db.Rebind(`SELECT id, account_id, username, role_id FROM users WHERE username=? AND role_id=?`)
	row := d.db.QueryRowx(query, username, models.OwnerRole)

	var user models.User
	err := row.StructScan(&user)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == nil {
		return nil, models.ErrAccountAlreadyExists
	}

	query = d.db.Rebind(`INSERT INTO users (id, account_id, username, password, role_id) VALUES (?, ?, ?, ?, ?);`)
	_, err = d.db.Exec(query, id, accountID, username, pass, models.OwnerRole)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, models.ErrAccountAlreadyExists
		}
		return nil, err
	}

	return &models.User{
		ID:        id,
		AccountID: accountID,
		Username:  username,
		RoleID:    models.OwnerRole,
	}, nil
}

func (d *MySQLStore) UpdateUser(user models.User) (*models.User, error) {
	query := d.db.Rebind(`UPDATE users SET username=:username, role_id=:role_id WHERE account_id=:account_id AND id=:id`)
	_, err := d.db.NamedExec(query, user)
	if err != nil {
		if isDuplicateKeyError(err) {
			return nil, models.ErrUserAlreadyExists
		}
		return nil, err
	}

	return &user, nil
}

func (d *MySQLStore) DeleteUser(accountID string, userID string) error {
	query := d.db.Rebind(`DELETE FROM users WHERE account_id=? AND id=?`)
	res, err := d.db.Exec(query, accountID, userID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return models.ErrUserNotFound
	}

	return nil
}

func (d *MySQLStore) SaveMeal(userID string, meal models.Meal) (*models.Meal, error) {
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, err
	}

	mealId := uuid.New().String()
	query := tx.Rebind(`INSERT INTO users_meals (id, user_id, name, date, time, calories) VALUES (?, ?, ?, ?, ?, ?);`)
	_, err = tx.Exec(query, mealId, userID, meal.Name, meal.Date, meal.Time, meal.Calories)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	caloriesDeficit, err := updateCaloriesDeficit(tx, userID, meal.Date)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &models.Meal{
		ID:              mealId,
		Date:            meal.Date,
		Time:            meal.Time,
		Name:            meal.Name,
		Calories:        meal.Calories,
		CaloriesDeficit: *caloriesDeficit,
	}, nil
}

func (d *MySQLStore) UpdateMeal(userID string, newMeal models.Meal) (*models.Meal, error) {

	query := d.db.Rebind(`SELECT date FROM users_meals WHERE user_id=? AND id=?`)
	row := d.db.QueryRowx(query, userID, newMeal.ID)

	var oldMealDateStr string
	err := row.Scan(&oldMealDateStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrMealNotFound
		}
		return nil, err
	}

	tx, err := d.db.Beginx()
	if err != nil {
		return nil, err
	}

	query = tx.Rebind(`UPDATE users_meals SET name=?, date=?, time=?, calories=? WHERE user_id=? AND id=?`)
	_, err = tx.Exec(query, newMeal.Name, newMeal.Date, newMeal.Time, newMeal.Calories, userID, newMeal.ID)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	caloriesDeficit, err := updateCaloriesDeficit(tx, userID, newMeal.Date)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	oldMealDateStr = oldMealDateStr[:10]
	if newMeal.Date != oldMealDateStr {
		_, err := updateCaloriesDeficit(tx, userID, oldMealDateStr)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &models.Meal{
		ID:              newMeal.ID,
		Date:            newMeal.Date,
		Time:            newMeal.Time,
		Name:            newMeal.Name,
		Calories:        newMeal.Calories,
		CaloriesDeficit: *caloriesDeficit,
	}, nil
}

func (d *MySQLStore) GetUsers(accountID string, page, perPage int, filter string) (models.UserSlice, error) {
	res := models.UserSlice{Items: make([]models.User, 0)}
	filterQuery, err := buildFilterMealsQuery(filter)
	if err != nil {
		return res, err
	}

	query := d.db.Rebind(fmt.Sprintf(`SELECT count(*)  
								FROM users
								WHERE account_id = ? %s`, filterQuery))
	row := d.db.QueryRowx(query, accountID)
	var total int
	err = row.Scan(&total)
	if err != nil {
		if isInvalidQuery(err) {
			return res, models.ErrInvalidQuery
		}
		return res, err
	}

	query = d.db.Rebind(fmt.Sprintf(`SELECT id, account_id, username, role_id    
								FROM users
								WHERE account_id = ? %s 
								ORDER BY id DESC 
								LIMIT ? OFFSET ?;`, filterQuery))
	rows, err := d.db.Queryx(query, accountID, perPage, page*perPage)
	if err != nil {
		if isInvalidQuery(err) {
			return res, models.ErrInvalidQuery
		}
		return res, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var user models.User
		err := rows.StructScan(&user)
		if err != nil {
			if err == sql.ErrNoRows {
				return res, nil
			}
			return res, err
		}
		res.Items = append(res.Items, user)
	}
	if err := rows.Err(); err != nil {
		return res, err
	}
	res.Total = total
	return res, nil
}

func (d *MySQLStore) GetMeals(userID string, page, perPage int, filter string) (models.MealSlice, error) {
	res := models.MealSlice{Items: make([]models.Meal, 0)}
	filterQuery, err := buildFilterMealsQuery(filter)
	if err != nil {
		return res, err
	}

	query := d.db.Rebind(fmt.Sprintf(`SELECT count(*)  
								FROM users_meals AS m
								LEFT JOIN users_calories AS c ON m.user_id = c.user_id AND m.date = c.date
								WHERE m.user_id=? %s`, filterQuery))
	row := d.db.QueryRowx(query, userID)
	var total int
	err = row.Scan(&total)
	if err != nil {
		if isInvalidQuery(err) {
			return res, models.ErrInvalidQuery
		}
		return res, err
	}

	query = d.db.Rebind(fmt.Sprintf(`SELECT m.id, m.date, m.time, m.name, m.calories, c.calories_deficit  
								FROM users_meals AS m 
								LEFT JOIN users_calories AS c ON m.user_id = c.user_id AND m.date = c.date
								WHERE m.user_id=? %s
								ORDER BY m.date, m.time DESC 
								LIMIT ? OFFSET ?;`, filterQuery))
	rows, err := d.db.Queryx(query, userID, perPage, page*perPage)
	if err != nil {
		if isInvalidQuery(err) {
			return res, models.ErrInvalidQuery
		}
		return res, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var meal models.Meal
		var dateStr, timeStr string
		err := rows.Scan(&meal.ID, &dateStr, &timeStr, &meal.Name, &meal.Calories, &meal.CaloriesDeficit)
		if err != nil {
			if err == sql.ErrNoRows {
				return res, nil
			}
			return res, err
		}
		meal.Date = dateStr[:10]
		meal.Time = timeStr
		res.Items = append(res.Items, meal)
	}
	if err := rows.Err(); err != nil {
		return res, err
	}
	res.Total = total
	return res, nil
}

func (d *MySQLStore) GetMeal(userID, mealID string) (*models.Meal, error) {
	query := d.db.Rebind(`SELECT m.id, m.date, m.time, m.name, m.calories, c.calories_deficit 
								FROM users_meals AS m
								LEFT JOIN users_calories AS c ON m.user_id = c.user_id AND m.date = c.date
								WHERE m.user_id=? AND m.id=?`)
	row := d.db.QueryRowx(query, userID, mealID)
	var meal models.Meal
	var dateStr, timeStr string
	err := row.Scan(&meal.ID, &dateStr, &timeStr, &meal.Name, &meal.Calories, &meal.CaloriesDeficit)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrMealNotFound
		}
		return nil, err
	}
	meal.Date = dateStr[:10]
	meal.Time = timeStr

	return &meal, nil
}

func (d *MySQLStore) DeleteMeal(userID string, mealID string) error {
	query := d.db.Rebind(`SELECT date FROM users_meals WHERE user_id=? AND id=?`)
	row := d.db.QueryRowx(query, userID, mealID)

	var oldMealDateStr string
	err := row.Scan(&oldMealDateStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.ErrMealNotFound
		}
		return err
	}

	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	query = d.db.Rebind(`DELETE FROM users_meals WHERE user_id=? AND id=?`)
	_, err = d.db.Exec(query, userID, mealID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = updateCaloriesDeficit(tx, userID, oldMealDateStr[:10])
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *MySQLStore) UpdateSettings(userID string, settings models.Settings) (*models.Settings, error) {
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, err
	}

	query := tx.Rebind(`SELECT id FROM users_settings WHERE user_id=?`)
	row := tx.QueryRowx(query, userID)
	var id string
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			query = tx.Rebind(`INSERT INTO users_settings (user_id, expected_daily_calories) VALUES (?, ?)`)
			_, err = tx.Exec(query, userID, settings.ExpectedDailyCalories)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
		} else {
			_ = tx.Rollback()
			return nil, err
		}
	} else {
		query := tx.Rebind(`UPDATE users_settings SET expected_daily_calories=? WHERE user_id=?`)
		_, err := tx.Exec(query, settings.ExpectedDailyCalories, userID)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
	}

	query = tx.Rebind(`UPDATE users_calories SET calories_deficit=1 WHERE user_id=? AND total_calories<?`)
	_, err = tx.Exec(query, userID, settings.ExpectedDailyCalories)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	query = tx.Rebind(`UPDATE users_calories SET calories_deficit=0 WHERE user_id=? AND total_calories>=?`)
	_, err = tx.Exec(query, userID, settings.ExpectedDailyCalories)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

func (d *MySQLStore) GetSettings(userID string) (*models.Settings, error) {
	query := d.db.Rebind(`SELECT expected_daily_calories FROM users_settings WHERE user_id=?`)
	row := d.db.QueryRowx(query, userID)

	var settings models.Settings
	_ = row.StructScan(&settings)

	return &settings, nil
}

func isDuplicateKeyError(err error) bool {
	switch mErr := err.(type) {
	case *mysql.MySQLError:
		if mErr.Number == 1062 {
			return true
		}
	}
	return false
}

func isInvalidQuery(err error) bool {
	switch mErr := err.(type) {
	case *mysql.MySQLError:
		if mErr.Number == 1054 || mErr.Number == 1064 {
			return true
		}
	}
	return false
}

func updateCaloriesDeficit(tx *sqlx.Tx, userID string, date string) (*bool, error) {
	query := tx.Rebind(`SELECT COALESCE(SUM(calories), 0) FROM users_meals WHERE user_id=? AND date=?`)
	row := tx.QueryRowx(query, userID, date)
	var totalCalories int
	err := row.Scan(&totalCalories)
	if err != nil {
		return nil, err
	}

	query = tx.Rebind(`SELECT expected_daily_calories FROM users_settings WHERE user_id=?`)
	row = tx.QueryRowx(query, userID)
	var settings models.Settings
	_ = row.StructScan(&settings)

	var caloriesDeficit bool
	if totalCalories < settings.ExpectedDailyCalories {
		caloriesDeficit = true
	}

	query = tx.Rebind(`SELECT id FROM users_calories WHERE user_id=? AND date=?`)
	row = tx.QueryRowx(query, userID, date)
	var id string
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			query = tx.Rebind(`INSERT INTO users_calories (user_id, date, total_calories, calories_deficit) VALUES (?, ?, ?, ?)`)
			_, err = tx.Exec(query, userID, date, totalCalories, caloriesDeficit)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		query = tx.Rebind(`UPDATE users_calories SET total_calories=?, calories_deficit=? WHERE user_id=? AND date=?`)
		_, err := tx.Exec(query, totalCalories, caloriesDeficit, userID, date)
		if err != nil {
			return nil, err
		}
	}

	return &caloriesDeficit, nil
}

func buildFilterMealsQuery(filter string) (string, error) {
	if filter == "" {
		return "", nil
	}

	res := strings.ReplaceAll(filter, "%20", " ")
	res = strings.ReplaceAll(res, "eq", "=")
	res = strings.ReplaceAll(res, "eq", "<>")
	res = strings.ReplaceAll(res, "gt", ">")
	res = strings.ReplaceAll(res, "lt", "<")
	res = strings.ReplaceAll(res, "date", "m.date")

	for _, c := range res {
		if c == '*' || c == ';' || c == '&' || c == '"' || c == '\\' {
			return "", models.ErrInvalidFilter
		}
	}

	res = " AND " + res
	return res, nil
}
