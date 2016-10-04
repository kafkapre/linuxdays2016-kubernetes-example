package main


import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mediocregopher/radix.v2/redis"
	"syscall"
	"os"
)

type Person struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func createKey(id string)  string {
	return "persons/" + id
}

func (p Person) storePerson(redis *redis.Client)  (string, error){
	key := createKey(p.Id)
	return redis.Cmd("HMSET", key, "id", p.Id, "name", p.Name, "surname", p.Surname).Str()
}

func fetchPerson(redis *redis.Client, key string) (Person, error) {
	res, err := redis.Cmd("HMGET", key, "id", "name", "surname").List()
	return Person{Id:res[0], Name:res[1], Surname:res[2]}, err
}

func existPerson(redis *redis.Client, id string)  (bool, error) {
	res, err := redis.Cmd("EXISTS", createKey(id)).Int()
	return (res == 1), err
}

func httpGetPersons(c *gin.Context, redis *redis.Client) {
	var (
		result gin.H
	)
	id := c.Param("id")

	exist, err := existPerson(redis, id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if  !exist {
		result = gin.H{"Msg": "Not found"}
		c.JSON(http.StatusNotFound, result)
	} else {
		person, err := fetchPerson(redis, createKey(id))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else {
			result = gin.H{"person": person}
			c.JSON(http.StatusOK, result)
		}
	}
}

func httpGetAllPersons(c *gin.Context, redis *redis.Client) {
	var (
		persons []Person
	)

	list, err := redis.Cmd("KEYS", "*").List()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	for _, key := range list{
		p, err:= fetchPerson(redis, key)
		if err != nil {
			fmt.Print(err.Error())
		}else{
			persons = append(persons, p)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"persons": persons,
		"count":  len(persons),
	})
}


func httpPostPerson(c *gin.Context, redis *redis.Client) {
	var p Person
	c.BindJSON(&p)

	if (len(p.Id) == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Id cannot be empty."})
	}else {
		exist, err := existPerson(redis, p.Id)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		} else if  exist {
			result := gin.H{"Msg": fmt.Sprintf("Person with id: %s already exists.", p.Id) }
			c.JSON(http.StatusNotFound, result)
		}else {
			httpStorePerson(c, redis, p)
		}
	}
}

func httpPutPerson(c *gin.Context, redis *redis.Client) {
	id := c.Param("id")

	if (len(id) == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Id cannot be empty."})
	}else {
		var p Person
		c.BindJSON(&p)
		p.Id = id

		httpStorePerson(c, redis, p)
	}
}

func httpStorePerson(c *gin.Context, redis *redis.Client, p Person) {
	_, err := p.storePerson(redis)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, gin.H{"stored_person": p})
}

func httpDeletePerson(c *gin.Context, redis *redis.Client) {
	id := c.Param("id")
	res, err := redis.Cmd("DEL", createKey(id)).Int()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	} else if res > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully deleted person with id: %s.", id),
		})
	}else {
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Person with id: %s was not found.", id),
		})
	}
}

func main() {
	redis := connectToRedis()
	if redis == nil {
		fmt.Println("Attempts to connect to redis faild. Server will be stoped")
		syscall.Exit(1)
	}

	router := gin.Default()

	// GET a person detail
	router.GET("/person/:id", func(c *gin.Context) {httpGetPersons(c, redis)})

	// GET all persons
	router.GET("/persons", func(c *gin.Context) {httpGetAllPersons(c, redis)})

	// POST new person details
	router.POST("/person", func(c *gin.Context) {httpPostPerson(c, redis)})
	router.POST("/person/", func(c *gin.Context) {httpPostPerson(c, redis)})

	// PUT - update a person details
	router.PUT("/person/:id", func(c *gin.Context) {httpPutPerson(c, redis)})

	// Delete resources
	router.DELETE("/person/:id", func(c *gin.Context) {httpDeletePerson(c, redis)})
	router.Run(":3000")
}

func connectToRedis()  (*redis.Client){
		url:=obtainRedisUrl()
		for i:=0; i < 20; i++ {
				redisClient, err := redis.Dial("tcp", url)
				if err == nil {
						fmt.Println("Connected to redis.")
						return redisClient
				}
				fmt.Printf("Tried to connect to redis [attempt: %d].\n", i)
				time.Sleep(5 * time.Second)
		}
		return nil
}

func obtainRedisUrl() string {
	redisIp := os.Getenv("REDIS_IP")
	redisPort := os.Getenv("REDIS_PORT")

	if len(redisIp) == 0 {
		redisIp = "localhost"
	}
	if len(redisPort) == 0 {
		redisPort="6379"
	}

	return redisIp + ":" + redisPort
}
