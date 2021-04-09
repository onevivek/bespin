package main

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"strconv"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/onevivek/bespin/go_lang/IMAP/database"
)

func main() {
	useStartTLS := false
	emailAddress := os.Getenv("IMAP_EMAIL_ID")
	password := os.Getenv("IMAP_APP_PASSWORD")
	serverName := os.Getenv("IMAP_SERVER_NAME")
	port := os.Getenv("IMAP_PORT")
	RedisAddress := os.Getenv("RedisAddress")
	noOfEmailsPerBatch := os.Getenv("noOfEmailsPerBatch")
	var noOfEmailsToFetch uint32 = 10
	if emailAddress == "" || password == "" || serverName == "" {
		panic("[IMAP_EMAIL_ID,IMAP_APP_PASSWORD,IMAP_SERVER_NAME]: Either one of these field or all the fields are missing in environment variables")
	}
	if RedisAddress == "" {
		RedisAddress = "localhost:6379"
	}

	// If the port is not mentioned in
	if port == "" {
		port = "993"
	}
	if noOfEmailsPerBatch != "" {
		if num, convErr := strconv.ParseUint(noOfEmailsPerBatch, 10, 32); convErr == nil {
			noOfEmailsToFetch = uint32(num)
		}
	}

	db, err := database.NewRedis(RedisAddress)
	if err != nil {
		log.Printf("cannot connect to redis: %v\n", err)
		os.Exit(1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("error closing db connection: %v\n", err)
		}
	}()
	ctx := context.Background()
	data, erro := db.Get(ctx, "uid")
	if erro != nil {
		log.Println("Error: ", erro)
	}
	log.Println("Data: ", data)
	hostport := serverName + ":" + port
	log.Println("Using host and port:-", hostport)
	log.Println("With User:-", emailAddress)

	log.Println("Connecting to server...")
	// Connect to server
	c, err := client.DialTLS(hostport, nil)
	//c.SetDebug(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")
	// Don't forget to logout

	defer c.Logout()

	if useStartTLS {
		ret, err := c.SupportStartTLS()
		if err != nil {
			log.Println("Error trying to determine whether StartTLS is supported.")
			log.Fatal(err)
		}
		if !ret {
			log.Println("StartTLS is not supported.")
			log.Fatal(err)
		}
		log.Println("Good, StartTLS is supported.")
		// Start a TLS session
		tlsConfig := &tls.Config{ServerName: serverName}
		if err := c.StartTLS(tlsConfig); err != nil {
			log.Fatal(err)
		}
		log.Println("TLS started")
	}

	// Login
	if err := c.Login(emailAddress, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)

		// Select mailBox
		mbox, err := c.Select(m.Name, false)
		if err != nil {
			log.Println("Mailbox selection error: ", err)
		}
		if mbox.Messages == 0 {
			log.Printf("No message in %s", m.Name)
			continue
		}
		log.Println("Flags for ", m.Name, mbox.Flags)

		// Get the last ${noOfEmailsToFetch} messages
		log.Printf("Fetching %d Message(s) from %s.", noOfEmailsToFetch, m.Name)

		uids := make([]uint32, 0)
		var lastFetchedUID = uint32(0)
		if mbox.Messages > noOfEmailsToFetch {
			// We're using unsigned integers here, only subtract if the result is > 0
			if lastFetchedUID < noOfEmailsToFetch {
				for a := lastFetchedUID; a < noOfEmailsToFetch; a++ {
					uids = append(uids, a+1)
				}
			} else if lastFetchedUID == noOfEmailsToFetch {
				uids = append(uids, lastFetchedUID+1)
			}
		} else if lastFetchedUID < mbox.Messages {
			for a := lastFetchedUID; a <= mbox.Messages; a++ {
				uids = append(uids, a+1)
			}
		}
		log.Printf("Uids: %d", uids)
		if len(uids) < 0 {
			continue
		}
		seqset := new(imap.SeqSet)
		seqset.AddNum(uids...)
		// Get the whole message RAW
		items := []imap.FetchItem{
			// imap.FetchBody,
			// imap.FetchBodyStructure,
			imap.FetchEnvelope,
			// imap.FetchFlags,
			// imap.FetchInternalDate,
			// imap.FetchRFC822,
			// imap.FetchRFC822Header,
			// imap.FetchRFC822Size,
			imap.FetchRFC822Text,
			imap.FetchUid,
		}
		messages := make(chan *imap.Message, noOfEmailsToFetch)
		done = make(chan error, 1)

		go func() {
			done <- c.UidFetch(seqset, items, messages)
		}()
		for msg := range messages {
			if msg != nil {
				log.Println("UID: ", msg.Uid)
				log.Println("To: ")
				for _, toAddrs := range msg.Envelope.To {
					log.Printf("%s", *toAddrs)
				}
				log.Println("Subject: ", msg.Envelope.Subject)
				for _, value := range msg.Body {
					if value != nil {
						len := value.Len()
						buf := make([]byte, len)
						n, err := value.Read(buf)
						if err != nil {
							log.Println("Error: ", err)
						}
						if n != len {
							log.Println("Error: Didn't read correct length")
						}
						println("+++++++++++++++++++++++Start of Raw Message +++++++++++++++++++++++++++++++")
						print(string(buf))
						println("+++++++++++++++++++++++End of Raw Message +++++++++++++++++++++++++++++++++")
					}
				}
			}
		}
		if erro := <-done; erro != nil {
			log.Println("Error in done: ", erro)
		}
	}
	log.Println("Done")
}
