package main

import "github.com/ChimeraCoder/anaconda"
import "fmt"
import "net/url"
import "strings"
import "time"
import "math/rand"
import "strconv"

var tag string = " #MAOUVS"
var deadcounter int = 0

func showTimeLine(api *anaconda.TwitterApi, v url.Values) {
  tweets, err := api.GetHomeTimeline(v)
  if err != nil {
    panic(err)
  }
  for _, tweet := range tweets {
    fmt.Println("tweet: ", tweet.Text)
    fmt.Println("id: ", tweet.Id)
  }
}

func gameGenerator(api *anaconda.TwitterApi) func(anaconda.Tweet) bool{
  var enemyHp int = 3 //HP is Hit Point
  var myHp int = 3
  var enemyTp int = 0  //Tp is Tame Point
  var myTp int = 0
  var turn int = 0
  var maouRoutines = [3][3][10]int{
  {
    {0,1,0,2,0,0,0,2,0,1},  //攻撃60%防御20%タメ20%
    {1,1,1,0,0,2,2,1,1,1},  //攻撃20%防御60%タメ20%
    {0,1,1,0,2,0,1,1,1,0},  //攻撃40%防御50%タメ10%
  },
  {
    {0,1,0,0,0,0,0,0,0,1},  //攻撃80%防御20%タメ0%
    {0,1,0,0,0,2,2,0,0,0},  //攻撃70%防御10%タメ20%
    {0,1,0,0,2,0,0,1,0,0},  //攻撃70%防御20%タメ10%
  },
  {
    {0,2,2,2,2,2,2,2,2,0},  //攻撃20%防御0%タメ80%
    {2,2,2,2,1,1,2,2,2,2},  //攻撃0%防御20%タメ80%
    {2,0,2,2,2,1,2,2,0,2},  //攻撃20%防御10%タメ70%
  },
  }
  var commands =[4]string{"こうげき","ぼうぎょ","ため","ひっさつ"}
  var lossRep = " あなたの負けです"
  var winRep = " あなたの勝ちです"

  maouRoutine := maouRoutines[ deadcounter%len(maouRoutines) ]

  //げーむるーちん
  return func (tweet anaconda.Tweet) bool{
    //fmt.Println("tw",tweet.Text)
    //fmt.Println("turn",)
    var act = -1

    for i := 0; i < len(commands); i++{
      if strings.Contains(tweet.Text,commands[i]){
        act = i
        break
      }
    }
    if(act < 0){return false}

    turn++
    turn %= len(maouRoutine)
    rand.Seed(time.Now().UnixNano())
    routine := rand.Intn(len(maouRoutine[0]))
    mact := maouRoutine[turn][routine]
    if enemyTp == 3{
      mact = 3
    }

    sendmess := "@"+tweet.User.ScreenName+" "
    sendmess += "勇者は" + commands[act] +"をした\n"
    sendmess += "魔王は" + commands[mact] + "をした\n"

    //以下、初心者あともすふぇあ満載のif祭り
    //こうげき
    if mact == 0 && act != 1{
      sendmess += "勇者はダメージを受けた\n"
      myHp--
    }
    if mact != 1 && act == 0{
      sendmess += "魔王はダメージを受けた\n"
      enemyHp--
    }

    //ため
    if act == 2{
      myTp++
      sendmess += "勇者タメ"+ strconv.Itoa(myTp)+"\n"
    }
    if mact == 2{
      enemyTp++
      sendmess += "魔王タメ"+ strconv.Itoa(enemyTp)+"\n"
    }

    //ひっさつ
    if mact == 3{
        sendmess += "勇者は必殺を受けた\n"
        if act == 1{
          myHp++
        }
        myHp -= 3
        enemyTp = 0
    }
    if act == 3{
        if myTp < 3{
          sendmess += "勇者はためが足りなかった\n"
        }else{
          sendmess += "魔王は必殺を受けた\n"
          if mact == 1{
            enemyHp++
          }
          enemyHp -= 3
        }
        myTp = 0
    }

    if enemyHp <= 0{
      sendmess += winRep
      deadcounter++
    }else if myHp <= 0{
      sendmess += lossRep
    }

    api.PostTweet(sendmess + "【"+tweet.IdStr+"】"+tag,nil)

    return enemyHp <= 0 || myHp <= 0
  }
}

func checkerCloser() func(anaconda.Tweet,*anaconda.TwitterApi){
  games := make(map[int64] func(anaconda.Tweet)bool)
  var words = [2]string{"ゲームスタート","@golang_bot"}
  var repray = " ゲームを始めます。やり方はbioを見てください。"

  return func (tweet anaconda.Tweet,api *anaconda.TwitterApi) {
      if strings.Contains(tweet.Text,words[1]){
        game,ok := games[tweet.User.Id]
        if strings.Contains(tweet.Text,words[0]){
          games[tweet.User.Id] = gameGenerator(api)
          api.PostTweet("@"+tweet.User.ScreenName+
            repray+"【"+tweet.IdStr+"】"+tag,nil)
        }else if ok{
          if game(tweet){
              fmt.Println("DeleteGame")
              delete(games,tweet.User.Id)
          }
        }
      }
    }
}
func main() {
	anaconda.SetConsumerKey("j62kyFWBVqCiJElWSHvofSz59")
	anaconda.SetConsumerSecret("FmA9RF7MAEwLbvXfvYRgZMMstinToo6kbR7CcwfpdhWrsEPyrg")
	api := anaconda.NewTwitterApi("3614143880-hYmuWlZc2UrXjC6lu18YUJ4xEPi60IFcZNrHZui","oEyx47xfEw1MFPAHmqpDqL5PCVtAYumTDCgy35wIFD2yv")

   v := url.Values{}
   v.Add("replies", "all")
   stream := api.UserStream(v)

   check := checkerCloser()
   for {
     select {
     case item := <-stream.C:
       switch status := item.(type) {
       case anaconda.Tweet:
         fmt.Printf("%v: %v\n%v\n", status.User.ScreenName, status.Text, status.InReplyToScreenName)
         check(status,api)
       default:
       }
     }
   }
}
