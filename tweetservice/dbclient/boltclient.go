package dbclient

import (
	"encoding/json"
	"fmt"
	"strconv"
	logrus "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"bytes"
	"github.com/djdnl13/twitter/tweetservice/model"
	"math/rand"
)

type IBoltClient interface {
	OpenBoltDb()
	QueryTweetsOffset(accountId string) ([]model.Tweet, error)
	QueryTweets(offset string) ([]model.Tweet, error)
	AddTweet(accountId string, text string, likesCount string) (bool, error)
	Seed()
	Check() bool
}

// Real implementation
type BoltClient struct {
	boltDB *bolt.DB
}

func (bc *BoltClient) OpenBoltDb() {
	var err error
	bc.boltDB, err = bolt.Open("tweets.db", 0600, nil)
	if err != nil {
		logrus.Fatal(err)
	}
}

func (bc *BoltClient) AddTweet(accountId string, text string, likesCount string) (bool, error) {

	j := 12
	keyTweet := strconv.Itoa(j)

	// Create an instance of our Tweet struct
	tweet := model.Tweet{
		Id:   accountId + "." + keyTweet,
		Text: text,
		LikesCount: likesCount,
		AccountId: accountId,
	}

	// Serialize the struct to JSON
	jsonBytes, _ := json.Marshal(tweet)

	// Write the data to the TweetBucket
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("TweetBucket"))
		err := b.Put([]byte(accountId + "." + keyTweet), jsonBytes)
		return err
	})

	return true, nil
}

func (bc *BoltClient) QueryTweetsOffset(offset string) ([]model.Tweet, error) {
	// Allocate an empty Account instance we'll let json.Unmarhal populate for us in a bit.
	tweet := model.Tweet{}
	var tweets []model.Tweet
	var n, _  = strconv.Atoi(offset)
	// Read an object from the bucket using boltDB.View
	err := bc.boltDB.View(func(tx *bolt.Tx) error {

		c := tx.Bucket([]byte("TweetBucket")).Cursor()
		i := 0
		for _, v := c.Last(); i < n; _, v = c.Prev() {
			json.Unmarshal(v, &tweet)
			tweets = append(tweets, tweet)
			i++
		}

		// Return nil to indicate nothing went wrong, e.g no error
		return nil
	})
	// If there were an error, return the error
	if err != nil {
		return tweets, err
	}
	// Return the Account struct and nil as error.
	return tweets, nil
}
func (bc *BoltClient) QueryTweets(accountId string) ([]model.Tweet, error) {
	// Allocate an empty Account instance we'll let json.Unmarhal populate for us in a bit.
	tweet := model.Tweet{}
	var tweets []model.Tweet

	// Read an object from the bucket using boltDB.View
	err := bc.boltDB.View(func(tx *bolt.Tx) error {

		c := tx.Bucket([]byte("TweetBucket")).Cursor()

		prefix := []byte(accountId + ".")

		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			json.Unmarshal(v, &tweet)
			tweets = append(tweets, tweet)
		}

		// Return nil to indicate nothing went wrong, e.g no error
		return nil
	})
	// If there were an error, return the error
	if err != nil {
		return tweets, err
	}
	// Return the Account struct and nil as error.
	return tweets, nil
}

// Start seeding accounts
func (bc *BoltClient) Seed() {
	bc.initializeBucket()
	bc.seedTweets()
}

// Naive healthcheck, just makes sure the DB connection has been initialized.
func (bc *BoltClient) Check() bool {
	return bc.boltDB != nil
}

// Creates an "AccountBucket" in our BoltDB. It will overwrite any existing bucket of the same name.
func (bc *BoltClient) initializeBucket() {
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("TweetBucket"))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
}

// Seed (n) make-believe account objects into the AcountBucket bucket.
func (bc *BoltClient) seedTweets() {

	total := 100
	for i := 0; i < total; i++ {

		// Generate a key 10000 or larger
		keyAccount := strconv.Itoa(10000 + i)
                tweetsNumber := rand.Intn(9) + 1
		for j := 0; j < tweetsNumber ; j++ {
			keyTweet := strconv.Itoa(j)
			// Create an instance of our Account struct
			tweet := model.Tweet{
				Id:   keyAccount+"."+keyTweet,
				Text: "Tweet message " +strconv.Itoa(j),
				LikesCount: strconv.Itoa(i),
				AccountId: keyAccount,
			}

			// Serialize the struct to JSON
			jsonBytes, _ := json.Marshal(tweet)

			// Write the data to the AccountBucket
			bc.boltDB.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("TweetBucket"))
				err := b.Put([]byte(keyAccount + "." + keyTweet), jsonBytes)
				return err
			})
		}
	}
	logrus.Infof("Seeded %v fake tweets...\n", total)
}
