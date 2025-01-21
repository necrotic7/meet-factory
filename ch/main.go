package main

import (
	"fmt"
	"log"
	"time"
)

type Meat struct {
	Name     string
	Time     int64
	Quantity int
}

type Employee struct {
	Name string
}

// 宣告員工清單
var EmployeeList []Employee = []Employee{
	{Name: "A"},
	{Name: "B"},
	{Name: "C"},
	{Name: "D"},
	{Name: "E"},
}

// 宣告肉品
var Beef Meat = Meat{Name: "牛肉", Time: 1, Quantity: 10}
var Pork Meat = Meat{Name: "豬肉", Time: 2, Quantity: 7}
var Chicken Meat = Meat{Name: "雞肉", Time: 3, Quantity: 5}

func main() {
	// 計算所有肉品數量
	allMeatCounts := Beef.Quantity + Pork.Quantity + Chicken.Quantity
	// 建立一個channel放所有的肉品
	meatsCh := make(chan Meat, allMeatCounts)
	// 將肉品放入channel
	go func() {
		for i := 0; i < Beef.Quantity; i++ {
			meatsCh <- Beef
		}
		for i := 0; i < Pork.Quantity; i++ {
			meatsCh <- Pork
		}
		for i := 0; i < Chicken.Quantity; i++ {
			meatsCh <- Chicken
		}
	}()

	done := make(chan int)

	for _, emp := range EmployeeList {
		go func(e Employee) {

			for {
				meat, ok := <-meatsCh

				if !ok {
					log.Printf("所有肉品已處理完畢，員工 %s 結束工作\n", e.Name)
					return
				}

				if err := meetProcessor(meat, e); err != nil {
					log.Print(err)
					meatsCh <- meat
				} else {
					done <- 1
				}
			}
		}(emp)
	}

	for i := 0; i < allMeatCounts; i++ {
		<-done
	}

	// close channel
	close(meatsCh)
	close(done)
	log.Printf("所有員工結束工作")
}

func meetProcessor(meat Meat, emp Employee) (err error) {
	log.Printf("%s 取得%s \n", emp.Name, meat.Name)
	if time.Now().Second()%2 == 0 {
		return fmt.Errorf("%s 處理%s失敗 \n", emp.Name, meat.Name)
	}
	// 處理肉品
	time.Sleep(time.Duration(meat.Time) * time.Second)
	log.Printf("%s 處理完%s \n", emp.Name, meat.Name)
	return
}
