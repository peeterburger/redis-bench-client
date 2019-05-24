package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/urfave/cli"
)

func main() {

	// Building CLI Interface
	var app = cli.NewApp()

	app.Name = "redis-bench-client"
	app.Usage = "A CLI to perform redis benches"
	app.Author = "Peter Burger"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		cli.Command{
			Name:    "execute",
			Aliases: []string{"exec"},
			Usage:   "Performs a bench",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host",
					Value: "127.0.0.1:6379",
					Usage: "Host to perform the bench on - (<host>:<port>)",
				},
				cli.IntFlag{
					Name:  "count",
					Value: 1,
					Usage: "The amount of SET operations to be performed",
				},
				cli.BoolFlag{
					Name:  "silent",
					Usage: "Hides all log information"},
			},
			Action: func(c *cli.Context) error {
				fmt.Fprintf(c.App.Writer, "Evaluating flags...\n")

				fmt.Fprintf(c.App.Writer, "[host] -> %s\n", c.String("host"))
				host := c.String("host")

				fmt.Fprintf(c.App.Writer, "[count] -> %d\n", c.Int("count"))
				count := c.Int("count")

				// fmt.Fprintf(c.App.Writer, "[silent] -> %b\n", c.Bool("silent"))
				// silent := c.Bool("silent")

				client := redis.NewClient(&redis.Options{
					Addr:     host,
					Password: "", // no password set
					DB:       0,  // use default DB
				})

				fmt.Fprintf(c.App.Writer, "Connecting to %s... ", host)
				pong, _ := client.Ping().Result()
				fmt.Fprintf(c.App.Writer, pong)

				fmt.Fprintf(c.App.Writer, "\n\n--- Benchmark Phase ---")

				t := time.Now()

				for i := 0; i <= count; i++ {
					fmt.Fprintf(c.App.Writer, "\nSET -> key: %d, val: %d", i, i)
					client.Set(strconv.Itoa(i), strconv.Itoa(i), 0)
				}

				end := time.Since(t)

				fmt.Fprintf(c.App.Writer, "\n\n --- Result ---\n")
				fmt.Fprintf(c.App.Writer, "Total operations: %d\n", count)
				fmt.Fprintf(c.App.Writer, "Total duration: %f s\n", end.Seconds())
				fmt.Fprintf(c.App.Writer, "Average operations per second: %f O/s\n", float64(count)/end.Seconds())
				fmt.Fprintf(c.App.Writer, "Average time per operation: %f s\n", end.Seconds()/float64(count))
				return nil
			},
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "%q not implemented.\n", command)
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	_ = app.Run(os.Args)

	/*
		client := redis.NewClient(&redis.Options{
			Addr:     host,
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		fmt.Print("\nConnecting to ", host, "...")
		pong, _ := client.Ping().Result()
		fmt.Println(pong)

		fmt.Print("\n\n--- Benchmark Phase ---\n\n")

		t := time.Now()

		for i := 0; i < count; i++ {
			fmt.Printf("\nSET -> key: %d, val: %d", i, i)
			client.Set(strconv.Itoa(i), strconv.Itoa(i), 0)
		}

		end := time.Since(t)

		fmt.Print("\n\n --- Bench --- \n\n")
		fmt.Printf("Total operations: %d\n", count)
		fmt.Printf("Total duration: %fs\n", end.Seconds())
		fmt.Printf("Average operations per second: %fO/s\n", float64(count)/end.Seconds())
		fmt.Printf("Average time per operation: %fs\n", end.Seconds()/float64(count))
	*/
}
