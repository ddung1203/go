# Banking

#### export 

Go가 private을 보호하는 방법과 export하는 방법은 간단하다. 오직 변수/함수 이름을 대문자로 시작하면 public으로 export 가능하다.

`accounts/accounts.go`
```go
package accounts

type Account struct {
	Owner string
	Balance int
}
```

`main.go`
```go
package main

import (
	"fmt"

	accounts "github.com/ddung1203/go/accounts"
)

func main() {
	account := accounts.Account{Owner: "Joongseok", Balance: 10000}
	fmt.Println(account)
}
```

하지만 상기와 같이 작성 시 public으로 누구나 임의로 balance를 작성할 수 있다.

> ex) `account.Owner = "foo"`

#### struct

Go의 흔한 패턴으로, constructor가 없기 때문에 function으로 construct하거나 struct를 만들도록 한다.

`accounts/accounts.go`
```go
package accounts

// Account struct
type Account struct {
	owner string
	balance int
}

// NewAccount create Account
func NewAccount(owner string) *Account{
	account := Account{owner: owner, balance: 0}
	return &account
}

// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// Balance of your account
func (a Account) Balance() int{
	return a.balance
}
```

`main.go`
```go
package main

import (
	"fmt"

	accounts "github.com/ddung1203/go/accounts"
)

func main() {
	account := accounts.NewAccount("Joongseok")
	account.Deposit(10)
	fmt.Println(account.Balance())
}
```

상기와 같이 function을 만들어서 object를 return 시킨다. 실제 메모리 address를 return 하며, 이와 같이 작성 시 account의 owner와 balance를 변경 불가능하다.

balance를 증가시키고 싶을 떈, method를 사용한다. method는 function과 달리 return값이 없으며, func와 function 사이에 receiver를 작성한다. 보통 struct의 첫 글자를 따서 작성한다.

여기서, `account.go`의 `Deposit()`에서 `*`을 작성하지 않으면 변경사항이 반영되지 않는다. 이유는, Go에서 object와 struct에 관여하는 부분 때문이다. Go는 function이나 method 등을 보내는 순간에 보안을 위해 복사본을 만든다. 이 경우 `account.Deposit()`을 실행 시 Go는 account를 받아 복사본을 만들며, 여기서 받는 `a`는 account이긴 하지만, `main.go`에서의 account의 복사본이다. 따라서 실제 account가 아니다.

`Balance()`의 경우 복사본의 여부와 상관없이 balance만 원하기 때문에 object를 `main.go`에서 복사해서 object의 balance를 받아온다.

즉, `Deposit()`의 경우 balance의 증가된 결과도 원하기 때문에 복사본을 만들지 않고 account의 balance를 증가시키는 것이며 Deposit method를 호출한 account를 사용하라는 것이다.

#### withdraw

이제, Withdraw 기능을 위해, 하기와 같이 작성한다.

Deposit의 개념과 같이 복사본을 만들지 않고 account의 balance를 감소시키도록 한다.

하지만, balance는 음수(`-`)가 될 수 없기 때문에 error를 발생하도록 한다. Go에서 error의 return은 error/nil 이 있다. none과 같은 개념이다.

또한 Go에서 exception이 없기 때문에 error를 다룰 코드를 작성해야 한다.

`accounts/accounts.go`
```go
var errNoMoney = errors.New("error: Can't withdraw")

// Withdraw * amount from yout account
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}
```

`main.go`
```go
func main() {
	account := accounts.NewAccount("Joongseok")
	account.Deposit(10)
	fmt.Println(account.Balance())
	err := account.Withdraw(20)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(account.Balance())
}
```

#### ChangeOwner, Owner

Go가 struct로 자동으로 호출해주는 method가 있다. 그 중 하나가 `String()`이다. 

return을 원하는 것을 작성하면 된다.

```go
// ChangeOwner of the account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// Owner of the account
func (a Account) Owner() string{
	return a.owner
}

func (a Account) String() string{
	return fmt.Sprint(a.Owner(), "'s account.\nHas: ", a.Balance())
}
```

```go
func main() {
	account := accounts.NewAccount("Joongseok")
	account.Deposit(10)
	err := account.Withdraw(5)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(account.String())
}
```

전체 코드는 하기와 같다.

`accounts/accounts.go`
```go
package accounts

import (
	"errors"
	"fmt"
)

// Account struct
type Account struct {
	owner string
	balance int
}

var errNoMoney = errors.New("error: Can't withdraw")

// NewAccount create Account
func NewAccount(owner string) *Account{
	account := Account{owner: owner, balance: 0}
	return &account
}

// Deposit x amount on your account
func (a *Account) Deposit(amount int) {
	a.balance += amount
}

// Balance of your account
func (a Account) Balance() int{
	return a.balance
}

// Withdraw * amount from yout account
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return errNoMoney
	}
	a.balance -= amount
	return nil
}

// ChangeOwner of the account
func (a *Account) ChangeOwner(newOwner string) {
	a.owner = newOwner
}

// Owner of the account
func (a Account) Owner() string{
	return a.owner
}

func (a Account) String() string{
	return fmt.Sprint(a.Owner(), "'s account.\nHas: ", a.Balance())
}
```

`main.go`
```go
package main

import (
	"fmt"
	"log"

	accounts "github.com/ddung1203/go/accounts"
)

func main() {
	account := accounts.NewAccount("Joongseok")
	account.Deposit(10)
	err := account.Withdraw(5)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(account.String())
}
```