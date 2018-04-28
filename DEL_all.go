package main

import (
    "fmt"
    "bufio"

    "github.com/garyburd/redigo/redis"
    "os"
)
func main() {

    c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialPassword("123456"))
    if err != nil {
        fmt.Println("Connect to redis error", err)
        return
    }
    defer c.Close()
    c.Do("SELECT", 1)

    fmt.Println("Enter the pattern of KEYs you want to delete, enter 'quit' to exit.")

    reader := bufio.NewReader(os.Stdin)
    for{
        fmt.Print("DEL KEYS >> ")
        KEYS, _, _ := reader.ReadLine()
        if(string(KEYS) == "quit"){
            break
        }
        values, _ := redis.Values(c.Do("KEYS", KEYS))
        //for i := 0; i < len(values); i++ {
        for i := 0; i < len(values); i++ {

            _, err := c.Do("DEL", string((values[i]).([]byte)))
            if err != nil {
                fmt.Println("DEL failed: ", err)
                return
            }

        }
        fmt.Println("OK")

    }




}