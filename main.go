package main

import (
	"log"
	"sync"
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
		// close channel
		close(meatsCh)
	}()

	var wg sync.WaitGroup
	var mu sync.Mutex // 確保肉品的安全分配

	for _, emp := range EmployeeList {
		// 每一個員工起一個go routine
		wg.Add(1)
		go func(e Employee) {
			// 完成後關閉go routine
			defer wg.Done()

			for {
				// 取得肉品前先上鎖，避免同時有員工取得肉品
				mu.Lock()
				meat, ok := <-meatsCh
				// 取得肉品後解鎖
				mu.Unlock()

				if !ok {
					log.Printf("所有肉品已處理完畢，員工 %s 結束工作", e.Name)
					return
				}

				meetProcessor(meat, e)
			}
		}(emp)
	}

	// 等待所有 go routine結束
	wg.Wait()
	log.Printf("所有員工結束工作")
}

func meetProcessor(meat Meat, emp Employee) {
	log.Printf("%s 取得%s \n", emp.Name, meat.Name)
	// 處理肉品
	time.Sleep(time.Duration(meat.Time) * time.Second)
	log.Printf("%s 處理完%s \n", emp.Name, meat.Name)
}
