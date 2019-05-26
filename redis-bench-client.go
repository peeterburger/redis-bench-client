package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/urfave/cli"
)

func main() {

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
				cli.BoolFlag{
					Name:  "middleware",
					Usage: "Sends raw data to server, which then accesses the database instead"},
			},
			Action: func(c *cli.Context) error {
				fmt.Fprintf(c.App.Writer, "Evaluating flags...\n")

				fmt.Fprintf(c.App.Writer, "[host] -> %s\n", c.String("host"))
				host := c.String("host")

				fmt.Fprintf(c.App.Writer, "[count] -> %d\n", c.Int("count"))
				count := c.Int("count")

				//fmt.Fprintf(c.App.Writer, "[silent] -> %d\n", c.Bool("silent"))
				//silent := c.Bool("silent")

				fmt.Fprintf(c.App.Writer, "[middleware] -> %b\n", c.Bool("middleware"))
				middleware := c.Bool("middleware")

				if middleware {

					tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
					checkError(err)

					conn, err := net.DialTCP("tcp", nil, tcpAddr)
					checkError(err)

					t := time.Now()

					for i := 0; i < count; i++ {
						fmt.Fprintf(c.App.Writer, "\nSET -> key: %d, val: %d", i, i)
						_, err = conn.Write([]byte(strconv.Itoa(i)))
						checkError(err)
					}

					end := time.Since(t)

					fmt.Fprintf(c.App.Writer, "\n\n --- Result ---\n")
					fmt.Fprintf(c.App.Writer, "Total operations: %d\n", count)
					fmt.Fprintf(c.App.Writer, "Total duration: %f s\n", end.Seconds())
					fmt.Fprintf(c.App.Writer, "Average operations per second: %f O/s\n", float64(count)/end.Seconds())
					fmt.Fprintf(c.App.Writer, "Average time per operation: %f s\n", end.Seconds()/float64(count))

					os.Exit(0)

				} else {

					client := redis.NewClient(&redis.Options{
						Addr:     host,
						Password: "",
						DB:       0,
					})

					fmt.Fprintf(c.App.Writer, "Connecting to %s... ", host)
					pong, _ := client.Ping().Result()
					fmt.Fprintf(c.App.Writer, pong)

					fmt.Fprintf(c.App.Writer, "\n\n--- Benchmark Phase ---")

					t := time.Now()

					for i := 0; i < count; i++ {
						fmt.Fprintf(c.App.Writer, "\nSET -> key: %d, val: %d", i, i)
						client.Set(strconv.Itoa(i), strconv.Itoa(i), 0)
					}

					end := time.Since(t)

					fmt.Fprintf(c.App.Writer, "\n\n --- Result ---\n")
					fmt.Fprintf(c.App.Writer, "Total operations: %d\n", count)
					fmt.Fprintf(c.App.Writer, "Total duration: %f s\n", end.Seconds())
					fmt.Fprintf(c.App.Writer, "Average operations per second: %f O/s\n", float64(count)/end.Seconds())
					fmt.Fprintf(c.App.Writer, "Average time per operation: %f s\n", end.Seconds()/float64(count))

				}

				return nil

			},
		},
	}

	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "%q not implemented.\n", command)
	}

	_ = app.Run(os.Args)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
