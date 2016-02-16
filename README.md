# Validations

Validations is a [GORM](https://github.com/jinzhu/gorm) extension, used to validate models when creating, updating

### Register GORM Callbacks

It is using [GORM](https://github.com/jinzhu/gorm)'s callbacks to handle validations, so you need to register callbacks to gorm DB first, do it like:

```go
import (
  "github.com/jinzhu/gorm"
  "github.com/qor/validations"
)

db, err := gorm.Open("sqlite3", "demo_db")

validations.RegisterCallbacks(db)
```

### Usage

After register callbacks, Creating, updating will trigger the `Validate` method defined for struct, if the method has added or returned any error, the process will be rollbacked.

```go
type User struct {
  gorm.Model
  Age uint
}

func (user User) Validate(db *gorm.DB) {
  if user.Age <= 18 {
    db.AddError(errors.New("age need to be 18+"))
  }
}

db.Create(&User{Age: 10})         // won't insert the record into database

var user User{Age: 20}
db.Create(&user)                  // user with age 20 will be inserted into database
db.Model(&user).Update("age", 10) // user's age won't be updated, and will return error `age need to be 18+`

// If you have added more than one error, could get them with `db.GetErrors()`

func (user User) Validate(db *gorm.DB) {
  if user.Age <= 18 {
    db.AddError(errors.New("age need to be 18+"))
  }
  if user.Name == "" {
    db.AddError(errors.New("name can't be blank"))
  }
}

db.Create(&User{}).GetErrors() // => []error{}
```

## [Qor Support](https://github.com/qor/qor)

[QOR](http://getqor.com) is architected from the ground up to accelerate development and deployment of Content Management Systems, E-commerce Systems, and Business Applications, and comprised of modules that abstract common features for such system.

Validations could be used alone, and it works nicely with QOR, if you have requirements to manage your application's data, be sure to check QOR out!

[QOR Demo:  http://demo.getqor.com/admin](http://demo.getqor.com/admin)

If you want to display errors for each form field in Qor Admin, you could register your error like this:

```go
func (user User) Validate(db *gorm.DB) {
  if user.Age <= 18 {
    db.AddError(validations.NewError(user, "Age", "age need to be 18+"))
  }
}
```

Checkout [http://demo.getqor.com/admin/products/1](http://demo.getqor.com/admin/products/1) as demo, change Name to be blank string and save to see what happens.

## License

Released under the [MIT License](http://opensource.org/licenses/MIT).