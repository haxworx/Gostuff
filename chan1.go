package main

import "fmt"
import "time"

const limit int = 32

func calculate(x int, c chan int)  {
        c <- x * x
}

func calc2(upper int, c chan int) {
        for i := 1; i <=  upper; i++ {
                c <- i * i
        }
}

func main() {
        c := make(chan int, limit)

        for i := 1; i <= limit; i++ {
               go calculate(i, c)
        }
        fmt.Println("All started!")

        time.Sleep(3 * time.Second)
        for len(c) > 0 {
                i := <-c
                fmt.Println(i)
        }

        c = make(chan int, limit)

        go calc2(cap(c), c)

        time.Sleep(3 * time.Second)

        for len(c) > 0 {
                i := <-c
                fmt.Println(i)
        }
}
