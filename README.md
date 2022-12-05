## Installation

```bash
go get -u github.com/LineoIT/sqlw
```

## How to use

- select all(\*)

```go
	q, args := sqlw.Select("authors").
		Where("id", "=", 11).
		Build()
	fmt.Println(sqlw.Debug(q, args...))
```

resultat:

```sql
select * from authors where id = 11
```

- insert query

```go
	type User struct {
		ID        int
		Email     string
		Phone     string
		Role      string
		Age       int
		UpdatedAt time.Time
	}

	user := User{
		Email: "abc@example.com",
		Phone: "+7939739473",
		Age:   29,
		ID:    3,
	}
    // insert
    q, args := sqlw.Insert("users").
		Value("email", user.Email).
		Value("phone", sqlw.Nullif(user.Phone, "''")).
		Value("age", sqlw.Coalesce(user.Age, "age")).
		Returns("id").
		Build()
	fmt.Println(sqlw.Debug(q, args...))

```

resultat:

```sql
insert into users(email,phone,age) values(abc@example.com,nullif(+7939739473,''),coalesce(29,age)) returning id
```

- Update query

```go
	// update
	q, args = sqlw.Update("users").
		Set("email", user.Email).
		Set("phone", sqlw.Nullif(user.Phone, "''")).
		Set("age", sqlw.Coalesce(user.Age, "age")).
		Set("role", sqlw.Coalesce(sqlw.Nullif(user.Role, "''"), "role")).
		Where("id", "=", user.ID).
		Returns("updated_at").
		Build()
	fmt.Println(sqlw.Debug(q, args...))
```

resultat:

```sql
update users set email=abc@example.com,phone=nullif(+7939739473,''),age=coalesce(29,age),role=coalesce(nullif(,''),role) returning updated_at where id = 3
```
